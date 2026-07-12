package registry

import (
	"staploy-cli/app/cmds"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
)

type RegistryListTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.RegistryListPackageCmd]
}

func (task *RegistryListTask) MainCmd() error {
	requestPacket := task.CreateDefPacket()

	registryRequest := &proto.RegistryRequestPacket{
		TaskType: proto.TaskRegistryTypes_LOCAL_PACKAGE_QUERY,
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
