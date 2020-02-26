package cmdtypes

import (
	"errors"
	"github.com/hellgate75/go-deploy/modules"
	"github.com/hellgate75/go-deploy/types/generic"
	"reflect"
)

/*
* Coverter interface, responsible to comvert raw interface from the parsing to a specific structure
 */
type Printable interface {
	/*
	* Traslates the object in printable version <BR/>
	* Return: <BR/>
	* (string) Representation of the structure<BR/>
	 */
	String() string
}

/*
* Unknown command structure
 */
type NilCommandConverter struct {
	CmdType string
}

func (nilCommand *NilCommandConverter) Convert(cmdValues interface{}) (interface{}, error) {
	return nil, errors.New("Not implemented type: " + nilCommand.CmdType)

}

var ERROR_TYPE reflect.Type = reflect.TypeOf(errors.New(""))

func NewConverter(cmdType string) generic.Converter {
	//	switch cmdType {
	//	case FEED_TYPE_SHELL:
	//		return &ShellCommand{}
	//	case FEED_TYPE_SERVICE:
	//		return &ServiceCommand{}
	//	case FEED_TYPE_COPY:
	//		return &CopyCommand{}
	//	}
	converter, err := modules.LoadConverterForModule(cmdType)
	if err != nil {
		return &NilCommandConverter{
			CmdType: cmdType,
		}
	}
	return converter

}
