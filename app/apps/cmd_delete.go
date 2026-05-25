package apps

import (
	"errors"
	"log"
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
	"staploy-cli/app/proto"
)

type DeleteCmdTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.DeleteCmd]
}

func (t *DeleteCmdTask) MainCmd() error {
	if t.CmdArgs.AppName == "" {
		return errors.New("AppName is required")
	}

	requestPacket := t.CreateDefPacket()
	requestPacket.TaskType = &proto.RequestPacket_AppsTaskType{AppsTaskType: proto.TaskAppsTypes_TYPE_APP_DELETE}

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
		log.Println("Delete Command Success")
	}
	return nil
}
