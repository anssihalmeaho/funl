package std

import (
	"github.com/anssihalmeaho/funl/funl"
)

type stdFuncType func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value)

type stdFuncInfo struct {
	Name       string
	Getter     func(name string) stdFuncType
	IsFunction bool
}

func setSTDFunctions(topFrame *funl.Frame, stdModuleName string, stdFuncs []stdFuncInfo, interpreter *funl.Interpreter) (err error) {
	nsSid := funl.SymIDMap.Add(stdModuleName)
	interpreter.NsDir.Put(nsSid, topFrame)

	for _, v := range stdFuncs {
		extProc := funl.ExtProcType{
			Impl:       v.Getter(stdModuleName + ":" + v.Name),
			IsFunction: v.IsFunction,
		}
		epVal := funl.Value{Kind: funl.ExtProcValue, Data: extProc}
		item := &funl.Item{Type: funl.ValueItem, Data: epVal}
		err = topFrame.Syms.Add(v.Name, item)
		if err != nil {
			return
		}
	}
	return
}

// StdFuncType exposed
type StdFuncType func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value)

// StdFuncInfo exposed
type StdFuncInfo struct {
	Name       string
	Getter     func(name string) StdFuncType
	IsFunction bool
}

// SetSTDFunctions exposed
func SetSTDFunctions(topFrame *funl.Frame, stdModuleName string, stdFuncs []StdFuncInfo, interpreter *funl.Interpreter) (err error) {
	nsSid := funl.SymIDMap.Add(stdModuleName)
	interpreter.NsDir.Put(nsSid, topFrame)

	for _, v := range stdFuncs {
		extProc := funl.ExtProcType{
			Impl:       v.Getter(stdModuleName + ":" + v.Name),
			IsFunction: v.IsFunction,
		}
		epVal := funl.Value{Kind: funl.ExtProcValue, Data: extProc}
		item := &funl.Item{Type: funl.ValueItem, Data: epVal}
		err = topFrame.Syms.Add(v.Name, item)
		if err != nil {
			return
		}
	}
	return
}

//InitSTD is used for initializing standard library
func InitSTD(interpreter *funl.Interpreter) (err error) {
	inits := []func(*funl.Interpreter) error{
		initSTDIO,
		initSTDTIME,
		initSTDBytes,
		initSTDFiles,
		initSTDJson,
		initSTDHttp,
		initSTDos,
		initSTDlog,
		initSTDStr,
		initSTDMath,
		initSTDAst,
		initSTDRPC,
		initSTDbase64,
		initSTDVar,
		initSTDRun,
	}
	for _, initf := range inits {
		err = initf(interpreter)
		if err != nil {
			return
		}
	}
	return
}
