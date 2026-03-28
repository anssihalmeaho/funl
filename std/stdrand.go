package std

import (
	"math/rand/v2"
	"strings"

	"github.com/anssihalmeaho/funl/funl"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func initSTDRand(interpreter *funl.Interpreter) (err error) {
	stdModuleName := "stdrand"
	topFrame := funl.NewTopFrameWithInterpreter(interpreter)
	stdFuncs := []stdFuncInfo{
		{
			Name:   "int",
			Getter: getStdRandInt,
		},
		{
			Name:   "intN",
			Getter: getStdRandIntN,
		},
		{
			Name:   "float",
			Getter: getStdRandFloat,
		},
		{
			Name:   "string",
			Getter: getStdRandString,
		},
		{
			Name:   "perm",
			Getter: getStdRandPerm,
		},
	}
	err = setSTDFunctions(topFrame, stdModuleName, stdFuncs, interpreter)
	return
}

func getStdRandPerm(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.IntValue {
			funl.RunTimeError2(frame, "%s: requires int value", name)
		}
		vlist := []funl.Value{}
		for _, v := range rand.Perm(arguments[0].Data.(int)) {
			vlist = append(vlist, funl.Value{Kind: funl.IntValue, Data: v})
		}
		return funl.MakeListOfValues(frame, vlist)
	}
}

func getStdRandString(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.IntValue {
			funl.RunTimeError2(frame, "%s: requires int value", name)
		}
		n := arguments[0].Data.(int)

		var sb strings.Builder
		for range n {
			// Pick a random character from the charset
			sb.WriteByte(charset[rand.IntN(len(charset))])
		}
		return funl.Value{Kind: funl.StringValue, Data: sb.String()}
	}
}

func getStdRandInt(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		return funl.Value{Kind: funl.IntValue, Data: rand.Int()}
	}
}

func getStdRandFloat(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		return funl.Value{Kind: funl.FloatValue, Data: rand.Float64()}
	}
}

func getStdRandIntN(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.IntValue {
			funl.RunTimeError2(frame, "%s: requires int value", name)
		}
		return funl.Value{Kind: funl.IntValue, Data: rand.IntN(arguments[0].Data.(int))}
	}
}
