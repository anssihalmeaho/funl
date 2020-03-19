package pmap

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

const (
	SOME = 1000 * 1000
)

func verifyThatHasThese(t *testing.T, shouldHave, shouldNotHave map[MKey]MValue, rbmap *RBMap) {
	if rbmap.IsEmpty() {
		t.Fatalf("should not be empty")
	}
	for mkey, mval := range shouldHave {
		val, found := rbmap.Get(mkey)
		if !found {
			t.Fatalf("Not found (%d)(%v)", mkey, mval)
		}
		if val != mval {
			t.Fatalf("values dont match (got: %v)(expect: %v)", val, mval)
		}
	}

	for mkey, mval := range shouldNotHave {
		val, found := rbmap.Get(mkey)
		if found {
			t.Fatalf("Key found, should have NOT (%d)(%v)(%v)", mkey, val, mval)
		}
	}
}

func checkDoubleREDS(node *Node) (string, bool) {
	if node == nil {
		return "", true
	}
	if node.Color == RED {
		if node.Right != nil && node.Right.Color == RED {
			return fmt.Sprintf("Right: %d", node.Key), false
		}
		if node.Left != nil && node.Left.Color == RED {
			return fmt.Sprintf("Left: %d", node.Key), false
		}
	}
	if s, ok := checkDoubleREDS(node.Right); !ok {
		return s, false
	}
	if s, ok := checkDoubleREDS(node.Left); !ok {
		return s, false
	}
	return "", true
}

func checkBlackAmounts(node *Node, bcount int, counts *[]int) (string, bool) {
	if node == nil {
		(*counts) = append((*counts), bcount)
		return "", true
	}
	if node.Color == BLACK {
		bcount++
	}
	if s, ok := checkBlackAmounts(node.Right, bcount, counts); !ok {
		return s, false
	}
	if s, ok := checkBlackAmounts(node.Left, bcount, counts); !ok {
		return s, false
	}
	return "", true
}

func checkBinaryTreeCond(node *Node) (string, bool) {
	if node == nil {
		return "", true
	}
	if node.Right != nil && node.Right.Key <= node.Key {
		return fmt.Sprintf("Wrong key order: %d (right: %d)", node.Key, node.Right.Key), false
	}
	if node.Left != nil && node.Left.Key >= node.Key {
		return fmt.Sprintf("Wrong key order: %d (left: %d)", node.Key, node.Left.Key), false
	}
	if s, ok := checkBinaryTreeCond(node.Right); !ok {
		return s, false
	}
	if s, ok := checkBinaryTreeCond(node.Left); !ok {
		return s, false
	}
	return "", true
}

func gatherKeyCounts(node *Node, keycounts map[MKey]int) {
	if node == nil {
		return
	}
	_, found := keycounts[node.Key]
	if !found {
		keycounts[node.Key] = 1
	} else {
		keycounts[node.Key]++
	}
	gatherKeyCounts(node.Right, keycounts)
	gatherKeyCounts(node.Left, keycounts)
}

func verifyRBInvariants(rbmap *RBMap) (string, bool) {
	//No red node has a red parent
	s, ok := checkDoubleREDS(rbmap.Root)
	if !ok {
		return fmt.Sprintf("Double RED: %s", s), false
	}

	//Every path from the root to an empty node has the same number of black nodes
	var counts []int
	s, ok = checkBlackAmounts(rbmap.Root, 0, &counts)
	if !ok {
		return fmt.Sprintf("Black count: %s", s), false
	}
	for _, v := range counts {
		if v != counts[0] {
			return fmt.Sprintf("Difference in Black count: %v", counts), false
		}
	}

	//Check that binary tree condition holds true
	s, ok = checkBinaryTreeCond(rbmap.Root)
	if !ok {
		return fmt.Sprintf("Binary tree condition failed: %s", s), false
	}

	//Check that there are no duplicate keys
	keycounts := make(map[MKey]int)
	gatherKeyCounts(rbmap.Root, keycounts)
	for k, v := range keycounts {
		if v != 1 {
			return fmt.Sprintf("Key multiple times: %d, %d", v, k), false
		}
	}

	return "", true
}

func checkThatAllFoundinGOMap(m map[int]string, node *Node) (string, bool) {
	if node == nil {
		return "", true
	}
	val, found := m[int(node.Key)]
	if !found {
		return fmt.Sprintf("Item not found in Go-map: %d", node.Key), false
	}
	if val != node.Val {
		return fmt.Sprintf("Item value different in Go-map: %v, %v", val, node.Val), false
	}
	if s, ok := checkThatAllFoundinGOMap(m, node.Right); !ok {
		return s, false
	}
	if s, ok := checkThatAllFoundinGOMap(m, node.Left); !ok {
		return s, false
	}
	return "", true
}

func compareMaps(m map[int]string, rbmap *RBMap) (string, bool) {
	// check that all go-map items are in rbmap
	for k, v := range m {
		val, found := rbmap.Get(MKey(k))
		if !found {
			return fmt.Sprintf("Not found (%d)(%v)", k, v), false
		}
		if val.(string) != v {
			return fmt.Sprintf("values dont match (got: %v)(expect: %v)", val, v), false
		}
	}

	// check that all rbmap items are in go-map
	return checkThatAllFoundinGOMap(m, rbmap.Root)
}

func convGomapToRB(m map[int]string) *RBMap {
	rbmap := NewRBMap()
	var nextRBmap *RBMap
	for k, v := range m {
		nextRBmap = rbmap.Put(MKey(k), v)
		rbmap = nextRBmap
	}
	return nextRBmap
}

func TestSeveralVersions(t *testing.T) {
	type mVersions struct {
		Gomap map[int]string
		Rbmap *RBMap
	}
	mversions := []mVersions{}
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	allValues := make(map[int]string)

	prevRbmap := NewRBMap()
	for i := 0; i < 10; i++ {
		newv := mVersions{
			Gomap: make(map[int]string),
			Rbmap: prevRbmap,
		}
		for j := 0; j < 20; j++ {
			rkey := r.Intn(1000)
			rvalue := fmt.Sprintf("val-%d", rkey)
			_, foundAlready := allValues[rkey]
			if !foundAlready {
				newv.Gomap[rkey] = rvalue
				newv.Rbmap = newv.Rbmap.Put(MKey(rkey), rvalue)
				allValues[rkey] = rvalue
			}
		}
		prevRbmap = newv.Rbmap
		mversions = append(mversions, newv)
	}

	//Check that RBMaps contain all cumulative key-values
	//Also check that item counts match
	checkThatRBMapContainsAllNeeded := func(rbm *RBMap, ind int) (string, bool) {
		itemcount := 0
		for i := 0; i <= ind; i++ {
			itemcount = itemcount + len(mversions[i].Gomap)
			for k, v := range mversions[i].Gomap {
				val, found := rbm.Get(MKey(k))
				if !found {
					return fmt.Sprintf("Not found (%d)(%v)", k, v), false
				}
				if v != val {
					return fmt.Sprintf("values dont match (got: %v)(expect: %v)", v, val), false
				}
			}
		}
		if cnt := rbm.Count(); cnt != itemcount {
			return fmt.Sprintf("Counts dont match, got: %d, expected: %d", cnt, itemcount), false
		}
		return "", true
	}

	for ind, v := range mversions {
		s, ok := verifyRBInvariants(v.Rbmap)
		if !ok {
			t.Errorf("invariants failed: %d: %s", ind, s)
			break
		}
		s, ok = checkThatRBMapContainsAllNeeded(v.Rbmap, ind)
		if !ok {
			t.Errorf("does not contain all needed: %d: %s", ind, s)
			break
		}
	}

	if false {
		for ind, v := range mversions {
			fmt.Println(fmt.Sprintf("=== %d =====================", ind))
			fmt.Println(v.Gomap)
			fmt.Println("---")
			v.Rbmap.Print()
		}
	}
}

func TestBigger(t *testing.T) {
	sourceMap := map[int]string{}
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	for i := 0; i < 30; i++ {
		rval := r.Intn(1000)
		sourceMap[rval] = fmt.Sprintf("val-%d", rval)
	}

	var rbmap *RBMap
	if false {
		fmt.Println("--------------------------")
		fmt.Println(sourceMap)
		rbmap = convGomapToRB(sourceMap)
		fmt.Println("--------------------------")
		rbmap.Print()
		fmt.Println("--------------------------")
	} else {
		rbmap = convGomapToRB(sourceMap)
	}

	if s, ok := verifyRBInvariants(rbmap); !ok {
		t.Errorf("invariant failed (%s)", s)
	}

	if s, ok := compareMaps(sourceMap, rbmap); !ok {
		t.Errorf("invariant failed (%s)", s)
	}
}

func TestEquals(t *testing.T) {
	rbmap1 := NewRBMap()
	rbmap2 := NewRBMap()
	rbmap3 := NewRBMap()
	rbmap4 := NewRBMap()
	rbmap5 := NewRBMap()

	if !(rbmap1.Equals(rbmap2) && rbmap2.Equals(rbmap1)) {
		t.Fatalf("Should be equal")
	}
	if !rbmap1.Equals(rbmap1) {
		t.Fatalf("Should be equal with itself")
	}

	rbmap1 = rbmap1.Put(MKey(30), "C")
	rbmap1 = rbmap1.Put(MKey(10), "A")
	rbmap1 = rbmap1.Put(MKey(40), "D")
	rbmap1 = rbmap1.Put(MKey(20), "B")

	rbmap2 = rbmap2.Put(MKey(20), "B")
	rbmap2 = rbmap2.Put(MKey(40), "D")
	rbmap2 = rbmap2.Put(MKey(30), "C")
	rbmap2 = rbmap2.Put(MKey(10), "A")

	rbmap3 = rbmap3.Put(MKey(10), "A")

	rbmap4 = rbmap4.Put(MKey(20), "B")
	rbmap4 = rbmap4.Put(MKey(40), "XXX")
	rbmap4 = rbmap4.Put(MKey(30), "C")
	rbmap4 = rbmap4.Put(MKey(10), "A")

	rbmap5 = rbmap5.Put(MKey(30), "C")
	rbmap5 = rbmap5.Put(MKey(1000), "A")
	rbmap5 = rbmap5.Put(MKey(40), "D")
	rbmap5 = rbmap5.Put(MKey(20), "B")

	if false {
		fmt.Println("--------------------------")
		rbmap1.Print()
		fmt.Println("--------------------------")
		rbmap2.Print()
		fmt.Println("--------------------------")
	}

	if !(rbmap1.Equals(rbmap2) && rbmap2.Equals(rbmap1)) {
		t.Fatalf("Should be equal")
	}
	if !rbmap1.Equals(rbmap1) {
		t.Fatalf("Should be equal with itself")
	}
	if rbmap1.Equals(rbmap3) {
		t.Fatalf("Should not be equal")
	}
	if rbmap1.Equals(rbmap4) {
		t.Fatalf("Should not be equal")
	}
	if rbmap1.Equals(rbmap5) {
		t.Fatalf("Should not be equal")
	}

	rbmap6 := rbmap1

	if !rbmap1.Equals(rbmap6) {
		t.Fatalf("Should be equal")
	}
}

func TestValues(t *testing.T) {
	var rbmap *RBMap
	rbmap = NewRBMap()

	emptyvalues := rbmap.Values()
	if l := len(emptyvalues); l != 0 {
		t.Fatalf("Should be empty (%d)", l)
	}

	rbmap = rbmap.Put(MKey(30), "C")
	rbmap = rbmap.Put(MKey(10), "A")
	rbmap = rbmap.Put(MKey(40), "D")
	rbmap = rbmap.Put(MKey(20), "B")

	allvalues := rbmap.Values()

	checkThatHasValue := func(s string, vl []MValue) {
		for _, v := range vl {
			if s == v {
				return
			}
		}
		t.Fatalf("%s not found", s)
	}

	checkThatHasValue("A", allvalues)
	checkThatHasValue("B", allvalues)
	checkThatHasValue("C", allvalues)
	checkThatHasValue("D", allvalues)
	if l := len(allvalues); l != 4 {
		t.Fatalf("Wrong length: %d", l)
	}

	rbmap2 := rbmap
	rbmap2 = rbmap2.Put(MKey(103), "G")
	rbmap2 = rbmap2.Put(MKey(100), "E")
	rbmap2 = rbmap2.Put(MKey(104), "H")
	rbmap2 = rbmap2.Put(MKey(101), "F")

	allvalues = rbmap.Values()
	allvalues2 := rbmap2.Values()

	checkThatHasValue("A", allvalues)
	checkThatHasValue("B", allvalues)
	checkThatHasValue("C", allvalues)
	checkThatHasValue("D", allvalues)
	if l := len(allvalues); l != 4 {
		t.Fatalf("Wrong length: %d", l)
	}

	checkThatHasValue("E", allvalues2)
	checkThatHasValue("F", allvalues2)
	checkThatHasValue("G", allvalues2)
	checkThatHasValue("H", allvalues2)
	if l := len(allvalues2); l != 8 {
		t.Fatalf("Wrong length: %d", l)
	}
}

func TestKeys(t *testing.T) {
	check := func(keys, allkeys []MKey, asmap map[MKey]bool) {
		if l1, l2 := len(keys), len(allkeys); l1 != l2 {
			t.Fatalf("Key count doesnt match, got: %d, expect: %d", l2, l1)
		}
		resmap := make(map[MKey]bool)
		for _, v := range allkeys {
			if !asmap[v] {
				t.Fatalf("Key not found: %d", v)
			}
			resmap[v] = true
		}
		if l1, l2 := len(resmap), len(keys); l1 != l2 {
			t.Fatalf("Not all same values in keys (%d)(%d)", l1, l2)
		}
	}

	var rbmap *RBMap
	rbmap = NewRBMap()

	emptykeys := rbmap.Keys()
	if l := len(emptykeys); l != 0 {
		t.Fatalf("Should be empty (%d)", l)
	}

	keys := []MKey{MKey(40), MKey(30), MKey(20), MKey(10)}
	asmap := make(map[MKey]bool)
	for _, onekey := range keys {
		rbmap = rbmap.Put(onekey, fmt.Sprintf("val-%d", onekey))
		asmap[onekey] = true
	}
	allkeys := rbmap.Keys()

	check(keys, allkeys, asmap)

	// lets make 2nd version
	rbmap2 := rbmap
	keys2 := []MKey{}
	asmap2 := make(map[MKey]bool)
	// lets copy old values to those
	for k, v := range asmap {
		keys2 = append(keys2, k)
		asmap2[k] = v
	}
	morekeys := []MKey{MKey(400), MKey(3), MKey(200), MKey(1)}
	keys2 = append(keys2, morekeys...)
	for _, onekey := range morekeys {
		rbmap2 = rbmap2.Put(onekey, fmt.Sprintf("val-%d", onekey))
		asmap2[onekey] = true
	}

	allkeys = rbmap.Keys()
	allkeys2 := rbmap2.Keys()

	check(keys, allkeys, asmap)
	check(keys2, allkeys2, asmap2)
}

func TestSome(t *testing.T) {
	var rbmap *RBMap
	rbmap = NewRBMap()

	rbmap1 := rbmap.Put(MKey(1), "value 1")
	rbmap2 := rbmap1.Put(MKey(10), "value 2")
	rbmap3 := rbmap2.Put(MKey(15), 10)
	rbmap4 := rbmap3.Put(MKey(3), 20)
	rbmap5 := rbmap4.Put(MKey(12), "value 3")

	if false {
		rbmap5.Print()
		fmt.Println("--------------------------")
		rbmap4.Print()
		fmt.Println("--------------------------")
	}

	if !rbmap.IsEmpty() {
		t.Fatalf("should be empty")
	}
	if cnt1, cnt2 := rbmap.Count(), 0; cnt1 != cnt2 {
		t.Errorf("Wrong counts, got: %d, expected: %d", cnt1, cnt2)
	}

	shouldHave := map[MKey]MValue{
		MKey(1): "value 1",
	}
	shouldNotHave := map[MKey]MValue{
		MKey(10): "value 2",
		MKey(15): 10,
		MKey(3):  20,
		MKey(12): "value 3",
	}
	verifyThatHasThese(t, shouldHave, shouldNotHave, rbmap1)
	if cnt1, cnt2 := rbmap1.Count(), len(shouldHave); cnt1 != cnt2 {
		t.Errorf("Wrong counts, got: %d, expected: %d", cnt1, cnt2)
	}

	shouldHave = map[MKey]MValue{
		MKey(1):  "value 1",
		MKey(10): "value 2",
	}
	shouldNotHave = map[MKey]MValue{
		MKey(15): 10,
		MKey(3):  20,
		MKey(12): "value 3",
	}
	verifyThatHasThese(t, shouldHave, shouldNotHave, rbmap2)
	if cnt1, cnt2 := rbmap2.Count(), len(shouldHave); cnt1 != cnt2 {
		t.Errorf("Wrong counts, got: %d, expected: %d", cnt1, cnt2)
	}

	shouldHave = map[MKey]MValue{
		MKey(1):  "value 1",
		MKey(10): "value 2",
		MKey(15): 10,
	}
	shouldNotHave = map[MKey]MValue{
		MKey(3):  20,
		MKey(12): "value 3",
	}
	verifyThatHasThese(t, shouldHave, shouldNotHave, rbmap3)
	if cnt1, cnt2 := rbmap3.Count(), len(shouldHave); cnt1 != cnt2 {
		t.Errorf("Wrong counts, got: %d, expected: %d", cnt1, cnt2)
	}

	shouldHave = map[MKey]MValue{
		MKey(1):  "value 1",
		MKey(10): "value 2",
		MKey(15): 10,
		MKey(3):  20,
	}
	shouldNotHave = map[MKey]MValue{
		MKey(12): "value 3",
	}
	verifyThatHasThese(t, shouldHave, shouldNotHave, rbmap4)
	if cnt1, cnt2 := rbmap4.Count(), len(shouldHave); cnt1 != cnt2 {
		t.Errorf("Wrong counts, got: %d, expected: %d", cnt1, cnt2)
	}

	shouldHave = map[MKey]MValue{
		MKey(1):  "value 1",
		MKey(10): "value 2",
		MKey(15): 10,
		MKey(3):  20,
		MKey(12): "value 3",
	}
	shouldNotHave = map[MKey]MValue{}
	verifyThatHasThese(t, shouldHave, shouldNotHave, rbmap5)
	if cnt1, cnt2 := rbmap5.Count(), len(shouldHave); cnt1 != cnt2 {
		t.Errorf("Wrong counts, got: %d, expected: %d", cnt1, cnt2)
	}
}

func findFromGomap(m map[int]int, key int) int {
	val, found := m[key]
	if !found {
		return 0
	}
	return val
}

func BenchmarkKeyFindGomap(b *testing.B) {
	m := getBenchMap()
	for n := 0; n < b.N; n++ {
		findFromGomap(m, 30)
		findFromGomap(m, 35)
		findFromGomap(m, 40)
		findFromGomap(m, 200)
		findFromGomap(m, 300)
	}
}

func findFromGomapWithMutex(rwmutex *sync.RWMutex, m map[int]int, key int) int {
	rwmutex.RLock()
	defer rwmutex.RUnlock()

	val, found := m[key]
	if !found {
		return 0
	}
	return val
}

func BenchmarkKeyFindGomapWithMutex(b *testing.B) {
	var mutexi sync.RWMutex
	m := getBenchMap()
	for n := 0; n < b.N; n++ {
		findFromGomapWithMutex(&mutexi, m, 30)
		findFromGomapWithMutex(&mutexi, m, 35)
		findFromGomapWithMutex(&mutexi, m, 40)
		findFromGomapWithMutex(&mutexi, m, 200)
		findFromGomapWithMutex(&mutexi, m, 300)
	}
}

func findFromRBmap(rbmap *RBMap, key int) int {
	val, found := rbmap.Get(MKey(key))
	if !found {
		return 0
	}
	return val.(int)
}

func getBenchMap() map[int]int {
	dum := 1
	m := make(map[int]int)
	for i := 0; i < 10*1000; i++ {
		m[i] = dum
	}
	/*
		s := rand.NewSource(time.Now().UnixNano())
		r := rand.New(s)
		m := make(map[int]int)
		for i := 0; i < 10 * 1000; i++ {
			rkey := r.Intn(1000 * 1000)
			m[rkey] = dum
		}
	*/
	return m
}

func BenchmarkKeyFindRBmap(b *testing.B) {
	rbmap := NewRBMap()
	for k, v := range getBenchMap() {
		rbmap = rbmap.Put(MKey(k), v)
	}

	for n := 0; n < b.N; n++ {
		findFromRBmap(rbmap, 30)
		findFromRBmap(rbmap, 35)
		findFromRBmap(rbmap, 40)
		findFromRBmap(rbmap, 200)
		findFromRBmap(rbmap, 300)
	}
}

func findFromSyncmap(sm *sync.Map, key int) int {
	val, found := sm.Load(key)
	if !found {
		return 0
	}
	return val.(int)
}

func BenchmarkKeyFindGoSyncmap(b *testing.B) {
	var sm sync.Map
	for k, v := range getBenchMap() {
		sm.Store(k, v)
	}

	for n := 0; n < b.N; n++ {
		findFromSyncmap(&sm, 30)
		findFromSyncmap(&sm, 35)
		findFromSyncmap(&sm, 40)
		findFromSyncmap(&sm, 200)
		findFromSyncmap(&sm, 300)
	}
}

func putItemsGomapWithMutex(rwmutex *sync.RWMutex, m map[int]int, key int, val int) {
	rwmutex.Lock()
	defer rwmutex.Unlock()

	m[key] = val
}

func BenchmarkPutItemGomapWithMutex(b *testing.B) {
	var mutexi sync.RWMutex
	m := getBenchMap()
	num := SOME
	for n := 0; n < b.N; n++ {
		putItemsGomapWithMutex(&mutexi, m, num, num+10)
		num++
	}
}

func putItemsGomap(m map[int]int, key int, val int) {
	m[key] = val
}

func BenchmarkPutItemGomap(b *testing.B) {
	m := getBenchMap()
	num := SOME
	for n := 0; n < b.N; n++ {
		putItemsGomap(m, num, num+10)
		num++
	}
}

func BenchmarkPutItemGomapDeepCopy(b *testing.B) {
	m := getBenchMap()
	num := SOME
	nextm := m

	makeDeepcopy := func(prevm map[int]int) map[int]int {
		newm := make(map[int]int)
		for k, v := range prevm {
			newm[k] = v
		}
		return newm
	}

	for n := 0; n < b.N; n++ {
		putItemsGomap(nextm, num, num+10)
		nextm = makeDeepcopy(nextm)
		num++
	}
}

func putItemsRBmap(rbmap *RBMap, key int, val int) *RBMap {
	return rbmap.Put(MKey(key), val)
}

func BenchmarkPutItemRBmap(b *testing.B) {
	rbmap := NewRBMap()
	for k, v := range getBenchMap() {
		rbmap = rbmap.Put(MKey(k), v)
	}

	num := SOME
	for n := 0; n < b.N; n++ {
		rbmap = putItemsRBmap(rbmap, num, num+10)
		num++
	}
}

func putItemsSyncmap(sm *sync.Map, key int, val int) {
	sm.Store(key, val)
}

func BenchmarkPutItemSyncmap(b *testing.B) {
	var sm sync.Map
	for k, v := range getBenchMap() {
		sm.Store(k, v)
	}

	num := SOME
	for n := 0; n < b.N; n++ {
		putItemsSyncmap(&sm, num, num+10)
		num++
	}
}

func TestParaRBmap(t *testing.T) {
	rbmap := NewRBMap()
	for k, v := range getBenchMap() {
		rbmap = rbmap.Put(MKey(k), v)
	}

	var wg sync.WaitGroup
	adder := func(myrbmap *RBMap) {
		defer wg.Done()

		num := SOME
		for i := 0; i < 500; i++ {
			myrbmap = putItemsRBmap(myrbmap, num, num+10)
			num++
		}
	}

	numOfGorotutines := 10
	for ind := 0; ind < numOfGorotutines; ind++ {
		wg.Add(1)
		go adder(rbmap)
	}
	wg.Wait()
}
