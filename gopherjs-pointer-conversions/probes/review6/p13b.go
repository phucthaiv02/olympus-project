package main

import (
	"fmt"
	"unsafe"
)

type UP unsafe.Pointer
type S struct{ A int64 }
type T struct{ A int64 }

func main() {
	var a [2]int64
	s := S{1}
	up := UP(&s)
	ua := UP(&a)
	fmt.Println((*T)(up) == (*T)(unsafe.Pointer(&s)), (*S)(up) == &s, (*[2]int64)(ua) == &a, (*T)(UP(nil)) == nil)
}
