package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type MyInt int
type P *int
type E struct{ F int64 }
type O struct {
	X int
	E
}
type MyInt64 int64

func up[T any](p *T) unsafe.Pointer { return unsafe.Pointer(p) }
func conv[T, U any](p *T) *U    { return (*U)(unsafe.Pointer(p)) }

func try(name string, f func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(name, "PANIC", r)
		}
	}()
	f()
}

func main() {
	try("G1", func() {
		x := 3
		fmt.Println("G1", up(&x) == unsafe.Pointer(&x), up((*MyInt)(&x)) == unsafe.Pointer(&x))
	})
	try("G2", func() {
		a := [2]MyInt64{1, 2}
		b := conv[[2]MyInt64, [2]int64](&a)
		var e interface{} = b[0]
		fmt.Printf("G2 %T %v\n", e, unsafe.Pointer(b) == unsafe.Pointer(&a))
	})
	try("R1", func() {
		s := []int{1, 2, 3}
		v := reflect.ValueOf(s).Index(1)
		fmt.Println("R1", v.Addr().Interface().(*int) == &s[1], v.UnsafeAddr() == uintptr(unsafe.Pointer(&s[1])))
	})
	try("R2", func() {
		var o O
		v := reflect.ValueOf(&o).Elem().FieldByName("F")
		fmt.Println("R2", v.Addr().Interface().(*int64) == &o.F, &o.F == &o.E.F)
	})
	try("R3", func() {
		x := 1
		var p P = &x
		v := reflect.ValueOf(p)
		fmt.Println("R3", v.Elem().Addr().Interface().(*int) == &x, v.Pointer() == uintptr(unsafe.Pointer(&x)))
	})
	try("R4", func() {
		x := 1
		v := reflect.NewAt(reflect.TypeOf(MyInt(0)), unsafe.Pointer(&x))
		fmt.Println("R4", v.Interface().(*MyInt) == (*MyInt)(unsafe.Pointer(&x)), unsafe.Pointer(v.Pointer()) == unsafe.Pointer(&x))
	})
	try("R5", func() {
		x := 1
		v := reflect.ValueOf(&x).Convert(reflect.TypeOf((*MyInt)(nil)))
		fmt.Println("R5", v.Interface().(*MyInt) == (*MyInt)(&x))
		var np *int
		w := reflect.ValueOf(np).Convert(reflect.TypeOf(P(nil)))
		fmt.Println("R5b", w.IsNil(), w.Interface().(P) == nil)
	})
	try("R6", func() {
		x := 1
		var p P
		reflect.ValueOf(&p).Elem().Set(reflect.ValueOf(&x))
		fmt.Println("R6", p == P(&x), (*int)(p) == &x)
	})
	try("R7", func() {
		arr := [3]int64{1, 2, 3}
		v := reflect.ValueOf(&arr).Elem().Index(1)
		q := (*MyInt64)(unsafe.Pointer(&arr[1]))
		fmt.Println("R7", v.Addr().Interface().(*int64) == &arr[1], unsafe.Pointer(q) == unsafe.Pointer(v.UnsafeAddr()))
		var e interface{} = *q
		fmt.Printf("R7b %T\n", e)
	})
	try("N1", func() {
		x := 1
		p := P(&x)
		fmt.Println("N1", &*p == &x, (*int)(p) == &*p)
	})
}
