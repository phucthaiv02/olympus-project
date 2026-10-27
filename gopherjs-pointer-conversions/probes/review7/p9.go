package main

import (
	"fmt"
	"io"
	"unsafe"
)

type Inner struct{ X, Y int }
type MyInner Inner
type Outer struct {
	In Inner
	Z  int
}

func main() {
	o := &Outer{In: Inner{1, 2}, Z: 3}
	fmt.Fprintf(io.Discard, "%p", o)
	v := (*MyInner)(unsafe.Pointer(&o.In))
	v.Y = 20
	fmt.Println(v.X, v.Y, o.In, o.Z)
	w := (*[1]Inner)(unsafe.Pointer(&o.In))
	fmt.Println(w[0].X, w[0].Y)
}
