package cmds

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"staploy-cli/app/consts"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"

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

type DefaultArgs struct {
	CmdTaskInterface
	Address         string
	Port            int
	Verbose         bool
	UseWorkerIdOnly bool
}

type CmdTask[T CmdTypes] struct {
	DefaultArgs    DefaultArgs
	CmdArgs        T
	TaskGroups     proto.TaskGroup
	WorkersIdCache map[string]string
}

func (a *CmdTask[T]) Init(defArgs DefaultArgs, cmdArgs T, group proto.TaskGroup) {
	a.DefaultArgs = defArgs
	a.CmdArgs = cmdArgs
	a.TaskGroups = group
	a.WorkersIdCache = make(map[string]string)
}

func (a *CmdTask[T]) CreateDefPacket(workers ...string) *proto.RequestPacket {
	if a.DefaultArgs.UseWorkerIdOnly {
		return a.CreateDefPacketIdOnly(workers...)
	}

	workerRealIds := make(map[string]string)
	var toQueryWorkers []string

findCache:
	for _, worker := range workers {
		for id, name := range a.WorkersIdCache {
			if id == worker || name == worker {
				workerRealIds[id] = name
				continue findCache
			}
		}
		toQueryWorkers = append(toQueryWorkers, worker)
	}

	workerListPacket := a.CreateDefPacketIdOnly()
	workerListPacket.TaskGroup = proto.TaskGroup_TASK_GROUP
	workerListPacket.TaskType = &proto.RequestPacket_GroupTaskType{GroupTaskType: &proto.GroupRequestPacket{
		GroupTaskTypes: proto.TaskGroupTypes_TYPE_QUERY_WORKER_IDS,
		Names:          toQueryWorkers,
	}}

	response, err := a.PostRequest(workerListPacket)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	groupValidMap := make(map[string]bool)
	for _, worker := range response.GroupResponse {
		if worker.GetGroupName() != "" {
			if worker.GetWorkerInfo() != nil {
				if !groupValidMap[worker.GetGroupName()] {
					groupValidMap[worker.GetGroupName()] = true
					if a.DefaultArgs.Verbose {
						logger.Tip("[DEBUG] Identified group: %s", worker.GetGroupName())
					}
				}

				workerRealIds[worker.GetWorkerInfo().GetWorkerId()] = worker.GetWorkerInfo().GetWorkerName()
				a.WorkersIdCache[worker.GetWorkerInfo().GetWorkerId()] = worker.GetWorkerInfo().GetWorkerName()
			} else {
				logger.Warn("Requested group name \"%s\" not exists. Skipping...", worker.GetRequestedName())
			}
			continue
		}

		if worker.GetWorkerInfo() != nil {
			workerRealIds[worker.GetWorkerInfo().GetWorkerId()] = worker.GetWorkerInfo().GetWorkerName()
			a.WorkersIdCache[worker.GetWorkerInfo().GetWorkerId()] = worker.GetWorkerInfo().GetWorkerName()

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

	keys := make([]string, 0, len(workerRealIds))
	for k := range workerRealIds {
		keys = append(keys, k)
	}
	return a.CreateDefPacketIdOnly(keys...)
}

func (a *CmdTask[T]) CreateDefPacketIdOnly(workers ...string) *proto.RequestPacket {
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
	var paths = fmt.Sprintf(consts.APIRouteSchema, "v1", consts.ConnTypeAdmin)
	var addr = fmt.Sprintf("http://%s:%d%s", a.DefaultArgs.Address, a.DefaultArgs.Port, paths)

	data, err := protojson.Marshal(requestPacket)
	if err != nil {
		return err
	}

	resp, err := http.Post(addr, "application/json", bytes.NewBuffer(data))
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

	resp, err := http.Post(a.GetServerAddr(), "application/json", bytes.NewBuffer(data))
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
	var paths = fmt.Sprintf(consts.APIRouteSchema, "v1", consts.ConnTypeAdmin)
	var addr = fmt.Sprintf("http://%s:%d%s", a.DefaultArgs.Address, a.DefaultArgs.Port, paths)
	return addr
}
