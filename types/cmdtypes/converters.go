package cmdtypes

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

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

/*
* Coverter interface, responsible to comvert raw interface from the parsing to a specific structure
 */
type Printable interface {
	/*
	* Traslates the object in printable version <BR/>
	* Return: <BR/>
	* (string) Representation of the structure<BR/>
	 */
	String() string
}

/*
* Unknown command structure
 */
type NilCommandConverter struct {
	CmdType int
}

func (nilCommand *NilCommandConverter) Convert(cmdValues interface{}) (interface{}, error) {
	return nil, errors.New("Not implemented type: " + strconv.Itoa(nilCommand.CmdType))

}

/*
* Shell command structure
 */
type ShellCommand struct {
	Exec     string
	RunAs    string
	AsRoot   bool
	WithVars []string
	WithList []string
}

var ERROR_TYPE reflect.Type = reflect.TypeOf(errors.New(""))

func (shell ShellCommand) String() string {
	return fmt.Sprintf("ShellCommand {Exec: %v, RunAs: %v, AsRoot: %v, WithVars: [%v], WithList: [%v]}", shell.Exec, shell.RunAs, strconv.FormatBool(shell.AsRoot), shell.WithVars, shell.WithList)
}

func (shell *ShellCommand) Convert(cmdValues interface{}) (interface{}, error) {
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
	return ShellCommand{
		Exec:     exec,
		RunAs:    runAs,
		AsRoot:   asRoot,
		WithVars: withVars,
	}, nil
}

/*
* Service command structure
 */
type ServiceCommand struct {
	Name     string
	State    string
	WithVars []string
	WithList []string
}

func (service ServiceCommand) String() string {
	return fmt.Sprintf("ServiceCommand {Name: %v, State: %v, WithVars: [%v], WithList: [%v]}", service.Name, service.State, service.WithVars, service.WithList)
}

func (service *ServiceCommand) Convert(cmdValues interface{}) (interface{}, error) {
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
	return ServiceCommand{
		Name:     name,
		State:    state,
		WithVars: withVars,
		WithList: withList,
	}, nil
}

/*
* Service command structure
 */
type CopyCommand struct {
	SourceDir      string
	DestinationDir string
	CreateDest     bool
	WithVars       []string
	WithList       []string
}

func (copyCmd CopyCommand) String() string {
	return fmt.Sprintf("ServiceCommand {SourceDir: %v, DestDir: %v, CreateDest: %v, WithVars: [%v], WithList: [%v]}", copyCmd.SourceDir, copyCmd.DestinationDir, strconv.FormatBool(copyCmd.CreateDest), copyCmd.WithVars, copyCmd.WithList)
}

func (copyCmd *CopyCommand) Convert(cmdValues interface{}) (interface{}, error) {
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
	var sourceDir, destDir string
	var withVars []string = make([]string, 0)
	var withList []string = make([]string, 0)
	var createDest bool = false
	var valType string = fmt.Sprintf("%v", cmdValues)
	if len(valType) > 3 && "map" == valType[0:3] {
		for key, value := range cmdValues.(map[string]interface{}) {
			var elemValType string = fmt.Sprintf("%v", value)
			if strings.ToLower(key) == "srcDir" {
				if elemValType == "string" {
					sourceDir = fmt.Sprintf("%v", value)
				} else {
					return nil, errors.New("Unable to parse command: service.srcDir, with aguments of type " + elemValType + ", expected type string")
				}
			} else if strings.ToLower(key) == "destDir" {
				if elemValType == "string" {
					destDir = fmt.Sprintf("%v", value)
				} else {
					return nil, errors.New("Unable to parse command: service.destDir, with aguments of type " + elemValType + ", expected type string")
				}
			} else if strings.ToLower(key) == "createIfMissing" {
				if elemValType == "string" {
					bl, err := strconv.ParseBool(fmt.Sprintf("%v", value))
					if err != nil {
						return nil, errors.New("Error parsing command: shell.createIfMissing, cause: " + err.Error())

					} else {
						createDest = bl
					}

				} else if elemValType == "bool" {
					createDest = value.(bool)
				} else {
					return nil, errors.New("Unable to parse command: shell.createIfMissing, with aguments of type " + elemValType + ", expected type bool or string")
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
	return CopyCommand{
		SourceDir:      sourceDir,
		DestinationDir: destDir,
		CreateDest:     createDest,
		WithVars:       withVars,
		WithList:       withList,
	}, nil
}

func NewConverter(cmdType int) Converter {
	switch cmdType {
	case FEED_TYPE_SHELL:
		return &ShellCommand{}
	case FEED_TYPE_SERVICE:
		return &ServiceCommand{}
	case FEED_TYPE_COPY:
		return &CopyCommand{}
	}
	return &NilCommandConverter{
		CmdType: cmdType,
	}

}
