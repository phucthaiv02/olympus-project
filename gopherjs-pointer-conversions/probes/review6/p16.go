package main

import (
	"fmt"
	"io"
	"unsafe"
)

type MyErr error
type MyInt int
type P *int

var Exp int
var unexp int

func main() {
	a := (*MyErr)(unsafe.Pointer(&io.EOF))
	b := (*MyErr)(unsafe.Pointer(&io.EOF))
	fmt.Println(1, a == b, unsafe.Pointer(a) == unsafe.Pointer(&io.EOF), *a == io.EOF)
	f := func() *int { return &Exp }
	fmt.Println(2, P(&Exp) == P(f()), (*MyInt)(unsafe.Pointer(&Exp)) == (*MyInt)(unsafe.Pointer(f())))
	g := func() *int { return &unexp }
	fmt.Println(3, P(&unexp) == P(g()), (*MyInt)(unsafe.Pointer(&unexp)) == (*MyInt)(unsafe.Pointer(g())))
	n := new(int)
	fmt.Println(4, P(n) == P(n), unsafe.Pointer((*MyInt)(unsafe.Pointer(n))) == unsafe.Pointer(n))
	cl := &[]int{1}
	type PS *[]int
	fmt.Println(5, PS(cl) == PS(cl))
}
