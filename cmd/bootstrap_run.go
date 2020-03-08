package cmd

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/hellgate75/go-deploy/io"
	"github.com/hellgate75/go-tcp-common/log"
	"github.com/hellgate75/go-deploy/net"
	"github.com/hellgate75/go-deploy/net/generic"
	"github.com/hellgate75/go-deploy/types/defaults"
	"github.com/hellgate75/go-deploy/types/module"
	"github.com/hellgate75/go-deploy/worker"
)

func (bootstrap *bootstrap) Run(feed *module.FeedExec, logger log.Logger) []error {
	var errList []error = make([]error, 0)
	hosts, errH := loadHostsFiles()
	if errH != nil || len(hosts) == 0 {
		Logger.Error("Unable to load hosts...")
		Logger.Error("Reason:", errH)
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
	configYaml, _ := module.RuntimeDeployConfig.Yaml()
	typeYaml, _ := module.RuntimeDeployType.Yaml()
	netYaml, _ := module.RuntimeNetworkType.Yaml()
	Logger.Debugf("Loaded:\nEnvironments: %s\nHosts: %s\nVariables: %s", envsYaml, hostsYaml, varsYaml)
	Logger.Debugf("\nConfig: %s\nType: %s\nNet: %s", configYaml, typeYaml, netYaml)

	Logger.Info("Connection Protocol: " + color.Yellow.Render(string(module.RuntimeNetworkType.NetProtocol)))
	var handler generic.ConnectionHandler = nil
	var isGoTCPClient bool = false
	if string(module.RuntimeNetworkType.NetProtocol) == string(module.NET_PROTOCOL_SSH) {
		handler = net.NewSshConnectionHandler(module.RuntimeDeployConfig.SingleSession)
	} else if string(module.RuntimeNetworkType.NetProtocol) == string(module.NET_PROTOCOL_GO_DEPLOY_CLIENT) {
		handler = net.NewGoTCPConnectionHandler(module.RuntimeDeployConfig.SingleSession)
		isGoTCPClient = true
	} else {
		Logger.Error("Unable to determine the Connection Handler for: " + string(module.RuntimeNetworkType.NetProtocol))
		panic("Unable to determine the Connection Handler")
	}
	if handler == nil {
		var message string = fmt.Sprintf("Unable to create ConnectionHandler for type: %s", string(module.RuntimeNetworkType.NetProtocol))
		Logger.Error(message)
		panic(message)
	}
	Logger.Infof("Connection Handler loaded: %s", color.Yellow.Render(fmt.Sprintf("%v", (handler != nil))))
	if isGoTCPClient {
		Logger.Warn("Using experimental GoTcp protocol instead of SSH ...")
		var missCertificate bool = module.RuntimeNetworkType == nil || module.RuntimeNetworkType.Certificate == ""
		var missKey bool = module.RuntimeNetworkType == nil || module.RuntimeNetworkType.KeyFile == ""
		if missCertificate || missKey {
			var message string = "Missing mandatory authentication TLS client key or certificate"
			Logger.Error(message)
			panic(message)
		}
	} else {
		Logger.Warn("Using SSH protocol ...")
		var missKey bool = module.RuntimeNetworkType == nil || module.RuntimeNetworkType.KeyFile == ""
		//		var missPassPhrase bool = module.RuntimeNetworkType == nil || module.RuntimeNetworkType.Passphrase == ""
		var missUser bool = module.RuntimeNetworkType == nil || module.RuntimeNetworkType.UserName == ""
		var missPassword bool = module.RuntimeNetworkType == nil || module.RuntimeNetworkType.Password == ""
		if missUser {
			var message string = "Missing mandatory authentication SSH username"
			Logger.Error(message)
			panic(message)
		} else if missPassword && missKey {
			var message string = "Missing mandatory authentication SSH passoword and rsa public key file"
			Logger.Error(message)
			panic(message)
		}
	}

	var sessionsMap map[string]module.Session = make(map[string]module.Session)

	for _, hg := range hosts {
		for _, hostValue := range hg.Hosts {
			var hostSessionMapKey string = hg.Name + "-" + hostValue.Name
			sessionsMap[hostSessionMapKey] = module.NewSession(module.NewSessionId())
			Logger.Debugf("Create session for host: %s -> Session Id: %s", color.Yellow.Render(hostValue.Name), color.Yellow.Render(sessionsMap[hostSessionMapKey].GetSessionId()))
			for _, variable := range vars {
				Logger.Debugf("Create session variable for host: %s -> Name: %s  Value: %s", color.Yellow.Render(hostValue.Name), variable.Name, variable.Value)
				sessionsMap[hostSessionMapKey].SetVar(variable.Name, variable.Value)
			}
			sessionsMap[hostSessionMapKey].SetSystemObject("connection-handler", handler.Clone())
			sessionsMap[hostSessionMapKey].SetSystemObject("rutime-config", module.RuntimeDeployConfig)
			sessionsMap[hostSessionMapKey].SetSystemObject("runtime-type", module.RuntimeDeployType)
			sessionsMap[hostSessionMapKey].SetSystemObject("runtime-net", module.RuntimeNetworkType)
			sessionsMap[hostSessionMapKey].SetSystemObject("host-groups", hosts)
			sessionsMap[hostSessionMapKey].SetSystemObject("envs", envs)
			sessionsMap[hostSessionMapKey].SetSystemObject("vars", vars)
			sessionsMap[hostSessionMapKey].SetSystemObject("system-logger", logger)
		}
	}
	Logger.Info("Starting Feed execution ...")
	execErrList := worker.ExecuteFeed(defaults.ConfigPattern{
		Config:     module.RuntimeDeployConfig,
		Type:       module.RuntimeDeployType,
		Net:        module.RuntimeNetworkType,
		Envs:       envs,
		HostGroups: hosts,
		Vars:       vars,
	}, feed, sessionsMap, logger)
	if len(execErrList) > 0 {
		errList = append(errList, execErrList...)
	}
	return errList
}
