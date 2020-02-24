package cmdtypes

import (
	"fmt"
	"errors"
	"strconv"
	"strings"
	"reflect"
)

type Converter interface {
	Convert(cmdValues interface{}) (interface{}, error)
}

type NilCommandConverter struct {
	CmdType int
}
func (nilCommand *NilCommandConverter) Convert(cmdValues interface{})  (interface{}, error) {
	return nil, errors.New("Not implemented type: " + strconv.Itoa(nilCommand.CmdType))
	
}

type ShellCommand struct {
	CmdName		string
	Exec 		string
	RunAs		string
	AsRoot		bool
	WithVars	[]string
}

func (shell ShellCommand) String() string {
	return fmt.Sprintf("ShellCommand {Exec: %v, RunAs: %v, AsRoot: %v, WithVars: [%v]}", shell.Exec, shell.RunAs, strconv.FormatBool(shell.AsRoot), shell.WithVars)
}

func (shell *ShellCommand) Convert(cmdValues interface{})  (interface{}, error) {
	var superError error = nil
	defer func() {
		if r := recover(); r != nil {
            if reflect.TypeOf(r).Kind().String() == "error" {
            	superError = r.(error)
            } else {
            	superError = errors.New(fmt.Sprintf("%v", r))
            }
        }
		
	}()
	var valType string = fmt.Sprintf("%v", cmdValues)
	var cmdName string = ""
	var exec string = ""
	var runAs string = ""
	var asRoot bool = false
	var withVars	[]string = make([]string, 0)
	if len(valType) > 3 && "map" == valType[0:3] {
		for key, value := range cmdValues.(map[string]interface{}) {
				var elemValType string = fmt.Sprintf("%v", value)
				if strings.ToLower(key) == "name" {
					if elemValType == "string" {
						cmdName = fmt.Sprintf("%v", value)
					} else {
						return nil, errors.New("Unable to parse command: shell.name, with aguments of type " + elemValType + ", expected type string")
					}
				} else if strings.ToLower(key) == "exec" {
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
	if (  superError != nil) {
		return nil, superError
	}
	return ShellCommand{
		CmdName: cmdName,
		Exec: exec,
		RunAs: runAs,
		AsRoot: asRoot,
		WithVars: withVars,
	}, nil
}

type ServiceCommand struct {
	CmdName		string
	Name 		string
	State		string
	WithVars	[]string
}

func (service ServiceCommand) String() string {
	return fmt.Sprintf("ServiceCommand {Name: %v, State: %v}", service.Name, service.State)
}

func (service *ServiceCommand) Convert(cmdValues interface{})  (interface{}, error) {
	var superError error = nil
	defer func() {
		if r := recover(); r != nil {
            if reflect.TypeOf(r).Kind().String() == "error" {
            	superError = r.(error)
            } else {
            	superError = errors.New(fmt.Sprintf("%v", r))
            }
        }
		
	}()
	var cmdName string = ""
	var name, state string
	var withVars []string = make([]string, 0)
	var valType string = fmt.Sprintf("%v", cmdValues)
	if len(valType) > 3 && "map" == valType[0:3] {
		for key, value := range cmdValues.(map[string]interface{}) {
				var elemValType string = fmt.Sprintf("%v", value)
				if strings.ToLower(key) == "name" {
					if elemValType == "string" {
						cmdName = fmt.Sprintf("%v", value)
					} else {
						return nil, errors.New("Unable to parse command: shell.name, with aguments of type " + elemValType + ", expected type string")
					}
				} else if strings.ToLower(key) == "service" {
					if elemValType == "string" {
						name = fmt.Sprintf("%v", value)
					} else {
						return nil, errors.New("Unable to parse command: shell.service, with aguments of type " + elemValType + ", expected type string")
					}
				} else if strings.ToLower(key) == "state" {
					if elemValType == "string" {
						state = fmt.Sprintf("%v", value)
					} else {
						return nil, errors.New("Unable to parse command: shell.state, with aguments of type " + elemValType + ", expected type string")
					}
				} else if strings.ToLower(key) == "withvars" {
					if elemValType == "[]string" {
						for _, val := range value.([]string) {
							withVars = append(withVars, val) 
						}
					} else {
						return nil, errors.New("Unable to parse command: shell.asRoot, with aguments of type " + elemValType + ", expected type []string")
					}
				} else {
					return nil, errors.New("Unknown command: service." + key)
	            }	
		}
	} else {
		return nil, errors.New("Unable to parse command: service, with aguments of type " + valType + ", expected type map[string]interfce{}")
	}
	if (  superError != nil) {
		return nil, superError
	}
	return ServiceCommand{
		CmdName: cmdName,
		Name: name,
		State: state,
		WithVars: withVars,
	}, nil
}

func NewConverter(cmdType int) Converter {
	switch cmdType	{
		case FEED_TYPE_SHELL:
			return &ShellCommand{}
		case FEED_TYPE_SERVICE:
			return &ServiceCommand{}
		case FEED_TYPE_FACT:
	}
	return &NilCommandConverter{
		CmdType: cmdType,
	}
	
}
