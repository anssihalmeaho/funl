package std

import (
	"github.com/anssihalmeaho/funl"
	"strings"
	"unicode"
)

func initSTDStr() (err error) {
	stdModuleName := "stdstr"
	topFrame := &funl.Frame{
		Syms:     funl.NewSymt(),
		OtherNS:  make(map[funl.SymID]funl.ImportInfo),
		Imported: make(map[funl.SymID]*funl.Frame),
	}
	stdStrFuncs := []stdFuncInfo{
		{
			Name:       "lowercase",
			Getter:     getStdStrLowercase,
			IsFunction: true,
		},
		{
			Name:       "uppercase",
			Getter:     getStdStrUpperercase,
			IsFunction: true,
		},
		{
			Name:       "replace",
			Getter:     getStdStrReplace,
			IsFunction: true,
		},
		{
			Name:       "strip",
			Getter:     getStdStrStrip,
			IsFunction: true,
		},
		{
			Name:       "rstrip",
			Getter:     getStdStrRstrip,
			IsFunction: true,
		},
		{
			Name:       "lstrip",
			Getter:     getStdStrLstrip,
			IsFunction: true,
		},
		{
			Name:       "join",
			Getter:     getStdStrJoin,
			IsFunction: true,
		},
		{
			Name:       "startswith",
			Getter:     getStdStrStartswith,
			IsFunction: true,
		},
		{
			Name:       "endswith",
			Getter:     getStdStrEndswith,
			IsFunction: true,
		},
		{
			Name:       "is-space",
			Getter:     getStdStrIsSpace,
			IsFunction: true,
		},
		{
			Name:       "is-digit",
			Getter:     getStdStrIsDigit,
			IsFunction: true,
		},
		{
			Name:       "is-lower",
			Getter:     getStdStrIsLower,
			IsFunction: true,
		},
		{
			Name:       "is-upper",
			Getter:     getStdStrIsUpper,
			IsFunction: true,
		},
		{
			Name:       "is-alpha",
			Getter:     getStdStrIsAlpha,
			IsFunction: true,
		},
		{
			Name:       "is-alnum",
			Getter:     getStdStrIsAlnum,
			IsFunction: true,
		},
	}
	err = setSTDFunctions(topFrame, stdModuleName, stdStrFuncs)
	return
}

func isCondition(name string, frame *funl.Frame, arguments []funl.Value, cond func(r rune) bool) (retVal funl.Value) {
	if l := len(arguments); l != 1 {
		funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need one", name, l)
	}
	if arguments[0].Kind != funl.StringValue {
		funl.RunTimeError2(frame, "%s: requires string value as 1st argument", name)
	}
	runes := []rune(arguments[0].Data.(string))
	condApplies := true
	for _, r := range runes {
		if !cond(r) {
			condApplies = false
			break
		}
	}
	retVal = funl.Value{Kind: funl.BoolValue, Data: condApplies}
	return
}

func getStdStrIsAlnum(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		isAlnum := func(r rune) bool {
			if unicode.IsLetter(r) {
				return true
			}
			if unicode.IsNumber(r) {
				return true
			}
			return false
		}
		retVal = isCondition(name, frame, arguments, isAlnum)
		return
	}
}

func getStdStrIsAlpha(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		retVal = isCondition(name, frame, arguments, unicode.IsLetter)
		return
	}
}

func getStdStrIsUpper(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		isUpperCharacter := func(r rune) bool {
			return !unicode.IsLower(r)
		}
		retVal = isCondition(name, frame, arguments, isUpperCharacter)
		return
	}
}

func getStdStrIsLower(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		isLowerCharacter := func(r rune) bool {
			return !unicode.IsUpper(r)
		}
		retVal = isCondition(name, frame, arguments, isLowerCharacter)
		return
	}
}

func getStdStrIsDigit(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		retVal = isCondition(name, frame, arguments, unicode.IsDigit)
		return
	}
}

func getStdStrIsSpace(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		retVal = isCondition(name, frame, arguments, unicode.IsSpace)
		return
	}
}

func getStdStrEndswith(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 2 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need two", name, l)
		}
		if arguments[0].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: requires string value as 1st argument", name)
		}
		if arguments[1].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: requires string value as 2nd argument", name)
		}
		isIt := strings.HasSuffix(arguments[0].Data.(string), arguments[1].Data.(string))
		retVal = funl.Value{Kind: funl.BoolValue, Data: isIt}
		return
	}
}

func getStdStrStartswith(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 2 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need two", name, l)
		}
		if arguments[0].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: requires string value as 1st argument", name)
		}
		if arguments[1].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: requires string value as 2nd argument", name)
		}
		isIt := strings.HasPrefix(arguments[0].Data.(string), arguments[1].Data.(string))
		retVal = funl.Value{Kind: funl.BoolValue, Data: isIt}
		return
	}
}

func getStdStrJoin(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 2 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need two", name, l)
		}
		if arguments[0].Kind != funl.ListValue {
			funl.RunTimeError2(frame, "%s: requires list value as 1st argument", name)
		}
		if arguments[1].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: requires string value as 2nd argument", name)
		}
		lit := funl.NewListIterator(arguments[0])
		var strs []string
		for {
			sval := lit.Next()
			if sval == nil {
				break
			}
			if sval.Kind != funl.StringValue {
				funl.RunTimeError2(frame, "%s: requires string value for join", name)
			}
			strs = append(strs, sval.Data.(string))
		}
		resultStr := strings.Join(strs, arguments[1].Data.(string))
		retVal = funl.Value{Kind: funl.StringValue, Data: resultStr}
		return
	}
}

func getStdStrLstrip(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need one", name, l)
		}
		if arguments[0].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: requires string value", name)
		}
		wsChecker := func(r rune) bool {
			return unicode.IsSpace(r)
		}
		resultStr := strings.TrimLeftFunc(arguments[0].Data.(string), wsChecker)
		retVal = funl.Value{Kind: funl.StringValue, Data: resultStr}
		return
	}
}

func getStdStrRstrip(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need one", name, l)
		}
		if arguments[0].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: requires string value", name)
		}
		wsChecker := func(r rune) bool {
			return unicode.IsSpace(r)
		}
		resultStr := strings.TrimRightFunc(arguments[0].Data.(string), wsChecker)
		retVal = funl.Value{Kind: funl.StringValue, Data: resultStr}
		return
	}
}

func getStdStrStrip(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need one", name, l)
		}
		if arguments[0].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: requires string value", name)
		}
		retVal = funl.Value{Kind: funl.StringValue, Data: strings.TrimSpace(arguments[0].Data.(string))}
		return
	}
}

func getStdStrReplace(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 3 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need three", name, l)
		}
		for i := 0; i < 3; i++ {
			if arguments[i].Kind != funl.StringValue {
				funl.RunTimeError2(frame, "%s: requires string value (%d. argument)", name, i+1)
			}
		}
		resulStr := strings.Replace(arguments[0].Data.(string), arguments[1].Data.(string), arguments[2].Data.(string), -1)
		retVal = funl.Value{Kind: funl.StringValue, Data: resulStr}
		return
	}
}

func getStdStrUpperercase(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need one", name, l)
		}
		if arguments[0].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: requires string value", name)
		}
		retVal = funl.Value{Kind: funl.StringValue, Data: strings.ToUpper(arguments[0].Data.(string))}
		return
	}
}

func getStdStrLowercase(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need one", name, l)
		}
		if arguments[0].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: requires string value", name)
		}
		retVal = funl.Value{Kind: funl.StringValue, Data: strings.ToLower(arguments[0].Data.(string))}
		return
	}
}
