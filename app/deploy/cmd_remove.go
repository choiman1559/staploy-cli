package deploy

import (
	"log"
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
	"staploy-cli/app/proto"
)

type RemoveCmdTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.RemoveCmd]
}

func (task *RemoveCmdTask) MainCmd() error {
	appInfo := &proto.AppInfoFetch{
		App: &proto.AppInfo{
			AppName: task.CmdArgs.AppName,
		},
	}

	versionSpecified := false
	if task.CmdArgs.Version != "" {
		versionSpecified = true
		appInfo.AppVersion = append(appInfo.GetAppVersion(), &proto.Version{VersionName: task.CmdArgs.Version})
	}

	for _, workerId := range task.CmdArgs.WorkerId {
		request := task.CreateDefPacket(workerId)
		request.TaskType = &proto.RequestPacket_DeployTaskType{DeployTaskType: proto.TaskDeployTypes_TYPE_DEPLOY_DEL_VERSION}
		request.AppInfoFetch = append(request.AppInfoFetch, appInfo)

		response, err := task.PostRequest(request)
		if err != nil {
			return err
		}

		if response.GetStatus() == consts.StatusOK {
			if versionSpecified {
				log.Printf("removed package %s (%s) at worker %s", task.CmdArgs.AppName, task.CmdArgs.Version, workerId)
			} else {
				log.Printf("removed all package %s at worker %s", task.CmdArgs.AppName, workerId)
			}
		} else {
			log.Printf("failed to remove %s at worker %s", task.CmdArgs.AppName, workerId)
		}
	}
	return nil
}
