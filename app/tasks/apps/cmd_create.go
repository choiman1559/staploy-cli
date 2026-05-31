package apps

import (
	"errors"
	"staploy-cli/app/cmds"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
)

type CreateCmdTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.CreateCmd]
}

func (t *CreateCmdTask) MainCmd() error {
	requestPacket := t.CreateDefPacket()
	requestPacket.TaskType = &proto.RequestPacket_AppsTaskType{AppsTaskType: proto.TaskAppsTypes_TYPE_APP_REGISTER}

	if t.CmdArgs.AppName == "" {
		return errors.New("--app-name field is required")
	}

	requestPacket.AppInfoFetch = append(requestPacket.GetAppInfoFetch(), &proto.AppInfoFetch{
		App: &proto.AppInfo{
			AppName:        t.CmdArgs.AppName,
			AppDescription: &t.CmdArgs.AppDescription,
		},
	})

	response, err := t.PostRequest(requestPacket)
	if err != nil {
		return err
	}

	logger.Info(response.GetExtraData())
	return nil
}
