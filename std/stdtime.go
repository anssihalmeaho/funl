package std

import (
	"fmt"
	"github.com/anssihalmeaho/funl"
	"time"
)

func initSTDTIME() (err error) {
	stdModuleName := "stdtime"
	topFrame := &funl.Frame{
		Syms:     funl.NewSymt(),
		OtherNS:  make(map[funl.SymID]funl.ImportInfo),
		Imported: make(map[funl.SymID]*funl.Frame),
	}
	stdTimeFuncs := []stdFuncInfo{
		{
			Name:   "sleep",
			Getter: getStdTimeSleep,
		},
		{
			Name:   "nanosleep",
			Getter: getStdTimeNanoSleep,
		},
		{
			Name:   "newtimer",
			Getter: getStdNewTimer,
		},
		{
			Name:   "stoptimer",
			Getter: getStdStopTimer,
		},
		{
			Name:   "callafter",
			Getter: getStdCallAfter,
		},
		{
			Name:   "newticker",
			Getter: getStdNewTicker,
		},
		{
			Name:   "stopticker",
			Getter: getStdStopTicker,
		},
	}
	err = setSTDFunctions(topFrame, stdModuleName, stdTimeFuncs)
	return
}

func getStdTimeNanoSleep(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.IntValue {
			funl.RunTimeError2(frame, "%s: requires integer as input", name)
		}
		sleepTime := arguments[0].Data.(int)
		time.Sleep(time.Duration(sleepTime))
		retVal = funl.Value{Kind: funl.BoolValue, Data: true}
		return
	}
}

func getStdTimeSleep(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.IntValue {
			funl.RunTimeError2(frame, "%s: requires integer as input", name)
		}
		sleepTime := arguments[0].Data.(int)
		time.Sleep(time.Duration(sleepTime) * time.Second)
		retVal = funl.Value{Kind: funl.BoolValue, Data: true}
		return
	}
}

type opaqueTimer struct {
	t *time.Timer
}

func (ot *opaqueTimer) TypeName() string {
	return "timer"
}

func (ot *opaqueTimer) Str() string {
	return fmt.Sprintf("timer(%#v)", *ot)
}

func (ot *opaqueTimer) Equals(with funl.OpaqueAPI) bool {
	timerVal, ok := with.(*opaqueTimer)
	if !ok {
		return false
	}
	return *ot == *timerVal
}

func getStdStopTimer(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		timVal, ok := arguments[0].Data.(*opaqueTimer)
		if !ok {
			funl.RunTimeError2(frame, "%s: requires timer value", name)
		}
		retVal = funl.Value{Kind: funl.BoolValue, Data: timVal.t.Stop()}
		return
	}
}

func getStdCallAfter(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 2 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.IntValue {
			funl.RunTimeError2(frame, "%s: requires int value", name)
		}
		durationInt, ok := arguments[0].Data.(int)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not int value", name)
		}
		if arguments[1].Kind != funl.FunctionValue {
			funl.RunTimeError2(frame, "%s: requires func/proc as 2nd argument", name)
		}
		_, ok = arguments[1].Data.(funl.FuncValue)
		if !ok {
			funl.RunTimeError2(frame, "%s: not func/proc as 2nd argument", name)
		}
		wrapperFunc := func() {
			argsForCall := []*funl.Item{&funl.Item{Type: funl.ValueItem, Data: arguments[1]}}
			funl.HandleCallOP(frame, argsForCall)
		}
		timer := time.AfterFunc(time.Duration(durationInt), wrapperFunc)
		retVal = funl.Value{Kind: funl.OpaqueValue, Data: &opaqueTimer{t: timer}}
		return
	}
}

func getStdNewTimer(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 3 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.IntValue {
			funl.RunTimeError2(frame, "%s: requires int value", name)
		}
		durationInt, ok := arguments[0].Data.(int)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not int value", name)
		}
		if arguments[1].Kind != funl.ChanValue {
			funl.RunTimeError2(frame, "%s: requires channel as 2nd argument", name)
		}
		timer := time.NewTimer(time.Duration(durationInt))
		retVal = funl.Value{Kind: funl.OpaqueValue, Data: &opaqueTimer{t: timer}}
		go func() {
			<-timer.C
			ch := arguments[1].Data.(chan funl.Value)
			ch <- arguments[2]
		}()
		return
	}
}

type opaqueTicker struct {
	t    *time.Ticker
	done chan bool
}

func (ot *opaqueTicker) TypeName() string {
	return "ticker"
}

func (ot *opaqueTicker) Str() string {
	return fmt.Sprintf("ticker(%#v)", *ot)
}

func (ot *opaqueTicker) Equals(with funl.OpaqueAPI) bool {
	tickerVal, ok := with.(*opaqueTicker)
	if !ok {
		return false
	}
	return *ot == *tickerVal
}

func getStdStopTicker(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		tickVal, ok := arguments[0].Data.(*opaqueTicker)
		if !ok {
			funl.RunTimeError2(frame, "%s: requires ticker value", name)
		}
		tickVal.t.Stop()
		retVal = funl.Value{Kind: funl.BoolValue, Data: true}
		tickVal.done <- true
		return
	}
}

func getStdNewTicker(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 3 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.IntValue {
			funl.RunTimeError2(frame, "%s: requires int value", name)
		}
		durationInt, ok := arguments[0].Data.(int)
		if !ok {
			funl.RunTimeError2(frame, "%s: argument is not int value", name)
		}
		if arguments[1].Kind != funl.ChanValue {
			funl.RunTimeError2(frame, "%s: requires channel as 2nd argument", name)
		}
		ticker := time.NewTicker(time.Duration(durationInt))
		done := make(chan bool)
		retVal = funl.Value{Kind: funl.OpaqueValue, Data: &opaqueTicker{t: ticker, done: done}}
		go func() {
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					ch := arguments[1].Data.(chan funl.Value)
					ch <- arguments[2]
				}
			}
		}()
		return
	}
}
