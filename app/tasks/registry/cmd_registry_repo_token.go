package registry

import (
	"staploy-cli/app/cmds"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
)

type RegistryManageRepoTokenTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.RegistryManageRepoTokenCmd]
}

func (task *RegistryManageRepoTokenTask) MainCmd() error {
	requestPacket := task.CreateDefPacket()
	registryRequest := &proto.RegistryRequestPacket{
		TaskType:      proto.TaskRegistryTypes_LOCAL_MANAGE_TOKEN_REPOSITORY,
		RepositoryUrl: append(make([]string, 0), task.CmdArgs.RepoUrl),
		BlobId:        &task.CmdArgs.Token,
	}

	requestPacket.TaskType = &proto.RequestPacket_RegistryTaskType{RegistryTaskType: registryRequest}
	response, err := task.PostRequest(requestPacket)
	if err != nil {
		return err
	}

	logger.Info("Registry List Task Response: %v", response)
	return nil
}
