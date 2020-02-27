package generic

import ()

type Feed struct {
	Name    string                      `yaml:"name,omitempty" json:"name,omitempty" xml:"name,chardata,omitempty"`
	Options map[interface{}]interface{} `yaml:"options,omitempty" json:"options,omitempty" xml:"option,chardata,omitempty"`
}

type OptionsSet struct {
	Options map[interface{}]interface{} `yaml:",omitempty" json:"options,omitempty" xml:"option,chardata,omitempty"`
}
