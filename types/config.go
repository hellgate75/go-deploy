package types

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type DeployConfig struct {
	DeployName string
	UseHosts   []string
	UseVars    []string
	ConfigDir  string
}

func (dc *DeployConfig) Yaml() string {
	bytes, err := yaml.Marshal(dc)
	if err != nil {
		return "Error: " + err.Error()
	} else {
		return fmt.Sprintf("\n%s", bytes)
	}
}
