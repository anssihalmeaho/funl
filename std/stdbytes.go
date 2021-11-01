package std

import (
	"bytes"
	"fmt"

	"github.com/anssihalmeaho/funl/funl"
)

func initSTDBytes(interpreter *funl.Interpreter) (err error) {
	stdModuleName := "stdbytes"
	topFrame := funl.NewTopFrameWithInterpreter(interpreter)
	stdBytesFuncs := []stdFuncInfo{
		{
			Name:       "str-to-bytes",
			Getter:     getStdBytesStrToBytes,
			IsFunction: true,
		},
		{
			Name:       "string",
			Getter:     getStdBytesString,
			IsFunction: true,
		},
		{
			Name:       "count",
			Getter:     getStdBytesCount,
			IsFunction: true,
		},
		{
			Name:       "new",
			Getter:     getStdBytesNew,
			IsFunction: true,
		},
		{
			Name:       "split-by",
			Getter:     getStdBytesSplitBy,
			IsFunction: true,
		},
	}
	err = setSTDFunctions(topFrame, stdModuleName, stdBytesFuncs, interpreter)

	nlVal := &OpaqueByteArray{data: []byte("\n")}
	item := &funl.Item{Type: funl.ValueItem, Data: funl.Value{Kind: funl.OpaqueValue, Data: nlVal}}
	err = topFrame.Syms.Add("nl", item)
	if err != nil {
		return
	}

	return
}

type OpaqueByteArray struct {
	data []byte
}

// Creates new OpaqueByteArray
func NewOpaqueByteArray(data []byte) *OpaqueByteArray {
	return &OpaqueByteArray{data: data}
}

func (ob *OpaqueByteArray) TypeName() string {
	return "bytearray"
}

func (ob *OpaqueByteArray) Str() string {
	s := ""
	for _, v := range ob.data {
		s += fmt.Sprintf("%x ", v)
	}
	if s != "" {
		s = s[:len(s)-1]
	}
	return fmt.Sprintf("bytearray(%s)", s)
}

func (ob *OpaqueByteArray) Equals(with funl.OpaqueAPI) bool {
	other, ok := with.(*OpaqueByteArray)
	if !ok {
		return false
	}
	return bytes.Equal(ob.data, other.data)
}

func getStdBytesSplitBy(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 2 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need one", name, l)
		}
		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		byteArray, ok := arguments[0].Data.(*OpaqueByteArray)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not bytearray value", name)
		}
		if arguments[1].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		separatorByteArray, ok := arguments[1].Data.(*OpaqueByteArray)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not bytearray value", name)
		}

		var resultList []funl.Value
		for _, bslice := range bytes.Split(byteArray.data, separatorByteArray.data) {
			resultList = append(resultList, funl.Value{Kind: funl.OpaqueValue, Data: &OpaqueByteArray{data: bslice}})
		}
		retVal = funl.MakeListOfValues(frame, resultList)
		return
	}
}

func getStdBytesNew(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need one", name, l)
		}
		if arguments[0].Kind != funl.ListValue {
			funl.RunTimeError2(frame, "%s: requires list as argument", name)
		}
		byteArray := &OpaqueByteArray{}
		ite := funl.NewListIterator(arguments[0])
		for {
			val := ite.Next()
			if val == nil {
				break
			}
			switch val.Kind {
			case funl.IntValue:
				if num := val.Data.(int); num > 255 {
					funl.RunTimeError2(frame, "%s: expects int value less than 256 (%d)", name, num)
				} else if num < 0 {
					funl.RunTimeError2(frame, "%s: expects int value more than 0 (%d)", name, num)
				} else {
					byteVal := byte(num)
					byteArray.data = append(byteArray.data, byteVal)
				}
			case funl.StringValue:
				s := val.Data.(string)
				for _, byteVal := range []byte(s) {
					byteArray.data = append(byteArray.data, byteVal)
				}
			case funl.OpaqueValue:
				if byteArr, ok := val.Data.(*OpaqueByteArray); !ok {
					funl.RunTimeError2(frame, "%s: unsupported type: %s", name, val.Data)
				} else {
					for _, byteVal := range byteArr.data {
						byteArray.data = append(byteArray.data, byteVal)
					}
				}
			default:
				funl.RunTimeError2(frame, "%s: unsupported type", name)
			}
		}
		retVal = funl.Value{Kind: funl.OpaqueValue, Data: byteArray}
		return
	}
}

func getStdBytesStrToBytes(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need one", name, l)
		}
		if arguments[0].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: requires string value", name)
		}
		strVal, ok := arguments[0].Data.(string)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not string value", name)
		}
		buf := bytes.NewBufferString(strVal)
		retVal = funl.Value{Kind: funl.OpaqueValue, Data: &OpaqueByteArray{data: buf.Bytes()}}
		return
	}
}

func getStdBytesString(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need one", name, l)
		}
		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		byteArray, ok := arguments[0].Data.(*OpaqueByteArray)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not bytearray value", name)
		}
		buf := bytes.NewBuffer(byteArray.data)
		str := buf.String()
		if str == "<nil>" {
			funl.RunTimeError2(frame, "%s: could not convert bytearray to string", name)
		}
		retVal = funl.Value{Kind: funl.StringValue, Data: str}
		return
	}
}

func getStdBytesCount(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need one", name, l)
		}
		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		byteArray, ok := arguments[0].Data.(*OpaqueByteArray)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not bytearray value", name)
		}
		retVal = funl.Value{Kind: funl.IntValue, Data: len(byteArray.data)}
		return
	}
}
