package generic

import ()

type Feed struct {
	Name  string                        `yaml:"name,omitempty" json:"name,omitempty" xml:"name,chardata,omitempty"`
	Steps []map[interface{}]interface{} `yaml:"steps,omitempty" json:"steps,omitempty" xml:"steps,chardata,omitempty"`
}

type OptionsSet struct {
	Steps []map[interface{}]interface{} `yaml:",omitempty" json:"steps,omitempty" xml:"steps,chardata,omitempty"`
}
