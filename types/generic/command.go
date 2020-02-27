package generic

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/hellgate75/go-deploy/types/module"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"
)

func (oset *OptionsSet) Load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	var data []byte
	data, err = ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	if string(module.RuntimeDeployConfig.ConfigLang) == "YAML" {
		err = yaml.Unmarshal(data, oset)
	} else if string(module.RuntimeDeployConfig.ConfigLang) == "XML" {
		err = xml.Unmarshal(data, oset)
	} else if string(module.RuntimeDeployConfig.ConfigLang) == "JSON" {
		err = json.Unmarshal(data, oset)
	} else {
		return errors.New("OptionsSet.Load: Unavailable converter for type: " + string(module.RuntimeDeployConfig.ConfigLang))
	}
	if err != nil {
		return err
	}
	return nil
}

func (oset *OptionsSet) Save(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		os.Remove(path)
	}
	var data []byte
	if string(module.RuntimeDeployConfig.ConfigLang) == "YAML" {
		data, err = yaml.Marshal(oset)
	} else if string(module.RuntimeDeployConfig.ConfigLang) == "XML" {
		data, err = xml.Marshal(oset)
	} else if string(module.RuntimeDeployConfig.ConfigLang) == "JSON" {
		data, err = json.Marshal(oset)
	} else {
		return errors.New("OptionsSet.Save: Unavailable converter for type: " + string(module.RuntimeDeployConfig.ConfigLang))
	}
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, data, 0666)
	if err != nil {
		return err
	}
	return nil
}

func (feed OptionsSet) Validate() ([]*module.Step, []error) {
	var errors []error = make([]error, 0)
	var steps []*module.Step = make([]*module.Step, 0)
	for key, value := range feed.Options {
		stepsX, errorsX := EvaluateSteps(key, value)
		for _, stepX := range stepsX {
			steps = append(steps, stepX)
		}
		for _, errX := range errorsX {
			errors = append(errors, errX)
		}
	}
	return steps, errors
}

type IFeed interface {
	Load(path string) error
	Save(path string) error
	Validate() (*module.FeedExec, []error)
}

func (feed *Feed) Load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	var data []byte
	data, err = ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	if string(module.RuntimeDeployConfig.ConfigLang) == "YAML" {
		err = yaml.Unmarshal(data, feed)
	} else if string(module.RuntimeDeployConfig.ConfigLang) == "XML" {
		err = xml.Unmarshal(data, feed)
	} else if string(module.RuntimeDeployConfig.ConfigLang) == "JSON" {
		err = json.Unmarshal(data, feed)
	} else {
		return errors.New("Feed.Load: Unavailable converter for type: " + string(module.RuntimeDeployConfig.ConfigLang))
	}
	if err != nil {
		return err
	}
	return nil
}

func (feed *Feed) Save(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		os.Remove(path)
	}
	var data []byte
	if string(module.RuntimeDeployConfig.ConfigLang) == "YAML" {
		data, err = yaml.Marshal(feed)
	} else if string(module.RuntimeDeployConfig.ConfigLang) == "XML" {
		data, err = xml.Marshal(feed)
	} else if string(module.RuntimeDeployConfig.ConfigLang) == "JSON" {
		data, err = json.Marshal(feed)
	} else {
		return errors.New("Feed.Save: Unavailable converter for type: " + string(module.RuntimeDeployConfig.ConfigLang))
	}
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, data, 0666)
	if err != nil {
		return err
	}
	return nil
}

func (feed Feed) Validate() (*module.FeedExec, []error) {
	var errors []error = make([]error, 0)
	var steps []*module.Step = make([]*module.Step, 0)
	for key, value := range feed.Options {
		stepsX, errorsX := EvaluateSteps(key, value)
		for _, stepX := range stepsX {
			steps = append(steps, stepX)
		}
		for _, errX := range errorsX {
			errors = append(errors, errX)
		}
	}
	return &module.FeedExec{
		Steps: steps,
	}, errors
}

func EvaluateSteps(key interface{}, value interface{}) ([]*module.Step, []error) {
	var errorsList []error = make([]error, 0)
	var steps []*module.Step = make([]*module.Step, 0)
	var keyType string = fmt.Sprintf("%T", key)
	var valueType string = fmt.Sprintf("%T", value)
	var err error
	if keyType == "string" {
		if err != nil {
			errorsList = append(errorsList, err)
		} else {
			//New Step
			var keyVal string = fmt.Sprintf("%v", key)
			if strings.Index(valueType, "map[") > 0 {
				for key, value := range value.(map[interface{}]interface{}) {
					stepsX, errorsX := EvaluateSteps(key, value)
					for _, stepX := range stepsX {
						steps = append(steps, stepX)
					}
					for _, errX := range errorsX {
						errorsList = append(errorsList, errX)
					}
				}
			} else {

				if strings.ToLower(keyVal) == "import" {
					if valueType == "[]string" || valueType == "[]interface{}" {
						var feeds []*module.FeedExec = make([]*module.FeedExec, 0)
						var arr []string = make([]string, 0)
						if valueType == "[]string" {
							for _, str := range value.([]string) {
								arr = append(arr, str)
							}
						} else {
							for _, iface := range value.([]interface{}) {
								arr = append(arr, fmt.Sprintf("%v", iface))
							}
						}
						for _, path := range arr {
							var cfeed *Feed = &Feed{}
							err := cfeed.Load(path)
							if err != nil {
								errorsList = append(errorsList, err)
							} else {
								fEx, exceptions := cfeed.Validate()
								if len(exceptions) > 0 {
									for _, errX := range exceptions {
										errorsList = append(errorsList, errX)
									}
								} else {
									feeds = append(feeds, fEx)
								}
							}

						}
						steps = append(steps, NewImportStep(feeds))
					} else {
						errorsList = append(errorsList, errors.New(fmt.Sprintf("Invalid import type %v, expected []string", valueType)))
					}

				} else if strings.ToLower(keyVal) == "include" {
					if valueType == "[]string" || valueType == "[]interface{}" {
						var feeds []*module.FeedExec = make([]*module.FeedExec, 0)
						var arr []string = make([]string, 0)
						if valueType == "[]string" {
							for _, str := range value.([]string) {
								arr = append(arr, str)
							}
						} else {
							for _, iface := range value.([]interface{}) {
								arr = append(arr, fmt.Sprintf("%v", iface))
							}
						}
						for _, path := range arr {
							var oset *OptionsSet = &OptionsSet{}
							err := oset.Load(path)
							if err != nil {
								errorsList = append(errorsList, err)
							} else {
								fSteps, exceptions := oset.Validate()
								if len(exceptions) > 0 {
									for _, errX := range exceptions {
										errorsList = append(errorsList, errX)
									}
								} else {
									steps = append(steps, fSteps...)
								}
							}

						}
						steps = append(steps, NewImportStep(feeds))
					} else {
						errorsList = append(errorsList, errors.New(fmt.Sprintf("Invalid import type %v, expected []string", valueType)))
					}

				} else {
					step, err := NewStep(keyType, value)
					if err != nil {
						errorsList = append(errorsList, err)
					} else {
						steps = append(steps, step)
					}
				}
			}
		}
	} else {

	}

	return steps, errorsList
}

func NewFeed(defaultName string) IFeed {
	return &Feed{
		Name:    defaultName,
		Options: make(map[interface{}]interface{}),
	}
}
