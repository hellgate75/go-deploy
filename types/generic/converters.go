package generic

import (
	"errors"
	"fmt"
	"github.com/hellgate75/go-deploy/modules"
)

/*
* Unknown command structure
 */
type NilCommandConverter struct {
	CmdType string
}

func (nilCommand *NilCommandConverter) Convert(cmdValues interface{}) (interface{}, error) {
	return nil, errors.New("NilCommandConverter -> Not implemented type: " + nilCommand.CmdType)

}

func NewConverter(cmdType string) modules.Converter {
	//	switch cmdType {
	//	case FEED_TYPE_SHELL:
	//		return &ShellCommand{}
	//	case FEED_TYPE_SERVICE:
	//		return &ServiceCommand{}
	//	case FEED_TYPE_COPY:
	//		return &CopyCommand{}
	//	}
	Logger.Warn(fmt.Sprintf("NewConverter -> cmdType: %s", cmdType))
	converter, err := modules.LoadConverterForModule(cmdType)
	if err != nil {
		return &NilCommandConverter{
			CmdType: cmdType,
		}
	}
	return converter

}
