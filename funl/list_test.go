package funl

import (
	"testing"
)

func TestListIterator1(t *testing.T) {
	a := Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 10}}
	b := Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 11}}
	c := Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 12}}
	d := Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 13}}
	e := Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 14}}
	operands1 := []*Item{
		&a,
		&b,
		&c,
	}
	list1 := handleListOP(nil, operands1)
	if list1.Kind != ListValue {
		t.Errorf("List assumed, got: %d", list1.Kind)
	}
	appendArgs := []*Item{
		&Item{Type: ValueItem, Data: list1},
		&d,
		&e,
	}
	list2 := handleAppendOP(nil, appendArgs)
	if list2.Kind != ListValue {
		t.Errorf("List assumed, got: %d", list2.Kind)
	}

	checkDataSame := func(x *Item, item *Value) {
		if item == nil {
			t.Errorf("Item is nil")
		}
		if item.Kind != IntValue {
			t.Errorf("Unexpected kind")
		}
		expData := x.Data.(Value).Data.(int)
		if num := item.Data.(int); num != expData {
			t.Errorf("Unexpected data: expect: %d, got: %d", expData, num)
		}
	}

	it := NewListIterator(list2)
	item := it.Next()
	checkDataSame(&a, item)
	item = it.Next()
	checkDataSame(&b, item)
	item = it.Next()
	checkDataSame(&c, item)
	item = it.Next()
	checkDataSame(&d, item)
	item = it.Next()
	checkDataSame(&e, item)
	item = it.Next()
	if item != nil {
		t.Errorf("expecting nil")
	}
}

func TestListIterator2(t *testing.T) {
	a := Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 10}}
	b := Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 11}}
	c := Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 12}}
	operands1 := []*Item{
		&a,
		&b,
		&c,
	}
	list1 := handleListOP(nil, operands1)
	if list1.Kind != ListValue {
		t.Errorf("List assumed, got: %d", list1.Kind)
	}

	checkDataSame := func(x *Item, item *Value) {
		if item == nil {
			t.Errorf("Item is nil")
		}
		if item.Kind != IntValue {
			t.Errorf("Unexpected kind")
		}
		expData := x.Data.(Value).Data.(int)
		if num := item.Data.(int); num != expData {
			t.Errorf("Unexpected data: expect: %d, got: %d", expData, num)
		}
	}

	it := NewListIterator(list1)
	item := it.Next()
	checkDataSame(&a, item)
	item = it.Next()
	checkDataSame(&b, item)
	item = it.Next()
	checkDataSame(&c, item)
	item = it.Next()
	if item != nil {
		t.Errorf("expecting nil")
	}
}

func TestListIterator3(t *testing.T) {
	d := Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 13}}
	e := Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 14}}
	operands1 := []*Item{}
	list1 := handleListOP(nil, operands1)
	if list1.Kind != ListValue {
		t.Errorf("List assumed, got: %d", list1.Kind)
	}
	appendArgs := []*Item{
		&Item{Type: ValueItem, Data: list1},
		&d,
		&e,
	}
	list2 := handleAppendOP(nil, appendArgs)
	if list2.Kind != ListValue {
		t.Errorf("List assumed, got: %d", list2.Kind)
	}

	checkDataSame := func(x *Item, item *Value) {
		if item == nil {
			t.Errorf("Item is nil")
		}
		if item.Kind != IntValue {
			t.Errorf("Unexpected kind")
		}
		expData := x.Data.(Value).Data.(int)
		if num := item.Data.(int); num != expData {
			t.Errorf("Unexpected data: expect: %d, got: %d", expData, num)
		}
	}

	it := NewListIterator(list2)
	item := it.Next()
	checkDataSame(&d, item)
	item = it.Next()
	checkDataSame(&e, item)
	item = it.Next()
	if item != nil {
		t.Errorf("expecting nil")
	}
}

func TestListIterator4(t *testing.T) {
	operands1 := []*Item{}
	list1 := handleListOP(nil, operands1)
	if list1.Kind != ListValue {
		t.Errorf("List assumed, got: %d", list1.Kind)
	}

	it := NewListIterator(list1)
	item := it.Next()
	if item != nil {
		t.Errorf("expecting nil")
	}
}

func TestListEmptyList(t *testing.T) {
	retVal := handleListOP(nil, []*Item{})

	if retVal.Kind != ListValue {
		t.Errorf("List assumed, got: %d", retVal.Kind)
	}

	isEmptyVal := handleEmptyOP(nil, []*Item{&Item{Type: ValueItem, Data: retVal}})
	if isEmptyVal.Kind != BoolValue {
		t.Fatalf("Funny value, bool expected")
	}
	if isEmptyVal.Data.(bool) != true {
		t.Errorf("should be empty")
	}
}

func TestListAddToFront(t *testing.T) {
	a := Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 10}}
	b := Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 11}}
	c := Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 12}}
	d := Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 13}}
	e := Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 14}}
	operands1 := []*Item{
		&a,
		&b,
		&c,
	}
	retVal := handleListOP(nil, operands1)
	if retVal.Kind != ListValue {
		t.Errorf("List assumed, got: %d", retVal.Kind)
	}

	appendArgs := []*Item{
		&Item{Type: ValueItem, Data: retVal},
		&d,
		&e,
	}
	list2 := handleAddOP(nil, appendArgs)
	if list2.Kind != ListValue {
		t.Errorf("List assumed, got: %d", list2.Kind)
	}

	tailVal := handleLastOP(nil, []*Item{&Item{Type: ValueItem, Data: list2}})
	if tailVal != c.Data {
		t.Errorf("Not head that was assumed: expected: %#v, got: %#v", c, tailVal)
	}

	headVal := handleHeadOP(nil, []*Item{&Item{Type: ValueItem, Data: list2}})
	if headVal != d.Data {
		t.Errorf("Not head that was assumed: expected: %#v, got: %#v", d, headVal)
	}

	// lets check original remains same
	tailVal = handleLastOP(nil, []*Item{&Item{Type: ValueItem, Data: retVal}})
	if tailVal != c.Data {
		t.Errorf("Not head that was assumed: expected: %#v, got: %#v", c, tailVal)
	}

	headVal = handleHeadOP(nil, []*Item{&Item{Type: ValueItem, Data: retVal}})
	if headVal != a.Data {
		t.Errorf("Not head that was assumed: expected: %#v, got: %#v", a, headVal)
	}

	// lets check rest of new one
	restList1 := handleRestOP(nil, []*Item{&Item{Type: ValueItem, Data: list2}})
	isEmptyRest := handleEmptyOP(nil, []*Item{&Item{Type: ValueItem, Data: restList1}})
	if isEmptyRest.Data.(bool) != false {
		t.Errorf("should not be empty")
	}

	headOfRest := handleHeadOP(nil, []*Item{&Item{Type: ValueItem, Data: restList1}})
	if headOfRest != e.Data {
		t.Errorf("Not head that was assumed: expected: %#v, got: %#v", e, headOfRest)
	}

	tailOfRest := handleLastOP(nil, []*Item{&Item{Type: ValueItem, Data: restList1}})
	if tailOfRest != c.Data {
		t.Errorf("Not head that was assumed: expected: %#v, got: %#v", c, tailOfRest)
	}

	// lets check still further rest of new one
	restList2 := handleRestOP(nil, []*Item{&Item{Type: ValueItem, Data: restList1}})
	isEmptyRest = handleEmptyOP(nil, []*Item{&Item{Type: ValueItem, Data: restList2}})
	if isEmptyRest.Data.(bool) != false {
		t.Errorf("should not be empty")
	}

	headOfRest = handleHeadOP(nil, []*Item{&Item{Type: ValueItem, Data: restList2}})
	if headOfRest != a.Data {
		t.Errorf("Not head that was assumed: expected: %#v, got: %#v", a, headOfRest)
	}
}

func TestListBeingShared(t *testing.T) {
	a := Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 10}}
	b := Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 11}}
	c := Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 12}}
	d := Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 13}}
	e := Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 14}}
	f := Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 15}}
	g := Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 16}}
	operands1 := []*Item{
		&a,
		&b,
		&c,
	}

	retVal := handleListOP(nil, operands1)
	if retVal.Kind != ListValue {
		t.Errorf("List assumed, got: %d", retVal.Kind)
	}

	appendArgs := []*Item{
		&Item{Type: ValueItem, Data: retVal},
		&d,
	}
	list2 := handleAppendOP(nil, appendArgs)
	if list2.Kind != ListValue {
		t.Errorf("List assumed, got: %d", list2.Kind)
	}

	tailVal := handleLastOP(nil, []*Item{&Item{Type: ValueItem, Data: list2}})
	if tailVal != d.Data {
		t.Errorf("Not head that was assumed: expected: %#v, got: %#v", d, tailVal)
	}

	headVal := handleHeadOP(nil, []*Item{&Item{Type: ValueItem, Data: list2}})
	if headVal != a.Data {
		t.Errorf("Not head that was assumed: expected: %#v, got: %#v", a, headVal)
	}

	tailVal = handleLastOP(nil, []*Item{&Item{Type: ValueItem, Data: retVal}})
	if tailVal != c.Data {
		t.Errorf("Not head that was assumed: expected: %#v, got: %#v", d, tailVal)
	}

	headVal = handleHeadOP(nil, []*Item{&Item{Type: ValueItem, Data: retVal}})
	if headVal != a.Data {
		t.Errorf("Not head that was assumed: expected: %#v, got: %#v", a, headVal)
	}

	restList2 := handleRestOP(nil, []*Item{&Item{Type: ValueItem, Data: list2}})
	isEmptyRest := handleEmptyOP(nil, []*Item{&Item{Type: ValueItem, Data: restList2}})
	if isEmptyRest.Data.(bool) != false {
		t.Errorf("should not be empty")
	}

	headOfRest := handleHeadOP(nil, []*Item{&Item{Type: ValueItem, Data: restList2}})
	if headOfRest != b.Data {
		t.Errorf("Not head that was assumed: expected: %#v, got: %#v", b, headOfRest)
	}

	tailOfRest := handleLastOP(nil, []*Item{&Item{Type: ValueItem, Data: restList2}})
	if tailOfRest != d.Data {
		t.Errorf("Not head that was assumed: expected: %#v, got: %#v", d, tailOfRest)
	}

	restList1 := handleRestOP(nil, []*Item{&Item{Type: ValueItem, Data: retVal}})
	isEmptyRest = handleEmptyOP(nil, []*Item{&Item{Type: ValueItem, Data: restList1}})
	if isEmptyRest.Data.(bool) != false {
		t.Errorf("should not be empty")
	}

	headOfRest = handleHeadOP(nil, []*Item{&Item{Type: ValueItem, Data: restList1}})
	if headOfRest != b.Data {
		t.Errorf("Not head that was assumed: expected: %#v, got: %#v", b, headOfRest)
	}

	tailOfRest = handleLastOP(nil, []*Item{&Item{Type: ValueItem, Data: restList1}})
	if tailOfRest != c.Data {
		t.Errorf("Not head that was assumed: expected: %#v, got: %#v", c, tailOfRest)
	}

	// yet 3rd list from original
	appendArgs = []*Item{
		&Item{Type: ValueItem, Data: retVal},
		&e,
		&f,
	}
	list3 := handleAppendOP(nil, appendArgs)
	if list3.Kind != ListValue {
		t.Errorf("List assumed, got: %d", list3.Kind)
	}
	tailVal = handleLastOP(nil, []*Item{&Item{Type: ValueItem, Data: list3}})
	if tailVal != f.Data {
		t.Errorf("Not head that was assumed: expected: %#v, got: %#v", f, tailVal)
	}

	headVal = handleHeadOP(nil, []*Item{&Item{Type: ValueItem, Data: list3}})
	if headVal != a.Data {
		t.Errorf("Not head that was assumed: expected: %#v, got: %#v", a, headVal)
	}

	// yet 4th list from second
	appendArgs = []*Item{
		&Item{Type: ValueItem, Data: retVal},
		&g,
	}
	list4 := handleAppendOP(nil, appendArgs)
	if list4.Kind != ListValue {
		t.Errorf("List assumed, got: %d", list4.Kind)
	}
	tailVal = handleLastOP(nil, []*Item{&Item{Type: ValueItem, Data: list4}})
	if tailVal != g.Data {
		t.Errorf("Not head that was assumed: expected: %#v, got: %#v", g, tailVal)
	}

	headVal = handleHeadOP(nil, []*Item{&Item{Type: ValueItem, Data: list4}})
	if headVal != a.Data {
		t.Errorf("Not head that was assumed: expected: %#v, got: %#v", a, headVal)
	}
}

func TestListOK(t *testing.T) {
	a := Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 10}}
	b := Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 11}}
	operands := []*Item{
		&a,
		&b,
	}

	retVal := handleListOP(nil, operands)

	if retVal.Kind != ListValue {
		t.Errorf("List assumed, got: %d", retVal.Kind)
	}

	isEmptyVal := handleEmptyOP(nil, []*Item{&Item{Type: ValueItem, Data: retVal}})
	if isEmptyVal.Kind != BoolValue {
		t.Fatalf("Funny value, bool expected")
	}
	if isEmptyVal.Data.(bool) != false {
		t.Errorf("should not be empty")
	}

	headVal := handleHeadOP(nil, []*Item{&Item{Type: ValueItem, Data: retVal}})

	if headVal != a.Data {
		t.Errorf("Not head that was assumed: expected: %#v, got: %#v", a, headVal)
	}

	isEmptyVal = handleEmptyOP(nil, []*Item{&Item{Type: ValueItem, Data: retVal}})
	if isEmptyVal.Data.(bool) != false {
		t.Errorf("should not be empty")
	}

	tailVal := handleLastOP(nil, []*Item{&Item{Type: ValueItem, Data: retVal}})

	if tailVal != b.Data {
		t.Errorf("Not head that was assumed: expected: %#v, got: %#v", b, tailVal)
	}

	restList := handleRestOP(nil, []*Item{&Item{Type: ValueItem, Data: retVal}})

	isEmptyRest := handleEmptyOP(nil, []*Item{&Item{Type: ValueItem, Data: restList}})
	if isEmptyRest.Data.(bool) != false {
		t.Errorf("should not be empty")
	}

	headOfRest := handleHeadOP(nil, []*Item{&Item{Type: ValueItem, Data: restList}})
	if headOfRest != b.Data {
		t.Errorf("Not head that was assumed: expected: %#v, got: %#v", b, headOfRest)
	}

	tailOfRest := handleLastOP(nil, []*Item{&Item{Type: ValueItem, Data: restList}})
	if tailOfRest != b.Data {
		t.Errorf("Not head that was assumed: expected: %#v, got: %#v", b, tailOfRest)
	}
}
