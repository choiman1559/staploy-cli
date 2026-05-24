package nodes

import (
	"log"
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
	"staploy-cli/app/proto"
)

type DisConnCmdTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.DisconnCmd]
}

func (a *DisConnCmdTask) MainCmd() error {
	for _, ids := range a.CmdArgs.WorkerId {
		packet := a.CreateDefPacket(ids)
		packet.TaskType = &proto.RequestPacket_NodeTaskType{NodeTaskType: proto.TaskNodeTypes_TYPE_NODE_DISCONN_WORKER}

		response, err := a.PostRequest(packet)
		if err != nil {
			log.Fatal(err)
			return err
		}

		if response != nil && response.GetStatus() == consts.StatusOK {
			log.Print("Disconnected worker " + ids)
		}
	}
	return nil
}
