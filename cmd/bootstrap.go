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
