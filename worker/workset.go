package worker

import (
	"errors"
	"fmt"
	"github.com/gookit/color"
	"github.com/hellgate75/go-tcp-common/log"
	"github.com/hellgate75/go-tcp-common/pool"
	"runtime"
	"strings"
	
	"github.com/hellgate75/go-deploy/net/generic"
	"github.com/hellgate75/go-deploy/types/defaults"
	"github.com/hellgate75/go-deploy/types/module"
)

// Execute Feed, after definition of per Session and Host Client components
func ExecuteFeed(connectionConfig module.ConnectionConfig, config defaults.ConfigPattern, feed *module.FeedExec, sessionsMap map[string]module.Session, logger log.Logger) []error {
	var errList []error = make([]error, 0)
	var feedName string = feed.Name
	if feedName == "" {
		feedName = "<none>"
	}
	logger.Infof("Executing on feed : %s", feedName)
	logger.Infof("Hosts Group : %s", feed.HostGroup)
	var selectedHostGroup *defaults.HostGroups = &defaults.HostGroups{}
	for _, hg := range config.HostGroups {
		if strings.ToLower(hg.Name) == strings.ToLower(feed.HostGroup) {
			selectedHostGroup = &hg
			break
		}
	}
	if selectedHostGroup == nil || selectedHostGroup.Hosts == nil {
		errList = append(errList, errors.New("Unable to discover selected group in provided host groups ..."))
		return errList
	}
	logger.Info("Selected Hosts: ")
	for _, host := range selectedHostGroup.Hosts {
		//create host client and open connection ...
		logger.Infof("- %s", color.Yellow.Render(host.Name))
		sessMapId := fmt.Sprintf("%s-%s", selectedHostGroup.Name, host.Name)
		logger.Debugf("       -> session key: %s", color.Yellow.Render(sessMapId))
		if session, ok := sessionsMap[sessMapId]; ok {
			logger.Debugf("       -> session id: %s", session.GetSessionId())
			if _, ok := clientsCache[sessMapId]; !ok {
				itf, err := session.GetSystemObject("connection-handler")
				var handler generic.ConnectionHandler
				if err != nil {
					errList = append(errList, err)
					return errList
				}
				handler = itf.(generic.ConnectionHandler)
				if config.Plugins.EnableDeployClientCommandsPlugin {
					handler.UsePlugins(config.Plugins.DeployClientCommandsPluginExtension, config.Plugins.DeployClientCommandsPluginFolder)
				}
				logger.Debugf("       -> Handler Is present: %v", (handler != nil))
				var client generic.NetworkClient
				client, err = generic.ConnectHandlerViaConfig(connectionConfig, handler, host, config.Net, config.Config)
				if err != nil {
					errList = append(errList, err)
					return errList
				}
				defer client.Close()
				clientsCache[sessMapId] = client
				logger.Debugf("       -> Client Is present and connected: %v", (client != nil))
			} else {
				logger.Debugf("       -> Client already exists for group: %s host: %s", selectedHostGroup.Name, host.Name)
			}
		} else {
			errList = append(errList, errors.New("Session Map not present for group: "+selectedHostGroup.Name+" and host: "+host.Name))
			return errList
		}
	}
	threadPool := pool.NewThreadPool(config.Config.MaxThreads, config.Config.ParallelExecutions)
	threadPool.SetLogger(logger)
	errorsHandler := &ErrorHandler{
		errorList: make([]ErrorItem, 0),
	}
	threadPool.SetErrorHandler(errorsHandler)
	defer threadPool.Stop()
	errXList := ExecuteSteps("", feed.Steps, selectedHostGroup, threadPool,
						errorsHandler, config, sessionsMap, logger, connectionConfig)
	if len(errXList) > 0 {
		errList = append(errList, errXList...)
	}

	return errList
}

type ErrorItem struct {
	UUID  string
	Error error
}

type ErrorHandler struct {
	errorList []ErrorItem
}

func (handler *ErrorHandler) HandleError(uuid string, e error) {
	if e != nil {
		handler.errorList = append(handler.errorList, ErrorItem{
			UUID:  uuid,
			Error: e,
		})
	}
}

func (handler *ErrorHandler) Reset() {
	handler.errorList = make([]ErrorItem, 0)
	runtime.GC()
}

func (handler *ErrorHandler) GetAll() []ErrorItem {
	return handler.errorList
}
