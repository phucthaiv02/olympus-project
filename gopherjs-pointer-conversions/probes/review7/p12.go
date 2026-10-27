package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type B struct {
	_ int32
	_ int32
	X int32
}

func main() {
	b := &B{X: 3}
	v := reflect.ValueOf(b).Elem()
	a0 := v.Field(0).UnsafeAddr()
	a1 := v.Field(1).UnsafeAddr()
	fmt.Println("blank:", a0 == a1, a0 == uintptr(unsafe.Pointer(b)), v.Field(0).Addr().Pointer() == v.Field(1).Addr().Pointer())
}
