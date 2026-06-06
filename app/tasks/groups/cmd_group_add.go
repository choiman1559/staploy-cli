package groups

import (
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
)

type GroupAddTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.GroupAddCmd]
}

func (task *GroupAddTask) MainCmd() error {
	requestPacket := task.CreateDefPacket()
	workers, err := task.ParseWorkers(true, task.CmdArgs.WorkerId...)
	if err != nil {
		return err
	}

	requestPacket.TaskType = &proto.RequestPacket_GroupTaskType{GroupTaskType: &proto.GroupRequestPacket{
		GroupTaskTypes: proto.TaskGroupTypes_TYPE_GROUP_ADD_WORKER,
		GroupName:      &task.CmdArgs.GroupName,
		Names:          workers,
	}}

	response, err := task.PostRequest(requestPacket)
	if err != nil {
		return err
	}

	if response.GetStatus() == consts.StatusOK {
		if len(response.GetGroupResponse()) < 1 || response.GetGroupResponse()[0].GetGroupName() != task.CmdArgs.GroupName {
			logger.Error("Failed to add workers to group %s", task.CmdArgs.GroupName, task.CmdArgs.GroupName)
		} else {
			logger.Info("Added workers to group %s", task.CmdArgs.GroupName)
		}
	} else {
		logger.Error("Failed to add workers to group %s, cause: %s", task.CmdArgs.GroupName, response.GetExtraData())
	}
	return nil
}
