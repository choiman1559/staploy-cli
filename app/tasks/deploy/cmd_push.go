package deploy

import (
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
	"staploy-cli/app/logger"
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

	workers, err := task.ParseWorkers(false, task.CmdArgs.WorkerId...)
	if err != nil {
		return err
	}

	for _, workerId := range workers {
		request := task.CreateDefPacket(workerId)
		request.TaskType = &proto.RequestPacket_DeployTaskType{DeployTaskType: proto.TaskDeployTypes_TYPE_DEPLOY_PUSH_VERSION}
		request.AppInfoFetch = append(request.AppInfoFetch, appInfo)

		response, err := task.PostRequest(request)
		if err != nil {
			return err
		}

		if response.GetStatus() == consts.StatusOK && response.GetWorkerResponse()[0].GetTaskResult().GetResultSuccessful() {
			logger.Info("Pushed package: \"%s\" (%s) at worker %s", task.CmdArgs.AppName, logger.VersionNamePrefix(task.CmdArgs.Version), workerId)
		} else {
			logger.Error("Failed to push \"%s\" (%s) at worker %s", task.CmdArgs.AppName, logger.VersionNamePrefix(task.CmdArgs.Version), workerId)
			if len(response.GetWorkerResponse()) > 0 {
				logger.Error("Error cause is: %s", response.GetWorkerResponse()[0].GetTaskResult().GetErrorMessage())
			}
		}
	}
	return nil
}
