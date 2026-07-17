package build

import (
	"staploy-cli/app/cmds"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
	"staploy-cli/app/tasks/apps"
	"staploy-cli/app/tasks/registry"
)

func (a *StaFileTask) processManage(defArgs *cmds.DefaultArgs, manages []*Manage) error {
	for _, manage := range manages {
		logger.Info("Processing manage app_name \"%s\"", manage.AppName)
		logger.EnableTree()

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
			t := &apps.UploadCmdTask{}
			t.Init(*defArgs, cmds.UploadCmd{PackageFile: manage.Upload.PackageFile}, proto.TaskGroup_TASK_MANAGE_APPS)

			err := t.MainCmd()
			if err != nil {
				logger.Error("Error uploading app \"%s\": %v", manage.AppName, err)
			}
		}

		if manage.Delete != nil {
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
		logger.Tip("Finished manage app_name \"%s\"", manage.AppName)
	}
	return nil
}
