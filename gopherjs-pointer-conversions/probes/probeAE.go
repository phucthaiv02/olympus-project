package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type MyInt int
type MyInt64 int64
type In struct{ A int }
type S struct {
	X int
	Y string
}
type SN struct {
	In In
	B  int
}
type SA struct {
	Arr [2]int
	B   int
}
type SP struct {
	P *int
	B int
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

func main() {
	check("A1 (*int)(up(&a)) == &a[0]", func() bool {
		a := [3]int{1, 2, 3}
		return (*int)(unsafe.Pointer(&a)) == &a[0]
	})
	check("A1b &a[0] first", func() bool {
		a := [3]int{1, 2, 3}
		q := &a[0]
		return (*int)(unsafe.Pointer(&a)) == q
	})
	check("A2 up(&a[0]) == up(&a)", func() bool {
		a := [3]int{1, 2, 3}
		return unsafe.Pointer(&a[0]) == unsafe.Pointer(&a)
	})
	check("A2b up(&a) first", func() bool {
		a := [3]int{1, 2, 3}
		u := unsafe.Pointer(&a)
		return unsafe.Pointer(&a[0]) == u
	})
	check("A3 reflect Pointer", func() bool {
		a := [3]int{1, 2, 3}
		return reflect.ValueOf(&a[0]).Pointer() == reflect.ValueOf(&a).Pointer() &&
			reflect.ValueOf(&a).Elem().Index(0).Addr().Pointer() == reflect.ValueOf(&a).Pointer()
	})
	check("A4 (*MyInt)(up(&a)) == (*MyInt)(up(&a[0]))", func() bool {
		a := [3]int{1, 2, 3}
		return (*MyInt)(unsafe.Pointer(&a)) == (*MyInt)(unsafe.Pointer(&a[0])) && (*MyInt)(unsafe.Pointer(&a)) == (*MyInt)(&a[0])
	})
	check("A5 write through", func() bool {
		a := [3]int{1, 2, 3}
		p := (*int)(unsafe.Pointer(&a))
		*p = 9
		a[0]++
		return a[0] == 10 && *p == 10
	})
	check("A6 int64 typed", func() bool {
		a := [2]int64{1, 2}
		p := (*MyInt64)(unsafe.Pointer(&a))
		*p = 5
		return a[0] == 5 && p == (*MyInt64)(&a[0]) && fmt.Sprintf("%T", *p) == "main.MyInt64" && (*int64)(unsafe.Pointer(&a)) == &a[0]
	})
	check("A7 int32 typed array", func() bool {
		a := [2]int32{1, 2}
		return (*int32)(unsafe.Pointer(&a)) == &a[0] && unsafe.Pointer(&a[0]) == unsafe.Pointer(&a)
	})
	check("A8 string array", func() bool {
		a := [2]string{"x", "y"}
		return (*string)(unsafe.Pointer(&a)) == &a[0] && *(*string)(unsafe.Pointer(&a)) == "x"
	})
	check("A9 2D", func() bool {
		var a [2][2]int
		return (*int)(unsafe.Pointer(&a)) == &a[0][0] && (*[2]int)(unsafe.Pointer(&a)) == &a[0] && unsafe.Pointer(&a[0][0]) == unsafe.Pointer(&a)
	})
	check("A10 elem0 -> array -> elem0", func() bool {
		a := [3]int{1, 2, 3}
		pa := (*[3]int)(unsafe.Pointer(&a[0]))
		return pa == &a && (*int)(unsafe.Pointer(pa)) == &a[0]
	})
	check("A11 slice backing", func() bool {
		s := []int{1, 2, 3}
		return (*int)(unsafe.Pointer((*[3]int)(s))) == &s[0]
	})
	check("A12 map key", func() bool {
		a := [3]int{1, 2, 3}
		m := map[unsafe.Pointer]int{unsafe.Pointer(&a): 1}
		m[unsafe.Pointer(&a[0])]++
		return len(m) == 1
	})
	check("S1 (*int)(up(&s)) == &s.X", func() bool {
		s := S{1, "y"}
		return (*int)(unsafe.Pointer(&s)) == &s.X
	})
	check("S1b &s.X first", func() bool {
		s := S{1, "y"}
		q := &s.X
		return (*int)(unsafe.Pointer(&s)) == q
	})
	check("S2 up(&s.X) == up(&s)", func() bool {
		s := S{1, "y"}
		return unsafe.Pointer(&s.X) == unsafe.Pointer(&s)
	})
	check("S3 reflect Pointer", func() bool {
		s := S{1, "y"}
		return reflect.ValueOf(&s.X).Pointer() == reflect.ValueOf(&s).Pointer() &&
			reflect.ValueOf(&s).Elem().Field(0).Addr().Pointer() == reflect.ValueOf(&s).Pointer()
	})
	check("S4 (*MyInt)(up(&s))", func() bool {
		s := S{1, "y"}
		p := (*MyInt)(unsafe.Pointer(&s))
		*p = 3
		return s.X == 3 && p == (*MyInt)(&s.X)
	})
	check("S5 nested struct", func() bool {
		var s SN
		return (*int)(unsafe.Pointer(&s)) == &s.In.A && (*In)(unsafe.Pointer(&s)) == &s.In && unsafe.Pointer(&s.In.A) == unsafe.Pointer(&s)
	})
	check("S6 struct first field array", func() bool {
		var s SA
		return (*int)(unsafe.Pointer(&s)) == &s.Arr[0] && (*[2]int)(unsafe.Pointer(&s)) == &s.Arr && unsafe.Pointer(&s.Arr) == unsafe.Pointer(&s)
	})
	check("S7 first field pointer", func() bool {
		n := 1
		s := SP{&n, 2}
		return (**int)(unsafe.Pointer(&s)) == &s.P && *(**int)(unsafe.Pointer(&s)) == &n
	})
	check("S8 back from field to struct", func() bool {
		s := S{1, "y"}
		return (*S)(unsafe.Pointer(&s.X)) == &s
	})
	check("S9 array of structs elem0 field", func() bool {
		a := [2]S{{1, "a"}, {2, "b"}}
		return (*int)(unsafe.Pointer(&a)) == &a[0].X && unsafe.Pointer(&a[0].X) == unsafe.Pointer(&a)
	})
}
