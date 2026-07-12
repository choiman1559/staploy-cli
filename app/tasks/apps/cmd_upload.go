package apps

import (
	"errors"
	"fmt"
	"staploy-cli/app/cmds"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
	"strings"
)

type UploadCmdTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.UploadCmd]
}

func (task *UploadCmdTask) MainCmd() error {
	logger.Process("Uploading file " + task.CmdArgs.PackageFile)
	blobToken, err := task.UploadFile(task.CmdArgs.PackageFile)
	if err != nil {
		return err
	}

	requestPacket := task.CreateDefPacket()
	requestPacket.TaskType = &proto.RequestPacket_AppsTaskType{AppsTaskType: proto.TaskAppsTypes_TYPE_APP_PKG_PARSE}
	requestPacket.ExtraData = &blobToken

	response, err := task.PostRequest(requestPacket)
	if err != nil {
		return err
	}

	if len(response.WorkerResponse) < 1 {
		if response.ErrorCause != "" {
			return errors.New(response.ErrorCause)
		}
		return fmt.Errorf("package is corrupted or missing metadata")
	}

	packageMetadata := response.GetWorkerResponse()[0].GetWorkerInfo().GetInstalledApp()[0]
	logger.Info("Package identified: \"%s\" version %s", packageMetadata.GetApp().AppName, packageMetadata.GetCurrentVersion().GetVersionName())

	requestPacket.TaskType = &proto.RequestPacket_AppsTaskType{AppsTaskType: proto.TaskAppsTypes_TYPE_APP_PKG_CREATE}
	response, err = task.PostRequest(requestPacket)
	if err != nil {
		return err
	}

	if len(response.WorkerResponse) < 1 {
		if response.GetExtraData() != "" {
			return errors.New(response.GetExtraData())
		}
		return fmt.Errorf("error occurred while creating worker target package on server-side")
	}

	var packageArch string
	for _, workers := range response.GetWorkerResponse() {
		packageArch += workers.GetWorkerInfo().GetCpuArch().String() + ", "
	}

	packageArch = strings.TrimSuffix(packageArch, ", ")
	logger.Info("Supported architectures: %s", packageArch)
	logger.Info("Successfully registered at server!")
	logger.Tip("Tip: Check available packages using \"staploy-cli apps -n %s\"", packageMetadata.GetApp().GetAppName())

	return nil
}
