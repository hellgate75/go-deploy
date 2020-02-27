package modules

import (
	"errors"
	"fmt"
	//	"github.com/hellgate75/go-deploy/io"
	"github.com/hellgate75/go-deploy/log"
	"github.com/hellgate75/go-deploy/types/module"
	//	"github.com/hellgate75/go-deploy/utils"
	"strconv"
	"time"
)

var Logger log.Logger = nil

var ModulesFolder = "mod"

type Symbol interface{}

/*
* Coverter interface, responsible to comvert raw interface from the parsing to a specific structure
 */
type Converter interface {
	/*
	* Converts a raw interface element to a command qualified structure <BR/>
	* Paramameters: <BR/>
	* cmdValues (interface{}) Raw value from the feed file parsing
	* Return: <BR/>
	* (interface{}) Qualified structure <BR/>
	* (error) Error occured during any conversion <BR/>
	 */
	Convert(cmdValues interface{}) (interface{}, error)
}

type Executor interface {
	Execute(step module.Step) error
}

//linked to modules_seekModule in all child modules
func seekModule(module string, feature string, acceptance chan bool, featureAcceptance chan bool, response chan interface{}) {
}

const (
	moduleAcceptanceTimeoutInSeconds int = 3
)

type SeekRequest struct {
	Module string
	Symbol string
}

func seek(module string, symbol string) (interface{}, error) {
	var acceptance chan bool = make(chan bool)
	var featureAcceptance chan bool = make(chan bool)
	var response chan interface{} = make(chan interface{})
	defer func() {
		if r := recover(); r != mil {
			Logger.Error(fmt.Sprintf("modules.seek -> Error: %v", r))
		}
		close(acceptance)
		close(featureAcceptance)
		close(response)
	}()
	Logger.Warn("Before Call ...")
	seekModule(module, symbol, acceptance, featureAcceptance, response)
	Logger.Warn("After Call ...")
	var accepted bool = false
	select {
	case res := <-acceptance:
		Logger.Warn(fmt.Sprintf("Acceptance response received: %s", strconv.FormatBool(res)))
		if res {
			accepted = true
		}
	case <-time.After(time.Duration(moduleAcceptanceTimeoutInSeconds) * time.Second):
		Logger.Warn("Call to modules' acceptance timed out ...")
	}
	var featureAccepted bool = false
	if accepted {
		select {
		case res := <-featureAcceptance:
			Logger.Warn(fmt.Sprintf("Feature Acceptance response received: %s", strconv.FormatBool(res)))
			if res {
				featureAccepted = true
			}
		case <-time.After(time.Duration(moduleAcceptanceTimeoutInSeconds) * time.Second):
			Logger.Warn("Call to modules feature acceptance timed out ...")
		}
	}
	var outcome interface{}
	if featureAccepted {
		select {
		case res := <-response:
			Logger.Warn(fmt.Sprintf("Module Component response received: %v", res))
			outcome = res
		case <-time.After(time.Duration(moduleAcceptanceTimeoutInSeconds) * time.Second):
			Logger.Warn("Call to modules feature acceptance timed out ...")
		}
	}
	Logger.Warn(fmt.Sprintf("Found module library component: ", outcome))
	//	component, errL := plugin.Lookup(symbol)
	//	if errL != nil {
	//		return nil, errors.New(fmt.Sprintf("Errors looking up for \"%s\" in plugin module : \"%s\". Details: %s", symbol, module, errL.Error()))
	//	}
	var err error = nil
	if outcome == nil {
		err = errors.New(fmt.Sprintf("Unable to find component %s in module %s", symbol, module))
	}
	return outcome, err
}

func LoadExecutorForModule(module string) (Executor, error) {
	//	var path string = io.GetCurrentFolder() + io.GetPathSeparator() + ModulesFolder + io.GetPathSeparator() + module + io.GetPathSeparator() + module + utils.GetShareLibExt()
	//	var exists bool = io.ExistsFile(path)
	//	Logger.Warn(fmt.Sprintf("modules.LoadExecutorForModule -> Current Path: %s, exists: %s", path, strconv.FormatBool(exists)))
	//	if !exists {
	//		return nil, errors.New(fmt.Sprintf("plugin.LoadExecutorForModule -> File %s doesn't exist", path))
	//	}
	symCollector, err := seek(module, "Executor")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Errors fetching plugin module : \"%s\". Details: %s", module, err.Error()))
	}
	var executor Executor
	executor, ok := symCollector.(Executor)
	Logger.Warn(fmt.Sprintf("modules.LoadExecutorForModule -> On Module: %s, found Executor: %v", module, executor))
	if !ok {
		return nil, errors.New(fmt.Sprintf("Uanble to parse Executor for module: %s", module))
	}
	return executor, nil
}

func LoadConverterForModule(module string) (Converter, error) {
	//	var path string = io.GetCurrentFolder() + io.GetPathSeparator() + ModulesFolder + io.GetPathSeparator() + module + io.GetPathSeparator() + module + utils.GetShareLibExt()
	//	var exists bool = io.ExistsFile(path)
	//	Logger.Warn(fmt.Sprintf("modules.LoadConverterForModule -> Current Path: %s, exists: %s", path, strconv.FormatBool(exists)))
	//	if !exists {
	//		return nil, errors.New(fmt.Sprintf("plugin.LoadConverterForModule -> File %s doesn't exist", path))
	//	}
	symCollector, err := seek(module, "Converter")
	Logger.Warn(fmt.Sprintf("modules.LoadConverterForModule -> symCollector: %v", symCollector))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Errors fetching plugin module : \"%s\". Details: %s", module, err.Error()))
	}
	var converter Converter
	converter, ok := symCollector.(Converter)
	Logger.Warn(fmt.Sprintf("modules.LoadConverterForModule -> On Module: %s, found Converters: %v", module, converter))
	if !ok {
		return nil, errors.New(fmt.Sprintf("Uanble to parse Converter for module: %s", module))
	}
	return converter, nil
}
