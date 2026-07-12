package registry

import (
	"staploy-cli/app/cmds"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
)

type RegistryRemoveRepoTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.RegistryRemoveRepoCmd]
}

func (task *RegistryRemoveRepoTask) MainCmd() error {
	requestPacket := task.CreateDefPacket()
	registryRequest := &proto.RegistryRequestPacket{
		TaskType:      proto.TaskRegistryTypes_LOCAL_REMOVE_REPOSITORY,
		RepositoryUrl: task.CmdArgs.RepoUrl,
	}

	requestPacket.TaskType = &proto.RequestPacket_RegistryTaskType{RegistryTaskType: registryRequest}
	response, err := task.PostRequest(requestPacket)
	if err != nil {
		return err
	}

	logger.Info("Registry List Task Response: %v", response)
	return nil
}
