package groups

import (
	"bytes"
	"fmt"
	"staploy-cli/app/cmds"
	"staploy-cli/app/proto"
	"strings"
	"text/tabwriter"
)

type GroupListTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.GroupListCmd]
}

func (task *GroupListTask) MainCmd() error {
	requestPacket := task.CreateDefPacket()
	requestPacket.TaskType = &proto.RequestPacket_GroupTaskType{GroupTaskType: &proto.GroupRequestPacket{
		GroupTaskTypes: proto.TaskGroupTypes_TYPE_QUERY_GROUP_LIST,
		GroupName:      &task.CmdArgs.GroupName,
	}}

	response, err := task.PostRequest(requestPacket)
	if err != nil {
		return err
	}

	for _, str := range task.parseResult(response.GetGroupResponse(), task.CmdArgs.GroupName == "") {
		fmt.Printf(str)
	}

	return nil
}

func (task *GroupListTask) parseResult(groupResponse []*proto.GroupResponsePacket, allGroupOverview bool) []string {
	var resultData []string
	var sb strings.Builder

	if allGroupOverview {
		sb.WriteString(fmt.Sprintf("Total %d group(s) exist.\n\n", len(groupResponse)))

		for _, group := range groupResponse {
			sb.WriteString(fmt.Sprintf("* Group: %s\n", group.GetGroupName()))
			sb.WriteString(fmt.Sprintf(" └─ Active Workers: %d\n\n", group.GetWorkerInfo().GetCpuCoreCount()))
		}
		resultData = append(resultData, sb.String())
	} else {
		groupName := task.CmdArgs.GroupName
		sb.WriteString(fmt.Sprintf("Group: %s\n", groupName))
		sb.WriteString(fmt.Sprintf(" └─ Total Members: %d worker(s)\n\n", len(groupResponse)))

		var tableBuf bytes.Buffer
		w := tabwriter.NewWriter(&tableBuf, 0, 0, 3, ' ', 0)

		_, err := fmt.Fprintln(w, "   INDEX\tWORKER NAME\tWORKER ID")
		if err != nil {
			return nil
		}
		if len(groupResponse) == 0 {
			_, err := fmt.Fprintln(w, "   (No workers registered in this group)")
			if err != nil {
				return nil
			}
		}

		for index, data := range groupResponse {
			worker := data.GetWorkerInfo()
			workerId := worker.GetWorkerId()
			workerName := worker.GetWorkerName()

			_, err := fmt.Fprintf(w, "   #%d\t%s\t%s\n", index, workerName, workerId)
			if err != nil {
				return nil
			}
		}

		err = w.Flush()
		if err != nil {
			return nil
		}
		sb.Write(tableBuf.Bytes())
		resultData = append(resultData, sb.String())
	}
	return resultData
}
