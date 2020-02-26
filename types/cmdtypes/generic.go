package cmdtypes

import (
	"github.com/hellgate75/go-deploy/types/generic"
)

func NewStep(stepType string, stepData interface{}) (*generic.Step, error) {
	data, err := NewConverter(stepType).Convert(stepData)
	if err != nil {
		return nil, err
	}
	return &generic.Step{
		StepType: stepType,
		StepData: data,
		Children: make([]*generic.Step, 0),
		Feeds:    make([]*generic.FeedExec, 0),
	}, nil
}

func NewStepWtihChildren(stepType string, stepData interface{}, children []*generic.Step) (*generic.Step, error) {
	data, err := NewConverter(stepType).Convert(stepData)
	if err != nil {
		return nil, err
	}
	return &generic.Step{
		StepType: stepType,
		StepData: data,
		Children: children,
		Feeds:    make([]*generic.FeedExec, 0),
	}, nil
}

func NewImportStep(feeds []*generic.FeedExec) *generic.Step {
	return &generic.Step{
		StepType: "import",
		StepData: nil,
		Children: make([]*generic.Step, 0),
		Feeds:    feeds,
	}
}

func NewImportStepWithChildren(feeds []*generic.FeedExec, children []*generic.Step) *generic.Step {
	return &generic.Step{
		StepType: "import",
		StepData: nil,
		Children: children,
		Feeds:    feeds,
	}
}
