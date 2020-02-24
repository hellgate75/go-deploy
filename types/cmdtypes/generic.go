package cmdtypes

import (
	"errors"
	"strings"
	
)

const (
	FEED_TYPE_IMPORT =  iota + 1
	FEED_TYPE_SHELL
	FEED_TYPE_SERVICE
	FEED_TYPE_FACT
)




type FeedExec struct {
	Steps	[]*Step
}

func KeyToType(key string) (int, error) {
	switch strings.ToLower(key) {
		case "import":
			return FEED_TYPE_IMPORT, nil
		case "shell":
			return FEED_TYPE_IMPORT, nil
		default:
			return 0, errors.New("Unable to decode key: " + key)
	}
}

type Step struct {
	StepType	int
	StepData	interface{}
	Children	[]Step
	Feeds 		[]*FeedExec
}

func NewStep(stepType int, stepData	interface{}) *Step {
	return &Step {
		StepType: stepType,
		StepData: stepData,
		Children: make([]Step, 0),
		Feeds: make([]*FeedExec, 0),
	}
}

func NewStepWtihChildren(stepType int, stepData	interface{}, children []Step) *Step {
	return &Step {
		StepType: stepType,
		StepData: stepData,
		Children: children,
		Feeds: make([]*FeedExec, 0),
	}
}

func NewImportStep(feeds []*FeedExec) *Step {
	return &Step {
		StepType: FEED_TYPE_IMPORT,
		StepData: nil,
		Children: make([]Step, 0),
		Feeds: feeds,
	}
}

func NewImportStepWithChildren(feeds []*FeedExec, children []Step) *Step {
	return &Step {
		StepType: FEED_TYPE_IMPORT,
		StepData: nil,
		Children: children,
		Feeds: feeds,
	}
}

