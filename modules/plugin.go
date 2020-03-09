package modules

import (
	"errors"
	"fmt"
	"github.com/hellgate75/go-tcp-common/log"
	"github.com/hellgate75/go-deploy/modules/meta"
	"github.com/hellgate75/go-deploy/modules/proxy"
)

var Logger log.Logger = nil

var ModulesFolder = "mod"

const (
	moduleAcceptanceTimeoutInSeconds int = 3
)

// Interface that describe a Request in the Deploy Manager Plugin Explorer
type SeekRequest struct {
	Module string
	Symbol string
}

var proxyVia proxy.Proxy = nil

func seek(module string) (meta.Converter, error) {
	var errGlobal error = nil
	defer func() {
		if r := recover(); r != nil {
			Logger.Errorf("modules.seek -> Error: %v", r)
			errGlobal = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	if proxyVia == nil {
		proxyVia = proxy.NewProxy()
	}
	var mod proxy.Module
	var err error
	mod, err = proxyVia.DiscoverModule(module)
	Logger.Debugf("Module is present: %v", (mod != nil))
	if err != nil {
		return nil, err
	}
	var itf meta.Converter
	itf, err = mod.GetComponent()
	Logger.Debugf("Module (%s) Component: %v", module, itf)
	if err != nil {
		return nil, err
	}
	itf.SetLogger(Logger)
	Logger.Debugf("Module (%s) Logger: %v", module, (Logger != nil))
	return itf, errGlobal
}

// Load modules that matches with the command, and provide the most suitable Converter Component, or giver the reason of the error
func LoadConverterForModule(module string) (meta.Converter, error) {
	converter, err := seek(module)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Errors fetching plugin module : \"%s\". Details: %s", module, err.Error()))
	}
	Logger.Debugf("modules.LoadConverterForModule -> On Module: %s, found Converters: %v", module, converter)
	return converter, nil
}
