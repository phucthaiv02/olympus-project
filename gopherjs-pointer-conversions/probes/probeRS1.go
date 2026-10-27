package main

import (
	"fmt"
	"unsafe"
)

type S struct {
	A int32
	B uint8
}

type MyInt64 int64

var _zero uintptr

func try(name string, f func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(name, "PANIC", r)
		}
	}()
	f()
}

func cond(b bool) bool { return b }

func main() {
	try("T1 byte elem vs ptr var", func() {
		s := []byte{1, 2, 3}
		p := &s[1]
		fmt.Println("T1", unsafe.Pointer(&s[1]) == unsafe.Pointer(p))
	})
	try("T2 _zero vs ptr var", func() {
		p := &_zero
		fmt.Println("T2", unsafe.Pointer(&_zero) == unsafe.Pointer(p))
	})
	try("T3 struct elem then struct ptr", func() {
		ss := []S{{1, 2}, {3, 4}}
		a := unsafe.Pointer(&ss[0])
		q := &ss[0]
		b := unsafe.Pointer(q)
		fmt.Println("T3", a == b)
	})
	try("T4 struct ptr then struct elem", func() {
		ss := []S{{1, 2}, {3, 4}}
		q := &ss[0]
		b := unsafe.Pointer(q)
		a := unsafe.Pointer(&ss[0])
		fmt.Println("T4", a == b)
	})
	try("T5 short-circuit struct conv", func() {
		var ps *S
		ok := ps != nil && unsafe.Pointer(ps) != nil
		fmt.Println("T5", ok)
	})
	try("T6 stale delayed write-back", func() {
		x := S{A: 1}
		for i := 0; i < 2; i++ {
			if i == 1 {
				x.A = 100
			}
			c := cond(i == 0) && unsafe.Pointer(&x) != nil
			_ = c
		}
		fmt.Println("T6", x.A)
	})
	try("T7 int64 elem write through", func() {
		w := []int64{1, 2}
		p := unsafe.Pointer(&w[0])
		w[0] = 5
		fmt.Println("T7", *(*int64)(p), (*int64)(p) == &w[0])
	})
	try("T8 names-only deref in syscall", func() {
		x := [2]MyInt64{1, 2}
		y := *(*[2]int64)(unsafe.Pointer(&x))
		var e interface{} = y[0]
		fmt.Printf("T8 %T\n", e)
	})
	try("T9 struct ptr eq after field change", func() {
		x := S{A: 1}
		u := unsafe.Pointer(&x)
		x.A = 9
		fmt.Println("T9", (*S)(u).A, (*S)(u) == &x)
	})
	try("T10 elem of array var vs ptr", func() {
		var arr [4]byte
		p := &arr[2]
		fmt.Println("T10", unsafe.Pointer(&arr[2]) == unsafe.Pointer(p))
	})
	try("T11 struct elem roundtrip", func() {
		ss := []S{{1, 2}}
		r := (*S)(unsafe.Pointer(&ss[0]))
		fmt.Println("T11", r == &ss[0])
	})
}
