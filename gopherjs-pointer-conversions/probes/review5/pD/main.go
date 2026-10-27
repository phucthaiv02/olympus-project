package main

import (
	"fmt"
	"unsafe"
)

type MyInt int
type MyInt64 int64

func main() {
	m := map[string]int64{"a": 1}
	q := (*map[string]MyInt64)(unsafe.Pointer(&m))
	fmt.Printf("%T\n", (*q)["a"])
	(*q)["b"] = 2
	fmt.Printf("%T\n", m["b"])

	x := 1
	mp := map[string]*int{"a": &x}
	qp := (*map[string]*MyInt)(unsafe.Pointer(&mp))
	fmt.Println((*qp)["a"] == (*MyInt)(&x))

	c := make(chan int64, 1)
	c <- 7
	fmt.Printf("%T\n", <-*(*chan MyInt64)(unsafe.Pointer(&c)))

	f := func() int64 { return 3 }
	fmt.Printf("%T\n", (*(*func() MyInt64)(unsafe.Pointer(&f)))())
}
