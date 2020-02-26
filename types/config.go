package types

import (
	"fmt"
	"github.com/hellgate75/go-deploy/io"
)

type deploySourceTypeValue byte
type deployDescriptorTypeValue byte
type deployBehaviouralTypeValue byte
type DeploymentTypeValue string
type DescriptorTypeValue string
type StrategyTypeValue string
type RestMethodTypeValue string
type NetProtocolTypeValue string

const (
	unknownSource deploySourceTypeValue = 0
	fileSource    deploySourceTypeValue = iota + 1
	httpSource
	restSource
	pipeSource
	streamSource
	unknownDeployment deploySourceTypeValue      = 0
	oneShotDeployment deployBehaviouralTypeValue = iota + 1
	continuousDeployment
	onDemandDeployment
	periodicDeployment
	yamlDescriptorType deployDescriptorTypeValue = iota + 1
	jsonDescriptorType
	UNKNOWN                       DeploymentTypeValue  = "UNKNOWN"
	FILE_SOURCE                   DeploymentTypeValue  = "FILE_SOURCE"
	HTTP_SOURCE                   DeploymentTypeValue  = "HTTP_SOURCE"
	REST_SOURCE                   DeploymentTypeValue  = "REST_SOURCE"
	PIPE_SOURCE                   DeploymentTypeValue  = "PIPE_SOURCE"
	STREAM_SOURCE                 DeploymentTypeValue  = "STREAM_SOURCE"
	YAML_DESCRIPTOR               DescriptorTypeValue  = "YAML"
	JSON_DESCRIPTOR               DescriptorTypeValue  = "JSON"
	XML_DESCRIPTOR                DescriptorTypeValue  = "XML"
	ONE_SHOT_DEPLOYMENT           StrategyTypeValue    = "ONE_SHOT_DEPLOYMENT"
	CONTINUOUS_DEPLOYMENT         StrategyTypeValue    = "CONTINUOUS_DEPLOYMENT"
	ON_DEMAND_DEPLOYMENT          StrategyTypeValue    = "ON_DEMAND_DEPLOYMENT"
	PERIODIC_DEPLOYMENT           StrategyTypeValue    = "PERIODIC_DEPLOYMENT"
	REST_GET_REQUEST              RestMethodTypeValue  = "GET"
	REST_POST_REQUEST             RestMethodTypeValue  = "POST"
	NET_PROTOCOL_SSH              NetProtocolTypeValue = "SSH"
	NET_PROTOCOL_GO_DEPLOY_CLIENT NetProtocolTypeValue = "GO_DEPLOY"
	DEFAULT_CONFIG_FOLDER         string               = ".deploy"
)

type DeployType struct {
	DeploymentType DeploymentTypeValue `yaml:"deploymentType,omitempty" json:"deploymentType,omitempty" xml:"deployment-type,chardata,omitempty"`
	DescriptorType DescriptorTypeValue `yaml:"descriptorType,omitempty" json:"descriptorType,omitempty" xmk:"descriptor-type,chardata,omitempty"`
	StrategyType   StrategyTypeValue   `yaml:"strategyType,omitempty" json:"strategyType,omitempty" json:"strategy-type,chardata,omitempty"`
	Scheduled      string              `yaml:"scheduled,omitempty" json:"scheduled,omitempty" json:"scheduled,chardata,omitempty"`
	Method         RestMethodTypeValue `yaml:"restMethod,omitempty" json:"restMethod,omitempty" json:"restMethod,chardata,omitempty"`
	PostBody       string              `yaml:"postBody,omitempty" json:"postBody,omitempty" json:"post-body,chardata,omitempty"`
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

type NetProtocolType struct {
	NetProtocol NetProtocolTypeValue `yaml:"protocol,omitempty" json:"protocol,omitempty" xml:"protocol,chardata,omitempty"`
	UserName    string               `yaml:"userName,omitempty" json:"userName,omitempty" xml:"username,chardata,omitempty"`
	Password    string               `yaml:"password,omitempty" json:"password,omitempty" xml:"password,chardata,omitempty"`
	KeyFile     string               `yaml:"keyFile,omitempty" json:"keyFile,omitempty" xml:"keyfile,chardata,omitempty"`
	Passphrase  string               `yaml:"passphrase,omitempty" json:"passphrase,omitempty" xml:"passphrase,chardata,omitempty"`
}

func (npt *NetProtocolType) Merge(npt2 *NetProtocolType) *NetProtocolType {
	return &NetProtocolType{
		KeyFile:     bestString(npt2.KeyFile, npt.KeyFile),
		NetProtocol: NetProtocolTypeValue(bestString(string(npt2.KeyFile), string(npt.KeyFile))),
		Passphrase:  bestString(npt2.Passphrase, npt.Passphrase),
		UserName:    bestString(npt2.UserName, npt.UserName),
		Password:    bestString(npt2.Password, npt.Password),
	}
}

func (npt *NetProtocolType) String() string {
	return fmt.Sprintf("NetProtocolType{NetProtocol: \"%v\", UserName: \"%s\", Password: \"%s\", KeyFile: \"%s\", Passphrase: \"%s\"}",
		npt.NetProtocol, npt.UserName, npt.Password, npt.KeyFile, npt.Passphrase)
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

type DeployConfig struct {
	DeployName  string              `yaml:"deployName,omitempty" json:"deployName,omitempty" xml:"deploy-name,chardata,omitempty"`
	UseHosts    []string            `yaml:"useHosts,omitempty" json:"useHosts,omitempty" xml:"use-hosts,chardata,omitempty"`
	UseVars     []string            `yaml:"useVars,omitempty" json:"useVars,omitempty" xml:"use-vars,chardata,omitempty"`
	WorkDir     string              `yaml:"workDir,omitempty" json:"workDir,omitempty" xml:"work-dir,chardata,omitempty"`
	ConfigDir   string              `yaml:"configDir,omitempty" json:"configDir,omitempty" xml:"config-dir,chardata,omitempty"`
	ConfigLang  DescriptorTypeValue `yaml:"configLang,omitempty" json:"configLang,omitempty" xml:"config-lang,chardata,omitempty"`
	EnvSelector string              `yaml:"env,omitempty" json:"env,omitempty" xml:"env,chardata,omitempty"`
}

func (dc *DeployConfig) Merge(dc2 *DeployConfig) *DeployConfig {
	var useHosts []string = make([]string, 0)
	var useVars []string = make([]string, 0)
	for _, val := range dc.UseHosts {
		useHosts = append(useHosts, val)
	}
	for _, val := range dc2.UseHosts {
		useHosts = append(useHosts, val)
	}
	for _, val := range dc.UseVars {
		useVars = append(useVars, val)
	}
	for _, val := range dc2.UseVars {
		useVars = append(useVars, val)
	}
	return &DeployConfig{
		ConfigDir:   bestString(dc2.ConfigDir, dc.ConfigDir),
		WorkDir:     bestString(dc2.WorkDir, dc.WorkDir),
		ConfigLang:  DescriptorTypeValue(bestString(string(dc2.ConfigLang), string(dc.ConfigLang))),
		DeployName:  bestString(dc2.DeployName, dc.DeployName),
		EnvSelector: bestString(dc2.EnvSelector, dc.EnvSelector),
		UseHosts:    useHosts,
		UseVars:     useVars,
	}
}

func (dc *DeployConfig) String() string {
	return fmt.Sprintf("DeployConfig{DeployName: \"%s\", UseHosts: %v, UseVars: %v, WorkDir: \"%s\", ConfigDir: \"%s\", ConfigLang: \"%v\", EnvSelector: \"%s\"}",
		dc.DeployName, dc.UseHosts, dc.UseVars, dc.WorkDir, dc.ConfigDir, dc.ConfigLang, dc.EnvSelector)
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
