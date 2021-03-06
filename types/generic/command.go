package generic

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/hellgate75/go-deploy/cmd"
	"github.com/hellgate75/go-tcp-common/log"
	"github.com/hellgate75/go-deploy/types/module"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"
)

var Logger log.Logger = nil

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
	for _, command := range feed.Steps {
		var commandMap = map[interface{}]interface{}(command)
		var name string = ""
		if val, ok := commandMap["name"]; ok {
			name = fmt.Sprintf("%v", val)
		} else if val, ok := commandMap["NAME"]; ok {
			name = fmt.Sprintf("%v", val)
		}
		for key, value := range command {
			if key != "name" && key != "NAME" {
				stepsX, errorsX := EvaluateSteps(name, key, value)
				for _, stepX := range stepsX {
					steps = append(steps, stepX)
				}
				for _, errX := range errorsX {
					errors = append(errors, errX)
				}
			}
		}
	}
	return steps, errors
}

//Feed Interface, that describes the available option for the load of the file
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
	dFormat := cmd.GetFileFormatDescritor(path, module.RuntimeDeployConfig.ConfigLang)
	if dFormat == module.YAML_DESCRIPTOR {
		err = yaml.Unmarshal(data, feed)
	} else if dFormat == module.XML_DESCRIPTOR {
		err = xml.Unmarshal(data, feed)
	} else if dFormat == module.JSON_DESCRIPTOR {
		err = json.Unmarshal(data, feed)
	} else {
		return errors.New("Feed.Load: Unavailable converter for type: " + string(dFormat))
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
	var errorList []error = make([]error, 0)
	if feed.HostGroup == "" {
		errorList = append(errorList, errors.New("Uanble to validate a feed without hosts 'group'"))
	}
	var steps []*module.Step = make([]*module.Step, 0)
	for _, command := range feed.Steps {
		var commandMap = map[interface{}]interface{}(command)
		var name string = ""
		if val, ok := commandMap["name"]; ok {
			name = fmt.Sprintf("%v", val)
		} else if val, ok := commandMap["NAME"]; ok {
			name = fmt.Sprintf("%v", val)
		}
		for key, value := range command {
			if key != "name" && key != "NAME" {
				stepsX, errorsX := EvaluateSteps(name, key, value)
				for _, stepX := range stepsX {
					steps = append(steps, stepX)
				}
				for _, errX := range errorsX {
					errorList = append(errorList, errX)
				}
			}
		}
	}
	return &module.FeedExec{
		Name:      feed.Name,
		HostGroup: feed.HostGroup,
		Steps:     steps,
	}, errorList
}

//Internl function that transforms Blob data in list of module.Step Structure pointers
func EvaluateSteps(name string, key interface{}, value interface{}) ([]*module.Step, []error) {
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
			Logger.Tracef("valueType: %v", valueType)
			if strings.Index(valueType, "map[") == 0 {

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
						steps = append(steps, NewImportStep(name, feeds))
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
						steps = append(steps, NewImportStep(name, feeds))
					} else {
						errorsList = append(errorsList, errors.New(fmt.Sprintf("Invalid import type %v, expected []string", valueType)))
					}

				} else {
					step, err := NewStep(name, fmt.Sprintf("%v", key), value)
					if err != nil {
						errorsList = append(errorsList, err)
					} else {
						steps = append(steps, step)
					}
				}
			} else {
				//				for key, value := range map[interface{}]interface{}(value) {
				//					stepsX, errorsX := EvaluateSteps(key, value)
				//					for _, stepX := range stepsX {
				//						steps = append(steps, stepX)
				//					}
				//					for _, errX := range errorsX {
				//						errorsList = append(errorsList, errX)
				//					}
				//				}
				errorsList = append(errorsList, errors.New("Value type: "+valueType+" is not expected one (map[interface{}]interface{})"))
			}
		}
	} else {
		errorsList = append(errorsList, errors.New("Key type: "+keyType+" is not expected one (string)"))
	}

	return steps, errorsList
}

//Create new empty Feed interface, ready for load and/or save
func NewFeed(defaultName string) IFeed {
	return &Feed{
		Name:  defaultName,
		Steps: make([]map[interface{}]interface{}, 0),
	}
}
