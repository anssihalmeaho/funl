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

func setSTDFunctions(topFrame *funl.Frame, stdModuleName string, stdFuncs []stdFuncInfo) (err error) {
	nsSid := funl.SymIDMap.Add(stdModuleName)
	funl.GetNSDir().Put(nsSid, topFrame)

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
func InitSTD() (err error) {
	inits := []func() error{
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
	}
	for _, initf := range inits {
		err = initf()
		if err != nil {
			return
		}
	}
	return
}
