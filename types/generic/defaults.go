package generic

import (
	"fmt"
	"github.com/hellgate75/go-deploy/io"
)

type NameValue struct {
	Name  string `yaml:"name" json:"name" xml:"name,chardata"`
	Value string `yaml:"value,omitempty" json:"value,omitempty" xml:"value,chardata,omitempty"`
}

type Vars struct {
	Vars []NameValue `yaml:"vars,omitempty" json:"vars,omitempty" xml:"vars,chardata,omitempty"`
}

type Environments struct {
	Envs []NameValue `yaml:"environments,omitempty" json:"environments,omitempty" xml:"environments,chardata,omitempty"`
}

type HostValue struct {
	Name      string   `yaml:"name" json:"name" xml:"name,chardata"`
	IpAddress string   `yaml:"ipAddress,omitempty" json:"ipAddress,omitempty" xml:"ip-address,chardata,omitempty"`
	HostName  string   `yaml:"hostName,omitempty" json:"hostName,omitempty" xml:"host-name,chardata,omitempty"`
	Roles     []string `yaml:"roles,omitempty" json:"roles,omitempty" xml:"roles,chardata,omitempty"`
}

func (hv *HostValue) String() string {
	return fmt.Sprintf("HostValue{Name: \"%s\", IpAddress: \"%s\", HostName: \"%s\", Roles: %v}",
		hv.Name, hv.IpAddress, hv.HostName, hv.Roles)
}

type Hosts struct {
	Hosts []Hosts `yaml:"hosts,omitempty" json:"hosts,omitempty" xml:"hosts,chardata,omitempty"`
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

func (hosts *Hosts) Yaml() (string, error) {
	return io.ToYaml(hosts)
}

func (hosts *Hosts) FromYamlFile(path string) (*Hosts, error) {
	itf, err := io.FromYamlFile(path, hosts)
	if err != nil {
		return nil, err
	}
	var hostsObj *Hosts = itf.(*Hosts)
	return hostsObj, nil
}

func (hosts *Hosts) FromYamlCode(yamlCode string) (*Hosts, error) {
	itf, err := io.FromYamlCode(yamlCode, hosts)
	if err != nil {
		return nil, err
	}
	var hostsObj *Hosts = itf.(*Hosts)
	return hostsObj, nil
}

func (hosts *Hosts) Xml() (string, error) {
	return io.ToXml(hosts)
}

func (hosts *Hosts) FromXmlFile(path string) (*Hosts, error) {
	itf, err := io.FromXmlFile(path, hosts)
	if err != nil {
		return nil, err
	}
	var hostsObj *Hosts = itf.(*Hosts)
	return hostsObj, nil
}

func (hosts *Hosts) FromXmlCode(xmlCode string) (*Hosts, error) {
	itf, err := io.FromXmlCode(xmlCode, hosts)
	if err != nil {
		return nil, err
	}
	var hostsObj *Hosts = itf.(*Hosts)
	return hostsObj, nil
}

func (hosts *Hosts) Json() (string, error) {
	return io.ToJson(hosts)
}

func (hosts *Hosts) FromJsonFile(path string) (*Hosts, error) {
	itf, err := io.FromJsonFile(path, hosts)
	if err != nil {
		return nil, err
	}
	var hostsObj *Hosts = itf.(*Hosts)
	return hostsObj, nil
}

func (hosts *Hosts) FromJsonCode(jsonCode string) (*Hosts, error) {
	itf, err := io.FromJsonCode(jsonCode, hosts)
	if err != nil {
		return nil, err
	}
	var hostsObj *Hosts = itf.(*Hosts)
	return hostsObj, nil
}
