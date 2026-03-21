package std

import (
	"github.com/anssihalmeaho/funl/funl"
	"github.com/gofrs/uuid/v5"
)

func initSTDuuid(interpreter *funl.Interpreter) (err error) {
	stdModuleName := "stduuid"
	topFrame := funl.NewTopFrameWithInterpreter(interpreter)
	stdFuncs := []stdFuncInfo{
		{
			Name:   "gen-rand",
			Getter: getStdUUIDGenRand,
		},
	}
	err = setSTDFunctions(topFrame, stdModuleName, stdFuncs, interpreter)
	return
}

func getStdUUIDGenRand(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		uuid, err := uuid.NewV4()
		if err != nil {
			funl.RunTimeError2(frame, "%s: failed to generate UUID: %v", name, err)
		}
		retVal = funl.Value{Kind: funl.StringValue, Data: uuid.String()}
		return
	}
}
