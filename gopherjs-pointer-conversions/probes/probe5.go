package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type MyInt int
type MyInt64 int64
type MyC complex128
type P *int
type SA struct{ X int }
type SB SA

func check(name string, f func() bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("FAIL", name, "(panic)")
		}
	}()
	if f() {
		fmt.Println("ok  ", name)
	} else {
		fmt.Println("FAIL", name)
	}
}

func main() {
	// A: non-nil unsafe.Pointer to a different pointer type
	check("A1 unsafe cross type %T", func() bool { i := 1; return fmt.Sprintf("%T", (*MyInt)(unsafe.Pointer(&i))) == "*main.MyInt" })
	check("A2 unsafe cross type assert", func() bool { i := 1; _, ok := any((*MyInt)(unsafe.Pointer(&i))).(*MyInt); return ok })
	check("A3 unsafe cross type same as explicit", func() bool { i := 1; p := &i; return (*MyInt)(unsafe.Pointer(p)) == (*MyInt)(p) })
	check("A4 unsafe cross type twice equal", func() bool { i := 1; p := &i; return (*MyInt)(unsafe.Pointer(p)) == (*MyInt)(unsafe.Pointer(p)) })
	check("A5 unsafe cross write", func() bool { i := 1; q := (*MyInt)(unsafe.Pointer(&i)); *q = 5; return i == 5 && *q == 5 })
	check("A6 unsafe int64 pointee type", func() bool { x := int64(3); r := (*MyInt64)(unsafe.Pointer(&x)); return fmt.Sprintf("%T", any(*r)) == "main.MyInt64" })
	check("A7 unsafe complex pointee type", func() bool { c := 1 + 2i; r := (*MyC)(unsafe.Pointer(&c)); return fmt.Sprintf("%T", any(*r)) == "main.MyC" })
	check("A8 unsafe **T pointee type", func() bool {
		i := 1; p := &i; qq := (**MyInt)(unsafe.Pointer(&p)); return fmt.Sprintf("%T", *qq) == "*main.MyInt" && *qq == (*MyInt)(p)
	})
	check("A9 int64 as uint64", func() bool { x := int64(-1); return *(*uint64)(unsafe.Pointer(&x)) == 1<<64-1 })
	// B: unsafe.Pointer(&s[i]) and new(T)
	check("B1 slice elem round trip", func() bool { s := []int{1, 2, 3}; return (*int)(unsafe.Pointer(&s[1])) == &s[1] })
	check("B2 slice elem write", func() bool { s := []int{1, 2, 3}; *(*int)(unsafe.Pointer(&s[1])) = 7; return s[1] == 7 })
	check("B3 array elem round trip", func() bool { a := [3]int{1, 2, 3}; return (*int)(unsafe.Pointer(&a[2])) == &a[2] })
	check("B4 slice elem via var", func() bool { s := []int{1, 2, 3}; up := unsafe.Pointer(&s[1]); return (*int)(up) == &s[1] && *(*int)(up) == 2 })
	check("B5 struct slice elem", func() bool { s := []SA{{1}, {2}}; return (*SA)(unsafe.Pointer(&s[1])) == &s[1] })
	check("B6 byte slice elem", func() bool { b := []byte{1, 2}; return (*byte)(unsafe.Pointer(&b[1])) == &b[1] })
	check("B7 eval once", func() bool {
		n := 0; s := []int{1, 2, 3}
		f := func() []int { n++; return s }
		q := (*int)(unsafe.Pointer(&f()[1])); return n == 1 && q == &s[1]
	})
	check("B8 new round trip", func() bool { q := (*int)(unsafe.Pointer(new(int))); *q = 3; return *q == 3 })
	check("B9 new struct", func() bool { q := (*SA)(unsafe.Pointer(new(SA))); q.X = 3; return q.X == 3 })
	// C: reflect.Value.Convert
	check("C1 reflect type", func() bool { i := 1; return fmt.Sprintf("%T", reflect.ValueOf(&i).Convert(reflect.TypeOf(P(nil))).Interface()) == "main.P" })
	check("C2 reflect identity", func() bool { i := 1; p := &i; return reflect.ValueOf(p).Convert(reflect.TypeOf(P(nil))).Interface().(P) == P(p) })
	check("C3 reflect twice", func() bool {
		i := 1; p := &i; t := reflect.TypeOf((*MyInt)(nil))
		return reflect.ValueOf(p).Convert(t).Interface() == reflect.ValueOf(p).Convert(t).Interface()
	})
	check("C4 reflect nil", func() bool { var p *int; return reflect.ValueOf(p).Convert(reflect.TypeOf(P(nil))).IsNil() })
	check("C5 reflect struct alias", func() bool {
		a := &SA{1}; b := reflect.ValueOf(a).Convert(reflect.TypeOf((*SB)(nil))).Interface().(*SB); b.X = 7; return a.X == 7 && b == (*SB)(a)
	})
	check("C6 reflect int64", func() bool {
		x := int64(3); r := reflect.ValueOf(&x).Convert(reflect.TypeOf((*MyInt64)(nil))).Interface().(*MyInt64); return r == (*MyInt64)(&x)
	})
	check("C7 reflect round trip", func() bool {
		i := 1; p := &i; q := reflect.ValueOf(p).Convert(reflect.TypeOf((*MyInt)(nil))); return q.Convert(reflect.TypeOf(p)).Interface().(*int) == p
	})
	check("C8 reflect array ptr nil", func() bool {
		var a *[2]int; return reflect.ValueOf(a).Convert(reflect.TypeOf((*[2]int)(nil))).IsNil()
	})
	check("C9 reflect ptrptr", func() bool {
		i := 1; p := &i; type PP **int; r := reflect.ValueOf(&p).Convert(reflect.TypeOf(PP(nil))).Interface().(PP); return (**int)(r) == &p
	})
}
