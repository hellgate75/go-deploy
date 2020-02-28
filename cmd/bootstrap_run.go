package cmd

import (
	"fmt"
	"github.com/hellgate75/go-deploy/io"
	"github.com/hellgate75/go-deploy/log"
	"github.com/hellgate75/go-deploy/types/module"
)

func (bootstrap *bootstrap) Run(feed *module.FeedExec, logger log.Logger) []error {
	var errList []error = make([]error, 0)
	hosts, errH := loadHostsFiles()
	if errH != nil {
		Logger.Warn("Unable to load hosts...")
		Logger.Warn("Reason:", errH)
		panic("Exit the procedure!!")
	}
	envs, errE := loadEnvsFile()
	if errE != nil {
		Logger.Warn("Unable to load environments...")
		Logger.Warn("Reason:", errE)
		Logger.Warn("We trust whatever you pass as environment!!")
	}
	vars, errV := loadVarsFiles()
	if errV != nil {
		Logger.Warn("Unable to load Vars...")
		Logger.Warn("Reason:", errV)
		Logger.Warn("Continue without any initial Variable!!")
	}
	envsYaml, _ := io.ToYaml(envs)
	hostsYaml, _ := io.ToYaml(hosts)
	varsYaml, _ := io.ToYaml(vars)
	Logger.Info(fmt.Sprintf("Loaded:\nEnvironments: %s\nHosts: %s\nVariables: %s", envsYaml, hostsYaml, varsYaml))
	var sessionsMap map[string]module.Session = make(map[string]module.Session)
	for _, host := range hosts {
		sessionsMap[host.Name] = module.NewSession(module.NewSessionId())
		Logger.Info(fmt.Sprintf("Create session for host: %s -> Session Id: %s", host.Name, sessionsMap[host.Name].GetSessionId()))
		for _, variable := range vars {
			sessionsMap[host.Name].SetVar(variable.Name, variable.Value)
		}
	}
	return errList
}
