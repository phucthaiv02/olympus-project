package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type N int64

func c(name string, f func() bool) {
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
	c("X1 direct row[0] then view", func() bool {
		var a [2][2]int64
		r := (*[2]N)(unsafe.Pointer(&a[0]))
		v := (*[2][2]N)(unsafe.Pointer(&a))
		return r == &v[0]
	})
	c("X2 view then direct row[1]", func() bool {
		var a [2][2]int64
		v := (*[2][2]N)(unsafe.Pointer(&a))
		r := (*[2]N)(unsafe.Pointer(&a[1]))
		return r == &v[1]
	})
	c("X3 two direct row conversions", func() bool {
		var a [2][2]int64
		return (*[2]N)(unsafe.Pointer(&a[1])) == (*[2]N)(unsafe.Pointer(&a[1]))
	})
	c("X4 round trip direct row", func() bool {
		var a [2][2]int64
		r := (*[2]N)(unsafe.Pointer(&a[1]))
		return (*[2]int64)(unsafe.Pointer(r)) == &a[1]
	})
	c("X5 view row back", func() bool {
		var a [2][2]int64
		v := (*[2][2]N)(unsafe.Pointer(&a))
		return (*[2]int64)(unsafe.Pointer(&v[1])) == &a[1]
	})
	c("X6 reflect row vs direct-converted", func() bool {
		var a [2][2]int64
		v := (*[2][2]N)(unsafe.Pointer(&a))
		r := (*[2]N)(unsafe.Pointer(&a[1]))
		return reflect.ValueOf(v).Elem().Index(1).Addr().Interface().(*[2]N) == r
	})
	c("X7 reflect row vs view row", func() bool {
		var a [2][2]int64
		v := (*[2][2]N)(unsafe.Pointer(&a))
		return reflect.ValueOf(v).Elem().Index(1).Addr().Interface().(*[2]N) == &v[1]
	})
	c("X8 reflect on source vs view row back", func() bool {
		var a [2][2]int64
		q := reflect.ValueOf(&a).Elem().Index(1).Addr().Interface().(*[2]int64)
		v := (*[2][2]N)(unsafe.Pointer(&a))
		return (*[2]int64)(unsafe.Pointer(&v[1])) == q
	})
	c("X9 write through direct row", func() bool {
		a := [2][2]int64{{1, 2}, {3, 4}}
		r := (*[2]N)(unsafe.Pointer(&a[1]))
		v := (*[2][2]N)(unsafe.Pointer(&a))
		r[0] = 7
		return a[1][0] == 7 && v[1][0] == 7
	})
	c("X10 write through view row", func() bool {
		a := [2][2]int64{{1, 2}, {3, 4}}
		v := (*[2][2]N)(unsafe.Pointer(&a))
		r := (*[2]N)(unsafe.Pointer(&a[1]))
		(&v[1])[1] = 8
		return r[1] == 8 && a[1][1] == 8
	})
	c("X11 3D rows", func() bool {
		var a [2][2][2]int64
		v := (*[2][2][2]N)(unsafe.Pointer(&a))
		return (*[2]N)(unsafe.Pointer(&a[1][1])) == &v[1][1] && (*[2][2]N)(unsafe.Pointer(&a[1])) == &v[1]
	})
	c("X12 struct array field row", func() bool {
		type S struct{ m [2][2]int64 }
		type U struct{ m [2][2]N }
		var s S
		u := (*U)(unsafe.Pointer(&s))
		return (*[2]N)(unsafe.Pointer(&s.m[1])) == &u.m[1]
	})
	c("X13 named row array", func() bool {
		type R [2]int64
		var a [2][2]int64
		v := (*[2]R)(unsafe.Pointer(&a))
		return (*R)(unsafe.Pointer(&a[1])) == &v[1]
	})
}
