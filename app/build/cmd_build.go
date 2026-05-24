package build

import (
	"archive/tar"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
	"staploy-cli/app/proto"
	"strings"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
)

type PkgCmdTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.BuildCmd]
}

var baseAppInfo *proto.InstalledAppInfo
var TargetPathByArch map[proto.CpuArch]string
var TargetHashVersion map[proto.CpuArch]*proto.Version
var ShareExecHash map[string]bool

func (a *PkgCmdTask) MainCmd() error {
	TargetPathByArch = make(map[proto.CpuArch]string)
	TargetHashVersion = make(map[proto.CpuArch]*proto.Version)
	ShareExecHash = make(map[string]bool)

	err := a.preCheckDirExists()
	if err != nil {
		return err
	}

	baseAppInfo = &proto.InstalledAppInfo{
		App: &proto.AppInfo{
			AppName:        a.CmdArgs.AppName,
			AppDescription: &a.CmdArgs.AppDescription,
		},
		CurrentVersion: &proto.Version{
			VersionName: a.CmdArgs.VersionName,
			LibVersion:  &a.CmdArgs.LibVersion,
		},
	}

	err = a.startPacking()
	if err != nil {
		return err
	}
	return nil
}

func (a *PkgCmdTask) getArchivePath() string {
	fileName := fmt.Sprintf("%s_%s.tar", a.CmdArgs.AppName, a.CmdArgs.VersionName)
	return filepath.Join(a.CmdArgs.OutputDir, fileName)
}

func (a *PkgCmdTask) startPacking() error {
	tarFile, err := os.Create(a.getArchivePath())
	if err != nil {
		log.Fatalf("failed to create tar archive: %v", err)
	}

	defer func(tarFile *os.File) {
		err := tarFile.Close()
		if err != nil {
			return
		}
	}(tarFile)

	tw := tar.NewWriter(tarFile)
	defer func(tw *tar.Writer) {
		err := a.writePkgHeader(tw)
		if err != nil {
			return
		}

		err = tw.Close()
		if err != nil {
			return
		}
	}(tw)

	for arch := range TargetPathByArch {
		if err := a.calculateVerHash(arch); err != nil {
			return err
		}

		if err := a.addArchDir(arch, tw); err != nil {
			return err
		}
	}
	return nil
}

func (a *PkgCmdTask) writePkgHeader(tw *tar.Writer) error {
	header := &proto.PackageHeader{
		FormatVersion: consts.PACKAGE_FORMAT_VERSION,
		ShareOnly:     false,
		PackageInfo:   baseAppInfo,
	}

	content, err := protojson.Marshal(header)
	if err != nil {
		return err
	}

	tarHeader := &tar.Header{
		Name:     consts.PACKAGE_FILE_METADATA,
		Size:     int64(len(content)),
		Mode:     0644,
		Typeflag: tar.TypeReg,
		ModTime:  time.Now(),
	}

	if err := tw.WriteHeader(tarHeader); err != nil {
		return err
	}

	reader := strings.NewReader(string(content))
	if _, err := io.Copy(tw, reader); err != nil {
		return err
	}
	return nil
}

func (a *PkgCmdTask) calculateVerHash(arch proto.CpuArch) error {
	srcPath := TargetPathByArch[arch]
	var binaries []*proto.BinaryInfo

	for _, execFileName := range a.CmdArgs.Executable {
		hasSharedExec, ok := ShareExecHash[execFileName]
		if arch == proto.CpuArch_UNKNOWN || !(ok && hasSharedExec) {
			fileName := filepath.Join(srcPath, execFileName)
			fileHash, err := getSha1Hash(fileName)

			if err != nil {
				if arch == proto.CpuArch_UNKNOWN && os.IsNotExist(err) {
					ShareExecHash[execFileName] = false
					continue
				}
				return err
			}

			binaries = append(binaries, &proto.BinaryInfo{
				Name: execFileName,
				Hash: fileHash,
			})

			if arch == proto.CpuArch_UNKNOWN {
				ShareExecHash[execFileName] = true
			}
		}
	}

	TargetHashVersion[arch] = &proto.Version{
		EntryBinaries: binaries,
	}
	return nil
}

func getSha1Hash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)

	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	return hashString, nil
}

func (a *PkgCmdTask) addArchDir(arch proto.CpuArch, tw *tar.Writer) error {
	srcPath := TargetPathByArch[arch]
	version := TargetHashVersion[arch]

	tarDirPath := arch.String()
	if //goland:noinspection GoDfaConstantCondition
	arch == proto.CpuArch_UNKNOWN {
		tarDirPath = consts.PACKAGE_DIR_SHARE
	}

	err := filepath.WalkDir(srcPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(srcPath, path)
		if err != nil {
			return err
		}

		header.Name = filepath.ToSlash(filepath.Join(tarDirPath, relPath))

		if d.IsDir() {
			header.Name += "/"
		}

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if !d.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer func(file *os.File) {
				err := file.Close()
				if err != nil {
					return
				}
			}(file)

			if _, err := io.Copy(tw, file); err != nil {
				return err
			}
		}
		return nil
	})

	content, err := protojson.Marshal(version)
	if err != nil {
		return err
	}

	header := &tar.Header{
		Name:     fmt.Sprintf("%s/%s", tarDirPath, consts.PACKAGE_FILE_METADATA),
		Size:     int64(len(content)),
		Mode:     0644,
		Typeflag: tar.TypeReg,
		ModTime:  time.Now(),
	}

	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	reader := strings.NewReader(string(content))
	if _, err := io.Copy(tw, reader); err != nil {
		return err
	}

	return err
}

func (a *PkgCmdTask) preCheckDirExists() error {
	result, err := checkDirExists(a.CmdArgs.OutputDir)
	if err != nil || !result {
		return fmt.Errorf("output dir doesn't exists: %s", a.CmdArgs.OutputDir)
	}

	oneOfTargetSpecified := false
	checkTask := func(path string, arch proto.CpuArch) error {
		result, err := checkDirExists(path)
		if err != nil {
			return err
		} else if result {
			oneOfTargetSpecified = true
			TargetPathByArch[arch] = path
		}
		return nil
	}

	err = checkTask(a.CmdArgs.I386, proto.CpuArch_i386)
	if err != nil {
		return err
	}

	err = checkTask(a.CmdArgs.X86_64, proto.CpuArch_x86_64)
	if err != nil {
		return err
	}

	err = checkTask(a.CmdArgs.Arm, proto.CpuArch_arm)
	if err != nil {
		return err
	}

	err = checkTask(a.CmdArgs.Aarch64, proto.CpuArch_aarch64)
	if err != nil {
		return err
	}

	err = checkTask(a.CmdArgs.Riscv32, proto.CpuArch_riscv32)
	if err != nil {
		return err
	}

	err = checkTask(a.CmdArgs.Riscv64, proto.CpuArch_riscv64)
	if err != nil {
		return err
	}

	err = checkTask(a.CmdArgs.Mipsel, proto.CpuArch_mipsel)
	if err != nil {
		return err
	}

	err = checkTask(a.CmdArgs.Mips64el, proto.CpuArch_mips64el)
	if err != nil {
		return err
	}

	result, err = checkDirExists(a.CmdArgs.Share)
	if err != nil {
		return err
	} else if result {
		TargetPathByArch[proto.CpuArch_UNKNOWN] = a.CmdArgs.Share
	}

	if !oneOfTargetSpecified {
		return fmt.Errorf("no build target specified")
	}
	return nil
}

func checkDirExists(path string) (bool, error) {
	if path == "" {
		return false, nil
	}

	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, fmt.Errorf("%s does not exist", path)
	}
	if !stat.IsDir() {
		return false, fmt.Errorf("%s is not a directory", path)
	}
	return true, nil
}
