package main

import (
	"fmt"
	"unsafe"
)

type Inner struct{ X, Y int }
type MyInner Inner
type Outer struct {
	In Inner
	Z  int
}
type MyInt int
type S struct {
	A [2]int
	B int
}

func try(name string, f func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(name, "PANIC:", r)
		}
	}()
	f()
}

func main() {
	try("nested struct", func() {
		o := &Outer{In: Inner{1, 2}, Z: 3}
		v0 := (*MyInner)(unsafe.Pointer(&o.In))
		fmt.Println("before:", v0.X, v0.Y)
		m := map[unsafe.Pointer]bool{unsafe.Pointer(o): true}
		_ = m
		v := (*MyInner)(unsafe.Pointer(&o.In))
		fmt.Println("after:", v.X, v.Y, unsafe.Pointer(v) == unsafe.Pointer(&o.In))
		v.Y = 20
		fmt.Println("write:", o.In.Y, o.Z)
	})
	try("nested struct fresh", func() {
		o := &Outer{In: Inner{1, 2}, Z: 3}
		_ = unsafe.Pointer(o)
		v := (*MyInner)(unsafe.Pointer(&o.In))
		fmt.Println("fresh:", v.X, v.Y)
	})
	try("array of struct", func() {
		arr := [2]S{{A: [2]int{1, 2}, B: 3}, {A: [2]int{4, 5}, B: 6}}
		_ = unsafe.Pointer(&arr)
		v := (*[2]MyInt)(unsafe.Pointer(&arr[0].A))
		fmt.Println("aos:", v[0], v[1])
	})
}
