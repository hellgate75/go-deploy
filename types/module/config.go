package module

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hellgate75/go-tcp-common/log"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
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
	KeyFile     string               `yaml:"keyFile,omitempty" json:"keyFile,omitempty" xml:"key-file,chardata,omitempty"`
	CaCert      string               `yaml:"caCert,omitempty" json:"caCert,omitempty" xml:"ca-cert,chardata,omitempty"`
	Passphrase  string               `yaml:"passphrase,omitempty" json:"passphrase,omitempty" xml:"passphrase,chardata,omitempty"`
	Certificate string               `yaml:"certificate,omitempty" json:"certificate,omitempty" xml:"certificate,chardata,omitempty"`
}

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
	Name     string
	StepType string
	StepData interface{}
	Children []*Step
	Feeds    []*FeedExec
}

type FeedExec struct {
	Name      string
	HostGroup string
	Steps     []*Step
}

var sessionVars map[string]map[string]string = make(map[string]map[string]string)
var sessionsMap map[string]*session = make(map[string]*session)
var systemObjectsMap map[string]map[string]interface{} = make(map[string]map[string]interface{})

type Session interface {
	GetSessionId() string
	GetDeployType() *DeployType
	GetNetProtocolType() *NetProtocolType
	GetDeployConfig() *DeployConfig
	GetVar(name string) (string, error)
	SetVar(name string, value string) bool
	GetKeys() []string
	GetSystemObject(name string) (interface{}, error)
	SetSystemObject(name string, value interface{}) bool
	GetSystemKeys() []string
}

type session struct {
	sync.RWMutex
	sessionId       string
	deployType      *DeployType
	netProtocolType *NetProtocolType
	deployConfig    *DeployConfig
}

func (session *session) GetSessionId() string {
	return session.sessionId
}
func (session *session) GetDeployType() *DeployType {
	return session.deployType
}
func (session *session) GetNetProtocolType() *NetProtocolType {
	return session.netProtocolType
}
func (session *session) GetDeployConfig() *DeployConfig {
	return session.deployConfig
}
func (session *session) GetVar(name string) (string, error) {
	defer func() {
		if r := recover(); r != nil {
			Logger.Errorf("Session.GetVar : %v", r)
		}
		session.RUnlock()
	}()
	session.RLock()
	if value, ok := sessionVars[session.sessionId][name]; ok {
		return value, nil
	} else {
		return "", errors.New(fmt.Sprintf("Variable %s not found in session!!", name))
	}
}
func (session *session) SetVar(name string, value string) bool {
	var out bool = true
	defer func() {
		if r := recover(); r != nil {
			Logger.Errorf("Session.GetVar : %v", r)
		}
		out = false
		session.Unlock()
	}()
	session.Lock()
	sessionVars[session.sessionId][name] = value
	return out
}
func (session *session) GetKeys() []string {
	defer func() {
		if r := recover(); r != nil {
			Logger.Errorf("Session.GetVar : %v", r)
		}
		session.RUnlock()
	}()
	var keys []string = make([]string, 0)
	session.RLock()
	for k, _ := range sessionVars[session.sessionId] {
		keys = append(keys, k)
	}
	return keys
}
func (session *session) GetSystemObject(name string) (interface{}, error) {
	defer func() {
		if r := recover(); r != nil {
			Logger.Errorf("Session.GetVar : %v", r)
		}
		session.RUnlock()
	}()
	session.RLock()
	if value, ok := systemObjectsMap[session.sessionId][name]; ok {
		return value, nil
	} else {
		return "", errors.New(fmt.Sprintf("Variable %s not found in session!!", name))
	}
}
func (session *session) SetSystemObject(name string, value interface{}) bool {
	var out bool = true
	defer func() {
		if r := recover(); r != nil {
			Logger.Errorf("Session.GetVar : %v", r)
		}
		out = false
		session.Unlock()
	}()
	session.Lock()
	systemObjectsMap[session.sessionId][name] = value
	return out
}
func (session *session) GetSystemKeys() []string {
	defer func() {
		if r := recover(); r != nil {
			Logger.Errorf("Session.GetVar : %v", r)
		}
		session.RUnlock()
	}()
	var keys []string = make([]string, 0)
	session.RLock()
	for k, _ := range systemObjectsMap[session.sessionId] {
		keys = append(keys, k)
	}
	return keys
}

func NewSessionId() string {
	uuid, err := uuid.NewUUID()
	if err != nil {
		Logger.Errorf("NewSessionId unable to create google.UUID -> Reason: %s", err.Error())
		block1 := strconv.FormatInt(rand.Int63(), 16)
		block2 := strconv.FormatInt(rand.Int63(), 16)
		block3 := strconv.FormatInt(rand.Int63(), 16)
		block4 := strconv.FormatInt(rand.Int63(), 16)
		block5 := strconv.FormatInt(rand.Int63(), 16)
		return fmt.Sprintf("%s-%s-%s-%s-%s", block1, block2, block3, block4, block5)
	}
	return uuid.String()
}

func DestroySession(sessionId string) {
	if _, ok := sessionVars[sessionId]; ok {
		sessionVars[sessionId] = nil
	}
	if _, ok := systemObjectsMap[sessionId]; ok {
		systemObjectsMap[sessionId] = nil
	}
	if _, ok := sessionsMap[sessionId]; ok {
		sessionsMap[sessionId] = nil
	}
	runtime.GC()
}

func NewSession(sessionId string) Session {
	if _, ok := sessionVars[sessionId]; !ok {
		sessionVars[sessionId] = make(map[string]string)
	}
	if _, ok := systemObjectsMap[sessionId]; !ok {
		systemObjectsMap[sessionId] = make(map[string]interface{})
	}
	if sessionX, ok := sessionsMap[sessionId]; ok {
		return sessionX
	} else {
		sessionX := &session{
			sessionId:       sessionId,
			deployConfig:    RuntimeDeployConfig,
			deployType:      RuntimeDeployType,
			netProtocolType: RuntimeNetworkType,
		}
		sessionsMap[sessionId] = sessionX
		return sessionX
	}
}

type ConnectionConfig struct {
	UseUserPassword         bool
	UseUserKey              bool
	UseUserKeyPassphrase    bool
	UseSSHConfig            bool
	UseTLSCertificates      bool
}

