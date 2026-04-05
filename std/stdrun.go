package std

import (
	"os"

	"github.com/anssihalmeaho/funl/funl"
)

func initSTDRun(interpreter *funl.Interpreter) (err error) {
	stdModuleName := "stdrun"
	topFrame := funl.NewTopFrameWithInterpreter(interpreter)
	stdAstFuncs := []stdFuncInfo{
		{
			Name:   "add-to-mod-cache",
			Getter: getAddToModCache,
			//IsFunction: true,
		},
		{
			Name:       "backtrace",
			Getter:     getBacktrace,
			IsFunction: true,
		},
		{
			Name:   "run-sub-package",
			Getter: getRunSubPackage,
		},
		{
			Name:   "run-sub-pack-raw",
			Getter: getRunSubPackRaw,
		},
		{
			Name:   "run-sub",
			Getter: getRunSub,
		},
	}
	err = setSTDFunctions(topFrame, stdModuleName, stdAstFuncs, interpreter)
	return
}

func getRunSub(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		l := len(arguments)
		if l != 1 && l != 2 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need one or two", name, l)
		}
		if arguments[0].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: assuming string", name)
		}
		srcFileName := arguments[0].Data.(string)
		if l == 2 && arguments[1].Kind != funl.ListValue {
			funl.RunTimeError2(frame, "%s: assuming list as 2nd argument", name)
		}
		args := []*funl.Item{}
		if l == 2 {
			lit := funl.NewListIterator(arguments[1])
			for {
				nextv := lit.Next()
				if nextv == nil {
					break
				}
				args = append(args, &funl.Item{Type: funl.ValueItem, Data: *nextv})
			}
		}

		content, err := os.ReadFile(srcFileName)
		if err != nil {
			return makeTripletList(frame, false, err.Error(), funl.Value{Kind: funl.StringValue, Data: ""})
		}
		val, err := funl.FunlMainWithArgs(string(content), args, "main", srcFileName, InitSTD)
		if err != nil {
			return makeTripletList(frame, false, err.Error(), val)
		}
		return makeTripletList(frame, true, "", val)
	}
}

func getRunSubPackRaw(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		l := len(arguments)
		if l != 2 && l != 3 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need two or three", name, l)
		}
		if arguments[0].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: assuming string", name)
		}
		if arguments[1].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: assuming opaque value", name)
		}
		if l == 3 && arguments[2].Kind != funl.ListValue {
			funl.RunTimeError2(frame, "%s: assuming list as 2nd argument", name)
		}
		args := []*funl.Item{}
		if l == 3 {
			lit := funl.NewListIterator(arguments[2])
			for {
				nextv := lit.Next()
				if nextv == nil {
					break
				}
				args = append(args, &funl.Item{Type: funl.ValueItem, Data: *nextv})
			}
		}
		srcFileName := arguments[0].Data.(string)
		rawDataBytes := arguments[1].Data.(*OpaqueByteArray)
		val, err := funl.FunlMainWithPackageContent(rawDataBytes.data, args, "main", srcFileName, InitSTD)
		if err != nil {
			return makeTripletList(frame, false, err.Error(), val)
		}
		return makeTripletList(frame, true, "", val)
	}
}

func getRunSubPackage(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		l := len(arguments)
		if l != 1 && l != 2 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need one or two", name, l)
		}
		if arguments[0].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: assuming string", name)
		}
		if l == 2 && arguments[1].Kind != funl.ListValue {
			funl.RunTimeError2(frame, "%s: assuming list as 2nd argument", name)
		}
		args := []*funl.Item{}
		if l == 2 {
			lit := funl.NewListIterator(arguments[1])
			for {
				nextv := lit.Next()
				if nextv == nil {
					break
				}
				args = append(args, &funl.Item{Type: funl.ValueItem, Data: *nextv})
			}
		}

		importName := arguments[0].Data.(string)
		val, err := funl.FunlMainWithPackage(args, "main", importName, InitSTD)
		if err != nil {
			return makeTripletList(frame, false, err.Error(), val)
		}
		return makeTripletList(frame, true, "", val)
	}
}

func getBacktrace(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		stack := []funl.Value{}

		// if debug printing disabled in pure functions return just empty list
		if funl.PrintingDisabledInFunctions {
			retVal = funl.MakeListOfValues(frame, stack)
			return
		}

		prevFrame := frame
		for {
			if prevFrame == nil {
				break
			}
			mapv := funl.HandleMapOP(frame, []*funl.Item{})
			if prevFrame.FuncProto == nil {
				mapv = putToMap(frame, mapv, "line", funl.Value{Kind: funl.IntValue, Data: 0})
				mapv = putToMap(frame, mapv, "file", funl.Value{Kind: funl.StringValue, Data: "-"})
				mapv = putToMap(frame, mapv, "args", funl.MakeListOfValues(frame, []funl.Value{}))
			} else {
				mapv = putToMap(frame, mapv, "line", funl.Value{Kind: funl.IntValue, Data: prevFrame.FuncProto.Lineno})
				mapv = putToMap(frame, mapv, "file", funl.Value{Kind: funl.StringValue, Data: prevFrame.FuncProto.SrcFileName})
				mapv = putToMap(frame, mapv, "args", funl.MakeListOfValues(frame, prevFrame.EvaluatedArgs))
			}

			stack = append(stack, mapv)
			prevFrame = prevFrame.Previous
		}
		retVal = funl.MakeListOfValues(frame, stack)
		return
	}
}

func getAddToModCache(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 2 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need two", name, l)
		}
		if arguments[0].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: assuming string", name)
		}
		if arguments[1].Kind != funl.MapValue {
			funl.RunTimeError2(frame, "%s: assuming map", name)
		}

		importModName := arguments[0].Data.(string)
		nspace := &funl.NSpace{OtherNS: make(map[funl.SymID]funl.ImportInfo), Syms: funl.NewSymt()}

		// loop symbol to value mappings
		keyvals := funl.HandleKeyvalsOP(frame, []*funl.Item{{Type: funl.ValueItem, Data: arguments[1]}})
		kvListIter := funl.NewListIterator(keyvals)
		for {
			nextKV := kvListIter.Next()
			if nextKV == nil {
				break
			}
			kvIter := funl.NewListIterator(*nextKV)
			keyv := *(kvIter.Next())
			valv := *(kvIter.Next())
			if keyv.Kind != funl.StringValue {
				continue // just skip...
			}
			err := nspace.Syms.Add(keyv.Data.(string), &funl.Item{Type: funl.ValueItem, Data: valv})
			if err != nil {
				funl.RunTimeError2(frame, "%s: error in adding symbol: %v", name, err)
			}
		}

		// then add to module cache
		interpreter := frame.GetTopFrame().Interpreter
		funl.AddNStoCache(true, importModName, nspace, interpreter)
		retVal = funl.Value{Kind: funl.BoolValue, Data: true}
		return
	}
}
