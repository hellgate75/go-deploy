package module

import (
	"fmt"
	"github.com/hellgate75/go-deploy/io"
	"github.com/hellgate75/go-deploy/utils"
)

var RuntimeDeployConfig *DeployConfig = nil
var RuntimeDeployType *DeployType = nil
var RuntimeNetworkType *NetProtocolType = nil
var RuntimePluginsType *PluginsConfig = nil

var ChartsDescriptorFormat DescriptorTypeValue = DescriptorTypeValue("YAML")

func (pc *PluginsConfig) Merge(pc2 *PluginsConfig) *PluginsConfig {
	return &PluginsConfig{
		EnableDeployClientCommandsPlugin: pc.EnableDeployClientCommandsPlugin || pc2.EnableDeployClientCommandsPlugin,
		DeployClientCommandsPluginExtension: bestString(pc2.DeployClientCommandsPluginExtension, pc.DeployClientCommandsPluginExtension),
		DeployClientCommandsPluginFolder: bestString(pc2.DeployClientCommandsPluginFolder, pc.DeployClientCommandsPluginFolder),
		EnableDeployClientsPlugin: pc.EnableDeployClientsPlugin || pc2.EnableDeployClientsPlugin,
		DeployClientsPluginExtension: bestString(pc2.DeployClientsPluginExtension, pc.DeployClientsPluginExtension),
		DeployClientsPluginFolder: bestString(pc2.DeployClientsPluginFolder, pc.DeployClientsPluginFolder),
		EnableDeployCommandsPlugin: pc.EnableDeployCommandsPlugin || pc2.EnableDeployCommandsPlugin,
		DeployCommandsPluginExtension: bestString(pc2.DeployCommandsPluginExtension, pc.DeployCommandsPluginExtension),
		DeployCommandsPluginFolder: bestString(pc2.DeployCommandsPluginFolder, pc.DeployCommandsPluginFolder),
	}
}

func (pc *PluginsConfig) String() string {
	return fmt.Sprintf("PluginsConfig{EnableDeployClientCommandsPlugin: %v, DeployClientCommandsPluginExtension \"%s\", DeployClientCommandsPluginFolder: \"%s\", EnableDeployClientsPlugin: %v, DeployClientsPluginExtension: \"%s\", DeployClientsPluginFolder: \"%s\", EnableDeployCommandsPlugin: %v, DeployCommandsPluginExtension: \"%s\", DeployCommandsPluginFolder: \"%s\"}",
		pc.EnableDeployClientCommandsPlugin, pc.DeployClientCommandsPluginExtension, pc.DeployClientCommandsPluginFolder, pc.EnableDeployClientsPlugin,
		pc.DeployClientsPluginExtension, pc.DeployClientsPluginFolder, pc.EnableDeployCommandsPlugin, pc.DeployCommandsPluginExtension, pc.DeployCommandsPluginFolder)
}

func (pc *PluginsConfig) Yaml() (string, error) {
	return io.ToYaml(pc)
}

func (pc *PluginsConfig) FromYamlFile(path string) (*PluginsConfig, error) {
	itf, err := io.FromYamlFile(path, pc)
	if err != nil {
		return nil, err
	}
	var conf *PluginsConfig = itf.(*PluginsConfig)
	return conf, nil
}

func (pc *PluginsConfig) FromYamlCode(yamlCode string) (*PluginsConfig, error) {
	itf, err := io.FromYamlCode(yamlCode, pc)
	if err != nil {
		return nil, err
	}
	var conf *PluginsConfig = itf.(*PluginsConfig)
	return conf, nil
}

func (pc *PluginsConfig) Xml() (string, error) {
	return io.ToXml(pc)
}

func (pc *PluginsConfig) FromXmlFile(path string) (*PluginsConfig, error) {
	itf, err := io.FromXmlFile(path, pc)
	if err != nil {
		return nil, err
	}
	var conf *PluginsConfig = itf.(*PluginsConfig)
	return conf, nil
}

func (pc *PluginsConfig) FromXmlCode(xmlCode string) (*PluginsConfig, error) {
	itf, err := io.FromXmlCode(xmlCode, pc)
	if err != nil {
		return nil, err
	}
	var conf *PluginsConfig = itf.(*PluginsConfig)
	return conf, nil
}

func (pc *PluginsConfig) Json() (string, error) {
	return io.ToJson(pc)
}

func (pc *PluginsConfig) FromJsonFile(path string) (*PluginsConfig, error) {
	itf, err := io.FromJsonFile(path, pc)
	if err != nil {
		return nil, err
	}
	var conf *PluginsConfig = itf.(*PluginsConfig)
	return conf, nil
}

func (pc *PluginsConfig) FromJsonCode(jsonCode string) (*PluginsConfig, error) {
	itf, err := io.FromJsonCode(jsonCode, pc)
	if err != nil {
		return nil, err
	}
	var conf *PluginsConfig = itf.(*PluginsConfig)
	return conf, nil
}

func (dt *DeployType) Merge(dt2 *DeployType) *DeployType {
	return &DeployType{
		DeploymentType: DeploymentTypeValue(bestString(string(dt2.DeploymentType), string(dt.DeploymentType))),
		DescriptorType: DescriptorTypeValue(bestString(string(dt2.DescriptorType), string(dt.DescriptorType))),
		StrategyType:   StrategyTypeValue(bestString(string(dt2.StrategyType), string(dt.StrategyType))),
		Method:         RestMethodTypeValue(bestString(string(dt2.Method), string(dt.Method))),
		Scheduled:      bestString(dt2.Scheduled, dt.Scheduled),
		PostBody:       bestString(dt2.PostBody, dt.PostBody),
	}
}

func (dt *DeployType) String() string {
	return fmt.Sprintf("DeployType{DeploymentType: \"%v\", DescriptorType: %v, StrategyType: %v, Method: \"%v\", Scheduled: \"%v\", PostBody: \"%v\"}",
		dt.DeploymentType, dt.DescriptorType, dt.StrategyType, dt.Method, dt.Scheduled, dt.PostBody)
}

func (dt *DeployType) Yaml() (string, error) {
	return io.ToYaml(dt)
}

func (dt *DeployType) FromYamlFile(path string) (*DeployType, error) {
	itf, err := io.FromYamlFile(path, dt)
	if err != nil {
		return nil, err
	}
	var conf *DeployType = itf.(*DeployType)
	return conf, nil
}

func (dt *DeployType) FromYamlCode(yamlCode string) (*DeployType, error) {
	itf, err := io.FromYamlCode(yamlCode, dt)
	if err != nil {
		return nil, err
	}
	var conf *DeployType = itf.(*DeployType)
	return conf, nil
}

func (dt *DeployType) Xml() (string, error) {
	return io.ToXml(dt)
}

func (dt *DeployType) FromXmlFile(path string) (*DeployType, error) {
	itf, err := io.FromXmlFile(path, dt)
	if err != nil {
		return nil, err
	}
	var conf *DeployType = itf.(*DeployType)
	return conf, nil
}

func (dt *DeployType) FromXmlCode(xmlCode string) (*DeployType, error) {
	itf, err := io.FromXmlCode(xmlCode, dt)
	if err != nil {
		return nil, err
	}
	var conf *DeployType = itf.(*DeployType)
	return conf, nil
}

func (dt *DeployType) Json() (string, error) {
	return io.ToJson(dt)
}

func (dt *DeployType) FromJsonFile(path string) (*DeployType, error) {
	itf, err := io.FromJsonFile(path, dt)
	if err != nil {
		return nil, err
	}
	var conf *DeployType = itf.(*DeployType)
	return conf, nil
}

func (dt *DeployType) FromJsonCode(jsonCode string) (*DeployType, error) {
	itf, err := io.FromJsonCode(jsonCode, dt)
	if err != nil {
		return nil, err
	}
	var conf *DeployType = itf.(*DeployType)
	return conf, nil
}

func (npt *NetProtocolType) Merge(npt2 *NetProtocolType) *NetProtocolType {
	return &NetProtocolType{
		CaCert:      bestString(npt2.CaCert, npt.CaCert),
		KeyFile:     bestString(npt2.KeyFile, npt.KeyFile),
		NetProtocol: NetProtocolTypeValue(bestString(string(npt2.NetProtocol), string(npt.NetProtocol))),
		Passphrase:  bestString(npt2.Passphrase, npt.Passphrase),
		UserName:    bestString(npt2.UserName, npt.UserName),
		Password:    bestString(npt2.Password, npt.Password),
		Certificate: bestString(npt2.Certificate, npt.Certificate),
	}
}

func (npt *NetProtocolType) String() string {
	return fmt.Sprintf("NetProtocolType{NetProtocol: \"%v\", UserName: \"%s\", Password: \"%s\", KeyFile: \"%s\", CaCert: \"%s\", Passphrase: \"%s\"}",
		npt.NetProtocol, npt.UserName, npt.Password, npt.KeyFile, npt.CaCert, npt.Passphrase)
}

func (npt *NetProtocolType) Yaml() (string, error) {
	return io.ToYaml(npt)
}

func (npt *NetProtocolType) FromYamlFile(path string) (*NetProtocolType, error) {
	itf, err := io.FromYamlFile(path, npt)
	if err != nil {
		return nil, err
	}
	var conf *NetProtocolType = itf.(*NetProtocolType)
	return conf, nil
}

func (npt *NetProtocolType) FromYamlCode(yamlCode string) (*NetProtocolType, error) {
	itf, err := io.FromYamlCode(yamlCode, npt)
	if err != nil {
		return nil, err
	}
	var conf *NetProtocolType = itf.(*NetProtocolType)
	return conf, nil
}

func (npt *NetProtocolType) Xml() (string, error) {
	return io.ToXml(npt)
}

func (npt *NetProtocolType) FromXmlFile(path string) (*NetProtocolType, error) {
	itf, err := io.FromXmlFile(path, npt)
	if err != nil {
		return nil, err
	}
	var conf *NetProtocolType = itf.(*NetProtocolType)
	return conf, nil
}

func (npt *NetProtocolType) FromXmlCode(xmlCode string) (*NetProtocolType, error) {
	itf, err := io.FromXmlCode(xmlCode, npt)
	if err != nil {
		return nil, err
	}
	var conf *NetProtocolType = itf.(*NetProtocolType)
	return conf, nil
}

func (npt *NetProtocolType) Json() (string, error) {
	return io.ToJson(npt)
}

func (npt *NetProtocolType) FromJsonFile(path string) (*NetProtocolType, error) {
	itf, err := io.FromJsonFile(path, npt)
	if err != nil {
		return nil, err
	}
	var conf *NetProtocolType = itf.(*NetProtocolType)
	return conf, nil
}

func (npt *NetProtocolType) FromJsonCode(jsonCode string) (*NetProtocolType, error) {
	itf, err := io.FromJsonCode(jsonCode, npt)
	if err != nil {
		return nil, err
	}
	var conf *NetProtocolType = itf.(*NetProtocolType)
	return conf, nil
}

func (dc *DeployConfig) Merge(dc2 *DeployConfig) *DeployConfig {
	var useHosts []string = utils.StringSliceTrim(utils.StringSliceUnique(utils.StringSliceAppend(dc.UseHosts, dc2.UseHosts)))
	var useVars []string = utils.StringSliceTrim(utils.StringSliceUnique(utils.StringSliceAppend(dc.UseVars, dc2.UseVars)))
	return &DeployConfig{
		ModulesDir:         bestString(dc2.ModulesDir, dc.ModulesDir),
		ConfigDir:          bestString(dc2.ConfigDir, dc.ConfigDir),
		ChartsDir:          bestString(dc2.ChartsDir, dc.ChartsDir),
		SystemDir:          bestString(dc2.SystemDir, dc.SystemDir),
		WorkDir:            bestString(dc2.WorkDir, dc.WorkDir),
		LogVerbosity:       bestString(dc2.LogVerbosity, dc.LogVerbosity),
		ConfigLang:         DescriptorTypeValue(bestString(string(dc2.ConfigLang), string(dc.ConfigLang))),
		DeployName:         bestString(dc2.DeployName, dc.DeployName),
		EnvSelector:        bestString(dc2.EnvSelector, dc.EnvSelector),
		ParallelExecutions: dc2.ParallelExecutions || dc.ParallelExecutions,
		MaxThreads:         maxInt64(dc2.MaxThreads, dc.MaxThreads),
		SingleSession:		dc2.SingleSession || dc.SingleSession,
		ReadTimeout:        maxInt64(dc2.ReadTimeout, dc.ReadTimeout),
		UseHosts:           useHosts,
		UseVars:            useVars,
	}
}

func (dc *DeployConfig) String() string {
	return fmt.Sprintf("DeployConfig{DeployName: \"%s\", UseHosts: %v, UseVars: %v, WorkDir: \"%s\", ConfigDir: \"%s\", ChartsDir: \"%s\", SystemDir: \"%s\", ModulesDir: \"%s\", ConfigLang: \"%v\", LogVerbosity: \"%v\", EnvSelector: \"%s\", SingleSession: %v, ParallelExecutions: %v, MaxThreads: %vm ReadTimeout: %v}",
		dc.DeployName, dc.UseHosts, dc.UseVars, dc.WorkDir, dc.ConfigDir, dc.ChartsDir, dc.SystemDir, dc.ModulesDir, dc.ConfigLang, dc.LogVerbosity, dc.EnvSelector, dc.SingleSession, dc.ParallelExecutions, dc.MaxThreads, dc.ReadTimeout)
}

func (dc *DeployConfig) Yaml() (string, error) {
	return io.ToYaml(dc)
}

func (dc *DeployConfig) FromYamlFile(path string) (*DeployConfig, error) {
	itf, err := io.FromYamlFile(path, dc)
	if err != nil {
		return nil, err
	}
	var conf *DeployConfig = itf.(*DeployConfig)
	return conf, nil
}

func (dc *DeployConfig) FromYamlCode(yamlCode string) (*DeployConfig, error) {
	itf, err := io.FromYamlCode(yamlCode, dc)
	if err != nil {
		return nil, err
	}
	var conf *DeployConfig = itf.(*DeployConfig)
	return conf, nil
}

func (dc *DeployConfig) Xml() (string, error) {
	return io.ToXml(dc)
}

func (dc *DeployConfig) FromXmlFile(path string) (*DeployConfig, error) {
	itf, err := io.FromXmlFile(path, dc)
	if err != nil {
		return nil, err
	}
	var conf *DeployConfig = itf.(*DeployConfig)
	return conf, nil
}

func (dc *DeployConfig) FromXmlCode(xmlCode string) (*DeployConfig, error) {
	itf, err := io.FromXmlCode(xmlCode, dc)
	if err != nil {
		return nil, err
	}
	var conf *DeployConfig = itf.(*DeployConfig)
	return conf, nil
}

func (dc *DeployConfig) Json() (string, error) {
	return io.ToJson(dc)
}

func (dc *DeployConfig) FromJsonFile(path string) (*DeployConfig, error) {
	itf, err := io.FromJsonFile(path, dc)
	if err != nil {
		return nil, err
	}
	var conf *DeployConfig = itf.(*DeployConfig)
	return conf, nil
}

func (dc *DeployConfig) FromJsonCode(jsonCode string) (*DeployConfig, error) {
	itf, err := io.FromJsonCode(jsonCode, dc)
	if err != nil {
		return nil, err
	}
	var conf *DeployConfig = itf.(*DeployConfig)
	return conf, nil
}

func bestString(str1 string, str2 string) string {
	if str1 != "" {
		return str1
	}
	return str2
}

func maxInt64(a int64, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
