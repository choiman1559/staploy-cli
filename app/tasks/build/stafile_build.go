package build

import (
	"fmt"
	"os"
	"os/exec"
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
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
		logger.Task("Processing build \"%s\"", build.AppName)
		logger.EnableTree()
		postRun := make(map[string]string)
		var resolvedAliasName string

		if strings.HasPrefix(build.AppName, consts.STAFILE_ALIAS_PREFIX) {
			aliasName := strings.TrimPrefix(build.AppName, consts.STAFILE_ALIAS_PREFIX)
			foundAlias, err := a.hitAppAlias(aliasName, nil)
			if err != nil {
				logger.DisableTree(true)
				logger.Error("Cannot find alias for build \"%s\"", build.AppName)
				return nil
			}

			resolvedAliasName = aliasName
			build.AppName = foundAlias.Alias.AppName
			if foundAlias.Alias.Version != "" {
				build.Version = foundAlias.Alias.Version
			}
		}

		if build.Version == "" {
			logger.DisableTree(true)
			logger.Error("Version not specified, Abort build")
			return nil
		}

		outPutDir, err := a.parseExecArgs(build, build.OutputDir)
		if err != nil {
			logger.DisableTree(true)
			logger.Error("Error parsing output_dir, Abort build => %v", err)
			return nil
		}

		version, err := a.parseExecArgs(build, build.Version)
		if err != nil {
			logger.DisableTree(true)
			logger.Error("Error parsing version, Abort build => %v", err)
			return nil
		}

		libVersion, err := a.parseExecArgs(build, build.LibVersion)
		if err != nil {
			logger.DisableTree(true)
			logger.Error("Error parsing lib_version, Abort build => %v", err)
			return nil
		}

		if build.Environments != nil && len(*build.Environments) > 0 {
			var newEnv []string
			oldEnv := *build.Environments
			build.Environments = nil

			for _, env := range oldEnv {
				value := strings.SplitN(env, "=", 2)
				if len(value) != 2 || value[1] == "" {
					logger.Warn("Invalid environment format, ignoring => %s", value[0])
					continue
				}

				result, err := a.execShell(build, fmt.Sprintf("echo \"%s\"", value[1]))
				if err == nil {
					newEnv = append(newEnv, fmt.Sprintf("%s=%s", value[0], result))
				} else {
					logger.Warn("Invalid environment while parsing, ignoring => %s; %v", value[0], err)
				}
			}

			if a.DefaultArgs.Verbose {
				logger.Tip("[DEBUG] Detected environments: %+v", newEnv)
			}
			build.Environments = &newEnv
		}

		if build.PreBuild != "" {
			out, err := a.execShell(build, build.PreBuild)
			if err != nil {
				logger.DisableTree(true)
				logger.Error("Error executing global pre-build command, Abort build => %v", err)
				return nil
			}
			logger.Info("Executed global pre-build command => %s", out)
		}

		buildCmd := cmds.BuildCmd{
			AppName:     build.AppName,
			VersionName: version,
			LibVersion:  libVersion,
			OutputDir:   outPutDir,
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
			a.processArch(build, arch, postRun)
		}

		t := &PkgCmdTask{
			StaFileAliasName: resolvedAliasName,
			StaFileTask:      a,
		}

		t.Init(a.DefaultArgs, buildCmd, proto.TaskGroup_TASK_MANAGE_APPS)
		err = t.MainCmd()

		if err != nil {
			logger.Error("Error processing build \"%s\": %v", build.AppName, err)
		} else {
			for arch, postCmd := range postRun {
				out, err := a.execShell(build, postCmd)
				if err != nil {
					logger.Error("Error executing %s post-build command => %s: %v", arch, postCmd, err)
				}
				logger.Info("Executed %s post-build command => %s", arch, out)
			}

			if build.PostBuild != "" {
				out, err := a.execShell(build, build.PostBuild)
				if err != nil {
					logger.Error("Error executing global post-build command => %s: %v", build.PostBuild, err)
				}
				logger.Info("Executed global post-build command => %s", out)
			}
		}

		logger.DisableTree(true)
		logger.Task("Finished building \"%s\"", build.AppName)
	}
	return nil
}

func (a *StaFileTask) processArch(build *Build, arch ArchTarget, postRun map[string]string) {
	if arch.target == nil {
		return
	}

	parsedPath, err := a.parseExecArgs(build, arch.target.Path)
	if err != nil {
		logger.Error("Error parsing path \"%s\": %v", arch.target.Path, err)
		return
	}

	if arch.target.PreBuild != "" {
		out, err := a.execShell(build, arch.target.PreBuild)
		if err != nil {
			logger.Error("Error executing %s pre-build command, skipping => %v", arch.name, err)
			return
		}
		logger.Info("Executed %s pre-build => %s", arch.name, out)
	}

	arch.setCmdFn(parsedPath)
	if arch.target.PostBuild != "" {
		postRun[arch.name] = arch.target.PostBuild
	}
}

func (a *StaFileTask) execShell(build *Build, command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	if build.Environments != nil && len(*build.Environments) > 0 {
		cmd.Env = append(os.Environ(), *build.Environments...)
	}

	out, err := cmd.Output()
	return strings.TrimSuffix(string(out), "\n"), err
}

func (a *StaFileTask) parseExecArgs(build *Build, args string) (string, error) {
	if args == "" || !strings.HasPrefix(args, consts.STAFILE_SHELL_PREFIX) {
		return args, nil
	}
	return a.execShell(build, strings.TrimPrefix(args, consts.STAFILE_SHELL_PREFIX))
}
