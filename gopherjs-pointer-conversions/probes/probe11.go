package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type MyInt int
type MyInt64 int64
type Raw struct {
	n int64
	k int
}
type Named struct {
	n MyInt64
	k MyInt
}
type P *int

var upType = reflect.TypeOf(unsafe.Pointer(nil))

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
	check("C1 Convert proxy to unsafe", func() bool {
		n := 1
		q := (*MyInt)(&n)
		u := reflect.ValueOf(q).Convert(upType).Interface().(unsafe.Pointer)
		return u == unsafe.Pointer(&n)
	})
	check("C2 Convert ptr to unsafe", func() bool {
		n := 1
		u := reflect.ValueOf(&n).Convert(upType).Interface().(unsafe.Pointer)
		return u == unsafe.Pointer(&n) && u == unsafe.Pointer((*MyInt)(&n))
	})
	check("C3 Convert named P to unsafe", func() bool {
		n := 1
		u := reflect.ValueOf(P(&n)).Convert(upType).Interface().(unsafe.Pointer)
		return u == unsafe.Pointer(&n)
	})
	check("C4 Convert nil to unsafe", func() bool {
		var p *MyInt
		u := reflect.ValueOf(p).Convert(upType).Interface().(unsafe.Pointer)
		return u == nil
	})
	check("C5 Convert unsafe to ptr round trip", func() bool {
		n := 1
		u := reflect.ValueOf(unsafe.Pointer(&n))
		q := u.Convert(reflect.TypeOf((*MyInt)(nil))).Interface().(*MyInt)
		return q == (*MyInt)(&n) && u.Convert(reflect.TypeOf((*int)(nil))).Interface().(*int) == &n
	})
	check("C6 Convert zero unsafe to ptr", func() bool {
		u := reflect.ValueOf(unsafe.Pointer(nil))
		return u.Convert(reflect.TypeOf((*MyInt)(nil))).Interface().(*MyInt) == nil
	})
	check("C7 Convert struct proxy to unsafe", func() bool {
		s := &Raw{}
		v := (*Named)(unsafe.Pointer(s))
		return reflect.ValueOf(v).Convert(upType).Interface().(unsafe.Pointer) == unsafe.Pointer(s)
	})
	check("C8 Convert array ptr to unsafe", func() bool {
		a := [2]int64{}
		q := (*[2]MyInt64)(unsafe.Pointer(&a))
		return reflect.ValueOf(q).Convert(upType).Interface().(unsafe.Pointer) == unsafe.Pointer(&a)
	})
	check("C9 Convert slice elem ptr to unsafe", func() bool {
		s := []int{1, 2}
		return reflect.ValueOf((*MyInt)(&s[1])).Convert(upType).Interface().(unsafe.Pointer) == unsafe.Pointer(&s[1])
	})
	check("A1 Pointer of proxies", func() bool {
		n := 1
		return reflect.ValueOf((*MyInt)(&n)).Pointer() == reflect.ValueOf(&n).Pointer()
	})
	check("A2 UnsafePointer proxies", func() bool {
		n := 1
		return reflect.ValueOf((*MyInt)(&n)).UnsafePointer() == unsafe.Pointer(&n)
	})
	check("A3 Elem Addr of proxy", func() bool {
		n := 1
		q := (*MyInt)(&n)
		return reflect.ValueOf(q).Elem().Addr().Interface().(*MyInt) == q && reflect.ValueOf(q).Elem().Addr().UnsafePointer() == unsafe.Pointer(&n)
	})
	check("A4 UnsafeAddr of proxy elem", func() bool {
		n := 1
		q := (*MyInt)(&n)
		return reflect.ValueOf(q).Elem().UnsafeAddr() == reflect.ValueOf(&n).Elem().UnsafeAddr()
	})
	check("A5 struct view field Addr", func() bool {
		s := &Raw{}
		v := (*Named)(unsafe.Pointer(s))
		return reflect.ValueOf(v).Elem().Field(0).Addr().UnsafePointer() == unsafe.Pointer(&s.n)
	})
	check("A6 plain struct field Addr (baseline)", func() bool {
		s := &Raw{}
		return reflect.ValueOf(s).Elem().Field(0).Addr().UnsafePointer() == unsafe.Pointer(&s.n)
	})
	check("A7 struct view Pointer", func() bool {
		s := &Raw{}
		v := (*Named)(unsafe.Pointer(s))
		return reflect.ValueOf(v).Pointer() == reflect.ValueOf(s).Pointer()
	})
	check("A8 array view Pointer", func() bool {
		a := [2]int64{}
		q := (*[2]MyInt64)(unsafe.Pointer(&a))
		return reflect.ValueOf(q).Pointer() == reflect.ValueOf(&a).Pointer()
	})
	check("A9 array view elem Addr", func() bool {
		a := [2]int64{}
		q := (*[2]MyInt64)(unsafe.Pointer(&a))
		return reflect.ValueOf(q).Elem().Index(1).Addr().UnsafePointer() == unsafe.Pointer(&a[1])
	})
	check("A10 slice elem Addr", func() bool {
		s := []int{1, 2}
		return reflect.ValueOf(s).Index(1).Addr().UnsafePointer() == unsafe.Pointer(&s[1])
	})
	check("A11 view field read via reflect", func() bool {
		s := &Raw{n: 3}
		v := (*Named)(unsafe.Pointer(s))
		f := reflect.ValueOf(v).Elem().Field(0)
		return f.Type() == reflect.TypeOf(MyInt64(0)) && f.Int() == 3
	})
	check("A12 Value.Set unsafe from proxy", func() bool {
		n := 1
		var u unsafe.Pointer
		reflect.ValueOf(&u).Elem().Set(reflect.ValueOf(unsafe.Pointer((*MyInt)(&n))))
		return u == unsafe.Pointer(&n)
	})
	check("A13 SetPointer then Convert", func() bool {
		n := 1
		var u unsafe.Pointer
		reflect.ValueOf(&u).Elem().SetPointer(unsafe.Pointer((*MyInt)(&n)))
		return (*int)(u) == &n
	})
}
