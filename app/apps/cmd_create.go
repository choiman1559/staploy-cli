package apps

import (
	"errors"
	"log"
	"staploy-cli/app/cmds"
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
		return errors.New("AppName is required")
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

	log.Println(response.GetExtraData())
	return nil
}
