package deploy

import (
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
	"sync"
)

type PushCmdTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.PushCmd]

	AppInfo *proto.AppInfoFetch
}

func (task *PushCmdTask) MainCmd() error {
	task.AppInfo = &proto.AppInfoFetch{
		App: &proto.AppInfo{
			AppName: task.CmdArgs.AppName,
		},
	}

	if task.CmdArgs.Version != "" {
		task.AppInfo.AppVersion = append(task.AppInfo.GetAppVersion(), &proto.Version{VersionName: task.CmdArgs.Version})
	}

	workers, err := task.ParseWorkers(false, task.CmdArgs.WorkerId...)
	if err != nil {
		return err
	}

	if task.CmdArgs.MaxThread < 1 {
		for _, workerId := range workers {
			err := task.requestPush(workerId)
			if err != nil {
				return err
			}
		}
		return nil
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(workers))
	sem := make(chan struct{}, task.CmdArgs.MaxThread)

	for _, workerId := range workers {
		wg.Add(1)
		sem <- struct{}{}

		go func(wId string) {
			defer wg.Done()
			defer func() { <-sem }()

			err := task.requestPush(wId)
			if err != nil {
				errChan <- err
			}
		}(workerId)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return <-errChan
	}
	return nil
}

func (task *PushCmdTask) requestPush(workerId string) error {
	request := task.CreateDefPacket(workerId)
	request.TaskType = &proto.RequestPacket_DeployTaskType{DeployTaskType: proto.TaskDeployTypes_TYPE_DEPLOY_PUSH_VERSION}
	request.AppInfoFetch = append(request.AppInfoFetch, task.AppInfo)

	response, err := task.PostRequest(request)
	if err != nil {
		return err
	}

	if response.GetStatus() == consts.StatusOK && response.GetWorkerResponse()[0].GetTaskResult().GetResultSuccessful() {
		logger.Info("Pushed package: \"%s\" (%s) at worker %s", task.CmdArgs.AppName, logger.VersionNamePrefix(response.GetExtraData()), workerId)
	} else {
		if task.CmdArgs.Version != "" {
			logger.Error("Failed to push \"%s\" (%s) at worker %s", task.CmdArgs.AppName, logger.VersionNamePrefix(task.CmdArgs.Version), workerId)
		} else {
			logger.Error("Failed to push \"%s\" at worker %s", task.CmdArgs.AppName, workerId)
		}

		if len(response.GetWorkerResponse()) > 0 {
			logger.Error("Error cause is: %s", response.GetWorkerResponse()[0].GetTaskResult().GetErrorMessage())
		}
	}
	return nil
}
