package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type MyInt64 int64
type MyInt int
type S struct{ X int }
type SI64 struct{ A, B int64 }
type SM64 struct{ A, B MyInt64 }
type SM1 struct{ A MyInt64 }
type Inner struct{ A int }
type Outer struct{ In Inner }
type Emb struct{ Inner }
type Base struct{ ID int }
type Derived struct {
	Base
	Name string
}
type WithPtr struct{ P *MyInt }
type WithArr struct{ A [2]int }

func (s *S) Get() int { return s.X }

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
	check("A1 *S -> *[1]S elem", func() bool {
		s := S{5}
		a := (*[1]S)(unsafe.Pointer(&s))
		return a[0].X == 5
	})
	check("A2 *S -> *[1]S write+identity", func() bool {
		s := S{5}
		a := (*[1]S)(unsafe.Pointer(&s))
		a[0].X = 7
		return s.X == 7 && &a[0] == &s
	})
	check("A3 *[2]int -> *[1][2]int", func() bool {
		x := [2]int{1, 2}
		a := (*[1][2]int)(unsafe.Pointer(&x))
		return a[0][1] == 2 && fmt.Sprintf("%T", a[0]) == "[2]int"
	})
	check("A4 *int64 -> *SM1 field type", func() bool {
		x := int64(3)
		s := (*SM1)(unsafe.Pointer(&x))
		var i any = s.A
		_, ok := i.(MyInt64)
		return ok && s.A == 3
	})
	check("A5 *[2]int64 -> *SM64 field type", func() bool {
		x := [2]int64{3, 4}
		s := (*SM64)(unsafe.Pointer(&x))
		var i any = s.B
		_, ok := i.(MyInt64)
		return ok && s.B == 4
	})
	check("A6 *SI64 -> *[2]MyInt64 elem type", func() bool {
		x := SI64{3, 4}
		s := (*[2]MyInt64)(unsafe.Pointer(&x))
		var i any = s[1]
		_, ok := i.(MyInt64)
		return ok && s[1] == 4
	})
	check("A7 *SI64 -> *MyInt64 type", func() bool {
		x := SI64{3, 4}
		p := (*MyInt64)(unsafe.Pointer(&x))
		var i any = *p
		_, ok := i.(MyInt64)
		return ok && *p == 3
	})
	check("A8 *int -> *Outer nested", func() bool {
		x := 5
		o := (*Outer)(unsafe.Pointer(&x))
		return o.In.A == 5
	})
	check("A9 *Inner -> *Emb", func() bool {
		in := Inner{5}
		e := (*Emb)(unsafe.Pointer(&in))
		return e.A == 5 && e.Inner.A == 5
	})
	check("A10 container_of", func() bool {
		d := &Derived{Base{1}, "n"}
		b := &d.Base
		back := (*Derived)(unsafe.Pointer(b))
		return back.Name == "n" && back == d
	})
	check("A11 (*int)(unsafe(&arr)) == &arr[0]", func() bool {
		arr := [3]int{1, 2, 3}
		p := (*int)(unsafe.Pointer(&arr))
		return p == &arr[0]
	})
	check("A12 (*int)(unsafe(&s)) == &s.X", func() bool {
		s := S{1}
		p := (*int)(unsafe.Pointer(&s))
		return p == &s.X
	})
	check("A13 (*S)(unsafe(&s.X)) == &s", func() bool {
		s := S{1}
		p := (*S)(unsafe.Pointer(&s.X))
		return p == &s
	})
	check("A14 foreign struct method", func() bool {
		x := 9
		p := (*S)(unsafe.Pointer(&x))
		return p.Get() == 9
	})
	check("A15 foreign struct deref copy", func() bool {
		x := 9
		p := (*S)(unsafe.Pointer(&x))
		v := *p
		x = 10
		return v.X == 9 && v == S{9}
	})
	check("A16 foreign struct assign whole", func() bool {
		x := 9
		p := (*S)(unsafe.Pointer(&x))
		*p = S{11}
		return x == 11
	})
	check("A17 foreign struct reflect", func() bool {
		x := 9
		p := (*S)(unsafe.Pointer(&x))
		v := reflect.ValueOf(p).Elem()
		v.Field(0).SetInt(4)
		return x == 4 && fmt.Sprint(p) == "&{4}" && v.Field(0).Addr().Interface().(*int) == &p.X
	})
	check("A18 foreign &p.X identity", func() bool {
		x := 9
		p := (*S)(unsafe.Pointer(&x))
		return &p.X == &p.X && *&p.X == 9 && unsafe.Pointer(&p.X) == unsafe.Pointer(&x)
	})
	check("A19 *[1]*int -> *WithPtr type", func() bool {
		n := 3
		a := [1]*int{&n}
		w := (*WithPtr)(unsafe.Pointer(&a))
		var i any = w.P
		_, ok := i.(*MyInt)
		return ok && *w.P == 3
	})
	check("A20 *[2]int -> *WithArr", func() bool {
		a := [2]int{1, 2}
		w := (*WithArr)(unsafe.Pointer(&a))
		return w.A[1] == 2 && &w.A == &a
	})
	check("A21 *int -> *[1]MyInt typ", func() bool {
		n := 3
		a := (*[1]MyInt)(unsafe.Pointer(&n))
		return a[0] == 3 && fmt.Sprintf("%T", a) == "*[1]main.MyInt"
	})
	check("A22 *int64 -> *[1]MyInt64 elem type", func() bool {
		n := int64(3)
		a := (*[1]MyInt64)(unsafe.Pointer(&n))
		var i any = a[0]
		_, ok := i.(MyInt64)
		return ok && a[0] == 3
	})
	check("A23 *[1]int64 -> *MyInt64 type", func() bool {
		n := [1]int64{3}
		a := (*MyInt64)(unsafe.Pointer(&n))
		var i any = *a
		_, ok := i.(MyInt64)
		return ok && *a == 3
	})
	check("A24 NewAt MyInt64 from [1]int64", func() bool {
		n := [1]int64{3}
		v := reflect.NewAt(reflect.TypeOf(MyInt64(0)), unsafe.Pointer(&n))
		return v.Elem().Interface() == any(MyInt64(3))
	})
	check("A25 *S -> *[1]S fmt", func() bool {
		s := S{5}
		a := (*[1]S)(unsafe.Pointer(&s))
		return fmt.Sprint(a) == "&[{5}]"
	})
}
