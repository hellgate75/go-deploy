package generic

import ()

// Feed Strcuture that contains row data, It will be parsed and validated becoming a pointer to module.FeedEx (Executable Feed)
type Feed struct {
	Name      string                        `yaml:"name,omitempty" json:"name,omitempty" xml:"name,chardata,omitempty"`
	HostGroup string                        `yaml:"group,omitempty" json:"group,omitempty" xml:"group,chardata,omitempty"`
	Steps     []map[interface{}]interface{} `yaml:"steps,omitempty" json:"steps,omitempty" xml:"steps,chardata,omitempty"`
}

// Fragment of Steps blob data, intended to to be converted in Validation phase becoming a list of one or more module.Step
type OptionsSet struct {
	Steps []map[interface{}]interface{} `yaml:",omitempty" json:"steps,omitempty" xml:"steps,chardata,omitempty"`
}
