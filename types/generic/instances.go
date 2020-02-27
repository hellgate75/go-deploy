package generic

import (
	"github.com/hellgate75/go-deploy/types/module"
)

func NewStep(stepType string, stepData interface{}) (*module.Step, error) {
	data, err := NewConverter(stepType).Convert(stepData)
	if err != nil {
		return nil, err
	}
	return &module.Step{
		StepType: stepType,
		StepData: data,
		Children: make([]*module.Step, 0),
		Feeds:    make([]*module.FeedExec, 0),
	}, nil
}

func NewStepWtihChildren(stepType string, stepData interface{}, children []*module.Step) (*module.Step, error) {
	data, err := NewConverter(stepType).Convert(stepData)
	if err != nil {
		return nil, err
	}
	return &module.Step{
		StepType: stepType,
		StepData: data,
		Children: children,
		Feeds:    make([]*module.FeedExec, 0),
	}, nil
}

func NewImportStep(feeds []*module.FeedExec) *module.Step {
	return &module.Step{
		StepType: "import",
		StepData: nil,
		Children: make([]*module.Step, 0),
		Feeds:    feeds,
	}
}

func NewImportStepWithChildren(feeds []*module.FeedExec, children []*module.Step) *module.Step {
	return &module.Step{
		StepType: "import",
		StepData: nil,
		Children: children,
		Feeds:    feeds,
	}
}
