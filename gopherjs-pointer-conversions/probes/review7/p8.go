package main

import (
	"fmt"
	"unsafe"
)

type MyInt64 int64

func up[T any](p *T) unsafe.Pointer { return unsafe.Pointer(p) }
func back[T any](p unsafe.Pointer) *T { return (*T)(p) }

func try(name string, f func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(name, "PANIC:", r)
		}
	}()
	f()
}

func main() {
	try("generic iface arr", func() {
		arr := [2]any{int64(1), "x"}
		_ = up(&arr)
		p := &arr[0]
		*p = MyInt64(3)
		fmt.Printf("g1: %T %v\n", arr[0], up(&arr) == unsafe.Pointer(&arr[0]))
	})
	try("generic arr of arrptr", func() {
		x := [2]int{}
		arr := [2]*[2]int{&x, nil}
		a := up(&arr)
		b := unsafe.Pointer(&arr)
		fmt.Println("g2:", a == b, up(&arr) == unsafe.Pointer(&arr[0]))
	})
	try("generic arr of struct ptr", func() {
		type S struct{ X int }
		arr := [2]*S{{1}, nil}
		a := up(&arr)
		fmt.Println("g3:", a == unsafe.Pointer(&arr), a == unsafe.Pointer(&arr[0]), back[[2]*S](a) == &arr)
	})
	try("generic struct", func() {
		type S struct {
			A [2]*[2]int
			B int
		}
		s := &S{}
		a := up(s)
		fmt.Println("g4:", a == unsafe.Pointer(s), a == unsafe.Pointer(&s.A[0]))
	})
	try("generic back", func() {
		arr := [2]int{1, 2}
		q := back[[2]int](unsafe.Pointer(&arr[0]))
		fmt.Println("g5:", q == &arr, back[int](up(&arr)) == &arr[0])
	})
}
