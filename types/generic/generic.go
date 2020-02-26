package generic

import ()

/*
* Coverter interface, responsible to comvert raw interface from the parsing to a specific structure
 */
type Converter interface {
	/*
	* Converts a raw interface element to a command qualified structure <BR/>
	* Paramameters: <BR/>
	* cmdValues (interface{}) Raw value from the feed file parsing
	* Return: <BR/>
	* (interface{}) Qualified structure <BR/>
	* (error) Error occured during any conversion <BR/>
	 */
	Convert(cmdValues interface{}) (interface{}, error)
}

type Executor interface {
	Execute(step Step) error
}

type Step struct {
	StepType string
	StepData interface{}
	Children []*Step
	Feeds    []*FeedExec
}

type FeedExec struct {
	Steps []*Step
}
