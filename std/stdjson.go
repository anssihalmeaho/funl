package std

import (
	"encoding/json"
	"fmt"
	"github.com/anssihalmeaho/funl"
	"reflect"
)

func initSTDJson() (err error) {
	stdModuleName := "stdjson"
	topFrame := &funl.Frame{
		Syms:     funl.NewSymt(),
		OtherNS:  make(map[funl.SymID]funl.ImportInfo),
		Imported: make(map[funl.SymID]*funl.Frame),
	}
	stdFuncs := []stdFuncInfo{
		{
			Name:       "encode",
			Getter:     getStdJSONencode,
			IsFunction: true,
		},
		{
			Name:       "decode",
			Getter:     getStdJSONdecode,
			IsFunction: true,
		},
	}
	err = setSTDFunctions(topFrame, stdModuleName, stdFuncs)

	item := &funl.Item{Type: funl.ValueItem, Data: funl.Value{Kind: funl.OpaqueValue, Data: &OpaqueJSONnull{}}}
	err = topFrame.Syms.Add("null", item)
	if err != nil {
		return
	}

	return
}

type OpaqueJSONnull struct{}

func (opa *OpaqueJSONnull) TypeName() string {
	return "json-null"
}

func (opa *OpaqueJSONnull) Str() string {
	return "json-null"
}

func (opa *OpaqueJSONnull) Equals(with funl.OpaqueAPI) bool {
	_, ok := with.(*OpaqueJSONnull)
	return ok
}

// encode(<VALUE>) -> list(bool, string, opaque:bytearray)
func getStdJSONencode(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need one", name, l)
		}

		ok, errText, val := encodeJSON(name, frame, arguments[0])
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

// decode(opaque:bytearray) -> list(bool, string, <VALUE>)
func getStdJSONdecode(name string) stdFuncType {
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

		ok, errText, val := decodeJSON(name, frame, byteArray.data)
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

func encodeJSON(name string, frame *funl.Frame, inValue funl.Value) (ok bool, errText string, val funl.Value) {
	defer func() {
		if r := recover(); r != nil {
			err, _ := r.(error)
			ok, errText = false, fmt.Sprintf("%s: %v", name, err)
		}
	}()

	resultAsBytes := traverseValuesEncode(frame, inValue, make([]byte, 0))
	val = funl.Value{Kind: funl.OpaqueValue, Data: &OpaqueByteArray{data: resultAsBytes}}
	ok = true
	return
}

func decodeJSON(name string, frame *funl.Frame, indata []byte) (ok bool, errText string, val funl.Value) {
	defer func() {
		if r := recover(); r != nil {
			err, _ := r.(error)
			ok, errText = false, fmt.Sprintf("%s: %v", name, err)
		}
	}()

	var res interface{}
	err := json.Unmarshal(indata, &res)
	if err != nil {
		ok, errText = false, fmt.Sprintf("%s: error in unmarshal: %v", name, err)
		return
	}
	val = traverseValues(frame, res)
	ok = true
	return
}

func traverseValuesEncode(frame *funl.Frame, inValue funl.Value, prevsl []byte) (nextsl []byte) {
	if inValue.Kind == funl.OpaqueValue {
		_, isNullV := inValue.Data.(*OpaqueJSONnull)
		if isNullV {
			nextsl = append(prevsl, []byte("null")...)
			return
		}
	}
	switch inValue.Kind {

	case funl.BoolValue:
		if inValue.Data.(bool) {
			nextsl = append(prevsl, []byte("true")...)
			return
		}
		nextsl = append(prevsl, []byte("false")...)
		return

	case funl.StringValue:
		strVal := inValue.Data.(string)
		strAsBytes, err := json.Marshal(strVal)
		if err != nil {
			panic(err)
		}
		nextsl = append(prevsl, strAsBytes...)
		return

	case funl.IntValue:
		intVal := inValue.Data.(int)
		intAsBytes, err := json.Marshal(intVal)
		if err != nil {
			panic(err)
		}
		nextsl = append(prevsl, intAsBytes...)
		return

	case funl.FloatValue:
		floatVal := inValue.Data.(float64)
		floatAsBytes, err := json.Marshal(floatVal)
		if err != nil {
			panic(err)
		}
		nextsl = append(prevsl, floatAsBytes...)
		return

	case funl.ListValue:
		var listAsBytes []byte
		listAsBytes = append(listAsBytes, []byte("[")...)
		listIter := funl.NewListIterator(inValue)
		isFirstRound := true
		for {
			nextItem := listIter.Next()
			if nextItem == nil {
				break
			}
			if !isFirstRound {
				listAsBytes = append(listAsBytes, []byte(", ")...)
			}
			isFirstRound = false
			listAsBytes = traverseValuesEncode(frame, *nextItem, listAsBytes)
		}
		listAsBytes = append(listAsBytes, []byte("]")...)
		nextsl = append(prevsl, listAsBytes...)
		return

	case funl.MapValue:
		var mapAsBytes []byte
		mapAsBytes = append(mapAsBytes, []byte("{")...)
		keyvals := funl.HandleKeyvalsOP(frame, []*funl.Item{&funl.Item{Type: funl.ValueItem, Data: inValue}})
		kvListIter := funl.NewListIterator(keyvals)
		isFirstRound := true
		for {
			nextKV := kvListIter.Next()
			if nextKV == nil {
				break
			}
			if !isFirstRound {
				mapAsBytes = append(mapAsBytes, []byte(", ")...)
			}
			isFirstRound = false
			kvIter := funl.NewListIterator(*nextKV)
			keyv := *(kvIter.Next())
			valv := *(kvIter.Next())
			if keyv.Kind != funl.StringValue {
				panic(fmt.Errorf("JSON object key not a string: %v", keyv))
			}
			mapAsBytes = traverseValuesEncode(frame, keyv, mapAsBytes)
			mapAsBytes = append(mapAsBytes, []byte(": ")...)
			mapAsBytes = traverseValuesEncode(frame, valv, mapAsBytes)
		}
		mapAsBytes = append(mapAsBytes, []byte("}")...)
		nextsl = append(prevsl, mapAsBytes...)
		return
	}
	panic(fmt.Errorf("Unexpected type: %v", inValue))
}

func traverseValues(frame *funl.Frame, intf interface{}) funl.Value {
	val := reflect.ValueOf(intf)
	if intf == nil {
		return funl.Value{Kind: funl.OpaqueValue, Data: &OpaqueJSONnull{}}
	}
	switch val.Kind() {
	case reflect.Bool:
		return funl.Value{Kind: funl.BoolValue, Data: val.Bool()}
	case reflect.String:
		return funl.Value{Kind: funl.StringValue, Data: val.String()}
	case reflect.Int:
		return funl.Value{Kind: funl.IntValue, Data: val.Int()}
	case reflect.Float64:
		return funl.Value{Kind: funl.FloatValue, Data: val.Float()}
	case reflect.Slice:
		var values []funl.Value
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i)
			listItem := traverseValues(frame, item.Interface())
			values = append(values, listItem)
		}
		return funl.MakeListOfValues(frame, values)
	case reflect.Map:
		mapval := funl.HandleMapOP(frame, []*funl.Item{})
		for _, k := range val.MapKeys() {
			v := val.MapIndex(k)
			keyv := funl.Value{Kind: funl.StringValue, Data: k.String()}
			valv := traverseValues(frame, v.Interface())
			putArgs := []*funl.Item{
				&funl.Item{Type: funl.ValueItem, Data: mapval},
				&funl.Item{Type: funl.ValueItem, Data: keyv},
				&funl.Item{Type: funl.ValueItem, Data: valv},
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
