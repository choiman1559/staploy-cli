package build

import (
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
	"staploy-cli/app/tasks/apps"
	"staploy-cli/app/tasks/registry"
	"strings"
)

func (a *StaFileTask) processManage(defArgs *cmds.DefaultArgs, manages []*Manage) error {
	for _, manage := range manages {
		logger.Task("Processing manage app_name \"%s\"", manage.AppName)
		logger.EnableTree()

		var resolvedAppAlias *ResolvedAppAlias
		if strings.HasPrefix(manage.AppName, consts.STAFILE_ALIAS_PREFIX) {
			aliasName := strings.TrimPrefix(manage.AppName, consts.STAFILE_ALIAS_PREFIX)
			foundAlias, err := a.hitAppAlias(aliasName, nil)
			if err != nil {
				logger.DisableTree(true)
				logger.Error("Cannot find alias for manage \"%s\"", manage.AppName)
				return nil
			}

			resolvedAppAlias = foundAlias
			manage.AppName = foundAlias.Alias.AppName
		}

		if manage.Create != nil {
			t := &apps.CreateCmdTask{}
			t.Init(*defArgs, cmds.CreateCmd{
				AppName:        manage.AppName,
				AppDescription: manage.Create.AppDescription,
			}, proto.TaskGroup_TASK_MANAGE_APPS)

			err := t.MainCmd()
			if err != nil {
				logger.Error("Error processing manage name \"%s\": %v", manage.AppName, err)
			}
		}

		if manage.Upload != nil {
			if resolvedAppAlias != nil {
				if resolvedAppAlias.Build == nil {
					logger.Error("No build output associated with app alias \"%s\"", resolvedAppAlias.Alias)
				} else {
					manage.Upload.PackageFile = resolvedAppAlias.Build.OutputDir
				}
			}

			if manage.Upload.PackageFile == "" {
				logger.Error("No package file specified for upload of app \"%s\"", manage.AppName)
			} else {
				t := &apps.UploadCmdTask{}
				t.Init(*defArgs, cmds.UploadCmd{PackageFile: manage.Upload.PackageFile}, proto.TaskGroup_TASK_MANAGE_APPS)

				err := t.MainCmd()
				if err != nil {
					logger.Error("Error uploading app \"%s\": %v", manage.AppName, err)
				}
			}
		}

		if manage.Delete != nil {
			if resolvedAppAlias != nil {
				manage.Delete.Versions = append(manage.Delete.Versions, resolvedAppAlias.Alias.Version)
			}

			t := &apps.DeleteCmdTask{}
			t.Init(*defArgs, cmds.DeleteCmd{
				AppName:     manage.AppName,
				VersionName: manage.Delete.Versions,
			}, proto.TaskGroup_TASK_MANAGE_APPS)

			err := t.MainCmd()
			if err != nil {
				logger.Error("Error deleting app \"%s\": %v", manage.AppName, err)
			}
		}

		if manage.Pull != nil {
			if resolvedAppAlias != nil {
				manage.Pull.AppName = resolvedAppAlias.Alias.AppName
				manage.Pull.Version = resolvedAppAlias.Alias.Version
			}

			t := &registry.RegistryPullTask{}
			t.Init(*defArgs, cmds.RegistryPullCmd{
				AppName:    manage.Pull.AppName,
				Version:    manage.Pull.Version,
				Repository: manage.Pull.Repository,
			}, proto.TaskGroup_TASK_REGISTRY)

			err := t.MainCmd()
			if err != nil {
				logger.Error("Error pulling app \"%s\": %v", manage.AppName, err)
			}
		}

		logger.DisableTree(true)
		logger.Task("Finished manage app_name \"%s\"", manage.AppName)
	}
	return nil
}
