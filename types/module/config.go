package module

import ()

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
	xmlDescriptorType
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
)

type DeployType struct {
	DeploymentType DeploymentTypeValue `yaml:"deploymentType,omitempty" json:"deploymentType,omitempty" xml:"deployment-type,chardata,omitempty"`
	DescriptorType DescriptorTypeValue `yaml:"descriptorType,omitempty" json:"descriptorType,omitempty" xmk:"descriptor-type,chardata,omitempty"`
	StrategyType   StrategyTypeValue   `yaml:"strategyType,omitempty" json:"strategyType,omitempty" json:"strategy-type,chardata,omitempty"`
	Scheduled      string              `yaml:"scheduled,omitempty" json:"scheduled,omitempty" json:"scheduled,chardata,omitempty"`
	Method         RestMethodTypeValue `yaml:"restMethod,omitempty" json:"restMethod,omitempty" json:"restMethod,chardata,omitempty"`
	PostBody       string              `yaml:"postBody,omitempty" json:"postBody,omitempty" json:"post-body,chardata,omitempty"`
}

type NetProtocolType struct {
	NetProtocol NetProtocolTypeValue `yaml:"protocol,omitempty" json:"protocol,omitempty" xml:"protocol,chardata,omitempty"`
	UserName    string               `yaml:"userName,omitempty" json:"userName,omitempty" xml:"username,chardata,omitempty"`
	Password    string               `yaml:"password,omitempty" json:"password,omitempty" xml:"password,chardata,omitempty"`
	KeyFile     string               `yaml:"keyFile,omitempty" json:"keyFile,omitempty" xml:"keyfile,chardata,omitempty"`
	Passphrase  string               `yaml:"passphrase,omitempty" json:"passphrase,omitempty" xml:"passphrase,chardata,omitempty"`
}

type DeployConfig struct {
	DeployName   string              `yaml:"deployName,omitempty" json:"deployName,omitempty" xml:"deploy-name,chardata,omitempty"`
	LogVerbosity string              `yaml:"verbosity,omitempty" json:"verbosity,omitempty" xml:"verbosity,chardata,omitempty"`
	UseHosts     []string            `yaml:"useHosts,omitempty" json:"useHosts,omitempty" xml:"use-hosts,chardata,omitempty"`
	UseVars      []string            `yaml:"useVars,omitempty" json:"useVars,omitempty" xml:"use-vars,chardata,omitempty"`
	WorkDir      string              `yaml:"workDir,omitempty" json:"workDir,omitempty" xml:"work-dir,chardata,omitempty"`
	ModulesDir   string              `yaml:"modulesDir,omitempty" json:"modulesDir,omitempty" xml:"modules-dir,chardata,omitempty"`
	ConfigDir    string              `yaml:"configDir,omitempty" json:"configDir,omitempty" xml:"config-dir,chardata,omitempty"`
	ChartsDir    string              `yaml:"chartsDir,omitempty" json:"chartsDir,omitempty" xml:"charts-dir,chardata,omitempty"`
	SystemDir    string              `yaml:"systemDir,omitempty" json:"systemDir,omitempty" xml:"system-dir,chardata,omitempty"`
	ConfigLang   DescriptorTypeValue `yaml:"configLang,omitempty" json:"configLang,omitempty" xml:"config-lang,chardata,omitempty"`
	EnvSelector  string              `yaml:"env,omitempty" json:"env,omitempty" xml:"env,chardata,omitempty"`
}

/*
* Coverter interface, responsible to comvert raw interface from the parsing to a specific structure
 */
type Printable interface {
	/*
	* Traslates the object in printable version <BR/>
	* Return: <BR/>
	* (string) Representation of the structure<BR/>
	 */
	String() string
}

type Step struct {
	StepType string
	StepData interface{}
	Children []*Step
	Feeds    []*FeedExec
}

type FeedExec struct {
	Steps []*Step
}
