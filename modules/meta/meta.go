package meta

import (
	"github.com/hellgate75/go-deploy/types/module"
	"github.com/hellgate75/go-deploy/types/threads"
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
	Convert(cmdValues interface{}) (threads.StepRunnable, error)
}

type ProxyStub interface {
	Discover(module string) (Converter, error)
}
