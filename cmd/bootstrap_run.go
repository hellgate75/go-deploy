package cmd

import (
	"errors"
	"fmt"
	"github.com/gookit/color"
	"github.com/hellgate75/go-tcp-common/io"
	"github.com/hellgate75/go-deploy/net"
	"github.com/hellgate75/go-deploy/types/defaults"
	"github.com/hellgate75/go-deploy/types/module"
	"github.com/hellgate75/go-deploy/worker"
	"github.com/hellgate75/go-tcp-client/client/proxy"
	modproxy "github.com/hellgate75/go-deploy/modules/proxy"
	"github.com/hellgate75/go-tcp-common/log"
)

// Start Deploy Process.
func (bootstrap *bootstrap) Run(feed *module.FeedExec, logger log.Logger) []error {
	var errorsList []error = make([]error, 0)
	defer func() {
		if r := recover(); r != nil {
			var message string = fmt.Sprintf("cmd.Bootstrap.Run - Recovery:\n- %v", r)
			Logger.Error(message)
			errorsList = append(errorsList, errors.New(fmt.Sprintf("%v", r)))
		}
	}()
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
	handlerEnvelope, errHandler := net.DiscoverConnectionHandler(string(module.RuntimeNetworkType.NetProtocol))
	if errHandler != nil {
		Logger.Error("Unable to determine the Connection Handler for: " + string(module.RuntimeNetworkType.NetProtocol))
		panic("Unable to determine the Connection Handler: " + errHandler.Error())
	}
	if handlerEnvelope == nil {
		var message string = fmt.Sprintf("Unable to create ConnectionHandler for type: %s", string(module.RuntimeNetworkType.NetProtocol))
		Logger.Error(message)
		panic(message)
	}
	handler, handlerConfig := handlerEnvelope(module.RuntimeDeployConfig.SingleSession)
	if module.RuntimePluginsType.EnableDeployClientCommandsPlugin {
		handler.UsePlugins(module.RuntimePluginsType.DeployClientCommandsPluginExtension, module.RuntimePluginsType.DeployClientCommandsPluginFolder)
	}
	Logger.Infof("Connection Handler loaded: %s", color.Yellow.Render(fmt.Sprintf("%v", (handler != nil))))
	Logger.Warn("Using "+string(module.RuntimeNetworkType.NetProtocol)+" protocol ...")
	var missKey bool = module.RuntimeNetworkType == nil || module.RuntimeNetworkType.KeyFile == ""
	var missPassPhrase bool = module.RuntimeNetworkType == nil || module.RuntimeNetworkType.Passphrase == ""
	var missUser bool = module.RuntimeNetworkType == nil || module.RuntimeNetworkType.UserName == ""
	var missPassword bool = module.RuntimeNetworkType == nil || module.RuntimeNetworkType.Password == ""
	var missCertificate bool = module.RuntimeNetworkType == nil || module.RuntimeNetworkType.Certificate == ""
	var useUserPassword bool = false
	var useUserKey bool = false
	var useUserKeyPassphrase bool = false
	var useCertificates bool = false
	if !missUser && !missPassword && handlerConfig.UseUserPassword {
		useUserPassword = true
	} else if !missUser && !missKey && missPassPhrase && handlerConfig.UseAuthKey{
		useUserKey = true
	} else if !missUser && !missKey && !missPassPhrase && handlerConfig.UseAuthKeyPassphrase {
		useUserKeyPassphrase = true
	} else if !missKey && !missCertificate && handlerConfig.UseCertificates {
		useCertificates = true
	} else {
		var message string = "Missing mandatory authentication user and/or passoword and/or rsa public key / TLS Key file or certificates for client type: " + string(module.RuntimeNetworkType.NetProtocol)
		Logger.Error(message)
		panic(message)
	}
	
	var connectionConfig module.ConnectionConfig = module.ConnectionConfig{
		UseUserPassword: useUserPassword,
		UseUserKey: useUserKey,
		UseSSHConfig: false,
		UseUserKeyPassphrase: useUserKeyPassphrase,
		UseTLSCertificates: useCertificates,
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
	if module.RuntimePluginsType.EnableDeployClientsPlugin {
		//modproxy "github.com/hellgate75/go-deploy/modules/proxy"
		proxy.UsePlugins = module.RuntimePluginsType.EnableDeployClientsPlugin
		proxy.PluginLibrariesExtension = module.RuntimePluginsType.DeployClientsPluginExtension
		proxy.PluginLibrariesFolder = module.RuntimePluginsType.DeployClientsPluginFolder
	}

	if module.RuntimePluginsType.EnableDeployCommandsPlugin{
		modproxy.UsePlugins = module.RuntimePluginsType.EnableDeployCommandsPlugin
		modproxy.PluginLibrariesExtension = module.RuntimePluginsType.DeployCommandsPluginExtension
		modproxy.PluginLibrariesFolder = module.RuntimePluginsType.DeployCommandsPluginFolder
	}
	Logger.Info("Starting Feed execution ...")
	execErrList := worker.ExecuteFeed(connectionConfig, defaults.ConfigPattern{
		Config:     module.RuntimeDeployConfig,
		Type:       module.RuntimeDeployType,
		Net:        module.RuntimeNetworkType,
		Plugins:    module.RuntimePluginsType,
		Envs:       envs,
		HostGroups: hosts,
		Vars:       vars,
	}, feed, sessionsMap, logger)
	if len(execErrList) > 0 {
		errorsList = append(errorsList, execErrList...)
	}
	return errorsList
}

