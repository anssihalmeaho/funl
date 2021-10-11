package std

import (
	"encoding/base64"
	"fmt"

	"github.com/anssihalmeaho/funl/funl"
)

func initSTDbase64(interpreter *funl.Interpreter) (err error) {
	stdModuleName := "stdbase64"
	topFrame := &funl.Frame{
		Syms:     funl.NewSymt(),
		OtherNS:  make(map[funl.SymID]funl.ImportInfo),
		Imported: make(map[funl.SymID]*funl.Frame),
	}
	stdFuncs := []stdFuncInfo{
		{
			Name:       "encode",
			Getter:     getStdBase64encode,
			IsFunction: true,
		},
		{
			Name:       "decode",
			Getter:     getStdBase64decode,
			IsFunction: true,
		},
	}
	err = setSTDFunctions(topFrame, stdModuleName, stdFuncs, interpreter)
	return
}

// call(stdbase64.encode <opaque:bytearray>) -> list(ok err <string>)
func getStdBase64encode(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need one", name, l)
		}
		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: assuming opaque bytearray", name)
		}
		bytearray, convOK := arguments[0].Data.(*OpaqueByteArray)
		var errStr string
		var str string
		if convOK {
			str = base64.StdEncoding.EncodeToString(bytearray.data)
		} else {
			errStr = "assuming bytearray"
		}
		values := []funl.Value{
			{
				Kind: funl.BoolValue,
				Data: errStr == "",
			},
			{
				Kind: funl.StringValue,
				Data: errStr,
			},
			{
				Kind: funl.StringValue,
				Data: str,
			},
		}
		retVal = funl.MakeListOfValues(frame, values)
		return
	}
}

// call(stdbase64.decode <opaque:bytearray>) -> list(ok err <string>)
func getStdBase64decode(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need one", name, l)
		}
		if arguments[0].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: assuming string", name)
		}
		bytedata, err := base64.StdEncoding.DecodeString(arguments[0].Data.(string))
		var errStr string
		var bytearray *OpaqueByteArray
		if err != nil {
			errStr = fmt.Sprintf("Error in decoding: %v", err)
			bytearray = &OpaqueByteArray{data: []byte{}}
		} else {
			bytearray = &OpaqueByteArray{data: bytedata}
		}
		values := []funl.Value{
			{
				Kind: funl.BoolValue,
				Data: err == nil,
			},
			{
				Kind: funl.StringValue,
				Data: errStr,
			},
			{
				Kind: funl.OpaqueValue,
				Data: bytearray,
			},
		}
		retVal = funl.MakeListOfValues(frame, values)
		return
	}
}
