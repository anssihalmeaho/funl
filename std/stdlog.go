package std

import (
	"fmt"
	"log"
	"os"

	"github.com/anssihalmeaho/funl/funl"
)

func initSTDlog(interpreter *funl.Interpreter) (err error) {
	stdModuleName := "stdlog"
	topFrame := &funl.Frame{
		Syms:     funl.NewSymt(),
		OtherNS:  make(map[funl.SymID]funl.ImportInfo),
		Imported: make(map[funl.SymID]*funl.Frame),
	}
	stdLogFuncs := []stdFuncInfo{
		{
			Name:   "get-logger",
			Getter: getStdLogGetLogger,
		},
		{
			Name:   "get-default-logger",
			Getter: getStdLogGetDefaultLogger,
		},
	}
	err = setSTDFunctions(topFrame, stdModuleName, stdLogFuncs, interpreter)
	return
}

func getStdLogGetDefaultLogger(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		l := len(arguments)
		if (l != 0) && (l != 1) {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), needs one or two", name, l)
		}

		var prefix string
		var flag int
		var separator = ":"
		// lets check config map if there is such
		if l == 1 {
			if arguments[0].Kind != funl.MapValue {
				funl.RunTimeError2(frame, "%s: assuming map as 1st argument", name)
			}

			keyvals := funl.HandleKeyvalsOP(frame, []*funl.Item{&funl.Item{Type: funl.ValueItem, Data: arguments[0]}})
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
					funl.RunTimeError2(frame, "%s:config key not a string: %v", name, keyv)
				}
				switch keyStr := keyv.Data.(string); keyStr {
				case "prefix":
					if valv.Kind != funl.StringValue {
						funl.RunTimeError2(frame, "%s: %s value not string: %v", name, keyStr, keyv)
					}
					prefix = valv.Data.(string)
				case "separator":
					if valv.Kind != funl.StringValue {
						funl.RunTimeError2(frame, "%s: %s value not string: %v", name, keyStr, keyv)
					}
					separator = valv.Data.(string)
				case "date":
					if valv.Kind != funl.BoolValue {
						funl.RunTimeError2(frame, "%s: %s value not bool value: %v", name, keyStr, keyv)
					}
					if valv.Data.(bool) {
						flag |= log.Ldate
					}
				case "time":
					if valv.Kind != funl.BoolValue {
						funl.RunTimeError2(frame, "%s: %s value not bool value: %v", name, keyStr, keyv)
					}
					if valv.Data.(bool) {
						flag |= log.Ltime
					}
				case "microseconds":
					if valv.Kind != funl.BoolValue {
						funl.RunTimeError2(frame, "%s: %s value not bool value: %v", name, keyStr, keyv)
					}
					if valv.Data.(bool) {
						flag |= log.Lmicroseconds
					}
				case "UTC":
					if valv.Kind != funl.BoolValue {
						funl.RunTimeError2(frame, "%s: %s value not bool value: %v", name, keyStr, keyv)
					}
					if valv.Data.(bool) {
						flag |= log.LUTC
					}
				}
			}
		}

		logger := log.New(os.Stdout, prefix, flag)

		logWrapper := func(wFrame *funl.Frame, wArguments []funl.Value) funl.Value {
			if l := len(wArguments); l == 0 {
				funl.RunTimeError2(frame, "Logger assumes at least one argument")
			}
			var textOutput string
			for _, argval := range wArguments {
				textOutput += fmt.Sprintf("%s%s", separator, argval)
			}
			logger.Println(textOutput)
			return funl.Value{Kind: funl.BoolValue, Data: true}
		}

		extProc := funl.ExtProcType{
			Impl:       logWrapper,
			IsFunction: false,
		}
		retVal = funl.Value{Kind: funl.ExtProcValue, Data: extProc}
		return
	}
}

//get-logger (<proc/ext-proc>) -> ext-proc (stdlog handler which serailizes and calls handler)
func getStdLogGetLogger(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		l := len(arguments)
		if (l != 1) && (l != 2) {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), needs one or two", name, l)
		}

		switch arguments[0].Kind {
		case funl.FunctionValue, funl.ExtProcValue:
		default:
			funl.RunTimeError2(frame, "%s: assuming procedure as argument", name)
		}

		const defaultLogBufferSize = 1024
		bufSize := defaultLogBufferSize
		// lets check config map if there is such
		if l == 2 {
			if arguments[1].Kind != funl.MapValue {
				funl.RunTimeError2(frame, "%s: assuming map as 2nd argument", name)
			}

			keyvals := funl.HandleKeyvalsOP(frame, []*funl.Item{&funl.Item{Type: funl.ValueItem, Data: arguments[1]}})
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
					funl.RunTimeError2(frame, "%s:config key not a string: %v", name, keyv)
				}
				switch keyStr := keyv.Data.(string); keyStr {
				case "buffer-size":
					if valv.Kind != funl.IntValue {
						funl.RunTimeError2(frame, "%s: %s value not int: %v", name, keyStr, keyv)
					}
					bufSize = valv.Data.(int)
				}
			}
		}

		logCh := make(chan []*funl.Item, bufSize)

		go func() {
			for argsForCall := range logCh {
				funl.HandleCallOP(frame, argsForCall)
			}
		}()

		logWrapper := func(wFrame *funl.Frame, wArguments []funl.Value) funl.Value {
			if l := len(wArguments); l == 0 {
				funl.RunTimeError2(frame, "Logger assumes at least one argument")
			}
			argsForCall := []*funl.Item{
				&funl.Item{
					Type: funl.ValueItem,
					Data: arguments[0],
				},
				&funl.Item{
					Type: funl.ValueItem,
					Data: funl.MakeListOfValues(wFrame, wArguments),
				},
			}
			logCh <- argsForCall
			return funl.Value{Kind: funl.BoolValue, Data: true}
		}

		extProc := funl.ExtProcType{
			Impl:       logWrapper,
			IsFunction: false,
		}
		retVal = funl.Value{Kind: funl.ExtProcValue, Data: extProc}
		return
	}
}
