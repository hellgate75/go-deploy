package service

import (
	"errors"
	"fmt"
	"github.com/hellgate75/go-deploy/types/generic"
	"reflect"
	"strings"
)

type serviceExecutor struct {
}

func (shell *serviceExecutor) Execute(step generic.Step) error {
	return nil
}

var Executor generic.Executor = &serviceExecutor{}

var ERROR_TYPE reflect.Type = reflect.TypeOf(errors.New(""))

/*
* Service command structure
 */
type serviceCommand struct {
	Name     string
	State    string
	WithVars []string
	WithList []string
}

func (service serviceCommand) String() string {
	return fmt.Sprintf("serviceCommand {Name: %v, State: %v, WithVars: [%v], WithList: [%v]}", service.Name, service.State, service.WithVars, service.WithList)
}

func (service *serviceCommand) Convert(cmdValues interface{}) (interface{}, error) {
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
	var name, state string
	var withVars []string = make([]string, 0)
	var withList []string = make([]string, 0)
	var valType string = fmt.Sprintf("%v", cmdValues)
	if len(valType) > 3 && "map" == valType[0:3] {
		for key, value := range cmdValues.(map[string]interface{}) {
			var elemValType string = fmt.Sprintf("%v", value)
			if strings.ToLower(key) == "name" {
				if elemValType == "string" {
					name = fmt.Sprintf("%v", value)
				} else {
					return nil, errors.New("Unable to parse command: service.name, with aguments of type " + elemValType + ", expected type string")
				}
			} else if strings.ToLower(key) == "state" {
				if elemValType == "string" {
					state = fmt.Sprintf("%v", value)
				} else {
					return nil, errors.New("Unable to parse command: service.state, with aguments of type " + elemValType + ", expected type string")
				}
			} else if strings.ToLower(key) == "withvars" {
				if elemValType == "[]string" {
					for _, val := range value.([]string) {
						withVars = append(withVars, val)
					}
				} else {
					return nil, errors.New("Unable to parse command: service.asRoot, with aguments of type " + elemValType + ", expected type []string")
				}
			} else if strings.ToLower(key) == "withlist" {
				if elemValType == "[]string" {
					for _, val := range value.([]string) {
						withList = append(withList, val)
					}
				} else {
					return nil, errors.New("Unable to parse command: service.withList, with aguments of type " + elemValType + ", expected type []string")
				}
			} else {
				return nil, errors.New("Unknown command: service." + key)
			}
		}
	} else {
		return nil, errors.New("Unable to parse command: service, with aguments of type " + valType + ", expected type map[string]interfce{}")
	}
	if superError != nil {
		return nil, superError
	}
	return serviceCommand{
		Name:     name,
		State:    state,
		WithVars: withVars,
		WithList: withList,
	}, nil
}

var Converter generic.Converter = &serviceCommand{}

func main() {}
