package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type MyInt int
type MyInt64 int64
type P *int
type PM *MyInt

func (m *MyInt) T() string { return fmt.Sprintf("%T", m) }

type S struct{ X int }

func (s *S) T() string { return fmt.Sprintf("%T", s) }

type PS *S

type SA struct{ X int64 }
type SB struct{ X MyInt64 }

type E struct{ G int }
type OA struct{ E }
type OB struct{ E }

var gi = 1

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
	check("A1 MakeFunc result type", func() bool {
		f := reflect.MakeFunc(reflect.TypeOf((func() P)(nil)), func([]reflect.Value) []reflect.Value {
			return []reflect.Value{reflect.ValueOf(&gi)}
		}).Interface().(func() P)
		q := f()
		return fmt.Sprintf("%T", q) == "main.P" && q == P(&gi)
	})
	check("A2 MakeFunc nil result", func() bool {
		f := reflect.MakeFunc(reflect.TypeOf((func() P)(nil)), func([]reflect.Value) []reflect.Value {
			return []reflect.Value{reflect.ValueOf((*int)(nil))}
		}).Interface().(func() P)
		return f() == nil
	})
	check("B1 Elem.Addr of P type", func() bool {
		v := reflect.ValueOf(P(&gi)).Elem().Addr().Interface()
		_, ok := v.(*int)
		return ok && fmt.Sprintf("%T", v) == "*int"
	})
	check("C1 &*p of P", func() bool {
		p := P(&gi)
		var x any = &*p
		_, ok := x.(*int)
		return ok
	})
	check("C2 &*p of PS", func() bool {
		s := &S{}
		p := PS(s)
		var x any = &*p
		_, ok := x.(*S)
		return ok
	})
	check("D1 method via PM", func() bool {
		var m MyInt
		pm := PM(&m)
		return (*pm).T() == "*main.MyInt"
	})
	check("D2 method via PS", func() bool {
		ps := PS(&S{})
		return (*ps).T() == "*main.S"
	})
	check("E1 slice ptr-to-array nil write", func() bool {
		s := []*[1]int64{{1}}
		t := *(*[]*[1]MyInt64)(unsafe.Pointer(&s))
		t[0] = nil
		return s[0] == nil
	})
	check("E2 slice ptr-to-array elem addr", func() bool {
		s := []*[1]int64{{1}}
		t := *(*[]*[1]MyInt64)(unsafe.Pointer(&s))
		return unsafe.Pointer(&t[0]) == unsafe.Pointer(&s[0])
	})
	check("E3 slice ptr-to-array elem addr typed", func() bool {
		s := []*[1]int64{{1}}
		t := *(*[]*[1]MyInt64)(unsafe.Pointer(&s))
		return (**[1]int64)(unsafe.Pointer(&t[0])) == &s[0]
	})
	check("E4 array via slice of arrays nil write", func() bool {
		s := [][1]*[1]int64{{{1}}}
		t := *(*[][1]*[1]MyInt64)(unsafe.Pointer(&s))
		t[0][0] = nil
		return s[0][0] == nil
	})
	check("F1 unsafe.Slice of struct view", func() bool {
		as := []SA{{1}, {2}}
		q := (*SB)(unsafe.Pointer(&as[0]))
		sl := unsafe.Slice(q, 2)
		sl[1].X = 5
		return as[1].X == 5 && fmt.Sprintf("%T", any(sl[1].X)) == "main.MyInt64"
	})
	check("F2 unsafe.Slice of explicit struct conv", func() bool {
		type S2 S
		as := []S{{1}, {2}}
		sl := unsafe.Slice((*S2)(&as[0]), 2)
		sl[1].X = 5
		return as[1].X == 5
	})
	check("G1 unsafe.Slice of int64 conv type", func() bool {
		a := [2]int64{1, 2}
		sl := unsafe.Slice((*MyInt64)(unsafe.Pointer(&a[0])), 2)
		return fmt.Sprintf("%T", any(sl[1])) == "main.MyInt64"
	})
	check("G2 unsafe.Slice of int64 conv write", func() bool {
		a := [2]int64{1, 2}
		sl := unsafe.Slice((*MyInt64)(unsafe.Pointer(&a[0])), 2)
		sl[1] = 7
		return a[1] == 7 && fmt.Sprintf("%T", any(a[1])) == "int64"
	})
	check("G3 unsafe.Slice of view elem", func() bool {
		a := [2]int64{1, 2}
		q := (*[2]MyInt64)(unsafe.Pointer(&a))
		sl := unsafe.Slice(&q[0], 2)
		return fmt.Sprintf("%T", any(sl[1])) == "main.MyInt64"
	})
	check("H1 promoted field addr", func() bool {
		var s OA
		return &s.G == &s.E.G
	})
	check("H2 promoted field unsafe", func() bool {
		var s OA
		return unsafe.Pointer(&s.G) == unsafe.Pointer(&s.E.G)
	})
	check("H3 promoted field via view", func() bool {
		var s OA
		v := (*OB)(unsafe.Pointer(&s))
		return unsafe.Pointer(&v.G) == unsafe.Pointer(&s.G)
	})
	check("H4 promoted field via explicit view", func() bool {
		var s OA
		v := (*OB)(&s)
		return unsafe.Pointer(&v.G) == unsafe.Pointer(&s.E.G)
	})
	check("H5 reflect FieldByName addr", func() bool {
		var s OA
		return reflect.ValueOf(&s).Elem().FieldByName("G").Addr().Pointer() == uintptr(unsafe.Pointer(&s.G))
	})
	check("I1 Convert nil addressable", func() bool {
		var p *int
		v := reflect.ValueOf(&p).Elem().Convert(reflect.TypeOf(P(nil)))
		return v.IsNil() && !v.CanSet()
	})
	check("I2 Convert nil addressable Interface", func() bool {
		var p *int
		v := reflect.ValueOf(&p).Elem().Convert(reflect.TypeOf(P(nil)))
		return v.Interface().(P) == nil
	})
	check("J1 append grow struct view aliasing", func() bool {
		s := []SA{{1}}
		t := *(*[]SB)(unsafe.Pointer(&s))
		t2 := append(t, SB{2})
		t2[0].X = 9
		return s[0].X == 1
	})
	check("J0 append grow plain struct baseline", func() bool {
		s := []SA{{1}}
		s2 := append(s, SA{2})
		s2[0].X = 9
		return s[0].X == 1
	})
}
