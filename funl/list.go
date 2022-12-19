package funl

import (
	"fmt"
)

type ListObject struct {
	Val  *Value
	Next *ListObject
}

type List struct {
	Head *ListObject
	Tail *ListObject
}

func (l List) String() string {
	s := "list("
	it := NewListIterator(Value{Kind: ListValue, Data: &l})
	first := true
	for {
		v := it.Next()
		if v == nil {
			return s + ")"
		} else if !first {
			s += ", "
		}
		first = false
		s += fmt.Sprintf("%#v", *v)
	}
}

func (l List) GoString() string {
	return l.String()
}

type ListIterator struct {
	NextItem *ListObject
	HeadDone bool
	Tail     *ListObject
}

func NewListIterator(val Value) (lit *ListIterator) {
	if val.Kind != ListValue {
		return
	}
	list := val.Data.(*List)
	var newlist *List
	if list.Tail != nil {
		newlist = reverseCopy(list)
	} else {
		newlist = &List{}
	}
	lit = &ListIterator{}
	if list.Head == nil {
		lit.NextItem = newlist.Head
		lit.HeadDone = true
	} else {
		lit.NextItem = list.Head
	}
	lit.Tail = newlist.Head
	return
}

func (lit *ListIterator) Next() *Value {
	if lit.NextItem == nil {
		if lit.HeadDone {
			return nil
		}
		lit.NextItem = lit.Tail
		lit.HeadDone = true
		return lit.Next()
	}
	retv := lit.NextItem.Val
	lit.NextItem = lit.NextItem.Next
	return retv
}

func handleExtendOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "extend"

	var argvals []Value
	for _, v := range operands {
		var argval Value
		switch v.Type {
		case ValueItem:
			argval = v.Data.(Value)
		case SymbolPathItem, OperCallItem:
			argval = EvalItem(v, frame)
		default:
			runTimeError2(frame, "something wrong (%s)", opName)
		}
		if argval.Kind != ListValue {
			runTimeError2(frame, "%s: arguments assumed to be list type", opName)
		}
		argvals = append(argvals, argval)
	}
	argsLen := len(argvals)

	var prevnew *ListObject
	var nextnew *ListObject
	var newHead *ListObject
	for _, v := range argvals {
		it := NewListIterator(v)
		for {
			nextitem := it.Next()
			if nextitem == nil {
				break
			}
			if nextnew != nil {
				prevnew = nextnew
			}
			nextnew = &ListObject{Val: nextitem, Next: nil}
			if newHead == nil {
				newHead = nextnew
			}
			if prevnew != nil {
				prevnew.Next = nextnew
			}
		}
	}
	switch argsLen {
	case 0:
	case 1:
		newHead = argvals[0].Data.(*List).Head
	}
	retVal = Value{Kind: ListValue, Data: &List{Head: newHead, Tail: nil}}
	return
}

func handleReverseOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "reverse"
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

	if val.Kind != ListValue {
		runTimeError2(frame, "First argument not list in %s operator", opName)
	}
	list, convok := val.Data.(*List)
	if !convok {
		runTimeError2(frame, "First argument is not list in %s operator", opName)
	}
	retVal.Kind = ListValue

	headL := &List{Tail: list.Head} // put it to tail so that reverseCopy can be used...
	var reversedHeadCopy *List
	if list.Head != nil {
		reversedHeadCopy = reverseCopy(headL)
	} else {
		reversedHeadCopy = &List{}
	}

	nextitem := list.Tail
	var newitem *ListObject
	var prevnewitem *ListObject
	var newHead *ListObject
	for {
		if nextitem == nil {
			break
		}
		if newitem != nil {
			prevnewitem = newitem
		}
		newitem = &ListObject{Val: nextitem.Val, Next: nil}
		if newHead == nil {
			newHead = newitem
		}
		if prevnewitem != nil {
			prevnewitem.Next = newitem
		}
		nextitem = nextitem.Next
	}
	if newitem != nil {
		newitem.Next = reversedHeadCopy.Head
	}
	if newHead == nil {
		newHead = reversedHeadCopy.Head
	}
	retVal.Data = &List{Head: newHead, Tail: nil}
	return
}

func handleRrestOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "rrest"
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

	if val.Kind != ListValue {
		runTimeError2(frame, "First argument not list in %s operator", opName)
	}
	list, convok := val.Data.(*List)
	if !convok {
		runTimeError2(frame, "First argument is not list in %s operator", opName)
	}

	retVal.Kind = ListValue
	if list.Tail == nil && list.Head == nil {
		runTimeError2(frame, "Attempt to access empty list in %s operator", opName)
	}
	if list.Tail != nil {
		retVal.Data = &List{Head: list.Head, Tail: list.Tail.Next}
		return
	}
	if list.Head.Next == nil {
		retVal.Data = &List{}
		return
	}
	// ok, lets copy all from head and drop latest away
	nextitem := list.Head
	var newitem *ListObject
	var prevnewitem *ListObject
	var newHead *ListObject
	for {
		if nextitem == nil {
			break
		}
		if newitem != nil {
			prevnewitem = newitem
		}
		newitem = &ListObject{Val: nextitem.Val, Next: nil}
		if newHead == nil {
			newHead = newitem
		}
		if prevnewitem != nil {
			prevnewitem.Next = newitem
		}
		nextitem = nextitem.Next
	}
	prevnewitem.Next = nil
	retVal.Data = &List{Head: newHead, Tail: nil}
	return
}

//MakeListOfValues offers API for std to create list for values
func MakeListOfValues(frame *Frame, values []Value) (retVal Value) {
	var vitems []*Item
	for _, val := range values {
		vitem := &Item{Type: ValueItem, Data: val}
		vitems = append(vitems, vitem)
	}
	return handleListOP(frame, vitems)
}

func handleListOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "list"
	retVal.Kind = ListValue
	if len(operands) == 0 {
		retVal.Data = &List{Head: nil, Tail: nil}
		return
	}
	var lobj *ListObject
	var head *ListObject
	for _, v := range operands {
		var argval Value
		switch v.Type {
		case ValueItem:
			argval = v.Data.(Value)
		case SymbolPathItem, OperCallItem:
			argval = EvalItem(v, frame)
		default:
			runTimeError2(frame, "something wrong (%s)", opName)
		}
		if lobj == nil {
			lobj = &ListObject{Val: &argval, Next: nil}
			head = lobj
		} else {
			lob := ListObject{Val: &argval, Next: nil}
			lobj.Next = &lob
			lobj = &lob
		}
	}
	retVal.Data = &List{Head: head, Tail: nil}
	return
}

func handleAppendOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "append"
	retVal.Kind = ListValue
	if l := len(operands); l < 1 {
		runTimeError2(frame, "%s operator needs at least one argument (%d given)", opName, l)
	}

	var val Value
	switch v := operands[0]; v.Type {
	case ValueItem:
		val = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		val = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	if val.Kind != ListValue {
		runTimeError2(frame, "First argument not list in %s operator", opName)
	}
	list, convok := val.Data.(*List)
	if !convok {
		runTimeError2(frame, "First argument is not list in %s operator", opName)
	}

	// special case that there is nothing to append
	if len(operands) == 1 {
		retVal.Data = list
		return
	}

	var lobj *ListObject
	var head *ListObject
	realOperands := operands[1:]
	// lets loop operands in reversed order
	for i := len(realOperands) - 1; i >= 0; i-- {
		v := realOperands[i]
		var argval Value
		switch v.Type {
		case ValueItem:
			argval = v.Data.(Value)
		case SymbolPathItem, OperCallItem:
			argval = EvalItem(v, frame)
		default:
			runTimeError2(frame, "something wrong (%s)", opName)
		}
		if lobj == nil {
			lobj = &ListObject{Val: &argval, Next: nil}
			head = lobj
		} else {
			lob := ListObject{Val: &argval, Next: nil}
			lobj.Next = &lob
			lobj = &lob
		}
	}
	lobj.Next = list.Tail
	retVal.Data = &List{Head: list.Head, Tail: head}
	return
}

func handleAddOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "add"
	retVal.Kind = ListValue
	if l := len(operands); l < 1 {
		runTimeError2(frame, "%s operator needs at least one argument (%d given)", opName, l)
	}

	var val Value
	switch v := operands[0]; v.Type {
	case ValueItem:
		val = v.Data.(Value)
	case SymbolPathItem, OperCallItem:
		val = EvalItem(v, frame)
	default:
		runTimeError2(frame, "something wrong (%s)", opName)
	}

	if val.Kind != ListValue {
		runTimeError2(frame, "First argument not list in %s operator", opName)
	}
	list, convok := val.Data.(*List)
	if !convok {
		runTimeError2(frame, "First argument is not list in %s operator", opName)
	}

	// special case that there is nothing to append
	if len(operands) == 1 {
		retVal.Data = list
		return
	}

	var lobj *ListObject
	var head *ListObject
	for _, v := range operands[1:] {
		var argval Value
		switch v.Type {
		case ValueItem:
			argval = v.Data.(Value)
		case SymbolPathItem, OperCallItem:
			argval = EvalItem(v, frame)
		default:
			runTimeError2(frame, "something wrong (%s)", opName)
		}
		if lobj == nil {
			lobj = &ListObject{Val: &argval, Next: nil}
			head = lobj
		} else {
			lob := ListObject{Val: &argval, Next: nil}
			lobj.Next = &lob
			lobj = &lob
		}
	}
	lobj.Next = list.Head
	retVal.Data = &List{Head: head, Tail: list.Tail}
	return
}

func reverseCopy(src *List) *List {
	cur := src.Tail
	var prev *ListObject
	newList := List{Head: nil, Tail: nil}
	for {
		newItem := ListObject{Val: cur.Val, Next: prev}
		prev = &newItem
		cur = cur.Next
		if cur == nil {
			newList.Head = prev
			return &newList
		}
	}
}

func handleRestOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "rest"
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

	if val.Kind != ListValue {
		runTimeError2(frame, "First argument not list in %s operator", opName)
	}
	list, convok := val.Data.(*List)
	if !convok {
		runTimeError2(frame, "First argument is not list in %s operator", opName)
	}

	retVal.Kind = ListValue
	if list.Head == nil {
		// we need to revert copy tail to head
		if list.Tail != nil {
			rev := reverseCopy(list)
			if rev.Head == nil {
				runTimeError2(frame, "Fatal error in %s operator", opName)
			}
			retVal.Data = &List{Head: rev.Head.Next, Tail: rev.Tail}
			return
		}
		runTimeError2(frame, "Attempt to access empty list in %s operator", opName)
		return
	}
	retVal.Data = &List{Head: list.Head.Next, Tail: list.Tail}
	return
}

func handleLastOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "last"
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

	if val.Kind != ListValue {
		runTimeError2(frame, "First argument not list in %s operator", opName)
	}
	list, convok := val.Data.(*List)
	if !convok {
		runTimeError2(frame, "First argument is not list in %s operator", opName)
	}

	if list.Head == nil && list.Tail == nil {
		runTimeError2(frame, "Attempt to access empty list")
	} else if list.Tail != nil {
		retVal = *list.Tail.Val
		return
	} else if list.Head != nil {
		// ok, we need to find tail by hard way
		cur := list.Head
		for {
			if cur.Next == nil {
				retVal = *cur.Val
				return
			}
			cur = cur.Next
		}
	}
	return
}

func handleHeadOP(frame *Frame, operands []*Item) (retVal Value) {
	opName := "head"
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

	if val.Kind != ListValue {
		runTimeError2(frame, "First argument not list in %s operator", opName)
	}
	list, convok := val.Data.(*List)
	if !convok {
		runTimeError2(frame, "First argument is not list in %s operator", opName)
	}

	if list.Tail == nil && list.Head == nil {
		runTimeError2(frame, "Attempt to access empty list")
	} else if list.Head != nil {
		retVal = *list.Head.Val
		return
	} else if list.Tail != nil {
		// ok, we need to find tail by hard way
		cur := list.Tail
		for {
			if cur.Next == nil {
				retVal = *cur.Val
				return
			}
			cur = cur.Next
		}
	}
	return
}
