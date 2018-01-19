package memindex

import (
	"errors"
	"math"
	"sync"
)

const (
	RADIX_TREE_SHIFT = 6
)

var (
	SLOTS_SIZE = int32(math.Pow(2, RADIX_TREE_SHIFT))
)

// RadixTree root of radix tree
type RadixTree struct {
	mask  uint32 //0x3f
	rnode *RadixTreeNode
}

// RadixTreeNode node in a Radix Tree
type RadixTreeNode struct {
	sync.RWMutex
	parent *RadixTreeNode //parent of the node
	offset uint8          //offset in parent`s slots
	slots  []interface{}  // the element of the slots can be nil, *LocationList, *RadixTreeNode
}

// NewRadixTree create a radix tree
func NewRadixTree() *RadixTree {
	return &RadixTree{mask: 0x3f, rnode: &RadixTreeNode{parent: nil, offset: 0, slots: make([]interface{}, SLOTS_SIZE)}}
}

//InsertOrUpdate insert an item into a the tree
func (t *RadixTree) InsertOrUpdate(index uint32, item interface{}) error {
	node := t.getTargetNode(index)
	if node == nil {
		return errors.New("cannot get target node")
	}
	node.slots[index&t.mask] = item
	return nil
}

func (t *RadixTree) getTargetNode(index uint32) *RadixTreeNode {
	shift := t.getMaxShift(index)
	// fmt.Printf("MaxShift:%d\n", shift)
	curNode := t.rnode
	for shift > 0 {
		slotOffset := (index >> shift) & t.mask
		curNode.Lock()
		// this will not happen
		if curNode.slots == nil {
			curNode.slots = make([]interface{}, SLOTS_SIZE)
		}
		if curNode.slots[slotOffset] == nil {
			curNode.slots[slotOffset] = &RadixTreeNode{parent: curNode, offset: uint8(slotOffset), slots: make([]interface{}, SLOTS_SIZE)}
		}
		curNode.Unlock()
		// fmt.Println(slotOffset, curNode.slots)
		curNode = curNode.slots[slotOffset].(*RadixTreeNode)
		shift -= RADIX_TREE_SHIFT
	}
	return curNode
}

// Lookup lookup for a item on position index
func (t *RadixTree) Lookup(index uint32) (interface{}, error) {
	shift := t.getMaxShift(index)
	// fmt.Printf("MaxShift:%d\n", shift)
	curNode := t.rnode
	for shift > 0 {
		slotOffset := (index >> shift) & t.mask
		// fmt.Println(slotOffset, curNode.slots)
		if curNode.slots == nil {
			return nil, errors.New("Slots are empty")
		}
		var status bool
		curNode, status = curNode.slots[slotOffset].(*RadixTreeNode)
		if status == false {
			return curNode.slots[slotOffset], nil
		}

		shift -= RADIX_TREE_SHIFT
	}
	return curNode.slots[index&t.mask], nil
}

func (t *RadixTree) getMaxShift(index interface{}) uint32 {
	// fmt.Printf("unsafe.Sizeof(index):%d\n", unsafe.Sizeof(index))
	// return uint32(unsafe.Sizeof(index)) / RADIX_TREE_SHIFT * RADIX_TREE_SHIFT
	switch index.(type) {
	case uint32:
		return 32 / RADIX_TREE_SHIFT * RADIX_TREE_SHIFT
	case uint64:
		return 64 / RADIX_TREE_SHIFT * RADIX_TREE_SHIFT
	default:
		return 0
	}
}
