package pmap

import (
	"fmt"
	"strings"
)

const (
	//RED represents red color of node
	RED = 0
	//BLACK represents black color of node
	BLACK = 1
)

//MKey is int
type MKey int

//MValue is map value
type MValue interface{}

//Node is node of tree
type Node struct {
	Left  *Node
	Right *Node
	Color int
	Key   MKey
	Val   MValue
}

//RBMap is Red-Black tree
type RBMap struct {
	Root      *Node
	itemCount int
	Handler   NodeHandler
}

//Get finds item from map
func (rbm *RBMap) Get(key MKey) (MValue, bool) {
	if rbm.IsEmpty() {
		return nil, false
	}
	return find(rbm.Root, key)
}

//IsEmpty is true if map is empty, otherwise true
func (rbm *RBMap) IsEmpty() bool {
	if rbm.Root == nil {
		return true
	}
	return false
}

func printNext(node *Node, depth int, direction string) {
	indent := strings.Repeat("..", depth)
	color := "RED"
	if node.Color == BLACK {
		color = "BLACK"
	}
	fmt.Println(fmt.Sprintf("%s%s %p: key=%d, val=%v, colour=%s", indent, direction, node, node.Key, node.Val, color))
	if node.Left != nil {
		printNext(node.Left, depth+1, "LEFT ")
	}
	if node.Right != nil {
		printNext(node.Right, depth+1, "RIGHT")
	}
}

//Count return number of items in RBMap
func (rbm *RBMap) Count() int {
	return rbm.itemCount
}

//Print prints map
func (rbm *RBMap) Print() {
	if rbm.IsEmpty() {
		fmt.Println("Empty")
		return
	}
	printNext(rbm.Root, 0, "ROOT ")
}

//Equals return true if maps are equal, false otherwise
func (rbm *RBMap) Equals(other *RBMap) bool {
	if rbm == other {
		return true
	}
	if rbm.IsEmpty() && other.IsEmpty() {
		return true
	}
	if rbm.IsEmpty() || other.IsEmpty() {
		return false
	}
	if rbm.Count() != other.Count() {
		return false
	}
	m1 := make(map[MKey]MValue)
	m2 := make(map[MKey]MValue)
	getKVs(rbm.Root, &m1)
	getKVs(other.Root, &m2)

	isSubsetOf := func(set1, set2 map[MKey]MValue) bool {
		for k, v := range set1 {
			if v2, found := set2[k]; !found {
				return false
			} else if v != v2 { // equality check might depend...
				return false
			}
		}
		return true
	}
	return isSubsetOf(m1, m2) && isSubsetOf(m2, m1)
}

func getKVs(node *Node, m *map[MKey]MValue) {
	(*m)[node.Key] = node.Val
	if node.Right != nil {
		getKVs(node.Right, m)
	}
	if node.Left != nil {
		getKVs(node.Left, m)
	}
}

func getNextKey(node *Node, keys *[]MKey) {
	*keys = append(*keys, node.Key)
	if node.Right != nil {
		getNextKey(node.Right, keys)
	}
	if node.Left != nil {
		getNextKey(node.Left, keys)
	}
}

//Keys returns list of all keys in map
func (rbm *RBMap) Keys() (keys []MKey) {
	if rbm.IsEmpty() {
		return
	}
	getNextKey(rbm.Root, &keys)
	return
}

func getNextValue(node *Node, values *[]MValue) {
	*values = append(*values, node.Val)
	if node.Right != nil {
		getNextValue(node.Right, values)
	}
	if node.Left != nil {
		getNextValue(node.Left, values)
	}
}

//Values returns list of all values in map
func (rbm *RBMap) Values() (values []MValue) {
	if rbm.IsEmpty() {
		return
	}
	getNextValue(rbm.Root, &values)
	return
}

func visitNext(node *Node, visitor func(*Node)) {
	visitor(node)
	if node.Right != nil {
		visitNext(node.Right, visitor)
	}
	if node.Left != nil {
		visitNext(node.Left, visitor)
	}
}

//VisitAll visits all nodes and calls handler
func (rbm *RBMap) VisitAll(visitor func(*Node)) {
	if rbm.IsEmpty() {
		return
	}
	visitNext(rbm.Root, visitor)
}

/*
func visitNextUntil(node *Node, visitor func(*Node) bool) bool {
	ret := visitor(node)
	if !ret {
		return false
	}
	if node.Right != nil {
		if enough := visitNextUntil(node.Right, visitor); enough {
			return false
		}
	}
	if node.Left != nil {
		if enough := visitNextUntil(node.Left, visitor); enough {
			return false
		}
	}
	return true
}

//VisitAllUntil visits all nodes and calls handler until handler returns false
func (rbm *RBMap) VisitAllUntil(visitor func(*Node) bool) {
	if rbm.IsEmpty() {
		return
	}
	visitNextUntil(rbm.Root, visitor)
}
*/

func find(node *Node, key MKey) (MValue, bool) {
	if node == nil {
		return nil, false
	}
	if key > node.Key {
		return find(node.Right, key)
	} else if key < node.Key {
		return find(node.Left, key)
	} else {
		if node.Key == key {
			return node.Val, true
		}
		return nil, false
	}
}

func balance(node *Node) *Node {
	if node.Color == RED {
		return node
	}
	if node.Right != nil && node.Right.Color == RED {
		rchild := node.Right
		if rchild.Right != nil && rchild.Right.Color == RED {
			moved := rchild.Left
			rchild.Right.Color = BLACK

			rchild.Left = node
			node.Right = moved
			return rchild
		}
		if rchild.Left != nil && rchild.Left.Color == RED {
			moved1 := rchild.Left.Right
			moved2 := rchild.Left.Left
			newtop := rchild.Left
			rchild.Color = BLACK

			newtop.Right = rchild
			newtop.Left = node
			node.Right = moved2
			rchild.Left = moved1
			return newtop
		}
	}
	if node.Left != nil && node.Left.Color == RED {
		lchild := node.Left
		if lchild.Left != nil && lchild.Left.Color == RED {
			moved := lchild.Right
			lchild.Left.Color = BLACK

			lchild.Right = node
			node.Left = moved
			return lchild
		}
		if lchild.Right != nil && lchild.Right.Color == RED {
			moved1 := lchild.Right.Right
			moved2 := lchild.Right.Left
			newtop := lchild.Right
			lchild.Color = BLACK

			newtop.Right = node
			newtop.Left = lchild
			node.Left = moved1
			lchild.Right = moved2
			return newtop
		}
	}
	return node
}

func insert(node *Node, key MKey, val MValue, nhandler NodeHandler) *Node {
	if node == nil {
		return &Node{
			Key:   key,
			Val:   val,
			Color: RED,
		}
	}
	if key > node.Key {
		nodecopy := *node
		newnode := insert(nodecopy.Right, key, val, nhandler)
		nodecopy.Right = newnode
		return balance(&nodecopy)
	} else if key < node.Key {
		nodecopy := *node
		newnode := insert(nodecopy.Left, key, val, nhandler)
		nodecopy.Left = newnode
		return balance(&nodecopy)
	} else {
		nodecopy := nhandler.HandleSameKey(node, key, val)
		return nodecopy
	}
}

//Put puts value to map
func (rbm *RBMap) Put(key MKey, val MValue) *RBMap {
	if rbm.IsEmpty() {
		rootNode := &Node{
			Key:   key,
			Val:   val,
			Color: RED,
		}
		return &RBMap{Root: rootNode, itemCount: 1, Handler: rbm.Handler}
	}
	nodeCopy := *rbm.Root
	newnode := balance(insert(&nodeCopy, key, val, rbm.Handler))
	newnode.Color = BLACK
	return &RBMap{Root: newnode, itemCount: rbm.itemCount + 1, Handler: rbm.Handler}
}

func copyPath(node *Node, key MKey, val MValue, nhandler NodeHandler) (*Node, bool) {
	if node == nil {
		return nil, false
	}
	if key > node.Key {
		nodecopy := *node
		newnode, found := copyPath(nodecopy.Right, key, val, nhandler)
		nodecopy.Right = newnode
		return &nodecopy, found
	} else if key < node.Key {
		nodecopy := *node
		newnode, found := copyPath(nodecopy.Left, key, val, nhandler)
		nodecopy.Left = newnode
		return &nodecopy, found
	} else {
		return nhandler.MarkDeletion(node, key, val)
	}
}

//Modify copies path with modified value
func (rbm *RBMap) Modify(key MKey, val MValue) (*RBMap, bool) {
	if rbm.IsEmpty() {
		return rbm, false
	}
	nodeCopy := *rbm.Root
	newnode, found := copyPath(&nodeCopy, key, val, rbm.Handler)
	return &RBMap{Root: newnode, itemCount: rbm.itemCount, Handler: rbm.Handler}, found
}

//NodeHandler is interface
type NodeHandler interface {
	HandleSameKey(*Node, MKey, MValue) *Node
	MarkDeletion(*Node, MKey, MValue) (*Node, bool)
}

//NewRBMapWithHandler returns new map with handler
func NewRBMapWithHandler(nhandler NodeHandler) *RBMap {
	return &RBMap{Handler: nhandler}
}

//NewRBMap returns new map
func NewRBMap() *RBMap {
	return &RBMap{}
}
