package registry

import (
	"staploy-cli/app/cmds"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
)

type RegistryListRepoTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.RegistryListRepoCmd]
}

func (task *RegistryListRepoTask) MainCmd() error {
	requestPacket := task.CreateDefPacket()
	registryRequest := &proto.RegistryRequestPacket{
		TaskType: proto.TaskRegistryTypes_LOCAL_LIST_REPOSITORY,
	}

	requestPacket.TaskType = &proto.RequestPacket_RegistryTaskType{RegistryTaskType: registryRequest}
	response, err := task.PostRequest(requestPacket)
	if err != nil {
		return err
	}

	logger.Info("Registry List Task Response: %v", response)
	return nil
}
