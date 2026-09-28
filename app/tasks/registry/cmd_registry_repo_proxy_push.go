package registry

import (
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
	"strconv"
)

type RegistryPushProxyTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.RegistryPushLocalCmd]

	ProxyUrl         string
	ProxyDefaultArgs cmds.DefaultArgs
}

func (task *RegistryPushProxyTask) MainCmd() error {
	logger.Process("Uploading file %s on server %s:%d", task.CmdArgs.PackageFile, task.ProxyDefaultArgs.Address, task.ProxyDefaultArgs.Port)
	task.OverrideConnType(consts.ConnTypeRegistry)
	defaultArgs := task.DefaultArgs
	task.DefaultArgs = task.ProxyDefaultArgs

	blobToken, err := task.UploadFile(task.CmdArgs.PackageFile)
	if err != nil {
		return err
	}

	task.OverrideConnType(consts.ConnTypeAdmin)
	task.DefaultArgs = defaultArgs

	requestPacket := task.CreateDefPacket()
	registryRequest := &proto.RegistryRequestPacket{
		TaskType:      proto.TaskRegistryTypes_LOCAL_REMOTE_PROXY,
		BlobId:        &blobToken,
		RepositoryUrl: []string{task.ProxyUrl, strconv.Itoa(int(proto.TaskRegistryTypes_TASK_PUSH))},
	}

	requestPacket.TaskType = &proto.RequestPacket_RegistryTaskType{RegistryTaskType: registryRequest}
	response, err := task.PostRequest(requestPacket)
	if err != nil {
		return err
	}

	err = PrintRegistryPushResult(response)
	if err != nil {
		return err
	}
	return nil
}
