package main

import (
	"fmt"
	"memindex"
	//"unsafe"
)

type Location struct {
	ips []string
}

var hasErr bool
func main() {
	hasErr = false
	// tree := memindex.NewRadixTree()
	// tree.InsertOrUpdate(1, 5)
	// fmt.Println(tree.Lookup(1))

	//Location
	location := &Location{}
	location.ips = append(location.ips, "172.18.11.95")
	location.ips = append(location.ips, "172.18.11.97")

	// tree.InsertOrUpdate(10, location)
	// fmt.Println(tree.Lookup(10))

	tree := memindex.NewRadixTree()
	tree.InstertOrUpdateSegment(0x100540, 0x9080e3f, location)

	for key := 0x100540; key <= 0x9080e3f; key++ {
		item, err :=tree.Lookup(uint32(key));
		if err != nil{
			fmt.Printf("Lookup error: %d\n", key)
			hasErr = true
			break
		}
		if _, ok := item.(*Location); !ok{
			fmt.Printf("Vaule error:%x\n, value:%v" ,item)
			hasErr = true
			break
		}
	}

	//lookup some keys which was not put into the tree
	for key := 0;key<0x100540;key++{
		item, _ := tree.Lookup(uint32(key))
		if item != nil{
			fmt.Printf("The key should not be exist, key=%x\n", key)
			hasErr = true
			break
		}
	}
	for key := 0x9080e3f+1;key<0xffffffff;key++{
		item, _ := tree.Lookup(uint32(key))
		if item != nil{
			fmt.Printf("The key should not be exist, key=%x\n", key)
			hasErr = true
			break
		}
	}
	if !hasErr{
		fmt.Println("Passed!")
	}

	//var a uint32
	//var b uint64
	//var c interface{}
	//c = b
	//switch c.(type) {
	//case uint32:
	//	fmt.Printf("32\n")
	//case uint64:
	//	fmt.Printf("64\n")
	//default:
	//	break
	//}
	//fmt.Printf("uint32:%d, uint64:%d, interface{}:%d\n", uint32(unsafe.Sizeof(a)), uint32(unsafe.Sizeof(b)), uint32(unsafe.Sizeof(c)))
}
