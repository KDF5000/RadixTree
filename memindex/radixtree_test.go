package memindex

import (
	"log"
	"math"
	"testing"
)

type Location struct {
	ips []string
}

type Segment struct {
	start uint32
	end   uint32
}

var segCases []Segment

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

func TestInsertSegment(t *testing.T) {
	segCases = []Segment{Segment{0x0, 0x3f}, Segment{0x100540, 0x9080e3f}}
	//Location
	location := &Location{}
	location.ips = append(location.ips, "172.18.11.95")
	location.ips = append(location.ips, "172.18.11.97")

	for i := 0; i < len(segCases); i++ {
		t.Logf("Test 0x%x - 0x%x\n", segCases[i].start, segCases[i].end)
		tree := NewRadixTree()
		tree.InstertOrUpdateSegment(segCases[i].start, segCases[i].end, location)

		for key := segCases[i].start; key <= segCases[i].end; key++ {
			item, err := tree.Lookup(uint32(key))
			if err != nil {
				t.Fatalf("Lookup error: %d\n", key)
				break
			}
			if _, ok := item.(*Location); !ok {
				t.Fatalf("Vaule error:%x\n, value:%v", key, item)
				break
			}
		}

		//lookup some keys which was not put into the tree
		for key := segCases[i].end + 1; key <= 0xffffffff && key <= segCases[i].end+128; key++ {
			item, _ := tree.Lookup(uint32(key))
			if item != nil {
				t.Fatalf("The key should not be exist, key:0x%x, val:%v\n", key, item)
				break
			}
		}
	}
}
