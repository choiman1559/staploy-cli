package registry

import (
	"staploy-cli/app/cmds"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
)

type RegistryCmdTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.RegistryCmd]
}

func (task *RegistryCmdTask) MainCmd() error {
	var taskInterface cmds.CmdTaskInterface
	checkList := []struct {
		isSet bool
		init  func() cmds.CmdTaskInterface
	}{
		{cmds.Args.Registry.PushLocal != nil, func() cmds.CmdTaskInterface {
			t := &RegistryPushLocalTask{}
			t.Init(task.DefaultArgs, *cmds.Args.Registry.PushLocal, proto.TaskGroup_TASK_REGISTRY)
			return t
		}},
		{cmds.Args.Registry.RemoveLocal != nil, func() cmds.CmdTaskInterface {
			t := &RegistryRemoveLocalTask{}
			t.Init(task.DefaultArgs, *cmds.Args.Registry.RemoveLocal, proto.TaskGroup_TASK_REGISTRY)
			return t
		}},
		{cmds.Args.Registry.ListLocal != nil, func() cmds.CmdTaskInterface {
			t := &RegistryListLocalTask{}
			t.Init(task.DefaultArgs, *cmds.Args.Registry.ListLocal, proto.TaskGroup_TASK_REGISTRY)
			return t
		}},

		{cmds.Args.Registry.Pull != nil, func() cmds.CmdTaskInterface {
			t := &RegistryPullTask{}
			t.Init(task.DefaultArgs, *cmds.Args.Registry.Pull, proto.TaskGroup_TASK_REGISTRY)
			return t
		}},
		{cmds.Args.Registry.ListRepo != nil, func() cmds.CmdTaskInterface {
			t := &RegistryListRepoTask{}
			t.Init(task.DefaultArgs, *cmds.Args.Registry.ListRepo, proto.TaskGroup_TASK_REGISTRY)
			return t
		}},
		{cmds.Args.Registry.AddRepo != nil, func() cmds.CmdTaskInterface {
			t := &RegistryAddRepoTask{}
			t.Init(task.DefaultArgs, *cmds.Args.Registry.AddRepo, proto.TaskGroup_TASK_REGISTRY)
			return t
		}},
		{cmds.Args.Registry.RemoveRepo != nil, func() cmds.CmdTaskInterface {
			t := &RegistryRemoveRepoTask{}
			t.Init(task.DefaultArgs, *cmds.Args.Registry.RemoveRepo, proto.TaskGroup_TASK_REGISTRY)
			return t
		}},
		{cmds.Args.Registry.ManageRepoToken != nil, func() cmds.CmdTaskInterface {
			t := &RegistryManageRepoTokenTask{}
			t.Init(task.DefaultArgs, *cmds.Args.Registry.ManageRepoToken, proto.TaskGroup_TASK_REGISTRY)
			return t
		}},
		{cmds.Args.Registry.List != nil, func() cmds.CmdTaskInterface {
			t := &RegistryListTask{}
			t.Init(task.DefaultArgs, *cmds.Args.Registry.List, proto.TaskGroup_TASK_REGISTRY)
			return t
		}},
		{cmds.Args.Registry.UpdateCache != nil, func() cmds.CmdTaskInterface {
			t := &RegistryUpdateCacheTask{}
			t.Init(task.DefaultArgs, *cmds.Args.Registry.UpdateCache, proto.TaskGroup_TASK_REGISTRY)
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
