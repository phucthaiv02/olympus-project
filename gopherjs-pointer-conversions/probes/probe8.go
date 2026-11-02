package main

import (
	"fmt"
	"reflect"
	"sort"
	"sync/atomic"
	"unsafe"
)

type MyInt int
type MyByte byte
type P *int
type St struct{ X, Y int }
type St2 St
type PSt *St

func check(name string, f func() bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("FAIL", name, "(panic)", r)
		}
	}()
	if f() {
		fmt.Println("ok  ", name)
	} else {
		fmt.Println("FAIL", name)
	}
}

func main() {
	check("X1 NewAt type", func() bool { i := 1; v := reflect.NewAt(reflect.TypeOf(MyInt(0)), unsafe.Pointer(&i)).Interface(); return fmt.Sprintf("%T", v) == "*main.MyInt" })
	check("X2 NewAt identity", func() bool { i := 1; v := reflect.NewAt(reflect.TypeOf(MyInt(0)), unsafe.Pointer(&i)).Interface().(*MyInt); return v == (*MyInt)(&i) })
	check("X3 NewAt write", func() bool { i := 1; reflect.NewAt(reflect.TypeOf(MyInt(0)), unsafe.Pointer(&i)).Elem().SetInt(4); return i == 4 })
	check("X4 NewAt same type", func() bool { i := 1; return reflect.NewAt(reflect.TypeOf(0), unsafe.Pointer(&i)).Interface().(*int) == &i })
	check("X5 NewAt struct", func() bool { s := &St{}; v := reflect.NewAt(reflect.TypeOf(St2{}), unsafe.Pointer(s)).Interface().(*St2); v.X = 2; return s.X == 2 && v == (*St2)(s) })
	check("X6 UnsafePointer", func() bool { i := 1; return reflect.ValueOf((*MyInt)(&i)).UnsafePointer() == unsafe.Pointer(&i) })
	check("X7 UnsafePointer same", func() bool { i := 1; return (*int)(reflect.ValueOf(&i).UnsafePointer()) == &i })
	check("X8 SetPointer", func() bool { i := 1; var u unsafe.Pointer; reflect.ValueOf(&u).Elem().SetPointer(unsafe.Pointer(&i)); return (*int)(u) == &i })
	check("X9 Addr identity", func() bool { i := 1; q := (*MyInt)(&i); return reflect.ValueOf(q).Elem().Addr().Interface().(*MyInt) == q })
	check("X10 Elem of Convert Addr", func() bool { i := 1; v := reflect.ValueOf(&i).Convert(reflect.TypeOf((*MyInt)(nil))); return v.Elem().Addr().Interface().(*MyInt) == (*MyInt)(&i) })
	check("X11 atomic.Pointer", func() bool { i := 1; var ap atomic.Pointer[MyInt]; ap.Store((*MyInt)(&i)); return ap.Load() == (*MyInt)(&i) && ap.CompareAndSwap((*MyInt)(&i), nil) })
	check("X12 atomic StorePointer", func() bool { i := 1; var u unsafe.Pointer; atomic.StorePointer(&u, unsafe.Pointer((*MyInt)(&i))); return (*int)(atomic.LoadPointer(&u)) == &i })
	check("X13 byte via unsafe", func() bool { b := []byte{1, 2}; q := (*MyByte)(unsafe.Pointer(&b[1])); *q = 9; return b[1] == 9 && fmt.Sprintf("%T", q) == "*main.MyByte" })
	check("X14 byte unsafe identity", func() bool { b := []byte{1, 2}; return (*byte)(unsafe.Pointer(&b[1])) == &b[1] && unsafe.Pointer(&b[1]) == unsafe.Pointer(&b[1]) })
	check("X15 new unsafe eq", func() bool { p := new(int); return unsafe.Pointer(p) == unsafe.Pointer((*MyInt)(p)) })
	check("X16 Swapper []P", func() bool { i, j := 1, 2; v := []P{&i, &j}; reflect.Swapper(v)(0, 1); return v[0] == P(&j) && fmt.Sprintf("%T", v[0]) == "main.P" })
	check("X17 sort.Slice ptrs", func() bool { a, b := 2, 1; v := []P{&a, &b}; sort.Slice(v, func(i, j int) bool { return *v[i] < *v[j] }); return v[0] == P(&b) })
	check("X18 reflect.Copy", func() bool { i := 1; src := []P{&i}; dst := make([]P, 1); reflect.Copy(reflect.ValueOf(dst), reflect.ValueOf(src)); return dst[0] == P(&i) })
	check("X19 map of unsafe from struct", func() bool { s := &St{}; m := map[unsafe.Pointer]bool{unsafe.Pointer(s): true}; return m[unsafe.Pointer((*St2)(s))] })
	check("X20 SetMapIndex conv", func() bool { i := 1; m := map[string]P{}; reflect.ValueOf(m).SetMapIndex(reflect.ValueOf("a"), reflect.ValueOf(&i).Convert(reflect.TypeOf(P(nil)))); return m["a"] == P(&i) })
	check("X21 Value.Set from Convert", func() bool { i := 1; var q *MyInt; reflect.ValueOf(&q).Elem().Set(reflect.ValueOf(&i).Convert(reflect.TypeOf(q))); return q == (*MyInt)(&i) })
	check("X22 struct field reflect conv", func() bool { s := &St{}; v := reflect.ValueOf(&s.Y).Convert(reflect.TypeOf((*MyInt)(nil))).Interface().(*MyInt); *v = 3; return s.Y == 3 && v == (*MyInt)(&s.Y) })
	check("X23 Interface eq after Convert", func() bool { i := 1; a := reflect.ValueOf(&i).Convert(reflect.TypeOf(P(nil))).Interface(); return a == any(P(&i)) })
	check("X24 unsafe of field eq", func() bool { s := St{}; return unsafe.Pointer(&s.Y) == unsafe.Pointer((*MyInt)(&s.Y)) })
	check("X25 ptr to unsafe ptr", func() bool { i := 1; u := unsafe.Pointer(&i); q := (*unsafe.Pointer)(unsafe.Pointer(&u)); return (*int)(*q) == &i })
}
