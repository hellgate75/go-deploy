package cmd

import (
	"errors"
	"fmt"
	"github.com/hellgate75/go-deploy/io"
	"github.com/hellgate75/go-tcp-common/log"
	"github.com/hellgate75/go-deploy/types/module"
)

// Loads Config Type, Network Type and Plugins Type files and merge them, saving in the Runtime package variables
func (bootstrap *bootstrap) Load(baseDir string, suffix string, format module.DescriptorTypeValue, logger log.Logger) []error {
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
			var message string = fmt.Sprintf("cmd.Bootstrap.Load - Recovery:\n- %v", r)
			Logger.Error(message)
			errorsList = append(errorsList, errors.New(fmt.Sprintf("%v", r)))
		}
	}()
	var matcher func(string) bool = getMatcher(format)

	var dataFileObjectList []*module.DeployType = make([]*module.DeployType, 0)
	var netFileObjectList []*module.NetProtocolType = make([]*module.NetProtocolType, 0)
	var pluginFileObjectList []*module.PluginsConfig = make([]*module.PluginsConfig, 0)

	var dataFileList []string = io.FindFilesIn(baseDir, true, DEPLOY_DATA_FILE_NAME+suffixString)
	if "" != suffixString && len(dataFileList) == 0 {
		dataFileList = io.FindFilesIn(baseDir, true, DEPLOY_DATA_FILE_NAME)
	}
	for _, dataFilePath := range dataFileList {
		logger.Debug("dataFilePath:" + dataFilePath)
		if io.IsFolder(dataFilePath) {
			var files []string = io.GetMatchedFiles(dataFilePath, true, matcher)
			for _, dataFilePathX := range files {
				var dType *module.DeployType = &module.DeployType{}
				dformat := GetFileFormatDescritor(dataFilePathX, format)
				var errX error = nil
				if dformat == module.YAML_DESCRIPTOR {
					dType, errX = dType.FromYamlFile(dataFilePathX)
				} else if dformat == module.XML_DESCRIPTOR {
					dType, errX = dType.FromXmlFile(dataFilePathX)
				} else if dformat == module.JSON_DESCRIPTOR {
					dType, errX = dType.FromJsonFile(dataFilePathX)
				}
				if errX != nil {
					errorsList = append(errorsList, errX)
				} else {
					dataFileObjectList = append(dataFileObjectList, dType)
				}
			}
		} else {
			var dType *module.DeployType = &module.DeployType{}
			dformat := GetFileFormatDescritor(dataFilePath, format)
			var errX error = nil
			if dformat == module.YAML_DESCRIPTOR {
				dType, errX = dType.FromYamlFile(dataFilePath)
			} else if dformat == module.XML_DESCRIPTOR {
				dType, errX = dType.FromXmlFile(dataFilePath)
			} else if dformat == module.JSON_DESCRIPTOR {
				dType, errX = dType.FromJsonFile(dataFilePath)
			}
			if errX != nil {
				errorsList = append(errorsList, errX)
			} else {
				dataFileObjectList = append(dataFileObjectList, dType)
			}
		}
	}

	var netFileList []string = io.FindFilesIn(baseDir, true, DEPLOY_NET_FILE_NAME+suffixString)
	if "" != suffixString && len(netFileList) == 0 {
		netFileList = io.FindFilesIn(baseDir, true, DEPLOY_NET_FILE_NAME)
	}
	for _, netFilePath := range netFileList {
		logger.Debug("netFilePath:" + netFilePath)
		if io.IsFolder(netFilePath) {
			var files []string = io.GetMatchedFiles(netFilePath, true, matcher)
			for _, netFilePathX := range files {
				var nType *module.NetProtocolType = &module.NetProtocolType{}
				var errX error = nil
				dformat := GetFileFormatDescritor(netFilePathX, format)
				if dformat == module.YAML_DESCRIPTOR {
					nType, errX = nType.FromYamlFile(netFilePathX)
				} else if dformat == module.XML_DESCRIPTOR {
					nType, errX = nType.FromXmlFile(netFilePathX)
				} else if dformat == module.JSON_DESCRIPTOR {
					nType, errX = nType.FromJsonFile(netFilePathX)
				}
				if errX != nil {
					errorsList = append(errorsList, errX)
				} else {
					netFileObjectList = append(netFileObjectList, nType)
				}
			}
		} else {
			var nType *module.NetProtocolType = &module.NetProtocolType{}
			var errX error = nil
			dformat := GetFileFormatDescritor(netFilePath, format)
			if dformat == module.YAML_DESCRIPTOR {
				nType, errX = nType.FromYamlFile(netFilePath)
			} else if dformat == module.XML_DESCRIPTOR {
				nType, errX = nType.FromXmlFile(netFilePath)
			} else if dformat == module.JSON_DESCRIPTOR {
				nType, errX = nType.FromJsonFile(netFilePath)
			}
			if errX != nil {
				errorsList = append(errorsList, errX)
			} else {
				netFileObjectList = append(netFileObjectList, nType)
			}
		}
	}
	
	
	var pluginsFileList []string = io.FindFilesIn(baseDir, true, DEPLOY_PKUGINS_FILE_NAME+suffixString)
	if "" != suffixString && len(pluginsFileList) == 0 {
		pluginsFileList = io.FindFilesIn(baseDir, true, DEPLOY_PKUGINS_FILE_NAME)
	}
	for _, pluginsFilePath := range pluginsFileList {
		logger.Debug("pluginsFilePath:" + pluginsFilePath)
		if io.IsFolder(pluginsFilePath) {
			var files []string = io.GetMatchedFiles(pluginsFilePath, true, matcher)
			for _, netFilePathX := range files {
				var nPlugins *module.PluginsConfig = &module.PluginsConfig{}
				var errX error = nil
				dformat := GetFileFormatDescritor(netFilePathX, format)
				if dformat == module.YAML_DESCRIPTOR {
					nPlugins, errX = nPlugins.FromYamlFile(netFilePathX)
				} else if dformat == module.XML_DESCRIPTOR {
					nPlugins, errX = nPlugins.FromXmlFile(netFilePathX)
				} else if dformat == module.JSON_DESCRIPTOR {
					nPlugins, errX = nPlugins.FromJsonFile(netFilePathX)
				}
				if errX != nil {
					errorsList = append(errorsList, errX)
				} else {
					pluginFileObjectList = append(pluginFileObjectList, nPlugins)
				}
			}
		} else {
			var nPlugins *module.PluginsConfig = &module.PluginsConfig{}
			var errX error = nil
			dformat := GetFileFormatDescritor(pluginsFilePath, format)
			if dformat == module.YAML_DESCRIPTOR {
				nPlugins, errX = nPlugins.FromYamlFile(pluginsFilePath)
			} else if dformat == module.XML_DESCRIPTOR {
				nPlugins, errX = nPlugins.FromXmlFile(pluginsFilePath)
			} else if dformat == module.JSON_DESCRIPTOR {
				nPlugins, errX = nPlugins.FromJsonFile(pluginsFilePath)
			}
			if errX != nil {
				errorsList = append(errorsList, errX)
			} else {
				pluginFileObjectList = append(pluginFileObjectList, nPlugins)
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
	
	
	var pluginsType *module.PluginsConfig= nil
	
	for _, pluginTypeX := range pluginFileObjectList {
		if deployType == nil {
			pluginsType = pluginTypeX
		} else {
			pluginsType = pluginsType.Merge(pluginTypeX)
		}
	}
	
	bootstrap.pluginsType = pluginsType
	
	return errorsList
}
