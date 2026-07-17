package registry

import (
	"fmt"
	"os"
	"staploy-cli/app/cmds"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
	"text/tabwriter"
)

type RegistryListRepoTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.RegistryListRepoCmd]
}

func (task *RegistryListRepoTask) MainCmd() error {
	requestPacket := task.CreateDefPacket()
	registryRequest := &proto.RegistryRequestPacket{
		TaskType: proto.TaskRegistryTypes_LOCAL_LIST_REPOSITORY,
	}

	requestPacket.TaskType = &proto.RequestPacket_RegistryTaskType{RegistryTaskType: registryRequest}
	response, err := task.PostRequest(requestPacket)
	if err != nil {
		return err
	}

	printRepoLists(response.GetRegistryResponse().RepositoryUrl)
	return nil
}

func printRepoLists(repoUrls []string) {
	total := len(repoUrls)

	if total == 0 {
		logger.Error("No registered repository sources found.")
		logger.Tip("Tip: use \"staploy-cli repo add <URL>\" to register a new repository.")
		return
	}

	logger.Info("Found %d registered repository %s.\n", total, map[bool]string{true: "sources", false: "source"}[total > 1])
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)

	_, err2 := fmt.Fprintln(w, "INDEX\tREPOSITORY URL")
	if err2 != nil {
		return
	}

	for i, url := range repoUrls {
		_, err := fmt.Fprintf(w, "#%d\t%s\n", i+1, url)
		if err != nil {
			return
		}
	}

	err := w.Flush()
	if err != nil {
		return
	}
	fmt.Println()
}
