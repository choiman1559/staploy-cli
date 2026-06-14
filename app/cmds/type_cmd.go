package cmds

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"staploy-cli/app/consts"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"
)

type CmdTypes interface {
	AppsCmd | BashCmd | BuildCmd | CreateCmd | DeleteCmd | DisconnCmd |
		FetchCmd | ListCmd | PushCmd | RemoveCmd | SetCmd | UploadCmd | StaFileCmd |
		GroupCmd | GroupAddCmd | GroupCreateCmd | GroupDeleteCmd | GroupListCmd | GroupRemoveCmd
}

type TaskTypes interface {
	proto.TaskAppsTypes | proto.TaskNodeTypes | proto.TaskDeployTypes
}

type CmdTaskInterface interface {
	MainCmd() error
}

var DisableTls bool
var SkipValidation bool
var WorkersIdCache map[string]string
var GroupValidMap map[string]bool

type DefaultArgs struct {
	CmdTaskInterface
	Address         string
	Port            int
	Verbose         bool
	UseWorkerIdOnly bool
}

type CmdTask[T CmdTypes] struct {
	DefaultArgs DefaultArgs
	CmdArgs     T
	TaskGroups  proto.TaskGroup
}

func InitCache(disableTls bool, skipValidation bool) {
	DisableTls = disableTls
	SkipValidation = skipValidation
	WorkersIdCache = make(map[string]string)
	GroupValidMap = make(map[string]bool)
}

func (a *CmdTask[T]) Init(defArgs DefaultArgs, cmdArgs T, group proto.TaskGroup) {
	a.DefaultArgs = defArgs
	a.CmdArgs = cmdArgs
	a.TaskGroups = group
}

func (a *CmdTask[T]) ParseWorkers(disableGroup bool, workers ...string) ([]string, error) {
	if workers == nil || len(workers) == 0 {
		return []string{}, nil
	}

	workerRealIds := make(map[string]string)
	var toQueryWorkers []string

findCache:
	for _, worker := range workers {
		if disableGroup && strings.HasPrefix(worker, "group:") {
			return nil, fmt.Errorf("group is not supported on this type of task: %s", worker)
		}

		for id, name := range WorkersIdCache {
			if id == worker || name == worker {
				workerRealIds[id] = name
				continue findCache
			}
		}
		toQueryWorkers = append(toQueryWorkers, worker)
	}

	if len(toQueryWorkers) > 0 {
		workerListPacket := a.CreateDefPacket()
		workerListPacket.TaskGroup = proto.TaskGroup_TASK_GROUP
		workerListPacket.TaskType = &proto.RequestPacket_GroupTaskType{GroupTaskType: &proto.GroupRequestPacket{
			GroupTaskTypes: proto.TaskGroupTypes_TYPE_QUERY_WORKER_IDS,
			Names:          toQueryWorkers,
		}}

		response, err := a.PostRequest(workerListPacket)
		if err != nil {
			return nil, err
		}

		if response.GetStatus() == consts.StatusError {
			logger.Error("Server response reported exception: %s", response.GetErrorCause())
		}

		for _, worker := range response.GroupResponse {
			if worker.GetWorkerInfo() != nil && !worker.GetIsAlive() {
				logger.Warn("Worker %s (%s) identified but not alive, Ignoring...", logger.ShortHash(worker.GetWorkerInfo().GetWorkerName()), worker.GetWorkerInfo().GetWorkerName())
				continue
			}

			if worker.GetGroupName() != "" {
				if worker.GetWorkerInfo() != nil {
					if !GroupValidMap[worker.GetGroupName()] {
						GroupValidMap[worker.GetGroupName()] = true
						if a.DefaultArgs.Verbose {
							logger.Tip("[DEBUG] Identified group: %s", worker.GetGroupName())
						}
					}

					workerRealIds[worker.GetWorkerInfo().GetWorkerId()] = worker.GetWorkerInfo().GetWorkerName()
					WorkersIdCache[worker.GetWorkerInfo().GetWorkerId()] = worker.GetWorkerInfo().GetWorkerName()
				} else {
					logger.Warn("Requested group name \"%s\" not exists. Skipping...", worker.GetRequestedName())
				}
				continue
			}

			if worker.GetWorkerInfo() != nil {
				workerRealIds[worker.GetWorkerInfo().GetWorkerId()] = worker.GetWorkerInfo().GetWorkerName()
				WorkersIdCache[worker.GetWorkerInfo().GetWorkerId()] = worker.GetWorkerInfo().GetWorkerName()

				if a.DefaultArgs.Verbose {
					if worker.RequestedName == worker.GetWorkerInfo().GetWorkerId() {
						logger.Tip("[DEBUG] Identified worker as Id: %s", logger.ShortHash(worker.GetWorkerInfo().GetWorkerId()))
					} else {
						logger.Tip("[DEBUG] Identified worker as Name: %s", logger.ShortHash(worker.GetWorkerInfo().GetWorkerId()))
					}
				}
				continue
			}
			logger.Warn("Requested identify \"%s\" is nor id, name, group; Skipping...", worker.GetRequestedName())
		}
	}

	keys := make([]string, 0, len(workerRealIds))
	for k := range workerRealIds {
		keys = append(keys, k)
	}

	if len(keys) < 1 {
		return nil, fmt.Errorf("no worker identified for %v", workers)
	}
	return keys, nil
}

func (a *CmdTask[T]) CreateDefPacket(workers ...string) *proto.RequestPacket {
	packet := &proto.RequestPacket{
		TaskGroup: a.TaskGroups,
	}

	for _, worker := range workers {
		packet.Worker = append(packet.Worker, &proto.WorkerInfo{WorkerId: worker})
	}

	return packet
}

//goland:noinspection HttpUrlsUsage
func (a *CmdTask[T]) PostRequestOnly(requestPacket *proto.RequestPacket) error {
	data, err := protojson.Marshal(requestPacket)
	if err != nil {
		return err
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: SkipValidation,
	}

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	resp, err := client.Post(a.GetServerAddr(), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	return nil
}

func (a *CmdTask[T]) PostRequest(requestPacket *proto.RequestPacket) (*proto.ResponsePacket, error) {
	data, err := protojson.Marshal(requestPacket)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: SkipValidation,
	}

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	resp, err := client.Post(a.GetServerAddr(), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	responsePacket := &proto.ResponsePacket{}
	err = protojson.Unmarshal(body, responsePacket)
	if err != nil {
		return nil, err
	}
	return responsePacket, nil
}

//goland:noinspection HttpUrlsUsage
func (a *CmdTask[T]) GetServerAddr() string {
	httpPrefix := "https://"
	if DisableTls {
		httpPrefix = "http://"
	}

	var paths = fmt.Sprintf(consts.APIRouteSchema, "v1", consts.ConnTypeAdmin)
	var addr = fmt.Sprintf("%s%s:%d%s", httpPrefix, a.DefaultArgs.Address, a.DefaultArgs.Port, paths)
	return addr
}
