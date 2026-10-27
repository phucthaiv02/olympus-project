package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type MyInt64 int64
type Src [1]*[1]int64
type Dst [1]*[1]MyInt64

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

// nilWrite converts u to *Dst, writes nil and non-nil elements through it and
// checks that the original array keeps its own types.
func nilWrite(a *Src, u unsafe.Pointer) bool {
	a[0] = &[1]int64{1}
	q := (*Dst)(u)
	if q[0] == nil || q[0][0] != 1 || fmt.Sprintf("%T", q[0]) != "*[1]main.MyInt64" {
		return false
	}
	q[0] = nil
	if a[0] != nil || q[0] != nil {
		return false
	}
	q[0] = &[1]MyInt64{2}
	return a[0] != nil && a[0][0] == 2 && fmt.Sprintf("%T", a[0]) == "*[1]int64"
}

func toUnsafe[P any](p *P) unsafe.Pointer { return unsafe.Pointer(p) }

type ptrTo[E any] interface{ *E }

func toUnsafeT[E any, PE ptrTo[E]](p PE) unsafe.Pointer { return unsafe.Pointer(p) }

func main() {
	check("W1 direct", func() bool { a := &Src{}; return nilWrite(a, unsafe.Pointer(a)) })
	check("W2 reflect UnsafePointer", func() bool { a := &Src{}; return nilWrite(a, reflect.ValueOf(a).UnsafePointer()) })
	check("W3 reflect Pointer", func() bool {
		a := &Src{}
		return reflect.ValueOf(a).Pointer() == reflect.ValueOf(a).Pointer() && nilWrite(a, reflect.ValueOf(a).UnsafePointer())
	})
	check("W4 reflect UnsafeAddr", func() bool {
		var s struct{ a Src }
		v := reflect.ValueOf(&s).Elem().Field(0)
		return nilWrite(&s.a, unsafe.Pointer(v.UnsafeAddr()))
	})
	check("W5 generic *P", func() bool { a := &Src{}; return nilWrite(a, toUnsafe(a)) })
	check("W6 generic core type", func() bool { a := &Src{}; return nilWrite(a, toUnsafeT[Src](a)) })
	check("W7 reflect Elem.Addr.UnsafePointer", func() bool {
		a := &Src{}
		return nilWrite(a, reflect.ValueOf(a).Elem().Addr().UnsafePointer())
	})
	check("W8 reflect Convert unsafe->named unsafe", func() bool {
		type U unsafe.Pointer
		a := &Src{}
		u := reflect.ValueOf(unsafe.Pointer(a)).Convert(reflect.TypeOf(U(nil))).Interface().(U)
		return nilWrite(a, unsafe.Pointer(u))
	})
	check("W9 named unsafe type", func() bool {
		type U unsafe.Pointer
		a := &Src{}
		return nilWrite(a, unsafe.Pointer(U(unsafe.Pointer(a))))
	})
	check("W10 struct field array", func() bool {
		type RS struct{ a Src }
		type NS struct{ a Dst }
		s := &RS{}
		n := (*NS)(unsafe.Pointer(s))
		s.a[0] = &[1]int64{1}
		n.a[0] = nil
		return s.a[0] == nil
	})
	check("W11 via SetPointer", func() bool {
		a := &Src{}
		var u unsafe.Pointer
		reflect.ValueOf(&u).Elem().SetPointer(unsafe.Pointer(a))
		return nilWrite(a, u)
	})
	check("W12 nested array view", func() bool {
		var a [1][1]*[1]int64
		a[0][0] = &[1]int64{1}
		q := (*[1][1]*[1]MyInt64)(unsafe.Pointer(&a))
		q[0][0] = nil
		return a[0][0] == nil
	})
	check("R1 slice UnsafePointer == &s[0]", func() bool {
		s := []int{1, 2}
		return reflect.ValueOf(s).UnsafePointer() == unsafe.Pointer(&s[0])
	})
	check("R2 slice UnsafePointer == SliceData", func() bool {
		s := []int64{1, 2}
		return reflect.ValueOf(s).UnsafePointer() == unsafe.Pointer(unsafe.SliceData(s))
	})
	check("R3 named slice view UnsafePointer", func() bool {
		s := []int64{1, 2}
		ms := *(*[]MyInt64)(unsafe.Pointer(&s))
		return reflect.ValueOf(ms).UnsafePointer() == reflect.ValueOf(s).UnsafePointer()
	})
	check("R4 array view Index Addr UnsafeAddr", func() bool {
		a := [2]int64{}
		q := (*[2]MyInt64)(unsafe.Pointer(&a))
		return reflect.ValueOf(q).Elem().Index(1).UnsafeAddr() == uintptr(unsafe.Pointer(&a[1]))
	})
}
