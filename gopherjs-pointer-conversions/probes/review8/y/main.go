package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type N int64
type M int64
type R [2]int64
type S struct{ v int64 }
type S2 struct{ v N }

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
	c("Y1 slice of arrays via unsafe.Slice", func() bool {
		s := make([][2]int64, 2)
		v := unsafe.Slice((*[2]N)(unsafe.Pointer(&s[0])), 2)
		return (*[2]N)(unsafe.Pointer(&s[1])) == &v[1]
	})
	c("Y2 [1][2] over a row", func() bool {
		var a [2][2]int64
		x := (*[1][2]N)(unsafe.Pointer(&a[1]))
		return &x[0] == (*[2]N)(unsafe.Pointer(&a[1])) && (*[2]int64)(unsafe.Pointer(&x[0])) == &a[1]
	})
	c("Y3 range rows", func() bool {
		var a [2][2]int64
		v := (*[2][2]N)(unsafe.Pointer(&a))
		for i := range a {
			if (*[2]N)(unsafe.Pointer(&a[i])) != &v[i] {
				return false
			}
		}
		return true
	})
	c("Y4 assign whole row through pointer", func() bool {
		a := [2][2]int64{{1, 2}, {3, 4}}
		v := (*[2][2]N)(unsafe.Pointer(&a))
		r := &v[1]
		*r = [2]N{5, 6}
		var i interface{} = a[1][0]
		_, ok := i.(int64)
		return ok && a[1] == [2]int64{5, 6} && *r == [2]N{5, 6}
	})
	c("Y5 assign whole row by index", func() bool {
		a := [2][2]int64{{1, 2}, {3, 4}}
		v := (*[2][2]N)(unsafe.Pointer(&a))
		_ = &a[1]
		v[1] = [2]N{5, 6}
		var i interface{} = a[1][0]
		_, ok := i.(int64)
		return ok && a[1] == [2]int64{5, 6}
	})
	c("Y6 from named to raw", func() bool {
		var a [2][2]N
		v := (*[2][2]int64)(unsafe.Pointer(&a))
		return (*[2]int64)(unsafe.Pointer(&a[1])) == &v[1] && (*[2]N)(unsafe.Pointer(&v[1])) == &a[1]
	})
	c("Y7 two named views", func() bool {
		var a [2][2]int64
		v := (*[2][2]N)(unsafe.Pointer(&a))
		w := (*[2][2]M)(unsafe.Pointer(v))
		return (*[2]M)(unsafe.Pointer(&v[1])) == &w[1] && (*[2]N)(unsafe.Pointer(&w[1])) == &v[1]
	})
	c("Y8 row element pointers after row access", func() bool {
		var a [2][2]int64
		_ = &a[1]
		_ = &a[1][1]
		v := (*[2][2]N)(unsafe.Pointer(&a))
		return &v[1][1] == (*N)(unsafe.Pointer(&a[1][1])) && fmt.Sprintf("%T", &v[1][1]) == "*main.N"
	})
	c("Y9 array of structs rows", func() bool {
		var a [2][2]S
		v := (*[2][2]S2)(unsafe.Pointer(&a))
		return (*[2]S2)(unsafe.Pointer(&a[1])) == &v[1] && (*S2)(unsafe.Pointer(&a[1][1])) == &v[1][1]
	})
	c("Y10 array of pointers rows", func() bool {
		var a [2][2]*int
		v := (*[2][2]*N)(unsafe.Pointer(&a))
		return (*[2]*N)(unsafe.Pointer(&a[1])) == &v[1]
	})
	c("Y11 rows via reflect both orders", func() bool {
		var a [2][2]int64
		v := (*[2][2]N)(unsafe.Pointer(&a))
		q := reflect.ValueOf(v).Elem().Index(0).Addr().Interface().(*[2]N)
		return q == &v[0] && q == (*[2]N)(unsafe.Pointer(&a[0]))
	})
	c("Y12 row reflect Set through view", func() bool {
		a := [2][2]int64{{1, 2}, {3, 4}}
		v := (*[2][2]N)(unsafe.Pointer(&a))
		reflect.ValueOf(v).Elem().Index(1).Set(reflect.ValueOf([2]N{5, 6}))
		var i interface{} = a[1][0]
		_, ok := i.(int64)
		return ok && a[1] == [2]int64{5, 6}
	})
	c("Y13 named row type", func() bool {
		var a [2][2]int64
		v := (*[2]R)(unsafe.Pointer(&a))
		return (*[2]int64)(unsafe.Pointer(&v[1])) == &a[1]
	})
	c("Y14 first row and whole array stay apart", func() bool {
		var a [2][2]int64
		v := (*[2][2]N)(unsafe.Pointer(&a))
		return unsafe.Pointer(&v[0]) == unsafe.Pointer(v) && (*[2][2]int64)(unsafe.Pointer(&v[0])) == &a
	})
	c("Y15 slice over view row", func() bool {
		a := [2][2]int64{{1, 2}, {3, 4}}
		v := (*[2][2]N)(unsafe.Pointer(&a))
		s := v[1][:]
		s[0] = 9
		return a[1][0] == 9 && &s[1] == (*N)(unsafe.Pointer(&a[1][1]))
	})
	c("Y16 %T of &v[1]", func() bool {
		var a [2][2]int64
		_ = &a[1]
		v := (*[2][2]N)(unsafe.Pointer(&a))
		return fmt.Sprintf("%T", &v[1]) == "*[2]main.N" && fmt.Sprintf("%T", v[1][0]) == "main.N"
	})
}
