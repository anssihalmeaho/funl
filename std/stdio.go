package std

import (
	"bufio"
	"fmt"
	"os"

	"github.com/anssihalmeaho/funl/funl"
)

func initSTDIO(interpreter *funl.Interpreter) (err error) {
	stdModuleName := "stdio"
	topFrame := funl.NewTopFrameWithInterpreter(interpreter)
	stdIOFuncs := []stdFuncInfo{
		{
			Name:   "printf",
			Getter: getStdIOPrintf,
		},
		{
			Name:   "printout",
			Getter: getStdIOPrintout,
		},
		{
			Name:   "printline",
			Getter: getStdIOPrintline,
		},
		{
			Name:   "readinput",
			Getter: getStdIOReadinput,
		},
	}
	err = setSTDFunctions(topFrame, stdModuleName, stdIOFuncs, interpreter)
	return
}

func getStdIOPrintf(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		var operands []*funl.Item
		for _, v := range arguments {
			operands = append(operands, &funl.Item{Type: funl.ValueItem, Data: v})
		}
		formattedStrVal := funl.HandleSprintfOP(frame, operands)
		if formattedStrVal.Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: assuming string (%#v)", name, formattedStrVal)
		}

		fmt.Printf(formattedStrVal.Data.(string))
		retVal = funl.Value{Kind: funl.BoolValue, Data: true}
		return
	}
}

func getStdIOPrintline(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		var s string
		for _, v := range arguments {
			var sval string
			switch v.Kind {
			case funl.StringValue:
				sval = v.Data.(string)
			default:
				sval = fmt.Sprintf("%v", v)
			}
			s += sval
		}
		fmt.Printf("%s\n", s)
		retVal = funl.Value{Kind: funl.BoolValue, Data: true}
		return
	}
}

func getStdIOPrintout(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		var s string
		for _, v := range arguments {
			var sval string
			switch v.Kind {
			case funl.StringValue:
				sval = v.Data.(string)
			default:
				sval = fmt.Sprintf("%v", v)
			}
			s += sval
		}
		fmt.Printf("%s", s)
		retVal = funl.Value{Kind: funl.BoolValue, Data: true}
		return
	}
}

func getStdIOReadinput(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		retVal = funl.Value{Kind: funl.StringValue, Data: scanner.Text()}
		return
	}
}
