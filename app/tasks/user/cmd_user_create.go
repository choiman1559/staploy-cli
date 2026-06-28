package user

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
	"syscall"

	"golang.org/x/term"
)

type UserCreateTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.UserCreateCmd]
}

func (task *UserCreateTask) MainCmd() error {
	requestPacket := task.CreateDefPacket()

	if task.CmdArgs.UserName == "" {
		return errors.New("user name is required")
	}

	if task.CmdArgs.Refresh {
		requestPacket.TaskType = &proto.RequestPacket_UserTaskType{UserTaskType: &proto.UserRequestPacket{
			UserTaskTypes: proto.TaskUserTypes_TYPE_TOKEN_REFRESH,
			UserLoginInfo: &proto.UserLoginInfo{
				UserName: task.CmdArgs.UserName,
			},
		}}
	} else {
		password, err := ReadPasswordSecurely("New password: ")
		if err != nil {
			return err
		}

		passwdConfirm, err := ReadPasswordSecurely("Confirm password: ")
		if err != nil {
			return err
		}

		if !bytes.Equal(passwdConfirm, password) {
			return errors.New("passwords do not match")
		}

		requestPacket.TaskType = &proto.RequestPacket_UserTaskType{UserTaskType: &proto.UserRequestPacket{
			UserTaskTypes: proto.TaskUserTypes_TYPE_USER_CREATE,
			UserLoginInfo: &proto.UserLoginInfo{
				UserName:     task.CmdArgs.UserName,
				UserPassword: password,
			},
		}}
	}

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

func ReadPasswordSecurely(prompt string) ([]byte, error) {
	fmt.Print(prompt)
	fd := syscall.Stdin

	rawPassword, err := term.ReadPassword(fd)
	fmt.Println()
	if err != nil {
		return nil, fmt.Errorf("failed to read password securely: %w", err)
	}

	if len(rawPassword) == 0 {
		return nil, fmt.Errorf("password cannot be empty")
	}

	hasher := sha256.New()
	hasher.Write(rawPassword)
	clientHash := hasher.Sum(nil)

	for i := range rawPassword {
		rawPassword[i] = 0
	}
	return clientHash, nil
}
