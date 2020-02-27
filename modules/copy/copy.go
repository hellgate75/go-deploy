package copy

import (
	"errors"
	"fmt"
	"github.com/hellgate75/go-deploy/log"
	"github.com/hellgate75/go-deploy/modules"
	"github.com/hellgate75/go-deploy/types/module"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

var Logger log.Logger = log.NewLogger(log.VerbosityLevelFromString(module.RuntimeDeployConfig.LogVerbosity))

type copyCmdExecutor struct {
}

func (shell *copyCmdExecutor) Execute(step module.Step) error {
	Logger.Warn(fmt.Sprintf("copy.Executor.Execute -> Executing command : %s", step.StepType))
	return nil
}

var Executor modules.Executor = &copyCmdExecutor{}

var ERROR_TYPE reflect.Type = reflect.TypeOf(errors.New(""))

/*
* Service command structure
 */
type copyCommand struct {
	SourceDir      string
	DestinationDir string
	CreateDest     bool
	WithVars       []string
	WithList       []string
}

func (copyCmd copyCommand) String() string {
	return fmt.Sprintf("ServiceCommand {SourceDir: %v, DestDir: %v, CreateDest: %v, WithVars: [%v], WithList: [%v]}", copyCmd.SourceDir, copyCmd.DestinationDir, strconv.FormatBool(copyCmd.CreateDest), copyCmd.WithVars, copyCmd.WithList)
}

func (copyCmd *copyCommand) Convert(cmdValues interface{}) (interface{}, error) {
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
	return copyCommand{
		SourceDir:      sourceDir,
		DestinationDir: destDir,
		CreateDest:     createDest,
		WithVars:       withVars,
		WithList:       withList,
	}, nil
}

var Converter modules.Converter = &copyCommand{}

//go:linkname modules_seekModule [github.com/hellgate75/go-deploy/modules.seekModules]
func modules_seekModule(module string, feature string, acceptance chan bool, featureAcceptance chan bool, response chan interface{}) {
	if module == "copy" {
		acceptance <- true
		if feature == "Converter" {
			featureAcceptance <- true
			response <- &copyCommand{}
		} else if feature == "Executor" {
			featureAcceptance <- true
			response <- &copyCmdExecutor{}
		} else {
			featureAcceptance <- false
			response <- nil
		}

	} else {
		acceptance <- false
	}
}

// doInit is defined in package runtime
//go:linkname doInit runtime.doInit
func doInit(t unsafe.Pointer) // t should be a *runtime.initTask

func main() {}
