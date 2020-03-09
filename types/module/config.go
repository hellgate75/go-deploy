package module

import (
	"github.com/hellgate75/go-tcp-common/log"
)

var Logger log.Logger = nil

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
	// Unknown Feed(s) and activation Source
	UNKNOWN                       DeploymentTypeValue  = "UNKNOWN"
	// File Feed(s) and activation Source
	FILE_SOURCE                   DeploymentTypeValue  = "FILE_SOURCE"
	// HTTP Call Feed(s) and activation Source
	HTTP_SOURCE                   DeploymentTypeValue  = "HTTP_SOURCE"
	// Rest service Feed(s) and activation Source
	REST_SOURCE                   DeploymentTypeValue  = "REST_SOURCE"
	// Pipe channel Feed(s) and activation Source
	PIPE_SOURCE                   DeploymentTypeValue  = "PIPE_SOURCE"
	// Streaing engine Feed(s) and activation Source
	STREAM_SOURCE                 DeploymentTypeValue  = "STREAM_SOURCE"
	// YAML File Descriptor
	YAML_DESCRIPTOR               DescriptorTypeValue  = "YAML"
	// JSON File Descriptor
	JSON_DESCRIPTOR               DescriptorTypeValue  = "JSON"
	// XML File Descriptor
	XML_DESCRIPTOR                DescriptorTypeValue  = "XML"
	// ONE SHOT Deployment Strategy
	ONE_SHOT_DEPLOYMENT           StrategyTypeValue    = "ONE_SHOT_DEPLOYMENT"
	// Continuous Deployment Strategy
	CONTINUOUS_DEPLOYMENT         StrategyTypeValue    = "CONTINUOUS_DEPLOYMENT"
	// Awake On-Demand Deployment Strategy
	ON_DEMAND_DEPLOYMENT          StrategyTypeValue    = "ON_DEMAND_DEPLOYMENT"
	// Sheduled and periodic Deployment Strategy
	PERIODIC_DEPLOYMENT           StrategyTypeValue    = "PERIODIC_DEPLOYMENT"
	// REST Service Get Requiest Identifier
	REST_GET_REQUEST              RestMethodTypeValue  = "GET"
	// REST Service Post Requiest Identifier
	REST_POST_REQUEST             RestMethodTypeValue  = "POST"
	// SSH Built-in Client protocol
	NET_PROTOCOL_SSH              NetProtocolTypeValue = "SSH"
	// Go! TLS/TCP Client protocol
	NET_PROTOCOL_GO_DEPLOY_CLIENT NetProtocolTypeValue = "GO_DEPLOY"
)

// Deploy Behaviour Configuration Struture
type DeployType struct {
	DeploymentType DeploymentTypeValue `yaml:"deploymentType,omitempty" json:"deploymentType,omitempty" xml:"deployment-type,chardata,omitempty"`
	DescriptorType DescriptorTypeValue `yaml:"descriptorType,omitempty" json:"descriptorType,omitempty" xmk:"descriptor-type,chardata,omitempty"`
	StrategyType   StrategyTypeValue   `yaml:"strategyType,omitempty" json:"strategyType,omitempty" json:"strategy-type,chardata,omitempty"`
	Scheduled      string              `yaml:"scheduled,omitempty" json:"scheduled,omitempty" json:"scheduled,chardata,omitempty"`
	Method         RestMethodTypeValue `yaml:"restMethod,omitempty" json:"restMethod,omitempty" json:"restMethod,chardata,omitempty"`
	PostBody       string              `yaml:"postBody,omitempty" json:"postBody,omitempty" json:"post-body,chardata,omitempty"`
}

// Networking and Client Configuration Struture
type NetProtocolType struct {
	NetProtocol NetProtocolTypeValue `yaml:"protocol,omitempty" json:"protocol,omitempty" xml:"protocol,chardata,omitempty"`
	UserName    string               `yaml:"userName,omitempty" json:"userName,omitempty" xml:"username,chardata,omitempty"`
	Password    string               `yaml:"password,omitempty" json:"password,omitempty" xml:"password,chardata,omitempty"`
	KeyFile     string               `yaml:"keyFile,omitempty" json:"keyFile,omitempty" xml:"key-file,chardata,omitempty"`
	CaCert      string               `yaml:"caCert,omitempty" json:"caCert,omitempty" xml:"ca-cert,chardata,omitempty"`
	Passphrase  string               `yaml:"passphrase,omitempty" json:"passphrase,omitempty" xml:"passphrase,chardata,omitempty"`
	Certificate string               `yaml:"certificate,omitempty" json:"certificate,omitempty" xml:"certificate,chardata,omitempty"`
}

// Main Configuration Struture
type DeployConfig struct {
	DeployName         string              `yaml:"deployName,omitempty" json:"deployName,omitempty" xml:"deploy-name,chardata,omitempty"`
	LogVerbosity       string              `yaml:"verbosity,omitempty" json:"verbosity,omitempty" xml:"verbosity,chardata,omitempty"`
	UseHosts           []string            `yaml:"useHosts,omitempty" json:"useHosts,omitempty" xml:"use-hosts,chardata,omitempty"`
	UseVars            []string            `yaml:"useVars,omitempty" json:"useVars,omitempty" xml:"use-vars,chardata,omitempty"`
	WorkDir            string              `yaml:"workDir,omitempty" json:"workDir,omitempty" xml:"work-dir,chardata,omitempty"`
	ModulesDir         string              `yaml:"modulesDir,omitempty" json:"modulesDir,omitempty" xml:"modules-dir,chardata,omitempty"`
	ConfigDir          string              `yaml:"configDir,omitempty" json:"configDir,omitempty" xml:"config-dir,chardata,omitempty"`
	ChartsDir          string              `yaml:"chartsDir,omitempty" json:"chartsDir,omitempty" xml:"charts-dir,chardata,omitempty"`
	SystemDir          string              `yaml:"systemDir,omitempty" json:"systemDir,omitempty" xml:"system-dir,chardata,omitempty"`
	ConfigLang         DescriptorTypeValue `yaml:"configLang,omitempty" json:"configLang,omitempty" xml:"config-lang,chardata,omitempty"`
	EnvSelector        string              `yaml:"env,omitempty" json:"env,omitempty" xml:"env,chardata,omitempty"`
	ParallelExecutions bool                `yaml:"parallel,omitempty" json:"parallel,omitempty" xml:"parallel,chardata,omitempty"`
	MaxThreads         int64               `yaml:"maxThreads,omitempty" json:"maxThreads,omitempty" xml:"max-threads,chardata,omitempty"`
	SingleSession      bool                `yaml:"singleSession,omitempty" json:"singleSession,omitempty" xml:"single-session,chardata,omitempty"`
	ReadTimeout      int64                `yaml:"readTimeout,omitempty" json:"readTimeout,omitempty" xml:"read-timeout,chardata,omitempty"`
}

// Plugins Configuration Struture
type PluginsConfig struct {
	EnableDeployClientsPlugin           bool   `yaml:"enableDeployClientsPlugin,omitempty" json:"enableDeployClientsPlugin,omitempty" xml:"enable-deploy-clients-plugin,chardata,omitempty"`
	DeployClientsPluginExtension        string `yaml:"deployClientsPluginExtension,omitempty" json:"deployClientsPluginExtension,omitempty" xml:"deploy-clients-plugin-extension,chardata,omitempty"`
	DeployClientsPluginFolder           string `yaml:"deployClientsPluginFolder,omitempty" json:"deployClientsPluginFolder,omitempty" xml:"deploy-clients-plugin-folder,chardata,omitempty"`
	EnableDeployCommandsPlugin          bool   `yaml:"enableDeployCommandsPlugin,omitempty" json:"enableDeployCommandsPlugin,omitempty" xml:"enable-deploy-commands-plugin,chardata,omitempty"`
	DeployCommandsPluginExtension       string `yaml:"deployCommandsPluginExtension,omitempty" json:"deployCommandsPluginExtension,omitempty" xml:"deploy-commands-plugin-extension,chardata,omitempty"`
	DeployCommandsPluginFolder          string `yaml:"deployCommandsPluginFolder,omitempty" json:"deployCommandsPluginFolder,omitempty" xml:"deploy-commands-plugin-folder,chardata,omitempty"`
	EnableDeployClientCommandsPlugin    bool   `yaml:"enableDeployClientCommandsPlugin,omitempty" json:"enableDeployClientCommandsPlugin,omitempty" xml:"enable-deploy-client-commands-plugin,chardata,omitempty"`
	DeployClientCommandsPluginExtension string `yaml:"deployClientCommandsPluginExtension,omitempty" json:"deployClientCommandsPluginExtension,omitempty" xml:"deploy-client-commands-plugin-extension,chardata,omitempty"`
	DeployClientCommandsPluginFolder    string `yaml:"deployClientCommandsPluginFolder,omitempty" json:"deployClientCommandsPluginFolder,omitempty" xml:"deploy-client-commands-plugin-folder,chardata,omitempty"`
}

// Printable interface, allows system to print as sting any implementing components (almost all in this project)
type Printable interface {
	// Traslates the object in printable version <BR/>
	// Return: <BR/>
	// (string) Representation of the structure<BR/>
	String() string
}

// Step Structure
type Step struct {
	Name     string
	StepType string
	StepData interface{}
	Children []*Step
	Feeds    []*FeedExec
}

// Executable Feed Structure
type FeedExec struct {
	Name      string
	HostGroup string
	Steps     []*Step
}

// Session Interface
type Session interface {
	// Retrives Session Unique Id
	GetSessionId() string
	// Retrives Session Deploy Type
	GetDeployType() *DeployType
	// Retrives Session Network Protocol Type
	GetNetProtocolType() *NetProtocolType
	// Retrives Session Deploy Config
	GetDeployConfig() *DeployConfig
	// Retrives a Session Variable by key
	GetVar(name string) (string, error)
	// Sets a Session Variable
	SetVar(name string, value string) bool
	// Retrives all Session Variable keys
	GetKeys() []string
	// Retrives a Session Object by key
	GetSystemObject(name string) (interface{}, error)
	// Sets a Session Variable
	// Build-in are required variables present in Vars files
	SetSystemObject(name string, value interface{}) bool
	// Retrives all Session Object keys
	// Build-in session objects :
	// connection-handler -> Current Session ConnectionHandler
	// rutime-config -> Session module.DeployConfig
	// runtime-type -> Session module.DeployType
	// runtime-net -> Session module.NetworkType
	// host-groups -> Current Running defaults.HostGroup
	// envs -> Session Environemnts Row defaults.NameValuePair list
	// vars -> Session Variables Row defaults.NameValuePair list
	// system-logger -> Centrilized Go! Deploy logger instance
	GetSystemKeys() []string
}

// Client Connection Features Availabolity Configuration
type ConnectionConfig struct {
	UseUserPassword         bool
	UseUserKey              bool
	UseUserKeyPassphrase    bool
	UseSSHConfig            bool
	UseTLSCertificates      bool
}

