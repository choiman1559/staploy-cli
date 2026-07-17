package registry

import (
	"errors"
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
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

	if response.GetStatus() != consts.StatusOK {
		if response.GetErrorCause() != "" {
			return errors.New(response.GetErrorCause())
		}
		return errors.New("registry associate token to repository failed")
	}

	logger.Info("Registry token successfully associated: %S", task.CmdArgs.RepoUrl)
	return nil
}
