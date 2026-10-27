package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
	arr := [3]int{1, 2, 3}
	s := arr[:]
	fmt.Println("r1:", reflect.ValueOf(s).Slice(3, 3).Pointer() == reflect.ValueOf(s[3:]).Pointer(), unsafe.SliceData(s[3:]) == &arr[0])
	fmt.Println("r2:", reflect.ValueOf(&arr).Elem().Slice(3, 3).Pointer() == uintptr(unsafe.Pointer(&arr)))
	fmt.Println("r3:", reflect.ValueOf(s).Slice3(1, 1, 1).Pointer() == uintptr(unsafe.Pointer(unsafe.SliceData(s[1:1:1]))), unsafe.SliceData(s[1:1:1]) == &arr[0])
	fmt.Println("r4:", unsafe.SliceData(s[1:1]) == &arr[1], unsafe.SliceData(arr[3:]) == &arr[0])
	str := "abc"
	fmt.Println("r5:", str[3:] == "")
}
