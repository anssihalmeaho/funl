package funl

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func handleArgslistOP(frame *Frame, operands []*Item) (retVal Value) {
	retVal = MakeListOfValues(frame, frame.EvaluatedArgs)
	return
}

func HandleSprintfOP(frame *Frame, operands []*Item) (retVal Value) {
	return handleSprintfOP(frame, operands)
}

func handleSprintfOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "sprintf"
	if l := len(operands); l == 0 {
		runTimeError2(frame, "%s: wrong amount of arguments (%d), need at least one", opName, l)
	}

	v := operands[0]
	var fmtv Value
	switch v.Type {
	case ValueItem:
		fmtv = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		fmtv = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	if fmtv.Kind != StringValue {
		runTimeError2(frame, "%s: requires string value", opName)
	}
	formatStr, ok := fmtv.Data.(string)
	if !ok {
		runTimeError2(frame, "%s: argument is not string value", opName)
	}
	var args []interface{}
	for _, operand := range operands[1:] {

		var srcVal Value
		switch operand.Type {
		case ValueItem:
			srcVal = operand.Data.(Value)
		case SymbolPathItem, OperCallItem:
			srcVal = EvalItem(operand, frame)
		default:
			runTimeError2(frame, "something wrong (%s)", opName)
		}

		switch srcVal.Kind {
		case IntValue:
			args = append(args, srcVal.Data.(int))
		case StringValue:
			args = append(args, srcVal.Data.(string))
		case FloatValue:
			args = append(args, srcVal.Data.(float64))
		case BoolValue:
			args = append(args, srcVal.Data.(bool))
		case OpaqueValue:
			if !frame.inProcCall {
				runTimeError2(frame, "%s not allowed in function for opaque value", opName)
			}
			args = append(args, srcVal)
		default:
			args = append(args, srcVal)
		}
	}
	retVal = Value{Kind: StringValue, Data: fmt.Sprintf(formatStr, args...)}
	return
}

var handleHelpOP = func() OperHandler {
	operInfo := NewDefaultOperators()
	operDocs := NewOperatorDocs()
	// need to show operator names in fixed order
	// otherwise Go random map handling would cause
	// impure effect in functions
	operNames := []string{}
	for operatorName := range operDocs {
		operNames = append(operNames, operatorName)
	}
	sort.Strings(operNames)

	return func(frame *Frame, operands []*Item) (retVal Value) {
		opName := "help"
		argCount := len(operands)
		if argCount > 1 {
			runTimeError2(frame, "%s operator needs at most one argument (%d given)", opName, argCount)
		}
		var text string

		if argCount == 0 {
			text = `
This is online help.
Usage: help('..topic-name...')
  topic can be:
	- operator name, like 'plus'
	- 'operators'
`
			retVal = Value{Kind: StringValue, Data: text}
			return
		}
		v := operands[0]
		var srcVal Value
		switch v.Type {
		case ValueItem:
			srcVal = v.Data.(Value)
		case SymbolPathItem, OperCallItem:
			srcVal = EvalItem(v, frame)
		default:
			runTimeError2(frame, "something wrong (%s)", opName)
		}
		if srcVal.Kind != StringValue {
			runTimeError2(frame, "%s: assuming string value as argument", opName)
		}
		topicStr := srcVal.Data.(string)
		if topicStr == "operators" {
			var operNameValues []Value
			for _, operName := range operNames {
				operNameValues = append(operNameValues, Value{Kind: StringValue, Data: operName})
			}
			retVal = MakeListOfValues(frame, operNameValues)
			return
		} else if operInfo.isOperator(topicStr) {
			text = operDocs[topicStr]
		}

		retVal = Value{Kind: StringValue, Data: text}
		return
	}
}()

func handleFloatOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "eval"
	if l := len(operands); l != 1 {
		runTimeError2(frame, "%s operator needs one argument (%d given)", opName, l)
	}

	v := operands[0]
	var srcVal Value
	switch v.Type {
	case ValueItem:
		srcVal = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		srcVal = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	switch srcVal.Kind {
	case IntValue:
		intv := srcVal.Data.(int)
		retVal = Value{Kind: FloatValue, Data: float64(intv)}
	case FloatValue:
		retVal = srcVal
	default:
		runTimeError2(frame, "%s assumes int (or float) as argument", opName)
	}

	return
}

type evalErrHandler struct{}

func (eh *evalErrHandler) handleParseError(errorText string) {
	runTimeError(errorText)
}

func handleEvalOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "eval"
	if l := len(operands); l != 1 {
		runTimeError2(frame, "%s operator needs one argument (%d given)", opName, l)
	}

	v := operands[0]
	var srcVal Value
	switch v.Type {
	case ValueItem:
		srcVal = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		srcVal = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	if srcVal.Kind != StringValue {
		runTimeError2(frame, "%s assumes string as argument", opName)
	}

	parser := NewParser(NewDefaultOperators(), nil)
	parser.SetErrorHandler(&evalErrHandler{})
	evalItem, err := parser.ParseOneExpression(srcVal.Data.(string))
	if err != nil {
		runTimeError2(frame, "%s: error in parsing expression (%v)", opName, err)
	}
	retVal = EvalItem(evalItem, frame)
	return
}

func handleImpOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "imp"
	//TODO: could be 2 args also, import path
	if l := len(operands); l != 1 {
		runTimeError2(frame, "Wrong amount of arguments for %s (%d given)", opName, l)
	}

	var modName string
	switch v := operands[0]; v.Type {
	case ValueItem:
		srcVal := v.Data.(Value)
		if srcVal.Kind != StringValue {
			runTimeError2(frame, "%s: should be string as 1st argument (if value)", opName)
		}
		modName = srcVal.Data.(string)
	case SymbolPathItem:
		sp := v.Data.(SymbolPath)
		modName = sp.ToString()
	default:
		runTimeError2(frame, "%s: assuming symbol as 1st argument", opName)
	}

	if frame.FuncProto == nil {
		runTimeError2(frame, "%s: not supported at main level", opName)
		return
	}
	sid := SymIDMap.Add(modName)
	if _, modFound := frame.FuncProto.NSpace.OtherNS[sid]; !modFound {
		frame.FuncProto.NSpace.OtherNS[sid] = ImportInfo{} // path added when given
	}
	AddImportsToNamespace(&frame.FuncProto.NSpace, frame)

	symt, found := frame.GetSymItemsOfImportedModule(sid)
	var moperands []*Item
	if found {
		for k, v := range symt.AsMap() {
			symNameVal := &Item{Type: ValueItem, Data: Value{Kind: StringValue, Data: SymIDMap.AsString(k)}}
			moperands = append(moperands, symNameVal, v)
		}
	}
	retVal = handleMapOP(frame, moperands)
	return
}

func handleLetOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "let"
	if l := len(operands); l != 2 {
		runTimeError2(frame, "Wrong amount of arguments for %s (%d given)", opName, l)
	}

	var symbolName string
	switch v := operands[0]; v.Type {
	case SymbolPathItem:
		sp := v.Data.(SymbolPath)
		symbolName = sp.ToString()
	default:
		runTimeError2(frame, "%s: assuming symbol as 1st argument", opName)
	}

	sid := SymIDMap.Add(symbolName)
	// Overlapping symbol names not allowed (might cause variable -like effect)
	if _, symfound := frame.GetSymItem(sid); symfound {
		runTimeError2(frame, "%s: Duplicate symbol name in scope not allowed (%s)", opName, SymIDMap.AsString(sid))
	}

	// lets evaluate let def. to value first, using new frame already (as arguments are there)
	letvalitem := &Item{Type: ValueItem, Data: EvalItem(operands[1], frame)}
	if !frame.Syms.AddBySID(sid, letvalitem) {
		runTimeError2(frame, "Symbol add failed")
	}

	retVal = letvalitem.Data.(Value)
	return
}

func handleSymvalOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "symval"
	if l := len(operands); l != 1 {
		runTimeError2(frame, "%s operator needs one argument (%d given)", opName, l)
	}

	if !frame.inProcCall {
		runTimeError2(frame, "%s not allowed in function", opName)
	}

	v := operands[0]
	var srcVal Value
	switch v.Type {
	case ValueItem:
		srcVal = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		srcVal = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	if srcVal.Kind != StringValue {
		runTimeError2(frame, "%s assumes string as argument", opName)
	}
	var sp SymbolPath
	for _, symstr := range strings.Split(srcVal.Data.(string), ".") {
		symsid, found := SymIDMap.Get(symstr)
		if !found {
			runTimeError2(frame, "%s: symbol not found (%s)", opName, srcVal.Data.(string))
		}
		sp = append(sp, symsid)
	}

	var symItem *Item
	var symfound bool
	if len(sp) == 1 {
		symItem, symfound = frame.GetSymItem(sp[0])
	} else if len(sp) > 1 {
		symItem, symfound = frame.GetImportedSymItem(sp[0], sp[1:])
	} else {
		runTimeError2(frame, "symbol not found: %s", sp.ToString())
	}
	if !symfound {
		runTimeError2(frame, "symbol not found: %s", sp.ToString())
	}

	switch symItem.Type {
	case ValueItem:
		retVal = symItem.Data.(Value)
	case SymbolPathItem, OperCallItem:
		retVal = EvalItem(symItem, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	return
}

var PrintingDisabledInFunctions bool

func handlePrintOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "print"

	if PrintingDisabledInFunctions {
		if !frame.inProcCall {
			retVal = Value{Kind: BoolValue, Data: true}
			return
		}
	}

	text := ""
	for _, v := range operands {
		var val Value
		switch v.Type {
		case ValueItem:
			val = v.Data.(Value)
		case SymbolPathItem, OperCallItem:
			val = EvalItem(v, frame)
		default:
			runTimeError2(frame, "something wrong (%s)", opName)
		}
		switch val.Kind {
		case StringValue:
			text += val.Data.(string)
		default:
			text += val.String()
		}
	}
	fmt.Println(text)
	retVal = Value{Kind: BoolValue, Data: true}
	return
}

func handleErrorOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "error"

	errorText := ""
	for _, v := range operands {
		var val Value
		switch v.Type {
		case ValueItem:
			val = v.Data.(Value)
		case SymbolPathItem, OperCallItem:
			val = EvalItem(v, frame)
		default:
			runTimeError2(frame, "something wrong (%s)", opName)
		}
		switch val.Kind {
		case StringValue:
			errorText += val.Data.(string)
		default:
			errorText += val.String()
		}
	}
	runTimeError2(frame, "%s", errorText)
	return
}

func handleNameOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "name"
	if l := len(operands); l != 1 {
		runTimeError2(frame, "%s operator needs one argument (%d given)", opName, l)
	}

	v := operands[0]
	switch v.Type {
	case SymbolPathItem:
		sp := v.Data.(SymbolPath)
		retVal = Value{Kind: StringValue, Data: sp.ToString()}
	case ValueItem:
		runTimeError2(frame, "%s: symbol assumed as argument, not value", opName)
	case OperCallItem:
		runTimeError2(frame, "%s: symbol assumed as argument, not operator call", opName)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}
	return
}

func handleCondOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "cond"
	argCount := len(operands)
	if l := argCount; l < 2 {
		runTimeError2(frame, "Wrong amount of arguments for %s (%d given)", opName, l)
	}
	if argCount%2 == 0 {
		runTimeError2(frame, "%s: else expression missing", opName)
	}

	for i := 0; i < argCount-1; i += 2 {
		v := operands[i]
		var condVal Value
		switch v.Type {
		case ValueItem:
			condVal = v.Data.(Value)
		case SymbolPathItem, OperCallItem:
			condVal = EvalItem(v, frame)
		default:
			runTimeError2(frame, "something wrong (%s)", opName)
		}

		if condVal.Kind != BoolValue {
			runTimeError2(frame, "%s: compared value assumed to be bool value (%d)", opName, i)
		}
		if condVal.Data.(bool) == true {
			v = operands[i+1]
			switch v.Type {
			case ValueItem:
				retVal = v.Data.(Value)
			case SymbolPathItem, OperCallItem:
				retVal = EvalItem(v, frame)
			default:
				runTimeError2(frame, "something wrong (%s)", opName)
			}
			return
		}
	}
	// no matches so lets return else expression
	v := operands[argCount-1]
	switch v.Type {
	case ValueItem:
		retVal = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		retVal = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}
	return
}

func handleCaseOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "case"
	argCount := len(operands)
	if l := argCount; l < 2 {
		runTimeError2(frame, "Wrong amount of arguments for %s (%d given)", opName, l)
	}

	v := operands[0]
	var caseVal Value
	switch v.Type {
	case ValueItem:
		caseVal = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		caseVal = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	itemAsVal := &Item{Type: ValueItem, Data: caseVal}

	var hasDefault bool
	matchCounter := argCount
	if (argCount % 2) == 0 {
		hasDefault = true
		matchCounter--
	}
	for i := 1; i < matchCounter; i += 2 {
		v = operands[i]
		var matchVal Value
		switch v.Type {
		case ValueItem:
			matchVal = v.Data.(Value)
		case SymbolPathItem, OperCallItem:
			matchVal = EvalItem(v, frame)
		default:
			runTimeError2(frame, "something wrong (%s)", opName)
		}

		argsForEq := []*Item{
			&Item{
				Type: ValueItem,
				Data: matchVal,
			},
			itemAsVal,
		}
		eqResult := handleEqOP(frame, argsForEq)
		if eqResult.Kind != BoolValue {
			runTimeError2(frame, "Invalid result from eq")
		}
		if eqResult.Data.(bool) == true {
			v = operands[i+1]
			switch v.Type {
			case ValueItem:
				retVal = v.Data.(Value)
			case SymbolPathItem, OperCallItem:
				retVal = EvalItem(v, frame)
			default:
				runTimeError2(frame, "something wrong (%s)", opName)
			}
			return
		}

	}

	if hasDefault {
		v = operands[argCount-1]
		switch v.Type {
		case ValueItem:
			retVal = v.Data.(Value)
		case SymbolPathItem, OperCallItem:
			retVal = EvalItem(v, frame)
		default:
			runTimeError2(frame, "something wrong (%s)", opName)
		}
		return
	}

	runTimeError2(frame, "%s: could not evaluate value", opName)
	return
}

func handleConvOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "conv"
	if l := len(operands); l != 2 {
		runTimeError2(frame, "Wrong amount of arguments for %s (%d given)", opName, l)
	}

	v := operands[0]
	var srcVal Value
	switch v.Type {
	case ValueItem:
		srcVal = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		srcVal = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	v = operands[1]
	var targetTypeVal Value
	switch v.Type {
	case ValueItem:
		targetTypeVal = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		targetTypeVal = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	if targetTypeVal.Kind != StringValue {
		runTimeError2(frame, "%s: 2nd argument assumed to be string", opName)
	}

	switch trgType := targetTypeVal.Data; trgType {
	case "bool":
		switch srcVal.Kind {
		case BoolValue:
			retVal = srcVal
		default:
			runTimeError2(frame, "%s: unsupported source type for converting %s", opName, trgType)
		}
	case "string":
		if srcVal.Kind == StringValue {
			retVal = srcVal
		} else {
			retVal = Value{Kind: StringValue, Data: srcVal.String()}
		}
	case "list":
		switch srcVal.Kind {
		case ListValue:
			retVal = srcVal
		case StringValue:
			targetList := &List{}
			var previtem *ListObject
			var newitem *ListObject
			for _, c := range srcVal.Data.(string) {
				newitem = &ListObject{Val: &Value{Kind: StringValue, Data: string(c)}, Next: nil}
				if targetList.Head == nil {
					targetList.Head = newitem
				}
				if previtem != nil {
					previtem.Next = newitem
				}
				previtem = newitem
			}
			retVal = Value{Kind: ListValue, Data: targetList}
		default:
			runTimeError2(frame, "%s: unsupported source type for converting %s", opName, trgType)
		}
	case "float":
		switch srcVal.Kind {
		case FloatValue:
			retVal = srcVal
		case IntValue:
			retVal = Value{Kind: FloatValue, Data: float64(srcVal.Data.(int))}
		default:
			runTimeError2(frame, "%s: unsupported source type for converting %s", opName, trgType)
		}
	case "int":
		switch srcVal.Kind {
		case IntValue:
			retVal = srcVal
		case FloatValue:
			intv := int(srcVal.Data.(float64))
			retVal = Value{Kind: IntValue, Data: intv}
		case StringValue:
			intv, err := strconv.Atoi(srcVal.Data.(string))
			if err != nil {
				retVal = Value{Kind: StringValue, Data: "Not able to convert to int"}
			} else {
				retVal = Value{Kind: IntValue, Data: intv}
			}
		default:
			runTimeError2(frame, "%s: unsupported source type for converting %s", opName, trgType)
		}
	case "funcproto":
		switch srcVal.Kind {
		case FuncProtoValue:
			retVal = srcVal
		default:
			runTimeError2(frame, "%s: unsupported source type for converting %s", opName, trgType)
		}
	case "function":
		switch srcVal.Kind {
		case FunctionValue:
			retVal = srcVal
		default:
			runTimeError2(frame, "%s: unsupported source type for converting %s", opName, trgType)
		}
	case "channel":
		switch srcVal.Kind {
		case ChanValue:
			retVal = srcVal
		default:
			runTimeError2(frame, "%s: unsupported source type for converting %s", opName, trgType)
		}
	}
	return
}

func handleStrOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "str"
	if l := len(operands); l != 1 {
		runTimeError2(frame, "%s operator needs one argument (%d given)", opName, l)
	}

	v := operands[0]
	var val Value
	switch v.Type {
	case ValueItem:
		val = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		val = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	if val.Kind == StringValue {
		retVal = val
		return
	}

	if !frame.inProcCall && val.Kind == OpaqueValue {
		runTimeError2(frame, "%s not allowed in function for opaque value", opName)
	}

	retVal = Value{Kind: StringValue, Data: val.String()}
	return
}

func handleLtOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "lt"
	if l := len(operands); l != 2 {
		runTimeError2(frame, "Wrong amount of arguments for %s (%d given)", opName, l)
	}

	v := operands[0]
	var val1 Value
	switch v.Type {
	case ValueItem:
		val1 = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		val1 = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	v = operands[1]
	var val2 Value
	switch v.Type {
	case ValueItem:
		val2 = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		val2 = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	if val1.Kind == IntValue && val2.Kind == IntValue {
		retVal = Value{Kind: BoolValue, Data: val1.Data.(int) < val2.Data.(int)}
	} else if val1.Kind == FloatValue && val2.Kind == IntValue {
		retVal = Value{Kind: BoolValue, Data: val1.Data.(float64) < float64(val2.Data.(int))}
	} else if val1.Kind == IntValue && val2.Kind == FloatValue {
		retVal = Value{Kind: BoolValue, Data: float64(val1.Data.(int)) < val2.Data.(float64)}
	} else if val1.Kind == FloatValue && val2.Kind == FloatValue {
		retVal = Value{Kind: BoolValue, Data: val1.Data.(float64) < val2.Data.(float64)}
	} else {
		runTimeError2(frame, "%s: invalid types", opName)
	}
	return
}

func handleLeOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "le"
	if l := len(operands); l != 2 {
		runTimeError2(frame, "Wrong amount of arguments for %s (%d given)", opName, l)
	}

	v := operands[0]
	var val1 Value
	switch v.Type {
	case ValueItem:
		val1 = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		val1 = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	v = operands[1]
	var val2 Value
	switch v.Type {
	case ValueItem:
		val2 = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		val2 = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	if val1.Kind == IntValue && val2.Kind == IntValue {
		retVal = Value{Kind: BoolValue, Data: val1.Data.(int) <= val2.Data.(int)}
	} else if val1.Kind == FloatValue && val2.Kind == IntValue {
		retVal = Value{Kind: BoolValue, Data: val1.Data.(float64) <= float64(val2.Data.(int))}
	} else if val1.Kind == IntValue && val2.Kind == FloatValue {
		retVal = Value{Kind: BoolValue, Data: float64(val1.Data.(int)) <= val2.Data.(float64)}
	} else if val1.Kind == FloatValue && val2.Kind == FloatValue {
		retVal = Value{Kind: BoolValue, Data: val1.Data.(float64) <= val2.Data.(float64)}
	} else {
		runTimeError2(frame, "%s: invalid types", opName)
	}
	return
}

func handleGeOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "ge"
	if l := len(operands); l != 2 {
		runTimeError2(frame, "Wrong amount of arguments for %s (%d given)", opName, l)
	}

	v := operands[0]
	var val1 Value
	switch v.Type {
	case ValueItem:
		val1 = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		val1 = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	v = operands[1]
	var val2 Value
	switch v.Type {
	case ValueItem:
		val2 = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		val2 = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	if val1.Kind == IntValue && val2.Kind == IntValue {
		retVal = Value{Kind: BoolValue, Data: val1.Data.(int) >= val2.Data.(int)}
	} else if val1.Kind == FloatValue && val2.Kind == IntValue {
		retVal = Value{Kind: BoolValue, Data: val1.Data.(float64) >= float64(val2.Data.(int))}
	} else if val1.Kind == IntValue && val2.Kind == FloatValue {
		retVal = Value{Kind: BoolValue, Data: float64(val1.Data.(int)) >= val2.Data.(float64)}
	} else if val1.Kind == FloatValue && val2.Kind == FloatValue {
		retVal = Value{Kind: BoolValue, Data: val1.Data.(float64) >= val2.Data.(float64)}
	} else {
		runTimeError2(frame, "%s: invalid types", opName)
	}
	return
}

func handleGtOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "gt"
	if l := len(operands); l != 2 {
		runTimeError2(frame, "Wrong amount of arguments for %s (%d given)", opName, l)
	}

	v := operands[0]
	var val1 Value
	switch v.Type {
	case ValueItem:
		val1 = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		val1 = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	v = operands[1]
	var val2 Value
	switch v.Type {
	case ValueItem:
		val2 = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		val2 = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	if val1.Kind == IntValue && val2.Kind == IntValue {
		retVal = Value{Kind: BoolValue, Data: val1.Data.(int) > val2.Data.(int)}
	} else if val1.Kind == FloatValue && val2.Kind == IntValue {
		retVal = Value{Kind: BoolValue, Data: val1.Data.(float64) > float64(val2.Data.(int))}
	} else if val1.Kind == IntValue && val2.Kind == FloatValue {
		retVal = Value{Kind: BoolValue, Data: float64(val1.Data.(int)) > val2.Data.(float64)}
	} else if val1.Kind == FloatValue && val2.Kind == FloatValue {
		retVal = Value{Kind: BoolValue, Data: val1.Data.(float64) > val2.Data.(float64)}
	} else {
		runTimeError2(frame, "%s: invalid types", opName)
	}
	return
}

func handleSplitOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "split"
	argCount := len(operands)
	if l := argCount; (l != 2) && (l != 1) {
		runTimeError2(frame, "Wrong amount of arguments for %s (%d given)", opName, l)
	}

	v := operands[0]
	var strval Value
	switch v.Type {
	case ValueItem:
		strval = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		strval = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	var substr Value
	if argCount > 1 {
		v = operands[1]
		switch v.Type {
		case ValueItem:
			substr = v.Data.(Value)
		case SymbolPathItem, OperCallItem:
			substr = EvalItem(v, frame)
		default:
			runTimeError2(frame, "something wrong (%s)", opName)
		}
		if substr.Kind != StringValue {
			runTimeError2(frame, "%s: string value assumed for 2nd argument", opName)
		}
	}

	if strval.Kind != StringValue {
		runTimeError2(frame, "%s: string value assumed for 1st argument", opName)
	}

	var parts []string
	switch argCount {
	case 1:
		parts = strings.Fields(strval.Data.(string))
	case 2:
		parts = strings.Split(strval.Data.(string), substr.Data.(string))
	}

	var partsList List
	var prevnewitem *ListObject
	var newitem *ListObject
	for _, v := range parts {
		newitem = &ListObject{Val: &Value{Kind: StringValue, Data: v}, Next: nil}
		if partsList.Head == nil {
			partsList.Head = newitem
		}
		if prevnewitem != nil {
			prevnewitem.Next = newitem
		}
		prevnewitem = newitem
	}
	retVal = Value{Kind: ListValue, Data: &partsList}
	return
}

func handleSliceOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "slice"
	argCount := len(operands)
	if l := argCount; (l != 2) && (l != 3) {
		runTimeError2(frame, "Wrong amount of arguments for %s (%d given)", opName, l)
	}

	v := operands[0]
	var seqval Value
	switch v.Type {
	case ValueItem:
		seqval = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		seqval = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	v = operands[1]
	var startIndVal Value
	switch v.Type {
	case ValueItem:
		startIndVal = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		startIndVal = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	if startIndVal.Kind != IntValue {
		runTimeError2(frame, "%s: int value assumed for 2nd argument", opName)
	}

	var endIndVal Value
	if argCount > 2 {
		v = operands[2]
		switch v.Type {
		case ValueItem:
			endIndVal = v.Data.(Value)
		case SymbolPathItem, OperCallItem:
			endIndVal = EvalItem(v, frame)
		default:
			runTimeError2(frame, "something wrong (%s)", opName)
		}
		if endIndVal.Kind != IntValue {
			runTimeError2(frame, "%s: int value assumed for 2nd argument", opName)
		}
		if ind1, ind2 := endIndVal.Data.(int), startIndVal.Data.(int); ind1 < ind2 {
			runTimeError2(frame, "%s: 3rd index (%d) should not be less than 2nd one (%d)", opName, ind1, ind2)
		}
	} else {
		endIndVal = Value{Kind: IntValue, Data: -1}
	}

	switch seqval.Kind {
	case ListValue:
		it := NewListIterator(seqval)
		count := 0
		ind1, ind2 := startIndVal.Data.(int), endIndVal.Data.(int)
		var items []*Item
		for {
			nextItemval := it.Next()
			if nextItemval == nil {
				break
			}
			if argCount == 2 {
				if count >= ind1 {
					items = append(items, &Item{Type: ValueItem, Data: *nextItemval})
				}
			} else {
				if (count >= ind1) && (count <= ind2) {
					items = append(items, &Item{Type: ValueItem, Data: *nextItemval})
				}
			}
			count++
		}
		retVal = handleListOP(frame, items)

	case StringValue:
		s := seqval.Data.(string)
		from := startIndVal.Data.(int)
		slen := len(s)
		// TODO: check also lower bounds...
		if from >= slen {
			runTimeError2(frame, "%s: Index out of range (2nd: %d)", opName, from)
		}
		if argCount == 2 {
			retVal = Value{Kind: StringValue, Data: string([]byte(s)[from:])}
		} else {
			to := endIndVal.Data.(int)
			if to >= slen {
				runTimeError2(frame, "%s: Index out of range (3rd: %d)", opName, to)
			}
			retVal = Value{Kind: StringValue, Data: string([]byte(s)[from : to+1])}
		}
	default:
		runTimeError2(frame, "Not suitable type for %s operator", opName)
	}

	return
}

func handleFindOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "find"
	if l := len(operands); l != 2 {
		runTimeError2(frame, "Wrong amount of arguments for %s (%d given)", opName, l)
	}

	v := operands[0]
	var seqval Value
	switch v.Type {
	case ValueItem:
		seqval = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		seqval = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}
	v = operands[1]
	var itemval Value
	switch v.Type {
	case ValueItem:
		itemval = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		itemval = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	switch seqval.Kind {
	case ListValue:
		it := NewListIterator(seqval)
		itemAsVal := &Item{Type: ValueItem, Data: itemval}
		var matches []*Item
		count := 0
		for {
			nextItemval := it.Next()
			if nextItemval == nil {
				break
			}
			argsForEq := []*Item{
				&Item{
					Type: ValueItem,
					Data: *nextItemval,
				},
				itemAsVal,
			}
			eqResult := handleEqOP(frame, argsForEq)
			if eqResult.Kind != BoolValue {
				runTimeError2(frame, "Invalid result from eq")
			}
			if eqResult.Data.(bool) == true {
				matches = append(matches, &Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: count}})
			}
			count++
		}
		retVal = handleListOP(frame, matches)

	case StringValue:
		if itemval.Kind != StringValue {
			runTimeError2(frame, "Unsupported type for %s", opName)
		}
		s := seqval.Data.(string)
		slen := len(s)
		substr := itemval.Data.(string)
		sublen := len(substr)
		sloc := 0
		realoc := 0
		var litems []*Item
		if sublen == 0 {
			retVal = Value{Kind: ListValue, Data: &List{}}
			return
		}
		for {
			loc := strings.Index(s, substr)
			if loc == -1 {
				break
			}
			sloc = loc + sublen
			litems = append(litems, &Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: realoc + loc}})
			realoc += sloc
			if sloc >= slen {
				break
			}
			s = string([]byte(s)[sloc:])
		}
		retVal = handleListOP(frame, litems)
	default:
		runTimeError2(frame, "Not suitable type for %s operator", opName)
	}

	return
}

func handleIndOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "ind"
	if l := len(operands); l != 2 {
		runTimeError2(frame, "Wrong amount of arguments for %s (%d given)", opName, l)
	}

	v := operands[0]
	var seqval Value
	switch v.Type {
	case ValueItem:
		seqval = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		seqval = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}
	v = operands[1]
	var itemval Value
	switch v.Type {
	case ValueItem:
		itemval = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		itemval = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	if itemval.Kind != IntValue {
		runTimeError2(frame, "Unsupported type for %s", opName)
	}

	switch seqval.Kind {
	case ListValue:
		it := NewListIterator(seqval)
		loc := itemval.Data.(int)
		count := 0
		for {
			nextItemval := it.Next()
			if nextItemval == nil {
				runTimeError2(frame, "%s: index out of range (len=%d)(accessed=%d)", opName, count, loc)
			}
			if loc == count {
				retVal = *nextItemval
				return
			}
			count++
		}

	case StringValue:
		s := seqval.Data.(string)
		loc := itemval.Data.(int)
		if len(s) <= loc {
			runTimeError2(frame, "%s: index out of range (len=%d)(accessed=%d)", opName, len(s), loc)
		}
		retVal = Value{Kind: StringValue, Data: string([]byte(s)[loc])}
	default:
		runTimeError2(frame, "Not suitable type for %s operator", opName)
	}

	return
}

func handleInOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "in"
	if l := len(operands); l != 2 {
		runTimeError2(frame, "Wrong amount of arguments for %s (%d given)", opName, l)
	}

	v := operands[0]
	var seqval Value
	switch v.Type {
	case ValueItem:
		seqval = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		seqval = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}
	v = operands[1]
	var itemval Value
	switch v.Type {
	case ValueItem:
		itemval = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		itemval = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	switch seqval.Kind {
	case ListValue:
		it := NewListIterator(seqval)
		itemAsVal := &Item{Type: ValueItem, Data: itemval}
		for {
			nextItemval := it.Next()
			if nextItemval == nil {
				break
			}
			argsForEq := []*Item{
				&Item{
					Type: ValueItem,
					Data: *nextItemval,
				},
				itemAsVal,
			}
			eqResult := handleEqOP(frame, argsForEq)
			if eqResult.Kind != BoolValue {
				runTimeError2(frame, "Invalid result from eq")
			}
			if eqResult.Data.(bool) == true {
				retVal = Value{Kind: BoolValue, Data: true}
				return
			}
		}
		retVal = Value{Kind: BoolValue, Data: false}

	case MapValue:
		argsForGetl := []*Item{
			&Item{Type: ValueItem, Data: seqval},
			&Item{Type: ValueItem, Data: itemval},
		}
		lval := handleGetlOP(frame, argsForGetl)
		if lval.Kind != ListValue {
			runTimeError2(frame, "%s: expecting list", opName)
		}
		headVal := handleHeadOP(frame, []*Item{&Item{Type: ValueItem, Data: lval}})
		if headVal.Kind != BoolValue {
			runTimeError2(frame, "%s: expecting bool as first in list", opName)
		}
		retVal = headVal

	case StringValue:
		if itemval.Kind != StringValue {
			runTimeError2(frame, "Unsupported type for %s", opName)
		}
		s := seqval.Data.(string)
		substr := itemval.Data.(string)
		isIn := strings.Contains(s, substr)
		retVal = Value{Kind: BoolValue, Data: isIn}
	default:
		runTimeError2(frame, "Not suitable type for %s operator", opName)
	}
	return
}

func handleTypeOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "type"
	if l := len(operands); l != 1 {
		runTimeError2(frame, "%s operator needs one argument (%d given)", opName, l)
	}

	v := operands[0]
	var val Value
	switch v.Type {
	case ValueItem:
		val = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		val = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	var typeName string
	switch val.Kind {
	case IntValue:
		typeName = "int"
	case FloatValue:
		typeName = "float"
	case BoolValue:
		typeName = "bool"
	case StringValue:
		typeName = "string"
	case FuncProtoValue:
		typeName = "funcproto"
	case FunctionValue:
		typeName = "function"
	case ListValue:
		typeName = "list"
	case ChanValue:
		typeName = "channel"
	case OpaqueValue:
		if !frame.inProcCall {
			runTimeError2(frame, "%s not allowed in function for opaque value", opName)
		}
		typeName = "opaque:" + val.Data.(OpaqueAPI).TypeName()
	case MapValue:
		typeName = "map"
	case ExtProcValue:
		typeName = "ext-proc"
	default:
		runTimeError2(frame, "%s: unknown type (%d)", opName, val.Kind)
	}
	retVal = Value{Kind: StringValue, Data: typeName}
	return
}

func handleLenOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "len"
	if l := len(operands); l != 1 {
		runTimeError2(frame, "%s operator needs one argument (%d given)", opName, l)
	}

	v := operands[0]
	var val Value
	switch v.Type {
	case ValueItem:
		val = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		val = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	var length int
	switch val.Kind {
	case ListValue:
		list, convok := val.Data.(*List)
		if !convok {
			runTimeError2(frame, "First argument is not list in %s operator", opName)
		}
		if list.Head != nil {
			length++
			cur := list.Head
			for {
				if cur.Next == nil {
					break
				}
				cur = cur.Next
				length++
			}
		}
		if list.Tail != nil {
			length++
			cur := list.Tail
			for {
				if cur.Next == nil {
					break
				}
				cur = cur.Next
				length++
			}
		}

	case MapValue:
		mapv, convok := val.Data.(*PMap)
		if !convok {
			runTimeError2(frame, "%s: inconsistency in map data", opName)
		}
		length = mapv.itemCount

	case StringValue:
		length = len(val.Data.(string))

	default:
		runTimeError2(frame, "Invalid type for %s", opName)
	}
	retVal = Value{Kind: IntValue, Data: length}
	return
}

func handleAndOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "and"
	if l := len(operands); l < 2 {
		runTimeError2(frame, "Not enough arguments for %s (%d given)", opName, l)
	}
	for _, v := range operands {
		var argval Value
		switch v.Type {
		case ValueItem:
			argval = v.Data.(Value)
		case SymbolPathItem, OperCallItem:
			argval = EvalItem(v, frame)
		default:
			runTimeError2(frame, "something wrong")
		}
		if argval.Kind != BoolValue {
			runTimeError2(frame, "Invalid type for %s", opName)
		}
		if !argval.Data.(bool) {
			retVal = Value{Kind: BoolValue, Data: false}
			return
		}
	}
	retVal = Value{Kind: BoolValue, Data: true}
	return
}

func handleOrOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "or"
	if l := len(operands); l < 2 {
		runTimeError2(frame, "Not enough arguments for %s (%d given)", opName, l)
	}
	for _, v := range operands {
		var argval Value
		switch v.Type {
		case ValueItem:
			argval = v.Data.(Value)
		case SymbolPathItem, OperCallItem:
			argval = EvalItem(v, frame)
		default:
			runTimeError2(frame, "something wrong")
		}
		if argval.Kind != BoolValue {
			runTimeError2(frame, "Invalid type for %s", opName)
		}
		if argval.Data.(bool) {
			retVal = Value{Kind: BoolValue, Data: true}
			return
		}
	}
	retVal = Value{Kind: BoolValue, Data: false}
	return
}

func handleNotOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "not"
	if l := len(operands); l != 1 {
		runTimeError2(frame, "Wrong amount of arguments for %s (%d given)", opName, l)
	}
	var argval Value
	switch cond := operands[0]; cond.Type {
	case ValueItem:
		argval = cond.Data.(Value)
	case SymbolPathItem, OperCallItem:
		argval = EvalItem(cond, frame)
	default:
		runTimeError2(frame, "something wrong")
	}
	if argval.Kind != BoolValue {
		runTimeError2(frame, "%s: condition should be boolean expression", opName)
	}
	retVal = Value{Kind: BoolValue, Data: !argval.Data.(bool)}
	return
}

func handleIfOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "if"
	if l := len(operands); l != 3 {
		runTimeError2(frame, "Wrong amount of arguments for %s (%d given)", opName, l)
	}
	var argval Value
	switch cond := operands[0]; cond.Type {
	case ValueItem:
		argval = cond.Data.(Value)
	case SymbolPathItem, OperCallItem:
		argval = EvalItem(cond, frame)
	default:
		runTimeError2(frame, "something wrong")
	}
	if argval.Kind != BoolValue {
		runTimeError2(frame, "%s: condition should be boolean expression", opName)
	}
	if argval.Data.(bool) {
		retVal = EvalItem(operands[1], frame)
	} else {
		retVal = EvalItem(operands[2], frame)
	}
	return
}

func handleEqOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "eq"
	var argType ValueType
	var comparedValue interface{}

	if l := len(operands); l < 2 {
		runTimeError2(frame, "Not enough arguments for %s (%d given)", opName, l)
	}
	for i, v := range operands {
		var argval Value
		switch v.Type {
		case ValueItem:
			argval = v.Data.(Value)
		case SymbolPathItem, OperCallItem:
			argval = EvalItem(v, frame)
		default:
			runTimeError2(frame, "something wrong")
		}
		if i == 0 {
			argType = argval.Kind
			comparedValue = argval.Data
		}
		if argType != argval.Kind {
			retVal = Value{Kind: BoolValue, Data: false}
			return
		}
		switch argval.Kind {
		case OpaqueValue:
			if !comparedValue.(OpaqueAPI).Equals(argval.Data.(OpaqueAPI)) {
				retVal = Value{Kind: BoolValue, Data: false}
				return
			}
		case IntValue:
			if comparedValue.(int) != argval.Data.(int) {
				retVal = Value{Kind: BoolValue, Data: false}
				return
			}
		case FloatValue:
			if comparedValue.(float64) != argval.Data.(float64) {
				retVal = Value{Kind: BoolValue, Data: false}
				return
			}
		case StringValue:
			if comparedValue.(string) != argval.Data.(string) {
				retVal = Value{Kind: BoolValue, Data: false}
				return
			}
		case BoolValue:
			if comparedValue.(bool) != argval.Data.(bool) {
				retVal = Value{Kind: BoolValue, Data: false}
				return
			}
		case MapValue:
			m1, convok := comparedValue.(*PMap)
			if !convok {
				runTimeError2(frame, "%s: invalid map", opName)
			}
			m2, convok := argval.Data.(*PMap)
			if !convok {
				runTimeError2(frame, "%s: invalid map", opName)
			}
			if !areEqualMaps(frame, m1, m2) {
				retVal = Value{Kind: BoolValue, Data: false}
				return
			}
		case ListValue:
			listv1 := Value{Kind: ListValue, Data: comparedValue.(*List)}
			listv2 := Value{Kind: ListValue, Data: argval.Data.(*List)}
			it1 := NewListIterator(listv1)
			it2 := NewListIterator(listv2)
			for {
				itemval1 := it1.Next()
				itemval2 := it2.Next()
				if (itemval1 == nil) && (itemval2 == nil) {
					break
				} else if (itemval1 != nil) && (itemval2 == nil) {
					retVal = Value{Kind: BoolValue, Data: false}
					return
				} else if (itemval1 == nil) && (itemval2 != nil) {
					retVal = Value{Kind: BoolValue, Data: false}
					return
				}
				argsForEq := []*Item{
					&Item{
						Type: ValueItem,
						Data: *itemval1,
					},
					&Item{
						Type: ValueItem,
						Data: *itemval2,
					},
				}
				eqResult := handleEqOP(frame, argsForEq)
				if eqResult.Kind != BoolValue {
					runTimeError2(frame, "Invalid result from eq")
				}
				if eqResult.Data.(bool) == false {
					retVal = Value{Kind: BoolValue, Data: false}
					return
				}
			}
		case FunctionValue:
			runTimeError2(frame, "%s: ...to be implemented", opName)
		default:
			runTimeError2(frame, "Invalid type for %s", opName)
		}
	}
	retVal = Value{Kind: BoolValue, Data: true}
	return
}

func handleMinusOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "minus"

	if l := len(operands); l != 2 {
		runTimeError2(frame, "%s operator should have exactly two arguments (%d given)", opName, l)
	}

	var argval1 Value
	var argval2 Value

	switch arg1 := operands[0]; arg1.Type {
	case ValueItem:
		argval1 = arg1.Data.(Value)
	case SymbolPathItem, OperCallItem:
		argval1 = EvalItem(arg1, frame)
	default:
		runTimeError2(frame, "something wrong")
	}

	switch arg2 := operands[1]; arg2.Type {
	case ValueItem:
		argval2 = arg2.Data.(Value)
	case SymbolPathItem, OperCallItem:
		argval2 = EvalItem(arg2, frame)
	default:
		runTimeError2(frame, "something wrong")
	}

	if argval1.Kind == IntValue && argval2.Kind == IntValue {
		num1 := argval1.Data.(int)
		num2 := argval2.Data.(int)
		retVal = Value{Kind: IntValue, Data: num1 - num2}
	} else if argval1.Kind == FloatValue && argval2.Kind == FloatValue {
		num1 := argval1.Data.(float64)
		num2 := argval2.Data.(float64)
		retVal = Value{Kind: FloatValue, Data: num1 - num2}
	} else {
		runTimeError2(frame, "Invalid types for %s", opName)
	}
	return
}

func handleDivOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "div"

	if l := len(operands); l != 2 {
		runTimeError2(frame, "%s operator should have exactly two arguments (%d given)", opName, l)
	}

	var dividend Value
	switch operands[0].Type {
	case ValueItem:
		dividend = operands[0].Data.(Value)
	case SymbolPathItem, OperCallItem:
		dividend = EvalItem(operands[0], frame)
	default:
		runTimeError2(frame, "something wrong")
	}

	var divisor Value
	switch operands[1].Type {
	case ValueItem:
		divisor = operands[1].Data.(Value)
	case SymbolPathItem, OperCallItem:
		divisor = EvalItem(operands[1], frame)
	default:
		runTimeError2(frame, "something wrong")
	}

	if dividend.Kind == IntValue && divisor.Kind == FloatValue {
		floatDivisor := divisor.Data.(float64)
		if floatDivisor == 0 {
			runTimeError2(frame, "division by zero")
		}
		retVal = Value{Kind: FloatValue, Data: float64(dividend.Data.(int)) / floatDivisor}
	} else if dividend.Kind == FloatValue && divisor.Kind == IntValue {
		intDivisor := divisor.Data.(int)
		if intDivisor == 0 {
			runTimeError2(frame, "division by zero")
		}
		retVal = Value{Kind: FloatValue, Data: dividend.Data.(float64) / float64(intDivisor)}
	} else if dividend.Kind == IntValue && divisor.Kind == IntValue {
		intDivisor := divisor.Data.(int)
		if intDivisor == 0 {
			runTimeError2(frame, "division by zero")
		}
		retVal = Value{Kind: IntValue, Data: dividend.Data.(int) / intDivisor}
	} else if dividend.Kind == FloatValue && divisor.Kind == FloatValue {
		floatDivisor := divisor.Data.(float64)
		if floatDivisor == 0 {
			runTimeError2(frame, "division by zero")
		}
		retVal = Value{Kind: FloatValue, Data: dividend.Data.(float64) / floatDivisor}
	} else {
		runTimeError2(frame, "Invalid type for %s", opName)
	}
	return
}

func handleModOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "mod"
	values := []int{}

	if l := len(operands); l != 2 {
		runTimeError2(frame, "%s operator should have exactly two arguments (%d given)", opName, l)
	}
	for _, v := range operands {
		var argval Value
		switch v.Type {
		case ValueItem:
			argval = v.Data.(Value)
		case SymbolPathItem, OperCallItem:
			argval = EvalItem(v, frame)
		default:
			runTimeError2(frame, "something wrong")
		}
		if argval.Kind != IntValue {
			runTimeError2(frame, "Invalid type for %s", opName)
		}
		values = append(values, argval.Data.(int))
	}
	if values[1] == 0 {
		runTimeError2(frame, "division by zero")
	}
	retVal = Value{Kind: IntValue, Data: values[0] % values[1]}
	return
}

func handleMulOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "mul"
	result := 1
	var resultF float64 = 1
	var isAnyFloat bool

	if l := len(operands); l < 2 {
		runTimeError2(frame, "Not enough arguments for %s (%d given)", opName, l)
	}
	for _, v := range operands {
		var argval Value
		switch v.Type {
		case ValueItem:
			argval = v.Data.(Value)
		case SymbolPathItem, OperCallItem:
			argval = EvalItem(v, frame)
		default:
			runTimeError2(frame, "something wrong")
		}
		switch argval.Kind {
		case IntValue:
			result *= argval.Data.(int)
		case FloatValue:
			isAnyFloat = true
			resultF *= argval.Data.(float64)
		default:
			runTimeError2(frame, "Invalid type for %s", opName)
		}
	}
	if isAnyFloat {
		retVal = Value{Kind: FloatValue, Data: resultF * float64(result)}
	} else {
		retVal = Value{Kind: IntValue, Data: result}
	}
	return
}

func handlePlusOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "plus"
	sum := 0
	var sumF float64 = 0
	text := ""
	var argType ValueType
	hasFloats := false

	if l := len(operands); l < 2 {
		runTimeError2(frame, "Not enough arguments for %s (%d given)", opName, l)
	}
	for i, v := range operands {
		var argval Value
		switch v.Type {
		case ValueItem:
			argval = v.Data.(Value)
		case SymbolPathItem, OperCallItem:
			argval = EvalItem(v, frame)
		default:
			runTimeError2(frame, "something wrong")
		}
		if i == 0 {
			argType = argval.Kind
		}
		if argType != argval.Kind {
			runTimeError2(frame, "mismatching types as arguments")
		}
		switch argval.Kind {
		case IntValue:
			intv := argval.Data.(int)
			sum += intv
		case FloatValue:
			hasFloats = true
			floatv := argval.Data.(float64)
			sumF += floatv
		case StringValue:
			strv := argval.Data.(string)
			text += strv
		default:
			runTimeError2(frame, "Invalid type for %s", opName)
		}
	}
	if hasFloats {
		retVal = Value{Kind: FloatValue, Data: sumF + float64(sum)}
		return
	}
	switch argType {
	case IntValue:
		retVal = Value{Kind: IntValue, Data: sum}
	case StringValue:
		retVal = Value{Kind: StringValue, Data: text}
	default:
		runTimeError2(frame, "Invalid type for %s", opName)
	}
	return
}

func handleEmptyOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "empty"
	if l := len(operands); l != 1 {
		runTimeError2(frame, "%s operator needs one argument (%d given)", opName, l)
	}

	v := operands[0]
	var argval Value
	switch v.Type {
	case ValueItem:
		argval = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		argval = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	switch argval.Kind {
	case ListValue:
		list, convok := argval.Data.(*List)
		if !convok {
			runTimeError2(frame, "First argument is not list in %s operator", opName)
		}
		if list.Head == nil && list.Tail == nil {
			retVal = Value{Kind: BoolValue, Data: true}
		} else {
			retVal = Value{Kind: BoolValue, Data: false}
		}
	case MapValue:
		mapv, convok := argval.Data.(*PMap)
		if !convok {
			runTimeError2(frame, "%s: inconsistency in map data", opName)
		}
		// NOTE. this doesnt work as some items are just marked as deleted
		//retVal = Value{Kind: BoolValue, Data: mapv.Rbm.IsEmpty()}
		retVal = Value{Kind: BoolValue, Data: mapv.itemCount == 0}
	default:
		runTimeError2(frame, "%s: unsupported type", opName)
	}
	return
}
