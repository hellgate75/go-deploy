package defaults

import (
	"fmt"
	"github.com/hellgate75/go-deploy/io"
)

type NameValue struct {
	Name  string `yaml:"name" json:"name" xml:"name,chardata"`
	Value string `yaml:"value,omitempty" json:"value,omitempty" xml:"value,chardata,omitempty"`
}

func (nv *NameValue) String() string {
	return fmt.Sprintf("NameValue{Name: %s, Value: %s}",
		nv.Name, nv.Value)
}

type Vars struct {
	Vars []NameValue `yaml:"vars,omitempty" json:"vars,omitempty" xml:"vars,chardata,omitempty"`
}

func (vars *Vars) String() string {
	var hostsVal string = "["
	for _, varX := range vars.Vars {
		prefix := ", "
		if len(hostsVal) == 0 {
			prefix = ""
		}
		hostsVal += prefix + varX.String()
	}
	hostsVal += "]"
	return fmt.Sprintf("Vars{Vars: %s}",
		hostsVal)
}

func (vars *Vars) Yaml() (string, error) {
	return io.ToYaml(vars)
}

func (vars *Vars) FromYamlFile(path string) (*Vars, error) {
	itf, err := io.FromYamlFile(path, vars)
	if err != nil {
		return nil, err
	}
	var varsObj *Vars = itf.(*Vars)
	return varsObj, nil
}

func (vars *Vars) FromYamlCode(yamlCode string) (*Vars, error) {
	itf, err := io.FromYamlCode(yamlCode, vars)
	if err != nil {
		return nil, err
	}
	var varsObj *Vars = itf.(*Vars)
	return varsObj, nil
}

func (vars *Vars) Xml() (string, error) {
	return io.ToXml(vars)
}

func (vars *Vars) FromXmlFile(path string) (*Vars, error) {
	itf, err := io.FromXmlFile(path, vars)
	if err != nil {
		return nil, err
	}
	var varsObj *Vars = itf.(*Vars)
	return varsObj, nil
}

func (vars *Vars) FromXmlCode(xmlCode string) (*Vars, error) {
	itf, err := io.FromXmlCode(xmlCode, vars)
	if err != nil {
		return nil, err
	}
	var varsObj *Vars = itf.(*Vars)
	return varsObj, nil
}

func (vars *Vars) Json() (string, error) {
	return io.ToJson(vars)
}

func (vars *Vars) FromJsonFile(path string) (*Vars, error) {
	itf, err := io.FromJsonFile(path, vars)
	if err != nil {
		return nil, err
	}
	var varsObj *Vars = itf.(*Vars)
	return varsObj, nil
}

func (vars *Vars) FromJsonCode(jsonCode string) (*Vars, error) {
	itf, err := io.FromJsonCode(jsonCode, vars)
	if err != nil {
		return nil, err
	}
	var varsObj *Vars = itf.(*Vars)
	return varsObj, nil
}

type Environments struct {
	Envs []NameValue `yaml:"environments,omitempty" json:"environments,omitempty" xml:"environments,chardata,omitempty"`
}

func (envs *Environments) String() string {
	var hostsVal string = "["
	for _, varX := range envs.Envs {
		prefix := ", "
		if len(hostsVal) == 0 {
			prefix = ""
		}
		hostsVal += prefix + varX.String()
	}
	hostsVal += "]"
	return fmt.Sprintf("Environments{Envs: %s}",
		hostsVal)
}

func (envs *Environments) Yaml() (string, error) {
	return io.ToYaml(envs)
}

func (envs *Environments) FromYamlFile(path string) (*Environments, error) {
	itf, err := io.FromYamlFile(path, envs)
	if err != nil {
		return nil, err
	}
	var envsObj *Environments = itf.(*Environments)
	return envsObj, nil
}

func (envs *Environments) FromYamlCode(yamlCode string) (*Environments, error) {
	itf, err := io.FromYamlCode(yamlCode, envs)
	if err != nil {
		return nil, err
	}
	var envsObj *Environments = itf.(*Environments)
	return envsObj, nil
}

func (envs *Environments) Xml() (string, error) {
	return io.ToXml(envs)
}

func (envs *Environments) FromXmlFile(path string) (*Environments, error) {
	itf, err := io.FromXmlFile(path, envs)
	if err != nil {
		return nil, err
	}
	var envsObj *Environments = itf.(*Environments)
	return envsObj, nil
}

func (envs *Environments) FromXmlCode(xmlCode string) (*Environments, error) {
	itf, err := io.FromXmlCode(xmlCode, envs)
	if err != nil {
		return nil, err
	}
	var envsObj *Environments = itf.(*Environments)
	return envsObj, nil
}

func (envs *Environments) Json() (string, error) {
	return io.ToJson(envs)
}

func (envs *Environments) FromJsonFile(path string) (*Environments, error) {
	itf, err := io.FromJsonFile(path, envs)
	if err != nil {
		return nil, err
	}
	var envsObj *Environments = itf.(*Environments)
	return envsObj, nil
}

func (envs *Environments) FromJsonCode(jsonCode string) (*Environments, error) {
	itf, err := io.FromJsonCode(jsonCode, envs)
	if err != nil {
		return nil, err
	}
	var envsObj *Environments = itf.(*Environments)
	return envsObj, nil
}

type Hosts struct {
	Hosts []HostValue `yaml:"hosts,omitempty" json:"hosts,omitempty" xml:"hosts,chardata,omitempty"`
}

func (hosts *Hosts) String() string {
	var hostsVal string = "["
	for _, host := range hosts.Hosts {
		prefix := ", "
		if len(hostsVal) == 0 {
			prefix = ""
		}
		hostsVal += prefix + host.String()
	}
	hostsVal += "]"
	return fmt.Sprintf("Hosts{Hosts: %s}",
		hostsVal)
}

type HostGroups struct {
	Name  string  `yaml:"name" json:"name" xml:"name,chardata"`
	Hosts []Hosts `yaml:"hosts" json:"hosts" xml:"hosts,chardata"`
}

func (hg *HostGroups) String() string {
	return fmt.Sprintf("HostValue{Name: \"%s\", Hosts: \"%v\"}",
		hg.Name, hg.Hosts)
}

type HostValue struct {
	Name      string   `yaml:"name" json:"name" xml:"name,chardata"`
	IpAddress string   `yaml:"ipAddress,omitempty" json:"ipAddress,omitempty" xml:"ip-address,chardata,omitempty"`
	HostName  string   `yaml:"hostName,omitempty" json:"hostName,omitempty" xml:"host-name,chardata,omitempty"`
	Port      string   `yaml:"port,omitempty" json:"port,omitempty" xml:"port,chardata,omitempty"`
	Roles     []string `yaml:"roles,omitempty" json:"roles,omitempty" xml:"roles,chardata,omitempty"`
}

func (hv *HostValue) String() string {
	return fmt.Sprintf("HostValue{Name: \"%s\", IpAddress: \"%s\", HostName: \"%s\", Roles: %v}",
		hv.Name, hv.IpAddress, hv.HostName, hv.Roles)
}

type HostGroupsConfig struct {
	HostGroups []HostGroups `yaml:"groups" json:"groups" xml:"groups,chardata"`
}

func (hgc *HostGroupsConfig) String() string {
	var hostsVal string = "["
	for _, host := range hgc.HostGroups {
		prefix := ", "
		if len(hostsVal) == 0 {
			prefix = ""
		}
		hostsVal += prefix + host.String()
	}
	hostsVal += "]"
	return fmt.Sprintf("HostGroupsConfig{HostGroups: %s}",
		hostsVal)
}

func (hgc *HostGroupsConfig) Yaml() (string, error) {
	return io.ToYaml(hgc)
}

func (hgc *HostGroupsConfig) FromYamlFile(path string) (*HostGroupsConfig, error) {
	itf, err := io.FromYamlFile(path, hgc)
	if err != nil {
		return nil, err
	}
	var hostsObj *HostGroupsConfig = itf.(*HostGroupsConfig)
	return hostsObj, nil
}

func (hgc *HostGroupsConfig) FromYamlCode(yamlCode string) (*HostGroupsConfig, error) {
	itf, err := io.FromYamlCode(yamlCode, hgc)
	if err != nil {
		return nil, err
	}
	var hostsObj *HostGroupsConfig = itf.(*HostGroupsConfig)
	return hostsObj, nil
}

func (hgc *HostGroupsConfig) Xml() (string, error) {
	return io.ToXml(hgc)
}

func (hgc *HostGroupsConfig) FromXmlFile(path string) (*HostGroupsConfig, error) {
	itf, err := io.FromXmlFile(path, hgc)
	if err != nil {
		return nil, err
	}
	var hostsObj *HostGroupsConfig = itf.(*HostGroupsConfig)
	return hostsObj, nil
}

func (hgc *HostGroupsConfig) FromXmlCode(xmlCode string) (*HostGroupsConfig, error) {
	itf, err := io.FromXmlCode(xmlCode, hgc)
	if err != nil {
		return nil, err
	}
	var hostsObj *HostGroupsConfig = itf.(*HostGroupsConfig)
	return hostsObj, nil
}

func (hgc *HostGroupsConfig) Json() (string, error) {
	return io.ToJson(hgc)
}

func (hgc *HostGroupsConfig) FromJsonFile(path string) (*HostGroupsConfig, error) {
	itf, err := io.FromJsonFile(path, hgc)
	if err != nil {
		return nil, err
	}
	var hostsObj *HostGroupsConfig = itf.(*HostGroupsConfig)
	return hostsObj, nil
}

func (hgc *HostGroupsConfig) FromJsonCode(jsonCode string) (*HostGroupsConfig, error) {
	itf, err := io.FromJsonCode(jsonCode, hgc)
	if err != nil {
		return nil, err
	}
	var hostsObj *HostGroupsConfig = itf.(*HostGroupsConfig)
	return hostsObj, nil
}
