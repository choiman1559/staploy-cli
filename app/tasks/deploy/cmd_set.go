package deploy

import (
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
)

type SetCmdTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.SetCmd]
}

func (task *SetCmdTask) MainCmd() error {
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
		request.TaskType = &proto.RequestPacket_DeployTaskType{DeployTaskType: proto.TaskDeployTypes_TYPE_DEPLOY_SET_VERSION}
		request.AppInfoFetch = append(request.AppInfoFetch, appInfo)

		response, err := task.PostRequest(request)
		if err != nil {
			return err
		}

		if response.GetStatus() == consts.StatusOK {
			if response.GetWorkerResponse()[0].GetTaskResult().GetResultSuccessful() {
				if versionSpecified {
					logger.Info("Setting triggered package %s (%s) at worker %s", task.CmdArgs.AppName, task.CmdArgs.Version, workerId)
				} else {
					logger.Tip("Untriggered package %s at worker %s", task.CmdArgs.AppName, workerId)
				}
			} else {
				logger.Error("error: %v", response.GetWorkerResponse()[0].GetTaskResult().GetErrorMessage())
			}
		} else {
			logger.Error("failed to setting trigger for %s at worker %s", task.CmdArgs.AppName, workerId)
		}
	}
	return nil
}
