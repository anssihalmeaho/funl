package std

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/anssihalmeaho/funl/funl"
	"gopkg.in/yaml.v3"
)

func initSTDYaml(interpreter *funl.Interpreter) (err error) {
	stdModuleName := "stdyaml"
	topFrame := funl.NewTopFrameWithInterpreter(interpreter)
	stdFuncs := []stdFuncInfo{
		{
			Name:       "decode",
			Getter:     getStdYAMLdecode,
			IsFunction: true,
		},
	}
	item := &funl.Item{Type: funl.ValueItem, Data: funl.Value{Kind: funl.OpaqueValue, Data: &OpaqueYAMLnull{}}}
	err = topFrame.Syms.Add("null", item)
	if err != nil {
		return
	}
	err = setSTDFunctions(topFrame, stdModuleName, stdFuncs, interpreter)
	return
}

type OpaqueYAMLnull struct{}

func (opa *OpaqueYAMLnull) TypeName() string {
	return "yaml-null"
}

func (opa *OpaqueYAMLnull) Str() string {
	return "yaml-null"
}

func (opa *OpaqueYAMLnull) Equals(with funl.OpaqueAPI) bool {
	_, ok := with.(*OpaqueYAMLnull)
	return ok
}

// decode(opaque:bytearray) -> list(bool, string, <VALUE>)
func getStdYAMLdecode(name string) stdFuncType {
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

		ok, errText, val := decodeYAML(name, frame, byteArray.data)
		if !ok {
			val = funl.Value{Kind: funl.StringValue, Data: ""}
		}
		values := []funl.Value{
			{
				Kind: funl.BoolValue,
				Data: ok,
			},
			{
				Kind: funl.StringValue,
				Data: errText,
			},
			val,
		}
		retVal = funl.MakeListOfValues(frame, values)
		return
	}
}

func decodeYAML(name string, frame *funl.Frame, indata []byte) (ok bool, errText string, val funl.Value) {
	defer func() {
		if r := recover(); r != nil {
			err, _ := r.(error)
			ok, errText = false, fmt.Sprintf("%s: %v", name, err)
		}
	}()

	var res any
	yamlDecoder := yaml.NewDecoder(bytes.NewReader(indata))
	err := yamlDecoder.Decode(&res)
	if err != nil {
		ok, errText = false, fmt.Sprintf("%s: error in unmarshal: %v", name, err)
		return
	}
	val = traverseYAML(frame, res)
	ok = true
	return
}

func traverseYAML(frame *funl.Frame, intf any) funl.Value {
	val := reflect.ValueOf(intf)
	if intf == nil {
		return funl.Value{Kind: funl.OpaqueValue, Data: &OpaqueYAMLnull{}}
	}
	switch val.Kind() {
	case reflect.Bool:
		return funl.Value{Kind: funl.BoolValue, Data: val.Bool()}
	case reflect.String:
		return funl.Value{Kind: funl.StringValue, Data: val.String()}
	case reflect.Int:
		return funl.Value{Kind: funl.IntValue, Data: int(val.Int())}
	case reflect.Float64:
		return funl.Value{Kind: funl.FloatValue, Data: val.Float()}
	case reflect.Slice:
		var values []funl.Value
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i)
			listItem := traverseYAML(frame, item.Interface())
			values = append(values, listItem)
		}
		return funl.MakeListOfValues(frame, values)
	case reflect.Map:
		mapval := funl.HandleMapOP(frame, []*funl.Item{})
		for _, k := range val.MapKeys() {
			v := val.MapIndex(k)
			keyv := funl.Value{Kind: funl.StringValue, Data: k.String()}
			valv := traverseYAML(frame, v.Interface())
			putArgs := []*funl.Item{
				{Type: funl.ValueItem, Data: mapval},
				{Type: funl.ValueItem, Data: keyv},
				{Type: funl.ValueItem, Data: valv},
			}
			mapval = funl.HandlePutOP(frame, putArgs)
			if mapval.Kind != funl.MapValue {
				panic(fmt.Errorf("failed to put key-value to map"))
			}
		}
		return mapval
	}
	panic(fmt.Errorf("Unexpected type: %v (%v)", val, val.Kind()))
}
