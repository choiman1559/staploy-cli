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
	AppsCmd | BashCmd | BuildCmd | CreateCmd |
		DeleteCmd | DisconnCmd | FetchCmd | ListCmd |
		PushCmd | RemoveCmd | SetCmd | UploadCmd
}

type TaskTypes interface {
	proto.TaskAppsTypes | proto.TaskNodeTypes | proto.TaskDeployTypes
}

type CmdTaskInterface interface {
	MainCmd() error
}

type DefaultArgs struct {
	CmdTaskInterface
	Address       string
	Port          int
	Verbose       bool
	UseWorkerName bool
}

type CmdTask[T CmdTypes] struct {
	DefaultArgs DefaultArgs
	CmdArgs     T
	TaskGroups  proto.TaskGroup
}

func (a *CmdTask[T]) Init(defArgs DefaultArgs, cmdArgs T, group proto.TaskGroup) {
	a.DefaultArgs = defArgs
	a.CmdArgs = cmdArgs
	a.TaskGroups = group
}

func (a *CmdTask[T]) CreateDefPacket(workers ...string) *proto.RequestPacket {
	if !a.DefaultArgs.UseWorkerName {
		return a.CreateDefPacketIdOnly(workers...)
	}

	workerListPacket := a.CreateDefPacketIdOnly()
	workerListPacket.TaskGroup = proto.TaskGroup_TASK_MANAGE_NODE
	workerListPacket.TaskType = &proto.RequestPacket_NodeTaskType{NodeTaskType: proto.TaskNodeTypes_TYPE_NODE_CONNECTED}

	response, err := a.PostRequest(workerListPacket)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	var workerRealIds []string
	for _, workerInfo := range response.WorkerResponse {
		for _, workerRaw := range workers {
			if workerInfo.GetWorkerInfo().GetWorkerId() == workerRaw {
				workerRealIds = append(workerRealIds, workerRaw)
			} else if workerInfo.GetWorkerInfo().GetWorkerName() == workerRaw {
				workerRealIds = append(workerRealIds, workerInfo.GetWorkerInfo().GetWorkerId())
			} else {
				logger.Warn("cannot determine given element (%s) is id or name. check given id worker is connected to server", workerRaw)
				os.Exit(1)
			}
		}
	}

	return a.CreateDefPacketIdOnly(workerRealIds...)
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
