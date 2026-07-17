package registry

import (
	"errors"
	"fmt"
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
	"strings"
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

	if len(response.GetRegistryResponse().GetAppInfo()) < 1 {
		if response.GetErrorCause() != "" {
			return errors.New(response.GetErrorCause())
		}
		return errors.New("package is corrupted or missing metadata")
	}

	packageMetadata := response.GetRegistryResponse().GetAppInfo()[0]
	packagePrinter := fmt.Sprintf("Package identified: \"%s\" version %s", packageMetadata.GetApp().AppName, packageMetadata.GetCurrentVersion().GetVersionName())

	var packageArch string
	for _, arches := range packageMetadata.GetCurrentVersion().GetSupportedArch() {
		if arches == proto.CpuArch_UNKNOWN {
			packageArch += "share, "
			continue
		}
		packageArch += arches.String() + ", "
	}

	packageArch = strings.TrimSuffix(packageArch, ", ")
	logger.Info("%s, Supported architectures: %s", packagePrinter, packageArch)
	logger.Info("Successfully pushed to registry!")
	return nil
}
