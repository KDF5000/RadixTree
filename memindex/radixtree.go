package memindex

import (
	"errors"
	"log"
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
	mask    uint32 //0x3f
	segMask uint32 //the mask for segment insert or get
	rnode   *RadixTreeNode
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
	return &RadixTree{mask: 0x3f, segMask: 0x3f, rnode: &RadixTreeNode{parent: nil, offset: 0, slots: make([]interface{}, SLOTS_SIZE)}}
}

//InsertOrUpdate insert an item into a the tree
func (t *RadixTree) InsertOrUpdate(index uint32, item interface{}) error {
	node := t.getTargetNode(index, 0)
	if node == nil {
		return errors.New("cannot get target node")
	}
	node.slots[index&t.mask] = item
	return nil
}

func (t *RadixTree) getTargetNode(index uint32, minShift uint32) *RadixTreeNode {
	shift := t.getMaxShift(index)
	// fmt.Printf("MaxShift:%d\n", shift)
	curNode := t.rnode
	for shift > minShift {
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

//countDigitOne count the tmes of 1 in number
//the number must be the format: 0000 1111， the 1 appears from the end of the number continuously
func (t *RadixTree) countDigitOne(number uint32) uint32 {
	count := uint32(0)
	for number&0x1 == 1 {
		count++
		number = number >> 1
	}
	return count
}

// InstertOrUpdateSegment update a segment of the whole range
// the segment length must be the times of 2 ^ 6  = 64
// @param start start point of the segment  ****000000
// @param end end point of the segment   ****111111
// @return error
func (t *RadixTree) InstertOrUpdateSegment(start uint32, end uint32, item interface{}) error {
	if (end < start) || (start&t.segMask != 0) || (end&t.segMask != t.mask) {
		return errors.New("the start or end is not the times of segMast")
	}
	// shiftBits := t.countDigitOne(t.segMask)
	// newStart := start >> shiftBits
	// newEnd := end >> shiftBits

	// rangeNum := end - start
	// shiftBits := uint32(0)
	// for rangeNum>>RADIX_TREE_SHIFT > 0 {
	// 	shiftBits += RADIX_TREE_SHIFT
	// 	rangeNum = rangeNum >> RADIX_TREE_SHIFT
	// }
	// newStart := start >> shiftBits
	// newEnd := end >> shiftBits
	// for key := newStart; key <= newEnd; key++ {
	// 	if err := t.InsertOrUpdate(key, item); err != nil {
	// 		return err
	// 	}
	// }
	maxShift := t.getMaxShift(start)
	t.recurseSet(start, end, item, maxShift)
	return nil
}

func (t *RadixTree) setRangeByShift(start uint32, end uint32, item interface{}, shift uint32) error {
	if end < start {
		return nil
	}
	node := t.getTargetNode(start, shift)
	if node == nil {
		return errors.New("cannot get target node")
	}
	for key := (start >> shift) & t.mask; key <= (end>>shift)&t.mask; key++ {
		node.slots[key] = item
	}
	return nil
}

func (t *RadixTree) recurseSet(start uint32, end uint32, item interface{}, shift uint32) {
	if shift < RADIX_TREE_SHIFT {
		return
	}

	var startPiece, endPiece uint32

	for {
		startPiece = (start >> shift) & t.mask
		endPiece = (end >> shift) & t.mask
		if endPiece < startPiece {
			// will not happpen
			log.Fatal("endPiect < startPiece")
			return
		}
		if startPiece != endPiece || shift <= RADIX_TREE_SHIFT {
			break
		}
		shift -= RADIX_TREE_SHIFT
	}

	if shift == RADIX_TREE_SHIFT {
		//the end
		t.setRangeByShift(start, end, item, shift)
	} else {
		highValue := start >> shift << shift
		//first segament continue to recurseSet
		t.recurseSet(start, highValue|(startPiece+1)<<shift-1, item, shift)
		//mid segament can be set directlly
		t.setRangeByShift(highValue|(startPiece+1)<<shift, highValue|(endPiece<<shift)-1, item, shift)
		//the last segament continue to recurseSet
		t.recurseSet(highValue|(endPiece<<shift), end, item, shift)
	}
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
		oldNode := curNode
		curNode, status = curNode.slots[slotOffset].(*RadixTreeNode)
		if status == false {
			return oldNode.slots[slotOffset], nil
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

func (t *RadixTree) printTree() {
	//TODO
}
