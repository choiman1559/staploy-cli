package registry

import (
	"errors"
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
	"strings"
)

type RegistryUpdateCacheTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.RegistryUpdateCacheCmd]
}

func (task *RegistryUpdateCacheTask) MainCmd() error {
	requestPacket := task.CreateDefPacket()
	registryRequest := &proto.RegistryRequestPacket{
		TaskType:      proto.TaskRegistryTypes_LOCAL_PACKAGE_CACHE_UPDATE,
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
		return errors.New("update repository package cache failed")
	}

	for _, repoUrl := range response.GetRegistryResponse().RepositoryUrl {
		data := strings.Split(repoUrl, "$")
		if data[1] != consts.StatusOK {
			logger.Error("Fetch failed: %s (cause: %s)", data[0], data[1])
			continue
		}
		logger.Info("Fetch ok: %s", data[0])
	}

	logger.Info("Successfully finished update repository cache")
	return nil
}
