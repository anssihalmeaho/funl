package funl

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"funlang/pmap"
	"hash/fnv"
)

//PMap is persistent map
type PMap struct {
	Rbm       *pmap.RBMap
	itemCount int
	delCount  int
}

func areEqualMaps(m1, m2 *PMap) bool {
	/*	NOTE. this check is removed as deletion marking makes it invalid
		if m1.Rbm == m2.Rbm {
			return true
		}
	*/
	v1 := &Item{Type: ValueItem, Data: Value{Kind: MapValue, Data: m1}}
	v2 := &Item{Type: ValueItem, Data: Value{Kind: MapValue, Data: m2}}
	vl1 := handleLenOP(nil, []*Item{v1})
	vl2 := handleLenOP(nil, []*Item{v2})
	l1 := vl1.Data.(int)
	l2 := vl2.Data.(int)

	if (l1 == 0) && (l2 == 0) {
		return true
	}
	if (l1 == 0) || (l2 == 0) {
		return false
	}
	if l1 != l2 {
		return false
	}

	kvm1 := make(map[pmap.MKey]pmap.MValue)
	kvm2 := make(map[pmap.MKey]pmap.MValue)
	getVisitor := func(kvm *map[pmap.MKey]pmap.MValue) func(node *pmap.Node) {
		return func(node *pmap.Node) {
			if node.Val.(NodeValue).Deleted {
				return
			}
			(*kvm)[node.Key] = node.Val
		}
	}
	m1.Rbm.VisitAll(getVisitor(&kvm1))
	m2.Rbm.VisitAll(getVisitor(&kvm2))

	equalValues := func(kval1, kval2 KeyVal) bool {
		keyitem1 := &Item{Type: ValueItem, Data: kval1.Key}
		keyitem2 := &Item{Type: ValueItem, Data: kval2.Key}
		retv := handleEqOP(nil, []*Item{keyitem1, keyitem2})
		if same := retv.Data.(bool); same {
			valitem1 := &Item{Type: ValueItem, Data: kval1.Val}
			valitem2 := &Item{Type: ValueItem, Data: kval2.Val}
			retvv := handleEqOP(nil, []*Item{valitem1, valitem2})
			return retvv.Data.(bool)
		}
		return false
	}

	isSameKV := func(kv1, kv2 pmap.MValue) bool {
		// lets first handle usual case that there is just one key-value
		if lkv1, lkv2 := len(kv1.(NodeValue).SameKeyValues), len(kv2.(NodeValue).SameKeyValues); (lkv1 + lkv2) == 0 {
			return equalValues(kv1.(NodeValue).Val, kv2.(NodeValue).Val)
		} else if lkv1 != lkv2 {
			return false // overflow list length needs to be same
		}

		// this is for rare case that there would several items for same hash-value
		keyvalues1 := []KeyVal{kv1.(NodeValue).Val}
		for _, v := range kv1.(NodeValue).SameKeyValues {
			keyvalues1 = append(keyvalues1, v)
		}
		keyvalues2 := []KeyVal{kv2.(NodeValue).Val}
		for _, v := range kv2.(NodeValue).SameKeyValues {
			keyvalues2 = append(keyvalues2, v)
		}

		findEqual := func(tofound KeyVal, from []KeyVal) bool {
			for _, onev := range from {
				if equalValues(tofound, onev) {
					return true
				}
			}
			return false
		}

		for _, keyval1 := range keyvalues1 {
			if !findEqual(keyval1, keyvalues2) {
				return false
			}
		}
		return true
	}

	for k, v := range kvm1 {
		if v2, found := kvm2[k]; !found {
			return false
		} else if !isSameKV(v, v2) {
			return false
		}
	}
	return true
}

func (pm PMap) String() string {
	s := "map("
	first := true
	visitor := func(node *pmap.Node) {
		item, ok := node.Val.(NodeValue)
		if !ok {
			runTimeError("not abe to convert map item")
		}
		if item.Deleted {
			return
		}
		if !first {
			s += ", "
		}
		first = false
		s += fmt.Sprintf("%#v : %#v", item.Val.Key, item.Val.Val)
		for _, nextitem := range item.SameKeyValues {
			s += fmt.Sprintf("%#v : %#v, ", nextitem.Key, nextitem.Val)
		}
	}
	pm.Rbm.VisitAll(visitor)
	return s + ")"
}

func (pm PMap) GoString() string {
	return pm.String()
}

type KeyVal struct {
	Key Value
	Val Value
}

type NodeValue struct {
	Val           KeyVal
	SameKeyValues []KeyVal
	Deleted       bool
}

func NewNodeValue(key, val Value) NodeValue {
	return NodeValue{Val: KeyVal{Key: key, Val: val}}
}

type NodeHandler struct{}

func isEqualItems(item1, item2 *Item) bool {
	eqResult := handleEqOP(nil, []*Item{item1, item2})
	if eqResult.Kind != BoolValue {
		runTimeError("Invalid result from eq")
	}
	return eqResult.Data.(bool)
}

func (nh *NodeHandler) MarkDeletion(srcNode *pmap.Node, key pmap.MKey, actualKey pmap.MValue) (trgNode *pmap.Node, keyFound bool) {
	nval, convok := srcNode.Val.(NodeValue)
	if !convok {
		runTimeError("invalid value (%v)", srcNode.Val)
	}

	if !nval.Deleted {
		keyInNode := &Item{Type: ValueItem, Data: srcNode.Val.(NodeValue).Val.Key}
		actualKeyItem := &Item{Type: ValueItem, Data: actualKey.(Value)}
		samekeysLen := len(nval.SameKeyValues)

		// case when there is only one value there -> can be marked as deleted
		if samekeysLen == 0 {
			if isEqualItems(keyInNode, actualKeyItem) {
				nval.Deleted = true
				keyFound = true
			}
		} else {
			// so there are several values, there remain at least one so it cannot be marked as deleted
			// NOTE. this should be very rare case

			// lets first check if its .Val
			// it was .Val, so lets move last of .SameKeyValues to .Val
			if isEqualItems(keyInNode, actualKeyItem) {
				// note. we knwo taht slice is not empty, so it wont panic
				moved := nval.SameKeyValues[samekeysLen-1]
				nval.SameKeyValues = nval.SameKeyValues[:samekeysLen-1]
				nval.Val = moved
				keyFound = true
			} else {
				// lets loop if any of .SameKeyValues equal to given key
				var newSameKVs []KeyVal
				for idx := range nval.SameKeyValues {
					keyItem := &Item{Type: ValueItem, Data: nval.SameKeyValues[idx].Key}
					if isEqualItems(keyItem, actualKeyItem) {
						keyFound = true
						break // assuming key is unique
					} else {
						newSameKVs = append(newSameKVs, nval.SameKeyValues[idx])
					}
				}
				if keyFound {
					nval.SameKeyValues = newSameKVs
				}
			}
		}
	}
	trgNode = &pmap.Node{
		Left:  srcNode.Left,
		Right: srcNode.Right,
		Color: srcNode.Color,
		Key:   srcNode.Key,
		Val:   nval,
	}
	return
}

func (nh *NodeHandler) HandleSameKey(srcNode *pmap.Node, key pmap.MKey, val pmap.MValue) (trgNode *pmap.Node) {
	trgNode = &pmap.Node{
		Left:  srcNode.Left,
		Right: srcNode.Right,
		Color: srcNode.Color,
		Key:   srcNode.Key,
		Val:   srcNode.Val,
	}
	keyvalues := []KeyVal{srcNode.Val.(NodeValue).Val}
	for _, v := range srcNode.Val.(NodeValue).SameKeyValues {
		keyvalues = append(keyvalues, v)
	}

	nval, convok := val.(NodeValue)
	if !convok {
		runTimeError("invalid value (%v)", val)
	}
	newKeyVal := nval.Val
	newKeyItem := &Item{
		Type: ValueItem,
		Data: newKeyVal.Key,
	}

	// lets first check if this is marked as deleted and could be reused
	if srcNode.Val.(NodeValue).Deleted {
		// just sanity checking
		if l := len(srcNode.Val.(NodeValue).SameKeyValues); l != 0 {
			runTimeError("same key value list should be empty when marked as deleted (len=%d)", l)
		}

		newNV := NodeValue{
			Val:           val.(NodeValue).Val,
			SameKeyValues: []KeyVal{},
			Deleted:       false,
		}
		trgNode.Val = newNV
		return
	}

	// lets check if there is same value by using eq -operator
	for _, v := range keyvalues {
		argsForEq := []*Item{
			&Item{
				Type: ValueItem,
				Data: v.Key,
			},
			newKeyItem,
		}
		eqResult := handleEqOP(nil, argsForEq)
		if eqResult.Kind != BoolValue {
			runTimeError("Invalid result from eq")
		}
		if eqResult.Data.(bool) == true {
			runTimeError("Key already exists (key: %d) (val: %v)", key, val)
		}
	}
	// ok, not duplicate found, lets add to slice
	newNodeVal := NodeValue{
		Val:           srcNode.Val.(NodeValue).Val,
		SameKeyValues: append(srcNode.Val.(NodeValue).SameKeyValues, newKeyVal),
	}
	trgNode.Val = newNodeVal
	return
}

func getMatchingValue(nval *NodeValue, keyVal Value) (retVal Value, found bool) {
	keyvalues := []KeyVal{nval.Val}
	for _, v := range nval.SameKeyValues {
		keyvalues = append(keyvalues, v)
	}
	keyItem := &Item{
		Type: ValueItem,
		Data: keyVal,
	}
	for _, v := range keyvalues {
		argsForEq := []*Item{
			&Item{
				Type: ValueItem,
				Data: v.Key,
			},
			keyItem,
		}
		eqResult := handleEqOP(nil, argsForEq)
		if eqResult.Kind != BoolValue {
			runTimeError("Invalid result from eq")
		}
		if eqResult.Data.(bool) == true {
			found = true
			retVal = v.Val
			return
		}
	}
	return
}

//HandleMapOP is for std usage
func HandleMapOP(frame *Frame, operands []*Item) (retVal Value) {
	return handleMapOP(frame, operands)
}

func handleMapOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "map"
	retVal.Kind = MapValue
	rbmap := pmap.NewRBMapWithHandler(&NodeHandler{})
	argCount := len(operands)
	if argCount == 0 {
		retVal.Data = &PMap{Rbm: rbmap}
		return
	}

	if (argCount % 2) != 0 {
		runTimeError2(frame, "%s: uneven amount of arguments (%d)", opName, argCount)
	}
	var mapval Value
	mapval = Value{Kind: MapValue, Data: &PMap{Rbm: rbmap}}
	for i := 0; i < argCount; i += 2 {
		// evaluate key
		v := operands[i]
		var keyv Value
		switch v.Type {
		case ValueItem:
			keyv = v.Data.(Value)
		case SymbolPathItem, OperCallItem:
			keyv = EvalItem(v, frame)
		default:
			runTimeError2(frame, "something wrong (%s)", opName)
		}

		// evaluate value
		v = operands[i+1]
		var valv Value
		switch v.Type {
		case ValueItem:
			valv = v.Data.(Value)
		case SymbolPathItem, OperCallItem:
			valv = EvalItem(v, frame)
		default:
			runTimeError2(frame, "something wrong (%s)", opName)
		}

		// lets add key-value to map
		putArgs := []*Item{
			&Item{Type: ValueItem, Data: mapval},
			&Item{Type: ValueItem, Data: keyv},
			&Item{Type: ValueItem, Data: valv},
		}
		mapval = handlePutOP(frame, putArgs)
		if mapval.Kind != MapValue {
			runTimeError2(frame, "%s: failed to put key-value to map", opName)
		}
	}
	retVal = mapval
	return
}

func getHashedBool(from bool) int {
	var asint int
	if from {
		asint = 1
	}

	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(asint))
	hash := fnv.New64a()
	_, err := hash.Write(b)
	if err != nil {
		runTimeError("Hash error (%v)", err)
	}
	return int(hash.Sum64())

}

func getHashedInt(from int) int {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(from))
	hash := fnv.New64a()
	_, err := hash.Write(b)
	if err != nil {
		runTimeError("Hash error (%v)", err)
	}
	return int(hash.Sum64())
}

func getHashedString(from string) int {
	b := []byte(from)
	hash := fnv.New64a()
	_, err := hash.Write(b)
	if err != nil {
		runTimeError("Hash error (%v)", err)
	}
	return int(hash.Sum64())
}

func getHashedFloat(from float64) int {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, from)

	hash := fnv.New64a()
	_, err := hash.Write(buf.Bytes())
	if err != nil {
		runTimeError("Hash error (%v)", err)
	}
	return int(hash.Sum64())
}

func getHashedList(from *List) int {
	buf := new(bytes.Buffer)
	lit := NewListIterator(Value{Kind: ListValue, Data: from})
	if lit == nil {
		runTimeError("Unable to hash list")
	}
	for {
		nextitem := lit.Next()
		if nextitem == nil {
			break
		}
		hashedKey, err := hashOfValue(*nextitem)
		if err != nil {
			runTimeError("illegal type for map key")
		}

		int64v := int64(hashedKey)
		err = binary.Write(buf, binary.LittleEndian, int64v)
		if err != nil {
			runTimeError("Hash error (%v)", err)
		}
	}
	hash := fnv.New64a()
	hash.Write(buf.Bytes())
	return int(hash.Sum64())
}

func hashOfValue(keyVal Value) (hashedKey int, err error) {
	switch keyVal.Kind {
	case IntValue:
		hashedKey = getHashedInt(keyVal.Data.(int))
	case StringValue:
		hashedKey = getHashedString(keyVal.Data.(string))
	case FloatValue:
		hashedKey = getHashedFloat(keyVal.Data.(float64))
	case ListValue:
		hashedKey = getHashedList(keyVal.Data.(*List))
	case BoolValue:
		hashedKey = getHashedBool(keyVal.Data.(bool))
	default:
		err = fmt.Errorf("illegal type for map key")
	}
	return
}

func handleKeysOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "keys"
	if l := len(operands); l != 1 {
		runTimeError2(frame, "%s operator needs one argument (%d given)", opName, l)
	}

	var mapVal Value
	switch v := operands[0]; v.Type {
	case ValueItem:
		mapVal = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		mapVal = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	if mapVal.Kind != MapValue {
		runTimeError2(frame, "First argument not map in %s operator", opName)
	}
	mapv, convok := mapVal.Data.(*PMap)
	if !convok {
		runTimeError2(frame, "First argument is not map in %s operator", opName)
	}

	var keyItems []*Item
	visitor := func(node *pmap.Node) {
		item, ok := node.Val.(NodeValue)
		if !ok {
			runTimeError2(frame, "%s: not able to convert map item", opName)
		}
		if !item.Deleted {
			keyItems = append(keyItems, &Item{Type: ValueItem, Data: item.Val.Key})
			for _, nextitem := range item.SameKeyValues {
				keyItems = append(keyItems, &Item{Type: ValueItem, Data: nextitem.Key})
			}
		}
	}
	mapv.Rbm.VisitAll(visitor)
	retVal = handleListOP(frame, keyItems)
	return
}

func handleValsOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "vals"
	if l := len(operands); l != 1 {
		runTimeError2(frame, "%s operator needs one argument (%d given)", opName, l)
	}

	var mapVal Value
	switch v := operands[0]; v.Type {
	case ValueItem:
		mapVal = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		mapVal = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	if mapVal.Kind != MapValue {
		runTimeError2(frame, "First argument not map in %s operator", opName)
	}
	mapv, convok := mapVal.Data.(*PMap)
	if !convok {
		runTimeError2(frame, "First argument is not map in %s operator", opName)
	}

	var keyItems []*Item
	visitor := func(node *pmap.Node) {
		item, ok := node.Val.(NodeValue)
		if !ok {
			runTimeError2(frame, "%s: not able to convert map item", opName)
		}
		if !item.Deleted {
			keyItems = append(keyItems, &Item{Type: ValueItem, Data: item.Val.Val})
			for _, nextitem := range item.SameKeyValues {
				keyItems = append(keyItems, &Item{Type: ValueItem, Data: nextitem.Val})
			}
		}
	}
	mapv.Rbm.VisitAll(visitor)
	retVal = handleListOP(frame, keyItems)
	return
}

//HandleKeyvalsOP is for std usage
func HandleKeyvalsOP(frame *Frame, operands []*Item) (retVal Value) {
	return handleKeyvalsOP(frame, operands)
}

func handleKeyvalsOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "keyvals"
	if l := len(operands); l != 1 {
		runTimeError2(frame, "%s operator needs one argument (%d given)", opName, l)
	}

	var mapVal Value
	switch v := operands[0]; v.Type {
	case ValueItem:
		mapVal = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		mapVal = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	if mapVal.Kind != MapValue {
		runTimeError2(frame, "First argument not map in %s operator", opName)
	}
	mapv, convok := mapVal.Data.(*PMap)
	if !convok {
		runTimeError2(frame, "First argument is not map in %s operator", opName)
	}

	var keyItems []*Item
	visitor := func(node *pmap.Node) {
		item, ok := node.Val.(NodeValue)
		if !ok {
			runTimeError2(frame, "%s: not able to convert map item", opName)
		}
		if !item.Deleted {
			kvpair := []*Item{
				&Item{Type: ValueItem, Data: item.Val.Key},
				&Item{Type: ValueItem, Data: item.Val.Val},
			}
			keyItems = append(keyItems, &Item{Type: ValueItem, Data: handleListOP(frame, kvpair)})
			for _, nextitem := range item.SameKeyValues {
				kvpair := []*Item{
					&Item{Type: ValueItem, Data: nextitem.Key},
					&Item{Type: ValueItem, Data: nextitem.Val},
				}
				keyItems = append(keyItems, &Item{Type: ValueItem, Data: handleListOP(frame, kvpair)})
			}
		}
	}
	mapv.Rbm.VisitAll(visitor)
	retVal = handleListOP(frame, keyItems)
	return
}

func handleGetlOP(frame *Frame, operands []*Item) (retVal Value) {
	argCount := len(operands)
	if argCount != 2 {
		opName := "getl"
		runTimeError2(frame, "%s operator needs two arguments (%d given)", opName, argCount)
	}
	retVal = commonGetOP("getl", argCount, true, frame, operands)
	return
}

func handleGetOP(frame *Frame, operands []*Item) (retVal Value) {
	argCount := len(operands)
	if argCount != 2 && argCount != 3 {
		opName := "get"
		runTimeError2(frame, "%s operator needs two or three arguments (%d given)", opName, argCount)
	}
	retVal = commonGetOP("get", argCount, false, frame, operands)
	return
}

func commonGetOP(opName string, argCount int, isGetl bool, frame *Frame, operands []*Item) (retVal Value) {
	var mapVal Value
	switch v := operands[0]; v.Type {
	case ValueItem:
		mapVal = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		mapVal = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	if mapVal.Kind != MapValue {
		runTimeError2(frame, "First argument not map in %s operator", opName)
	}
	mapv, convok := mapVal.Data.(*PMap)
	if !convok {
		runTimeError2(frame, "First argument is not map in %s operator", opName)
	}

	var keyVal Value
	switch v := operands[1]; v.Type {
	case ValueItem:
		keyVal = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		keyVal = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	hashedKey, err := hashOfValue(keyVal)
	if err != nil {
		runTimeError2(frame, "%s: illegal type for map key", opName)
	}

	nodeval, found := mapv.Rbm.Get(pmap.MKey(hashedKey))

	if found {
		nval, convok := nodeval.(NodeValue)
		if !convok {
			runTimeError2(frame, "%s: invalid value (%v)", opName, keyVal.Data)
		}
		if nval.Deleted {
			found = false
		}
	}

	if !found {
		if isGetl {
			argsForReturnList := []*Item{
				&Item{
					Type: ValueItem,
					Data: Value{Kind: BoolValue, Data: false},
				},
				&Item{
					Type: ValueItem,
					Data: Value{Kind: BoolValue, Data: false},
				},
			}
			retVal = handleListOP(frame, argsForReturnList)
			return
		}
		if argCount == 3 {
			switch v := operands[2]; v.Type {
			case ValueItem:
				retVal = v.Data.(Value)
			case SymbolPathItem, OperCallItem:
				retVal = EvalItem(v, frame)
			default:
				runTimeError2(frame, "something wrong (%s)", opName)
			}
			return
		}
		runTimeError2(frame, "%s: key not found (%v)", opName, keyVal.Data)
	}
	nval, convok := nodeval.(NodeValue)
	if !convok {
		runTimeError2(frame, "%s: invalid value (%v)", opName, keyVal.Data)
	}
	if len(nval.SameKeyValues) == 0 {
		//lets check still that keys are actually equal, although hashes are...
		argsForEq := []*Item{
			&Item{
				Type: ValueItem,
				Data: nval.Val.Key,
			},
			&Item{
				Type: ValueItem,
				Data: keyVal,
			},
		}
		eqResult := handleEqOP(nil, argsForEq)
		if eqResult.Kind != BoolValue {
			runTimeError2(frame, "Invalid result from eq")
		}
		if eqResult.Data.(bool) == false {
			// not same
			goto NotFound
		}
		// keys are matching, return value
		if isGetl {
			argsForReturnList := []*Item{
				&Item{
					Type: ValueItem,
					Data: Value{Kind: BoolValue, Data: true},
				},
				&Item{
					Type: ValueItem,
					Data: nval.Val.Val,
				},
			}
			retVal = handleListOP(frame, argsForReturnList)
		} else {
			retVal = nval.Val.Val
		}
		return
	}
	retVal, found = getMatchingValue(&nval, keyVal)
	if found {
		if isGetl {
			argsForReturnList := []*Item{
				&Item{
					Type: ValueItem,
					Data: Value{Kind: BoolValue, Data: true},
				},
				&Item{
					Type: ValueItem,
					Data: retVal,
				},
			}
			retVal = handleListOP(frame, argsForReturnList)
		}
		return
	}
NotFound:
	if isGetl {
		argsForReturnList := []*Item{
			&Item{
				Type: ValueItem,
				Data: Value{Kind: BoolValue, Data: false},
			},
		}
		retVal = handleListOP(frame, argsForReturnList)
		return
	}
	if argCount == 3 {
		switch v := operands[2]; v.Type {
		case ValueItem:
			retVal = v.Data.(Value)
		case SymbolPathItem, OperCallItem:
			retVal = EvalItem(v, frame)
		default:
			runTimeError2(frame, "something wrong (%s)", opName)
		}
		return
	}
	runTimeError2(frame, "%s: key not found (%v)", opName, keyVal.Data)
	return
}

func handleDellOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "dell"
	newMapVal, keyFound := delCommon(false, opName, frame, operands)
	argsForReturnList := []*Item{
		&Item{
			Type: ValueItem,
			Data: Value{Kind: BoolValue, Data: keyFound},
		},
		&Item{
			Type: ValueItem,
			Data: newMapVal,
		},
	}
	retVal = handleListOP(frame, argsForReturnList)
	return
}

func handleDelOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "del"
	retVal, _ = delCommon(true, opName, frame, operands)
	return
}

// DeletionsLimit is limit for amount of items marked as deleted in map, before making new map
const DeletionsLimit = 1000

func delCommon(rteIfNotExist bool, opName string, frame *Frame, operands []*Item) (retVal Value, keyFound bool) {
	retVal.Kind = MapValue

	if l := len(operands); l != 2 {
		runTimeError2(frame, "%s operator needs two arguments (%d given)", opName, l)
	}

	var mapVal Value
	switch v := operands[0]; v.Type {
	case ValueItem:
		mapVal = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		mapVal = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	if mapVal.Kind != MapValue {
		runTimeError2(frame, "First argument not map in %s operator", opName)
	}
	mapv, convok := mapVal.Data.(*PMap)
	if !convok {
		runTimeError2(frame, "First argument is not map in %s operator", opName)
	}

	var keyVal Value
	switch v := operands[1]; v.Type {
	case ValueItem:
		keyVal = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		keyVal = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	if (mapv.itemCount < 1) && rteIfNotExist {
		runTimeError2(frame, "empty map for %s operator", opName)
	}

	hashedKey, err := hashOfValue(keyVal)
	if err != nil {
		runTimeError2(frame, "%s: illegal type for map key", opName)
	}

	newmap, itemFound := mapv.Rbm.Modify(pmap.MKey(hashedKey), keyVal)
	var newCount int
	var newDelCount int
	if itemFound {
		newCount = mapv.itemCount - 1
		newDelCount = mapv.delCount + 1
	} else {
		if rteIfNotExist {
			runTimeError2(frame, "%s: key not found", opName)
		}
		newCount = mapv.itemCount
		newDelCount = mapv.delCount
	}

	newMapVal := &PMap{Rbm: newmap, itemCount: newCount, delCount: newDelCount}

	// lets check if there are so many deleted ones that new map should be created
	if newMapVal.delCount > DeletionsLimit {
		var keyItems []*Item
		var deletedCount int
		visitor := func(node *pmap.Node) {
			item, ok := node.Val.(NodeValue)
			if !ok {
				runTimeError2(frame, "%s: not able to convert map item", opName)
			}
			if !item.Deleted {
				keyItems = append(keyItems, &Item{Type: ValueItem, Data: item.Val.Key})
				keyItems = append(keyItems, &Item{Type: ValueItem, Data: item.Val.Val})
				for _, nextitem := range item.SameKeyValues {
					keyItems = append(keyItems, &Item{Type: ValueItem, Data: nextitem.Key})
					keyItems = append(keyItems, &Item{Type: ValueItem, Data: nextitem.Val})
				}
			} else {
				deletedCount++
			}
		}
		newMapVal.Rbm.VisitAll(visitor)

		// check if amount of delted really exceeds the limit
		if deletedCount > DeletionsLimit {
			retVal = handleMapOP(frame, keyItems)
			keyFound = itemFound
			return
		}
		newMapVal.delCount = 0 // lets clear del counter so that traversal is not done all the time
	}

	retVal.Data = newMapVal
	keyFound = itemFound
	return
}

//HandlePutOP is for std usage
func HandlePutOP(frame *Frame, operands []*Item) (retVal Value) {
	return handlePutOP(frame, operands)
}

func handlePutOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "put"
	retVal.Kind = MapValue

	if l := len(operands); l != 3 {
		runTimeError2(frame, "%s operator needs three arguments (%d given)", opName, l)
	}

	var mapVal Value
	switch v := operands[0]; v.Type {
	case ValueItem:
		mapVal = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		mapVal = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	if mapVal.Kind != MapValue {
		runTimeError2(frame, "First argument not map in %s operator", opName)
	}
	mapv, convok := mapVal.Data.(*PMap)
	if !convok {
		runTimeError2(frame, "First argument is not map in %s operator", opName)
	}

	var keyVal Value
	switch v := operands[1]; v.Type {
	case ValueItem:
		keyVal = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		keyVal = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	hashedKey, err := hashOfValue(keyVal)
	if err != nil {
		runTimeError2(frame, "%s: illegal type for map key", opName)
	}

	var valueVal Value
	switch v := operands[2]; v.Type {
	case ValueItem:
		valueVal = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		valueVal = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	nodeval := NewNodeValue(keyVal, valueVal)
	newmap := mapv.Rbm.Put(pmap.MKey(hashedKey), nodeval)
	retVal.Data = &PMap{Rbm: newmap, itemCount: mapv.itemCount + 1, delCount: mapv.delCount}
	return
}
