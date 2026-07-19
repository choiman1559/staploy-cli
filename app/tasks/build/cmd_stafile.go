package build

import (
	"fmt"
	"os"
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type StaFileTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.StaFileCmd]
}

func (a *StaFileTask) MainCmd() error {
	_, err := os.Stat(a.CmdArgs.ConfigFile)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("config file %s does not exist", a.CmdArgs.ConfigFile)
		}
	}

	bytes, err := os.ReadFile(a.CmdArgs.ConfigFile)
	if err != nil {
		return fmt.Errorf("file read error: %v", err)
	}

	var staFile StaployFile
	file, diags := hclsyntax.ParseConfig(bytes, a.CmdArgs.ConfigFile, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return fmt.Errorf("parse error at %v", diags)
	}

	if file == nil {
		return fmt.Errorf("config file %s does not exist", a.CmdArgs.ConfigFile)
	}

	diags = gohcl.DecodeBody(file.Body, nil, &staFile)
	if diags.HasErrors() {
		logger.Tip("%+v", diags)
		return fmt.Errorf("decode error at %v", err)
	}

	if a.DefaultArgs.Verbose {
		logger.Tip("[DEBUG] Configuration is %+v", staFile)
	}

	if staFile.Config != nil {
		a.DefaultArgs.UseWorkerIdOnly = staFile.Config.UseIdOnly
		a.DefaultArgs.Port = staFile.Config.Port
		a.DefaultArgs.Address = staFile.Config.Address
	}

	return a.ParseStaFile(&staFile)
}

func (a *StaFileTask) FormatWorker(workerIdOrName string) string {
	realId, err := a.ParseWorkers(true, workerIdOrName)
	if err != nil || len(realId) != 1 {
		return ""
	}

	if strings.HasPrefix(workerIdOrName, "group:") {
		return workerIdOrName
	}
	return fmt.Sprintf("%s (%s)", logger.ShortHash(realId[0]), cmds.WorkersIdCache[realId[0]])
}

func (a *StaFileTask) ParseStaFile(staployFile *StaployFile) error {
	logger.Process("Parsing Staployfile... Found %d build, %d manages and %d targets", len(staployFile.Builds), len(staployFile.Manages), len(staployFile.Targets))

	if len(staployFile.Builds) > 0 {
		err := a.processBuild(staployFile.Builds)
		if err != nil {
			return err
		}
	}

	if len(staployFile.Manages) < 1 && len(staployFile.Targets) < 1 {
		return nil
	} else if staployFile.Config == nil {
		return fmt.Errorf("manage and target block requires configure block")
	}

	defaultArgs := &cmds.DefaultArgs{
		Address:         staployFile.Config.Address,
		Port:            staployFile.Config.Port,
		UseWorkerIdOnly: staployFile.Config.UseIdOnly,
		Verbose:         a.DefaultArgs.Verbose,
	}

	if defaultArgs.Address == "" {
		return fmt.Errorf("address is required, Abort")
	}

	if defaultArgs.Port == 0 {
		return fmt.Errorf("port is required, Abort")
	}

	logger.Process("Syncing with Server (%s:%d)...", defaultArgs.Address, defaultArgs.Port)
	workerListPacket := a.CreateDefPacket()
	workerListPacket.TaskGroup = proto.TaskGroup_TASK_MANAGE_NODE
	workerListPacket.TaskType = &proto.RequestPacket_NodeTaskType{NodeTaskType: proto.TaskNodeTypes_TYPE_NODE_CONNECTED}

	response, err := a.PostRequest(workerListPacket)
	if err != nil {
		return fmt.Errorf("failed to connecting server: %s", err)
	}

	if response == nil {
		return fmt.Errorf("response is nil")
	}

	if response.Status != consts.StatusOK {
		return fmt.Errorf("server responded with %s", response.GetErrorCause())
	}

	if a.DefaultArgs.Verbose {
		logger.Tip("[DEBUG] Current workers are %+v", response.GetWorkerResponse())
	}

	err = a.processManage(defaultArgs, staployFile.Manages)
	if err != nil {
		return err
	}

	err = a.processTarget(defaultArgs, staployFile.Targets)
	if err != nil {
		return err
	}
	return nil
}
