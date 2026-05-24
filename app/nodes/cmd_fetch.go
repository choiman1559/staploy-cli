package nodes

import (
	"fmt"
	"staploy-cli/app/cmds"
	"staploy-cli/app/proto"
)

type FetchCmdTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.FetchCmd]
}

func (a *FetchCmdTask) MainCmd() error {
	packet := a.CreateDefPacket(a.CmdArgs.WorkerId)
	packet.TaskType = &proto.RequestPacket_NodeTaskType{NodeTaskType: proto.TaskNodeTypes_TYPE_NODE_REQ_APP_INFO}

	if a.CmdArgs.AppName != "" {
		appInfo := proto.AppInfoFetch{App: &proto.AppInfo{AppName: a.CmdArgs.AppName}}
		if len(a.CmdArgs.VersionName) > 0 {
			for _, v := range a.CmdArgs.VersionName {
				appInfo.AppVersion = append(appInfo.GetAppVersion(), &proto.Version{VersionName: v})
			}
		}

		packet.AppInfoFetch = append(packet.GetAppInfoFetch(), &appInfo)
	}

	responsePacket, err := a.PostRequest(packet)
	if err != nil {
		return err
	}

	for _, str := range AppInfoFormatter(responsePacket.GetWorkerResponse()[0].GetWorkerInfo(), a.CmdArgs.Detail) {
		fmt.Print(str)
	}
	return nil
}

func AppInfoFormatter(workerInfo *proto.WorkerInfo, detail bool) []string {
	var workerData []string

	for _, data := range workerInfo.InstalledApp {
		workerData = append(workerData, fmt.Sprintf("App %s\n", data.GetApp().GetAppName()))
		workerData = append(workerData, fmt.Sprintf("\tDescription: %s\n", safeBlankStr(data.GetApp().GetAppDescription())))

		if data.CurrentVersion != nil {
			workerData = append(workerData, fmt.Sprintf("\tCurrent Version: %s\n", data.GetCurrentVersion().GetVersionName()))
		}

		workerData = append(workerData, fmt.Sprintf("\tAvailable Versions:\n"))
		for _, version := range data.GetAvailableVersion() {
			workerData = append(workerData, fmt.Sprintf("\t\t\tVersion: %s\n", version.GetVersionName()))

			if detail {
				if version.GetLibVersion() != "" {
					workerData = append(workerData, fmt.Sprintf("\t\t\t\tUsed Library: %s\n", safeBlankStr(version.GetLibVersion())))
				}

				if len(version.GetEntryBinaries()) > 0 {
					workerData = append(workerData, fmt.Sprintf("\t\t\t\tReported entry binaries:\n"))
					for i, execs := range version.GetEntryBinaries() {
						workerData = append(workerData, fmt.Sprintf("\t\t\t\t\tBinary #%d: %s (Hash: %s)\n", i, execs.GetName(), execs.GetHash()))
					}
				}
			}
		}
	}

	return workerData
}

func safeBlankStr(str string) string {
	if str == "" {
		return "(none)"
	}
	return str
}
