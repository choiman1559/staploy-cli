package registry

import (
	"staploy-cli/app/cmds"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
)

type RegistryAddRepoTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.RegistryAddRepoCmd]
}

func (task *RegistryAddRepoTask) MainCmd() error {
	requestPacket := task.CreateDefPacket()
	registryRequest := &proto.RegistryRequestPacket{
		TaskType:      proto.TaskRegistryTypes_LOCAL_ADD_REPOSITORY,
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
