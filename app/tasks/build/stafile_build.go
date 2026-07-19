package build

import (
	"os/exec"
	"staploy-cli/app/cmds"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
	"strings"
)

type ArchTarget struct {
	name     string
	target   *BuildTarget
	setCmdFn func(path string)
}

func (a *StaFileTask) processBuild(builds []*Build) error {
	for _, build := range builds {
		logger.Info("Processing build \"%s\"", build.AppName)
		logger.EnableTree()
		postRun := make(map[string]string)

		if build.PreBuild != "" {
			out, err := execShell(build.PreBuild)
			if err != nil {
				logger.DisableTree(true)
				logger.Error("Error executing global pre-build command, Abort build => %v", err)
				return nil
			}
			logger.Info("Executed global pre-build command => %s", strings.TrimSuffix(out, "\n"))
		}

		buildCmd := cmds.BuildCmd{
			AppName:     build.AppName,
			VersionName: build.Version,
			LibVersion:  build.LibVersion,
			OutputDir:   build.OutputDir,
			Executable:  build.Executables,
		}

		arches := []ArchTarget{
			{"share", build.Target_share, func(p string) { buildCmd.Share = p }},
			{"i386", build.Target_i386, func(p string) { buildCmd.I386 = p }},
			{"x86_64", build.Target_x86_64, func(p string) { buildCmd.X86_64 = p }},
			{"arm", build.Target_arm, func(p string) { buildCmd.Arm = p }},
			{"aarch64", build.Target_aarch64, func(p string) { buildCmd.Aarch64 = p }},
			{"riscv32", build.Target_riscv32, func(p string) { buildCmd.Riscv32 = p }},
			{"riscv64", build.Target_riscv64, func(p string) { buildCmd.Riscv64 = p }},
			{"mipsel", build.Target_mipsel, func(p string) { buildCmd.Mipsel = p }},
			{"mips64el", build.Target_mips64el, func(p string) { buildCmd.Mips64el = p }},
		}

		for _, arch := range arches {
			processArch(arch, postRun)
		}

		t := &PkgCmdTask{}
		t.Init(a.DefaultArgs, buildCmd, proto.TaskGroup_TASK_MANAGE_APPS)
		err := t.MainCmd()

		if err != nil {
			logger.Error("Error processing build \"%s\": %v", build.AppName, err)
		} else {
			for arch, postCmd := range postRun {
				out, err := execShell(postCmd)
				if err != nil {
					logger.Error("Error executing %s post-build command => %s: %v", arch, postCmd, err)
				}
				logger.Info("Executed %s post-build command => %s", arch, strings.TrimSuffix(out, "\n"))
			}

			out, err := execShell(build.PostBuild)
			if err != nil {
				logger.Error("Error executing global post-build command => %s: %v", build.PostBuild, err)
			}
			logger.Info("Executed global post-build command => %s", strings.TrimSuffix(out, "\n"))
		}

		logger.DisableTree(true)
		logger.Tip("Finished building \"%s\"", build.AppName)
	}
	return nil
}

func processArch(arch ArchTarget, postRun map[string]string) {
	if arch.target == nil {
		return
	}

	if arch.target.PreBuild != "" {
		out, err := execShell(arch.target.PreBuild)
		if err != nil {
			logger.Error("Error executing %s pre-build command, skipping => %v", arch.name, err)
			return
		}
		logger.Info("Executed %s pre-build => %s", arch.name, strings.TrimSuffix(out, "\n"))
	}

	arch.setCmdFn(arch.target.Path)
	if arch.target.PostBuild != "" {
		postRun[arch.name] = arch.target.PostBuild
	}
}

func execShell(command string) (string, error) {
	out, err := exec.Command("bash", "-c", command).Output()
	return string(out), err
}
