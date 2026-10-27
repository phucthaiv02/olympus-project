package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type MyInt int
type MyInt64 int64
type P *int
type UP unsafe.Pointer
type H struct {
	F P
	U unsafe.Pointer
	Q *MyInt
}

func main() {
	var x int
	px := &x
	u1 := (*MyInt)(unsafe.Pointer(px))
	var a [3]int64
	pa := (*[3]MyInt64)(unsafe.Pointer(&a))
	fmt.Println("1", reflect.ValueOf(px).UnsafePointer() == unsafe.Pointer(&x), reflect.ValueOf(u1).UnsafePointer() == unsafe.Pointer(&x))
	fmt.Println("2", reflect.ValueOf(&a).UnsafePointer() == unsafe.Pointer(pa), reflect.ValueOf(pa).Pointer() == uintptr(unsafe.Pointer(&a)))
	fmt.Println("3", reflect.ValueOf(&a[1]).Pointer() == reflect.ValueOf(&pa[1]).Pointer())
	fmt.Println("4", reflect.ValueOf(&a).Elem().Index(1).UnsafeAddr() == uintptr(unsafe.Pointer(&pa[1])))
	fmt.Println("5", reflect.ValueOf(pa).Elem().Index(2).Addr().Interface().(*MyInt64) == &pa[2])
	var h H
	hv := reflect.ValueOf(&h).Elem()
	hv.Field(0).Set(reflect.ValueOf(px))
	hv.Field(1).Set(reflect.ValueOf(unsafe.Pointer(px)))
	hv.Field(2).Set(reflect.ValueOf(u1))
	fmt.Println("6", h.F == P(px), h.U == unsafe.Pointer(&x), h.Q == u1)
	hv.Field(1).SetPointer(unsafe.Pointer(u1))
	fmt.Println("7", h.U == unsafe.Pointer(px))
	fmt.Println("8", hv.Field(0).Addr().Interface().(*P) == &h.F, hv.Field(1).UnsafeAddr() == uintptr(unsafe.Pointer(&h.U)))
	up := UP(unsafe.Pointer(px))
	cv := reflect.ValueOf(unsafe.Pointer(px)).Convert(reflect.TypeOf(UP(nil))).Interface().(UP)
	fmt.Println("9", cv == up, unsafe.Pointer(cv) == unsafe.Pointer(u1))
	fmt.Println("10", reflect.ValueOf(P(nil)).IsNil(), reflect.ValueOf(P(px)).Pointer() == uintptr(unsafe.Pointer(&x)))
	fmt.Println("11", reflect.ValueOf(px).Equal(reflect.ValueOf(px)), reflect.ValueOf(px).Convert(reflect.TypeOf(P(nil))).Equal(reflect.ValueOf(P(px))))
	// NewAt to array & struct
	na := reflect.NewAt(reflect.TypeOf([3]MyInt64{}), unsafe.Pointer(&a)).Interface().(*[3]MyInt64)
	fmt.Println("12", na == pa, reflect.TypeOf(na[0]))
	na[0] = 5
	fmt.Println("13", a[0], reflect.TypeOf(a[0]))
	// reflect.Value.Elem of named pointer then Addr
	ev := reflect.ValueOf(P(px)).Elem()
	fmt.Println("14", ev.Addr().Interface().(*int) == px, ev.CanAddr())
	// Convert *[3]int64 -> named *[3]int64
	type PA *[3]int64
	ca := reflect.ValueOf(&a).Convert(reflect.TypeOf(PA(nil))).Interface().(PA)
	fmt.Println("15", ca == PA(&a), reflect.ValueOf((*[3]int64)(nil)).Convert(reflect.TypeOf(PA(nil))).IsNil())
	fmt.Println("16", fmt.Sprintf("%p", pa) == fmt.Sprintf("%p", &a), fmt.Sprintf("%p", &pa[1]) == fmt.Sprintf("%p", &a[1]))
}
