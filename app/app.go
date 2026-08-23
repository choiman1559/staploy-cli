package app

import (
	"staploy-cli/app/cmds"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
	"staploy-cli/app/tasks/apps"
	"staploy-cli/app/tasks/build"
	"staploy-cli/app/tasks/deploy"
	"staploy-cli/app/tasks/groups"
	"staploy-cli/app/tasks/nodes"
	"staploy-cli/app/tasks/registry"
	"staploy-cli/app/tasks/user"

	"github.com/alexflint/go-arg"
)

//goland:noinspection DuplicatedCode
func HandleProcessInvoke() {
	arg.MustParse(&cmds.Args)
	logger.InitLogger(cmds.Args.NoColor)

	defaultArgs := cmds.DefaultArgs{
		Address:         cmds.Args.Address,
		Port:            cmds.Args.Port,
		Verbose:         cmds.Args.Verbose,
		UseWorkerIdOnly: cmds.Args.UseIdOnly,
	}

	cmds.InitCache(cmds.Args.DisableTls, cmds.Args.SkipValidation, cmds.Args.UserJwtToken)
	var taskInterface cmds.CmdTaskInterface

	checkList := []struct {
		isSet bool
		init  func() cmds.CmdTaskInterface
	}{
		/// Build Tasks
		{cmds.Args.Build != nil, func() cmds.CmdTaskInterface {
			t := &build.PkgCmdTask{}
			t.Init(defaultArgs, *cmds.Args.Build, proto.TaskGroup_TASK_NONE)
			return t
		}},
		{cmds.Args.File != nil, func() cmds.CmdTaskInterface {
			t := &build.StaFileTask{}
			t.Init(defaultArgs, *cmds.Args.File, proto.TaskGroup_TASK_NONE)
			return t
		}},
		{cmds.Args.Group != nil, func() cmds.CmdTaskInterface {
			t := &groups.GroupCmdTask{}
			t.Init(defaultArgs, *cmds.Args.Group, proto.TaskGroup_TASK_NONE)
			return t
		}},

		/// TaskAppsTypes
		{cmds.Args.Create != nil, func() cmds.CmdTaskInterface {
			t := &apps.CreateCmdTask{}
			t.Init(defaultArgs, *cmds.Args.Create, proto.TaskGroup_TASK_MANAGE_APPS)
			return t
		}},
		{cmds.Args.Delete != nil, func() cmds.CmdTaskInterface {
			t := &apps.DeleteCmdTask{}
			t.Init(defaultArgs, *cmds.Args.Delete, proto.TaskGroup_TASK_MANAGE_APPS)
			return t
		}},
		{cmds.Args.Apps != nil, func() cmds.CmdTaskInterface {
			t := &apps.AppsCmdTask{}
			t.Init(defaultArgs, *cmds.Args.Apps, proto.TaskGroup_TASK_MANAGE_APPS)
			return t
		}},
		{cmds.Args.Upload != nil, func() cmds.CmdTaskInterface {
			t := &apps.UploadCmdTask{}
			t.Init(defaultArgs, *cmds.Args.Upload, proto.TaskGroup_TASK_MANAGE_APPS)
			return t
		}},

		/// TaskNodeTypes
		{cmds.Args.List != nil, func() cmds.CmdTaskInterface {
			t := &nodes.ListCmdTask{}
			t.Init(defaultArgs, *cmds.Args.List, proto.TaskGroup_TASK_MANAGE_NODE)
			return t
		}},
		{cmds.Args.Fetch != nil, func() cmds.CmdTaskInterface {
			t := &nodes.FetchCmdTask{}
			t.Init(defaultArgs, *cmds.Args.Fetch, proto.TaskGroup_TASK_MANAGE_NODE)
			return t
		}},
		{cmds.Args.Bash != nil, func() cmds.CmdTaskInterface {
			t := &nodes.BashCmdTask{}
			t.Init(defaultArgs, *cmds.Args.Bash, proto.TaskGroup_TASK_MANAGE_NODE)
			return t
		}},
		{cmds.Args.Disconn != nil, func() cmds.CmdTaskInterface {
			t := &nodes.DisConnCmdTask{}
			t.Init(defaultArgs, *cmds.Args.Disconn, proto.TaskGroup_TASK_MANAGE_NODE)
			return t
		}},

		/// TaskDeployTypes
		{cmds.Args.Push != nil, func() cmds.CmdTaskInterface {
			t := &deploy.PushCmdTask{}
			t.Init(defaultArgs, *cmds.Args.Push, proto.TaskGroup_TASK_DEPLOY)
			return t
		}},
		{cmds.Args.Remove != nil, func() cmds.CmdTaskInterface {
			t := &deploy.RemoveCmdTask{}
			t.Init(defaultArgs, *cmds.Args.Remove, proto.TaskGroup_TASK_DEPLOY)
			return t
		}},
		{cmds.Args.Set != nil, func() cmds.CmdTaskInterface {
			t := &deploy.SetCmdTask{}
			t.Init(defaultArgs, *cmds.Args.Set, proto.TaskGroup_TASK_DEPLOY)
			return t
		}},
		{cmds.Args.User != nil, func() cmds.CmdTaskInterface {
			t := &user.UserCmdTask{}
			t.Init(defaultArgs, *cmds.Args.User, proto.TaskGroup_TASK_USER)
			return t
		}},
		{cmds.Args.Registry != nil, func() cmds.CmdTaskInterface {
			t := &registry.RegistryCmdTask{}
			t.Init(defaultArgs, *cmds.Args.Registry, proto.TaskGroup_TASK_REGISTRY)
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
		logger.Error("no command arg specified. abort")
		return
	}

	err := taskInterface.MainCmd()
	if err != nil {
		logger.Error("%v", err)
	}
}
