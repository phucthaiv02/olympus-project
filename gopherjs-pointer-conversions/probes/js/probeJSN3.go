package main

import (
	"fmt"
	"unsafe"

	"github.com/gopherjs/gopherjs/js"
)

type P *js.Object

func main() {
	o := js.Global.Get("Object").New()
	o2 := js.Global.Get("Object").New()
	a := [2]P{P(o), nil}
	v := (*[2]*js.Object)(unsafe.Pointer(&a))
	fmt.Println("read0", v[0] == o, "read1nil", v[1] == nil)
	v[1] = o2
	fmt.Println("write", a[1] == P(o2), a[1] != nil)
	js.Global.Get("console").Call("log", js.InternalObject(a))
}
