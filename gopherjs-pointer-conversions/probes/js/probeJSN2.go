package main

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/gopherjs/gopherjs/js"
)

type P *js.Object
type SO struct{ O *js.Object }
type SP struct{ O P }

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
	o2 := js.Global.Get("Object").New()
	check("K1 [2]*js.Object as [2]P read", func() bool {
		a := [2]*js.Object{o, nil}
		v := (*[2]P)(unsafe.Pointer(&a))
		return v[0] == P(o) && v[1] == nil && fmt.Sprintf("%T", v[0]) == "main.P"
	})
	check("K2 [2]P write keeps raw", func() bool {
		a := [2]*js.Object{o, nil}
		v := (*[2]P)(unsafe.Pointer(&a))
		v[1] = P(o2)
		v[0] = nil
		return a[1] == o2 && a[0] == nil && a[1].Get("constructor") == js.Global.Get("Object")
	})
	check("K3 [2]P as [2]*js.Object", func() bool {
		a := [2]P{P(o), nil}
		v := (*[2]*js.Object)(unsafe.Pointer(&a))
		ok := v[0] == o && v[1] == nil
		v[1] = o2
		return ok && a[1] == P(o2)
	})
	check("K4 struct view field", func() bool {
		s := &SO{o}
		v := (*SP)(unsafe.Pointer(s))
		ok := v.O == P(o) && fmt.Sprintf("%T", v.O) == "main.P"
		v.O = P(o2)
		ok = ok && s.O == o2
		v.O = nil
		return ok && s.O == nil
	})
	check("K5 js.Undefined as P", func() bool {
		u := js.Undefined
		q := P(u)
		return q != nil && (*js.Object)(q) == u && P(unsafe.Pointer(u)) == q
	})
	check("K6 primitive as P", func() bool {
		n := o.Get("n")
		q := P(n)
		return q != nil && (*js.Object)(q).Int() == 1 && P(n) == q
	})
	check("K7 fmt of *js.Object unchanged", func() bool { return fmt.Sprintf("%T", o) == "*js.Object" })
	check("K8 reflect Elem Field of P", func() bool {
		f := reflect.ValueOf(P(o)).Elem().Field(0)
		return f.Interface().(*js.Object) == o
	})
	check("K9 reflect Elem Field of *js.Object", func() bool {
		f := reflect.ValueOf(o).Elem().Field(0)
		return f.Interface().(*js.Object) == o
	})
	check("K10 NewAt js.Object", func() bool {
		v := reflect.NewAt(reflect.TypeOf(js.Object{}), unsafe.Pointer(o))
		return v.Interface().(*js.Object) == o
	})
	check("K11 reflect Convert P -> *js.Object", func() bool {
		v := reflect.ValueOf(P(o)).Convert(reflect.TypeOf((*js.Object)(nil)))
		return v.Interface().(*js.Object) == o
	})
	check("K12 reflect Set *js.Object from P", func() bool {
		var q *js.Object
		reflect.ValueOf(&q).Elem().Set(reflect.ValueOf(o))
		return q == o
	})
	check("K13 DeepEqual", func() bool { return reflect.DeepEqual(P(o), P(o)) })
	check("K14 nil P via unsafe zero", func() bool {
		var u unsafe.Pointer
		return P(u) == nil && (*js.Object)(u) == nil
	})
	check("K15 slice of *js.Object as []P via unsafe.Slice", func() bool {
		a := [2]*js.Object{o, o2}
		s := unsafe.Slice((*P)(unsafe.Pointer(&a[0])), 2)
		return s[1] == P(o2)
	})
	fmt.Println("print", P(o) != nil, fmt.Sprint(o.Get("n")))
}
