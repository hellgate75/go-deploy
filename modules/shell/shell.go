package main

import (
	"errors"
	"fmt"
	"github.com/hellgate75/go-deploy/log"
	"github.com/hellgate75/go-deploy/modules"
	"github.com/hellgate75/go-deploy/types/module"
	"reflect"
	"strconv"
	"strings"
)

var Logger log.Logger = log.NewLogger(log.VerbosityLevelFromString(module.RuntimeDeployConfig.LogVerbosity))

type shellExecutor struct {
}

func (shell *shellExecutor) Execute(step module.Step) error {
	Logger.Warn(fmt.Sprintf("shell.Executor.Execute -> Executing command : %s", step.StepType))
	return nil
}

var Executor modules.Executor = &shellExecutor{}

var ERROR_TYPE reflect.Type = reflect.TypeOf(errors.New(""))

func (shell shellCommand) String() string {
	return fmt.Sprintf("ShellCommand {Exec: %v, RunAs: %v, AsRoot: %v, WithVars: [%v], WithList: [%v]}", shell.Exec, shell.RunAs, strconv.FormatBool(shell.AsRoot), shell.WithVars, shell.WithList)
}

/*
* Shell command structure
 */
type shellCommand struct {
	Exec     string
	RunAs    string
	AsRoot   bool
	WithVars []string
	WithList []string
}

func (shell *shellCommand) Convert(cmdValues interface{}) (interface{}, error) {
	var superError error = nil
	defer func() {
		if r := recover(); r != nil {
			if ERROR_TYPE.AssignableTo(reflect.TypeOf(r)) {
				superError = r.(error)
			} else {
				superError = errors.New(fmt.Sprintf("%v", r))
			}
		}

	}()
	var valType string = fmt.Sprintf("%v", cmdValues)
	var exec string = ""
	var runAs string = ""
	var asRoot bool = false
	var withVars []string = make([]string, 0)
	var withList []string = make([]string, 0)
	if len(valType) > 3 && "map" == valType[0:3] {
		for key, value := range cmdValues.(map[string]interface{}) {
			var elemValType string = fmt.Sprintf("%v", value)
			if strings.ToLower(key) == "exec" {
				if elemValType == "string" {
					exec = fmt.Sprintf("%v", value)
				} else if elemValType == "[]string" {
					strings.Join(value.([]string), " ")
				} else {
					return nil, errors.New("Unable to parse command: shell.exec, with aguments of type " + elemValType + ", expected type string or []string")
				}
			} else if strings.ToLower(key) == "runas" {
				if elemValType == "string" {
					runAs = fmt.Sprintf("%v", value)
				} else {
					return nil, errors.New("Unable to parse command: shell.runAs, with aguments of type " + elemValType + ", expected type string")
				}
			} else if strings.ToLower(key) == "asroot" {
				if elemValType == "string" {
					bl, err := strconv.ParseBool(fmt.Sprintf("%v", value))
					if err != nil {
						return nil, errors.New("Error parsing command: shell.asRoot, cause: " + err.Error())

					} else {
						asRoot = bl
					}

				} else if elemValType == "bool" {
					asRoot = value.(bool)
				} else {
					return nil, errors.New("Unable to parse command: shell.asRoot, with aguments of type " + elemValType + ", expected type bool or string")
				}
			} else if strings.ToLower(key) == "withvars" {
				if elemValType == "[]string" {
					for _, val := range value.([]string) {
						withVars = append(withVars, val)
					}
				} else {
					return nil, errors.New("Unable to parse command: shell.asRoot, with aguments of type " + elemValType + ", expected type []string")
				}
			} else if strings.ToLower(key) == "withlist" {
				if elemValType == "[]string" {
					for _, val := range value.([]string) {
						withList = append(withList, val)
					}
				} else {
					return nil, errors.New("Unable to parse command: shell.withList, with aguments of type " + elemValType + ", expected type []string")
				}
			} else {
				return nil, errors.New("Unknown command: shell." + key)

			}
		}
	} else {
		return nil, errors.New("Unable to parse command: shell, with aguments of type " + valType + ", expected type map[string]interfce{}")
	}
	if exec == "" {
		return nil, errors.New("Missing command: shell.exec -> mandatory field")

	}
	if superError != nil {
		return nil, superError
	}
	return shellCommand{
		Exec:     exec,
		RunAs:    runAs,
		AsRoot:   asRoot,
		WithVars: withVars,
	}, nil
}

var Converter modules.Converter = &shellCommand{}

func main() {}
