package groups

import (
	"staploy-cli/app/cmds"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
)

type GroupCmdTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.GroupCmd]
}

func (task *GroupCmdTask) MainCmd() error {
	var taskInterface cmds.CmdTaskInterface
	checkList := []struct {
		isSet bool
		init  func() cmds.CmdTaskInterface
	}{
		{cmds.Args.Group.GroupCreate != nil, func() cmds.CmdTaskInterface {
			t := &GroupCreateTask{}
			t.Init(task.DefaultArgs, *cmds.Args.Group.GroupCreate, proto.TaskGroup_TASK_GROUP)
			return t
		}},
		{cmds.Args.Group.GroupAdd != nil, func() cmds.CmdTaskInterface {
			t := &GroupAddTask{}
			t.Init(task.DefaultArgs, *cmds.Args.Group.GroupAdd, proto.TaskGroup_TASK_GROUP)
			return t
		}},
		{cmds.Args.Group.GroupList != nil, func() cmds.CmdTaskInterface {
			t := &GroupListTask{}
			t.Init(task.DefaultArgs, *cmds.Args.Group.GroupList, proto.TaskGroup_TASK_GROUP)
			return t
		}},
		{cmds.Args.Group.GroupDelete != nil, func() cmds.CmdTaskInterface {
			t := &GroupDeleteTask{}
			t.Init(task.DefaultArgs, *cmds.Args.Group.GroupDelete, proto.TaskGroup_TASK_GROUP)
			return t
		}},
		{cmds.Args.Group.GroupRemove != nil, func() cmds.CmdTaskInterface {
			t := &GroupRemoveTask{}
			t.Init(task.DefaultArgs, *cmds.Args.Group.GroupRemove, proto.TaskGroup_TASK_GROUP)
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
		logger.Error("no group command arg specified. abort")
	}

	err := taskInterface.MainCmd()
	return err
}
