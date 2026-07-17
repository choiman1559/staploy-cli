package registry

import (
	"fmt"
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
)

type RegistryListLocalTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.RegistryListLocalCmd]
}

func (task *RegistryListLocalTask) MainCmd() error {
	task.OverrideConnType(consts.ConnTypeRegistry)
	requestPacket := task.CreateDefPacket()

	registryRequest := &proto.RegistryRequestPacket{
		TaskType: proto.TaskRegistryTypes_TASK_LIST,
	}

	if task.CmdArgs.AppName != "" {
		appInfo := &proto.AppInfoFetch{
			App: &proto.AppInfo{
				AppName: task.CmdArgs.AppName,
			},
		}
		registryRequest.AppInfo = appInfo
	}

	requestPacket.TaskType = &proto.RequestPacket_RegistryTaskType{RegistryTaskType: registryRequest}
	response, err := task.PostRequest(requestPacket)
	if err != nil {
		return err
	}

	if response.GetStatus() != consts.StatusOK {
		return fmt.Errorf(response.GetErrorCause())
	}
	installedApp := response.GetRegistryResponse().AppInfo

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
			lines = append(lines, fmt.Sprintf("    %s (supports %v)\n", logger.VersionNamePrefix(ver.VersionName), FormatArchArray(ver.GetSupportedArch())))
		}
	}
	return lines
}

func FormatArchArray(archArray []proto.CpuArch) []string {
	arches := make([]string, len(archArray))
	for i, arch := range archArray {
		if arch == proto.CpuArch_UNKNOWN {
			arches[i] = "share"
			continue
		}
		arches[i] = arch.String()
	}
	return arches
}
