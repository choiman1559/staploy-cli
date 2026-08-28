package build

import (
	"fmt"
	"staploy-cli/app/cmds"
	"staploy-cli/app/consts"
	"staploy-cli/app/logger"
	"strings"
)

type ResolvedAppAlias struct {
	Alias *AppAlias
	Build *cmds.BuildCmd
}

func (a *StaFileTask) processAlias(_ *cmds.DefaultArgs, alias *Alias) error {
	logger.Process("Parsing Alias... found %d worker alias and %d app alias", len(alias.Worker), len(alias.App))

	for _, workerAlias := range alias.Worker {
		logger.Process("Processing worker alias \"%s\"", workerAlias.Alias)
		logger.EnableTree()
		resolvedWorkers := workerAlias.WorkerIds

		if workerAlias.Where != nil {
			var filteredWorkers []string
			parsedWorkers, err := a.ParseWorkerAlias(false, workerAlias.WorkerIds...)

			if err != nil {
				logger.DisableTree(true)
				logger.Error("Failed to parse worker, cause: %s", err.Error())
				continue
			}

			for _, worker := range parsedWorkers {
				workerInfo, err := a.askWorkerInfoData(worker)
				if err == nil {
					passedFilter, err := a.evalFilter(workerInfo, workerAlias.Where)
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
				resolvedWorkers = filteredWorkers
			} else {
				logger.DisableTree(true)
				logger.Warn("No workers selected for worker alias \"%s\", skipping.", workerAlias.Alias)
				continue
			}
		}

		a.ResolvedWorkerAlias[workerAlias.Alias] = &WorkerAlias{WorkerIds: resolvedWorkers}
		logger.DisableTree(true)
		logger.Process("Finished worker alias \"%s\"", workerAlias.Alias)
	}

	for _, appAlias := range alias.App {
		logger.Process("Processing app alias \"%s\"", appAlias.Alias)
		logger.EnableTree()

		parsedAppName, err := a.parseExecArgs(&Build{}, appAlias.AppName)
		if err != nil {
			logger.DisableTree(true)
			logger.Error("Failed to parse app name \"%s\", cause: %s", appAlias.AppName, err.Error())
			continue
		}

		parsedVersion, err := a.parseExecArgs(&Build{}, appAlias.Version)
		if err != nil {
			logger.DisableTree(true)
			logger.Error("Failed to parse app version \"%s\", cause: %s", appAlias.AppName, err.Error())
			continue
		}

		a.ResolvedAppAlias[appAlias.Alias] = &ResolvedAppAlias{Alias: &AppAlias{
			Alias:   appAlias.Alias,
			AppName: parsedAppName,
			Version: logger.TrimVersion(parsedVersion),
		}}

		logger.DisableTree(true)
		logger.Process("Finished app alias \"%s\"", appAlias.Alias)
	}

	return nil
}

func (a *StaFileTask) hitAppAlias(appAliasName string, build *cmds.BuildCmd) (*ResolvedAppAlias, error) {
	appAlias := a.ResolvedAppAlias[appAliasName]
	if appAlias == nil {
		return nil, fmt.Errorf("app alias \"%s\" not found", appAliasName)
	}

	if build != nil {
		if appAlias.Build != nil {
			return appAlias, fmt.Errorf("build metadata of app alias \"%s\" already exists", appAliasName)
		}
		appAlias.Build = build
	}

	return appAlias, nil
}

func (a *StaFileTask) resolveAppAlias(appName *string, version *string) error {
	if !strings.HasPrefix(*appName, consts.STAFILE_ALIAS_PREFIX) {
		return nil
	}

	aliasName := strings.TrimPrefix(*appName, consts.STAFILE_ALIAS_PREFIX)
	foundAlias, err := a.hitAppAlias(aliasName, nil)
	if err != nil {
		return err
	}

	*appName = foundAlias.Alias.AppName
	if foundAlias.Alias.Version != "" && version != nil {
		*version = foundAlias.Alias.Version
	}

	return nil
}
