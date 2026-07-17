package groups

import (
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
)

type GroupCreateTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.GroupCreateCmd]
}

func (task *GroupCreateTask) MainCmd() error {
	requestPacket := task.CreateDefPacket()
	requestPacket.TaskType = &proto.RequestPacket_GroupTaskType{GroupTaskType: &proto.GroupRequestPacket{
		GroupTaskTypes: proto.TaskGroupTypes_TYPE_GROUP_CREATE,
		GroupName:      &task.CmdArgs.GroupName,
	}}

	response, err := task.PostRequest(requestPacket)
	if err != nil {
		return err
	}

	if response.GetStatus() == consts.StatusOK {
		if len(response.GetGroupResponse()) < 1 || response.GetGroupResponse()[0].GetGroupName() != task.CmdArgs.GroupName {
			logger.Error("Failed to create group %s, group %s already exists", task.CmdArgs.GroupName, task.CmdArgs.GroupName)
		} else {
			logger.Info("Created group %s", task.CmdArgs.GroupName)
		}
	} else {
		logger.Error("Failed to create group %s, cause: %s", task.CmdArgs.GroupName, CollectErrorMessage(response))
	}
	return nil
}
