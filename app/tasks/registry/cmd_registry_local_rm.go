package registry

import (
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

	logger.Info("Registry List Task Response: %v", response)
	return nil
}
