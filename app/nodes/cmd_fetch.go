package nodes

import (
	"bytes"
	"fmt"
	"staploy-cli/app/cmds"
	"staploy-cli/app/proto"
	"strings"
	"text/tabwriter"
)

type FetchCmdTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.FetchCmd]
}

func (a *FetchCmdTask) MainCmd() error {
	packet := a.CreateDefPacket(a.CmdArgs.WorkerId)
	packet.TaskType = &proto.RequestPacket_NodeTaskType{NodeTaskType: proto.TaskNodeTypes_TYPE_NODE_REQ_APP_INFO}

	var isDetail = a.CmdArgs.Detail
	if a.CmdArgs.AppName != "" {
		appInfo := proto.AppInfoFetch{App: &proto.AppInfo{AppName: a.CmdArgs.AppName}}
		if len(a.CmdArgs.VersionName) > 0 {
			for _, v := range a.CmdArgs.VersionName {
				appInfo.AppVersion = append(appInfo.GetAppVersion(), &proto.Version{VersionName: v})
			}
		}

		packet.AppInfoFetch = append(packet.GetAppInfoFetch(), &appInfo)
	} else {
		isDetail = false
	}

	responsePacket, err := a.PostRequest(packet)
	if err != nil {
		return err
	}

	if len(responsePacket.GetWorkerResponse()) < 1 {
		return fmt.Errorf("fetch response is empty. check worker-id is correct")
	}

	for _, str := range AppInfoFormatter(responsePacket.GetWorkerResponse()[0].GetWorkerInfo(), isDetail) {
		fmt.Print(str)
	}
	return nil
}

func AppInfoFormatter(workerInfo *proto.WorkerInfo, detail bool) []string {
	var workerData []string

	for _, data := range workerInfo.InstalledApp {
		var sb strings.Builder
		appName := data.GetApp().GetAppName()
		appDesc := safeBlankStr(data.GetApp().GetAppDescription())

		sb.WriteString(fmt.Sprintf("App: %s\n", appName))
		sb.WriteString(fmt.Sprintf(" └─ Description: %s\n", appDesc))

		currentVer := "(none)"
		if data.CurrentVersion != nil && data.CurrentVersion.VersionName != "" {
			currentVer = data.GetCurrentVersion().GetVersionName()
		}
		sb.WriteString(fmt.Sprintf(" └─ Current Version: %s\n\n", currentVer))

		var tableBuf bytes.Buffer
		w := tabwriter.NewWriter(&tableBuf, 0, 0, 3, ' ', 0)

		if detail {
			_, err := fmt.Fprintln(w, "   VERSION\tSTATUS\tLIB RUNTIME\tENTRY BINARIES (NAME / HASH)")
			if err != nil {
				return nil
			}
		} else {
			_, err := fmt.Fprintln(w, "   VERSION\tSTATUS")
			if err != nil {
				return nil
			}
		}

		availableVersions := data.GetAvailableVersion()
		if len(availableVersions) == 0 {
			_, err := fmt.Fprintln(w, "   (No available versions reported from server)")
			if err != nil {
				return nil
			}
		}

		for _, version := range availableVersions {
			verName := version.GetVersionName()

			status := "-"
			if data.CurrentVersion != nil && verName == currentVer {
				status = "[Active]"
			}

			if detail {
				libVer := safeBlankStr(version.GetLibVersion())
				execs := version.GetEntryBinaries()

				if len(execs) == 0 {
					_, err := fmt.Fprintf(w, "   %s\t%s\t%s\t%s\n", verName, status, libVer, "(none)")
					if err != nil {
						return nil
					}
				} else {
					firstBin := fmt.Sprintf("%s (%s)", execs[0].GetName(), shortHash(execs[0].GetHash()))
					_, err := fmt.Fprintf(w, "   %s\t%s\t%s\t%s\n", verName, status, libVer, firstBin)
					if err != nil {
						return nil
					}

					for i := 1; i < len(execs); i++ {
						nextBin := fmt.Sprintf("%s (%s)", execs[i].GetName(), shortHash(execs[i].GetHash()))
						_, err := fmt.Fprintf(w, "   \t\t\t%s\n", nextBin)
						if err != nil {
							return nil
						}
					}
				}
			} else {
				_, err := fmt.Fprintf(w, "   %s\t%s\n", verName, status)
				if err != nil {
					return nil
				}
			}
		}

		err := w.Flush()
		if err != nil {
			return nil
		}
		sb.Write(tableBuf.Bytes())
		sb.WriteString("\n")

		workerData = append(workerData, sb.String())
	}
	return workerData
}

func shortHash(hash string) string {
	if len(hash) > 10 {
		return hash[:10] + "..."
	}
	return hash
}

func safeBlankStr(str string) string {
	if str == "" {
		return "(none)"
	}
	return str
}
