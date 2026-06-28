package user

import (
	"staploy-cli/app/cmds"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
)

type UserCmdTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.UserCmd]
}

func (task *UserCmdTask) MainCmd() error {
	var taskInterface cmds.CmdTaskInterface
	checkList := []struct {
		isSet bool
		init  func() cmds.CmdTaskInterface
	}{
		{cmds.Args.User.UserPerm != nil, func() cmds.CmdTaskInterface {
			t := &UserPermTask{}
			t.Init(task.DefaultArgs, *cmds.Args.User.UserPerm, proto.TaskGroup_TASK_USER)
			return t
		}},
		{cmds.Args.User.UserCreate != nil, func() cmds.CmdTaskInterface {
			t := &UserCreateTask{}
			t.Init(task.DefaultArgs, *cmds.Args.User.UserCreate, proto.TaskGroup_TASK_USER)
			return t
		}},
		{cmds.Args.User.UserLogin != nil, func() cmds.CmdTaskInterface {
			t := &UserLoginTask{}
			t.Init(task.DefaultArgs, *cmds.Args.User.UserLogin, proto.TaskGroup_TASK_USER)
			return t
		}},
		{cmds.Args.User.UserRemove != nil, func() cmds.CmdTaskInterface {
			t := &UserRemoveTask{}
			t.Init(task.DefaultArgs, *cmds.Args.User.UserRemove, proto.TaskGroup_TASK_USER)
			return t
		}},
		{cmds.Args.User.UserList != nil, func() cmds.CmdTaskInterface {
			t := &UserListTask{}
			t.Init(task.DefaultArgs, *cmds.Args.User.UserList, proto.TaskGroup_TASK_USER)
			return t
		}},
		{cmds.Args.User.UserAudit != nil, func() cmds.CmdTaskInterface {
			t := &UserAuditTask{}
			t.Init(task.DefaultArgs, *cmds.Args.User.UserAudit, proto.TaskGroup_TASK_USER)
			return t
		}},
	}

	for _, item := range checkList {
		if item.isSet {
			taskInterface = item.init()
			break
		}
	}

	if taskInterface == nil {
		logger.Error("no user command arg specified. abort")
	}

	err := taskInterface.MainCmd()
	return err
}
