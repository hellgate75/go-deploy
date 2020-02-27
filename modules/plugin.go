package modules

import (
	"errors"
	"fmt"
	"github.com/hellgate75/go-deploy/io"
	"github.com/hellgate75/go-deploy/types/module"
	"plugin"
)

var ModulesFolder = "mod"

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

func seek(module string, symbol string) (plugin.Symbol, error) {
	plugin, err := plugin.Open(module)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Errors loading plugin module : \"%s\". Details: %s", module, err.Error()))
	}
	component, errL := plugin.Lookup(symbol)
	if errL != nil {
		return nil, errors.New(fmt.Sprintf("Errors looking up for \"%s\" in plugin module : \"%s\". Details: %s", symbol, module, errL.Error()))
	}
	return component, nil
}

func LoadExecutorForModule(module string) (Executor, error) {
	var path string = "./" + io.GetPathSeparator() + ModulesFolder + io.GetPathSeparator() + module + io.GetPathSeparator() + module + io.GetShareLibExt()
	symCollector, err := seek(path, "Executor")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Errors fetching plugin module : \"%s\". Details: %s", module, err.Error()))
	}
	var executor Executor
	executor, ok := symCollector.(Executor)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Uanble to parse Executor for module: %s", module))
	}
	return executor, nil
}

func LoadConverterForModule(module string) (Converter, error) {
	var path string = "./" + io.GetPathSeparator() + ModulesFolder + io.GetPathSeparator() + module + io.GetPathSeparator() + module + io.GetShareLibExt()
	symCollector, err := seek(path, "Executor")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Errors fetching plugin module : \"%s\". Details: %s", module, err.Error()))
	}
	var converter Converter
	converter, ok := symCollector.(Converter)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Uanble to parse Converter for module: %s", module))
	}
	return converter, nil
}
