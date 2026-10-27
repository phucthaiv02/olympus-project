package main

import (
	"fmt"
	"unsafe"
)

type MyInt int
type MyInt64 int64
type P *int
type S struct{ A, B int }
type S2 struct{ A, B MyInt }

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
	check("D10a", func() bool { s := &S{1, 2}; return (*S2)(unsafe.Pointer(s)) == (*S2)(unsafe.Pointer(s)) })
	check("D10b", func() bool { s := &S{1, 2}; return (*S2)(unsafe.Pointer(&s.A)) == (*S2)(unsafe.Pointer(&s.A)) })
	check("D10c", func() bool { s := &S{1, 2}; a := (*S2)(unsafe.Pointer(s)); return fmt.Sprintf("%T", any(a.A)) == "main.MyInt" })
	check("MP1 map ptr identity", func() bool {
		m := map[string]int64{}
		q := (*map[string]MyInt64)(unsafe.Pointer(&m))
		return q == (*map[string]MyInt64)(unsafe.Pointer(&m)) && (*map[string]int64)(unsafe.Pointer(q)) == &m
	})
	check("B1 blocking P conv", func() bool {
		x := 1
		ch := make(chan *int, 1)
		done := make(chan bool)
		go func() {
			p := P(<-ch)
			var n *int
			ch2 := make(chan *int, 1)
			ch2 <- n
			q := P(<-ch2)
			done <- p == P(&x) && q == nil
		}()
		ch <- &x
		return <-done
	})
	check("B2 map[P] lookup by *int", func() bool {
		x := 1
		m := map[P]int{&x: 5}
		return m[&x] == 5 && m[P(&x)] == 5
	})
	check("B3 map[unsafe.Pointer] key", func() bool {
		x := 1
		m := map[unsafe.Pointer]int{unsafe.Pointer(&x): 5}
		return m[unsafe.Pointer((*MyInt)(&x))] == 5 && m[unsafe.Pointer(P(&x))] == 5
	})
	check("B4 map[unsafe.Pointer] key array", func() bool {
		a := [2]int64{}
		m := map[unsafe.Pointer]int{unsafe.Pointer(&a): 5}
		return m[unsafe.Pointer((*[2]MyInt64)(unsafe.Pointer(&a)))] == 5
	})
	check("B5 interface holding unsafe ptr", func() bool {
		x := 1
		var i any = unsafe.Pointer(&x)
		return i == any(unsafe.Pointer((*MyInt)(&x)))
	})
}
