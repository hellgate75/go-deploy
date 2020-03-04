package worker

import (
	"github.com/hellgate75/go-deploy/log"
	"github.com/hellgate75/go-deploy/net/generic"
	"github.com/hellgate75/go-deploy/types/defaults"
	"github.com/hellgate75/go-deploy/types/module"
)

var ClientCache map[string]generic.NetworkClient = make(map[string]generic.NetworkClient)

func ExecuteFeed(config defaults.ConfigPattern, feed *module.FeedExec, logger log.Logger) []error {
	var errList []error = make([]error, 0)
	var feedName string = feed.Name
	if feedName == "" {
		feedName = "<none>"
	}
	logger.Infof("Executing on feed : %s", feedName)
	//	itf, err := shell.session.GetSystemObject("connection-handler")
	//	var handler generic.ConnectionHandler
	//	if err != nil {
	//		return err
	//	}
	//	handler = itf.(generic.ConnectionHandler)
	//	Logger.Debugf("Handler present: %v", (handler != nil))
	//	client, err = generic.ConnectHandlerViaConfig(handler, shell.host, shell.config.Net, shell.config.Config)
	//	if err != nil {
	//		return err
	//	}
	return errList
}
