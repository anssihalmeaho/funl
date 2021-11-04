package funl

import (
	"fmt"
	"sync"
)

type ExtProcType struct {
	Impl       func(*Frame, []Value) Value
	IsFunction bool
}

var debugPrintOn = false

func DebugPrint(format string, args ...interface{}) {
	if debugPrintOn {
		fmt.Println()
		fmt.Printf(format, args...)
	}
}

// NSAccess is access point to namespaces
type NSAccess struct {
	nsMap map[SymID]*NSTopInfo
	sync.RWMutex
}

type NSTopInfo struct {
	TopFrame         *Frame
	SymbolsEvaluated bool
}

func (nsa *NSAccess) Print() (s string) {
	nsa.RLock()
	defer nsa.RUnlock()

	for k, v := range nsa.nsMap {
		s += fmt.Sprintf("\n  ns: %s, frame: %v, ", SymIDMap.AsString(k), v.TopFrame.Syms.Print(0))
		for imp := range v.TopFrame.Imported {
			s += fmt.Sprintf("im: %s, ", SymIDMap.AsString(imp))
		}
	}
	return
}

func (nsa *NSAccess) Put(sid SymID, frame *Frame) {
	nsa.Lock()
	defer nsa.Unlock()

	nsa.nsMap[sid] = &NSTopInfo{TopFrame: frame}
}

func (nsa *NSAccess) GetTopFrameBySID(sid SymID) (frame *Frame, found bool) {
	nsa.RLock()
	defer nsa.RUnlock()

	var topinfo *NSTopInfo
	topinfo, found = nsa.nsMap[sid]
	if !found {
		return
	}
	frame = topinfo.TopFrame
	return
}

func (nsa *NSAccess) HasNS(sid SymID) bool {
	nsa.RLock()
	defer nsa.RUnlock()

	_, found := nsa.nsMap[sid]
	return found
}

func (nsa *NSAccess) FillFromAstNSpaceAndStore(frame *Frame, sid SymID, nspace *NSpace) {
	evalAndAssignValuesForSymbolsInFrameForNS(frame, nspace)
	topInfo := &NSTopInfo{
		TopFrame:         frame,
		SymbolsEvaluated: true, // is this useless..?
	}

	nsa.Lock()
	defer nsa.Unlock()
	nsa.nsMap[sid] = topInfo
}

func newTopFrameForNS(ns *NSpace, interpreter *Interpreter) *Frame {
	return &Frame{
		Syms:        NewSymt(),
		OtherNS:     ns.OtherNS,
		Imported:    make(map[SymID]*Frame),
		Interpreter: interpreter,
	}
}

func AddImportsToNamespace(nspace *NSpace, frame *Frame, interpreter *Interpreter) {
	AddImportsToNamespaceSub(nspace, frame, interpreter)
}

func AddImportsToNamespaceSub(nspace *NSpace, frame *Frame, interpreter *Interpreter) {
	for sid, importInfo := range nspace.OtherNS {
		importedFrame, found := interpreter.NsDir.GetTopFrameBySID(sid)
		if !found {
			var err error
			importedFrame, found, err = readModuleFromFile(frame.inProcCall, sid, importInfo.importPath, interpreter)
			if err != nil {
				importedFrame, found, err = readExtModuleFromFile(sid, importInfo.importPath, interpreter)
			}
			if err != nil {
				runTimeError2(frame, err.Error())
			}
			if !found {
				runTimeError2(frame, "Namespace top frame not found: %s", SymIDMap.AsString(sid))
			}
		}
		frame.Imported[sid] = importedFrame
	}
}

func EvalAndAssignValuesForSymbolsInFrameForNS(frame *Frame, ns *NSpace) {
	evalAndAssignValuesForSymbolsInFrameForNS(frame, ns)
}

func evalAndAssignValuesForSymbolsInFrameForNS(frame *Frame, ns *NSpace) {
	symbolMap := ns.Syms.AsMap()
	for _, sid := range ns.Syms.Keys() {
		item := symbolMap[sid]
		switch item.Type {
		case ValueItem:
			vi := item.Data.(Value)
			switch vi.Kind {
			case IntValue, FloatValue, StringValue, BoolValue, ListValue:
				if !frame.Syms.AddBySID(sid, item) {
					runTimeError2(frame, "Symbol adding failed")
				}
			case FuncProtoValue:
				newVal := Value{
					Kind: FunctionValue,
					Data: FuncValue{
						FuncProto:  vi.Data.(*Function),
						AccessLink: frame,
					},
				}
				newItem := &Item{
					Type: ValueItem,
					Data: newVal,
				}
				if !frame.Syms.AddBySID(sid, newItem) {
					runTimeError2(frame, "Symbol adding failed")
				}
				subframe := &Frame{
					Syms:     NewSymt(),
					OtherNS:  make(map[SymID]ImportInfo),
					Imported: make(map[SymID]*Frame),
				}
				fillSymsForFunc(subframe.Syms, vi.Data.(*Function))
			}
		}
	}
	// new round to evaluate symbols and operator calls
	symbolMap = ns.Syms.AsMap()
	for _, sid := range ns.Syms.Keys() {
		item := symbolMap[sid]
		switch item.Type {
		case ValueItem:
			// do nothing for values
		case SymbolPathItem, OperCallItem:
			evaluatedItem := &Item{Type: ValueItem, Data: EvalItem(item, frame)}
			if !frame.Syms.AddBySID(sid, evaluatedItem) {
				runTimeError2(frame, "Symbol adding failed")
			}
		}
	}
}

func fillSymsForFunc(syms *Symt, funcProto *Function) {
	symbolMap := funcProto.NSpace.Syms.AsMap()
	for _, sid := range funcProto.NSpace.Syms.Keys() {
		item := symbolMap[sid]
		switch item.Type {
		case ValueItem:
			vi := item.Data.(Value)
			switch vi.Kind {
			case IntValue, FloatValue, StringValue, BoolValue, ListValue:
				if !syms.AddBySID(sid, item) {
					runTimeError("Symbol adding failed (%s)", SymIDMap.AsString(sid))
				}
			case FuncProtoValue:
				subframe := &Frame{
					Syms:     NewSymt(),
					OtherNS:  make(map[SymID]ImportInfo),
					Imported: make(map[SymID]*Frame),
				}
				fillSymsForFunc(subframe.Syms, vi.Data.(*Function))
			default:
				runTimeError("Data corrupted")
			}
		case SymbolPathItem, OperCallItem:
			if !syms.AddBySID(sid, item) {
				runTimeError("Symbol adding failed")
			}
		default:
			runTimeError("Data corrupted")
		}
	}
}

func GetArgs(argsAsStr string) (argsItems []*Item, err error) {
	parser := NewParser(NewDefaultOperators(), nil)
	argsParsed, err1 := parser.ParseArgs(argsAsStr + ")")
	if err1 != nil {
		err = fmt.Errorf("Args parsing failed: %v", err1)
		return
	}
	for _, a := range argsParsed {
		argsItems = append(argsItems, a)
	}
	return
}

var initExtensions []func(*Interpreter) error

// AddExtensionInitializer can be used for registering initializer
// for some extension module (registered in init -function)
func AddExtensionInitializer(initializer func(*Interpreter) error) {
	initExtensions = append(initExtensions, initializer)
}

type Interpreter struct {
	NsDir *NSAccess
}

func NewInterpreter() *Interpreter {
	interpreter := &Interpreter{
		NsDir: &NSAccess{nsMap: make(map[SymID]*NSTopInfo)},
	}
	return interpreter
}

func FunlMainWithArgs(content string, argsItems []*Item, name, srcFileName string, initSTD func(*Interpreter) error) (retValue Value, err error) {
	parser := NewParser(NewDefaultOperators(), &srcFileName)
	var nsName string
	var nspace *NSpace
	nsName, nspace, err = parser.Parse(string(content))
	if err != nil {
		return
	}

	interpreter := NewInterpreter()

	// first create top frame for namespace and put to nsDir
	topframe := newTopFrameForNS(nspace, interpreter)
	nsSid := SymIDMap.Add(nsName)

	if err = initSTD(interpreter); err != nil {
		runTimeError("Error in std-lib init (%v)", err)
	}

	if err = initFunSourceSTD(interpreter); err != nil {
		runTimeError("Error in std-lib (fun source) init (%v)", err)
	}

	// call possible extension inits
	for _, initializer := range initExtensions {
		if err := initializer(interpreter); err != nil {
			runTimeError("Error in extension module init (%v)", err)
		}
	}

	// then put imports to all namespaces
	topframe.inProcCall = true // NOTE. this was added later as otherwise proc calls failed at main level
	AddImportsToNamespace(nspace, topframe, interpreter)

	// then evaluate and assign symbols of namespaces
	interpreter.NsDir.FillFromAstNSpaceAndStore(topframe, nsSid, nspace)

	mainSid, found := SymIDMap.Get("main")
	if !found {
		runTimeError("Main module not found")
	}
	initFrame, found := interpreter.NsDir.GetTopFrameBySID(mainSid)
	if !found {
		runTimeError("Main frame not found")
	}
	// lets set main function
	mainItem, mainfound := initFrame.Syms.GetByName(name)
	if !mainfound {
		runTimeError("main function not found (%s)", name)
	}
	if mainItem.Type != ValueItem {
		runTimeError("main function not found (%s)", "not function value")
	}
	mainVal, isok := mainItem.Data.(Value)
	if !isok {
		runTimeError("cannot access main function")
	}
	if mainVal.Kind != FunctionValue {
		runTimeError("main function not found (%s)", "cannot convert to value")
	}
	_, ok := mainVal.Data.(FuncValue)
	if !ok {
		runTimeError("main function not found (%s)", "cannot convert to function")
	}

	// lets create root -function which calls actually main
	mainFuncValueAsArgList := []*Item{mainItem}
	if len(argsItems) > 0 {
		mainFuncValueAsArgList = append(mainFuncValueAsArgList, argsItems...)
	}
	rootProto := Function{IsProc: true, Body: &Item{Type: OperCallItem, Data: OpCall{OperID: CallOP, Operands: mainFuncValueAsArgList}}}
	rootFunc := FuncValue{FuncProto: &rootProto, AccessLink: initFrame}
	initFrame.FuncProto = rootFunc.FuncProto
	initFrame.inProcCall = true

	retValue = evalMain(initFrame)
	return
}

var stdfunMap = map[string]string{}

// GetReplCode return REPL code
func GetReplCode() string {
	return stdfunMap["repl"]
}

// InitFunSourceSTD provided for initializing std externally
func InitFunSourceSTD(interpreter *Interpreter) (err error) {
	return initFunSourceSTD(interpreter)
}

func initFunSourceSTD(interpreter *Interpreter) (err error) {
	type funModInfo struct {
		name    string
		content string
	}
	var funmodNames = []string{
		"stdfu",
		"stdset",
		"stddbc",
		"stdfilu", //note. this needs stdfu
		"stdser",
		"stdmeta",
		"stdpp",
		"stdpr",
		"stdsort",
	}
	for _, funmodName := range funmodNames {
		err = AddFunModToNamespace(funmodName, []byte(stdfunMap[funmodName]), interpreter)
		if err != nil {
			return
		}
	}
	return
}
