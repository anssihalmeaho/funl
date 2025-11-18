package std

import (
	"bytes"
	"encoding/csv"

	"github.com/anssihalmeaho/funl/funl"
)

func initSTDCsv(interpreter *funl.Interpreter) (err error) {
	stdModuleName := "stdcsv"
	topFrame := funl.NewTopFrameWithInterpreter(interpreter)
	stdCsvFuncs := []stdFuncInfo{
		{
			Name:       "read-all",
			Getter:     getCSVReadAll,
			IsFunction: true,
		},
		{
			Name:       "write-all",
			Getter:     getCSVWriteAll,
			IsFunction: true,
		},
	}
	err = setSTDFunctions(topFrame, stdModuleName, stdCsvFuncs, interpreter)
	return
}

func getCSVWriteAll(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), needs one", name, l)
		}
		if arguments[0].Kind != funl.ListValue {
			funl.RunTimeError2(frame, "%s: requires list value", name)
		}
		buf := bytes.NewBuffer([]byte{})
		w := csv.NewWriter(buf)
		listIter := funl.NewListIterator(arguments[0])
		for {
			nextItem := listIter.Next()
			if nextItem == nil {
				break
			}
			if nextItem.Kind != funl.ListValue {
				funl.RunTimeError2(frame, "%s: not a list value", name)
			}
			strList := []string{}
			oneIter := funl.NewListIterator(*nextItem)
			for {
				nextStr := oneIter.Next()
				if nextStr == nil {
					break
				}
				if nextStr.Kind != funl.StringValue {
					funl.RunTimeError2(frame, "%s: not a string value", name)
				}
				strList = append(strList, nextStr.Data.(string))
			}
			w.Write(strList)
		}
		w.Flush()
		if err := w.Error(); err != nil {
			return funl.MakeListOfValues(frame, []funl.Value{
				{
					Kind: funl.BoolValue,
					Data: false,
				},
				{
					Kind: funl.StringValue,
					Data: err.Error(),
				},
				funl.MakeListOfValues(frame, []funl.Value{}),
			})
		}
		return funl.MakeListOfValues(frame, []funl.Value{
			{
				Kind: funl.BoolValue,
				Data: true,
			},
			{
				Kind: funl.StringValue,
				Data: "",
			},
			{
				Kind: funl.OpaqueValue,
				Data: &OpaqueByteArray{data: buf.Bytes()},
			},
		})
	}
}

func getCSVReadAll(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), needs one", name, l)
		}
		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		byteArray, ok := arguments[0].Data.(*OpaqueByteArray)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not bytearray value", name)
		}
		reader := csv.NewReader(bytes.NewReader(byteArray.data))
		records, err := reader.ReadAll()
		if err != nil {
			return funl.MakeListOfValues(frame, []funl.Value{
				{
					Kind: funl.BoolValue,
					Data: false,
				},
				{
					Kind: funl.StringValue,
					Data: err.Error(),
				},
				funl.MakeListOfValues(frame, []funl.Value{}),
			})
		}

		recordList := []funl.Value{}
		for _, record := range records {
			oneList := []funl.Value{}
			for _, oneVal := range record {
				oneList = append(oneList, funl.Value{Kind: funl.StringValue, Data: oneVal})
			}
			recordList = append(recordList, funl.MakeListOfValues(frame, oneList))
		}
		return funl.MakeListOfValues(frame, []funl.Value{
			{
				Kind: funl.BoolValue,
				Data: true,
			},
			{
				Kind: funl.StringValue,
				Data: "",
			},
			funl.MakeListOfValues(frame, recordList),
		})
	}
}
