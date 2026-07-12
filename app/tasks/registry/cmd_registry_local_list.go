package registry

import (
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
)

type RegistryListLocalTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.RegistryListLocalCmd]
}

func (task *RegistryListLocalTask) MainCmd() error {
	task.OverrideConnType(consts.ConnTypeRegistry)
	requestPacket := task.CreateDefPacket()

	registryRequest := &proto.RegistryRequestPacket{
		TaskType: proto.TaskRegistryTypes_TASK_LIST,
	}

	if task.CmdArgs.AppName != "" {
		appInfo := &proto.AppInfoFetch{
			App: &proto.AppInfo{
				AppName: task.CmdArgs.AppName,
			},
		}
		registryRequest.AppInfo = appInfo
	}

	requestPacket.TaskType = &proto.RequestPacket_RegistryTaskType{RegistryTaskType: registryRequest}
	response, err := task.PostRequest(requestPacket)
	if err != nil {
		return err
	}

	logger.Info("Registry List Task Response: %v", response)
	return nil
}
