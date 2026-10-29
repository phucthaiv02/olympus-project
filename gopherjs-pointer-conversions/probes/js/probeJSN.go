package main

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/gopherjs/gopherjs/js"
)

type P *js.Object
type Q *js.Object

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
	o := js.Global.Get("Object").New()
	o.Set("n", 1)
	check("J1 P(o) == P(unsafe(o))", func() bool { return P(o) == P(unsafe.Pointer(o)) })
	check("J2 P(o) == P(o)", func() bool { return P(o) == P(o) })
	check("J3 %T P(o)", func() bool { return fmt.Sprintf("%T", P(o)) == "main.P" })
	check("J4 %T P(unsafe(o))", func() bool { return fmt.Sprintf("%T", P(unsafe.Pointer(o))) == "main.P" })
	check("J5 back (*js.Object)(P(o)) == o", func() bool { return (*js.Object)(P(o)) == o })
	check("J6 back via unsafe", func() bool { return (*js.Object)(unsafe.Pointer(P(o))) == o })
	check("J7 unsafe(P(o)) == unsafe(o)", func() bool { return unsafe.Pointer(P(o)) == unsafe.Pointer(o) })
	check("J8 iface assert P", func() bool { var i any = P(o); q, ok := i.(P); return ok && q == P(o) })
	check("J9 iface assert not *js.Object", func() bool { var i any = P(o); _, ok := i.(*js.Object); return !ok })
	check("J10 P -> Q", func() bool { return Q(P(o)) == Q(o) && fmt.Sprintf("%T", Q(P(o))) == "main.Q" })
	check("J11 implicit var P = o", func() bool { var q P = o; return q == P(o) })
	check("J12 implicit back var *js.Object = P", func() bool { var q *js.Object = P(o); return q == o })
	check("J13 nil P", func() bool { var n *js.Object; return P(n) == nil && P(unsafe.Pointer(n)) == nil && (*js.Object)(P(nil)) == nil })
	check("J14 reflect Convert", func() bool {
		v := reflect.ValueOf(o).Convert(reflect.TypeOf(P(nil)))
		q, ok := v.Interface().(P)
		return ok && q == P(o)
	})
	check("J15 reflect Set", func() bool {
		var q P
		reflect.ValueOf(&q).Elem().Set(reflect.ValueOf(o))
		return q == P(o)
	})
	check("J16 reflect NewAt-free Pointer eq", func() bool {
		return reflect.ValueOf(P(o)).Pointer() == reflect.ValueOf(o).Pointer()
	})
	check("J17 map key P", func() bool {
		m := map[P]int{}
		m[P(o)]++
		m[P(unsafe.Pointer(o))]++
		return len(m) == 1
	})
	check("J18 deref field through P", func() bool { return (*js.Object)(P(o)).Get("n").Int() == 1 })
	check("J19 o unchanged keys", func() bool { return js.Global.Get("Object").Call("keys", o).Length() == 1 })
}
