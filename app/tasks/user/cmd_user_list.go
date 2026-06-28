package user

import (
	"fmt"
	"os"
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
	"strings"
	"text/tabwriter"
)

type UserListTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.UserListCmd]
}

func (task *UserListTask) MainCmd() error {
	requestPacket := task.CreateDefPacket()

	userParameter := &proto.UserRequestPacket{
		UserTaskTypes: proto.TaskUserTypes_TYPE_USER_LISTS,
	}

	if task.CmdArgs.UserName != "" {
		userParameter.UserLoginInfo = &proto.UserLoginInfo{
			UserName: task.CmdArgs.UserName,
		}
	}

	requestPacket.TaskType = &proto.RequestPacket_UserTaskType{UserTaskType: userParameter}

	response, err := task.PostRequest(requestPacket)
	if err != nil {
		return err
	}

	if response.GetStatus() == consts.StatusOK {
		PrintUserMetadataList(response.UserResponse.GetUserMetaDatas())
	} else {
		return fmt.Errorf("%s", response.GetErrorCause())
	}
	return nil
}

func ParsePermissionsList(perm uint32) string {
	if perm&uint32(proto.PermissionFlag_SYSTEM_ADMIN) != 0 {
		return "[SYSTEM_ADMIN]"
	}

	var active []string

	if perm&uint32(proto.PermissionFlag_APP_CREATE) != 0 {
		active = append(active, "APP_CREATE")
	}
	if perm&uint32(proto.PermissionFlag_APP_UPLOAD) != 0 {
		active = append(active, "APP_UPLOAD")
	}
	if perm&uint32(proto.PermissionFlag_APP_DELETE) != 0 {
		active = append(active, "APP_DELETE")
	}
	if perm&uint32(proto.PermissionFlag_NODE_BASH) != 0 {
		active = append(active, "NODE_BASH")
	}
	if perm&uint32(proto.PermissionFlag_NODE_PUSH) != 0 {
		active = append(active, "NODE_PUSH")
	}
	if perm&uint32(proto.PermissionFlag_NODE_SET) != 0 {
		active = append(active, "NODE_SET")
	}
	if perm&uint32(proto.PermissionFlag_NODE_REMOVE) != 0 {
		active = append(active, "NODE_REMOVE")
	}
	if perm&uint32(proto.PermissionFlag_QUERY_ENDPOINT) != 0 {
		active = append(active, "QUERY_ENDPOINT")
	}
	if perm&uint32(proto.PermissionFlag_NODE_REMOVE) != 0 {
		active = append(active, "NODE_DISCONN")
	}
	if perm&uint32(proto.PermissionFlag_GROUP_MANAGE) != 0 {
		active = append(active, "GROUP_MANAGE")
	}
	if perm&uint32(proto.PermissionFlag_USER_MANAGE) != 0 {
		active = append(active, "USER_MANAGE")
	}

	if len(active) == 0 {
		return "[NONE]"
	}

	return "[" + strings.Join(active, ", ") + "]"
}

func PrintUserMetadataList(userMetadataList []*proto.UserMetadata) {
	totalRecords := len(userMetadataList)
	logger.Info("Fetched %d user metadata %s successfully.\n",
		totalRecords, map[bool]string{true: "records", false: "record"}[totalRecords > 1])

	if totalRecords == 0 {
		fmt.Println("No user metadata found.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	_, err := fmt.Fprintln(w, "UUID\tUSERNAME\tVERSION\tROLE\tPERMISSIONS")
	if err != nil {
		return
	}

	for _, user := range userMetadataList {
		if user == nil {
			continue
		}

		shortUUID := user.GetUuid()
		if len(shortUUID) > 8 {
			shortUUID = shortUUID[:8]
		}

		role := "none"
		if user.RoleName != nil && user.GetRoleName() != "" {
			role = user.GetRoleName()
		}

		permStr := ParsePermissionsList(user.GetPermissions())

		_, err2 := fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\n",
			shortUUID,
			user.GetUserName(),
			user.GetVersion(),
			role,
			permStr,
		)
		if err2 != nil {
			return
		}
	}

	err = w.Flush()
	if err != nil {
		return
	}
}
