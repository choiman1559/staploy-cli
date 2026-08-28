package build

import (
	"staploy-cli/app/cmds"
	"staploy-cli/app/logger"
	"staploy-cli/app/proto"
	"staploy-cli/app/tasks/deploy"
	"staploy-cli/app/tasks/nodes"
)

func (a *StaFileTask) processTarget(defArgs *cmds.DefaultArgs, targets []*Target) error {
	for _, target := range targets {
		logger.Task("Processing target \"%s\"", target.TargetName)
		logger.EnableTree()

		resolvedWorkers, err := a.ParseWorkerAlias(false, target.WorkerIds...)
		if err != nil {
			logger.DisableTree(true)
			logger.Error("Failed to parse workers, cause: %s", err.Error())
			continue
		}
		target.WorkerIds = resolvedWorkers

		if target.Where != nil {
			var filteredWorkers []string
			parsedWorkers, err := a.ParseWorkers(false, target.WorkerIds...)

			if err != nil {
				logger.DisableTree(true)
				logger.Error("Failed to parse worker, cause: %s", err.Error())
				continue
			}

			for _, worker := range parsedWorkers {
				workerInfo, err := a.askWorkerInfoData(worker)
				if err == nil {
					passedFilter, err := a.evalFilter(workerInfo, target.Where)
					if err != nil {
						logger.Warn("Failed to eval filter for worker \"%s\", cause: %s", worker, err.Error())
						continue
					}

					if passedFilter {
						if a.DefaultArgs.Verbose {
							logger.Tip("[DEBUG] worker \"%s\" passed filter", worker)
						}
						filteredWorkers = append(filteredWorkers, workerInfo.WorkerId)
					}
				} else {
					logger.Warn(err.Error())
					continue
				}
			}

			if len(filteredWorkers) > 0 {
				target.WorkerIds = filteredWorkers
			} else {
				logger.DisableTree(true)
				logger.Warn("No workers selected for target \"%s\", skipping.", target.TargetName)
				continue
			}
		}

		if target.Unset != nil {
			a.EvalTask(defArgs, &BashCommandTask{
				workerId:   target.WorkerIds,
				preDeploy:  target.Unset.PreDeploy,
				postDeploy: target.Unset.PostDeploy,
				deployTask: func() {
					err := a.resolveAppAlias(&target.Unset.AppName, nil)
					if err != nil {
						logger.Error("Cannot resolve app alias for unset: %q", target.Unset.AppName)
						return
					}

					for _, workerId := range target.WorkerIds {
						t := &deploy.UnsetCmdTask{}
						t.Init(*defArgs, cmds.UnsetCmd{
							WorkerId: []string{workerId}, AppName: target.Unset.AppName,
						}, proto.TaskGroup_TASK_DEPLOY)

						err := t.MainCmd()
						if err != nil {
							logger.Warn("Worker %s Executing untriggering: %s -> Failed for: %v", a.FormatWorker(workerId), target.Unset.AppName, err)
						}
					}
				},
			})
		}

		if target.Remove != nil {
			a.EvalTask(defArgs, &BashCommandTask{
				workerId:   target.WorkerIds,
				preDeploy:  target.Remove.PreDeploy,
				postDeploy: target.Remove.PostDeploy,
				deployTask: func() {
					err := a.resolveAppAlias(&target.Remove.AppName, &target.Remove.Version)
					if err != nil {
						logger.Error("Cannot resolve app alias for remove: %q", target.Remove.AppName)
						return
					}

					if target.Remove.Version == "" {
						logger.Error("Version not specified, Abort Remove")
						return
					}

					for _, workerId := range target.WorkerIds {
						t := &deploy.RemoveCmdTask{}
						t.Init(*defArgs, cmds.RemoveCmd{
							WorkerId: []string{workerId}, AppName: target.Remove.AppName, Version: target.Remove.Version, AutoRemove: target.Remove.AutoRemove,
						}, proto.TaskGroup_TASK_DEPLOY)

						err := t.MainCmd()
						if err != nil {
							logger.Warn("Worker %s Executing activating: %s (%s) -> Failed for: %v", a.FormatWorker(workerId), target.Remove.AppName, logger.VersionNamePrefix(target.Remove.Version), err)
						}
					}
				},
			})
		}

		if target.Deploy != nil {
			a.EvalTask(defArgs, &BashCommandTask{
				workerId:   target.WorkerIds,
				preDeploy:  target.Deploy.PreDeploy,
				postDeploy: target.Deploy.PostDeploy,
				deployTask: func() {
					err := a.resolveAppAlias(&target.Deploy.AppName, &target.Deploy.Version)
					if err != nil {
						logger.Error("Cannot resolve app alias for deploy: %q", target.Deploy.AppName)
						return
					}

					if target.Deploy.Version == "" {
						logger.Error("Version not specified, Abort Deploy")
						return
					}

					for _, workerId := range target.WorkerIds {
						worker := []string{workerId}

						pt := &deploy.PushCmdTask{}
						pt.Init(*defArgs, cmds.PushCmd{
							WorkerId: worker, AppName: target.Deploy.AppName,
							Version: target.Deploy.Version, MaxThread: 0,
						}, proto.TaskGroup_TASK_DEPLOY)

						err := pt.MainCmd()
						if err != nil {
							logger.Warn("Worker %s Executing deploying: %s (%s) -> Failed for: %v", a.FormatWorker(workerId), target.Deploy.AppName, logger.VersionNamePrefix(target.Deploy.Version), err)
							continue
						}

						st := &deploy.SetCmdTask{}
						st.Init(*defArgs, cmds.SetCmd{
							WorkerId: worker, AppName: target.Deploy.AppName, Version: target.Deploy.Version,
						}, proto.TaskGroup_TASK_DEPLOY)

						err = st.MainCmd()
						if err != nil {
							logger.Warn("Worker %s Executing activating: %s (%s) -> Failed for: %v", a.FormatWorker(workerId), target.Deploy.AppName, logger.VersionNamePrefix(target.Deploy.Version), err)
						}
					}
				},
			})
		}

		if target.Push != nil {
			a.EvalTask(defArgs, &BashCommandTask{
				workerId:   target.WorkerIds,
				preDeploy:  target.Push.PreDeploy,
				postDeploy: target.Push.PostDeploy,
				deployTask: func() {
					err := a.resolveAppAlias(&target.Push.AppName, &target.Push.Version)
					if err != nil {
						logger.Error("Cannot resolve app alias for push: %q", target.Push.AppName)
						return
					}

					if target.Push.Version == "" {
						logger.Error("Version not specified, Abort Push")
						return
					}

					for _, workerId := range target.WorkerIds {
						t := &deploy.PushCmdTask{}
						t.Init(*defArgs, cmds.PushCmd{
							WorkerId: []string{workerId}, AppName: target.Push.AppName,
							Version: target.Push.Version, MaxThread: 0,
						}, proto.TaskGroup_TASK_DEPLOY)

						err := t.MainCmd()
						if err != nil {
							logger.Warn("Worker %s Executing pushing: %s (%s) -> Failed for: %v", a.FormatWorker(workerId), target.Push.AppName, logger.VersionNamePrefix(target.Push.Version), err)
						}
					}
				},
			})
		}

		if target.Set != nil {
			a.EvalTask(defArgs, &BashCommandTask{
				workerId:   target.WorkerIds,
				preDeploy:  target.Set.PreDeploy,
				postDeploy: target.Set.PostDeploy,
				deployTask: func() {
					err := a.resolveAppAlias(&target.Set.AppName, &target.Set.Version)
					if err != nil {
						logger.Error("Cannot resolve app alias for set: %q", target.Set.AppName)
						return
					}

					for _, workerId := range target.WorkerIds {
						t := &deploy.SetCmdTask{}
						t.Init(*defArgs, cmds.SetCmd{
							WorkerId: []string{workerId}, AppName: target.Set.AppName, Version: target.Set.Version,
						}, proto.TaskGroup_TASK_DEPLOY)

						err := t.MainCmd()
						if err != nil {
							logger.Warn("Worker %s Executing activating: %s (%s) -> Failed for: %v", a.FormatWorker(workerId), target.Set.AppName, logger.VersionNamePrefix(target.Set.Version), err)
						}
					}
				},
			})
		}

		if target.Disconn != nil {
			a.EvalTask(defArgs, &BashCommandTask{
				workerId:   target.WorkerIds,
				preDeploy:  target.Disconn.PreDeploy,
				postDeploy: "",
				deployTask: func() {
					t := &nodes.DisConnCmdTask{}
					t.Init(*defArgs, cmds.DisconnCmd{WorkerId: target.WorkerIds}, proto.TaskGroup_TASK_MANAGE_NODE)
					err := t.MainCmd()
					if err != nil {
						logger.Warn("Disconnected failed: %v", err)
					}
				},
			})
		}

		logger.DisableTree(true)
		logger.Task("Finished processing target \"%s\"", target.TargetName)
	}
	return nil
}

type BashCommandTask struct {
	workerId   []string
	preDeploy  string
	postDeploy string
	deployTask func()
}

// EvalTask TODO: continue without error worker
func (a *StaFileTask) EvalTask(defArgs *cmds.DefaultArgs, task *BashCommandTask) {
	if task.preDeploy != "" {
		for _, target := range task.workerId {
			err := a.callBash(defArgs, target, task.preDeploy)
			if err != nil {
				logger.Error("Worker %s Executing pre-deploy: \"%s\" -> Failed for: %v", a.FormatWorker(target), task.preDeploy, err)
				return
			}
			logger.Hook("Worker %s Executing pre-deploy: \"%s\" -> Success.", a.FormatWorker(target), task.preDeploy)
		}
	}
	task.deployTask()
	if task.postDeploy != "" {
		for _, target := range task.workerId {
			err := a.callBash(defArgs, target, task.postDeploy)
			if err != nil {
				logger.Error("Worker %s Executing post-deploy: \"%s\" -> Failed for: %v", a.FormatWorker(target), task.postDeploy, err)
				continue
			}
			logger.Hook("Worker %s Executing post-deploy: \"%s\" -> Success.", a.FormatWorker(target), task.postDeploy)
		}
	}
}

func (a *StaFileTask) callBash(defArgs *cmds.DefaultArgs, workerId string, command string) error {
	t := &nodes.BashCmdTask{}
	t.Init(*defArgs, cmds.BashCmd{
		WorkerId: workerId, Command: command, NoOutput: true,
	}, proto.TaskGroup_TASK_MANAGE_NODE)
	return t.MainCmd()
}
