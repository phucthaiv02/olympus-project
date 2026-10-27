package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type S struct{ X, Y int }

func check(name string, f func() bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("FAIL", name, "panic:", r)
		}
	}()
	if f() {
		fmt.Println("ok", name)
	} else {
		fmt.Println("FAIL", name)
	}
}

func main() {
	check("F1 *int -> *[1]int slice write", func() bool {
		n := 3
		a := (*[1]int)(unsafe.Pointer(&n))
		s := a[:]
		s[0] = 5
		return n == 5
	})
	check("F2 *int -> *[1]int deref assign", func() bool {
		n := 3
		a := (*[1]int)(unsafe.Pointer(&n))
		*a = [1]int{7}
		return n == 7
	})
	check("F3 *int -> *[1]int copy", func() bool {
		n := 3
		a := (*[1]int)(unsafe.Pointer(&n))
		copy(a[:], []int{8})
		return n == 8
	})
	check("F4 *int -> *[1]int reflect Index set", func() bool {
		n := 3
		a := (*[1]int)(unsafe.Pointer(&n))
		reflect.ValueOf(a).Elem().Index(0).SetInt(9)
		return n == 9
	})
	check("F5 &s.X -> *[2]int field elems", func() bool {
		s := S{1, 2}
		a := (*[2]int)(unsafe.Pointer(&s.X))
		return a[0] == 1 && a[1] == 2
	})
	check("F6 unsafe.Slice(&n,1) baseline", func() bool {
		n := 3
		s := unsafe.Slice(&n, 1)
		s[0] = 4
		return n == 4
	})
	check("F7 *int -> *[1]int &a[0] == &n", func() bool {
		n := 3
		a := (*[1]int)(unsafe.Pointer(&n))
		return &a[0] == &n
	})
	check("F8 *int -> *[1]int reflect Elem Index Addr", func() bool {
		n := 3
		a := (*[1]int)(unsafe.Pointer(&n))
		return reflect.ValueOf(a).Elem().Index(0).Addr().Interface().(*int) == &n
	})
}
