package std

import (
	"fmt"
	"sync"

	"github.com/anssihalmeaho/funl/funl"
)

func initSTDVar(interpreter *funl.Interpreter) (err error) {
	stdModuleName := "stdvar"
	topFrame := &funl.Frame{
		Syms:     funl.NewSymt(),
		OtherNS:  make(map[funl.SymID]funl.ImportInfo),
		Imported: make(map[funl.SymID]*funl.Frame),
	}
	stdFuncs := []stdFuncInfo{
		{
			Name:   "new",
			Getter: getStdVarNew,
		},
		{
			Name:   "value",
			Getter: getStdVarValue,
		},
		{
			Name:   "set",
			Getter: getStdVarSet,
		},
		{
			Name:   "change",
			Getter: getStdVarChange,
		},
		{
			Name:   "change-v2",
			Getter: getStdVarChangeV2,
		},
	}
	err = setSTDFunctions(topFrame, stdModuleName, stdFuncs, interpreter)
	return
}

// OpaqueVarRef ...
type OpaqueVarRef struct {
	ValRef *funl.Value
	sync.RWMutex
}

// TypeName ...
func (ref *OpaqueVarRef) TypeName() string {
	return "var-ref"
}

// Str ...
func (ref *OpaqueVarRef) Str() string {
	ref.RLock()
	val := *(ref.ValRef)
	ref.RUnlock()
	return fmt.Sprintf("var-ref(%v)", val)
}

// Equals ...
func (ref *OpaqueVarRef) Equals(with funl.OpaqueAPI) bool {
	other, ok := with.(*OpaqueVarRef)
	if !ok {
		return false
	}
	other.RLock()
	ref.RLock()
	isSame := (ref.ValRef == other.ValRef)
	ref.RUnlock()
	other.RUnlock()
	return isSame
}

func getStdVarChangeV2(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l < 2 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need two", name, l)
		}
		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: assuming opaque var-ref", name)
		}
		varref, convOK := arguments[0].Data.(*OpaqueVarRef)
		if !convOK {
			funl.RunTimeError2(frame, "%s: assuming var-ref", name)
		}

		switch arguments[1].Kind {
		case funl.FunctionValue:
			funcVal, funcOK := arguments[1].Data.(funl.FuncValue)
			if !funcOK {
				funl.RunTimeError2(frame, "%s: not func as 2nd argument", name)
			}
			if funcVal.FuncProto.IsProc {
				funl.RunTimeError2(frame, "%s: proc not allowed", name)
			}
		case funl.ExtProcValue:
			extFuncVal, extfuncOK := arguments[1].Data.(funl.ExtProcType)
			if !extfuncOK {
				funl.RunTimeError2(frame, "%s: not ext-func as 2nd argument", name)
			}
			if !extFuncVal.IsFunction {
				funl.RunTimeError2(frame, "%s: ext-proc not allowed", name)
			}
		default:
			funl.RunTimeError2(frame, "%s: assuming function as argument", name)
		}

		varref.Lock()
		defer varref.Unlock()

		oldVal := *(varref.ValRef)

		argsForCall := []*funl.Item{
			{
				Type: funl.ValueItem,
				Data: arguments[1],
			},
			{
				Type: funl.ValueItem,
				Data: oldVal,
			},
		}
		if len(arguments) == 3 {
			argsForCall = append(argsForCall, &funl.Item{Type: funl.ValueItem, Data: arguments[2]})
		}
		retListVal, callErr := func() (rv funl.Value, errDesc error) {
			defer func() {
				if r := recover(); r != nil {
					var rtestr string
					if err, isError := r.(error); isError {
						rtestr = err.Error()
					}
					errDesc = fmt.Errorf("%s", rtestr)
				}
			}()
			return funl.HandleCallOP(frame, argsForCall), nil
		}()

		var rval funl.Value
		var addRetVal funl.Value
		var errtext string
		if callErr == nil && retListVal.Kind != funl.ListValue {
			callErr = fmt.Errorf("List value expected")
		}
		if callErr == nil {
			lit := funl.NewListIterator(retListVal)
			newVal := lit.Next()
			if newVal == nil {
				callErr = fmt.Errorf("Too short list received (empty)")
				errtext = callErr.Error()
				rval = funl.Value{Kind: funl.StringValue, Data: ""}
				addRetVal = funl.Value{Kind: funl.StringValue, Data: ""}
			} else {
				nextVal := lit.Next()
				if nextVal == nil {
					callErr = fmt.Errorf("Too short list received (one item)")
					errtext = callErr.Error()
					rval = funl.Value{Kind: funl.StringValue, Data: ""}
					addRetVal = funl.Value{Kind: funl.StringValue, Data: ""}
				} else {
					addRetVal = *nextVal
					rval = *newVal
					varref.ValRef = newVal
				}
			}
		} else {
			errtext = callErr.Error()
			rval = funl.Value{Kind: funl.StringValue, Data: ""}
			addRetVal = funl.Value{Kind: funl.StringValue, Data: ""}
		}

		values := []funl.Value{
			{
				Kind: funl.BoolValue,
				Data: callErr == nil,
			},
			{
				Kind: funl.StringValue,
				Data: errtext,
			},
			rval,
			addRetVal,
		}
		retVal = funl.MakeListOfValues(frame, values)
		return
	}
}

func getStdVarChange(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 2 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need two", name, l)
		}
		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: assuming opaque var-ref", name)
		}
		varref, convOK := arguments[0].Data.(*OpaqueVarRef)
		if !convOK {
			funl.RunTimeError2(frame, "%s: assuming var-ref", name)
		}

		switch arguments[1].Kind {
		case funl.FunctionValue:
			funcVal, funcOK := arguments[1].Data.(funl.FuncValue)
			if !funcOK {
				funl.RunTimeError2(frame, "%s: not func as 2nd argument", name)
			}
			if funcVal.FuncProto.IsProc {
				funl.RunTimeError2(frame, "%s: proc not allowed", name)
			}
		case funl.ExtProcValue:
			extFuncVal, extfuncOK := arguments[1].Data.(funl.ExtProcType)
			if !extfuncOK {
				funl.RunTimeError2(frame, "%s: not ext-func as 2nd argument", name)
			}
			if !extFuncVal.IsFunction {
				funl.RunTimeError2(frame, "%s: ext-proc not allowed", name)
			}
		default:
			funl.RunTimeError2(frame, "%s: assuming function as argument", name)
		}

		varref.Lock()
		defer varref.Unlock()

		oldVal := *(varref.ValRef)

		argsForCall := []*funl.Item{
			{
				Type: funl.ValueItem,
				Data: arguments[1],
			},
			{
				Type: funl.ValueItem,
				Data: oldVal,
			},
		}
		newVal, callErr := func() (rv funl.Value, errDesc error) {
			defer func() {
				if r := recover(); r != nil {
					var rtestr string
					if err, isError := r.(error); isError {
						rtestr = err.Error()
					}
					errDesc = fmt.Errorf("%s", rtestr)
				}
			}()
			return funl.HandleCallOP(frame, argsForCall), nil
		}()
		var rval funl.Value
		var errtext string
		if callErr == nil {
			rval = newVal
			varref.ValRef = &newVal
		} else {
			errtext = callErr.Error()
			rval = funl.Value{Kind: funl.StringValue, Data: ""}
		}

		values := []funl.Value{
			{
				Kind: funl.BoolValue,
				Data: callErr == nil,
			},
			{
				Kind: funl.StringValue,
				Data: errtext,
			},
			rval,
		}
		retVal = funl.MakeListOfValues(frame, values)
		return
	}
}

func getStdVarSet(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 2 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need two", name, l)
		}
		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: assuming opaque var-ref", name)
		}
		varref, convOK := arguments[0].Data.(*OpaqueVarRef)
		if !convOK {
			funl.RunTimeError2(frame, "%s: assuming var-ref", name)
		}
		newval := arguments[1]
		varref.Lock()
		varref.ValRef = &newval
		varref.Unlock()
		retVal = funl.Value{Kind: funl.BoolValue, Data: true}
		return
	}
}

func getStdVarValue(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need one", name, l)
		}
		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: assuming opaque var-ref", name)
		}
		varref, convOK := arguments[0].Data.(*OpaqueVarRef)
		if !convOK {
			funl.RunTimeError2(frame, "%s: assuming var-ref", name)
		}
		varref.RLock()
		retVal = *(varref.ValRef)
		varref.RUnlock()
		return
	}
}

func getStdVarNew(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need one", name, l)
		}
		val := arguments[0]
		varref := &OpaqueVarRef{ValRef: &val}
		retVal = funl.Value{Kind: funl.OpaqueValue, Data: varref}
		return
	}
}
