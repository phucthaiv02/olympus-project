package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type MyInt64 int64
type MyInt int
type A [2]int64
type S struct{ X, Y int }
type S2 struct{ X, Y MyInt }
type B struct{ X int }
type Outer struct{ *B }

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

func gderef[PT ~*int](p PT) *int { return &*p }

func main() {
	check("P1 elem0 array ptr == &a (typed)", func() bool {
		var a [4]int
		return (*[4]int)(unsafe.Pointer(&a[0])) == &a
	})
	check("P2 elem0 array ptr == &a (string)", func() bool {
		var a [2]string
		return (*[2]string)(unsafe.Pointer(&a[0])) == &a
	})
	check("P3 elem0 array ptr == &a (int64)", func() bool {
		var a [2]int64
		return (*[2]int64)(unsafe.Pointer(&a[0])) == &a
	})
	check("P4 unsafe(&a) == unsafe(&a[0]) via array ptr", func() bool {
		var a [4]int
		p := (*[4]int)(unsafe.Pointer(&a[0]))
		return unsafe.Pointer(p) == unsafe.Pointer(&a)
	})
	check("P5 known-type order: (*[2]int64)(&a) twice", func() bool {
		var a A
		p1 := (*[2]int64)(&a)
		_ = unsafe.Pointer(&a)
		p2 := (*[2]int64)(&a)
		return p1 == p2
	})
	check("P6 round trip after named record", func() bool {
		var a [2]int64
		pa := (*A)(&a)
		_ = unsafe.Pointer(pa)
		p := (*[2]int64)(unsafe.Pointer(&a))
		return p == &a
	})
	check("P7 round trip after reflect Pointer named", func() bool {
		var a [2]int64
		_ = reflect.ValueOf((*A)(&a)).Pointer()
		p := (*[2]int64)(unsafe.Pointer(&a))
		return p == &a
	})
	check("P8 overlap copy through offset ptr int64", func() bool {
		a := [4]int64{1, 2, 3, 4}
		p := (*[3]int64)(unsafe.Pointer(&a[1]))
		copy(p[:], a[:3])
		return a == [4]int64{1, 1, 2, 3}
	})
	check("P9 overlap copy through view", func() bool {
		a := [4]int64{1, 2, 3, 4}
		v := (*[4]MyInt64)(unsafe.Pointer(&a))
		copy(v[1:], v[:3])
		return a == [4]int64{1, 1, 2, 3}
	})
	check("P10 overlap copy view vs orig", func() bool {
		a := [4]int64{1, 2, 3, 4}
		v := (*[4]int64)(unsafe.Pointer((*[4]MyInt64)(unsafe.Pointer(&a))))
		_ = v
		s := (*[4]MyInt64)(unsafe.Pointer(&a))[:]
		t := *(*[]int64)(unsafe.Pointer(&s))
		copy(t[1:], a[:3])
		return a == [4]int64{1, 1, 2, 3}
	})
	check("P11 *[2]int elem of 2D -> *[2][2]int", func() bool {
		var a [3][2]int
		a[2][1] = 7
		p := (*[2][2]int)(unsafe.Pointer(&a[1]))
		return p[1][1] == 7 && &p[0] == &a[1]
	})
	check("P12 *S -> *[1]S", func() bool {
		s := S{1, 2}
		p := (*[1]S)(unsafe.Pointer(&s))
		return p[0].Y == 2 && &p[0] == &s
	})
	check("P13 *[2]int -> *[1][2]int", func() bool {
		x := [2]int{1, 2}
		p := (*[1][2]int)(unsafe.Pointer(&x))
		return p[0][1] == 2
	})
	check("P14 typed subarray cached elem ptr then wider conv", func() bool {
		a := [4]int{1, 2, 3, 4}
		p := (*[2]int)(unsafe.Pointer(&a[1]))
		q := &p[1]
		r := (*[2]int)(unsafe.Pointer(&a[2]))
		return q == &a[2] && r[1] == 4
	})
	check("P15 typed subarray cached elem ptr then unsafe.Slice", func() bool {
		a := [4]int{1, 2, 3, 4}
		p := (*[2]int)(unsafe.Pointer(&a[1]))
		_ = &p[1]
		s := unsafe.Slice(&a[2], 2)
		return s[1] == 4
	})
	check("P16 map value type via unsafe", func() bool {
		m := map[string]int64{"a": 1}
		pm := (*map[string]MyInt64)(unsafe.Pointer(&m))
		v := (*pm)["a"]
		_, ok := any(v).(MyInt64)
		return ok
	})
	check("P17 generic &*p", func() bool {
		type P *int
		x := 1
		var p P = &x
		q := gderef(p)
		return q == &x && *q == 1
	})
	check("P18 promoted field via embedded ptr", func() bool {
		o := Outer{&B{}}
		p := &o.X
		*p = 3
		return p == &o.B.X && o.B.X == 3
	})
	check("P19 promoted field via embedded ptr unsafe", func() bool {
		o := Outer{&B{}}
		return unsafe.Pointer(&o.X) == unsafe.Pointer(&o.B.X)
	})
	check("P20 reflect Field Addr reserved name", func() bool {
		type R struct{ new, x int }
		var r R
		return reflect.ValueOf(&r).Elem().Field(0).Addr().UnsafePointer() == unsafe.Pointer(&r.new)
	})
	check("P21 *[]int -> *[]MyInt then subslice data ptr", func() bool {
		s := []int64{1, 2, 3}
		ps := (*[]MyInt64)(unsafe.Pointer(&s))
		t := (*ps)[1:]
		return unsafe.Pointer(&t[0]) == unsafe.Pointer(&s[1])
	})
	check("P22 view slice cap0 SliceData", func() bool {
		s := []int{1, 2, 3}
		t := s[3:]
		return unsafe.SliceData(t) == unsafe.SliceData(s)
	})
	check("P23 reflect Slice3 cap0 Pointer", func() bool {
		s := []int{1, 2, 3}
		t := s[2:2:2]
		u := reflect.ValueOf(s).Slice3(2, 2, 2)
		return u.Pointer() == reflect.ValueOf(t).Pointer()
	})
	check("P24 struct view S2 field addr eq", func() bool {
		s := S{1, 2}
		v := (*S2)(unsafe.Pointer(&s))
		return unsafe.Pointer(&v.Y) == unsafe.Pointer(&s.Y)
	})
}
