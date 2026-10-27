package main

import (
	"fmt"
	"reflect"
	"sync/atomic"
	"unsafe"
)

type MyInt int
type MyInt64 int64
type P *int
type U unsafe.Pointer
type S struct{ A, B int }
type S2 struct{ A, B MyInt }

var gx int
var gp P = &gx
var gq = P(&gx)
var gu = unsafe.Pointer(&gx)

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

func takeP(p P) string { return fmt.Sprintf("%T", p) }

type Holder struct{ p P }

func (h *Holder) set(p P) { h.p = p }

func main() {
	check("G1 package-level P", func() bool { return gp == gq && fmt.Sprintf("%T", gp) == "main.P" && gu == unsafe.Pointer(gp) })
	check("U1 named unsafe from ptr", func() bool { x := 1; return U(&x) == U((*MyInt)(&x)) })
	check("U2 named unsafe back", func() bool {
		x := 1
		u := U(&x)
		return (*MyInt)(unsafe.Pointer(u)) == (*MyInt)(&x)
	})
	check("D1 defer arg P", func() bool {
		x := 1
		var got string
		func() { defer func(p P) { got = fmt.Sprintf("%T %v", p, p == P(&x)) }(&x) }()
		return got == "main.P true"
	})
	check("D2 method value arg", func() bool {
		x := 1
		h := &Holder{}
		f := h.set
		f(&x)
		return h.p == P(&x) && fmt.Sprintf("%T", h.p) == "main.P"
	})
	check("D3 nil *int to P in func value", func() bool {
		var f func(P) bool = func(p P) bool { return p == nil }
		var n *int
		return f(n)
	})
	check("D4 closure return P nil", func() bool {
		return func() P { var n *int; return n }() == nil
	})
	check("D5 array of P", func() bool { x := 1; a := [...]P{&x, nil}; return a[0] == P(&x) && a[1] == nil })
	check("D6 append P", func() bool {
		x := 1
		var n *int
		s := append([]P(nil), &x, n)
		return s[0] == P(&x) && s[1] == nil
	})
	check("D7 atomic.Pointer", func() bool {
		x := MyInt(1)
		var ap atomic.Pointer[MyInt]
		ap.Store(&x)
		y := 0
		_ = y
		return ap.CompareAndSwap((*MyInt)(unsafe.Pointer(&x)), nil) && ap.Load() == nil
	})
	check("D8 switch on P", func() bool {
		x := 1
		p := P(&x)
		switch p {
		case &x:
			return true
		}
		return false
	})
	check("D9 escaping var unsafe", func() bool {
		x := 1
		f := func() { x++ }
		f()
		return (*MyInt)(unsafe.Pointer(&x)) == (*MyInt)(&x) && unsafe.Pointer(P(&x)) == unsafe.Pointer(&x)
	})
	check("D10 *S2 from *S equal", func() bool {
		s := &S{1, 2}
		a := (*S2)(unsafe.Pointer(s))
		b := (*S2)(unsafe.Pointer(&s.A))
		return a == b && fmt.Sprintf("%T", any(a.A)) == "main.MyInt"
	})
	check("R1 reflect Convert to U", func() bool {
		x := 1
		v := reflect.ValueOf(&x).Convert(reflect.TypeOf(U(nil))).Interface().(U)
		return v == U(&x)
	})
	check("R2 reflect Convert U to *MyInt", func() bool {
		x := 1
		v := reflect.ValueOf(U(&x)).Convert(reflect.TypeOf(unsafe.Pointer(nil))).Interface().(unsafe.Pointer)
		return (*MyInt)(v) == (*MyInt)(&x)
	})
	check("R3 Value.Set unsafe.Pointer field", func() bool {
		x := 1
		var u unsafe.Pointer
		reflect.ValueOf(&u).Elem().Set(reflect.ValueOf(unsafe.Pointer((*MyInt)(&x))))
		return u == unsafe.Pointer(&x)
	})
	check("R4 reflect Pointer of P", func() bool {
		x := 1
		return reflect.ValueOf(P(&x)).Pointer() == reflect.ValueOf(&x).Pointer()
	})
	check("R5 reflect Field Addr of P field", func() bool {
		x := 1
		h := &Holder{p: &x}
		return reflect.ValueOf(h).Elem().Field(0).Addr().Interface().(*P) == &h.p
	})
	check("R6 reflect Index Addr of []P", func() bool {
		x := 1
		s := []P{&x}
		return reflect.ValueOf(s).Index(0).Addr().Interface().(*P) == &s[0]
	})
	check("R7 Convert *[]int64 to *IS", func() bool {
		type IS []int64
		s := []int64{1}
		q := reflect.ValueOf(&s).Convert(reflect.TypeOf((*IS)(nil))).Interface().(*IS)
		return q == (*IS)(&s)
	})
}
