package build

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type StaFileTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.StaFileCmd]

	WorkerInfos         map[string]*proto.WorkerInfo
	ResolvedWorkerAlias map[string]*WorkerAlias
	ResolvedAppAlias    map[string]*ResolvedAppAlias
	ServerPort          int
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
		serverAddress, err := a.parseExecArgs(&Build{}, staFile.Config.Address)
		if err != nil {
			return fmt.Errorf("parse server address error: %v", err)
		}

		IsValidAddress := func(address string) bool {
			if net.ParseIP(address) != nil {
				return true
			}

			var hostnameRegex = regexp.MustCompile(`^(?=.{1,253}$)(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)*[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?$`)
			return hostnameRegex.MatchString(address)
		}

		if serverAddress == "" || !IsValidAddress(serverAddress) {
			return fmt.Errorf("invalid server address: %s", serverAddress)
		}

		serverPortStr, err := a.parseExecArgs(&Build{}, staFile.Config.Port)
		if serverPortStr == "" || err != nil {
			return fmt.Errorf("parse server port error: %v", err)
		}

		serverPort, err := strconv.Atoi(serverPortStr)
		if err != nil {
			return fmt.Errorf("cannot convert server port value into integer error: %v", err)
		}

		staFile.Config.Address = serverAddress
		staFile.Config.Port = serverPortStr

		a.ServerPort = serverPort
		a.DefaultArgs.UseWorkerIdOnly = staFile.Config.UseIdOnly
		a.DefaultArgs.Port = a.ServerPort
		a.DefaultArgs.Address = staFile.Config.Address
	}

	a.WorkerInfos = make(map[string]*proto.WorkerInfo)
	a.ResolvedWorkerAlias = make(map[string]*WorkerAlias)
	a.ResolvedAppAlias = make(map[string]*ResolvedAppAlias)

	return a.ParseStaFile(&staFile)
}

func (a *StaFileTask) ParseWorkerAlias(disableGroup bool, workers ...string) ([]string, error) {
	var aliases, restAction []string
	for _, worker := range workers {
		if strings.HasPrefix(worker, consts.STAFILE_ALIAS_PREFIX) {
			aliasName := strings.TrimPrefix(worker, consts.STAFILE_ALIAS_PREFIX)
			workerAlias := a.ResolvedWorkerAlias[aliasName]

			if workerAlias != nil {
				aliases = append(aliases, workerAlias.WorkerIds...)
			} else {
				return nil, fmt.Errorf("worker alias %s is not found", aliasName)
			}
		} else {
			restAction = append(restAction, worker)
		}
	}

	restResult, err := a.ParseWorkers(disableGroup, restAction...)
	if err != nil {
		return nil, err
	}

	return append(aliases, restResult...), nil
}

func (a *StaFileTask) FormatWorker(workerIdOrName string) string {
	realId, err := a.ParseWorkerAlias(true, workerIdOrName)
	if err != nil || len(realId) != 1 {
		return ""
	}

	if strings.HasPrefix(workerIdOrName, consts.STAFILE_GROUP_PREFIX) {
		return workerIdOrName
	}
	return fmt.Sprintf("%s (%s)", logger.ShortHash(realId[0]), cmds.WorkersIdCache[realId[0]])
}

func (a *StaFileTask) ParseStaFile(staployFile *StaployFile) error {
	logger.Process("Parsing Staployfile... Found %d build, %d manages and %d targets", len(staployFile.Builds), len(staployFile.Manages), len(staployFile.Targets))

	if staployFile.Alias == nil && len(staployFile.Builds) > 0 {
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
		Port:            a.ServerPort,
		UseWorkerIdOnly: staployFile.Config.UseIdOnly,
		Verbose:         a.DefaultArgs.Verbose,
	}

	if defaultArgs.Address == "" {
		return fmt.Errorf("address is required, Abort")
	}

	if defaultArgs.Port == 0 {
		return fmt.Errorf("port is required, Abort")
	}

	logger.Process("Syncing with Server... (%s:%d)", defaultArgs.Address, defaultArgs.Port)
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

	if staployFile.Alias != nil {
		err := a.processAlias(defaultArgs, staployFile.Alias)
		if err != nil {
			return err
		}
	}

	if staployFile.Alias != nil && len(staployFile.Builds) > 0 {
		err := a.processBuild(staployFile.Builds)
		if err != nil {
			return err
		}
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

func (a *StaFileTask) askWorkerInfoData(worker string) (*proto.WorkerInfo, error) {
	if a.WorkerInfos[worker] != nil {
		return a.WorkerInfos[worker], nil
	}

	packet := a.CreateDefPacket(worker)
	packet.TaskGroup = proto.TaskGroup_TASK_MANAGE_NODE
	packet.TaskType = &proto.RequestPacket_NodeTaskType{NodeTaskType: proto.TaskNodeTypes_TYPE_NODE_REQ_WORKER_INFO}

	response, err := a.PostRequest(packet)
	if err != nil {
		return nil, fmt.Errorf("failed to post request for worker \"%s\", cause: %s", worker, err.Error())
	}

	if len(response.GetWorkerResponse()) != 1 {
		return nil, fmt.Errorf("failed to get worker (\"%s\") information", worker)
	}

	workerInfo := response.GetWorkerResponse()[0].GetWorkerInfo()
	a.WorkerInfos[workerInfo.GetWorkerId()] = workerInfo
	return workerInfo, nil
}
