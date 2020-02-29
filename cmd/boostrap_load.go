package cmd

import (
	"github.com/hellgate75/go-deploy/io"
	"github.com/hellgate75/go-deploy/log"
	"github.com/hellgate75/go-deploy/types/module"
)

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
