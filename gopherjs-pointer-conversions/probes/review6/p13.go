package main

import (
	"fmt"
	"unsafe"
)

type UP unsafe.Pointer
type MyInt int

func main() {
	x := 1
	up := UP(&x)
	p := (*MyInt)(up)
	q := (*int)(up)
	fmt.Println(p == (*MyInt)(unsafe.Pointer(&x)), q == &x, (*int)(UP(nil)) == nil)
}
