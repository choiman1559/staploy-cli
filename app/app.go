package app

import (
	"log"
	"staploy-cli/app/apps"
	"staploy-cli/app/build"
	"staploy-cli/app/cmds"
	"staploy-cli/app/deploy"
	"staploy-cli/app/nodes"
	"staploy-cli/app/proto"

	"github.com/alexflint/go-arg"
)

//goland:noinspection DuplicatedCode
func HandleProcessInvoke() {
	arg.MustParse(&cmds.Args)

	defaultArgs := cmds.DefaultArgs{
		Address:       cmds.Args.Address,
		Port:          cmds.Args.Port,
		Verbose:       cmds.Args.Verbose,
		UseWorkerName: cmds.Args.UseName,
	}

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
	}

	for _, item := range checkList {
		if item.isSet {
			taskInterface = item.init()
			break
		}
	}

	if taskInterface == nil {
		log.Fatal("no command arg specified")
	}

	err := taskInterface.MainCmd()
	if err != nil {
		log.Fatal(err)
	}
}
