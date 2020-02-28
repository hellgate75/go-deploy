package cmd

import (
	"github.com/hellgate75/go-deploy/io"
	"github.com/hellgate75/go-deploy/log"
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

func (bootstrap *bootstrap) Run(feed *module.FeedExec, logger log.Logger) []error {
	var errList []error = make([]error, 0)

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
