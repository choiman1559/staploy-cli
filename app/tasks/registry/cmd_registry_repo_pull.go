package registry

import (
	"errors"
	"fmt"
	"staploy-cli/app/cmds"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
)

type RegistryPullTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.RegistryPullCmd]
}

func (task *RegistryPullTask) MainCmd() error {
	requestPacket := task.CreateDefPacket()
	if task.CmdArgs.AppName == "" {
		return fmt.Errorf("app name is required")
	}

	appInfo := &proto.AppInfoFetch{
		App: &proto.AppInfo{
			AppName: task.CmdArgs.AppName,
		},
	}

	if task.CmdTask.CmdArgs.Version != "" {
		appInfo.AppVersion = []*proto.Version{{VersionName: task.CmdTask.CmdArgs.Version}}
	}

	registryRequest := &proto.RegistryRequestPacket{
		TaskType: proto.TaskRegistryTypes_LOCAL_PULL_PACKAGE,
		AppInfo:  appInfo,
	}

	if task.CmdTask.CmdArgs.Repository != "" {
		registryRequest.RepositoryUrl = []string{task.CmdTask.CmdArgs.Repository}
	}

	requestPacket.TaskType = &proto.RequestPacket_RegistryTaskType{RegistryTaskType: registryRequest}
	response, err := task.PostRequest(requestPacket)
	if err != nil {
		return err
	}

	if len(response.GetRegistryResponse().GetAppInfo()) < 1 {
		if response.GetErrorCause() != "" {
			return errors.New(response.GetErrorCause())
		}
		return errors.New("could not find app info: " + task.CmdTask.CmdArgs.AppName)
	}

	packageMetadata := response.GetRegistryResponse().GetAppInfo()[0]
	packagePrinter := fmt.Sprintf("Package pulled: \"%s\" version %s from %s", packageMetadata.GetApp().AppName, packageMetadata.GetCurrentVersion().GetVersionName(), response.RegistryResponse.RepositoryUrl[0])

	logger.Info("%s", packagePrinter)
	return nil
}
