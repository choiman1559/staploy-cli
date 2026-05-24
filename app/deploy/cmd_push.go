package deploy

import (
	"log"
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
	"staploy-cli/app/proto"
)

type PushCmdTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.PushCmd]
}

func (task *PushCmdTask) MainCmd() error {
	appInfo := &proto.AppInfoFetch{
		App: &proto.AppInfo{
			AppName: task.CmdArgs.AppName,
		},
	}
	appInfo.AppVersion = append(appInfo.GetAppVersion(), &proto.Version{VersionName: task.CmdArgs.Version})

	for _, workerId := range task.CmdArgs.WorkerId {
		request := task.CreateDefPacket(workerId)
		request.TaskType = &proto.RequestPacket_DeployTaskType{DeployTaskType: proto.TaskDeployTypes_TYPE_DEPLOY_PUSH_VERSION}
		request.AppInfoFetch = append(request.AppInfoFetch, appInfo)

		response, err := task.PostRequest(request)
		if err != nil {
			return err
		}

		if response.GetStatus() == consts.StatusOK {
			log.Printf("pushed package %s (%s) at worker %s", task.CmdArgs.AppName, task.CmdArgs.Version, workerId)
		} else {
			log.Printf("failed to push %s (%s) at worker %s", task.CmdArgs.AppName, task.CmdArgs.Version, workerId)
		}
	}
	return nil
}
