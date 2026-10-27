package main

import (
	"fmt"
	"reflect"
	"sort"
	"unsafe"
)

type MyInt int
type MyInt64 int64
type S []int64
type SA struct{ X int64 }
type SB struct{ X MyInt64 }
type A64 [2]int64
type B64 [2]MyInt64

func (b *B64) Set() string { b[0] = 5; return fmt.Sprintf("%T", b[0]) }

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
	check("Z1 named slice conv keeps data ptr (cap 0)", func() bool {
		s := []int64{1, 2, 3}
		x := s[1:][1:1:1]
		return unsafe.SliceData(S(x)) == unsafe.SliceData(x)
	})
	check("Z2 implicit named slice keeps data ptr (cap 0)", func() bool {
		s := []int64{1, 2, 3}
		x := s[1:][1:1:1]
		var y S = x
		return reflect.ValueOf(y).Pointer() == reflect.ValueOf(x).Pointer()
	})
	check("Z3 unsafe reinterpret keeps data ptr (cap 0)", func() bool {
		s := []int64{1, 2, 3}
		x := s[1:][1:1:1]
		t := *(*[]MyInt64)(unsafe.Pointer(&x))
		return unsafe.Pointer(unsafe.SliceData(t)) == unsafe.Pointer(unsafe.SliceData(x))
	})
	check("Z4 unsafe.Slice len 0 then reinterpret", func() bool {
		s := []int{1, 2, 3}
		x := unsafe.Slice(&s[2], 0)
		t := *(*[]MyInt)(unsafe.Pointer(&x))
		return unsafe.Pointer(unsafe.SliceData(t)) == unsafe.Pointer(&s[2])
	})
	check("Z5 reflect Convert slice keeps ptr (cap 0)", func() bool {
		s := []int64{1, 2, 3}
		x := s[1:][1:1:1]
		y := reflect.ValueOf(x).Convert(reflect.TypeOf(S(nil)))
		return y.Pointer() == reflect.ValueOf(x).Pointer()
	})
	check("Q1 Swapper on struct view slice", func() bool {
		s := []SA{{1}, {2}}
		t := *(*[]SB)(unsafe.Pointer(&s))
		reflect.Swapper(t)(0, 1)
		return s[0].X == 2 && s[1].X == 1
	})
	check("Q2 sort.Slice on struct view slice", func() bool {
		s := []SA{{2}, {1}}
		t := *(*[]SB)(unsafe.Pointer(&s))
		sort.Slice(t, func(i, j int) bool { return t[i].X < t[j].X })
		return s[0].X == 1 && s[1].X == 2
	})
	check("Q3 Swapper on nested array view slice", func() bool {
		s := [][1]int64{{1}, {2}}
		t := *(*[][1]MyInt64)(unsafe.Pointer(&s))
		reflect.Swapper(t)(0, 1)
		return s[0][0] == 2 && s[1][0] == 1
	})
	check("Q4 tuple swap on struct view slice", func() bool {
		s := []SA{{1}, {2}}
		t := *(*[]SB)(unsafe.Pointer(&s))
		t[0], t[1] = t[1], t[0]
		return s[0].X == 2 && s[1].X == 1
	})
	check("R1 method on named array view", func() bool {
		a := A64{1, 2}
		q := (*B64)(unsafe.Pointer(&a))
		r := q.Set()
		return r == "main.MyInt64" && a[0] == 5 && fmt.Sprintf("%T", any(a[0])) == "int64"
	})
	check("R2 SliceData of struct view slice + unsafe.Slice", func() bool {
		s := []SA{{1}, {2}}
		t := *(*[]SB)(unsafe.Pointer(&s))
		u := unsafe.Slice(unsafe.SliceData(t), 2)
		return u[1].X == 2
	})
	check("R3 reflect Index Addr of struct view elem + unsafe.Slice", func() bool {
		s := []SA{{1}, {2}}
		t := *(*[]SB)(unsafe.Pointer(&s))
		p := reflect.ValueOf(t).Index(0).Addr().Interface().(*SB)
		u := unsafe.Slice(p, 2)
		return u[1].X == 2
	})
}
