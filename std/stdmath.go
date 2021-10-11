package std

import (
	"math"

	"github.com/anssihalmeaho/funl/funl"
)

func initSTDMath(interpreter *funl.Interpreter) (err error) {
	stdModuleName := "stdmath"
	topFrame := &funl.Frame{
		Syms:     funl.NewSymt(),
		OtherNS:  make(map[funl.SymID]funl.ImportInfo),
		Imported: make(map[funl.SymID]*funl.Frame),
	}
	stdMathFuncs := []stdFuncInfo{
		{
			Name:       "is-nan",
			Getter:     getStdMathIsNaN,
			IsFunction: true,
		},
		{
			Name:       "is-inf",
			Getter:     getStdMathIsInf,
			IsFunction: true,
		},
		{
			Name:       "abs",
			Getter:     getStdMathAbs,
			IsFunction: true,
		},
		{
			Name:       "acos",
			Getter:     getStdMathAcos,
			IsFunction: true,
		},
		{
			Name:       "acosh",
			Getter:     getStdMathAcosh,
			IsFunction: true,
		},
		{
			Name:       "asin",
			Getter:     getStdMathAsin,
			IsFunction: true,
		},
		{
			Name:       "asinh",
			Getter:     getStdMathAsinh,
			IsFunction: true,
		},
		{
			Name:       "atan",
			Getter:     getStdMathAtan,
			IsFunction: true,
		},
		{
			Name:       "atanh",
			Getter:     getStdMathAtanh,
			IsFunction: true,
		},
		{
			Name:       "cos",
			Getter:     getStdMathCos,
			IsFunction: true,
		},
		{
			Name:       "cosh",
			Getter:     getStdMathCosh,
			IsFunction: true,
		},
		{
			Name:       "sin",
			Getter:     getStdMathSin,
			IsFunction: true,
		},
		{
			Name:       "sinh",
			Getter:     getStdMathSinh,
			IsFunction: true,
		},
		{
			Name:       "tan",
			Getter:     getStdMathTan,
			IsFunction: true,
		},
		{
			Name:       "tanh",
			Getter:     getStdMathTanh,
			IsFunction: true,
		},
		{
			Name:       "ceil",
			Getter:     getStdMathCeil,
			IsFunction: true,
		},
		{
			Name:       "floor",
			Getter:     getStdMathFloor,
			IsFunction: true,
		},
		{
			Name:       "trunc",
			Getter:     getStdMathTrunc,
			IsFunction: true,
		},
		{
			Name:       "frexp",
			Getter:     getStdMathFrexp,
			IsFunction: true,
		},
		{
			Name:       "ldexp",
			Getter:     getStdMathLdexp,
			IsFunction: true,
		},
		{
			Name:       "modf",
			Getter:     getStdMathModf,
			IsFunction: true,
		},
		{
			Name:       "remainder",
			Getter:     getStdMathRemainder,
			IsFunction: true,
		},
		{
			Name:       "exp",
			Getter:     getStdMathExp,
			IsFunction: true,
		},
		{
			Name:       "exp2",
			Getter:     getStdMathExp2,
			IsFunction: true,
		},
		{
			Name:       "expm1",
			Getter:     getStdMathExpm1,
			IsFunction: true,
		},
		{
			Name:       "log",
			Getter:     getStdMathLog,
			IsFunction: true,
		},
		{
			Name:       "log10",
			Getter:     getStdMathLog10,
			IsFunction: true,
		},
		{
			Name:       "log1p",
			Getter:     getStdMathLog1p,
			IsFunction: true,
		},
		{
			Name:       "log2",
			Getter:     getStdMathLog2,
			IsFunction: true,
		},
		{
			Name:       "logb",
			Getter:     getStdMathLogb,
			IsFunction: true,
		},
		{
			Name:       "pow",
			Getter:     getStdMathPow,
			IsFunction: true,
		},
		{
			Name:       "sqrt",
			Getter:     getStdMathSqrt,
			IsFunction: true,
		},
		{
			Name:       "cbrt",
			Getter:     getStdMathCbrt,
			IsFunction: true,
		},
	}
	err = setSTDFunctions(topFrame, stdModuleName, stdMathFuncs, interpreter)

	item := &funl.Item{Type: funl.ValueItem, Data: funl.Value{Kind: funl.FloatValue, Data: math.Pi}}
	err = topFrame.Syms.Add("pi", item)
	if err != nil {
		return
	}

	return
}

func getStdMathIsNaN(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need one", name, l)
		}
		if arguments[0].Kind != funl.FloatValue {
			funl.RunTimeError2(frame, "%s: requires float value as argument", name)
		}
		retVal = funl.Value{Kind: funl.BoolValue, Data: math.IsNaN(arguments[0].Data.(float64))}
		return
	}
}

func getStdMathIsInf(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		switch l := len(arguments); l {
		case 1:
			fval := arguments[0].Data.(float64)
			retVal = funl.Value{Kind: funl.BoolValue, Data: math.IsInf(fval, 0)}
			return

		case 2:
			fval := arguments[0].Data.(float64)
			if arguments[1].Kind != funl.StringValue {
				funl.RunTimeError2(frame, "%s: requires string value as argument", name)
			}
			switch arguments[1].Data.(string) {
			case "+":
				retVal = funl.Value{Kind: funl.BoolValue, Data: math.IsInf(fval, 1)}
			case "-":
				retVal = funl.Value{Kind: funl.BoolValue, Data: math.IsInf(fval, -1)}
			default:
				funl.RunTimeError2(frame, "%s: requires '+' or '-' as input argument", name)
			}
			return

		default:
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		return
	}
}

func getStdMathPow(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 2 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need two", name, l)
		}
		if arguments[0].Kind != funl.FloatValue {
			funl.RunTimeError2(frame, "%s: requires float value as 1st argument", name)
		}
		if arguments[1].Kind != funl.FloatValue {
			funl.RunTimeError2(frame, "%s: requires float value as 2nd argument", name)
		}
		result := math.Pow(arguments[0].Data.(float64), arguments[1].Data.(float64))
		retVal = funl.Value{Kind: funl.FloatValue, Data: result}
		return
	}
}

func getStdMathLdexp(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 2 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need two", name, l)
		}
		if arguments[0].Kind != funl.FloatValue {
			funl.RunTimeError2(frame, "%s: requires float value as 1st argument", name)
		}
		if arguments[1].Kind != funl.IntValue {
			funl.RunTimeError2(frame, "%s: requires int value as 2nd argument", name)
		}
		result := math.Ldexp(arguments[0].Data.(float64), arguments[1].Data.(int))
		retVal = funl.Value{Kind: funl.FloatValue, Data: result}
		return
	}
}

func getStdMathRemainder(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 2 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need two", name, l)
		}
		if arguments[0].Kind != funl.FloatValue {
			funl.RunTimeError2(frame, "%s: requires float value as 1st argument", name)
		}
		if arguments[1].Kind != funl.FloatValue {
			funl.RunTimeError2(frame, "%s: requires float value as 2nd argument", name)
		}
		result := math.Mod(arguments[0].Data.(float64), arguments[1].Data.(float64))
		retVal = funl.Value{Kind: funl.FloatValue, Data: result}
		return
	}
}

func getStdMathModf(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need one", name, l)
		}
		if arguments[0].Kind != funl.FloatValue {
			funl.RunTimeError2(frame, "%s: requires float value as argument", name)
		}
		intPart, frac := math.Modf(arguments[0].Data.(float64))
		values := []funl.Value{
			{
				Kind: funl.FloatValue,
				Data: intPart,
			},
			{
				Kind: funl.FloatValue,
				Data: frac,
			},
		}
		retVal = funl.MakeListOfValues(frame, values)
		return
	}
}

func getStdMathFrexp(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need one", name, l)
		}
		if arguments[0].Kind != funl.FloatValue {
			funl.RunTimeError2(frame, "%s: requires float value as argument", name)
		}
		frac, exp := math.Frexp(arguments[0].Data.(float64))
		values := []funl.Value{
			{
				Kind: funl.FloatValue,
				Data: frac,
			},
			{
				Kind: funl.IntValue,
				Data: exp,
			},
		}
		retVal = funl.MakeListOfValues(frame, values)
		return
	}
}

func getMathFunction(name string, mathFunc func(float64) float64) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need one", name, l)
		}
		if arguments[0].Kind != funl.FloatValue {
			funl.RunTimeError2(frame, "%s: requires float value as argument", name)
		}
		fval := mathFunc(arguments[0].Data.(float64))
		retVal = funl.Value{Kind: funl.FloatValue, Data: fval}
		return
	}
}

func getStdMathCbrt(name string) stdFuncType {
	return getMathFunction(name, math.Cbrt)
}

func getStdMathSqrt(name string) stdFuncType {
	return getMathFunction(name, math.Sqrt)
}

func getStdMathLogb(name string) stdFuncType {
	return getMathFunction(name, math.Logb)
}

func getStdMathLog2(name string) stdFuncType {
	return getMathFunction(name, math.Log2)
}

func getStdMathLog1p(name string) stdFuncType {
	return getMathFunction(name, math.Log1p)
}

func getStdMathLog10(name string) stdFuncType {
	return getMathFunction(name, math.Log10)
}

func getStdMathLog(name string) stdFuncType {
	return getMathFunction(name, math.Log)
}

func getStdMathExpm1(name string) stdFuncType {
	return getMathFunction(name, math.Expm1)
}

func getStdMathExp(name string) stdFuncType {
	return getMathFunction(name, math.Exp)
}

func getStdMathExp2(name string) stdFuncType {
	return getMathFunction(name, math.Exp2)
}

func getStdMathTrunc(name string) stdFuncType {
	return getMathFunction(name, math.Trunc)
}

func getStdMathFloor(name string) stdFuncType {
	return getMathFunction(name, math.Floor)
}

func getStdMathCeil(name string) stdFuncType {
	return getMathFunction(name, math.Ceil)
}

func getStdMathTanh(name string) stdFuncType {
	return getMathFunction(name, math.Tanh)
}

func getStdMathTan(name string) stdFuncType {
	return getMathFunction(name, math.Tan)
}

func getStdMathSinh(name string) stdFuncType {
	return getMathFunction(name, math.Sinh)
}

func getStdMathSin(name string) stdFuncType {
	return getMathFunction(name, math.Sin)
}

func getStdMathCosh(name string) stdFuncType {
	return getMathFunction(name, math.Cosh)
}

func getStdMathCos(name string) stdFuncType {
	return getMathFunction(name, math.Cos)
}

func getStdMathAtanh(name string) stdFuncType {
	return getMathFunction(name, math.Atanh)
}

func getStdMathAtan(name string) stdFuncType {
	return getMathFunction(name, math.Atan)
}

func getStdMathAsinh(name string) stdFuncType {
	return getMathFunction(name, math.Asinh)
}

func getStdMathAsin(name string) stdFuncType {
	return getMathFunction(name, math.Asin)
}

func getStdMathAcosh(name string) stdFuncType {
	return getMathFunction(name, math.Acosh)
}

func getStdMathAcos(name string) stdFuncType {
	return getMathFunction(name, math.Acos)
}

func getStdMathAbs(name string) stdFuncType {
	return getMathFunction(name, math.Abs)
}
