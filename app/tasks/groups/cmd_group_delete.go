package groups

import (
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
)

type GroupDeleteTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.GroupDeleteCmd]
}

func (task *GroupDeleteTask) MainCmd() error {
	requestPacket := task.CreateDefPacketIdOnly()
	requestPacket.TaskType = &proto.RequestPacket_GroupTaskType{GroupTaskType: &proto.GroupRequestPacket{
		GroupTaskTypes: proto.TaskGroupTypes_TYPE_GROUP_DELETE,
		GroupName:      &task.CmdArgs.GroupName,
	}}

	response, err := task.PostRequest(requestPacket)
	if err != nil {
		return err
	}

	if response.GetStatus() == consts.StatusOK {
		if len(response.GetGroupResponse()) < 1 || response.GetGroupResponse()[0].GetGroupName() != task.CmdArgs.GroupName {
			logger.Error("Group %s not exists", task.CmdArgs.GroupName)
		} else {
			logger.Tip("Deleted group %s", task.CmdArgs.GroupName)
		}
	} else {
		logger.Error("Failed to delete group %s, cause: %s", task.CmdArgs.GroupName, response.GetExtraData())
	}
	return nil
}
