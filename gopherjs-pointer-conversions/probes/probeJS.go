package main

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/gopherjs/gopherjs/js"
)

type U struct {
	X int
	O *js.Object
}
type T struct {
	O *js.Object
	X int
}

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

func newObj(n int) *js.Object {
	o := js.Global.Get("Object").New()
	o.Set("n", n)
	return o
}

func main() {
	check("F1 Field Addr == &s.O", func() bool {
		s := U{1, newObj(1)}
		p := reflect.ValueOf(&s).Elem().Field(1).Addr().Interface().(**js.Object)
		return p == &s.O && *p == s.O
	})
	check("F2 Field UnsafeAddr", func() bool {
		s := U{1, newObj(1)}
		return unsafe.Pointer(reflect.ValueOf(&s).Elem().Field(1).UnsafeAddr()) == unsafe.Pointer(&s.O)
	})
	check("F3 &s.O first, then reflect", func() bool {
		s := U{1, newObj(1)}
		q := &s.O
		return reflect.ValueOf(&s).Elem().Field(1).Addr().Interface().(**js.Object) == q
	})
	check("F4 Field read/Interface", func() bool {
		s := U{1, newObj(7)}
		return reflect.ValueOf(&s).Elem().Field(1).Interface().(*js.Object) == s.O && reflect.ValueOf(&s).Elem().Field(1).Interface().(*js.Object).Get("n").Int() == 7
	})
	check("F5 Field Set", func() bool {
		s := U{1, newObj(1)}
		o := newObj(2)
		reflect.ValueOf(&s).Elem().Field(1).Set(reflect.ValueOf(o))
		return s.O == o
	})
	check("F6 Addr().Elem() read/write", func() bool {
		s := U{1, newObj(1)}
		e := reflect.ValueOf(&s).Elem().Field(1).Addr().Elem()
		o := newObj(3)
		ok := e.Interface().(*js.Object) == s.O
		e.Set(reflect.ValueOf(o))
		return ok && s.O == o
	})
	check("F7 write through Addr pointer", func() bool {
		s := U{1, newObj(1)}
		p := reflect.ValueOf(&s).Elem().Field(1).Addr().Interface().(**js.Object)
		o := newObj(4)
		*p = o
		return s.O == o
	})
	check("F8 first field js.Object", func() bool {
		s := T{newObj(1), 2}
		v := reflect.ValueOf(&s).Elem().Field(0)
		return v.Addr().Interface().(**js.Object) == &s.O && unsafe.Pointer(v.UnsafeAddr()) == unsafe.Pointer(&s) && unsafe.Pointer(&s.O) == unsafe.Pointer(&s)
	})
	check("A1 array Index Addr", func() bool {
		a := [2]*js.Object{newObj(1), newObj(2)}
		v := reflect.ValueOf(&a).Elem().Index(1)
		return v.Addr().Interface().(**js.Object) == &a[1] && unsafe.Pointer(v.UnsafeAddr()) == unsafe.Pointer(&a[1]) && v.Interface().(*js.Object) == a[1]
	})
	check("A2 array Index Set + Addr().Elem()", func() bool {
		a := [2]*js.Object{newObj(1), newObj(2)}
		o := newObj(5)
		reflect.ValueOf(&a).Elem().Index(1).Addr().Elem().Set(reflect.ValueOf(o))
		reflect.ValueOf(&a).Elem().Index(0).Set(reflect.ValueOf(o))
		return a[1] == o && a[0] == o
	})
	check("S1 slice Index Addr", func() bool {
		s := []*js.Object{newObj(1), newObj(2)}
		v := reflect.ValueOf(s).Index(1)
		return v.Addr().Interface().(**js.Object) == &s[1] && unsafe.Pointer(v.UnsafeAddr()) == unsafe.Pointer(&s[1]) && v.Interface().(*js.Object) == s[1]
	})
	check("S2 slice Index Set", func() bool {
		s := []*js.Object{newObj(1), newObj(2)}
		o := newObj(6)
		reflect.ValueOf(s).Index(1).Set(reflect.ValueOf(o))
		*reflect.ValueOf(s).Index(0).Addr().Interface().(**js.Object) = o
		return s[1] == o && s[0] == o
	})
	check("N1 nil js.Object field", func() bool {
		var s U
		v := reflect.ValueOf(&s).Elem().Field(1)
		return v.IsNil() && v.Addr().Interface().(**js.Object) == &s.O
	})
}
