package main

import (
	"reflect"
	"unsafe"

	"github.com/gopherjs/gopherjs/js"
)

func try(name string, f func()) {
	defer func() {
		if r := recover(); r != nil {
			println(name, "PANIC:", r)
		}
	}()
	f()
}

type S struct {
	O *js.Object
	N int
}

func main() {
	try("nil jsobj -> unsafe", func() {
		var o *js.Object
		p := unsafe.Pointer(o)
		println("nil jsobj unsafe==nil:", p == nil)
	})
	try("nonnil jsobj roundtrip", func() {
		o := js.Global.Get("Object").New()
		p := unsafe.Pointer(o)
		println("roundtrip:", (*js.Object)(p) == o, p != nil)
	})
	try("js array roundtrip", func() {
		o := js.Global.Get("Array").New(1, 2, 3)
		p := unsafe.Pointer(o)
		println("arr roundtrip:", (*js.Object)(p) == o)
		println("arr has $ptr:", o.Get("$ptr") != js.Undefined, o.Get("$proxies") != js.Undefined)
	})
	try("reflect Pointer nil jsobj", func() {
		var o *js.Object
		println("refl ptr nil:", reflect.ValueOf(&o).Elem().Pointer())
	})
	try("reflect Pointer jsobj", func() {
		o := js.Global.Get("Object").New()
		println("refl ptr eq:", reflect.ValueOf(o).Pointer() == uintptr(unsafe.Pointer(o)))
	})
	try("nil unsafe -> jsobj", func() {
		var p unsafe.Pointer
		println("jsobj nil:", (*js.Object)(p) == nil)
	})
	try("struct with jsobj first", func() {
		s := &S{}
		println("&s vs &s.O:", unsafe.Pointer(s) == unsafe.Pointer(&s.O))
		println("back:", (*S)(unsafe.Pointer(&s.O)) == s)
	})
	try("jsobj in map key", func() {
		o := js.Global.Get("Object").New()
		m := map[unsafe.Pointer]int{}
		m[unsafe.Pointer(o)] = 1
		println("map:", m[unsafe.Pointer(o)])
	})
	try("ValueOf nil jsobj Convert unsafe", func() {
		var o *js.Object
		v := reflect.ValueOf(&o).Elem().Convert(reflect.TypeOf(unsafe.Pointer(nil)))
		println("conv nil:", v.Interface().(unsafe.Pointer) == nil)
	})
	try("slice of jsobj via unsafe", func() {
		a := [2]*js.Object{js.Global, nil}
		type P = *js.Object
		b := (*[2]P)(unsafe.Pointer(&a))
		println("arr view:", b[0] == js.Global, b[1] == nil, unsafe.Pointer(b) == unsafe.Pointer(&a[0]))
	})
}
