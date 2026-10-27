package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type MyInt int
type MyInt64 int64
type MyC complex128
type E0 struct{}
type F0 struct{}
type S struct{ A int64 }
type T struct{ A MyInt64 }
type PA0 *[0]int

func main() {
	var e [0]int
	p := &e
	fmt.Println(1, (*[0]MyInt)(unsafe.Pointer(p)) != nil, PA0(p) != nil, unsafe.Pointer(p) != nil)
	var z E0
	fmt.Println(2, (*F0)(unsafe.Pointer(&z)) == (*F0)(unsafe.Pointer(&z)), (*F0)(unsafe.Pointer(&z)) != nil)
	var pa *[2]int
	q := (**[2]MyInt)(unsafe.Pointer(&pa))
	fmt.Println(3, *q == nil)
	*q = &[2]MyInt{1, 2}
	fmt.Println(4, pa != nil, reflect.TypeOf(pa[0]), reflect.TypeOf((*q)[0]), unsafe.Pointer(pa) == unsafe.Pointer(*q))
	*q = nil
	fmt.Println(5, pa == nil)
	var v int64
	pv := (*MyInt64)(unsafe.Pointer(&v))
	*pv = 5
	fmt.Println(6, reflect.TypeOf(v), reflect.TypeOf(*pv), pv == (*MyInt64)(unsafe.Pointer(&v)), v)
	var c complex128 = 1 + 2i
	pc := (*MyC)(unsafe.Pointer(&c))
	*pc = 3
	fmt.Println(7, reflect.TypeOf(c), reflect.TypeOf(*pc), c)
	var ps *S
	qq := (**T)(unsafe.Pointer(&ps))
	fmt.Println(8, *qq == nil)
	s := S{1}
	ps = &s
	fmt.Println(9, *qq == (*T)(unsafe.Pointer(&s)), reflect.TypeOf((*qq).A))
	*qq = &T{7}
	fmt.Println(10, ps.A, reflect.TypeOf(ps.A), reflect.TypeOf(*ps))
	*qq = nil
	fmt.Println(11, ps == nil)
	var arr [2]*S
	arr2 := (*[2]*T)(unsafe.Pointer(&arr))
	fmt.Println(12, arr2[0] == nil, arr2[1] == nil)
	arr[1] = &s
	fmt.Println(13, arr2[1] == (*T)(unsafe.Pointer(&s)))
	arr2[0] = &T{3}
	fmt.Println(14, reflect.TypeOf(*arr[0]), arr[0].A)
	// pointer to int64 inside interface boxed
	var ii interface{} = (*MyInt64)(unsafe.Pointer(&v))
	fmt.Println(15, ii == interface{}(pv))
	// 2D arrays
	var g [2][2]int64
	g2 := (*[2][2]MyInt64)(unsafe.Pointer(&g))
	g2[1][1] = 4
	fmt.Println(16, g, reflect.TypeOf(g[1][1]), reflect.TypeOf(g2[1][1]), &g2[1] == (*[2]MyInt64)(unsafe.Pointer(&g[1])))
	row := g2[0]
	row[0] = 99
	fmt.Println(17, g[0][0])
	g2[0] = [2]MyInt64{8, 9}
	fmt.Println(18, g[0], reflect.TypeOf(g[0][1]))
}
