package main

import (
	"fmt"
	"unsafe"

	"github.com/gopherjs/gopherjs/js"
)

type MyInt64 int64
type S struct{ X int64 }
type S2 struct{ X MyInt64 }
type S3 S

func main() {
	s := S{1}
	p := &s
	o := js.InternalObject(p)
	fmt.Println("constructor is *S:", o.Get("constructor") == js.InternalObject(p).Get("constructor").Get("elem").Get("ptr"), "elem.kind:", o.Get("constructor").Get("elem").Get("kind").Int())
	fmt.Println("isStructPointer:", js.Global.Call("$isStructPointer", o).Bool())
	q := (*S2)(unsafe.Pointer(p))
	q.X = 7
	fmt.Printf("unsafe view: s.X=%v %T, q.X=%T, back==p %v\n", s.X, s.X, q.X, (*S)(unsafe.Pointer(q)) == p)
	r := (*S3)(p)
	r.X = 9
	fmt.Println("explicit (*S3)(p):", s.X, (*S)(r) == p)
	ps := []S{{1}, {2}}
	e := (*S2)(unsafe.Pointer(&ps[1]))
	e.X = 5
	fmt.Println("slice elem view:", ps[1].X, (*S)(unsafe.Pointer(e)) == &ps[1])
	var g S
	h := new(S)
	fmt.Println("var/new:", js.Global.Call("$isStructPointer", js.InternalObject(&g)).Bool(), js.Global.Call("$isStructPointer", js.InternalObject(h)).Bool())
}
