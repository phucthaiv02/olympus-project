package main

import (
	"reflect"
	"unsafe"

	"github.com/gopherjs/gopherjs/js"
)

type S struct {
	*js.Object
	N int    `js:"n"`
	T string `js:"t"`
}

func main() {
	s := &S{Object: js.Global.Get("Object").New()}
	s.N = 5
	p := &s.N
	rv := reflect.ValueOf(s).Elem().Field(1)
	println("addr==&s.N:", rv.Addr().Interface().(*int) == p)
	println("unsafeaddr stable:", rv.UnsafeAddr() == reflect.ValueOf(s).Elem().Field(1).UnsafeAddr())
	println("unsafeaddr == unsafe(&s.N):", rv.UnsafeAddr() == uintptr(unsafe.Pointer(p)))
	println("&s.N == &s.N:", &s.N == &s.N)
	*p = 7
	println("s.N:", s.N, rv.Int())
	type MyInt int
	q := (*MyInt)(unsafe.Pointer(p))
	*q = 9
	println("s.N after:", s.N, unsafe.Pointer(q) == unsafe.Pointer(p))
	// FieldByName
	println("fbn==&s.N:", reflect.ValueOf(s).Elem().FieldByName("N").Addr().Interface().(*int) == p)
}
