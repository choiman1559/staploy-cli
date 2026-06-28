package user

import (
	"fmt"
	"os"
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
	"text/tabwriter"
	"time"

	protos "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/known/anypb"
)

type UserAuditTask struct {
	cmds.CmdTaskInterface
	cmds.CmdTask[cmds.UserAuditCmd]
}

var auditUnmarshalOpts = protos.UnmarshalOptions{
	Resolver:     protoregistry.GlobalTypes,
	AllowPartial: true,
}

func (task *UserAuditTask) MainCmd() error {
	requestPacket := task.CreateDefPacket()

	requestPacket.TaskType = &proto.RequestPacket_UserTaskType{UserTaskType: &proto.UserRequestPacket{
		UserTaskTypes: proto.TaskUserTypes_TYPE_USER_AUDIT,
	}}

	response, err := task.PostRequest(requestPacket)
	if err != nil {
		return err
	}

	if response.GetStatus() == consts.StatusOK {
		userResponse := response.GetUserResponse()
		if userResponse != nil {
			PrintAuditTable(userResponse.GetAuditData())
		} else {
			logger.Info("[+] Login successes but no audit data records found.")
		}
	} else {
		return fmt.Errorf("%s", response.GetErrorCause())
	}
	return nil
}

func resolvePacketDetails(audit *proto.AuditLogData) string {
	if audit.Response != nil && audit.Response.MessageIs((*proto.ResponsePacket)(nil)) {
		msg, err := anypb.UnmarshalNew(audit.Response, auditUnmarshalOpts)
		if err == nil {
			res := msg.(*proto.ResponsePacket)
			if res.GetStatus() == "error" {
				return res.GetErrorCause()
			}
		}
	}

	if audit.Request != nil && audit.Request.MessageIs((*proto.RequestPacket)(nil)) {
		msg, err := anypb.UnmarshalNew(audit.Request, auditUnmarshalOpts)
		if err == nil {
			req := msg.(*proto.RequestPacket)

			if req.GetUserTaskType() != nil {
				userTask := req.GetUserTaskType()
				if audit.Response != nil && audit.Response.MessageIs((*proto.ResponsePacket)(nil)) {
					if m, e := anypb.UnmarshalNew(audit.Response, auditUnmarshalOpts); e == nil {
						res := m.(*proto.ResponsePacket)
						if res.GetExtraData() != "" {
							return fmt.Sprintf("User Task: %s -> %s", userTask.GetUserTaskTypes().String(), res.GetExtraData())
						}
					}
				}
				return fmt.Sprintf("User Task: %s", userTask.GetUserTaskTypes().String())
			}

			if req.GetGroupTaskType() != nil {
				return fmt.Sprintf("Group Task: %s", req.GetGroupTaskType().GetGroupTaskTypes().String())
			}

			if req.GetNodeTaskType() != proto.TaskNodeTypes_TYPE_NODE_NONE {
				nodeTask := req.GetNodeTaskType()
				shortWorkerID := ""
				if req.GetWorker() != nil && len(req.GetWorker()) > 0 && req.GetWorker()[0].GetWorkerId() != "" {
					wid := req.GetWorker()[0].GetWorkerId()
					if len(wid) > 8 {
						shortWorkerID = " (Worker: " + wid[:8] + ")"
					} else {
						shortWorkerID = " (Worker: " + wid + ")"
					}
				}

				switch nodeTask {
				case proto.TaskNodeTypes_TYPE_NODE_EXECUTE_SHELL:
					return fmt.Sprintf("[Node Execution] Remote Bash Shell: '%s' -> Fin: successful", req.GetExtraData())
				case proto.TaskNodeTypes_TYPE_NODE_DISCONN_WORKER:
					return "[Node Topology] Terminated Connection for Worker" + shortWorkerID
				default:
					return fmt.Sprintf("[Node Action] %s", nodeTask.String())
				}
			}

			if req.GetAppsTaskType() != proto.TaskAppsTypes_TYPE_APP_NONE {
				appTask := req.GetAppsTaskType()
				appNameInfo := ""
				if req.GetExtraData() != "" {
					appNameInfo = " -> Asset ID: " + req.GetExtraData()
				}

				switch appTask {
				case proto.TaskAppsTypes_TYPE_APP_REGISTER:
					return "[App Configuration] Registered New Asset Definition" + appNameInfo
				case proto.TaskAppsTypes_TYPE_APP_DELETE:
					return "[App Action] Purged Application Asset Permanent Block" + appNameInfo
				case proto.TaskAppsTypes_TYPE_APP_PKG_PARSE:
					return "[App Package] Parsed Tarball Configuration Spec" + appNameInfo
				case proto.TaskAppsTypes_TYPE_APP_PKG_CREATE:
					return "[App Package] Created Deployment Tarball Bundle -> SHA256 Match"
				default:
					return fmt.Sprintf("[App Action] %s", appTask.String())
				}
			}
		}
	}

	return "No discrete parameters"
}

func PrintAuditTable(logs []*proto.AuditLogData) {
	total := len(logs)
	logger.Info("Fetched %d infrastructure audit log %s successfully.\n",
		total, map[bool]string{true: "trails", false: "trail"}[total > 1])

	if total == 0 {
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	_, err := fmt.Fprintln(w, "TIMESTAMP\tOPERATOR\tACTION\tSTATUS\tDETAILS / ERROR CAUSE")
	if err != nil {
		return
	}

	for _, audit := range logs {
		if audit == nil {
			continue
		}

		tm := time.UnixMilli(int64(audit.GetTimestamp()))
		timeStr := tm.Format("2006-01-02 15:04:05")

		statusStr := "OK"
		if audit.Response != nil && audit.Response.MessageIs((*proto.ResponsePacket)(nil)) {
			if m, err := anypb.UnmarshalNew(audit.Response, auditUnmarshalOpts); err == nil {
				if m.(*proto.ResponsePacket).GetStatus() == "error" {
					statusStr = "ERROR"
				}
			}
		}

		operator := audit.GetOperator()
		if operator == "$NO_USER" {
			operator = "(Anonymous)"
		}

		detailStr := resolvePacketDetails(audit)
		_, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			timeStr,
			operator,
			audit.GetAction().String(),
			statusStr,
			detailStr,
		)
		if err != nil {
			return
		}
	}

	err = w.Flush()
	if err != nil {
		return
	}
	fmt.Println()
}
