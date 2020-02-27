package generic

import (
	"errors"
	"github.com/hellgate75/go-deploy/modules"
	"github.com/hellgate75/go-deploy/types/module"
)

/*
* Unknown command structure
 */
type NilCommandExecutor struct {
	CmdType string
}

func (nilCommand *NilCommandExecutor) Execute(step module.Step) error {
	return errors.New("Not implemented type: " + nilCommand.CmdType)

}

func NewExecutor(cmdType string) modules.Executor {
	//	switch cmdType {
	//	case FEED_TYPE_SHELL:
	//		return &ShellExecutor{}
	//	case FEED_TYPE_SERVICE:
	//		return &ServiceExecutor{}
	//	case FEED_TYPE_COPY:
	//		return &CopyExecutor{}
	//	}
	executor, err := modules.LoadExecutorForModule(cmdType)
	if err != nil {
		return &NilCommandExecutor{
			CmdType: cmdType,
		}
	}
	return executor

}
