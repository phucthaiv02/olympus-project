package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type MyInt int
type MyInt64 int64
type S struct{ A int64 }
type T struct{ A MyInt64 }
type P *int

func to[T any](p unsafe.Pointer) *T     { return (*T)(p) }
func from[T any](p *T) unsafe.Pointer   { return unsafe.Pointer(p) }
func cv[T, U any](p *T) *U              { return (*U)(unsafe.Pointer(p)) }

func main() {
	var a [2]int64
	pa := (*[2]MyInt64)(unsafe.Pointer(&a))
	fmt.Println(1, to[[2]MyInt64](unsafe.Pointer(&a)) == pa, from(&a) == unsafe.Pointer(pa))
	var s S
	fmt.Println(2, to[T](unsafe.Pointer(&s)) == (*T)(unsafe.Pointer(&s)), from(&s) == unsafe.Pointer((*T)(unsafe.Pointer(&s))))
	var b [2]*int
	pb := cv[[2]*int, [2]P](&b)
	fmt.Println(3, pb == (*[2]P)(unsafe.Pointer(&b)), reflect.TypeOf(pb[0]))
	x := 1
	pb[0] = &x
	fmt.Println(4, b[0] == &x, reflect.TypeOf(b[0]))
	var c [2]int
	pc := cv[[2]int, [2]MyInt](&c)
	fmt.Println(5, pc == (*[2]MyInt)(unsafe.Pointer(&c)), reflect.TypeOf(pc[0]), &pc[1] == (*MyInt)(unsafe.Pointer(&c[1])))
	u := uintptr(unsafe.Pointer(&x))
	fmt.Println(6, (*MyInt)(unsafe.Pointer(u)) == (*MyInt)(unsafe.Pointer(&x)))
	var e [2][]int64
	pe := cv[[2][]int64, [2][]MyInt64](&e)
	fmt.Println(7, pe == (*[2][]MyInt64)(unsafe.Pointer(&e)), pe[0] == nil)
	var d [2]*[2]int
	pd := cv[[2]*[2]int, [2]*[2]MyInt](&d)
	fmt.Println(8, pd == (*[2]*[2]MyInt)(unsafe.Pointer(&d)), pd[0] == nil)
	d[1] = &[2]int{1, 2}
	fmt.Println(9, reflect.TypeOf(pd[1][0]), unsafe.Pointer(pd[1]) == unsafe.Pointer(d[1]))
	pd[0] = &[2]MyInt{3, 4}
	fmt.Println(10, reflect.TypeOf(d[0][0]), d[0][0])
}
