package generic

import (
	"github.com/hellgate75/go-deploy/types/module"
)

// Create New module.Step by given name, step classifier, step data (blob of data to be converted)
func NewStep(name string, stepType string, stepData interface{}) (*module.Step, error) {
	data, err := NewConverter(stepType).Convert(stepData)
	if err != nil {
		return nil, err
	}
	return &module.Step{
		Name:     name,
		StepType: stepType,
		StepData: data,
		Children: make([]*module.Step, 0),
		Feeds:    make([]*module.FeedExec, 0),
	}, nil
}

// Create New module.Step by given name, step classifier, step data (blob of data to be converted) and children steps
func NewStepWtihChildren(name string, stepType string, stepData interface{}, children []*module.Step) (*module.Step, error) {
	data, err := NewConverter(stepType).Convert(stepData)
	if err != nil {
		return nil, err
	}
	return &module.Step{
		Name:     name,
		StepType: stepType,
		StepData: data,
		Children: children,
		Feeds:    make([]*module.FeedExec, 0),
	}, nil
}

// Create New module.Step by given name, and generic.Feed children list
func NewImportStep(name string, feeds []*module.FeedExec) *module.Step {
	return &module.Step{
		Name:     name,
		StepType: "import",
		StepData: nil,
		Children: make([]*module.Step, 0),
		Feeds:    feeds,
	}
}

// Create New module.Step by given name, and generic.Feed children list, with children module.Step elements
func NewImportStepWithChildren(name string, feeds []*module.FeedExec, children []*module.Step) *module.Step {
	return &module.Step{
		Name:     name,
		StepType: "import",
		StepData: nil,
		Children: children,
		Feeds:    feeds,
	}
}
