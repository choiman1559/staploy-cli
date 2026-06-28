package user

import (
	"errors"
	"fmt"
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
)

type UserPermTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.UserPermissionsCmd]
}

func (task *UserPermTask) MainCmd() error {
	requestPacket := task.CreateDefPacket()

	if task.CmdArgs.UserName == "" {
		return errors.New("user name is required")
	}

	if len(task.CmdTask.CmdArgs.PermCmd) == 0 {
		logger.Info("No permission specified, setting default permissions (PERMISSION_NONE)")
	}

	requestPacket.TaskType = &proto.RequestPacket_UserTaskType{UserTaskType: &proto.UserRequestPacket{
		UserTaskTypes: proto.TaskUserTypes_TYPE_USER_RBAC,
		UserLoginInfo: &proto.UserLoginInfo{
			UserName: task.CmdArgs.UserName,
		},
		Permissions: new(uint32(getPermFlags(task.CmdArgs.PermCmd))),
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

func getPermFlags(flags []cmds.CommandPermEnum) int32 {
	var flagValue = int32(proto.PermissionFlag_NONE)

	for _, flag := range flags {
		switch flag {
		case cmds.Create:
			flagValue = flagValue | int32(proto.PermissionFlag_APP_CREATE)
			break

		case cmds.Delete:
			flagValue = flagValue | int32(proto.PermissionFlag_APP_DELETE)
			break

		case cmds.Upload:
			flagValue = flagValue | int32(proto.PermissionFlag_APP_UPLOAD)
			break

		case cmds.Bash:
			flagValue = flagValue | int32(proto.PermissionFlag_NODE_BASH)
			break

		case cmds.Disconn:
			flagValue = flagValue | int32(proto.PermissionFlag_NODE_DISCONN)
			break

		case cmds.Push:
			flagValue = flagValue | int32(proto.PermissionFlag_NODE_PUSH)
			break

		case cmds.Remove:
			flagValue = flagValue | int32(proto.PermissionFlag_NODE_REMOVE)
			break

		case cmds.Set:
			flagValue = flagValue | int32(proto.PermissionFlag_NODE_SET)
			break

		case cmds.Query:
			flagValue = flagValue | int32(proto.PermissionFlag_QUERY_ENDPOINT)
			break

		case cmds.Group:
			flagValue = flagValue | int32(proto.PermissionFlag_GROUP_MANAGE)
			break

		case cmds.User:
			flagValue = flagValue | int32(proto.PermissionFlag_USER_MANAGE)
			break

		case cmds.Admin:
			flagValue = int32(proto.PermissionFlag_SYSTEM_ADMIN)
			break
		}
	}

	return flagValue
}
