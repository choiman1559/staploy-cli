package user

import (
	"errors"
	"fmt"
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
)

type UserRemoveTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.UserRemoveCmd]
}

func (task *UserRemoveTask) MainCmd() error {
	requestPacket := task.CreateDefPacket()

	if task.CmdArgs.UserName == "" {
		return errors.New("user name is required")
	}

	requestPacket.TaskType = &proto.RequestPacket_UserTaskType{UserTaskType: &proto.UserRequestPacket{
		UserTaskTypes: proto.TaskUserTypes_TYPE_USER_REMOVE,
		UserLoginInfo: &proto.UserLoginInfo{
			UserName: task.CmdArgs.UserName,
		},
	}}

	response, err := task.PostRequest(requestPacket)
	if err != nil {
		return err
	}

	if response.GetStatus() == consts.StatusOK {
		logger.Info("%s", response.GetExtraData())
	} else {
		return fmt.Errorf("%s", response.GetErrorCause())
	}
	return nil
}
