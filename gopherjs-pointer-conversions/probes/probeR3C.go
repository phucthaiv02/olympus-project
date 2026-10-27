package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type Pair struct{ A, B int32 }
type PI struct{ A, B int }
type S struct{ X int }

func try(name string, f func() string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(name, "panic:", r)
		}
	}()
	fmt.Println(name, f())
}

func main() {
	try("C1 %T %v", func() string {
		x := Pair{1, 2}
		p := (*[2]int32)(unsafe.Pointer(&x))
		return fmt.Sprintf("%T %v", p, p)
	})
	try("C2 p[:]", func() string {
		x := Pair{1, 2}
		p := (*[2]int32)(unsafe.Pointer(&x))
		s := p[:]
		return fmt.Sprint(len(s), s)
	})
	try("C3 &p[1]", func() string {
		x := Pair{1, 2}
		p := (*[2]int32)(unsafe.Pointer(&x))
		return fmt.Sprint(&p[1] == &p[1], unsafe.Pointer(&p[1]) == unsafe.Pointer(&x.B), *&p[1])
	})
	try("C4 reflect", func() string {
		x := Pair{1, 2}
		p := (*[2]int32)(unsafe.Pointer(&x))
		v := reflect.ValueOf(p).Elem()
		v.Index(1).SetInt(6)
		return fmt.Sprint(x.B, v.Len(), v.Index(0).Int())
	})
	try("C5 int elems %v", func() string {
		x := PI{1, 2}
		p := (*[2]int)(unsafe.Pointer(&x))
		s := p[:]
		s[0] = 5
		return fmt.Sprintf("%v %v %d", p, s, x.A)
	})
	try("C6 int elems copy", func() string {
		x := PI{1, 2}
		p := (*[2]int)(unsafe.Pointer(&x))
		*p = [2]int{7, 8}
		return fmt.Sprint(x)
	})
	try("C7 int -> [2]int %v", func() string {
		x := 3
		p := (*[1]int)(unsafe.Pointer(&x))
		return fmt.Sprintf("%v", p)
	})
	try("C8 struct->[1]int32 slice", func() string {
		x := Pair{1, 2}
		p := (*[1]int32)(unsafe.Pointer(&x))
		return fmt.Sprintf("%v %v", p, p[:])
	})
	try("C9 &s.X ->[1]int", func() string {
		s := S{4}
		p := (*[1]int)(unsafe.Pointer(&s.X))
		return fmt.Sprintf("%v %v", p, p[:])
	})
	try("C10 s -> [1]int", func() string {
		s := S{4}
		p := (*[1]int)(unsafe.Pointer(&s))
		return fmt.Sprintf("%v %v %v", p, p[:], &p[0] == &s.X)
	})
}
