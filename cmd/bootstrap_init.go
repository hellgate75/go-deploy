package cmd

import (
	"github.com/hellgate75/go-deploy/io"
	"github.com/hellgate75/go-deploy/log"
	"github.com/hellgate75/go-deploy/types/module"
)

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
