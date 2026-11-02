// Probes for the expert review of 2026-11-02. Build tag free: js parts are in js.go for GopherJS.
package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type T struct{ x int }
type PT *T
type PI *int
type MyInt int
type N int64
type Row *[2]N
type Grid *[2][2]N

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
	// &*p on non-nil named pointers.
	check("D1 %T &*PT", func() bool { x := T{1}; p := PT(&x); return fmt.Sprintf("%T", &*p) == "*main.T" })
	check("D2 &*PT == &x", func() bool { x := T{1}; p := PT(&x); return &*p == &x })
	check("D3 %T &*PI", func() bool { n := 1; p := PI(&n); return fmt.Sprintf("%T", &*p) == "*int" })
	check("D4 &*PI == &n", func() bool { n := 1; p := PI(&n); return &*p == &n })
	check("D5 &*PI == (*int)(p)", func() bool { n := 1; p := PI(&n); return &*p == (*int)(p) })
	check("D6 &*PI write", func() bool { n := 1; p := PI(&n); q := &*p; *q = 5; return n == 5 })

	// Writes through a converted **T.
	check("U1 **MyInt write keeps *int", func() bool {
		n, m := 1, MyInt(2)
		v := &n
		pp := (**MyInt)(unsafe.Pointer(&v))
		*pp = &m
		var i interface{} = v
		_, ok := i.(*int)
		return ok && fmt.Sprintf("%T", v) == "*int" && *v == 2
	})
	check("U2 **MyInt read has *MyInt", func() bool {
		n := 1
		v := &n
		pp := (**MyInt)(unsafe.Pointer(&v))
		return fmt.Sprintf("%T", *pp) == "*main.MyInt" && int(**pp) == 1
	})
	check("U3 **MyInt write aliases and identity", func() bool {
		n, m := 1, MyInt(2)
		v := &n
		pp := (**MyInt)(unsafe.Pointer(&v))
		*pp = &m
		return *pp == &m && (*int)(unsafe.Pointer(&m)) == v
	})
	check("U4 **T write of nil", func() bool {
		n := 1
		v := &n
		pp := (**MyInt)(unsafe.Pointer(&v))
		*pp = nil
		return v == nil && *pp == nil && fmt.Sprintf("%T", v) == "*int"
	})
	check("U5 struct field **", func() bool {
		type A struct{ p *int }
		type B struct{ p *MyInt }
		n, m := 1, MyInt(2)
		a := A{&n}
		b := (*B)(unsafe.Pointer(&a))
		b.p = &m
		return fmt.Sprintf("%T", a.p) == "*int" && *a.p == 2
	})

	// Elem().Addr() from a named pointer.
	check("A1 Elem().Addr() type PT", func() bool { x := T{1}; p := PT(&x); return reflect.ValueOf(p).Elem().Addr().Type() == reflect.TypeOf(&x) })
	check("A2 Elem().Addr() == &x PT", func() bool { x := T{1}; p := PT(&x); q, ok := reflect.ValueOf(p).Elem().Addr().Interface().(*T); return ok && q == &x })
	check("A3 Elem().Addr() type PI", func() bool { n := 1; p := PI(&n); return reflect.ValueOf(p).Elem().Addr().Type() == reflect.TypeOf(&n) })
	check("A4 Elem().Addr() == &n PI", func() bool { n := 1; p := PI(&n); q, ok := reflect.ValueOf(p).Elem().Addr().Interface().(*int); return ok && q == &n })
	check("A5 named array", func() bool {
		a := [2]N{1, 2}
		p := Row(&a)
		q, ok := reflect.ValueOf(p).Elem().Addr().Interface().(*[2]N)
		return ok && q == &a
	})

	// Rows of a converted multidimensional array.
	check("M1 rows of 2 conversions", func() bool {
		var a [2][2]int64
		v1 := (*[2][2]N)(unsafe.Pointer(&a))
		v2 := (*[2][2]N)(unsafe.Pointer(&a))
		return &v1[0] == &v2[0] && &v1[1] == &v2[1]
	})
	check("M2 row vs row converted directly", func() bool {
		var a [2][2]int64
		v := (*[2][2]N)(unsafe.Pointer(&a))
		r := (*[2]N)(unsafe.Pointer(&a[0]))
		return r == &v[0]
	})
	check("M3 row back to source", func() bool {
		var a [2][2]int64
		v := (*[2][2]N)(unsafe.Pointer(&a))
		return (*[2]int64)(unsafe.Pointer(&v[1])) == &a[1]
	})
	check("M4 row unsafe equal", func() bool {
		var a [2][2]int64
		v := (*[2][2]N)(unsafe.Pointer(&a))
		return unsafe.Pointer(&v[1]) == unsafe.Pointer(&a[1])
	})
	check("M5 reflect Index.Addr vs direct", func() bool {
		var a [2][2]int64
		v := (*[2][2]N)(unsafe.Pointer(&a))
		q, ok := reflect.ValueOf(v).Elem().Index(1).Addr().Interface().(*[2]N)
		return ok && q == &v[1]
	})
	check("M6 reflect Index.Addr pointer value", func() bool {
		var a [2][2]int64
		v := (*[2][2]N)(unsafe.Pointer(&a))
		return reflect.ValueOf(v).Elem().Index(1).Addr().Pointer() == reflect.ValueOf(&a[1]).Pointer()
	})
	check("M7 elements share and read typed", func() bool {
		a := [2][2]int64{{1, 2}, {3, 4}}
		v := (*[2][2]N)(unsafe.Pointer(&a))
		v[1][1] = 9
		var i interface{} = a[1][1]
		_, ok := i.(int64)
		return ok && a[1][1] == 9 && fmt.Sprintf("%T", v[1][1]) == "main.N"
	})
	check("M8 element pointers", func() bool {
		var a [2][2]int64
		v := (*[2][2]N)(unsafe.Pointer(&a))
		return (*int64)(unsafe.Pointer(&v[1][1])) == &a[1][1]
	})
	check("M9 named Grid rows", func() bool {
		var a [2][2]int64
		v := Grid(unsafe.Pointer(&a))
		w := (*[2][2]N)(unsafe.Pointer(&a))
		return &v[1] == &w[1]
	})
	check("M10 unconverted 2D baseline", func() bool {
		var a [2][2]int64
		p := &a
		return &p[1] == &a[1] && reflect.ValueOf(p).Elem().Index(1).Addr().Interface().(*[2]int64) == &a[1]
	})
}
