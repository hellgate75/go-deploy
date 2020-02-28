package cmd

import (
	"errors"
	"fmt"
	"github.com/hellgate75/go-deploy/io"
	"github.com/hellgate75/go-deploy/log"
	"github.com/hellgate75/go-deploy/types/defaults"
	"github.com/hellgate75/go-deploy/types/module"
	"os"
	"runtime"
	"strings"
)

const (
	DEPLOY_CONFIG_FILE_NAME string = "deploy-config"
	DEPLOY_DATA_FILE_NAME   string = "deploy-data"
	DEPLOY_NET_FILE_NAME    string = "deploy-net"
	DEPLOY_ENVS_FILE_NAME   string = "deploy-envs"
	DEFAULT_CONFIG_FOLDER   string = "env"
	DEFAULT_CHARTS_FOLDER   string = "charts"
	DEFAULT_MODULES_FOLDER  string = "mod"
	DEFAULT_SYSTEM_FOLDER   string = ".go-deploy"
)

var Logger log.Logger = nil

type Bootstrap interface {
	Init(baseDir string, suffix string, format module.DescriptorTypeValue, logger log.Logger) []error
	Load(baseDir string, suffix string, format module.DescriptorTypeValue, logger log.Logger) []error
	Run(feed *module.FeedExec, logger log.Logger) []error
	GetDeployConfig() *module.DeployConfig
	GetDeployType() *module.DeployType
	GetNetType() *module.NetProtocolType
	GetDefaultDeployConfig() *module.DeployConfig
	GetDefaultDeployType() *module.DeployType
	GetDefaultNetType() *module.NetProtocolType
}

type bootstrap struct {
	deployConfig *module.DeployConfig
	deployType   *module.DeployType
	netType      *module.NetProtocolType
}

func (bootstrap *bootstrap) GetDeployConfig() *module.DeployConfig {
	return bootstrap.deployConfig
}

func (bootstrap *bootstrap) GetDeployType() *module.DeployType {
	return bootstrap.deployType
}

func (bootstrap *bootstrap) GetNetType() *module.NetProtocolType {
	return bootstrap.netType
}

func (bootstrap *bootstrap) GetDefaultDeployConfig() *module.DeployConfig {
	dt, err := ParseArguments()
	if err != nil {
		return nil
	}
	return dt
}

func loadVarsFiles() ([]defaults.NameValue, error) {
	var varList []defaults.NameValue = make([]defaults.NameValue, 0)
	var configDir string = module.RuntimeDeployConfig.ConfigDir
	var varfiles []string = module.RuntimeDeployConfig.UseVars
	for _, varFile := range varfiles {
		var index int = strings.Index(varFile, ".")
		var ext string = ""
		if index < 1 {
			ext = strings.ToLower(string(module.RuntimeDeployConfig.ConfigLang))
			varFile = varFile + "." + ext
			index = strings.Index(varFile, ".")
			//return varList, errors.New("bootstrap.loadVarsFiles -> Unable to discover extension for file:" + varFile)
		}
		ext = strings.ToLower(varFile[index+1:])
		if ext != "yml" && ext != "yaml" && ext != "xml" && ext != "json" {
			return varList, errors.New("bootstrap.loadVarsFiles -> Unable to parser for extension:" + ext)
		}
		var varFileFullPath = configDir + io.GetPathSeparator() + varFile
		Logger.Info("Loading vars file: " + varFileFullPath)
		var varsFileObj *defaults.Vars = &defaults.Vars{}
		var err error
		if ext == "yml" || ext == "yaml" {
			varsFileObj, err = varsFileObj.FromYamlFile(varFileFullPath)
			if err != nil {
				return varList, errors.New("bootstrap.loadVarsFiles -> Cause:" + err.Error())
			}
		} else if ext == "json" {
			varsFileObj, err = varsFileObj.FromJsonFile(varFileFullPath)
			if err != nil {
				return varList, errors.New("bootstrap.loadVarsFiles -> Cause:" + err.Error())
			}
		} else if ext == "xml" {
			varsFileObj, err = varsFileObj.FromXmlFile(varFileFullPath)
			if err != nil {
				return varList, errors.New("bootstrap.loadVarsFiles -> Cause:" + err.Error())
			}
		}
		varList = append(varList, varsFileObj.Vars...)
	}
	return varList, nil
}

func loadHostsFiles() ([]defaults.HostValue, error) {
	var hostsList []defaults.HostValue = make([]defaults.HostValue, 0)
	var configDir string = module.RuntimeDeployConfig.ConfigDir
	var hostsfiles []string = module.RuntimeDeployConfig.UseHosts
	for _, hostsFile := range hostsfiles {
		var index int = strings.Index(hostsFile, ".")
		var ext string = ""
		if index < 1 {
			ext = strings.ToLower(string(module.RuntimeDeployConfig.ConfigLang))
			hostsFile = hostsFile + "." + ext
			index = strings.Index(hostsFile, ".")
			//return hostsList, errors.New("bootstrap.loadHostsFiles -> Unable to discover extension for file:" + hostsFile)
		}
		ext = strings.ToLower(hostsFile[index+1:])
		if ext != "yml" && ext != "yaml" && ext != "xml" && ext != "json" {
			return hostsList, errors.New("bootstrap.loadHostsFiles -> Unable to parser for extension:" + ext)
		}
		var hostsFileFullPath = configDir + io.GetPathSeparator() + hostsFile
		Logger.Info("Loading hosts file: " + hostsFileFullPath)
		var hostsFileObj *defaults.Hosts = &defaults.Hosts{}
		var err error
		if ext == "yml" || ext == "yaml" {
			hostsFileObj, err = hostsFileObj.FromYamlFile(hostsFileFullPath)
			if err != nil {
				return hostsList, errors.New("bootstrap.loadHostsFiles -> Cause:" + err.Error())
			}
		} else if ext == "json" {
			hostsFileObj, err = hostsFileObj.FromJsonFile(hostsFileFullPath)
			if err != nil {
				return hostsList, errors.New("bootstrap.loadHostsFiles -> Cause:" + err.Error())
			}
		} else if ext == "xml" {
			hostsFileObj, err = hostsFileObj.FromXmlFile(hostsFileFullPath)
			if err != nil {
				return hostsList, errors.New("bootstrap.loadHostsFiles -> Cause:" + err.Error())
			}
		}
		hostsList = append(hostsList, hostsFileObj.Hosts...)
	}
	return hostsList, nil
}

func loadEnvsFile() ([]defaults.NameValue, error) {
	var envsList []defaults.NameValue = make([]defaults.NameValue, 0)
	var configDir string = module.RuntimeDeployConfig.ConfigDir
	var envFile = DEPLOY_ENVS_FILE_NAME
	var ext string = strings.ToLower(string(module.RuntimeDeployConfig.ConfigLang))
	var envsFileFullPath = configDir + io.GetPathSeparator() + envFile + "." + ext
	Logger.Info("Loading environment file: " + envsFileFullPath)
	var envsFileObj *defaults.Environments = &defaults.Environments{}
	var err error
	if ext == "yml" || ext == "yaml" {
		envsFileObj, err = envsFileObj.FromYamlFile(envsFileFullPath)
		if err != nil {
			return envsList, errors.New("bootstrap.loadEnvsFile -> Cause:" + err.Error())
		}
	} else if ext == "json" {
		envsFileObj, err = envsFileObj.FromJsonFile(envsFileFullPath)
		if err != nil {
			return envsList, errors.New("bootstrap.loadEnvsFile -> Cause:" + err.Error())
		}
	} else if ext == "xml" {
		envsFileObj, err = envsFileObj.FromXmlFile(envsFileFullPath)
		if err != nil {
			return envsList, errors.New("bootstrap.loadEnvsFile -> Cause:" + err.Error())
		}
	}
	envsList = append(envsList, envsFileObj.Envs...)
	return envsList, nil
}

func (bootstrap *bootstrap) Run(feed *module.FeedExec, logger log.Logger) []error {
	var errList []error = make([]error, 0)
	envs, errE := loadEnvsFile()
	hosts, errH := loadHostsFiles()
	if errH != nil {
		Logger.Warn("Unable to load hosts...")
		Logger.Warn("Reason:", errH)
		panic("Exit the procedure!!")
	}
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
	return errList
}

func (bootstrap *bootstrap) GetDefaultDeployType() *module.DeployType {
	return &module.DeployType{
		DeploymentType: module.FILE_SOURCE,
		DescriptorType: module.YAML_DESCRIPTOR,
		Method:         "",
		PostBody:       "",
		StrategyType:   module.ONE_SHOT_DEPLOYMENT,
	}
}

func (bootstrap *bootstrap) GetDefaultNetType() *module.NetProtocolType {
	return &module.NetProtocolType{
		NetProtocol: module.NET_PROTOCOL_SSH,
		UserName:    "docker",
		KeyFile:     userHomeDir() + io.GetPathSeparator() + ".ssh" + io.GetPathSeparator() + "id_rsa.pub",
	}
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

func (bootstrap *bootstrap) Init(baseDir string, suffix string, format module.DescriptorTypeValue, logger log.Logger) []error {
	if baseDir == "" {
		baseDir = "./" + DEFAULT_CONFIG_FOLDER
	}
	var suffixString string = ""
	if suffix != "" {
		suffixString = "-" + suffix
	}
	var errors []error = make([]error, 0)
	var matcher func(string) bool = getMatcher(format)

	var configFileObjectList []*module.DeployConfig = make([]*module.DeployConfig, 0)

	var configFileList []string = io.FindFilesIn(baseDir, true, DEPLOY_CONFIG_FILE_NAME+suffixString)
	for _, configFilePath := range configFileList {
		logger.Debug("configFilePath:" + configFilePath)
		if io.IsFolder(configFilePath) {
			var files []string = io.GetMatchedFiles(configFilePath, true, matcher)
			for _, configFilePathX := range files {
				var config *module.DeployConfig = &module.DeployConfig{}
				var errX error = nil
				if format == module.YAML_DESCRIPTOR {
					config, errX = config.FromYamlFile(configFilePathX)
				} else if format == module.XML_DESCRIPTOR {
					config, errX = config.FromXmlFile(configFilePathX)
				} else if format == module.JSON_DESCRIPTOR {
					config, errX = config.FromJsonFile(configFilePathX)
				}
				if errX != nil {
					errors = append(errors, errX)
				} else {
					configFileObjectList = append(configFileObjectList, config)
				}
			}
		} else {
			var config *module.DeployConfig = &module.DeployConfig{}
			var errX error = nil
			if format == module.YAML_DESCRIPTOR {
				config, errX = config.FromYamlFile(configFilePath)
			} else if format == module.XML_DESCRIPTOR {
				config, errX = config.FromXmlFile(configFilePath)
			} else if format == module.JSON_DESCRIPTOR {
				config, errX = config.FromJsonFile(configFilePath)
			}
			if errX != nil {
				errors = append(errors, errX)
			} else {
				configFileObjectList = append(configFileObjectList, config)
			}
		}
	}

	var deployConfig *module.DeployConfig = nil

	for _, deployConfigX := range configFileObjectList {
		if deployConfig == nil {
			deployConfig = deployConfigX
		} else {
			deployConfig = deployConfig.Merge(deployConfigX)
		}
	}

	bootstrap.deployConfig = deployConfig

	return errors
}

func (bootstrap *bootstrap) Load(baseDir string, suffix string, format module.DescriptorTypeValue, logger log.Logger) []error {
	if baseDir == "" {
		baseDir = "./" + DEFAULT_CONFIG_FOLDER
	}
	var suffixString string = ""
	if suffix != "" {
		suffixString = "-" + suffix
	}
	var errors []error = make([]error, 0)
	var matcher func(string) bool = getMatcher(format)

	var dataFileObjectList []*module.DeployType = make([]*module.DeployType, 0)
	var netFileObjectList []*module.NetProtocolType = make([]*module.NetProtocolType, 0)

	var dataFileList []string = io.FindFilesIn(baseDir, true, DEPLOY_DATA_FILE_NAME+suffixString)
	for _, dataFilePath := range dataFileList {
		logger.Debug("dataFilePath:" + dataFilePath)
		if io.IsFolder(dataFilePath) {
			var files []string = io.GetMatchedFiles(dataFilePath, true, matcher)
			for _, dataFilePathX := range files {
				var dType *module.DeployType = &module.DeployType{}
				var errX error = nil
				if format == module.YAML_DESCRIPTOR {
					dType, errX = dType.FromYamlFile(dataFilePathX)
				} else if format == module.XML_DESCRIPTOR {
					dType, errX = dType.FromXmlFile(dataFilePathX)
				} else if format == module.JSON_DESCRIPTOR {
					dType, errX = dType.FromJsonFile(dataFilePathX)
				}
				if errX != nil {
					errors = append(errors, errX)
				} else {
					dataFileObjectList = append(dataFileObjectList, dType)
				}
			}
		} else {
			var dType *module.DeployType = &module.DeployType{}
			var errX error = nil
			if format == module.YAML_DESCRIPTOR {
				dType, errX = dType.FromYamlFile(dataFilePath)
			} else if format == module.XML_DESCRIPTOR {
				dType, errX = dType.FromXmlFile(dataFilePath)
			} else if format == module.JSON_DESCRIPTOR {
				dType, errX = dType.FromJsonFile(dataFilePath)
			}
			if errX != nil {
				errors = append(errors, errX)
			} else {
				dataFileObjectList = append(dataFileObjectList, dType)
			}
		}
	}

	var netFileList []string = io.FindFilesIn(baseDir, true, DEPLOY_NET_FILE_NAME+suffixString)
	for _, netFilePath := range netFileList {
		logger.Debug("netFilePath:" + netFilePath)
		if io.IsFolder(netFilePath) {
			var files []string = io.GetMatchedFiles(netFilePath, true, matcher)
			for _, netFilePathX := range files {
				var nType *module.NetProtocolType = &module.NetProtocolType{}
				var errX error = nil
				if format == module.YAML_DESCRIPTOR {
					nType, errX = nType.FromYamlFile(netFilePathX)
				} else if format == module.XML_DESCRIPTOR {
					nType, errX = nType.FromXmlFile(netFilePathX)
				} else if format == module.JSON_DESCRIPTOR {
					nType, errX = nType.FromJsonFile(netFilePathX)
				}
				if errX != nil {
					errors = append(errors, errX)
				} else {
					netFileObjectList = append(netFileObjectList, nType)
				}
			}
		} else {
			var nType *module.NetProtocolType = &module.NetProtocolType{}
			var errX error = nil
			if format == module.YAML_DESCRIPTOR {
				nType, errX = nType.FromYamlFile(netFilePath)
			} else if format == module.XML_DESCRIPTOR {
				nType, errX = nType.FromXmlFile(netFilePath)
			} else if format == module.JSON_DESCRIPTOR {
				nType, errX = nType.FromJsonFile(netFilePath)
			}
			if errX != nil {
				errors = append(errors, errX)
			} else {
				netFileObjectList = append(netFileObjectList, nType)
			}
		}
	}

	var deployType *module.DeployType = nil

	for _, deployTypeX := range dataFileObjectList {
		if deployType == nil {
			deployType = deployTypeX
		} else {
			deployType = deployType.Merge(deployTypeX)
		}
	}

	bootstrap.deployType = deployType

	var netType *module.NetProtocolType = nil

	for _, netTypeX := range netFileObjectList {
		if netType == nil {
			netType = netTypeX
		} else {
			netType = netType.Merge(netTypeX)
		}
	}

	bootstrap.netType = netType

	return errors
}

func NewBootStrap() Bootstrap {
	return &bootstrap{}
}

func getMatcher(format module.DescriptorTypeValue) func(string) bool {
	return func(name string) bool {
		if format == module.JSON_DESCRIPTOR {
			if idx := strings.Index(name, "."); idx > 0 {
				return strings.ToLower(name[idx+1:]) == "json"
			}
		} else if format == module.XML_DESCRIPTOR {
			if idx := strings.Index(name, "."); idx > 0 {
				return strings.ToLower(name[idx+1:]) == "xml"
			}
		} else if format == module.YAML_DESCRIPTOR {
			if idx := strings.Index(name, "."); idx > 0 {
				return strings.ToLower(name[idx+1:]) == "yml" ||
					strings.ToLower(name[idx+1:]) == "yaml"
			}
		}
		return false
	}
}
