package apps

import (
	"errors"
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
)

type DeleteCmdTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.DeleteCmd]
}

func (t *DeleteCmdTask) MainCmd() error {
	if t.CmdArgs.AppName == "" {
		return errors.New("--app-name field is required")
	}

	requestPacket := t.CreateDefPacket()
	requestPacket.TaskType = &proto.RequestPacket_AppsTaskType{AppsTaskType: proto.TaskAppsTypes_TYPE_APP_DELETE}
	t.CmdArgs.VersionName = logger.TrimVersions(t.CmdArgs.VersionName)

	versionList := make([]*proto.Version, len(t.CmdArgs.VersionName))
	for i, v := range t.CmdArgs.VersionName {
		versionList[i] = &proto.Version{
			VersionName: v,
		}
	}

	requestPacket.AppInfoFetch = append(requestPacket.GetAppInfoFetch(), &proto.AppInfoFetch{
		App:        &proto.AppInfo{AppName: t.CmdArgs.AppName},
		AppVersion: versionList,
	})

	response, err := t.PostRequest(requestPacket)
	if err != nil {
		return err
	}

	if response.GetStatus() == consts.StatusOK {
		if len(t.CmdArgs.VersionName) == 0 {
			logger.Info("Deleting app entry \"%s\" finished successfully", t.CmdArgs.AppName)
		} else {
			logger.Info("Deleting app entry \"%s\" (version %v) finished successfully", t.CmdArgs.AppName, t.CmdArgs.VersionName)
		}
	} else {
		logger.Error("Deleting app entry \"%s\" finished with error: %s", response.GetErrorCause())
	}
	return nil
}
