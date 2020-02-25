package types

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

type DeployConfig struct {
	DeployName string
	UseHosts   []string
	UseVars    []string
	ConfigDir  string
}

func (dc *DeployConfig) Yaml() (string, error) {
	return toYaml(dc)
}

func (dc *DeployConfig) FromFile(path string) (*DeployConfig, error) {
	itf, err := fromFile(path, dc)
	if err != nil {
		return nil, err
	}
	var conf *DeployConfig = itf.(*DeployConfig)
	return conf, nil
}

func fromFile(path string, itf interface{}) (interface{}, error) {
	_, errS := os.Stat(path)
	if errS != nil {
		return nil, errors.New("fromFile::Stats: " + errS.Error())
	}
	file, errF := os.Open(path)
	if errF != nil {
		return nil, errors.New("fromFile::OpenFile: " + errF.Error())
	}
	bytes, errR := ioutil.ReadAll(file)
	if errR != nil {
		return nil, errors.New("fromFile::ReadFile: " + errR.Error())
	}
	err := yaml.Unmarshal(bytes, itf)
	if err != nil {
		return nil, errors.New("fromFile::Unmarshal: " + err.Error())
	} else {
		return itf, nil
	}
}

func toYaml(itf interface{}) (string, error) {
	bytes, err := yaml.Marshal(itf)
	if err != nil {
		return "", errors.New("toYaml::Marchal: " + err.Error())
	} else {
		return fmt.Sprintf("\n%s", bytes), nil
	}
}
