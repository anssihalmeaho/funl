package funl

import (
	"fmt"
	"math"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

type ValueType int

type ItemType int

type OperType int

const (
	IntValue ValueType = iota
	StringValue
	BoolValue
	FuncProtoValue
	FunctionValue
	ListValue
	ExtProcValue
	ChanValue
	FloatValue
	OpaqueValue
	MapValue

	ValueItem ItemType = iota
	SymbolPathItem
	OperCallItem

	AndOP OperType = iota
	OrOP
	NotOP
	CallOP
	EqOP
	IfOP
	PlusOP
	MinusOP
	MulOP
	DivOP
	ModOP
	ListOP
	EmptyOP
	HeadOP
	LastOP
	RestOP
	AppendOP
	AddOP
	LenOP
	TypeOP
	InOP
	IndOP
	FindOP
	SliceOP
	RrestOP
	ReverseOP
	ExtendOP
	SplitOP
	GtOP
	LtOP
	LeOP
	GeOP
	StrOP
	ConvOP
	CaseOP
	NameOP
	ErrorOP
	PrintOP
	SpawnOP
	ChanOP
	SendOP
	RecvOP
	SymvalOP
	TryOP
	TrylOP
	SelectOP
	EvalOP
	WhileOP
	FloatOP
	MapOP
	PutOP
	GetOP
	GetlOP
	KeysOP
	ValsOP
	KeyvalsOP
	LetOP
	ImpOP
	DelOP
	DellOP
	SprintfOP
	ArgslistOP
	CondOP
	HelpOP
	RecwithOP
	MaximumOP
)

//OpaqueAPI is interface for opaque type
type OpaqueAPI interface {
	TypeName() string
	Str() string
	Equals(with OpaqueAPI) bool
}

func (it ItemType) GoString() string {
	switch it {
	case ValueItem:
		return "ValueItem"
	case SymbolPathItem:
		return "SymbolPathItem"
	case OperCallItem:
		return "OperCallItem"
	}
	return "Unknown Item type"
}

func (it ItemType) String() string {
	return it.GoString()
}

func (vt ValueType) GoString() string {
	switch vt {
	case IntValue:
		return "Int-Value"
	case StringValue:
		return "String-Value"
	case BoolValue:
		return "Bool-Value"
	case FuncProtoValue:
		return "Func-Proto-Value"
	case FunctionValue:
		return "Func-Value"
	case ListValue:
		return "List-Value"
	}
	return "Unknown-Value"
}

func (vt ValueType) String() string {
	return vt.GoString()
}

func operTypeFromIntToString(ot OperType) string {
	str, ok := map[OperType]string{
		AndOP:      "and",
		OrOP:       "or",
		NotOP:      "not",
		CallOP:     "call",
		EqOP:       "eq",
		IfOP:       "if",
		PlusOP:     "plus",
		MinusOP:    "minus",
		MulOP:      "mul",
		DivOP:      "div",
		ModOP:      "mod",
		ListOP:     "list",
		EmptyOP:    "empty",
		HeadOP:     "head",
		LastOP:     "last",
		RestOP:     "rest",
		AppendOP:   "append",
		AddOP:      "add",
		LenOP:      "len",
		TypeOP:     "type",
		InOP:       "in",
		IndOP:      "ind",
		FindOP:     "find",
		SliceOP:    "slice",
		RrestOP:    "rrest",
		ReverseOP:  "reverse",
		ExtendOP:   "extend",
		SplitOP:    "split",
		GtOP:       "gt",
		LtOP:       "lt",
		LeOP:       "le",
		GeOP:       "ge",
		StrOP:      "str",
		ConvOP:     "conv",
		CaseOP:     "case",
		NameOP:     "name",
		ErrorOP:    "error",
		PrintOP:    "print",
		SpawnOP:    "spawn",
		ChanOP:     "chan",
		SendOP:     "send",
		RecvOP:     "recv",
		SymvalOP:   "symval",
		TryOP:      "try",
		TrylOP:     "tryl",
		SelectOP:   "select",
		EvalOP:     "eval",
		WhileOP:    "while",
		FloatOP:    "float",
		MapOP:      "map",
		PutOP:      "put",
		GetOP:      "get",
		GetlOP:     "getl",
		KeysOP:     "keys",
		ValsOP:     "vals",
		KeyvalsOP:  "keyvals",
		LetOP:      "let",
		ImpOP:      "imp",
		DelOP:      "del",
		DellOP:     "dell",
		SprintfOP:  "sprintf",
		ArgslistOP: "argslist",
		CondOP:     "cond",
		HelpOP:     "help",
		RecwithOP:  "recwith",
		MaximumOP:  "MAX",
	}[ot]
	if !ok {
		return fmt.Sprintf("Unknown (%d)", int(ot))
	}
	return str
}

func (ot OperType) GoString() string {
	return operTypeFromIntToString(ot)
}

func (ot OperType) String() string {
	return operTypeFromIntToString(ot)
}

type SymID int

func (sid SymID) String() string {
	return SymIDMap.AsString(sid)
}

func (sid SymID) GoString() string {
	return sid.String()
}

type SymbolToIDConverter struct {
	counter  SymID
	symIDmap map[string]SymID
	sync.RWMutex
}

func NewSymbolToIDConverter() *SymbolToIDConverter {
	return &SymbolToIDConverter{counter: 1, symIDmap: make(map[string]SymID)}
}

func (sidc *SymbolToIDConverter) SymbolCount() int {
	sidc.RLock()
	defer sidc.RUnlock()

	return int(sidc.counter)
}

func (sidc *SymbolToIDConverter) AsString(sid SymID) string {
	sidc.RLock()
	defer sidc.RUnlock()

	for k, v := range sidc.symIDmap {
		if v == sid {
			return k
		}
	}
	return ""
}

func (sidc *SymbolToIDConverter) Get(symbol string) (sid SymID, found bool) {
	sidc.RLock()
	defer sidc.RUnlock()

	sid, found = sidc.symIDmap[symbol]
	return
}

func (sidc *SymbolToIDConverter) Add(symbol string) SymID {
	sidc.Lock()
	defer sidc.Unlock()

	sid, found := sidc.symIDmap[symbol]
	if found {
		return sid
	}
	sidc.symIDmap[symbol] = sidc.counter
	sid = sidc.counter
	sidc.counter++

	if false {
		debug.PrintStack()
	}

	return sid
}

var SymIDMap = NewSymbolToIDConverter()
var AnySymSid = SymIDMap.Add("_")

var wasteCounter uint64
var wasteMutex = &sync.Mutex{}
var getWastedName func() string

func init() {
	if (strconv.IntSize == 64) && (runtime.GOARCH == "amd64") {
		getWastedName = func() string {
			atomic.AddUint64(&wasteCounter, 1)
			return fmt.Sprintf("__waste_%d", atomic.LoadUint64(&wasteCounter))
		}
	} else {
		getWastedName = func() string {
			wasteMutex.Lock()
			defer wasteMutex.Unlock()

			wasteCounter++
			return fmt.Sprintf("__waste_%d", wasteCounter)
		}
	}
}

// GetWastedName is for stdast use
func GetWastedName() string {
	return getWastedName()
}

type Symt struct {
	ordered []SymID // keys
	mapped  map[SymID]*Item
	sync.RWMutex
}

func NewSymt() *Symt {
	return &Symt{mapped: make(map[SymID]*Item)}
}

func (sym *Symt) MakeCopy() *Symt {
	sym.RLock()
	defer sym.RUnlock()

	newsyms := NewSymt()
	for k, v := range sym.mapped {
		newsyms.mapped[k] = v
	}
	for _, v := range sym.ordered {
		newsyms.ordered = append(newsyms.ordered, v)
	}
	return newsyms
}

type SymbolPath []SymID

func (sp *SymbolPath) ToString() string {
	var targetStr string
	pathLen := len(*sp)
	for ind, part := range *sp {
		targetStr = targetStr + SymIDMap.AsString(part)
		if ind < pathLen-1 {
			targetStr += "."
		}
	}
	return targetStr
}

func (sym *Symt) Has(key SymID) bool {
	sym.RLock()
	defer sym.RUnlock()

	_, found := sym.mapped[key]
	return found
}

func (sym *Symt) Keys() []SymID {
	sym.RLock()
	defer sym.RUnlock()

	return sym.ordered
}

func (sym *Symt) AsMap() map[SymID]*Item {
	sym.RLock()
	defer sym.RUnlock()

	return sym.mapped
}

func (sym *Symt) Print(depth int) (s string) {
	for k, v := range sym.mapped {
		s += fmt.Sprintf("%ssym: %s: %v\n", depthPrint(depth), SymIDMap.AsString(k), v.Print(depth))
	}
	return
}

func (sym *Symt) FindFuncSID(fptr *Function) (sid SymID, found bool) {
	sym.RLock()
	defer sym.RUnlock()

	for k, v := range sym.mapped {
		if v.Type == ValueItem {
			val := v.Data.(Value)
			if val.Kind == FunctionValue {
				fval := val.Data.(FuncValue)
				if fval.FuncProto == fptr {
					return k, true
				}
			}
		}
	}
	return
}

func (sym *Symt) GetByName(symbol string) (*Item, bool) {
	sid, found := SymIDMap.Get(symbol)
	if !found {
		return nil, false
	}
	return sym.GetBySID(sid)
}

func (sym *Symt) GetBySID(sid SymID) (*Item, bool) {
	sym.RLock()
	defer sym.RUnlock()

	item, found := sym.mapped[sid]
	return item, found
}

func (sym *Symt) AddBySIDByOverwriteIfNeeded(sid SymID, item *Item) bool {
	if sid == AnySymSid {
		sid = SymIDMap.Add(getWastedName())
	}

	sym.Lock()
	defer sym.Unlock()

	sym.mapped[sid] = item
	sym.ordered = append(sym.ordered, sid)
	return true
}

func (sym *Symt) AddBySID(sid SymID, item *Item) bool {
	if sid == AnySymSid {
		sid = SymIDMap.Add(getWastedName())
	}

	sym.Lock()
	defer sym.Unlock()

	if _, found := sym.mapped[sid]; found {
		return false
	}
	sym.mapped[sid] = item
	sym.ordered = append(sym.ordered, sid)
	return true
}

func (sym *Symt) Add(symbol string, item *Item) error {
	if symbol == "_" {
		symbol = getWastedName()
	}

	sid := SymIDMap.Add(symbol)

	sym.Lock()
	defer sym.Unlock()

	if _, found := sym.mapped[sid]; found {
		return fmt.Errorf("Symbol (%s) already found", symbol)
	}
	sym.mapped[sid] = item
	sym.ordered = append(sym.ordered, sid)
	return nil
}

type ImportInfo struct {
	importPath string // "" if not defined
}

func (imp *ImportInfo) Path() string {
	return imp.importPath
}

func (imp *ImportInfo) SetPath(path string) {
	imp.importPath = path
}

func NewDocStrings(val []string) *DocStrings {
	return &DocStrings{strList: val}
}

type DocStrings struct {
	strList []string
}

func (ds *DocStrings) AsText() string {
	if ds.strList == nil {
		return ""
	}
	return strings.Join(ds.strList, "\n")
}

type NSpace struct {
	Syms    *Symt
	OtherNS map[SymID]ImportInfo
	Docs    *DocStrings
}

func depthPrint(depth int) (s string) {
	for i := 0; i < depth; i++ {
		s = s + ".."
	}
	return
}

func (ns *NSpace) Print(depth int) (s string) {
	s += fmt.Sprintf("%sns:\n", depthPrint(depth))
	s += fmt.Sprintf("%s  syms:\n%s\n", depthPrint(depth), ns.Syms.Print(depth+1))
	s += fmt.Sprintf("%s  peer-ns: %v\n", depthPrint(depth), ns.OtherNS)
	return
}

type Item struct {
	Type             ItemType
	Data             interface{}
	Expand           bool
	ExpandArgIndexes map[int]bool
}

func (item *Item) Print(depth int) (s string) {
	switch item.Type {
	case ValueItem:
		s = "VALUE"
		vi := item.Data.(Value)
		switch vi.Kind {
		case IntValue:
			intv := vi.Data.(int)
			s = fmt.Sprintf("%d", intv)
		case StringValue:
			s = vi.Data.(string)
		case OpaqueValue:
			s = "opaque:" + vi.Data.(OpaqueAPI).Str()
		case BoolValue:
			boolv := vi.Data.(bool)
			if boolv {
				s = "true"
			} else {
				s = "false"
			}
		case FuncProtoValue:
			funcv := vi.Data.(*Function)
			s = fmt.Sprintf("func-value: file: %s line: %d pos: %d", funcv.SrcFileName, funcv.Lineno, funcv.Pos)
		case ListValue:
			s = "LIST"
		case FloatValue:
			floatv := vi.Data.(float64)
			s = fmt.Sprintf("%v", floatv)
		default:
			s = "UNKNOWN VALUE"
		}
	case SymbolPathItem:
		sp := item.Data.(SymbolPath)
		s = sp.ToString()
	case OperCallItem:
		opc := item.Data.(OpCall)
		s2 := ""
		for _, v := range opc.Operands {
			s2 += (v.Print(depth) + ", ")
		}
		s += fmt.Sprintf("op-call: %d, (operands: %s)", opc.OperID, s2)
	default:
		s = "UNKNOWN ITEM"
	}
	return
}

type Value struct {
	Kind ValueType
	Data interface{}
}

func (val Value) String() string {
	switch val.Kind {
	case IntValue:
		intv := val.Data.(int)
		return fmt.Sprintf("%d", intv)
	case FloatValue:
		floatv := val.Data.(float64)
		if math.Trunc(floatv) == floatv {
			// its whole number
			return fmt.Sprintf("%.1f", floatv)
		}
		return fmt.Sprintf("%v", floatv)
	case StringValue:
		return fmt.Sprintf("'%s'", val.Data.(string))
	case OpaqueValue:
		oif := val.Data.(OpaqueAPI)
		return fmt.Sprintf("opaque(%s)", oif.Str())
	case BoolValue:
		boolv := val.Data.(bool)
		if boolv {
			return "true"
		}
		return "false"
	case FuncProtoValue:
		funcv := val.Data.(*Function)
		return fmt.Sprintf("func-value: file: %s line: %d pos: %d", funcv.SrcFileName, funcv.Lineno, funcv.Pos)
	case FunctionValue:
		funcv := val.Data.(FuncValue)
		fproto := funcv.FuncProto
		return fmt.Sprintf("func-value: file: %s line: %d pos: %d", fproto.SrcFileName, fproto.Lineno, fproto.Pos)
	case ListValue:
		return val.Data.(*List).String()
	case MapValue:
		return val.Data.(*PMap).String()
	case ChanValue:
		return "chan-value"
	case ExtProcValue:
		return "ext-proc"
	default:
		return "UNKNOWN VALUE"
	}
}

func (val Value) GoString() string {
	return val.String()
}

type FuncValue struct {
	FuncProto  *Function
	AccessLink *Frame // nil if root
}

type OpCall struct {
	OperID   OperType
	Operands []*Item
}

type Function struct {
	IsProc      bool
	ArgNames    []SymID
	Body        *Item
	NSpace      NSpace
	Lineno      int
	Pos         int
	SrcFileName string
}

func (f *Function) Show() {
	fmt.Printf("\nFunction:\n")
	for i, sid := range f.ArgNames {
		fmt.Printf("  arg %d: %s\n", i, SymIDMap.AsString(sid))
	}
	fmt.Printf("  body: %#v\n", *f.Body)
}
