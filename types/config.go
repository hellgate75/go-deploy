package types

import (

)

type DeployConfig struct {
	DeployName	string
	UseHosts	[]string
	UseVars		[]string
	ConfigDir	string
}
