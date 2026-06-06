package nodes

import (
	"fmt"
	"log"
	"staploy-cli/app/cmds"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
)

type ListCmdTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.ListCmd]
}

func (a *ListCmdTask) MainCmd() error {
	workers, err := a.ParseWorkers(false, a.CmdArgs.WorkerId...)
	if err != nil {
		return err
	}

	packet := a.CreateDefPacket(workers...)
	if a.CmdArgs.Refresh {
		if len(a.CmdArgs.WorkerId) <= 0 {
			return fmt.Errorf("using --refresh option requires at least one worker id")
		}
		packet.TaskType = &proto.RequestPacket_NodeTaskType{NodeTaskType: proto.TaskNodeTypes_TYPE_NODE_REQ_WORKER_INFO}
	} else {
		packet.TaskType = &proto.RequestPacket_NodeTaskType{NodeTaskType: proto.TaskNodeTypes_TYPE_NODE_CONNECTED}
	}

	response, err := a.PostRequest(packet)
	if err != nil {
		log.Fatal(err)
		return err
	}

	for i, worker := range response.WorkerResponse {
		for _, str := range WorkerInfoFormatter(worker.GetWorkerInfo(), a.CmdArgs.Detail, i) {
			fmt.Print(str)
		}

		if a.CmdArgs.Detail && i < len(response.WorkerResponse)-1 {
			fmt.Print("\n")
		}
	}
	return nil
}

func WorkerInfoFormatter(workerInfo *proto.WorkerInfo, detail bool, index int) []string {
	var workerData []string
	if detail {
		workerData = append(workerData, fmt.Sprintf("* Worker #%d %s\n", index, workerInfo.GetWorkerId()))
		workerData = append(workerData, fmt.Sprintf(" ├─ Name: %s\n", workerInfo.GetWorkerName()))
		workerData = append(workerData, fmt.Sprintf(" ├─ Cpu: %s (%d Core(s))\n", workerInfo.GetCpuArch().String(), workerInfo.GetCpuCoreCount()))
		workerData = append(workerData, fmt.Sprintf(" ├─ Memory: %s (%s bytes)\n", logger.FormatBytes(workerInfo.GetMemoryInBytes()), logger.FormatWithCommas(workerInfo.GetMemoryInBytes())))
		workerData = append(workerData, fmt.Sprintf(" ├─ Working directory: %s\n", workerInfo.GetBinLocation()))
		workerData = append(workerData, fmt.Sprintf(" └─ Additional flags\n"))
		workerData = append(workerData, fmt.Sprintf("     ├─ Default buffer size: %s bytes\n", logger.FormatWithCommas(workerInfo.WorkerFlags.BUFFER_SIZE)))
		workerData = append(workerData, fmt.Sprintf("     ├─ Remote shell enabled: %v\n", workerInfo.WorkerFlags.USE_REMOTE_SHELL))
		workerData = append(workerData, fmt.Sprintf("     └─ Skip executable integrity check: %v\n", workerInfo.WorkerFlags.SKIP_HASH_VERIFICATION))
	} else {
		workerData = append(workerData, fmt.Sprintf("Worker #%d %s (%s)\n", index, workerInfo.GetWorkerId(), workerInfo.GetWorkerName()))
	}
	return workerData
}
