package modules

import (
	"errors"
	"fmt"
	"github.com/hellgate75/go-deploy/io"
	"github.com/hellgate75/go-deploy/types/generic"
	"plugin"
)

var ModulesFolder = "mod"

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

func LoadExecutorForModule(module string) (generic.Executor, error) {
	var path string = "./" + io.GetPathSeparator() + ModulesFolder + io.GetPathSeparator() + module + io.GetPathSeparator() + module + io.GetShareLibExt()
	symCollector, err := seek(path, "Executor")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Errors fetching plugin module : \"%s\". Details: %s", module, err.Error()))
	}
	var executor generic.Executor
	executor, ok := symCollector.(generic.Executor)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Uanble to parse Executor for module: %s", module))
	}
	return executor, nil
}

func LoadConverterForModule(module string) (generic.Converter, error) {
	var path string = "./" + io.GetPathSeparator() + ModulesFolder + io.GetPathSeparator() + module + io.GetPathSeparator() + module + io.GetShareLibExt()
	symCollector, err := seek(path, "Executor")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Errors fetching plugin module : \"%s\". Details: %s", module, err.Error()))
	}
	var converter generic.Converter
	converter, ok := symCollector.(generic.Converter)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Uanble to parse Converter for module: %s", module))
	}
	return converter, nil
}
