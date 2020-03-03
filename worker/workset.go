package worker

import (
	"github.com/hellgate75/go-deploy/log"
	"github.com/hellgate75/go-deploy/types/defaults"
	"github.com/hellgate75/go-deploy/types/module"
)

func ExecuteFeed(config defaults.ConfigPattern, feed *module.FeedExec, logger log.Logger) []error {
	var errList []error = make([]error, 0)
	return errList
}
