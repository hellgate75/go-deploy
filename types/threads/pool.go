package threads

import (
	"github.com/hellgate75/go-deploy/net/generic"
	"github.com/hellgate75/go-deploy/types/defaults"
	"github.com/hellgate75/go-deploy/types/module"
	"github.com/hellgate75/go-tcp-common/pool"
)

// Step specific Runnable interface, implementing pool.Runnable
type StepRunnable interface {
	pool.Runnable
	// Clone the runnable
	Clone() StepRunnable
	//Set Host target
	SetClient(client generic.NetworkClient)
	//Set Host target
	SetHost(host defaults.HostValue)
	//Set specific session component
	SetSession(session module.Session)
	//Set configuration data
	SetConfig(config defaults.ConfigPattern)
	// Verify equality between StepRunnable instances
	Equals(r StepRunnable) bool
}
