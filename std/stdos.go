package std

import (
	"bytes"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/anssihalmeaho/funl/funl"
)

func initSTDos(interpreter *funl.Interpreter) (err error) {
	stdModuleName := "stdos"
	topFrame := funl.NewTopFrameWithInterpreter(interpreter)
	stdOSFuncs := []stdFuncInfo{
		{
			Name:   "reg-signal-handler",
			Getter: getStdOSregSignalHandler,
		},
		{
			Name:   "exit",
			Getter: getStdOSexit,
		},
		{
			Name:   "getenv",
			Getter: getStdOSGetEnv,
		},
		{
			Name:   "setenv",
			Getter: getStdOSSetEnv,
		},
		{
			Name:   "unsetenv",
			Getter: getStdOSUnSetEnv,
		},
		{
			Name:   "exec",
			Getter: getStdOSExec,
		},
		{
			Name:   "exec-with",
			Getter: getStdOSExecWith,
		},
	}
	err = setSTDFunctions(topFrame, stdModuleName, stdOSFuncs, interpreter)
	return
}

func getStdOSExecWith(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		l := len(arguments)
		if l < 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), needs at least one", name, l)
		}

		var stdinBytearr *OpaqueByteArray
		if arguments[0].Kind != funl.MapValue {
			funl.RunTimeError2(frame, "%s: requires map value", name)
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
				funl.RunTimeError2(frame, "%s: info key not a string: %v", name, keyv)
			}
			switch keyStr := keyv.Data.(string); keyStr {
			case "stdin":
				if valv.Kind != funl.OpaqueValue {
					funl.RunTimeError2(frame, "%s: %s value not opaque value: %v", name, keyStr, keyv)
				}
				var convOK bool
				stdinBytearr, convOK = valv.Data.(*OpaqueByteArray)
				if !convOK {
					funl.RunTimeError2(frame, "%s: %s value not opaque: %v", name, keyStr, keyv)
				}
			}
		}

		var argStrs []string
		for i, argval := range arguments[1:] {
			if argval.Kind != funl.StringValue {
				funl.RunTimeError2(frame, "%s: assuming string as key (%d)", name, i)
			}
			argStrs = append(argStrs, argval.Data.(string))
		}

		var outStd bytes.Buffer
		var outErr bytes.Buffer
		cmd := exec.Command(argStrs[0], argStrs[1:]...)
		cmd.Stdout = &outStd
		cmd.Stderr = &outErr

		if stdinBytearr != nil {
			cmd.Stdin = bytes.NewBuffer(stdinBytearr.data)
		}

		err := cmd.Run()

		var errText string
		if err != nil {
			errText = err.Error()
		}
		values := []funl.Value{
			{
				Kind: funl.BoolValue,
				Data: err == nil,
			},
			{
				Kind: funl.StringValue,
				Data: errText,
			},
			{
				Kind: funl.OpaqueValue,
				Data: &OpaqueByteArray{data: outStd.Bytes()},
			},
			{
				Kind: funl.OpaqueValue,
				Data: &OpaqueByteArray{data: outErr.Bytes()},
			},
		}
		retVal = funl.MakeListOfValues(frame, values)
		return
	}
}

func getStdOSExec(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		l := len(arguments)
		if l < 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), needs at least one", name, l)
		}
		var argStrs []string
		for i, argval := range arguments {
			if argval.Kind != funl.StringValue {
				funl.RunTimeError2(frame, "%s: assuming string as key (%d)", name, i)
			}
			argStrs = append(argStrs, argval.Data.(string))
		}

		var outStd bytes.Buffer
		var outErr bytes.Buffer
		cmd := exec.Command(argStrs[0], argStrs[1:]...)
		cmd.Stdout = &outStd
		cmd.Stderr = &outErr
		err := cmd.Run()

		var errText string
		if err != nil {
			errText = err.Error()
		}
		values := []funl.Value{
			{
				Kind: funl.BoolValue,
				Data: err == nil,
			},
			{
				Kind: funl.StringValue,
				Data: errText,
			},
			{
				Kind: funl.OpaqueValue,
				Data: &OpaqueByteArray{data: outStd.Bytes()},
			},
			{
				Kind: funl.OpaqueValue,
				Data: &OpaqueByteArray{data: outErr.Bytes()},
			},
		}
		retVal = funl.MakeListOfValues(frame, values)
		return
	}
}

func getStdOSUnSetEnv(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		l := len(arguments)
		if l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), needs one", name, l)
		}

		if arguments[0].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: assuming string as key", name)
		}
		err := os.Unsetenv(arguments[0].Data.(string))
		var errText string
		if err != nil {
			errText = err.Error()
		}
		retVal = funl.Value{Kind: funl.StringValue, Data: errText}
		return
	}
}

func getStdOSSetEnv(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		l := len(arguments)
		if l != 2 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), needs two", name, l)
		}

		if arguments[0].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: assuming string as key", name)
		}
		if arguments[1].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: assuming string as value", name)
		}
		err := os.Setenv(arguments[0].Data.(string), arguments[1].Data.(string))
		var errText string
		if err != nil {
			errText = err.Error()
		}
		retVal = funl.Value{Kind: funl.StringValue, Data: errText}
		return
	}
}

func getStdOSGetEnv(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		l := len(arguments)
		if l > 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need at most one", name, l)
		}

		// if no arguments are given then returns map with all key-values
		if l == 0 {
			envMapval := funl.HandleMapOP(frame, []*funl.Item{})
			for _, e := range os.Environ() {
				pair := strings.Split(e, "=")
				if len(pair) != 2 {
					continue
				}

				keyv := funl.Value{Kind: funl.StringValue, Data: pair[0]}
				valv := funl.Value{Kind: funl.StringValue, Data: pair[1]}
				putArgs := []*funl.Item{
					&funl.Item{Type: funl.ValueItem, Data: envMapval},
					&funl.Item{Type: funl.ValueItem, Data: keyv},
					&funl.Item{Type: funl.ValueItem, Data: valv},
				}
				envMapval = funl.HandlePutOP(frame, putArgs)
				if envMapval.Kind != funl.MapValue {
					funl.RunTimeError2(frame, "%s: failed to put key-value to map", name)
				}

			}
			retVal = envMapval
			return
		}
		// otherwise, lets find value for key given as argument
		if arguments[0].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: assuming string as argument", name)
		}
		envVal, found := os.LookupEnv(arguments[0].Data.(string))
		retValues := []funl.Value{
			{
				Kind: funl.BoolValue,
				Data: found,
			},
			{
				Kind: funl.StringValue,
				Data: envVal,
			},
		}
		retVal = funl.MakeListOfValues(frame, retValues)
		return
	}
}

// call(stdos.reg-signal-handler, proc(<signal-str>) end, <int>, <int>,...)
// -> ext-proc : canceller
func getStdOSregSignalHandler(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		l := len(arguments)
		if l < 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need at least one", name, l)
		}
		if arguments[0].Kind != funl.FunctionValue {
			// what about ext-procs...?
			funl.RunTimeError2(frame, "%s: requires func/proc value", name)
		}
		var signums []os.Signal
		if l > 1 {
			for _, signalNumVal := range arguments[1:] {
				if signalNumVal.Kind != funl.IntValue {
					funl.RunTimeError2(frame, "%s: requires int value for signal number", name)
				}
				signalNumAsInt, ok := signalNumVal.Data.(int)
				if !ok {
					funl.RunTimeError2(frame, "%s: signal number could not be read", name)
				}
				signalNum := syscall.Signal(signalNumAsInt)
				signums = append(signums, signalNum)
			}
		}
		c := make(chan os.Signal, 1)
		if l > 1 {
			signal.Notify(c, signums...)
		} else {
			signal.Notify(c)
		}
		go func() {
			for {
				sig := <-c
				sigNum, _ := sig.(syscall.Signal)

				argsForCall := []*funl.Item{
					&funl.Item{
						Type: funl.ValueItem,
						Data: arguments[0],
					},
					&funl.Item{
						Type: funl.ValueItem,
						Data: funl.Value{
							Kind: funl.IntValue,
							Data: int(sigNum),
						},
					},
					&funl.Item{
						Type: funl.ValueItem,
						Data: funl.Value{
							Kind: funl.StringValue,
							Data: sig.String(),
						},
					},
				}
				funl.HandleCallOP(frame, argsForCall)
			}
		}()

		retVal = funl.Value{Kind: funl.BoolValue, Data: true} // temporary
		return
	}
}

func getStdOSexit(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		l := len(arguments)
		if l == 1 {
			if arguments[0].Kind != funl.IntValue {
				exitCode, ok := arguments[0].Data.(int)
				if ok {
					os.Exit(exitCode)
				}
			}
		}
		os.Exit(0)
		return
	}
}
