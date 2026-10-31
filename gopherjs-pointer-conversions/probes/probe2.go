package main

import (
	"fmt"
	"unsafe"
)

type MyInt int
type MyInt64 int64
type MyC complex128
type S struct{ a, b int }

func check(name string, ok bool) {
	if ok {
		fmt.Println("ok  ", name)
	} else {
		fmt.Println("FAIL", name)
	}
}

func main() {
	i := 5
	p := &i
	q := (*MyInt)(unsafe.Pointer(p))
	check("type", fmt.Sprintf("%T", q) == "*main.MyInt")
	_, isMy := any(q).(*MyInt)
	check("assert", isMy)
	check("read", *q == 5)
	*q = 7
	check("write", i == 7)
	check("round trip", (*int)(unsafe.Pointer(q)) == p)
	check("same twice", (*MyInt)(unsafe.Pointer(p)) == q)
	check("same as explicit", (*MyInt)(p) == q)
	check("pointee type", fmt.Sprintf("%T", any(*q)) == "main.MyInt")

	x := int64(1 << 40)
	r := (*MyInt64)(unsafe.Pointer(&x))
	check("int64 pointee type", fmt.Sprintf("%T %v", any(*r), *r) == "main.MyInt64 1099511627776")
	*r = 3
	check("int64 write", x == 3 && fmt.Sprintf("%T", any(x)) == "int64")
	u := (*uint64)(unsafe.Pointer(&x))
	x = -1
	check("int64 as uint64", *u == 1<<64-1)

	c := complex(1, 2)
	mc := (*MyC)(unsafe.Pointer(&c))
	check("complex", fmt.Sprintf("%T %v", any(*mc), *mc) == "main.MyC (1+2i)")

	pp := &p
	qq := (**MyInt)(unsafe.Pointer(pp))
	check("ptrptr type", fmt.Sprintf("%T %T", qq, *qq) == "**main.MyInt *main.MyInt")
	check("ptrptr read", **qq == 7 && *qq == q)
	var nilp *int
	np := &nilp
	check("ptrptr nil", *(**MyInt)(unsafe.Pointer(np)) == nil)

	s := &S{1, 2}
	check("struct same", (*S)(unsafe.Pointer(s)) == s)
	var up unsafe.Pointer
	check("nil", (*MyInt)(up) == nil && (*S)(up) == nil)
	arr := [3]int{1, 2, 3}
	ap := (*[3]MyInt)(unsafe.Pointer(&arr))
	check("array", ap[1] == 2)
}
