package cmds

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"staploy-cli/app/consts"
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
	Address string
	Port    int
	Verbose bool
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
