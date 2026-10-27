package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type P *int
type MyInt int

func main() {
	x := 1
	px := &x
	var q *int
	reflect.ValueOf(&q).Elem().Set(reflect.ValueOf(P(px)))
	fmt.Println(1, q == px)
	c := reflect.ValueOf(P(px)).Convert(reflect.TypeOf(px)).Interface().(*int)
	fmt.Println(2, c == px)
	var ifc interface{}
	reflect.ValueOf(&ifc).Elem().Set(reflect.ValueOf(P(px)))
	fmt.Println(3, ifc == interface{}(P(px)))
	m := map[unsafe.Pointer]int{unsafe.Pointer(px): 1}
	fmt.Println(4, m[unsafe.Pointer((*MyInt)(unsafe.Pointer(px)))], m[unsafe.Pointer(P(px))])
	mi := map[interface{}]int{(*MyInt)(unsafe.Pointer(px)): 2}
	fmt.Println(5, mi[(*MyInt)(unsafe.Pointer(&x))])
	var nilP P
	fmt.Println(6, reflect.ValueOf(nilP).Convert(reflect.TypeOf(px)).IsNil(), reflect.ValueOf(nilP).Convert(reflect.TypeOf(px)).Interface().(*int) == nil)
	v := reflect.New(reflect.TypeOf(P(nil))).Elem()
	v.Set(reflect.ValueOf(px))
	fmt.Println(7, v.Interface().(P) == P(px), v.Pointer() == uintptr(unsafe.Pointer(px)))
}
