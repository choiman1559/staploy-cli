package deploy

import (
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
	"staploy-cli/app/tasks/registry"
)

type DeployCmdTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.DeployCmd]
}

func (task *DeployCmdTask) MainCmd() error {
	hasPkgOnServer := task.checkPkgOnLocalServer()
	if !hasPkgOnServer && !task.CmdArgs.LocalOnly {
		logger.Tip("App %s not available on local server, trying pulling from remote repository...", task.CmdArgs.AppName)
		t := &registry.RegistryPullTask{}

		t.Init(task.DefaultArgs, cmds.RegistryPullCmd{
			AppName: task.CmdArgs.AppName,
			Version: task.CmdTask.CmdArgs.Version,
		}, proto.TaskGroup_TASK_REGISTRY)
		err := t.MainCmd()

		if err != nil {
			return err
		}
	}

	workers, err := task.ParseWorkers(false, task.CmdArgs.WorkerId...)
	if err != nil {
		return err
	}

	for _, worker := range workers {
		pushTask := &PushCmdTask{}
		pushTask.Init(task.DefaultArgs, cmds.PushCmd{
			AppName: task.CmdArgs.AppName,
			Version: task.CmdTask.CmdArgs.Version,
		}, proto.TaskGroup_TASK_DEPLOY)

		pushTask.AppInfo = &proto.AppInfoFetch{
			App: &proto.AppInfo{
				AppName: task.CmdArgs.AppName,
			},
		}

		if task.CmdArgs.Version != "" {
			pushTask.CmdArgs.Version = logger.TrimVersion(pushTask.CmdArgs.Version)
			pushTask.AppInfo.AppVersion = append(pushTask.AppInfo.GetAppVersion(), &proto.Version{VersionName: task.CmdArgs.Version})
		}

		version, err := pushTask.requestPush(worker)
		if err != nil {
			logger.Error("Failed to push app %s to worker %s, %s", task.CmdArgs.AppName, worker, err)
			continue
		}

		setTask := &SetCmdTask{}
		setTask.Init(task.DefaultArgs, cmds.SetCmd{
			AppName:  task.CmdArgs.AppName,
			Version:  version,
			WorkerId: []string{worker},
		}, proto.TaskGroup_TASK_DEPLOY)

		err = setTask.MainCmd()
		if err != nil {
			logger.Error("Failed to set app %s (v%s) to worker %s, Error: %s", task.CmdArgs.AppName, version, worker, err)
		}
	}

	return nil
}

func (task *DeployCmdTask) checkPkgOnLocalServer() bool {
	requestPacket := task.CreateDefPacket()
	requestPacket.TaskGroup = proto.TaskGroup_TASK_MANAGE_APPS
	requestPacket.TaskType = &proto.RequestPacket_AppsTaskType{AppsTaskType: proto.TaskAppsTypes_TYPE_APP_LISTS}

	if task.CmdArgs.AppName != "" {
		requestPacket.AppInfoFetch = append(requestPacket.GetAppInfoFetch(), &proto.AppInfoFetch{
			App: &proto.AppInfo{AppName: task.CmdArgs.AppName},
		})
	}

	response, err := task.PostRequest(requestPacket)
	if err != nil {
		return false
	}

	if response.GetStatus() != consts.StatusOK {
		return false
	}

	installedApp := response.GetWorkerResponse()[0].GetWorkerInfo().InstalledApp
	for _, installedInfo := range installedApp {
		if installedInfo.GetApp().GetAppName() == task.CmdArgs.AppName {
			return true
		}
	}
	return false
}
