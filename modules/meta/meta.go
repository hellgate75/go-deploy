package meta

import (
	"github.com/hellgate75/go-deploy/types/module"
	"github.com/hellgate75/go-deploy/types/threads"
	"github.com/hellgate75/go-tcp-common/log"
)

type Symbol interface{}

// Get Required system logger verbosity
func GetVerbosity() string {
	if module.RuntimeDeployConfig != nil {
		return module.RuntimeDeployConfig.LogVerbosity
	} else {
		return "INFO"
	}

}


// Coverter interface, responsible to comvert raw interface from the parsing to a specific structure
type Converter interface {
	// Converts a raw interface element to a command qualified structure <BR/>
	// Paramameters: <BR/>
	// cmdValues (interface{}) Raw value from the feed file parsing
	// Return: <BR/>
	// (interface{}) Qualified structure <BR/>
	// (error) Error occured during any conversion <BR/>
	Convert(cmdValues interface{}) (threads.StepRunnable, error)

	// Set locally into the converter the logger to use
	SetLogger(l log.Logger)
}

// Interface that defines the Proxy Stub Components behaviors
type ProxyStub interface {
	Discover(module string) (Converter, error)
}
