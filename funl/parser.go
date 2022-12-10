package funl

import (
	"fmt"
	"os"
	"strconv"
)

type tokenIter struct {
	location int
	tokens   []token
	lexer    lexerAPI
}

func (iter *tokenIter) tokenize(sourceText string) (err error) {
	//iter.tokens, err = iter.lexer.scan(sourceText)
	var tokens []token
	iter.tokens = []token{}
	tokens, err = iter.lexer.scan(sourceText)
	if err != nil {
		return
	}
	for _, token := range tokens {
		switch token.Type {
		case tokenMultiLineComment, tokenLineComment:
		default:
			iter.tokens = append(iter.tokens, token)
		}
	}
	return
}

func (iter *tokenIter) throwAway() {
	iter.location++
}

func (iter *tokenIter) lookAhead(args ...int) (token, bool) {
	offset := 0
	if len(args) > 0 {
		offset = args[0]
	}
	if (iter.location + offset) < len(iter.tokens) {
		return iter.tokens[iter.location+offset], true
	}
	return token{}, false
}

func (iter *tokenIter) next() (token, bool) {
	if iter.location >= len(iter.tokens) {
		return token{}, false
	}
	token := iter.tokens[iter.location]
	iter.location++
	return token, true
}

func newTokenIter(lexer lexerAPI) *tokenIter {
	return &tokenIter{lexer: lexer}
}

// NodeType defines type of AST node
type NodeType int

const (
	NodeTypeRoot NodeType = iota
	NodeTypeValue
	NodeTypeLet
	NodeTypeSymbolPath
	NodeTypeOperCall
)

// Node represents one syntax node in syntax tree
type Node struct {
	NodeType NodeType
	Childs   []*Node
}

// Parser implements syntax tree parser
type Parser struct {
	tokenIter        *tokenIter
	root             Node
	operators        Operators
	canEndWithSymbol bool
	srcFileName      *string
	errorHandler     ParseErrorHandler
}

type ParseErrorHandler interface {
	HandleParseError(errorText string)
}

func (p *Parser) SetErrorHandler(peh ParseErrorHandler) {
	p.errorHandler = peh
}

func (p *Parser) stopOnError(line interface{}, format string, args ...interface{}) {
	lineStr := ""
	if line != nil {
		lineStr = fmt.Sprintf("line %d: ", line.(int)) // atoi ??
	}
	var fileName string
	if p.srcFileName == nil {
		fileName = ""
	} else {
		fileName = fmt.Sprintf("%s: ", *(p.srcFileName))
	}
	errorText := fmt.Sprintf("Syntax error: "+fileName+lineStr+format, args...)
	if p.errorHandler != nil {
		p.errorHandler.HandleParseError(errorText)
		return
	}
	fmt.Println(errorText)
	os.Exit(-1)
}

func (p *Parser) SetCanEndWithSymbol(canSet bool) {
	p.canEndWithSymbol = canSet
}

// ParseImport parses import definition
func (p *Parser) ParseImport() (modName string, importInfo ImportInfo) {
	importPath := ""

	parseShortImport := func() {
		nextToken, any := p.tokenIter.next()
		if !any {
			p.stopOnError(nil, "Invalid import")
		}
		if nextToken.Type != tokenSymbol {
			p.stopOnError(nextToken.Lineno, "Invalid import, assumed symbol (%s)", nextToken.Value)
		}
		modName = nextToken.Value
	}

	parseLongImport := func() {
		nextToken, any := p.tokenIter.next()
		if !any {
			p.stopOnError(nil, "Invalid import")
		}
		if nextToken.Type != tokenString {
			p.stopOnError(nextToken.Lineno, "Invalid import, assumed string (%s)", nextToken.Value)
		}
		importPath = nextToken.Value
		nextToken, any = p.tokenIter.next()
		if !any {
			p.stopOnError(nil, "Invalid import")
		}
		if nextToken.Type != tokenAs {
			p.stopOnError(nextToken.Lineno, "as missing")
		}
		parseShortImport()
	}

	p.tokenIter.throwAway()
	token, hasAny := p.tokenIter.lookAhead()
	if !hasAny {
		p.stopOnError(nil, "Invalid import")
	}

	switch token.Type {
	case tokenString:
		parseLongImport()
	case tokenSymbol:
		parseShortImport()
	default:
		p.stopOnError(token.Lineno, "Invalid token in import: %s", token.Value)
	}

	importInfo.importPath = importPath
	DebugPrint("import:%s:%s \n", importPath, modName)
	return
}

func isValue(token token) bool {
	switch token.Type {
	case tokenNumber, tokenTrue, tokenFalse, tokenString, tokenFuncBegin, tokenProcBegin:
		return true
	}
	return false
}

func (p *Parser) isExpandedLetDef() bool {
	var counter int
	for {
		nextToken, anyFound := p.tokenIter.lookAhead(counter)
		if !anyFound {
			return false
		}
		switch nextToken.Type {
		case tokenSymbol:
		case tokenEqualsSign:
			return true
		default:
			return false
		}
		counter++
	}
}

func (p *Parser) ParseExpandedLet() (symbols []string, item *Item) {
	for {
		nextToken, anyFound := p.tokenIter.next()
		if !anyFound {
			p.stopOnError(nil, "Invalid (expanded) let definition, no symbol found")
		}
		switch nextToken.Type {
		case tokenSymbol:
			symbols = append(symbols, nextToken.Value)
		case tokenEqualsSign:
			item = p.ParseExpr()
			DebugPrint("expanded lets: %v", symbols)
			return
		default:
			p.stopOnError(nextToken.Lineno, "Invalid let definition, symbol assumed (%s)", nextToken.Value)
		}
	}
}

func (p *Parser) ParseFuncValue(procOrFuncToken token) (funcData *Function) {
	DebugPrint("func start")

	funcData = &Function{NSpace: NSpace{Syms: NewSymt(), OtherNS: make(map[SymID]ImportInfo)}}

	switch procOrFuncToken.Type {
	case tokenProcBegin:
		funcData.IsProc = true
	case tokenFuncBegin:
		funcData.IsProc = false
	default:
		p.stopOnError(procOrFuncToken.Lineno, "Invalid function keyword : %#v", procOrFuncToken)
	}

	token, hasAny := p.tokenIter.next()
	if !hasAny {
		p.stopOnError(nil, "Invalid function value, nothing found")
	}
	if token.Type != tokenOpenBracket {
		p.stopOnError(token.Lineno, "Invalid function value, starting bracket assumed, found: %s", token.Value)
	}

	var argSymbolList []string
	// lets read arguments:
	// lets check empty argument list case beforehand
	if tok, _ := p.tokenIter.lookAhead(); tok.Type == tokenClosingBracket {
		p.tokenIter.throwAway()
	} else {
	ArgLoop:
		for {
			token, hasAny := p.tokenIter.lookAhead()
			if !hasAny {
				p.stopOnError(nil, "Invalid argument list in function")
			}
			switch token.Type {
			case tokenClosingBracket:
				p.tokenIter.throwAway()
				break ArgLoop
			case tokenComma:
				p.tokenIter.throwAway()
			case tokenStartNS:
				p.stopOnError(token.Lineno, "Namespace definition not allowed in function")
			case tokenEndNS:
				p.stopOnError(token.Lineno, "Namespace end definition not allowed in function")
			default:
				argToken, hasArg := p.tokenIter.next()
				if !hasArg {
					p.stopOnError(nil, "No function argument found")
				}
				if argToken.Type != tokenSymbol {
					p.stopOnError(argToken.Lineno, "Invalid functiona argument (%s)", argToken.Value)
				}
				argSymbolList = append(argSymbolList, argToken.Value)

				sid := SymIDMap.Add(argToken.Value)
				funcData.ArgNames = append(funcData.ArgNames, sid)
			}
		}
	}

	// then parse other parts
	var bodyFound bool
	setBodyFound := func(lineno int) {
		if bodyFound {
			p.stopOnError(lineno, "More than one body for function found")
		}
		bodyFound = true
	}
BlockLoop:
	for {
		token, hasAny := p.tokenIter.lookAhead()
		if !hasAny {
			p.stopOnError(nil, "Invalid function value")
		}

		switch token.Type {
		case tokenFuncEnd:
			p.tokenIter.throwAway()
			DebugPrint("func end (args: %v)", argSymbolList)
			break BlockLoop
		case tokenFuncBegin, tokenProcBegin:
			funcData.Body = p.ParseExpr()
			setBodyFound(token.Lineno)
		case tokenImport:
			modName, item := p.ParseImport()
			sid := SymIDMap.Add(modName)
			if _, modFound := funcData.NSpace.OtherNS[sid]; modFound {
				p.stopOnError(token.Lineno, "Module already imported (%s)", modName)
			}
			funcData.NSpace.OtherNS[sid] = item
		case tokenSymbol:
			secondToken, hasSecond := p.tokenIter.lookAhead(1)
			if !hasSecond {
				p.stopOnError(nil, "Invalid function structure")
			}
			if secondToken.Type == tokenEqualsSign {
				letName, item := p.ParseLet()

				var isExpanderCase bool
				// check if its special expander case like: x = list(1):
				expToken, hasAny := p.tokenIter.lookAhead()
				if hasAny && (expToken.Type == tokenExpander) {
					isExpanderCase = true
				}

				if isExpanderCase {
					p.tokenIter.throwAway()
					wasteName := getWastedName()
					wasteSymID := SymIDMap.Add(wasteName)
					if err := funcData.NSpace.Syms.Add(wasteName, item); err != nil {
						p.stopOnError(secondToken.Lineno, "Failed to add let def. to symbol table (%s)", wasteName)
					}

					wasteSymbolItem := &Item{Type: SymbolPathItem, Data: SymbolPath{wasteSymID}}
					indexItem := &Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 0}}
					opc := OpCall{OperID: IndOP, Operands: []*Item{wasteSymbolItem, indexItem}}
					indOPcallItem := &Item{Type: OperCallItem, Data: opc}
					if err := funcData.NSpace.Syms.Add(letName, indOPcallItem); err != nil {
						p.stopOnError(secondToken.Lineno, "Failed to add let def. to symbol table (%s)", letName)
					}
				} else {
					if err := funcData.NSpace.Syms.Add(letName, item); err != nil {
						p.stopOnError(secondToken.Lineno, "Failed to add let def. to symbol table (%s)", letName)
					}
				}
			} else if p.isExpandedLetDef() {
				letNames, item := p.ParseExpandedLet()
				wasteName := getWastedName()
				wasteSymID := SymIDMap.Add(wasteName)
				if err := funcData.NSpace.Syms.Add(wasteName, item); err != nil {
					p.stopOnError(secondToken.Lineno, "Failed to add let def. to symbol table (%s)", wasteName)
				}

				for letSymIndex, letName := range letNames {
					wasteSymbolItem := &Item{Type: SymbolPathItem, Data: SymbolPath{wasteSymID}}
					indexItem := &Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: letSymIndex}}
					opc := OpCall{OperID: IndOP, Operands: []*Item{wasteSymbolItem, indexItem}}
					indOPcallItem := &Item{Type: OperCallItem, Data: opc}
					if err := funcData.NSpace.Syms.Add(letName, indOPcallItem); err != nil {
						p.stopOnError(secondToken.Lineno, "Failed to add let def. to symbol table (%s)", letName)
					}
				}

				// read expander token too
				expToken, hasAny := p.tokenIter.lookAhead()
				if !hasAny {
					p.stopOnError(nil, "assuming expander, found nothing")
				}
				if expToken.Type != tokenExpander {
					p.stopOnError(expToken.Lineno, "assuming expander")
				}
				p.tokenIter.throwAway()
			} else {
				funcData.Body = p.ParseExpr()
				setBodyFound(token.Lineno)
			}
		case tokenStartNS:
			p.stopOnError(token.Lineno, "Namespace definition not allowed in function")
		case tokenEndNS:
			p.stopOnError(token.Lineno, "Namespace end definition not allowed in function")
		default:
			funcData.Body = p.ParseExpr()
			setBodyFound(token.Lineno)
		}
	}
	if !bodyFound {
		p.stopOnError(token.Lineno, "Function body not found")
	}
	return
}

// ParseValue parses value
func (p *Parser) ParseValue() (item *Item) {
	token, hasAny := p.tokenIter.next()
	if !hasAny {
		p.stopOnError(nil, "Invalid value, nothing found")
	}
	var value Value
	switch token.Type {
	case tokenNumber:
		// lets first check if its float value
		firstlookAhead, hasSome := p.tokenIter.lookAhead()
		if hasSome && firstlookAhead.Type == tokenDot {
			secondlookAhead, hasStillSome := p.tokenIter.lookAhead(1)
			if hasStillSome && secondlookAhead.Type == tokenNumber {
				floatStr := token.Value + "." + secondlookAhead.Value
				floatv, pfErr := strconv.ParseFloat(floatStr, 64)
				if pfErr != nil {
					p.stopOnError(token.Lineno, "Float value reading failed (%s)", floatStr)
				}
				value.Kind = FloatValue
				value.Data = floatv
				p.tokenIter.throwAway()
				p.tokenIter.throwAway()
			}
		}
		if value.Kind != FloatValue {
			// it wasnt float, lets assume its int
			numval, err := strconv.Atoi(token.Value)
			if err != nil {
				p.stopOnError(token.Lineno, "Invalid numeric value: %s", token.Value)
			}
			DebugPrint("num value: %d", numval)
			value.Kind = IntValue
			value.Data = numval
		}
	case tokenTrue:
		DebugPrint("bool value: true")
		value.Kind = BoolValue
		value.Data = true
	case tokenFalse:
		DebugPrint("bool value: false")
		value.Kind = BoolValue
		value.Data = false
	case tokenString:
		DebugPrint("string value: %s", token.Value)
		value.Kind = StringValue
		value.Data = token.Value
	case tokenFuncBegin, tokenProcBegin:
		funcData := p.ParseFuncValue(token)
		funcData.Lineno, funcData.Pos = token.Lineno, token.Pos
		if p.srcFileName != nil {
			funcData.SrcFileName = *(p.srcFileName)
		}
		DebugPrint("func value: %s", token.Value)
		value.Kind = FuncProtoValue
		value.Data = funcData
	}
	item = &Item{Type: ValueItem, Data: value}
	return
}

// OperNameToID is for std usage
func OperNameToID(opName string) (op OperType, ok bool) {
	return operNameToID(opName)
}

func operNameToID(opName string) (op OperType, ok bool) {
	switch opName {
	case "and":
		op = AndOP
	case "or":
		op = OrOP
	case "not":
		op = NotOP
	case "call":
		op = CallOP
	case "eq":
		op = EqOP
	case "if":
		op = IfOP
	case "plus":
		op = PlusOP
	case "minus":
		op = MinusOP
	case "mul":
		op = MulOP
	case "div":
		op = DivOP
	case "mod":
		op = ModOP
	case "list":
		op = ListOP
	case "empty":
		op = EmptyOP
	case "head":
		op = HeadOP
	case "last":
		op = LastOP
	case "rest":
		op = RestOP
	case "append":
		op = AppendOP
	case "add":
		op = AddOP
	case "len":
		op = LenOP
	case "type":
		op = TypeOP
	case "in":
		op = InOP
	case "ind":
		op = IndOP
	case "find":
		op = FindOP
	case "slice":
		op = SliceOP
	case "rrest":
		op = RrestOP
	case "reverse":
		op = ReverseOP
	case "extend":
		op = ExtendOP
	case "split":
		op = SplitOP
	case "gt":
		op = GtOP
	case "lt":
		op = LtOP
	case "le":
		op = LeOP
	case "ge":
		op = GeOP
	case "str":
		op = StrOP
	case "conv":
		op = ConvOP
	case "case":
		op = CaseOP
	case "name":
		op = NameOP
	case "error":
		op = ErrorOP
	case "print":
		op = PrintOP
	case "spawn":
		op = SpawnOP
	case "chan":
		op = ChanOP
	case "send":
		op = SendOP
	case "recv":
		op = RecvOP
	case "symval":
		op = SymvalOP
	case "try":
		op = TryOP
	case "tryl":
		op = TrylOP
	case "select":
		op = SelectOP
	case "eval":
		op = EvalOP
	case "while":
		op = WhileOP
	case "float":
		op = FloatOP
	case "map":
		op = MapOP
	case "put":
		op = PutOP
	case "get":
		op = GetOP
	case "getl":
		op = GetlOP
	case "keys":
		op = KeysOP
	case "vals":
		op = ValsOP
	case "keyvals":
		op = KeyvalsOP
	case "let":
		op = LetOP
	case "imp":
		op = ImpOP
	case "del":
		op = DelOP
	case "dell":
		op = DellOP
	case "sprintf":
		op = SprintfOP
	case "argslist":
		op = ArgslistOP
	case "cond":
		op = CondOP
	case "help":
		op = HelpOP
	case "recwith":
		op = RecwithOP
	default:
		return
	}
	ok = true
	return
}

// ParseOperCall parses operator call
func (p *Parser) ParseOperCall() (item *Item) {
	operName := ""

	token, hasAny := p.tokenIter.next()
	if !hasAny {
		p.stopOnError(nil, "Invalid operator call, nothing found")
	}
	operName = token.Value
	opid, ok := operNameToID(operName)
	if !ok {
		p.stopOnError(nil, "Invalid operator call, operator not found (%s)", operName)
	}
	opc := OpCall{OperID: opid, Operands: []*Item{}}

	token, hasAny = p.tokenIter.next()
	if !hasAny {
		p.stopOnError(nil, "Invalid operator call, starting bracket assumed, nothing found")
	}
	if token.Type != tokenOpenBracket {
		p.stopOnError(token.Lineno, "Invalid operator call, starting bracket assumed, found: %s", token.Value)
	}

	var expandOperands bool
	var argCounter int
	var expandIndexes map[int]bool
	expandIndexes = make(map[int]bool)
	// lets check empty argument list case beforehand
	if tok, _ := p.tokenIter.lookAhead(); tok.Type == tokenClosingBracket {
		p.tokenIter.throwAway()
	} else {
	ArgLoop:
		for {
			token, hasAny := p.tokenIter.lookAhead()
			if !hasAny {
				p.stopOnError(nil, "Invalid operator call")
			}
			switch token.Type {
			case tokenClosingBracket:
				p.tokenIter.throwAway()
				break ArgLoop
			case tokenComma:
				p.tokenIter.throwAway()
			case tokenStartNS:
				p.stopOnError(token.Lineno, "Namespace definition not allowed in operator call")
			case tokenEndNS:
				p.stopOnError(token.Lineno, "Namespace end definition not allowed in operator call")
			default:
				argItem := p.ParseExpr()
				opc.Operands = append(opc.Operands, argItem)

				// lets check if there is expansion
				tok, hasSome := p.tokenIter.lookAhead()
				if hasSome && tok.Type == tokenExpander {
					expandOperands = true
					expandIndexes[argCounter] = true
					p.tokenIter.throwAway()
				}

				argCounter++
			}
		}
	}
	DebugPrint("operator-call: %s", operName)
	item = &Item{Type: OperCallItem, Data: opc, Expand: expandOperands, ExpandArgIndexes: expandIndexes}
	return
}

// ParseSymbolPath
func (p *Parser) ParseSymbolPath() (item *Item) {
	var symPath []string

	sym, _ := p.tokenIter.next()
	symPath = append(symPath, sym.Value)
	for {
		token, hasSome := p.tokenIter.lookAhead()
		if !hasSome {
			if p.canEndWithSymbol {
				break
			}
			p.stopOnError(token.Lineno, "invalid symbol path (%s)", token.Value)
		}
		if token.Type == tokenDot {
			p.tokenIter.throwAway()
		} else {
			break
		}
		token, hasSome = p.tokenIter.lookAhead()
		if !hasSome {
			p.stopOnError(token.Lineno, "invalid symbol path, missing symbol (%s)", token.Value)
		}
		if token.Type != tokenSymbol {
			p.stopOnError(token.Lineno, "invalid symbol path, symbol assumed (%s)", token.Value)
		}
		sym, _ = p.tokenIter.next()
		symPath = append(symPath, sym.Value)
	}
	var symbolIDPath SymbolPath
	for _, v := range symPath {
		sid := SymIDMap.Add(v)
		symbolIDPath = append(symbolIDPath, sid)
	}
	item = &Item{Type: SymbolPathItem, Data: symbolIDPath}
	DebugPrint("symbol path: %v", symPath)
	return
}

// ParseExpr parses expression
func (p *Parser) ParseExpr() (item *Item) {
	token, hasAny := p.tokenIter.lookAhead()
	if !hasAny {
		p.stopOnError(nil, "Invalid expression, nothing found")
	}
	if isValue(token) {
		item = p.ParseValue()
	} else if token.Type == tokenSymbol {
		nextToken, hasAny := p.tokenIter.lookAhead(1)
		if hasAny && nextToken.Type == tokenOpenBracket {
			item = p.ParseOperCall()
		} else {
			item = p.ParseSymbolPath()
		}
	} else {
		p.stopOnError(token.Lineno, "Invalid expression (%s)", token.Value)
	}
	return
}

// ParseLet parses let definition
func (p *Parser) ParseLet() (symbol string, item *Item) {
	letSymbol := ""

	token, hasAny := p.tokenIter.next()
	if !hasAny {
		p.stopOnError(nil, "Invalid let definition, no symbol found")
	}
	if token.Type != tokenSymbol {
		p.stopOnError(token.Lineno, "Invalid let definition, symbol assumed (%s)", token.Value)
	}
	letSymbol = token.Value

	token, hasAny = p.tokenIter.next()
	if !hasAny {
		p.stopOnError(nil, "Invalid let definition, no equals sign (=) found")
	}
	if token.Type != tokenEqualsSign {
		p.stopOnError(token.Lineno, "Invalid let definition, equals sign (=) assumed (%s)", token.Value)
	}

	item = p.ParseExpr()
	symbol = letSymbol

	DebugPrint("let:%s", letSymbol)
	return
}

// ParseNamespace parses namespace definition
func (p *Parser) ParseNamespace() (nsName string, ns *NSpace) {
	p.tokenIter.throwAway()
	token, hasAny := p.tokenIter.next()
	if !hasAny {
		p.stopOnError(nil, "No name for namespace")
	}
	if token.Type != tokenSymbol {
		p.stopOnError(token.Lineno, "No namespace found")
	}
	nsName = token.Value
	DebugPrint("namespace start: %s", nsName)
	ns = &NSpace{OtherNS: make(map[SymID]ImportInfo), Syms: NewSymt()}
	for {
		token, hasAny = p.tokenIter.lookAhead()
		if !hasAny {
			p.stopOnError(nil, "No end of namespace found (%s)", nsName)
		}
		if token.Type == tokenEndNS {
			p.tokenIter.throwAway()
			break
		}
		switch token.Type {
		case tokenImport:
			modName, importInfo := p.ParseImport()
			sid := SymIDMap.Add(modName)
			ns.OtherNS[sid] = importInfo
		case tokenSymbol:
			secondToken, hasSecond := p.tokenIter.lookAhead(1)
			if !hasSecond {
				p.stopOnError(nil, "Invalid function structure")
			}
			if secondToken.Type == tokenEqualsSign {
				sym, item := p.ParseLet()

				var isExpanderCase bool
				// check if its special expander case like: x = list(1):
				expToken, hasAny := p.tokenIter.lookAhead()
				if hasAny && (expToken.Type == tokenExpander) {
					isExpanderCase = true
				}

				if isExpanderCase {
					p.tokenIter.throwAway()
					wasteName := getWastedName()
					wasteSymID := SymIDMap.Add(wasteName)
					if err := ns.Syms.Add(wasteName, item); err != nil {
						p.stopOnError(secondToken.Lineno, "Failed to add let def. to symbol table (%s)", wasteName)
					}

					wasteSymbolItem := &Item{Type: SymbolPathItem, Data: SymbolPath{wasteSymID}}
					indexItem := &Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 0}}
					opc := OpCall{OperID: IndOP, Operands: []*Item{wasteSymbolItem, indexItem}}
					indOPcallItem := &Item{Type: OperCallItem, Data: opc}
					if err := ns.Syms.Add(sym, indOPcallItem); err != nil {
						p.stopOnError(secondToken.Lineno, "Failed to add let def. to symbol table (%s)", sym)
					}
				} else {
					err := ns.Syms.Add(sym, item)
					if err != nil {
						p.stopOnError(token.Lineno, "Unexpected error : %v", err)
					}
				}
			} else if p.isExpandedLetDef() {
				letNames, item := p.ParseExpandedLet()
				wasteName := getWastedName()
				wasteSymID := SymIDMap.Add(wasteName)
				if err := ns.Syms.Add(wasteName, item); err != nil {
					p.stopOnError(secondToken.Lineno, "Failed to add let def. to symbol table (%s)", wasteName)
				}

				for letSymIndex, letName := range letNames {
					wasteSymbolItem := &Item{Type: SymbolPathItem, Data: SymbolPath{wasteSymID}}
					indexItem := &Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: letSymIndex}}
					opc := OpCall{OperID: IndOP, Operands: []*Item{wasteSymbolItem, indexItem}}
					indOPcallItem := &Item{Type: OperCallItem, Data: opc}
					if err := ns.Syms.Add(letName, indOPcallItem); err != nil {
						p.stopOnError(secondToken.Lineno, "Failed to add let def. to symbol table (%s)", letName)
					}
				}

				// read expander token too
				expToken, hasAny := p.tokenIter.lookAhead()
				if !hasAny {
					p.stopOnError(nil, "assuming expander, found nothing")
				}
				if expToken.Type != tokenExpander {
					p.stopOnError(expToken.Lineno, "assuming expander")
				}
				p.tokenIter.throwAway()
			} else {
				p.stopOnError(token.Lineno, "Unexpected token : %s", token.Value)
			}
		default:
			p.stopOnError(token.Lineno, "invalid let definition : %s", token.Value)
		}
	}
	DebugPrint("namespace end: %s", nsName)
	return
}

// Parse calls token iterator to scan source text and based
// on tokens it forms syntax tree
func (p *Parser) Parse(source string) (nsName string, ns *NSpace, err error) {
	err = p.tokenIter.tokenize(source)
	if err != nil {
		return
	}

	if false {
		DebugPrint("")
		for {
			token, isAny := p.tokenIter.next()
			if !isAny {
				break
			}
			DebugPrint("%#v", token)
		}
	}

	// error handling for actual parsing
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	if err != nil {
		return
	}

	for {
		token, hasAny := p.tokenIter.lookAhead()
		if !hasAny {
			break
		}
		if token.Type == tokenStartNS {
			nsName, ns = p.ParseNamespace()
			return
		}
		err = fmt.Errorf("%d: No namespace found", token.Lineno)
		return
	}
	err = fmt.Errorf("invalid namespace")
	return
}

func (p *Parser) ParseOneExpression(source string) (item *Item, err error) {
	err = p.tokenIter.tokenize(source)
	if err != nil {
		return
	}

	// error handling for actual parsing
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	if err != nil {
		return
	}

	p.SetCanEndWithSymbol(true) // a bit of a hack I guess...
	item = p.ParseExpr()
	return
}

func (p *Parser) ParseArgs(source string) (argsItems []*Item, err error) {
	err = p.tokenIter.tokenize(source)
	if err != nil {
		return
	}

	// error handling for actual parsing
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	if err != nil {
		return
	}

	if tok, _ := p.tokenIter.lookAhead(); tok.Type == tokenClosingBracket {
		p.tokenIter.throwAway()
	} else {
	ArgLoop:
		for {
			token, hasAny := p.tokenIter.lookAhead()
			if !hasAny {
				p.stopOnError(nil, "Invalid operator call")
			}
			switch token.Type {
			case tokenClosingBracket:
				p.tokenIter.throwAway()
				break ArgLoop
			case tokenComma:
				p.tokenIter.throwAway()
			case tokenStartNS:
				p.stopOnError(token.Lineno, "Namespace definition not allowed in operator call")
			case tokenEndNS:
				p.stopOnError(token.Lineno, "Namespace end definition not allowed in operator call")
			default:
				argItem := p.ParseExpr()
				argsItems = append(argsItems, argItem)
			}
		}
	}
	return
}

// NewParser returns new parser instance
func NewParser(operators Operators, srcFileName *string) *Parser {
	scanner := newTokenizer(operators)
	if scanner == nil {
		return nil
	}
	tokenIterator := newTokenIter(scanner)
	if tokenIterator == nil {
		return nil
	}
	return &Parser{tokenIter: tokenIterator, operators: operators, root: Node{NodeType: NodeTypeRoot}, srcFileName: srcFileName}
}
