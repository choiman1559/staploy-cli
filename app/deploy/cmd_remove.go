package deploy

import (
	"fmt"
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
	if task.CmdArgs.AutoRemove {
		return task.RemoveAuto()
	}
	return task.RemoveSpecified()
}

func (task *RemoveCmdTask) RemoveAuto() error {
	for _, workerId := range task.CmdArgs.WorkerId {
		packet := task.CreateDefPacket(workerId)
		packet.TaskGroup = proto.TaskGroup_TASK_MANAGE_NODE
		packet.TaskType = &proto.RequestPacket_NodeTaskType{NodeTaskType: proto.TaskNodeTypes_TYPE_NODE_REQ_APP_INFO}

		responsePacket, err := task.PostRequest(packet)
		if err != nil {
			return err
		}

		if len(responsePacket.GetWorkerResponse()) < 1 {
			fmt.Printf("fetch response for worker %s is empty. check worker-id is correct\n", workerId)
			continue
		}

		removePacket := task.CreateDefPacketIdOnly(packet.GetWorker()[0].GetWorkerId())
		removePacket.TaskType = &proto.RequestPacket_DeployTaskType{DeployTaskType: proto.TaskDeployTypes_TYPE_DEPLOY_DEL_VERSION}

		for _, v := range responsePacket.GetWorkerResponse()[0].GetWorkerInfo().GetInstalledApp() {
			if task.CmdArgs.AppName != "" && task.CmdArgs.AppName != v.GetApp().GetAppName() {
				continue
			}

			appInfo := &proto.AppInfoFetch{
				App: &proto.AppInfo{
					AppName: v.GetApp().GetAppName(),
				},
			}

			if v.GetCurrentVersion().GetVersionName() != "" && len(v.GetAvailableVersion()) > 0 {
				for _, version := range v.GetAvailableVersion() {
					if version.GetVersionName() != v.GetCurrentVersion().GetVersionName() {
						appInfo.AppVersion = append(appInfo.GetAppVersion(), version)
						log.Printf("removing unused package %s (%s) at worker %s", v.GetApp().GetAppName(), version.GetVersionName(), workerId)
					}
				}
				removePacket.AppInfoFetch = append(removePacket.AppInfoFetch, appInfo)
			}
		}

		removeResponse, err := task.PostRequest(removePacket)
		if err != nil {
			return err
		}

		if len(removeResponse.GetWorkerResponse()) < 1 {
			fmt.Printf("fetch response for worker %s is empty. check worker-id is correct\n", workerId)
		} else if removeResponse.GetWorkerResponse()[0].GetTaskResult().GetResultSuccessful() {
			log.Println("cleaned all unused packages on worker " + workerId)
		} else {
			return fmt.Errorf("error on worker " + workerId + " => " + removeResponse.GetWorkerResponse()[0].GetTaskResult().GetErrorMessage())
		}
	}
	return nil
}

func (task *RemoveCmdTask) RemoveSpecified() error {
	if task.CmdArgs.AppName == "" {
		return fmt.Errorf("--app-name must be specified")
	}

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
