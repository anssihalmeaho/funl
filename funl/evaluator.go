package funl

import (
	"fmt"
)

type OperHandler func(*Frame, []*Item) Value

var operTbl [MaximumOP]OperHandler

func init() {
	operTbl[AndOP] = handleAndOP
	operTbl[OrOP] = handleOrOP
	operTbl[NotOP] = handleNotOP
	operTbl[CallOP] = handleCallOP
	operTbl[EqOP] = handleEqOP
	operTbl[IfOP] = handleIfOP
	operTbl[PlusOP] = handlePlusOP
	operTbl[MinusOP] = handleMinusOP
	operTbl[MulOP] = handleMulOP
	operTbl[DivOP] = handleDivOP
	operTbl[ModOP] = handleModOP
	operTbl[ListOP] = handleListOP
	operTbl[EmptyOP] = handleEmptyOP
	operTbl[HeadOP] = handleHeadOP
	operTbl[LastOP] = handleLastOP
	operTbl[RestOP] = handleRestOP
	operTbl[AppendOP] = handleAppendOP
	operTbl[AddOP] = handleAddOP
	operTbl[LenOP] = handleLenOP
	operTbl[TypeOP] = handleTypeOP
	operTbl[InOP] = handleInOP
	operTbl[IndOP] = handleIndOP
	operTbl[FindOP] = handleFindOP
	operTbl[SliceOP] = handleSliceOP
	operTbl[RrestOP] = handleRrestOP
	operTbl[ReverseOP] = handleReverseOP
	operTbl[ExtendOP] = handleExtendOP
	operTbl[SplitOP] = handleSplitOP
	operTbl[GtOP] = handleGtOP
	operTbl[LtOP] = handleLtOP
	operTbl[LeOP] = handleLeOP
	operTbl[GeOP] = handleGeOP
	operTbl[StrOP] = handleStrOP
	operTbl[ConvOP] = handleConvOP
	operTbl[CaseOP] = handleCaseOP
	operTbl[NameOP] = handleNameOP
	operTbl[ErrorOP] = handleErrorOP
	operTbl[PrintOP] = handlePrintOP
	operTbl[SpawnOP] = handleSpawnOP
	operTbl[ChanOP] = handleChanOP
	operTbl[SendOP] = handleSendOP
	operTbl[RecvOP] = handleRecvOP
	operTbl[SymvalOP] = handleSymvalOP
	operTbl[TryOP] = handleTryOP
	operTbl[SelectOP] = handleSelectOP
	operTbl[EvalOP] = handleEvalOP
	operTbl[WhileOP] = handleWhileOP
	operTbl[FloatOP] = handleFloatOP
	operTbl[MapOP] = handleMapOP
	operTbl[PutOP] = handlePutOP
	operTbl[GetOP] = handleGetOP
	operTbl[GetlOP] = handleGetlOP
	operTbl[KeysOP] = handleKeysOP
	operTbl[ValsOP] = handleValsOP
	operTbl[KeyvalsOP] = handleKeyvalsOP
	operTbl[LetOP] = handleLetOP
	operTbl[ImpOP] = handleImpOP
	operTbl[DelOP] = handleDelOP
	operTbl[DellOP] = handleDellOP
	operTbl[SprintfOP] = handleSprintfOP
	operTbl[ArgslistOP] = handleArgslistOP
	operTbl[CondOP] = handleCondOP
	operTbl[HelpOP] = handleHelpOP
	operTbl[RecwithOP] = handleRecwithOP
}

func RunTimeError(format string, args ...interface{}) {
	runTimeError2(nil, format, args...)
}

func runTimeError(format string, args ...interface{}) {
	runTimeError2(nil, format, args...)
}

func RunTimeError2(frame *Frame, format string, args ...interface{}) {
	runTimeError2(frame, format, args...)
}

var PrintingRTElocationAndScopeEnabled bool

func runTimeError2(frame *Frame, format string, args ...interface{}) {
	if false {
		fmt.Printf("\nruntime error: "+format+"\n", args...)
	}
	if PrintingRTElocationAndScopeEnabled {
		if frame != nil {
			fmt.Printf("call scope (RTE):\n")
			fDebugData := frame.GetFuncDebugInfos([]fdebugInfo{})
			idx := 0
			for i := len(fDebugData) - 1; i >= 0; i-- {
				finfo := fDebugData[i].function
				argvalues := fDebugData[i].argvalues
				idx++
				if finfo != nil {
					fmt.Printf("  %d: File: %s, Line %d, Pos: %d\n", idx, finfo.SrcFileName, finfo.Lineno, finfo.Pos)
					fmt.Printf("       args:\n")
					for argid, argvalue := range argvalues {
						fmt.Printf("       %d. %s\n", argid+1, argvalue)
					}
					fmt.Println()
					fmt.Printf("       symbol values:\n")
					for k, v := range fDebugData[i].syms.AsMap() {
						valAsStr := "-"
						if v.Type == ValueItem {
							valAsStr = fmt.Sprintf("%s", v.Data.(Value))
						}
						fmt.Printf("       %s: %s\n", SymIDMap.AsString(k), valAsStr)
					}
				} else {
					fmt.Printf("  %d: -\n", idx)
				}
			}
		}
	}
	panic(fmt.Errorf(format, args...))
}

type Frame struct {
	FuncProto     *Function
	Syms          *Symt
	OtherNS       map[SymID]ImportInfo
	AccessLink    *Frame // nil if root
	Imported      map[SymID]*Frame
	inProcCall    bool
	EvaluatedArgs []Value
	Interpreter   *Interpreter
}

type fdebugInfo struct {
	function  *Function
	argvalues []Value
	syms      *Symt
}

// GetTopFrame gets top frame
func (fr *Frame) GetTopFrame() *Frame {
	if fr.AccessLink == nil {
		return fr
	}
	return fr.AccessLink.GetTopFrame()
}

// SetInProcCall sets inProcCall in Frame
func (fr *Frame) SetInProcCall(v bool) {
	fr.inProcCall = v
}

// GetFuncDebugInfos gets function infos for backtrace
func (fr *Frame) GetFuncDebugInfos(prev []fdebugInfo) []fdebugInfo {
	fdeb := fdebugInfo{
		function:  fr.FuncProto,
		argvalues: fr.EvaluatedArgs,
		syms:      fr.Syms,
	}
	if fr.AccessLink == nil {
		return append(prev, fdeb)
	}
	return fr.AccessLink.GetFuncDebugInfos(append(prev, fdeb))
}

//GetSymItemsOfImportedModule finds symbol table for module
func (fr *Frame) GetSymItemsOfImportedModule(modSid SymID) (*Symt, bool) {
	frame, found := fr.Imported[modSid]
	if !found {
		if fr.AccessLink != nil {
			return fr.AccessLink.GetSymItemsOfImportedModule(modSid)
		}
		return nil, false
	}
	return frame.Syms, true
}

// GetImportedSymItem gets item from imported namespaces
func (fr *Frame) GetImportedSymItem(modSid SymID, rest SymbolPath) (*Item, bool) {
	frame, found := fr.Imported[modSid]
	if !found {
		if fr.AccessLink != nil {
			return fr.AccessLink.GetImportedSymItem(modSid, rest)
		}
		return nil, false
	}
	if l := len(rest); l == 1 {
		return frame.GetSymItem(rest[0])
	} else if l > 1 {
		return frame.GetImportedSymItem(rest[0], rest[1:])
	}
	runTimeError2(frame, "invalid symbol path")
	return nil, false // should never get here
}

// GetSymItem gets item ralated to symbol
func (fr *Frame) GetSymItem(sid SymID) (*Item, bool) {
	item, found := fr.Syms.GetBySID(sid)
	if found {
		return item, true
	}
	if fr.AccessLink != nil {
		return fr.AccessLink.GetSymItem(sid)
	}
	return nil, false
}

func (fr *Frame) FindFuncSID(fptr *Function) (sid SymID, found bool) {
	sid, found = fr.Syms.FindFuncSID(fptr)
	if found {
		return
	}
	if fr.AccessLink != nil {
		return fr.AccessLink.FindFuncSID(fptr)
	}
	return
}

func evalMain(frame *Frame) Value {
	item := frame.FuncProto.Body
	return EvalItemV2(item, frame, &AddInfo{evaluatingBody: true})
}

func handleWhileOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "while"
	argsCount := len(operands)
	if l := argsCount; l < 2 {
		runTimeError2(frame, "Wrong amount of arguments for %s (%d given)", opName, l)
	}

	// create new frame, which is copy of previous one
	nFrame := Frame{}
	nFrame.FuncProto = frame.FuncProto
	nFrame.AccessLink = frame.AccessLink
	nFrame.Syms = frame.Syms.MakeCopy()
	nFrame.inProcCall = frame.inProcCall
	nFrame.Imported = frame.Imported
	// nFrame.Interpreter = frame.Interpreter , only in top frame

	nextFrame := &nFrame
	for {
		var argval Value
		switch cond := operands[0]; cond.Type {
		case ValueItem:
			argval = cond.Data.(Value)
		case SymbolPathItem, OperCallItem:
			argval = EvalItem(cond, nextFrame)
		default:
			runTimeError2(frame, "something wrong")
		}
		if argval.Kind != BoolValue {
			runTimeError2(frame, "%s: condition should be boolean expression", opName)
		}

		// if condition is false, then return result
		if !argval.Data.(bool) {
			retVal = EvalItem(operands[argsCount-1], nextFrame)
			return
		}

		// lets first evaluate arguments
		var evaluatedArgs []Value
		for _, argitem := range operands[1 : argsCount-1] {
			if argitem.Type == ValueItem {
				evaluatedArgs = append(evaluatedArgs, argitem.Data.(Value))
			} else {
				evaluatedArgs = append(evaluatedArgs, EvalItem(argitem, nextFrame))
			}
		}

		// then add arguments to symbol table
		if len(evaluatedArgs) < len(nextFrame.FuncProto.ArgNames) {
			runTimeError2(frame, "Not enough arguments for function call, need: %d, got:%d", len(evaluatedArgs), len(nextFrame.FuncProto.ArgNames))
		}
		for index, sid := range frame.FuncProto.ArgNames {
			if (index + 0) >= len(evaluatedArgs) {
				runTimeError2(frame, "Too short argument list in function call")
			}
			argitem := &Item{Type: ValueItem, Data: evaluatedArgs[index+0]}
			if !nextFrame.Syms.AddBySIDByOverwriteIfNeeded(sid, argitem) {
				runTimeError2(frame, "Argument add failed")
			}
			index++
		}
		// put args also to separate slice so that those can be accessed by argslist -operator
		nextFrame.EvaluatedArgs = evaluatedArgs

		// NOTE. Imports should remain

		// lets fill let defintions to frame symbol table
		symbolMap := nextFrame.FuncProto.NSpace.Syms.AsMap()
		for _, sid := range nextFrame.FuncProto.NSpace.Syms.Keys() {
			letitem := symbolMap[sid]
			letvalitem := &Item{Type: ValueItem, Data: EvalItem(letitem, nextFrame)}
			if !nextFrame.Syms.AddBySIDByOverwriteIfNeeded(sid, letvalitem) {
				runTimeError2(frame, "Symbol add failed")
			}
		}
	}
}

//HandleCallOP for std lib usage
func HandleCallOP(frame *Frame, operands []*Item) (retVal Value) {
	return handleCallOP(frame, operands)
}

func handleCallOP(frame *Frame, operands []*Item) (retVal Value) {
	if l := len(operands); l == 0 {
		runTimeError2(frame, "Empty argument list for call operator")
	}
	// lets first evaluate arguments
	var evaluatedArgs []Value
	for _, argitem := range operands {
		if argitem.Type == ValueItem {
			evaluatedArgs = append(evaluatedArgs, argitem.Data.(Value))
		} else {
			evaluatedArgs = append(evaluatedArgs, EvalItem(argitem, frame))
		}
	}
	// create new frame
	nextFrame := Frame{
		Syms:     NewSymt(),
		OtherNS:  nil, // TODO: needs to be something more...
		Imported: make(map[SymID]*Frame),
		// Interpreter: frame.Interpreter, only in top frame
	}
	isExtProcCall := false
	// lets take function from first argument
	funcitem := evaluatedArgs[0]
	switch funcitem.Kind {
	case FuncProtoValue:
		nextFrame.FuncProto = funcitem.Data.(*Function)
		nextFrame.AccessLink = frame
	case FunctionValue:
		nextFrame.FuncProto = funcitem.Data.(FuncValue).FuncProto
		nextFrame.AccessLink = funcitem.Data.(FuncValue).AccessLink
	case ExtProcValue:
		isExtProcCall = true
		nextFrame.FuncProto = nil
		nextFrame.AccessLink = frame
	default:
		runTimeError2(frame, "First argument for call is not function (%v)", funcitem)
	}

	if frame.inProcCall {
		if isExtProcCall {
			nextFrame.inProcCall = true
		} else if nextFrame.FuncProto.IsProc {
			nextFrame.inProcCall = true
		} else {
			nextFrame.inProcCall = false
		}
	} else {
		if isExtProcCall {
			if !funcitem.Data.(ExtProcType).IsFunction {
				runTimeError2(frame, "external proc call not allowed from function")
			}
		} else if nextFrame.FuncProto.IsProc {
			runTimeError2(frame, "proc call not allowed from func")
		}
	}

	if isExtProcCall {
		extp := funcitem.Data.(ExtProcType)
		retVal = extp.Impl(frame, evaluatedArgs[1:])
		return
	}

	// then add arguments to symbol table
	if len(evaluatedArgs) < len(nextFrame.FuncProto.ArgNames) {
		runTimeError2(frame, "Not enough arguments for function call, need: %d, got:%d", len(evaluatedArgs), len(nextFrame.FuncProto.ArgNames))
	}
	for index, sid := range nextFrame.FuncProto.ArgNames {
		if (index + 1) >= len(evaluatedArgs) {
			runTimeError2(frame, "Too short argument list in function call")
		}
		argitem := &Item{Type: ValueItem, Data: evaluatedArgs[index+1]}

		// Overlapping symbol names not allowed (might cause variable -like effect)
		if _, symfound := nextFrame.GetSymItem(sid); symfound {
			runTimeError2(frame, "Duplicate symbol name in scope not allowed (%s)", SymIDMap.AsString(sid))
		}

		if !nextFrame.Syms.AddBySID(sid, argitem) {
			runTimeError2(frame, "Argument add failed")
		}
		index++
	}
	// put args also to separate slice so that those can be accessed by argslist -operator
	nextFrame.EvaluatedArgs = evaluatedArgs[1:]
	// lets handle local imports
	interpreter := frame.GetTopFrame().Interpreter
	AddImportsToNamespace(&nextFrame.FuncProto.NSpace, &nextFrame, interpreter)

	// lets fill let defintions to frame symbol table
	symbolMap := nextFrame.FuncProto.NSpace.Syms.AsMap()
	for _, sid := range nextFrame.FuncProto.NSpace.Syms.Keys() {
		letitem := symbolMap[sid]

		// Overlapping symbol names not allowed (might cause variable -like effect)
		if _, symfound := nextFrame.GetSymItem(sid); symfound {
			runTimeError2(frame, "Duplicate symbol name in scope not allowed (%s)", SymIDMap.AsString(sid))
		}

		// lets evaluate let def. to value first, using new frame already (as arguments are there)
		letvalitem := &Item{Type: ValueItem, Data: EvalItem(letitem, &nextFrame)}
		if !nextFrame.Syms.AddBySID(sid, letvalitem) {
			runTimeError2(frame, "Symbol add failed")
		}
	}
	retVal = EvalItemV2(nextFrame.FuncProto.Body, &nextFrame, &AddInfo{evaluatingBody: true})
	return
}

type AddInfo struct {
	evaluatingBody bool
}

// TODO: its copy-paste, but no wrapper because it would add additional overhead
func EvalItemV2(item *Item, frame *Frame, addInfo *AddInfo) (retVal Value) {
	switch item.Type {
	case ValueItem:
		var ok bool
		retVal, ok = item.Data.(Value)
		if !ok {
			runTimeError2(frame, "Data corrupted")
		}
		// lets form function value here
		if retVal.Kind == FuncProtoValue {
			fv := FuncValue{FuncProto: retVal.Data.(*Function), AccessLink: frame}
			retVal = Value{Kind: FunctionValue, Data: fv}
		}
	case SymbolPathItem:
		sp := item.Data.(SymbolPath)
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
		return EvalItem(symItem, frame)
	case OperCallItem:
		opcall, ok := item.Data.(OpCall)
		if !ok {
			runTimeError2(frame, "Data corrupted")
		}

		// special check that while-operator is used only in top of function body
		if opcall.OperID == WhileOP {
			if addInfo != nil && !addInfo.evaluatingBody {
				runTimeError2(frame, "while -operator should be used only as immediate operator in func/proc body")
			}
		}

		// lets change function-prototypes to function values here
		var newOperands []*Item
		for _, v := range opcall.Operands {
			var argval Value
			switch v.Type {
			case ValueItem:
				argval = v.Data.(Value)
				if argval.Kind == FuncProtoValue {
					fv := FuncValue{FuncProto: argval.Data.(*Function), AccessLink: frame}
					argval = Value{Kind: FunctionValue, Data: fv}
					newOperands = append(newOperands, &Item{Type: ValueItem, Data: argval})
				} else {
					newOperands = append(newOperands, v)
				}
			default:
				newOperands = append(newOperands, v)
			}
		}

		if !item.Expand {
			return operTbl[opcall.OperID](frame, newOperands)
		}
		var expandedOperands []*Item
		for idx, v := range newOperands {
			if !item.ExpandArgIndexes[idx] {
				expandedOperands = append(expandedOperands, v)
				continue
			}
			listVal := EvalItem(v, frame)
			if listVal.Kind != ListValue {
				runTimeError2(frame, "Expansion requires list value")
			}
			lit := NewListIterator(listVal)
			for {
				nextv := lit.Next()
				if nextv == nil {
					break
				}
				expandedOperands = append(expandedOperands, &Item{Type: ValueItem, Data: *nextv})
			}
		}
		return operTbl[opcall.OperID](frame, expandedOperands)
	default:
		runTimeError2(frame, "Data corrupted")
	}
	return
}

func EvalItem(item *Item, frame *Frame) (retVal Value) {
	switch item.Type {
	case ValueItem:
		var ok bool
		retVal, ok = item.Data.(Value)
		if !ok {
			runTimeError2(frame, "Data corrupted")
		}
		// lets form function value here
		if retVal.Kind == FuncProtoValue {
			fv := FuncValue{FuncProto: retVal.Data.(*Function), AccessLink: frame}
			retVal = Value{Kind: FunctionValue, Data: fv}
		}
	case SymbolPathItem:
		sp := item.Data.(SymbolPath)
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
		return EvalItem(symItem, frame)
	case OperCallItem:
		opcall, ok := item.Data.(OpCall)
		if !ok {
			runTimeError2(frame, "Data corrupted")
		}

		// special check that while-operator is used only in top of function body
		if opcall.OperID == WhileOP {
			runTimeError2(frame, "while -operator should be used only as immediate operator in func/proc body")
		}

		// lets change function-prototypes to function values here
		var newOperands []*Item
		for _, v := range opcall.Operands {
			var argval Value
			switch v.Type {
			case ValueItem:
				argval = v.Data.(Value)
				if argval.Kind == FuncProtoValue {
					fv := FuncValue{FuncProto: argval.Data.(*Function), AccessLink: frame}
					argval = Value{Kind: FunctionValue, Data: fv}
					newOperands = append(newOperands, &Item{Type: ValueItem, Data: argval})
				} else {
					newOperands = append(newOperands, v)
				}
			default:
				newOperands = append(newOperands, v)
			}
		}

		if !item.Expand {
			return operTbl[opcall.OperID](frame, newOperands)
		}
		var expandedOperands []*Item
		for idx, v := range newOperands {
			if !item.ExpandArgIndexes[idx] {
				expandedOperands = append(expandedOperands, v)
				continue
			}
			listVal := EvalItem(v, frame)
			if listVal.Kind != ListValue {
				runTimeError2(frame, "Expansion requires list value")
			}
			lit := NewListIterator(listVal)
			for {
				nextv := lit.Next()
				if nextv == nil {
					break
				}
				expandedOperands = append(expandedOperands, &Item{Type: ValueItem, Data: *nextv})
			}
		}
		return operTbl[opcall.OperID](frame, expandedOperands)
	default:
		runTimeError2(frame, "Data corrupted")
	}
	return
}
