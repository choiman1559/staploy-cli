package nodes

import (
	"errors"
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
	workers, err := a.ParseWorkers(false, a.CmdArgs.WorkerId)
	if err != nil {
		return err
	}

	for _, worker := range workers {
		packet := a.CreateDefPacket(worker)
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
			logger.Info("Shell executed on worker %s, Output:\n%s", logger.ShortHash(worker), response.WorkerResponse[0].GetPacketInfo().GetExtraData())
			return nil
		} else if response.GetErrorCause() != "" {
			return errors.New(response.GetErrorCause())
		}
		return fmt.Errorf("specified worker is not valid or not connected to server")
	}
	return nil
}
