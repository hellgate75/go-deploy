package threads

import (
	"github.com/hellgate75/go-deploy/types/defaults"
	"github.com/hellgate75/go-deploy/types/module"
	"github.com/hellgate75/go-deploy/worker/pool"
)

// Step specific Runnable interface, implementing pool.Runnable
type StepRunnable interface {
	pool.Runnable
	// Clone the runnable
	Clone() StepRunnable
	//Set Host target
	SetHost(host defaults.HostValue)
	//Set specific session component
	SetSession(session module.Session)
	//Set configuration data
	SetConfig(config defaults.ConfigPattern)
}
