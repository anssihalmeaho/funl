package funl

import (
	"fmt"
	"reflect"
)

func handleTryOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "try"

	if !frame.inProcCall {
		runTimeError2(frame, "%s not allowed in function", opName)
	}

	operCount := len(operands)
	if (operCount != 2) && (operCount != 1) {
		runTimeError2(frame, "Wrong amount of arguments for %s (%d given)", opName, operCount)
	}

	var rteText string
	isFailure := false
	v := operands[0]
	switch v.Type {
	case ValueItem:
		retVal = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		func() {
			defer func() {
				if isFailure {
					if r := recover(); r != nil {
						rteText = r.(error).Error()
					}
				}
			}()
			isFailure = true
			retVal = EvalItem(v, frame)
			isFailure = false
		}()
	default:
		isFailure = true
	}

	if !isFailure {
		return
	}

	if operCount == 1 {
		retVal = Value{Kind: StringValue, Data: "RTE:" + rteText}
		return
	}

	v = operands[1]
	switch v.Type {
	case ValueItem:
		retVal = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		retVal = EvalItem(v, frame)
	}

	return
}

func handleSpawnOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "spawn"

	if !frame.inProcCall {
		runTimeError2(frame, "%s not allowed in function", opName)
	}

	for _, operand := range operands {
		go func(it *Item) {
			isFailure := false
			defer func() {
				if isFailure {
					var rteText string
					if r := recover(); r != nil {
						rteText = r.(error).Error()
					}
					fmt.Println()
					fmt.Println("Fiber died, Runtime error: ", rteText)
				}
			}()
			isFailure = true
			EvalItem(it, frame)
			isFailure = false
		}(operand)
	}
	retVal = Value{Kind: BoolValue, Data: true}
	return
}

func handleChanOP(frame *Frame, operands []*Item) (retVal Value) {
	chData := make(chan Value)
	retVal = Value{Kind: ChanValue, Data: chData}
	return
}

func handleSendOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "send"

	if !frame.inProcCall {
		runTimeError2(frame, "%s not allowed in function", opName)
	}

	if l := len(operands); l != 2 {
		runTimeError2(frame, "Wrong amount of arguments for %s (%d given)", opName, l)
	}

	v := operands[0]
	var chVal Value
	switch v.Type {
	case ValueItem:
		chVal = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		chVal = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}
	if chVal.Kind != ChanValue {
		runTimeError2(frame, "Expecting channel as 1st arg for %s", opName)
	}

	v = operands[1]
	var dataVal Value
	switch v.Type {
	case ValueItem:
		dataVal = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		dataVal = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	chVal.Data.(chan Value) <- dataVal

	retVal = Value{Kind: BoolValue, Data: true}
	return
}

func handleRecvOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "recv"

	if !frame.inProcCall {
		runTimeError2(frame, "%s not allowed in function", opName)
	}

	if l := len(operands); l != 1 {
		runTimeError2(frame, "Wrong amount of arguments for %s (%d given)", opName, l)
	}

	v := operands[0]
	var chVal Value
	switch v.Type {
	case ValueItem:
		chVal = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		chVal = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}
	if chVal.Kind != ChanValue {
		runTimeError2(frame, "Expecting channel as 1st arg for %s", opName)
	}

	retVal = <-chVal.Data.(chan Value)
	return
}

func handleSelectOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "select"

	if !frame.inProcCall {
		runTimeError2(frame, "%s not allowed in function", opName)
	}

	operCount := len(operands)
	if (operCount == 0) || (operCount%2) != 0 {
		runTimeError2(frame, "Wrong amount of arguments for %s (%d given)", opName, operCount)
	}

	var chans []chan Value
	var procs []FuncValue

	var isTwoListCase bool
	if operCount == 2 {
		isTwoListCase = func() bool {
			v1 := operands[0]
			var chlistVal Value
			switch v1.Type {
			case ValueItem:
				chlistVal = v1.Data.(Value)
			case SymbolPathItem, OperCallItem:
				chlistVal = EvalItem(v1, frame)
			default:
				runTimeError2(frame, "something wrong (%s)", opName)
			}
			if chlistVal.Kind != ListValue {
				return false
			}

			v2 := operands[1]
			var flistVal Value
			switch v2.Type {
			case ValueItem:
				flistVal = v2.Data.(Value)
			case SymbolPathItem, OperCallItem:
				flistVal = EvalItem(v2, frame)
			default:
				runTimeError2(frame, "something wrong (%s)", opName)
			}
			if flistVal.Kind != ListValue {
				return false
			}

			chit := NewListIterator(chlistVal)
			for {
				nextch := chit.Next()
				if nextch == nil {
					break
				}
				if nextch.Kind != ChanValue {
					runTimeError2(frame, "%s: assuming channel in 1st list", opName)
				}
				chans = append(chans, nextch.Data.(chan Value))
			}

			fit := NewListIterator(flistVal)
			for {
				nextfu := fit.Next()
				if nextfu == nil {
					break
				}
				if nextfu.Kind != FunctionValue {
					runTimeError2(frame, "%s: assuming func/proc in 2nd list", opName)
				}
				procs = append(procs, nextfu.Data.(FuncValue))
			}

			if l1, l2 := len(chans), len(procs); l1 != l2 {
				runTimeError2(frame, "%s: lists have not same length (1st: %d)(2nd: %d)", opName, l1, l2)
			}
			return true
		}()
	}

	if !isTwoListCase {
		for i, v := range operands {
			var val Value
			switch v.Type {
			case ValueItem:
				val = v.Data.(Value)
			case SymbolPathItem, OperCallItem:
				val = EvalItem(v, frame)
			default:
				runTimeError2(frame, "something wrong (%s)", opName)
			}

			if (i % 2) == 0 {
				if val.Kind != ChanValue {
					runTimeError2(frame, "Expecting channel as arg for %s", opName)
				}
				chans = append(chans, val.Data.(chan Value))
			} else {
				//TODO: should external procs be supported too ?
				if val.Kind != FunctionValue {
					runTimeError2(frame, "Expecting func/proc as arg for %s", opName)
				}
				procs = append(procs, val.Data.(FuncValue))
			}
		}
	}

	var cases []reflect.SelectCase
	for _, ch := range chans {
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
			Send: reflect.Value{},
		})
	}
	i, receivedVal, ok := reflect.Select(cases)
	if !ok {
		runTimeError2(frame, "%s: error in receiving", opName)
	}
	recval := receivedVal.Interface().(Value)
	fitem := &Item{Type: ValueItem, Data: Value{Kind: FunctionValue, Data: procs[i]}}
	argsForCall := []*Item{
		fitem,
		&Item{Type: ValueItem, Data: recval},
	}

	retVal = handleCallOP(frame, argsForCall)
	return
}
