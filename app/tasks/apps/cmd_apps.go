package apps

import (
	"fmt"
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
)

type AppsCmdTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.AppsCmd]
}

func (t *AppsCmdTask) MainCmd() error {
	requestPacket := t.CreateDefPacket()
	requestPacket.TaskType = &proto.RequestPacket_AppsTaskType{AppsTaskType: proto.TaskAppsTypes_TYPE_APP_LISTS}

	if t.CmdArgs.AppName != "" {
		requestPacket.AppInfoFetch = append(requestPacket.GetAppInfoFetch(), &proto.AppInfoFetch{
			App: &proto.AppInfo{AppName: t.CmdArgs.AppName},
		})
	}

	response, err := t.PostRequest(requestPacket)
	if err != nil {
		return err
	}

	if response.GetStatus() != consts.StatusOK {
		return fmt.Errorf(response.GetErrorCause())
	}

	installedApp := response.GetWorkerResponse()[0].GetWorkerInfo().InstalledApp
	for _, installedInfo := range installedApp {
		for _, str := range formatAppInfo(installedInfo) {
			fmt.Printf(str)
		}
		fmt.Println()
	}
	return nil
}

func formatAppInfo(installedInfo *proto.InstalledAppInfo) []string {
	var lines []string

	lines = append(lines, fmt.Sprintf("* App: %s\n", installedInfo.GetApp().GetAppName()))
	lines = append(lines, fmt.Sprintf(" ├─ Description: %s\n", logger.SafeBlankStr(installedInfo.GetApp().GetAppDescription())))
	lines = append(lines, " └─ Available Versions:\n")

	if len(installedInfo.AvailableVersion) < 1 {
		lines = append(lines, fmt.Sprintf("    (No uploaded package version yet)\n"))
	} else {
		for _, ver := range installedInfo.AvailableVersion {
			lines = append(lines, fmt.Sprintf("    %s (supports %v)\n", logger.VersionNamePrefix(ver.VersionName), ver.GetSupportedArch()))
		}
	}
	return lines
}
