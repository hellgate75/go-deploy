package cmd

import (
	"errors"
	"fmt"
	"github.com/hellgate75/go-tcp-common/io"
	"github.com/hellgate75/go-tcp-common/log"
	"github.com/hellgate75/go-deploy/types/module"
)

// Loads Main Deploy Config files and merge them, saving in the Runtime package variables
func (bootstrap *bootstrap) Init(baseDir string, suffix string, format module.DescriptorTypeValue, logger log.Logger) []error {
	if baseDir == "" {
		baseDir = "./" + DEFAULT_CONFIG_FOLDER
	}
	var suffixString string = ""
	if suffix != "" {
		suffixString = "-" + suffix
	}
	var errorsList []error = make([]error, 0)
	defer func() {
		if r := recover(); r != nil {
			var message string = fmt.Sprintf("cmd.Bootstrap.Init - Recovery:\n- %v", r)
			Logger.Error(message)
			errorsList = append(errorsList, errors.New(fmt.Sprintf("%v", r)))
		}
	}()
	var matcher func(string) bool = getMatcher(format)

	var configFileObjectList []*module.DeployConfig = make([]*module.DeployConfig, 0)

	var configFileList []string = io.FindFilesIn(baseDir, true, DEPLOY_CONFIG_FILE_NAME+suffixString)
	if "" != suffixString && len(configFileList) == 0 {
		configFileList = io.FindFilesIn(baseDir, true, DEPLOY_CONFIG_FILE_NAME)
	}
	for _, configFilePath := range configFileList {
		logger.Debug("configFilePath:" + configFilePath)
		if io.IsFolder(configFilePath) {
			var files []string = io.GetMatchedFiles(configFilePath, true, matcher)
			for _, configFilePathX := range files {
				dformat := GetFileFormatDescritor(configFilePathX, format)
				var config *module.DeployConfig = &module.DeployConfig{}
				var errX error = nil
				if dformat == module.YAML_DESCRIPTOR {
					config, errX = config.FromYamlFile(configFilePathX)
				} else if dformat == module.XML_DESCRIPTOR {
					config, errX = config.FromXmlFile(configFilePathX)
				} else if dformat == module.JSON_DESCRIPTOR {
					config, errX = config.FromJsonFile(configFilePathX)
				}
				if errX != nil {
					errorsList = append(errorsList, errX)
				} else {
					configFileObjectList = append(configFileObjectList, config)
				}
			}
		} else {
			var config *module.DeployConfig = &module.DeployConfig{}
			var errX error = nil
			dformat := GetFileFormatDescritor(configFilePath, format)
			if dformat == module.YAML_DESCRIPTOR {
				config, errX = config.FromYamlFile(configFilePath)
			} else if dformat == module.XML_DESCRIPTOR {
				config, errX = config.FromXmlFile(configFilePath)
			} else if dformat == module.JSON_DESCRIPTOR {
				config, errX = config.FromJsonFile(configFilePath)
			}
			if errX != nil {
				errorsList = append(errorsList, errX)
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

	return errorsList
}
