package main

import (
	"fmt"
	"memindex"
	//"unsafe"
)

type Location struct{
	ips []string
}
func main() {
	tree := memindex.NewRadixTree()
	tree.InsertOrUpdate(1, 5)
	fmt.Println(tree.Lookup(1))

	//Location
	location := &Location{}
	location.ips = append(location.ips, "172.18.11.95")
	location.ips = append(location.ips, "172.18.11.97")

	tree.InsertOrUpdate(10, location)
	fmt.Println(tree.Lookup(10))
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
