package std

import (
	"fmt"

	"github.com/anssihalmeaho/funl/funl"
)

func initSTDLex(interpreter *funl.Interpreter) (err error) {
	stdModuleName := "stdlex"
	topFrame := funl.NewTopFrameWithInterpreter(interpreter)
	stdAstFuncs := []stdFuncInfo{
		{
			Name:       "tokenize",
			Getter:     getTokenize,
			IsFunction: true,
		},
	}
	err = setSTDFunctions(topFrame, stdModuleName, stdAstFuncs, interpreter)
	return
}

func getTokenize(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), needs one", name, l)
		}
		if arguments[0].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: assuming string", name)
		}
		tokenizer := funl.NewTokenizer(funl.NewDefaultOperators())
		tokens, err := tokenizer.Scan(arguments[0].Data.(string))
		if err != nil {
			values := []funl.Value{
				{
					Kind: funl.BoolValue,
					Data: false,
				},
				{
					Kind: funl.StringValue,
					Data: err.Error(),
				},
				{
					Kind: funl.StringValue,
					Data: "",
				},
			}
			retVal = funl.MakeListOfValues(frame, values)
			return
		}

		tokenValSlice := []funl.Value{}
		for _, token := range tokens {
			mapval := funl.HandleMapOP(frame, []*funl.Item{})

			putArgs := []*funl.Item{
				&funl.Item{Type: funl.ValueItem, Data: mapval},
				&funl.Item{Type: funl.ValueItem, Data: funl.Value{Kind: funl.StringValue, Data: "value"}},
				&funl.Item{Type: funl.ValueItem, Data: funl.Value{Kind: funl.StringValue, Data: token.Value}},
			}
			mapval = funl.HandlePutOP(frame, putArgs)

			putArgs = []*funl.Item{
				&funl.Item{Type: funl.ValueItem, Data: mapval},
				&funl.Item{Type: funl.ValueItem, Data: funl.Value{Kind: funl.StringValue, Data: "type"}},
				&funl.Item{Type: funl.ValueItem, Data: funl.Value{
					Kind: funl.StringValue,
					Data: fmt.Sprintf("%s", token.Type),
				}},
			}
			mapval = funl.HandlePutOP(frame, putArgs)

			putArgs = []*funl.Item{
				&funl.Item{Type: funl.ValueItem, Data: mapval},
				&funl.Item{Type: funl.ValueItem, Data: funl.Value{Kind: funl.StringValue, Data: "line"}},
				&funl.Item{Type: funl.ValueItem, Data: funl.Value{
					Kind: funl.IntValue,
					Data: token.Lineno,
				}},
			}
			mapval = funl.HandlePutOP(frame, putArgs)

			putArgs = []*funl.Item{
				&funl.Item{Type: funl.ValueItem, Data: mapval},
				&funl.Item{Type: funl.ValueItem, Data: funl.Value{Kind: funl.StringValue, Data: "pos"}},
				&funl.Item{Type: funl.ValueItem, Data: funl.Value{
					Kind: funl.IntValue,
					Data: token.Pos,
				}},
			}
			mapval = funl.HandlePutOP(frame, putArgs)

			tokenValSlice = append(tokenValSlice, mapval)
		}
		values := []funl.Value{
			{
				Kind: funl.BoolValue,
				Data: true,
			},
			{
				Kind: funl.StringValue,
				Data: "",
			},
			funl.MakeListOfValues(frame, tokenValSlice),
		}
		retVal = funl.MakeListOfValues(frame, values)
		return
	}
}
