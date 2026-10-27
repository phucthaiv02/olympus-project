package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type MyInt64 int64
type I1 struct{ V int64 }
type I2 struct{ V MyInt64 }
type S struct {
	In  I1
	Arr [2]int64
	P   *int64
}
type T struct {
	In  I2
	Arr [2]MyInt64
	P   *MyInt64
}

func main() {
	var arr [2]S
	arr2 := (*[2]T)(unsafe.Pointer(&arr))
	fmt.Println("1", &arr2[1] == (*T)(unsafe.Pointer(&arr[1])), unsafe.Pointer(&arr2[1]) == unsafe.Pointer(&arr[1]))
	fmt.Println("2", (*S)(unsafe.Pointer(&arr2[1])) == &arr[1])
	pt := (*T)(unsafe.Pointer(&arr[0]))
	fmt.Println("3", &pt.In == (*I2)(unsafe.Pointer(&arr[0].In)), reflect.TypeOf(pt.In.V), &pt.In.V == (*MyInt64)(unsafe.Pointer(&arr[0].In.V)))
	pt.In = I2{5}
	fmt.Println("4", arr[0].In.V, reflect.TypeOf(arr[0].In.V))
	pt.Arr[1] = 7
	fmt.Println("5", arr[0].Arr, reflect.TypeOf(arr[0].Arr[1]), reflect.TypeOf(pt.Arr[1]))
	fmt.Println("6", &pt.Arr == (*[2]MyInt64)(unsafe.Pointer(&arr[0].Arr)), &pt.Arr[1] == (*MyInt64)(unsafe.Pointer(&arr[0].Arr[1])))
	fmt.Println("7", (*int64)(unsafe.Pointer(&pt.Arr[1])) == &arr[0].Arr[1])
	var z int64 = 3
	arr[0].P = &z
	fmt.Println("8", pt.P == (*MyInt64)(unsafe.Pointer(&z)), reflect.TypeOf(pt.P), reflect.TypeOf(*pt.P))
	var w MyInt64 = 4
	pt.P = &w
	fmt.Println("9", arr[0].P == (*int64)(unsafe.Pointer(&w)), reflect.TypeOf(*arr[0].P), *arr[0].P)
	pt.P = nil
	fmt.Println("10", arr[0].P == nil)
	// slice of structs
	ss := []S{{}, {}}
	ts := unsafe.Slice((*T)(unsafe.Pointer(&ss[0])), 2)
	fmt.Println("11", &ts[1] == (*T)(unsafe.Pointer(&ss[1])), &ts[1].In.V == (*MyInt64)(unsafe.Pointer(&ss[1].In.V)))
	ts[1].In.V = 9
	fmt.Println("12", ss[1].In.V, reflect.TypeOf(ss[1].In.V))
	// copy struct value read through converted pointer
	cp := *pt
	cp.In.V = 100
	fmt.Println("13", arr[0].In.V, reflect.TypeOf(cp.In.V))
	*pt = T{In: I2{11}}
	fmt.Println("14", arr[0].In.V, reflect.TypeOf(arr[0].In.V), arr[0].P == nil)
	// reflect over converted
	rv := reflect.ValueOf(pt).Elem().Field(0).Field(0)
	fmt.Println("15", rv.Type(), rv.Addr().Interface().(*MyInt64) == &pt.In.V)
	rv.SetInt(77)
	fmt.Println("16", arr[0].In.V, reflect.TypeOf(arr[0].In.V))
	fmt.Println("17", reflect.ValueOf(&arr[0]).Pointer() == reflect.ValueOf(pt).Pointer())
	na := reflect.NewAt(reflect.TypeOf(T{}), unsafe.Pointer(&arr[1])).Interface().(*T)
	fmt.Println("18", na == &arr2[1])
}
