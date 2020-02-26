package types

import (
	"errors"
	"fmt"
	"github.com/hellgate75/go-deploy/types/cmdtypes"
	"github.com/hellgate75/go-deploy/types/generic"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"
)

type Feed struct {
	Name    string                      `yaml:"name,omitempty"`
	Options map[interface{}]interface{} `yaml:",omitempty"`
}

type IFeed interface {
	Load(path string) error
	Save(path string) error
	Validate() bool
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
	err = yaml.Unmarshal(data, feed)
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
	data, err = yaml.Marshal(feed)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, data, 0666)
	if err != nil {
		return err
	}
	return nil
}

func (feed Feed) Validate() (*generic.FeedExec, []error) {
	var errors []error = make([]error, 0)
	var steps []*generic.Step = make([]*generic.Step, 0)
	for key, value := range feed.Options {
		stepsX, errorsX := EvaluateSteps(key, value)
		for _, stepX := range stepsX {
			steps = append(steps, stepX)
		}
		for _, errX := range errorsX {
			errors = append(errors, errX)
		}
	}
	return &generic.FeedExec{
		Steps: steps,
	}, errors
}

func EvaluateSteps(key interface{}, value interface{}) ([]*generic.Step, []error) {
	var errorsList []error = make([]error, 0)
	var steps []*generic.Step = make([]*generic.Step, 0)
	var keyType string = fmt.Sprintf("%T", key)
	var valueType string = fmt.Sprintf("%T", value)
	var err error
	if keyType == "string" {
		//keyStepType, err = cmdtypes.KeyToType(fmt.Sprintf("%v", key))
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
						var feeds []*generic.FeedExec = make([]*generic.FeedExec, 0)
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
						steps = append(steps, cmdtypes.NewImportStep(feeds))
					} else {
						errorsList = append(errorsList, errors.New(fmt.Sprintf("Invalid import type %v, expected []string", valueType)))
					}

				} else {
					step, err := cmdtypes.NewStep(keyType, value)
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
