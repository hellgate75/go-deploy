package worker

import (
	"github.com/hellgate75/go-deploy/types/defaults"
	"github.com/hellgate75/go-deploy/types/module"
	"os"
	"runtime"
)

type ConfigPattern struct {
	Config     *module.DeployConfig
	Type       *module.DeployType
	Net        *module.NetProtocolType
	Envs       []defaults.NameValue
	Vars       []defaults.NameValue
	HostGroups []defaults.HostGroups
}

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	} else if runtime.GOOS == "linux" {
		home := os.Getenv("XDG_CONFIG_HOME")
		if home != "" {
			return home
		}
	}
	return os.Getenv("HOME")
}

func exxecutionDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	return dir
}
