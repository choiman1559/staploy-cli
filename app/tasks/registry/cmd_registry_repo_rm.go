package registry

import (
	"errors"
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
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

	if response.GetStatus() != consts.StatusOK {
		if response.GetErrorCause() != "" {
			return errors.New(response.GetErrorCause())
		}
		return errors.New("registry remove repository failed")
	}

	if len(response.GetRegistryResponse().RepositoryUrl) > 0 {
		logger.Info("Registry url successfully removed: %v", response.GetRegistryResponse().RepositoryUrl)
	} else {
		return errors.New("requested repository is not exists")
	}
	return nil
}
