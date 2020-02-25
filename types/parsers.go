package types

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

func fromYamlCode(yamlCode string, itf interface{}) (interface{}, error) {
	err := yaml.Unmarshal([]byte(yamlCode), itf)
	if err != nil {
		return nil, errors.New("fromYamlCode::Unmarshal: " + err.Error())
	} else {
		return itf, nil
	}
}

func fromJsonCode(jsonCode string, itf interface{}) (interface{}, error) {
	err := json.Unmarshal([]byte(jsonCode), itf)
	if err != nil {
		return nil, errors.New("fromJsonCode::Unmarshal: " + err.Error())
	} else {
		return itf, nil
	}
}

func fromXmlCode(xmlCode string, itf interface{}) (interface{}, error) {
	err := xml.Unmarshal([]byte(xmlCode), itf)
	if err != nil {
		return nil, errors.New("fromXmlCode::Unmarshal: " + err.Error())
	} else {
		return itf, nil
	}
}

func fromYamlFile(path string, itf interface{}) (interface{}, error) {
	_, errS := os.Stat(path)
	if errS != nil {
		return nil, errors.New("fromYamlFile::Stats: " + errS.Error())
	}
	file, errF := os.Open(path)
	if errF != nil {
		return nil, errors.New("fromYamlFile::OpenFile: " + errF.Error())
	}
	bytes, errR := ioutil.ReadAll(file)
	if errR != nil {
		return nil, errors.New("fromYamlFile::ReadFile: " + errR.Error())
	}
	err := yaml.Unmarshal(bytes, itf)
	if err != nil {
		return nil, errors.New("fromYamlFile::Unmarshal: " + err.Error())
	} else {
		return itf, nil
	}
}

func fromJsonFile(path string, itf interface{}) (interface{}, error) {
	_, errS := os.Stat(path)
	if errS != nil {
		return nil, errors.New("fromJsonFile::Stats: " + errS.Error())
	}
	file, errF := os.Open(path)
	if errF != nil {
		return nil, errors.New("fromJsonFile::OpenFile: " + errF.Error())
	}
	bytes, errR := ioutil.ReadAll(file)
	if errR != nil {
		return nil, errors.New("fromJsonFile::ReadFile: " + errR.Error())
	}
	err := json.Unmarshal(bytes, itf)
	if err != nil {
		return nil, errors.New("fromJsonFile::Unmarshal: " + err.Error())
	} else {
		return itf, nil
	}
}

func fromXmlFile(path string, itf interface{}) (interface{}, error) {
	_, errS := os.Stat(path)
	if errS != nil {
		return nil, errors.New("fromXmlFile::Stats: " + errS.Error())
	}
	file, errF := os.Open(path)
	if errF != nil {
		return nil, errors.New("fromXmlFile::OpenFile: " + errF.Error())
	}
	bytes, errR := ioutil.ReadAll(file)
	if errR != nil {
		return nil, errors.New("fromXmlFile::ReadFile: " + errR.Error())
	}
	err := xml.Unmarshal(bytes, itf)
	if err != nil {
		return nil, errors.New("fromXmlFile::Unmarshal: " + err.Error())
	} else {
		return itf, nil
	}
}

func toYaml(itf interface{}) (string, error) {
	bytes, err := yaml.Marshal(itf)
	if err != nil {
		return "", errors.New("toYaml::Marshal: " + err.Error())
	} else {
		return fmt.Sprintf("\n%s", bytes), nil
	}
}

func toJson(itf interface{}) (string, error) {
	bytes, err := json.Marshal(itf)
	if err != nil {
		return "", errors.New("toJson::Marshal: " + err.Error())
	} else {
		return fmt.Sprintf("\n%s", bytes), nil
	}
}

func toXml(itf interface{}) (string, error) {
	bytes, err := xml.MarshalIndent(itf, "", "  ")
	if err != nil {
		return "", errors.New("toJson::Marshal: " + err.Error())
	} else {
		return fmt.Sprintf("\n%s", bytes), nil
	}
}
