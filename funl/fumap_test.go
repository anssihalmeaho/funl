package funl

import (
	"fmt"
	"testing"
)

func TestItemsWithCollidingKeys(t *testing.T) {
	mapval := handleMapOP(nil, []*Item{})
	mapval = handlePutOP(nil, []*Item{
		&Item{Type: ValueItem, Data: mapval},
		&Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 0}},
		&Item{Type: ValueItem, Data: Value{Kind: StringValue, Data: "int-0"}},
	})
	mapval = handlePutOP(nil, []*Item{
		&Item{Type: ValueItem, Data: mapval},
		&Item{Type: ValueItem, Data: Value{Kind: FloatValue, Data: 0.0}},
		&Item{Type: ValueItem, Data: Value{Kind: StringValue, Data: "float-0"}},
	})
	//t.Logf("%s", mapval)
}

func TestReuseDeletedItemInPut(t *testing.T) {
	defaultVal := &Item{Type: ValueItem, Data: Value{Kind: StringValue, Data: "not found"}}
	nums := []Item{
		Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 10}},
		Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 20}},
		Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 30}},
	}

	mapval := handleMapOP(nil, []*Item{})
	if mapval.Kind != MapValue {
		t.Errorf("Map assumed, got: %d", mapval.Kind)
	}

	for _, num := range nums {
		val := Item{Type: ValueItem, Data: Value{Kind: StringValue, Data: fmt.Sprintf("val-%d", num.Data.(Value).Data)}}
		mapItem := Item{Type: ValueItem, Data: mapval}
		mapval = handlePutOP(nil, []*Item{&mapItem, &num, &val})
		if mapval.Kind != MapValue {
			t.Errorf("Map assumed, got: %d", mapval.Kind)
		}
	}

	mapItem := Item{Type: ValueItem, Data: mapval}

	// lets delete one item
	newMapVal := handleDelOP(nil, []*Item{&mapItem, &nums[1]})
	if newMapVal.Kind != MapValue {
		t.Errorf("Map value expected, got: %d", newMapVal.Kind)
	}
	newMapItem := Item{Type: ValueItem, Data: newMapVal}

	// lets put new value with same key as deleted item
	newValSameKey := Item{Type: ValueItem, Data: Value{Kind: StringValue, Data: "new-20"}}
	latestMapval := handlePutOP(nil, []*Item{&newMapItem, &nums[1], &newValSameKey})
	if latestMapval.Kind != MapValue {
		t.Errorf("Map assumed, got: %d", latestMapval.Kind)
	}
	latestMapItem := Item{Type: ValueItem, Data: latestMapval}

	// new value should be found with the key
	latestVal := handleGetOP(nil, []*Item{&latestMapItem, &nums[1], defaultVal})
	if latestVal.Kind != StringValue {
		t.Errorf("String value expected, got: %d", latestVal.Kind)
	}
	if s := latestVal.Data.(string); s != "new-20" {
		t.Errorf("Unexpected value: %s", s)
	}
}

func TestMapKeyDeletion(t *testing.T) {
	nums := []Item{
		Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 10}},
		Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 20}},
		Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 30}},
	}

	mapval := handleMapOP(nil, []*Item{})
	if mapval.Kind != MapValue {
		t.Errorf("Map assumed, got: %d", mapval.Kind)
	}

	for _, num := range nums {
		val := Item{Type: ValueItem, Data: Value{Kind: StringValue, Data: fmt.Sprintf("val-%d", num.Data.(Value).Data)}}
		mapItem := Item{Type: ValueItem, Data: mapval}
		mapval = handlePutOP(nil, []*Item{&mapItem, &num, &val})
		if mapval.Kind != MapValue {
			t.Errorf("Map assumed, got: %d", mapval.Kind)
		}
	}

	mapItem := Item{Type: ValueItem, Data: mapval}

	// lets delete one item
	newMapVal := handleDelOP(nil, []*Item{&mapItem, &nums[1]})
	if newMapVal.Kind != MapValue {
		t.Errorf("Map value expected, got: %d", newMapVal.Kind)
	}
	newMapItem := Item{Type: ValueItem, Data: newMapVal}

	defaultVal := &Item{Type: ValueItem, Data: Value{Kind: StringValue, Data: "not found"}}

	// deleted key should not be found
	val := handleGetOP(nil, []*Item{&newMapItem, &nums[1], defaultVal})
	if val.Kind != StringValue {
		t.Errorf("String value expected, got: %d", val.Kind)
	}
	if s := val.Data.(string); s != "not found" {
		t.Errorf("Unexpected value: %s", s)
	}

	// deleted key should be found from previous version of map
	val = handleGetOP(nil, []*Item{&mapItem, &nums[1], defaultVal})
	if val.Kind != StringValue {
		t.Errorf("String value expected, got: %d", val.Kind)
	}
	if s := val.Data.(string); s != "val-20" {
		t.Errorf("Unexpected value: %s", s)
	}

	// len of new map should be 2
	lenv := handleLenOP(nil, []*Item{&newMapItem})
	if lenv.Kind != IntValue {
		t.Errorf("Unexpected value: %s", lenv)
	}
	if mlen := lenv.Data.(int); mlen != 2 {
		t.Errorf("Unexpected length: %d", mlen)
	}

	// len of previous map should be 3
	lenv = handleLenOP(nil, []*Item{&mapItem})
	if lenv.Kind != IntValue {
		t.Errorf("Unexpected value: %s", lenv)
	}
	if mlen := lenv.Data.(int); mlen != 3 {
		t.Errorf("Unexpected length: %d", mlen)
	}

	// keys should not contain deleted key
	keysv := handleKeysOP(nil, []*Item{&newMapItem})
	if keysv.Kind != ListValue {
		t.Errorf("Unexpected value: %v", keysv)
	}
	liter := NewListIterator(keysv)
	keysCount := 0
	hasDeletedOne := false
	for {
		nextv := liter.Next()
		if nextv == nil {
			break
		}
		if nextv.Kind == IntValue && nextv.Data.(int) == 20 {
			hasDeletedOne = true
		}
		keysCount++
	}
	if keysCount != 2 {
		t.Errorf("Unexpected count: %d", keysCount)
	}
	if hasDeletedOne {
		t.Errorf("Deleted item still included")
	}

	// vals should not contain deleted key
	valsv := handleValsOP(nil, []*Item{&newMapItem})
	if valsv.Kind != ListValue {
		t.Errorf("Unexpected value: %v", valsv)
	}
	liter = NewListIterator(valsv)
	valsCount := 0
	hasDeletedOne = false
	for {
		nextv := liter.Next()
		if nextv == nil {
			break
		}
		if nextv.Kind == StringValue && nextv.Data.(string) == "val-20" {
			hasDeletedOne = true
		}
		valsCount++
	}
	if valsCount != 2 {
		t.Errorf("Unexpected count: %d", valsCount)
	}
	if hasDeletedOne {
		t.Errorf("Deleted item still included")
	}

	// keyvals should not contain deleted key
	keyvalsv := handleKeyvalsOP(nil, []*Item{&newMapItem})
	if keyvalsv.Kind != ListValue {
		t.Errorf("Unexpected value: %v", keyvalsv)
	}
	liter = NewListIterator(keyvalsv)
	keyvalsCount := 0
	for {
		nextv := liter.Next()
		if nextv == nil {
			break
		}
		keyvalsCount++
	}
	if keyvalsCount != 2 {
		t.Errorf("Unexpected count: %d", keyvalsCount)
	}

	// in -operator should not found deleted item
	isInVal := handleInOP(nil, []*Item{&newMapItem, &nums[1]})
	if isInVal.Kind != BoolValue {
		t.Errorf("unexpected type: %v", isInVal)
	}
	if isItemInMap := isInVal.Data.(bool); isItemInMap {
		t.Errorf("in operator finds deleted item")
	}
}

func TestMapOKfloatKeys(t *testing.T) {
	nums := []Item{
		Item{Type: ValueItem, Data: Value{Kind: FloatValue, Data: 0.01}},
		Item{Type: ValueItem, Data: Value{Kind: FloatValue, Data: 0.2}},
		Item{Type: ValueItem, Data: Value{Kind: FloatValue, Data: 3.0}},
	}

	mapval := handleMapOP(nil, []*Item{})
	if mapval.Kind != MapValue {
		t.Errorf("Map assumed, got: %d", mapval.Kind)
	}

	for _, num := range nums {
		val := Item{Type: ValueItem, Data: Value{Kind: StringValue, Data: fmt.Sprintf("val-%v", num.Data.(Value).Data)}}
		mapItem := Item{Type: ValueItem, Data: mapval}
		mapval = handlePutOP(nil, []*Item{&mapItem, &num, &val})
		if mapval.Kind != MapValue {
			t.Errorf("Map assumed, got: %d", mapval.Kind)
		}
	}

	mapItem := Item{Type: ValueItem, Data: mapval}
	val := handleGetOP(nil, []*Item{&mapItem, &nums[1]})
	if val.Kind != StringValue {
		t.Errorf("String value expected, got: %d", val.Kind)
	}
	if s := val.Data.(string); s != "val-0.2" {
		t.Errorf("Unexpected value: %s", s)
	}
}

func TestMapOKintKeys(t *testing.T) {
	nums := []Item{
		Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 10}},
		Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 20}},
		Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 30}},
	}

	mapval := handleMapOP(nil, []*Item{})
	if mapval.Kind != MapValue {
		t.Errorf("Map assumed, got: %d", mapval.Kind)
	}

	for _, num := range nums {
		val := Item{Type: ValueItem, Data: Value{Kind: StringValue, Data: fmt.Sprintf("val-%d", num.Data.(Value).Data)}}
		mapItem := Item{Type: ValueItem, Data: mapval}
		mapval = handlePutOP(nil, []*Item{&mapItem, &num, &val})
		if mapval.Kind != MapValue {
			t.Errorf("Map assumed, got: %d", mapval.Kind)
		}
	}

	mapItem := Item{Type: ValueItem, Data: mapval}
	val := handleGetOP(nil, []*Item{&mapItem, &nums[1]})
	if val.Kind != StringValue {
		t.Errorf("String value expected, got: %d", val.Kind)
	}
	if s := val.Data.(string); s != "val-20" {
		t.Errorf("Unexpected value: %s", s)
	}
}

func TestMapOKstringKeys(t *testing.T) {
	nums := []Item{
		Item{Type: ValueItem, Data: Value{Kind: StringValue, Data: "something ABC"}},
		Item{Type: ValueItem, Data: Value{Kind: StringValue, Data: "something DEF"}},
		Item{Type: ValueItem, Data: Value{Kind: StringValue, Data: "something GHI"}},
	}

	mapval := handleMapOP(nil, []*Item{})
	if mapval.Kind != MapValue {
		t.Errorf("Map assumed, got: %d", mapval.Kind)
	}

	for _, num := range nums {
		val := Item{Type: ValueItem, Data: Value{Kind: StringValue, Data: fmt.Sprintf("val-%s", num.Data.(Value).Data)}}
		mapItem := Item{Type: ValueItem, Data: mapval}
		mapval = handlePutOP(nil, []*Item{&mapItem, &num, &val})
		if mapval.Kind != MapValue {
			t.Errorf("Map assumed, got: %d", mapval.Kind)
		}
	}

	mapItem := Item{Type: ValueItem, Data: mapval}
	val := handleGetOP(nil, []*Item{&mapItem, &nums[1]})
	if val.Kind != StringValue {
		t.Errorf("String value expected, got: %d", val.Kind)
	}
	if s := val.Data.(string); s != "val-something DEF" {
		t.Errorf("Unexpected value: %s", s)
	}
}

func TestMapOKListAsKeys(t *testing.T) {
	operands1 := []*Item{
		&Item{Type: ValueItem, Data: Value{Kind: StringValue, Data: "something ABC"}},
		&Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 20}},
		&Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 30}},
	}
	list1 := handleListOP(nil, operands1)
	if list1.Kind != ListValue {
		t.Errorf("List assumed, got: %d", list1.Kind)
	}

	operands2 := []*Item{
		&Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 10}},
		&Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 20}},
		&Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 30}},
	}
	list2 := handleListOP(nil, operands2)
	if list2.Kind != ListValue {
		t.Errorf("List assumed, got: %d", list1.Kind)
	}

	operands3 := []*Item{
		&Item{Type: ValueItem, Data: Value{Kind: StringValue, Data: "something JKL"}},
		&Item{Type: ValueItem, Data: Value{Kind: IntValue, Data: 40}},
	}
	list3 := handleListOP(nil, operands3)
	if list3.Kind != ListValue {
		t.Errorf("List assumed, got: %d", list1.Kind)
	}

	nums := []Item{
		Item{Type: ValueItem, Data: list1},
		Item{Type: ValueItem, Data: list2},
		Item{Type: ValueItem, Data: list3},
	}

	mapval := handleMapOP(nil, []*Item{})
	if mapval.Kind != MapValue {
		t.Errorf("Map assumed, got: %d", mapval.Kind)
	}

	for _, num := range nums {
		val := Item{Type: ValueItem, Data: Value{Kind: StringValue, Data: fmt.Sprintf("val-%s", num.Data.(Value).Data)}}
		mapItem := Item{Type: ValueItem, Data: mapval}
		mapval = handlePutOP(nil, []*Item{&mapItem, &num, &val})
		if mapval.Kind != MapValue {
			t.Errorf("Map assumed, got: %d", mapval.Kind)
		}
	}

	mapItem := Item{Type: ValueItem, Data: mapval}
	val := handleGetOP(nil, []*Item{&mapItem, &nums[1]})
	if val.Kind != StringValue {
		t.Errorf("String value expected, got: %d", val.Kind)
	}
	if s := val.Data.(string); s != "val-list(10, 20, 30)" {
		t.Errorf("Unexpected value: %s", s)
	}
}
