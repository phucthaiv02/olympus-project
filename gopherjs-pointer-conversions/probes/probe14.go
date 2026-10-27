package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type S struct{ X int }
type S2 struct{ A, B int32 }
type S3 struct{ A, B, C int }

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
	st := reflect.TypeOf(S{})
	check("N1 NewAt struct from int", func() bool {
		n := 7
		v := reflect.NewAt(st, unsafe.Pointer(&n))
		p, ok := v.Interface().(*S)
		return ok && fmt.Sprintf("%T", v.Interface()) == "*main.S" && p.X == 7 && v.Type() == reflect.PtrTo(st)
	})
	check("N2 NewAt write through", func() bool {
		n := 7
		p := reflect.NewAt(st, unsafe.Pointer(&n)).Interface().(*S)
		p.X = 9
		return n == 9
	})
	check("N3 NewAt round trip", func() bool {
		n := 7
		v := reflect.NewAt(st, unsafe.Pointer(&n))
		return v.UnsafePointer() == unsafe.Pointer(&n) && (*int)(v.UnsafePointer()) == &n
	})
	check("N4 Convert unsafe -> *S", func() bool {
		n := 7
		v := reflect.ValueOf(unsafe.Pointer(&n)).Convert(reflect.TypeOf(unsafe.Pointer(nil)))
		p := (*S)(v.Interface().(unsafe.Pointer))
		return fmt.Sprintf("%T", p) == "*main.S" && p.X == 7
	})
	check("N5 direct (*S)(unsafe(&n))", func() bool {
		n := 7
		p := (*S)(unsafe.Pointer(&n))
		var i any = p
		_, ok := i.(*S)
		return ok && p.X == 7 && (*int)(unsafe.Pointer(p)) == &n && p == (*S)(unsafe.Pointer(&n))
	})
	check("N6 (*int)(unsafe(&s)) first field", func() bool {
		s := S3{1, 2, 3}
		p := (*int)(unsafe.Pointer(&s))
		*p = 5
		return fmt.Sprintf("%T", p) == "*int" && s.A == 5 && *p == 5 && (*S3)(unsafe.Pointer(p)) == &s
	})
	check("N7 struct -> array", func() bool {
		s := S2{1, 2}
		a := (*[2]int32)(unsafe.Pointer(&s))
		a[1] = 7
		return fmt.Sprintf("%T", a) == "*[2]int32" && a[0] == 1 && s.B == 7 && (*S2)(unsafe.Pointer(a)) == &s
	})
	check("N8 array -> struct", func() bool {
		a := [2]int32{3, 4}
		s := (*S2)(unsafe.Pointer(&a))
		s.A = 9
		return fmt.Sprintf("%T", s) == "*main.S2" && s.B == 4 && a[0] == 9 && (*[2]int32)(unsafe.Pointer(s)) == &a
	})
	check("N9 element -> struct", func() bool {
		a := []int32{1, 2, 3, 4}
		s := (*S2)(unsafe.Pointer(&a[2]))
		return s.A == 3 && s.B == 4 && unsafe.Pointer(s) == unsafe.Pointer(&a[2])
	})
	check("N10 struct element -> int", func() bool {
		ss := []S{{1}, {2}}
		p := (*int)(unsafe.Pointer(&ss[1]))
		return *p == 2 && fmt.Sprintf("%T", p) == "*int"
	})
	check("N11 NewAt int from struct", func() bool {
		s := S{4}
		v := reflect.NewAt(reflect.TypeOf(0), unsafe.Pointer(&s))
		return fmt.Sprintf("%T", v.Interface()) == "*int" && v.Elem().Int() == 4
	})
	check("N12 *int -> *[2]int extra elems", func() bool {
		n := 1
		a := (*[2]int)(unsafe.Pointer(&n))
		return fmt.Sprintf("%T", a) == "*[2]int" && a[0] == 1
	})
	check("N13 nil stays nil", func() bool {
		var n *int
		return (*S)(unsafe.Pointer(n)) == nil && reflect.NewAt(st, unsafe.Pointer(n)).IsNil()
	})
}
