package registry

import (
	"errors"
	"fmt"
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
)

type RegistryRemoveLocalTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.RegistryRemoveLocalCmd]
}

func (task *RegistryRemoveLocalTask) MainCmd() error {
	task.OverrideConnType(consts.ConnTypeRegistry)
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
		TaskType: proto.TaskRegistryTypes_TASK_REMOVE,
		AppInfo:  appInfo,
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
	logger.Tip("Package removed: \"%s\" version %s", packageMetadata.GetApp().AppName, packageMetadata.GetCurrentVersion().GetVersionName())
	return nil
}
