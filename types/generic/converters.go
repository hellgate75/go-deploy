package generic

import (
	"errors"
	"github.com/hellgate75/go-deploy/modules"
	"github.com/hellgate75/go-deploy/modules/meta"
	"github.com/hellgate75/go-deploy/types/threads"
	"github.com/hellgate75/go-tcp-common/log"
)

// Unknown Converter Structure
type NilCommandConverter struct {
	CmdType string
}

func (nilCommand *NilCommandConverter) SetLogger(l log.Logger) {
}

func (nilCommand *NilCommandConverter) Convert(cmdValues interface{}) (threads.StepRunnable, error) {
	return nil, errors.New("NilCommandConverter -> Not implemented type: " + nilCommand.CmdType)

}

var convertersMap map[string]meta.Converter = make(map[string]meta.Converter)

// Creates a new meta.Converter based on the command name, crowling in the bult-in and custom libraries (if available)
func NewConverter(cmdType string) meta.Converter {
	Logger.Debugf("NewConverter -> cmdType: %s", cmdType)
	//Verify local coverters cache
	if _, ok := convertersMap[cmdType]; ok {
		return convertersMap[cmdType]
	}
	converter, err := modules.LoadConverterForModule(cmdType)
	if err != nil {
		return &NilCommandConverter{
			CmdType: cmdType,
		}
	}
	if _, ok := convertersMap[cmdType]; !ok && converter != nil {
		//Store in local coverters cache
		convertersMap[cmdType] = converter
	}
	return converter

}
