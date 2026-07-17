package registry

import (
	"errors"
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
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

	if response.GetStatus() != consts.StatusOK {
		if response.GetErrorCause() != "" {
			return errors.New(response.GetErrorCause())
		}
		return errors.New("registry add repository failed")
	}

	if len(response.GetRegistryResponse().RepositoryUrl) > 0 {
		logger.Info("Registry url successfully added: %v", response.GetRegistryResponse().RepositoryUrl)
	} else {
		return errors.New("requested repository is already exist on server")
	}
	return nil
}
