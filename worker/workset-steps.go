package worker

import (
	"errors"
	"fmt"
	"github.com/hellgate75/go-deploy/log"
	"github.com/hellgate75/go-deploy/net/generic"
	"github.com/hellgate75/go-deploy/types/defaults"
	"github.com/hellgate75/go-deploy/types/module"
	"github.com/hellgate75/go-deploy/types/threads"
	"github.com/hellgate75/go-deploy/worker/pool"
)

var clientsCache map[string]generic.NetworkClient = make(map[string]generic.NetworkClient)

func ExecuteSteps(prefix string, steps []*module.Step,
	selectedHostGroup *defaults.HostGroups, threadPool pool.ThreadPool,
	errorsHandler *ErrorHandler, config defaults.ConfigPattern,
	sessionsMap map[string]module.Session, logger log.Logger) []error {
	var errList []error = make([]error, 0)
	for _, step := range steps {
		errorsHandler.Reset()
		stepName := step.Name
		if stepName == "" {
			stepName = "<none>"
		}
		logger.Warnf("%s[ %s ]", prefix, stepName)
		if step.StepData != nil {
			thread := step.StepData.(threads.StepRunnable)
			var threadsMap map[string]threads.StepRunnable = make(map[string]threads.StepRunnable)
			for _, host := range selectedHostGroup.Hosts {
				hostThread := thread.Clone()
				sessMapId := fmt.Sprintf("%s-%s", selectedHostGroup.Name, host.Name)
				if client, ok := clientsCache[sessMapId]; ok {
					hostThread.SetClient(client)
				}
				if session, ok := sessionsMap[sessMapId]; ok {
					hostThread.SetSession(session)
				}
				hostThread.SetConfig(config)
				hostThread.SetHost(host)
				threadsMap[sessMapId] = hostThread
				logger.Debugf("Scheduling step process for %s - %s ...", selectedHostGroup.Name, host.Name)
				threadPool.Schedule(hostThread)
				logger.Debugf("Scheduled step process for %s - %s!!", selectedHostGroup.Name, host.Name)
			}
			threadPool.Start()
			err := threadPool.WaitFor()
			threadPool.Stop()
			if err != nil {
				errList = append(errList, errors.New(fmt.Sprintf("%v", err)))
				return errList
			}
			threadErrorsList := errorsHandler.GetAll()
			if len(threadErrorsList) > 0 {
				for _, host := range selectedHostGroup.Hosts {
					sessMapId := fmt.Sprintf("%s-%s", selectedHostGroup.Name, host.Name)
					var item ErrorItem = ErrorItem{}
					if threadX, ok := threadsMap[sessMapId]; ok {
						for _, errItem := range threadErrorsList {
							if threadX.UUID() == errItem.UUID {
								item = errItem
								break
							}
						}
						if item.UUID != "" && item.Error != nil {
							logger.Failuref("- [Host: %s, Process Id: %s, status: ko]\n Error: %s", host.Name, threadX.UUID(), item.Error.Error())
						} else {
							logger.Successf("- [Host: %s, Process Id: %s, status: ok]", host.Name, threadX.UUID())
						}
					} else {
						errList = append(errList, errors.New("Thread Map not present for group: "+selectedHostGroup.Name+" and host: "+host.Name))
						return errList
					}
				}
			} else {
				for _, host := range selectedHostGroup.Hosts {
					sessMapId := fmt.Sprintf("%s-%s", selectedHostGroup.Name, host.Name)
					if threadX, ok := threadsMap[sessMapId]; ok {
						logger.Successf("- [Host: %s, Process Id: %s, status: ok]", host.Name, threadX.UUID())
					} else {
						errList = append(errList, errors.New("Thread Map not present for group: "+selectedHostGroup.Name+" and host: "+host.Name))
						return errList
					}
				}
			}

		} else {
			logger.Warn("No step executable found, progressing with children or next step ...")
		}
		if step.Children != nil && len(step.Children) > 0 {
			var subPrefix string = fmt.Sprintf("%s [ %s ]", prefix, stepName)
			errXList := ExecuteSteps(subPrefix, step.Children, selectedHostGroup, threadPool, errorsHandler, config, sessionsMap, logger)
			if len(errXList) > 0 {
				errList = append(errList, errXList...)
			}
		}
		if step.Feeds != nil && len(step.Feeds) > 0 {
			for _, feed := range step.Feeds {
				var feedName string = feed.Name
				if feedName == "" {
					feedName = "<none>"
				}
				logger.Warnf("Executing Feed: %s children of Step %s", feedName, stepName)
				errXList := ExecuteFeed(config, feed, sessionsMap, logger)
				if len(errXList) > 0 {
					errList = append(errList, errXList...)
				}
			}
		}
	}

	return errList
}
