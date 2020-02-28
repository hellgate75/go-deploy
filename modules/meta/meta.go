package meta

import (
	"github.com/hellgate75/go-deploy/types/module"
)

type Symbol interface{}

func GetVerbosity() string {
	if module.RuntimeDeployConfig != nil {
		return module.RuntimeDeployConfig.LogVerbosity
	} else {
		return "INFO"
	}

}

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
	Execute(step module.Step, session module.Session) error
}

type ProxyStub interface {
	Discover(module string, component string) (interface{}, error)
}
