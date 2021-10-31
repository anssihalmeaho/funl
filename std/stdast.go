package std

import (
	"fmt"
	"strings"

	"github.com/anssihalmeaho/funl/funl"
)

func initSTDAst(interpreter *funl.Interpreter) (err error) {
	stdModuleName := "stdast"
	topFrame := &funl.Frame{
		Syms:     funl.NewSymt(),
		OtherNS:  make(map[funl.SymID]funl.ImportInfo),
		Imported: make(map[funl.SymID]*funl.Frame),
	}
	stdAstFuncs := []stdFuncInfo{
		{
			Name:       "func-value-to-ast",
			Getter:     getFuncValueToAst,
			IsFunction: true,
		},
		{
			Name:       "parse-ns-to-ast",
			Getter:     getParseNSToAst,
			IsFunction: true,
		},
		{
			Name:       "parse-expr-to-ast",
			Getter:     getParseExprToAst,
			IsFunction: true,
		},
		{
			Name:       "eval-ast",
			Getter:     getEvalAst,
			IsFunction: true,
		},
		{
			Name:       "p-eval-ast",
			Getter:     getEvalAst,
			IsFunction: false,
		},
		{
			Name:       "import-ast",
			Getter:     getImportAst,
			IsFunction: false,
		},
	}
	err = setSTDFunctions(topFrame, stdModuleName, stdAstFuncs, interpreter)
	return
}

func putToMap(frame *funl.Frame, prevmap funl.Value, name string, value funl.Value) funl.Value {
	putArgs := []*funl.Item{
		&funl.Item{Type: funl.ValueItem, Data: prevmap},
		&funl.Item{Type: funl.ValueItem, Data: funl.Value{Kind: funl.StringValue, Data: name}},
		&funl.Item{Type: funl.ValueItem, Data: value},
	}
	return funl.HandlePutOP(frame, putArgs)
}

func parseItem(frame *funl.Frame, item *funl.Item) funl.Value {
	mapv := funl.HandleMapOP(frame, []*funl.Item{})
	switch item.Type {
	case funl.ValueItem:
		v := item.Data.(funl.Value)
		switch v.Kind {

		case funl.IntValue, funl.StringValue, funl.BoolValue, funl.FloatValue:
			mapv = putToMap(frame, mapv, "val", v)

		case funl.FunctionValue:
			fv := v.Data.(funl.FuncValue)
			funcProtoVal := parseFunctionProto(frame, fv.FuncProto)
			mapv = putToMap(frame, mapv, "func", funcProtoVal)

		case funl.FuncProtoValue:
			fproto := v.Data.(*funl.Function)
			funcProtoVal := parseFunctionProto(frame, fproto)
			mapv = putToMap(frame, mapv, "func", funcProtoVal)

		default:
			funl.RunTimeError2(frame, "Unexpected value: %#v", v)
		}

	case funl.SymbolPathItem:
		sp := item.Data.(funl.SymbolPath)
		var symbolPathNames []funl.Value
		for _, sid := range sp {
			symnamev := funl.Value{Kind: funl.StringValue, Data: funl.SymIDMap.AsString(sid)}
			symbolPathNames = append(symbolPathNames, symnamev)
		}
		mapv = putToMap(frame, mapv, "sym", funl.MakeListOfValues(frame, symbolPathNames))

	case funl.OperCallItem:
		opcall := item.Data.(funl.OpCall)
		operatorName := fmt.Sprintf("%s", opcall.OperID)
		operList := []funl.Value{
			funl.Value{Kind: funl.StringValue, Data: operatorName},
		}
		for _, operand := range opcall.Operands {
			operList = append(operList, parseItem(frame, operand))
		}
		mapv = putToMap(frame, mapv, "op", funl.MakeListOfValues(frame, operList))

	default:
		funl.RunTimeError2(frame, "Invalid item: %v", item)
	}
	mapv = putToMap(frame, mapv, "expand", funl.Value{Kind: funl.BoolValue, Data: item.Expand})
	var indexes []funl.Value
	for k := range item.ExpandArgIndexes {
		indexes = append(indexes, funl.Value{Kind: funl.IntValue, Data: k})
	}
	mapv = putToMap(frame, mapv, "expand-idx", funl.MakeListOfValues(frame, indexes))
	return mapv
}

func parseNS(frame *funl.Frame, ns funl.NSpace) funl.Value {
	nsMap := funl.HandleMapOP(frame, []*funl.Item{})

	symbolMap := ns.Syms.AsMap()
	var letvalues []funl.Value
	for _, sid := range ns.Syms.Keys() {
		symvalue := parseItem(frame, symbolMap[sid])
		symname := funl.SymIDMap.AsString(sid)
		pairSlice := []funl.Value{funl.Value{Kind: funl.StringValue, Data: symname}, symvalue}
		pair := funl.MakeListOfValues(frame, pairSlice)
		letvalues = append(letvalues, pair)
	}
	nsMap = putToMap(frame, nsMap, "syms", funl.MakeListOfValues(frame, letvalues))

	importMap := funl.HandleMapOP(frame, []*funl.Item{})
	for k, v := range ns.OtherNS {
		symname := funl.SymIDMap.AsString(k)
		importMap = putToMap(frame, importMap, symname, funl.Value{Kind: funl.StringValue, Data: v.Path()})
	}
	nsMap = putToMap(frame, nsMap, "imports", importMap)

	return nsMap
}

func parseFunctionProto(frame *funl.Frame, funcProto *funl.Function) funl.Value {
	mapval := funl.HandleMapOP(frame, []*funl.Item{})

	var argNameVals []funl.Value
	for _, argSymID := range funcProto.ArgNames {
		argNameV := funl.Value{Kind: funl.StringValue, Data: funl.SymIDMap.AsString(argSymID)}
		argNameVals = append(argNameVals, argNameV)
	}
	argNameList := funl.MakeListOfValues(frame, argNameVals)
	mapval = putToMap(frame, mapval, "args", argNameList)
	mapval = putToMap(frame, mapval, "ns", parseNS(frame, funcProto.NSpace))
	mapval = putToMap(frame, mapval, "body", parseItem(frame, funcProto.Body))
	mapval = putToMap(frame, mapval, "is-proc", funl.Value{Kind: funl.BoolValue, Data: funcProto.IsProc})

	mapval = putToMap(frame, mapval, "line", funl.Value{Kind: funl.IntValue, Data: funcProto.Lineno})
	mapval = putToMap(frame, mapval, "pos", funl.Value{Kind: funl.IntValue, Data: funcProto.Pos})
	mapval = putToMap(frame, mapval, "file", funl.Value{Kind: funl.StringValue, Data: funcProto.SrcFileName})
	return mapval
}

type parserErrHandler struct{}

func (eh *parserErrHandler) HandleParseError(errorText string) {
	funl.RunTimeError(errorText)
}

func getLetNameValue(frame *funl.Frame, pairList *funl.Value) (string, funl.Value) {
	nvIter := funl.NewListIterator(*pairList)
	name := nvIter.Next()
	val := nvIter.Next()
	if nvIter.Next() != nil {
		funl.RunTimeError2(frame, "Invalid name value pair")
	}
	if name.Kind != funl.StringValue {
		funl.RunTimeError2(frame, "Name not string")
	}
	return name.Data.(string), *val
}

func getNameValue(frame *funl.Frame, m funl.Value, name string) (bool, funl.Value) {
	mItem := &funl.Item{Type: funl.ValueItem, Data: m}
	nameItem := &funl.Item{Type: funl.ValueItem, Data: funl.Value{Kind: funl.StringValue, Data: name}}
	retListV := funl.HandleGetlOP(frame, []*funl.Item{mItem, nameItem})
	var retl []*funl.Value
	listIter := funl.NewListIterator(retListV)
	for {
		v := listIter.Next()
		if v == nil {
			break
		}
		retl = append(retl, v)
	}
	hasName, convOK := retl[0].Data.(bool)
	if !(hasName && convOK) {
		return false, funl.Value{}
	}
	return true, *retl[1]
}

func importNS(frame *funl.Frame, astmap funl.Value) (ok bool, errText string) {
	hasIt, nspaceV := getNameValue(frame, astmap, "nspace")
	if !hasIt {
		funl.RunTimeError2(frame, "nspace not found")
	}
	if nspaceV.Kind != funl.MapValue {
		funl.RunTimeError2(frame, "nspace is invalid")
	}

	hasIt, modNameV := getNameValue(frame, astmap, "name")
	if !hasIt {
		funl.RunTimeError2(frame, "module name not found")
	}
	if modNameV.Kind != funl.StringValue {
		funl.RunTimeError2(frame, "module name is invalid")
	}

	nspace := makeNS(frame, nspaceV)

	interpreter := frame.GetTopFrame().Interpreter
	funl.AddNStoCache(true, modNameV.Data.(string), nspace, interpreter)
	return true, ""
}

func makeNS(frame *funl.Frame, nsV funl.Value) *funl.NSpace {
	// ns : syms
	hasIt, symsV := getNameValue(frame, nsV, "syms")
	if !hasIt {
		funl.RunTimeError2(frame, "syms not found")
	}
	if symsV.Kind != funl.ListValue {
		funl.RunTimeError2(frame, "syms is invalid")
	}
	symsIter := funl.NewListIterator(symsV)
	symt := funl.NewSymt()
	for {
		v := symsIter.Next()
		if v == nil {
			break
		}
		symName, symValue := getLetNameValue(frame, v)

		// to avoid overlapping wasted name with existing one
		// lets generate new wasted name
		if strings.HasPrefix(symName, "__waste_") {
			symName = funl.GetWastedName()
		}

		sid := funl.SymIDMap.Add(symName)
		letItem := makeItem(frame, symValue)
		symAddOK := symt.AddBySID(sid, letItem)
		if !symAddOK {
			funl.RunTimeError2(frame, "syms found already (%s)", symName)
		}
	}

	// ns : imports
	hasIt, importsV := getNameValue(frame, nsV, "imports")
	if !hasIt {
		funl.RunTimeError2(frame, "imports not found")
	}
	if importsV.Kind != funl.MapValue {
		funl.RunTimeError2(frame, "imports is invalid")
	}
	keyvals := funl.HandleKeyvalsOP(frame, []*funl.Item{&funl.Item{Type: funl.ValueItem, Data: importsV}})
	kvListIter := funl.NewListIterator(keyvals)
	otherNSMap := make(map[funl.SymID]funl.ImportInfo)
	for {
		nextKV := kvListIter.Next()
		if nextKV == nil {
			break
		}
		kvIter := funl.NewListIterator(*nextKV)
		keyv := *(kvIter.Next())
		valv := *(kvIter.Next())
		if keyv.Kind != funl.StringValue {
			funl.RunTimeError2(frame, "key not a string: %v", keyv)
		}
		if valv.Kind != funl.StringValue {
			funl.RunTimeError2(frame, "value not a string: %v", valv)
		}
		sid := funl.SymIDMap.Add(keyv.Data.(string))
		var importInfo funl.ImportInfo
		importInfo.SetPath(valv.Data.(string))
		otherNSMap[sid] = importInfo
	}

	return &funl.NSpace{Syms: symt, OtherNS: otherNSMap}
}

func makeFuncValue(frame *funl.Frame, funcMap funl.Value) *funl.Item {
	// is-proc
	hasIt, isProcV := getNameValue(frame, funcMap, "is-proc")
	if !hasIt {
		funl.RunTimeError2(frame, "is-proc not found")
	}
	if isProcV.Kind != funl.BoolValue {
		funl.RunTimeError2(frame, "is-proc is invalid")
	}

	// args
	hasIt, argNamesV := getNameValue(frame, funcMap, "args")
	if !hasIt {
		funl.RunTimeError2(frame, "args not found")
	}
	if argNamesV.Kind != funl.ListValue {
		funl.RunTimeError2(frame, "args is invalid")
	}
	argsIter := funl.NewListIterator(argNamesV)
	var argNames []funl.SymID
	for {
		v := argsIter.Next()
		if v == nil {
			break
		}
		if v.Kind != funl.StringValue {
			funl.RunTimeError2(frame, "arg not string")
		}
		sid, found := funl.SymIDMap.Get(v.Data.(string))
		if !found {
			sid = funl.SymIDMap.Add(v.Data.(string))
		}
		argNames = append(argNames, sid)
	}

	// body
	hasIt, bodyV := getNameValue(frame, funcMap, "body")
	if !hasIt {
		funl.RunTimeError2(frame, "body not found")
	}
	if bodyV.Kind != funl.MapValue {
		funl.RunTimeError2(frame, "body is invalid")
	}
	bodyItem := makeItem(frame, bodyV)

	// ns
	hasIt, nsV := getNameValue(frame, funcMap, "ns")
	if !hasIt {
		funl.RunTimeError2(frame, "ns not found")
	}
	if nsV.Kind != funl.MapValue {
		funl.RunTimeError2(frame, "ns is invalid")
	}
	nspace := makeNS(frame, nsV)

	// line
	hasIt, lineV := getNameValue(frame, funcMap, "line")
	if !hasIt {
		funl.RunTimeError2(frame, "line not found")
	}
	if lineV.Kind != funl.IntValue {
		funl.RunTimeError2(frame, "line is invalid")
	}

	// pos
	hasIt, posV := getNameValue(frame, funcMap, "pos")
	if !hasIt {
		funl.RunTimeError2(frame, "pos not found")
	}
	if posV.Kind != funl.IntValue {
		funl.RunTimeError2(frame, "pos is invalid")
	}

	// file
	hasIt, fileV := getNameValue(frame, funcMap, "file")
	if !hasIt {
		funl.RunTimeError2(frame, "file not found")
	}
	if fileV.Kind != funl.StringValue {
		funl.RunTimeError2(frame, "file is invalid")
	}

	fproto := funl.Function{
		IsProc:   isProcV.Data.(bool),
		ArgNames: argNames,
		Body:     bodyItem,
		//NSpace:      funl.NSpace{Syms: symt, OtherNS: otherNSMap},
		NSpace:      *nspace,
		Lineno:      lineV.Data.(int),
		Pos:         posV.Data.(int),
		SrcFileName: fileV.Data.(string),
	}
	/*
		fval := funl.FuncValue{
			FuncProto:  &fproto,
			AccessLink: frame,
		}
		funcVal := funl.Value{Kind: funl.FunctionValue, Data: fval}
	*/
	funcVal := funl.Value{Kind: funl.FuncProtoValue, Data: &fproto}

	exprItem := &funl.Item{
		Type:             funl.ValueItem,
		Data:             funcVal,
		Expand:           false,
		ExpandArgIndexes: make(map[int]bool),
		//Expand:           getExpandFlag(),
		//ExpandArgIndexes: getExpandIndexes(),
	}
	return exprItem
}

func getExpandFlag(frame *funl.Frame, astmap funl.Value) bool {
	hasFlag, val := getNameValue(frame, astmap, "expand")
	if !hasFlag {
		funl.RunTimeError2(frame, "expand flag not found")
	}
	if val.Kind != funl.BoolValue {
		funl.RunTimeError2(frame, "Invalid expand flag")
	}
	return val.Data.(bool)
}

func getExpandIndexes(frame *funl.Frame, astmap funl.Value) map[int]bool {
	hasIdxs, val := getNameValue(frame, astmap, "expand-idx")
	if !hasIdxs {
		funl.RunTimeError2(frame, "expand indexes not found")
	}
	if val.Kind != funl.ListValue {
		funl.RunTimeError2(frame, "Invalid expand indexes")
	}
	listIter := funl.NewListIterator(val)
	indexMap := make(map[int]bool)
	for {
		v := listIter.Next()
		if v == nil {
			break
		}
		if v.Kind != funl.IntValue {
			funl.RunTimeError2(frame, "Invalid expand index (should be int)(%#v)", v)
		}
		indexMap[v.Data.(int)] = true
	}
	return indexMap
}

func makeItem(frame *funl.Frame, astmap funl.Value) *funl.Item {
	var exprItem *funl.Item
	if hasVal, val := getNameValue(frame, astmap, "val"); hasVal {
		exprItem = &funl.Item{
			Type:             funl.ValueItem,
			Data:             val,
			Expand:           getExpandFlag(frame, astmap),
			ExpandArgIndexes: getExpandIndexes(frame, astmap),
		}
	} else if hasSym, symList := getNameValue(frame, astmap, "sym"); hasSym {
		var symPath funl.SymbolPath
		listIter := funl.NewListIterator(symList)
		for {
			v := listIter.Next()
			if v == nil {
				break
			}
			if v.Kind != funl.StringValue {
				funl.RunTimeError2(frame, "Invalid symbol (%#v)", v)
			}
			symStr, convOK := v.Data.(string)
			if !convOK {
				funl.RunTimeError2(frame, "Invalid symbol string (%#v)", v)
			}
			sid, found := funl.SymIDMap.Get(symStr)
			if !found {
				sid = funl.SymIDMap.Add(symStr)
			}
			symPath = append(symPath, sid)
		}

		exprItem = &funl.Item{
			Type:             funl.SymbolPathItem,
			Data:             symPath,
			Expand:           getExpandFlag(frame, astmap),
			ExpandArgIndexes: getExpandIndexes(frame, astmap),
		}
	} else if hasOP, opList := getNameValue(frame, astmap, "op"); hasOP {
		listIter := funl.NewListIterator(opList)
		var ops []*funl.Value
		for {
			v := listIter.Next()
			if v == nil {
				break
			}
			ops = append(ops, v)
		}
		if l := len(ops); l < 1 {
			funl.RunTimeError2(frame, "Too short operator call (%d)", l)
		}
		if ops[0].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "Invalid operator (%#v)", ops)
		}
		opName, convOK := ops[0].Data.(string)
		if !convOK {
			funl.RunTimeError2(frame, "Invalid symbol string (%#v)", ops)
		}

		opType, opFound := funl.OperNameToID(opName)
		if !opFound {
			funl.RunTimeError2(frame, "Operator not found (%s)", opName)
		}

		var operands []*funl.Item
		for _, argVal := range ops[1:len(ops)] {
			operandItem := makeItem(frame, *argVal)
			operands = append(operands, operandItem)
		}

		opCall := funl.OpCall{
			OperID:   opType,
			Operands: operands,
		}

		exprItem = &funl.Item{
			Type:             funl.OperCallItem,
			Data:             opCall,
			Expand:           getExpandFlag(frame, astmap),
			ExpandArgIndexes: getExpandIndexes(frame, astmap),
		}
	} else if hasFunc, funcMap := getNameValue(frame, astmap, "func"); hasFunc {
		exprItem = makeFuncValue(frame, funcMap)
	} else {
		funl.RunTimeError2(frame, "Invalid expression")
	}
	return exprItem
}

func evalExpr(frame *funl.Frame, astmap funl.Value) funl.Value {
	exprItem := makeItem(frame, astmap)
	return funl.EvalItem(exprItem, frame)
}

func getImportAst(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		checkValidity := func(frame *funl.Frame, arguments []funl.Value) (ok bool, errText string, astmap funl.Value) {
			if l := len(arguments); l != 1 {
				funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need one", name, l)
			}
			if arguments[0].Kind != funl.MapValue {
				ok, errText = false, fmt.Sprintf("%s: requires map value", name)
				return
			}
			return true, "", arguments[0]
		}

		ok, errText, astmap := checkValidity(frame, arguments)
		if !ok {
			values := []funl.Value{
				{
					Kind: funl.BoolValue,
					Data: false,
				},
				{
					Kind: funl.StringValue,
					Data: errText,
				},
			}
			retVal = funl.MakeListOfValues(frame, values)
			return
		}
		ok, errText = importNS(frame, astmap)
		values := []funl.Value{
			{
				Kind: funl.BoolValue,
				Data: ok,
			},
			{
				Kind: funl.StringValue,
				Data: errText,
			},
		}
		retVal = funl.MakeListOfValues(frame, values)
		return
	}
}

func getEvalAst(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		checkValidity := func(frame *funl.Frame, arguments []funl.Value) (ok bool, errText string, astmap funl.Value) {
			if l := len(arguments); l != 1 {
				funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need one", name, l)
			}
			if arguments[0].Kind != funl.MapValue {
				ok, errText = false, fmt.Sprintf("%s: requires map value", name)
				return
			}
			return true, "", arguments[0]
		}

		ok, errText, astmap := checkValidity(frame, arguments)
		if !ok {
			values := []funl.Value{
				{
					Kind: funl.BoolValue,
					Data: false,
				},
				{
					Kind: funl.StringValue,
					Data: errText,
				},
				funl.HandleMapOP(frame, []*funl.Item{}),
			}
			retVal = funl.MakeListOfValues(frame, values)
			return
		}
		values := []funl.Value{
			{
				Kind: funl.BoolValue,
				Data: ok,
			},
			{
				Kind: funl.StringValue,
				Data: errText,
			},
			evalExpr(frame, astmap),
		}
		retVal = funl.MakeListOfValues(frame, values)
		return
	}
}

func getParseExprToAst(name string) stdFuncType {
	checkValidity := func(frame *funl.Frame, arguments []funl.Value) (ok bool, errText string, content string) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need one", name, l)
		}
		if arguments[0].Kind != funl.StringValue {
			ok, errText = false, fmt.Sprintf("%s: requires string value", name)
			return
		}
		content, ok = arguments[0].Data.(string)
		if !ok {
			errText = fmt.Sprintf("%s: argument is not string value", name)
			return
		}
		return
	}

	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		ok, errText, content := checkValidity(frame, arguments)
		if !ok {
			values := []funl.Value{
				{
					Kind: funl.BoolValue,
					Data: false,
				},
				{
					Kind: funl.StringValue,
					Data: errText,
				},
				funl.HandleMapOP(frame, []*funl.Item{}),
			}
			retVal = funl.MakeListOfValues(frame, values)
			return
		}

		var srcFileName string
		parser := funl.NewParser(funl.NewDefaultOperators(), &srcFileName)
		parser.SetErrorHandler(&parserErrHandler{})
		exprItem, err := parser.ParseOneExpression(string(content))
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
				funl.HandleMapOP(frame, []*funl.Item{}),
			}
			retVal = funl.MakeListOfValues(frame, values)
			return
		}

		values := []funl.Value{
			{
				Kind: funl.BoolValue,
				Data: ok,
			},
			{
				Kind: funl.StringValue,
				Data: errText,
			},
			parseItem(frame, exprItem),
		}
		retVal = funl.MakeListOfValues(frame, values)
		return
	}
}

func getParseNSToAst(name string) stdFuncType {
	checkValidity := func(frame *funl.Frame, arguments []funl.Value) (ok bool, errText string, content string) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need one", name, l)
		}
		if arguments[0].Kind != funl.StringValue {
			ok, errText = false, fmt.Sprintf("%s: requires string value", name)
			return
		}
		content, ok = arguments[0].Data.(string)
		if !ok {
			errText = fmt.Sprintf("%s: argument is not string value", name)
			return
		}
		return
	}

	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		ok, errText, content := checkValidity(frame, arguments)
		if !ok {
			values := []funl.Value{
				{
					Kind: funl.BoolValue,
					Data: false,
				},
				{
					Kind: funl.StringValue,
					Data: errText,
				},
				funl.HandleMapOP(frame, []*funl.Item{}),
			}
			retVal = funl.MakeListOfValues(frame, values)
			return
		}

		var srcFileName string
		parser := funl.NewParser(funl.NewDefaultOperators(), &srcFileName)
		parser.SetErrorHandler(&parserErrHandler{})
		var nsName string
		var nspace *funl.NSpace
		var err error
		nsName, nspace, err = parser.Parse(string(content))
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
				funl.HandleMapOP(frame, []*funl.Item{}),
			}
			retVal = funl.MakeListOfValues(frame, values)
			return
		}

		mapval := funl.HandleMapOP(frame, []*funl.Item{})
		mapval = putToMap(frame, mapval, "nspace", parseNS(frame, *nspace))
		mapval = putToMap(frame, mapval, "name", funl.Value{Kind: funl.StringValue, Data: nsName})

		values := []funl.Value{
			{
				Kind: funl.BoolValue,
				Data: ok,
			},
			{
				Kind: funl.StringValue,
				Data: errText,
			},
			mapval,
		}
		retVal = funl.MakeListOfValues(frame, values)
		return
	}
}

func getFuncValueToAst(name string) stdFuncType {
	checkValidity := func(frame *funl.Frame, arguments []funl.Value) (ok bool, errText string, funcVal funl.FuncValue) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d), need one", name, l)
		}
		if arguments[0].Kind != funl.FunctionValue {
			ok, errText = false, fmt.Sprintf("%s: requires func/proc value", name)
			return
		}
		funcVal, ok = arguments[0].Data.(funl.FuncValue)
		if !ok {
			errText = fmt.Sprintf("%s: argument is not func/proc value", name)
			return
		}
		return
	}

	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		ok, errText, _ := checkValidity(frame, arguments)
		if !ok {
			values := []funl.Value{
				{
					Kind: funl.BoolValue,
					Data: false,
				},
				{
					Kind: funl.StringValue,
					Data: errText,
				},
				funl.HandleMapOP(frame, []*funl.Item{}),
			}
			retVal = funl.MakeListOfValues(frame, values)
			return
		}

		fItem := &funl.Item{
			Type: funl.ValueItem,
			Data: arguments[0],
		}
		values := []funl.Value{
			{
				Kind: funl.BoolValue,
				Data: ok,
			},
			{
				Kind: funl.StringValue,
				Data: errText,
			},
			parseItem(frame, fItem),
			//parseFunctionProto(frame, funcVal.FuncProto),
		}
		retVal = funl.MakeListOfValues(frame, values)
		return
	}
}
