package cmd

import (
	"github.com/hellgate75/go-deploy/io"
	"github.com/hellgate75/go-deploy/log"
	"github.com/hellgate75/go-deploy/types"
	"os"
	"runtime"
	"strings"
)

const (
	DEPLOY_CONFIG_FILE_NAME string = "deploy-config"
	DEPLOY_DATA_FILE_NAME   string = "deploy-data"
	DEPLOY_NET_FILE_NAME    string = "deploy-net"
)

type Bootstrap interface {
	Load(baseDir string, suffix string, format types.DescriptorTypeValue, logger log.Logger) []error
	GetDeployConfig() *types.DeployConfig
	GetDeployType() *types.DeployType
	GetNetType() *types.NetProtocolType
	GetDefaultDeployType() *types.DeployType
	GetDefaultNetType() *types.NetProtocolType
}

type bootstrap struct {
	deployConfig *types.DeployConfig
	deployType   *types.DeployType
	netType      *types.NetProtocolType
}

func (bootstrap *bootstrap) GetDeployConfig() *types.DeployConfig {
	return bootstrap.deployConfig
}

func (bootstrap *bootstrap) GetDeployType() *types.DeployType {
	return bootstrap.deployType
}

func (bootstrap *bootstrap) GetNetType() *types.NetProtocolType {
	return bootstrap.netType
}

func (bootstrap *bootstrap) GetDefaultDeployType() *types.DeployType {
	return &types.DeployType{
		DeploymentType: types.FILE_SOURCE,
		DescriptorType: types.YAML_DESCRIPTOR,
		Method:         "",
		PostBody:       "",
		StrategyType:   types.ONE_SHOT_DEPLOYMENT,
	}
}

func (bootstrap *bootstrap) GetDefaultNetType() *types.NetProtocolType {
	return &types.NetProtocolType{
		NetProtocol: types.NET_PROTOCOL_SSH,
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

func (bootstrap *bootstrap) Load(baseDir string, suffix string, format types.DescriptorTypeValue, logger log.Logger) []error {
	if baseDir == "" {
		baseDir = "./.godeploy"
	}
	var suffixString string = ""
	if suffix != "" {
		suffixString = "-" + suffix
	}
	var errors []error = make([]error, 0)
	var matcher func(string) bool = getMatcher(format)

	var configFileObjectList []*types.DeployConfig = make([]*types.DeployConfig, 0)
	var dataFileObjectList []*types.DeployType = make([]*types.DeployType, 0)
	var netFileObjectList []*types.NetProtocolType = make([]*types.NetProtocolType, 0)

	var configFileList []string = io.FindFilesIn(baseDir, true, DEPLOY_CONFIG_FILE_NAME+suffixString)
	for _, configFilePath := range configFileList {
		logger.Debug("configFilePath:" + configFilePath)
		if io.IsFolder(configFilePath) {
			var files []string = io.GetMatchedFiles(configFilePath, true, matcher)
			for _, configFilePathX := range files {
				var config *types.DeployConfig = &types.DeployConfig{}
				var errX error = nil
				if format == types.YAML_DESCRIPTOR {
					config, errX = config.FromYamlFile(configFilePathX)
				} else if format == types.XML_DESCRIPTOR {
					config, errX = config.FromXmlFile(configFilePathX)
				} else if format == types.JSON_DESCRIPTOR {
					config, errX = config.FromJsonFile(configFilePathX)
				}
				if errX != nil {
					errors = append(errors, errX)
				} else {
					configFileObjectList = append(configFileObjectList, config)
				}
			}
		} else {
			var config *types.DeployConfig = &types.DeployConfig{}
			var errX error = nil
			if format == types.YAML_DESCRIPTOR {
				config, errX = config.FromYamlFile(configFilePath)
			} else if format == types.XML_DESCRIPTOR {
				config, errX = config.FromXmlFile(configFilePath)
			} else if format == types.JSON_DESCRIPTOR {
				config, errX = config.FromJsonFile(configFilePath)
			}
			if errX != nil {
				errors = append(errors, errX)
			} else {
				configFileObjectList = append(configFileObjectList, config)
			}
		}
	}

	var dataFileList []string = io.FindFilesIn(baseDir, true, DEPLOY_DATA_FILE_NAME+suffixString)
	for _, dataFilePath := range dataFileList {
		logger.Debug("dataFilePath:" + dataFilePath)
		if io.IsFolder(dataFilePath) {
			var files []string = io.GetMatchedFiles(dataFilePath, true, matcher)
			for _, dataFilePathX := range files {
				var dType *types.DeployType = &types.DeployType{}
				var errX error = nil
				if format == types.YAML_DESCRIPTOR {
					dType, errX = dType.FromYamlFile(dataFilePathX)
				} else if format == types.XML_DESCRIPTOR {
					dType, errX = dType.FromXmlFile(dataFilePathX)
				} else if format == types.JSON_DESCRIPTOR {
					dType, errX = dType.FromJsonFile(dataFilePathX)
				}
				if errX != nil {
					errors = append(errors, errX)
				} else {
					dataFileObjectList = append(dataFileObjectList, dType)
				}
			}
		} else {
			var dType *types.DeployType = &types.DeployType{}
			var errX error = nil
			if format == types.YAML_DESCRIPTOR {
				dType, errX = dType.FromYamlFile(dataFilePath)
			} else if format == types.XML_DESCRIPTOR {
				dType, errX = dType.FromXmlFile(dataFilePath)
			} else if format == types.JSON_DESCRIPTOR {
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
				var nType *types.NetProtocolType = &types.NetProtocolType{}
				var errX error = nil
				if format == types.YAML_DESCRIPTOR {
					nType, errX = nType.FromYamlFile(netFilePathX)
				} else if format == types.XML_DESCRIPTOR {
					nType, errX = nType.FromXmlFile(netFilePathX)
				} else if format == types.JSON_DESCRIPTOR {
					nType, errX = nType.FromJsonFile(netFilePathX)
				}
				if errX != nil {
					errors = append(errors, errX)
				} else {
					netFileObjectList = append(netFileObjectList, nType)
				}
			}
		} else {
			var nType *types.NetProtocolType = &types.NetProtocolType{}
			var errX error = nil
			if format == types.YAML_DESCRIPTOR {
				nType, errX = nType.FromYamlFile(netFilePath)
			} else if format == types.XML_DESCRIPTOR {
				nType, errX = nType.FromXmlFile(netFilePath)
			} else if format == types.JSON_DESCRIPTOR {
				nType, errX = nType.FromJsonFile(netFilePath)
			}
			if errX != nil {
				errors = append(errors, errX)
			} else {
				netFileObjectList = append(netFileObjectList, nType)
			}
		}
	}

	var deployConfig *types.DeployConfig = nil

	for _, deployConfigX := range configFileObjectList {
		if deployConfig == nil {
			deployConfig = deployConfigX
		} else {
			deployConfig = deployConfig.Merge(deployConfigX)
		}
	}

	bootstrap.deployConfig = deployConfig

	var deployType *types.DeployType = nil

	for _, deployTypeX := range dataFileObjectList {
		if deployType == nil {
			deployType = deployTypeX
		} else {
			deployType = deployType.Merge(deployTypeX)
		}
	}

	bootstrap.deployType = deployType

	var netType *types.NetProtocolType = nil

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

func getMatcher(format types.DescriptorTypeValue) func(string) bool {
	return func(name string) bool {
		if format == types.JSON_DESCRIPTOR {
			if idx := strings.Index(name, "."); idx > 0 {
				return strings.ToLower(name[idx+1:]) == "json"
			}
		} else if format == types.XML_DESCRIPTOR {
			if idx := strings.Index(name, "."); idx > 0 {
				return strings.ToLower(name[idx+1:]) == "xml"
			}
		} else if format == types.YAML_DESCRIPTOR {
			if idx := strings.Index(name, "."); idx > 0 {
				return strings.ToLower(name[idx+1:]) == "yml" ||
					strings.ToLower(name[idx+1:]) == "yaml"
			}
		}
		return false
	}
}
