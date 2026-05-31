package nodes

import (
	"fmt"
	"log"
	"staploy-cli/app/cmds"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
)

type BashCmdTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.BashCmd]
}

func (a *BashCmdTask) MainCmd() error {
	packet := a.CreateDefPacket(a.CmdArgs.WorkerId)
	packet.TaskType = &proto.RequestPacket_NodeTaskType{NodeTaskType: proto.TaskNodeTypes_TYPE_NODE_EXECUTE_SHELL}
	packet.ExtraData = &a.CmdArgs.Command

	if a.CmdArgs.NoOutput {
		err := a.PostRequestOnly(packet)
		if err != nil {
			return err
		}
		return nil
	}

	response, err := a.PostRequest(packet)
	if err != nil {
		log.Fatal(err)
		return err
	}

	if len(response.WorkerResponse) > 0 {
		logger.Info("Shell executed, Output:\n", response.WorkerResponse[0].GetPacketInfo().GetExtraData())
	} else {
		return fmt.Errorf("specified worker is not valid or not connected to server")
	}
	return nil
}
