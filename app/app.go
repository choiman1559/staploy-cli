package app

import (
	"log"
	"staploy-cli/app/build"
	"staploy-cli/app/cmds"
	"staploy-cli/app/nodes"
	"staploy-cli/app/proto"

	"github.com/alexflint/go-arg"
)

func HandleProcessInvoke() {
	arg.MustParse(&cmds.Args)

	defaultArgs := cmds.DefaultArgs{
		Address: cmds.Args.Address,
		Port:    cmds.Args.Port,
		Verbose: cmds.Args.Verbose,
	}

	var taskInterface cmds.CmdTaskInterface

	if cmds.Args.Build != nil {
		cmdTask := &build.PkgCmdTask{}
		cmdTask.Init(defaultArgs, *cmds.Args.Build, proto.TaskGroup_TASK_NONE)
		taskInterface = cmdTask
	}

	if cmds.Args.List != nil {
		cmdTask := &nodes.ListCmdTask{}
		cmdTask.Init(defaultArgs, *cmds.Args.List, proto.TaskGroup_TASK_MANAGE_NODE)
		taskInterface = cmdTask
	}

	if cmds.Args.Fetch != nil {
		cmdTask := &nodes.FetchCmdTask{}
		cmdTask.Init(defaultArgs, *cmds.Args.Fetch, proto.TaskGroup_TASK_MANAGE_NODE)
		taskInterface = cmdTask
	}

	if cmds.Args.Bash != nil {
		cmdTask := &nodes.BashCmdTask{}
		cmdTask.Init(defaultArgs, *cmds.Args.Bash, proto.TaskGroup_TASK_MANAGE_NODE)
		taskInterface = cmdTask
	}

	if cmds.Args.Disconn != nil {
		cmdTask := &nodes.DisConnCmdTask{}
		cmdTask.Init(defaultArgs, *cmds.Args.Disconn, proto.TaskGroup_TASK_MANAGE_NODE)
		taskInterface = cmdTask
	}

	err := taskInterface.MainCmd()
	if err != nil {
		log.Fatal(err)
	}
}
