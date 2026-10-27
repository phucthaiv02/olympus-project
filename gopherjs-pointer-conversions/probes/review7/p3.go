package main

import (
	"unsafe"

	"github.com/gopherjs/gopherjs/js"
)

func main() {
	defer func() { println("PANIC", recover()) }()
	o := js.Global.Get("Object").New()
	p := unsafe.Pointer(o)
	js.Global.Get("console").Call("log", "p:", js.InternalObject(p))
	q := (*js.Object)(p)
	js.Global.Get("console").Call("log", "q:", js.InternalObject(q))
	var np unsafe.Pointer
	js.Global.Get("console").Call("log", "nil:", js.InternalObject((*js.Object)(np)))
	var n *js.Object
	js.Global.Get("console").Call("log", "nilobj->unsafe:", js.InternalObject(unsafe.Pointer(n)))
}
