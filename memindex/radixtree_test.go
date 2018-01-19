package memindex

import (
	"log"
	"math"
	"testing"
)

func TestBasicOp(t *testing.T) {
	//insert 2 ^32 numbers
	tree := NewRadixTree()
	for i := uint32(0); i <= uint32(math.Pow(2, 32)-1)/10; i++ {
		tree.InsertOrUpdate(i, i)
	}
	// get
	for i := uint32(0); i <= uint32(math.Pow(2, 32)-1)/10; i++ {
		item, err := tree.Lookup(i)
		if err != nil {
			log.Fatal("lookup error")
		}
		// t.Logf("Key=%d, value=%d\n", i, item)
		if item != i {
			t.Fatalf("Lookup result should be %d, but get %v\n", i, item)
		}
	}
}
