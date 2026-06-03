package groups

import (
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
)

type GroupRemoveTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.GroupRemoveCmd]
}

func (task *GroupRemoveTask) MainCmd() error {
	requestPacket := task.CreateDefPacketIdOnly()
	requestPacket.TaskType = &proto.RequestPacket_GroupTaskType{GroupTaskType: &proto.GroupRequestPacket{
		GroupTaskTypes: proto.TaskGroupTypes_TYPE_GROUP_REMOVE_WORKER,
		GroupName:      &task.CmdArgs.GroupName,
		Names:          task.CmdArgs.WorkerId,
	}}

	response, err := task.PostRequest(requestPacket)
	if err != nil {
		return err
	}

	if response.GetStatus() == consts.StatusOK {
		if len(response.GetGroupResponse()) < 1 || response.GetGroupResponse()[0].GetGroupName() != task.CmdArgs.GroupName {
			logger.Error("Failed to remove workers from group %s", task.CmdArgs.GroupName, task.CmdArgs.GroupName)
		} else {
			logger.Tip("Removed workers from group %s", task.CmdArgs.GroupName)
		}
	} else {
		logger.Error("Failed to remove workers from group %s, cause: %s", task.CmdArgs.GroupName, response.GetExtraData())
	}
	return nil
}
