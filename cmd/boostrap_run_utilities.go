package cmd

import (
	"errors"
	"github.com/hellgate75/go-tcp-common/io"
	"github.com/hellgate75/go-deploy/types/defaults"
	"github.com/hellgate75/go-deploy/types/module"
	"strings"
)

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
		Logger.Debug("Loading vars file: " + varFileFullPath)
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

func loadHostsFiles() ([]defaults.HostGroups, error) {
	var hostsList []defaults.HostGroups = make([]defaults.HostGroups, 0)
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
		Logger.Debug("Loading hosts file: " + hostsFileFullPath)
		var hostsFileObj *defaults.HostGroupsConfig = &defaults.HostGroupsConfig{}
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
		hostsList = append(hostsList, hostsFileObj.HostGroups...)
	}
	return hostsList, nil
}

func loadEnvsFile() ([]defaults.NameValue, error) {
	var envsList []defaults.NameValue = make([]defaults.NameValue, 0)
	var configDir string = module.RuntimeDeployConfig.ConfigDir
	var envFile = DEPLOY_ENVS_FILE_NAME
	var ext string = strings.ToLower(string(module.RuntimeDeployConfig.ConfigLang))
	var envsFileFullPath = configDir + io.GetPathSeparator() + envFile + "." + ext
	Logger.Debug("Loading environment file: " + envsFileFullPath)
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
