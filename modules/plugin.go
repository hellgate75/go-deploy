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
	var acceptance chan bool = make(chan bool)
	var featureAcceptance chan bool = make(chan bool)
	var response chan interface{} = make(chan interface{})
	defer func() {
		if r := recover(); r != nil {
			Logger.Errorf("modules.seek -> Error: %v", r)
		}
		close(acceptance)
		close(featureAcceptance)
		close(response)
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
	return itf, nil
	//	seekModule(module, symbol, acceptance, featureAcceptance, response)
	//	Logger.Warn("After Call ...")
	//	var accepted bool = false
	//	select {
	//	case res := <-acceptance:
	//		Logger.Warn(fmt.Sprintf("Acceptance response received: %s", strconv.FormatBool(res)))
	//		if res {
	//			accepted = true
	//		}
	//	case <-time.After(time.Duration(moduleAcceptanceTimeoutInSeconds) * time.Second):
	//		Logger.Warn("Call to modules' acceptance timed out ...")
	//	}
	//	var featureAccepted bool = false
	//	if accepted {
	//		select {
	//		case res := <-featureAcceptance:
	//			Logger.Warn(fmt.Sprintf("Feature Acceptance response received: %s", strconv.FormatBool(res)))
	//			if res {
	//				featureAccepted = true
	//			}
	//		case <-time.After(time.Duration(moduleAcceptanceTimeoutInSeconds) * time.Second):
	//			Logger.Warn("Call to modules feature acceptance timed out ...")
	//		}
	//	}
	//	var outcome interface{}
	//	if featureAccepted {
	//		select {
	//		case res := <-response:
	//			Logger.Warn(fmt.Sprintf("Module Component response received: %v", res))
	//			outcome = res
	//		case <-time.After(time.Duration(moduleAcceptanceTimeoutInSeconds) * time.Second):
	//			Logger.Warn("Call to modules feature acceptance timed out ...")
	//		}
	//	}
	//	Logger.Warn(fmt.Sprintf("Found module library component: ", outcome))
	//	//	component, errL := plugin.Lookup(symbol)
	//	//	if errL != nil {
	//	//		return nil, errors.New(fmt.Sprintf("Errors looking up for \"%s\" in plugin module : \"%s\". Details: %s", symbol, module, errL.Error()))
	//	//	}
	//	var err error = nil
	//	if outcome == nil {
	//		err = errors.New(fmt.Sprintf("Unable to find component %s in module %s", symbol, module))
	//	}
	//	return outcome, err
}

func LoadConverterForModule(module string) (meta.Converter, error) {
	//	var path string = io.GetCurrentFolder() + io.GetPathSeparator() + ModulesFolder + io.GetPathSeparator() + module + io.GetPathSeparator() + module + utils.GetShareLibExt()
	//	var exists bool = io.ExistsFile(path)
	//	Logger.Warn(fmt.Sprintf("modules.LoadConverterForModule -> Current Path: %s, exists: %s", path, strconv.FormatBool(exists)))
	//	if !exists {
	//		return nil, errors.New(fmt.Sprintf("plugin.LoadConverterForModule -> File %s doesn't exist", path))
	//	}
	converter, err := seek(module)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Errors fetching plugin module : \"%s\". Details: %s", module, err.Error()))
	}
	Logger.Debugf("modules.LoadConverterForModule -> On Module: %s, found Converters: %v", module, converter)
	return converter, nil
}
