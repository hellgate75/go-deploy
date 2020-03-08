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

type SeekRequest struct {
	Module string
	Symbol string
}

var proxyVia proxy.Proxy = proxy.NewProxy()

func seek(module string) (meta.Converter, error) {
	var errGlobal error = nil
	defer func() {
		if r := recover(); r != nil {
			Logger.Errorf("modules.seek -> Error: %v", r)
			errGlobal = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	var mod proxy.Module
	var err error
	mod, err = proxyVia.DiscoverModule(module)
	Logger.Debugf("Module is present: %v", (mod != nil))
	if err != nil {
		return nil, err
	}
	var itf meta.Converter
	itf, err = mod.GetComponent()
	Logger.Debugf("Module Component: %v", itf)
	if err != nil {
		return nil, err
	}
	itf.SetLogger(Logger)
	return itf, errGlobal
}

func LoadConverterForModule(module string) (meta.Converter, error) {
	converter, err := seek(module)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Errors fetching plugin module : \"%s\". Details: %s", module, err.Error()))
	}
	Logger.Debugf("modules.LoadConverterForModule -> On Module: %s, found Converters: %v", module, converter)
	return converter, nil
}
