package user

import (
	"errors"
	"fmt"
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
)

type UserLoginTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.UserLoginCmd]
}

func (task *UserLoginTask) MainCmd() error {
	requestPacket := task.CreateDefPacket()

	if task.CmdArgs.UserName == "" {
		return errors.New("user name is required")
	}

	password, err := ReadPasswordSecurely(fmt.Sprintf("[%s] password: ", task.CmdArgs.UserName))
	if err != nil {
		return err
	}

	requestPacket.TaskType = &proto.RequestPacket_UserTaskType{UserTaskType: &proto.UserRequestPacket{
		UserTaskTypes: proto.TaskUserTypes_TYPE_USER_LOGIN,
		UserLoginInfo: &proto.UserLoginInfo{
			UserName:     task.CmdArgs.UserName,
			UserPassword: password,
		},
	}}

	response, err := task.PostRequest(requestPacket)
	if err != nil {
		return err
	}

	if response.GetStatus() == consts.StatusOK {
		if !task.CmdArgs.PrintTokenOnly {
			logger.Info("Login successes.")
			logger.Tip("Tip: use \"export STAPLOY_USER_TOKEN=[YOUR_TOKEN]\" to use login credentials.")
		}
		fmt.Printf("%s\n", *response.GetUserResponse().UserLoginInfo.UserToken)
	} else {
		logger.Error("login error: %s", response.GetErrorCause())
	}
	return nil
}
