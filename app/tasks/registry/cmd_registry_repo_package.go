package registry

import (
	"fmt"
	"os"
	"staploy-cli/app/cmds"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
	"text/tabwriter"
)

type RegistryListTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.RegistryListPackageCmd]
}

type AppVersionWithRepo struct {
	*proto.Version
	RepoUrl string
}

func (task *RegistryListTask) MainCmd() error {
	requestPacket := task.CreateDefPacket()

	registryRequest := &proto.RegistryRequestPacket{
		TaskType: proto.TaskRegistryTypes_LOCAL_PACKAGE_QUERY,
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

	allApps := make(map[string][]*AppVersionWithRepo)
	for i, appInfo := range response.GetRegistryResponse().GetAppInfo() {
		repoUrl := response.GetRegistryResponse().GetRepositoryUrl()[i]
		for _, version := range appInfo.GetAvailableVersion() {
			allApps[appInfo.GetApp().GetAppName()] = append(allApps[appInfo.GetApp().GetAppName()], &AppVersionWithRepo{
				Version: version,
				RepoUrl: repoUrl,
			})
		}
	}

	printPackageLists(allApps)
	return nil
}

func printPackageLists(lists map[string][]*AppVersionWithRepo) {
	totalApps := len(lists)

	if totalApps == 0 {
		logger.Error("No packages or repository indexing assets found.")
		logger.Tip("Tip: use \"staploy-cli registry update\" to refresh your package repository cache.")
		return
	}

	logger.Info("Found %d application asset %s registered.\n", totalApps, map[bool]string{true: "packages", false: "package"}[totalApps > 1])
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	_, err := fmt.Fprintln(w, "APPLICATION\tVERSION\tLIB\tARCH\tSOURCE REPOSITORY")
	if err != nil {
		return
	}

	for appName, versions := range lists {
		if len(versions) == 0 {
			continue
		} else if len(versions) == 1 {
			v := versions[0]
			libStr := v.GetLibVersion()
			archStr := fmt.Sprintf("%v", FormatArchArray(v.GetSupportedArch()))

			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				appName,
				v.GetVersionName(),
				map[bool]string{true: "any", false: libStr}[libStr == ""],
				map[bool]string{true: "unknown", false: archStr}[archStr == ""],
				v.RepoUrl,
			)
			continue
		}

		_, err2 := fmt.Fprintf(w, "%s\t\t\t\t\n", appName)
		if err2 != nil {
			return
		}

		for i, v := range versions {
			if v == nil || v.Version == nil {
				continue
			}

			libStr := v.GetLibVersion()
			archStr := FormatArchArray(v.GetSupportedArch())

			if libStr == "" {
				libStr = "any"
			}

			treePrefix := v.GetVersionName()
			if i == len(versions)-1 {
				treePrefix = "└── " + treePrefix
			} else {
				treePrefix = "├── " + treePrefix
			}

			_, err3 := fmt.Fprintf(w, "\t%s\t%s\t%s\t%s\n",
				treePrefix,
				libStr,
				archStr,
				v.RepoUrl,
			)
			if err3 != nil {
				return
			}
		}
	}

	err = w.Flush()
	if err != nil {
		return
	}
	fmt.Println()
}
