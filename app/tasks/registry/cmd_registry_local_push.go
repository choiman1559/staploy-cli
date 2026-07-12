package registry

import (
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
)

type RegistryPushLocalTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.RegistryPushLocalCmd]
}

func (task *RegistryPushLocalTask) MainCmd() error {
	task.OverrideConnType(consts.ConnTypeRegistry)
	logger.Process("Uploading file " + task.CmdArgs.PackageFile)
	blobToken, err := task.UploadFile(task.CmdArgs.PackageFile)
	if err != nil {
		return err
	}

	requestPacket := task.CreateDefPacket()
	registryRequest := &proto.RegistryRequestPacket{
		TaskType: proto.TaskRegistryTypes_TASK_PUSH,
		BlobId:   &blobToken,
	}

	requestPacket.TaskType = &proto.RequestPacket_RegistryTaskType{RegistryTaskType: registryRequest}
	response, err := task.PostRequest(requestPacket)
	if err != nil {
		return err
	}

	logger.Info("Registry List Task Response: %v", response)
	return nil
}
