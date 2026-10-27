package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type MyInt32 int32
type MyF float64
type MyByte byte
type T struct {
	A [2]int32
	B float64
}
type T2 struct {
	A [2]MyInt32
	B MyF
}

func try(name string, f func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(name, "PANIC:", r)
		}
	}()
	f()
}

func main() {
	try("typed", func() {
		a := [4]int32{1, 2, 3, 4}
		v := (*[4]MyInt32)(unsafe.Pointer(&a))
		v[2] = 30
		fmt.Println("t1:", a[2], fmt.Sprintf("%T", v[1]), unsafe.Pointer(&v[2]) == unsafe.Pointer(&a[2]), (*MyInt32)(unsafe.Pointer(&a[2])) == &v[2])
		s := a[1:]
		fmt.Println("t2:", unsafe.Pointer(&s[0]) == unsafe.Pointer(&a[1]), unsafe.Pointer(unsafe.SliceData(s)) == unsafe.Pointer(&a[1]))
		ap := (*[3]int32)(s)
		fmt.Println("t3:", unsafe.Pointer(ap) == unsafe.Pointer(&a[1]), &ap[0] == &a[1], (*[3]int32)(unsafe.Pointer(ap)) == ap)
		mv := (*[3]MyInt32)(unsafe.Pointer(ap))
		mv[0] = 20
		fmt.Println("t4:", a[1], unsafe.Pointer(mv) == unsafe.Pointer(&a[1]), &mv[1] == (*MyInt32)(unsafe.Pointer(&a[2])))
		ms := v[1:]
		fmt.Println("t5:", &ms[0] == &v[1], unsafe.Pointer(&ms[1]) == unsafe.Pointer(&a[2]), reflect.ValueOf(ms).Pointer() == uintptr(unsafe.Pointer(&a[1])))
	})
	try("2d", func() {
		m := [2][2]int32{{1, 2}, {3, 4}}
		fmt.Println("d1:", unsafe.Pointer(&m) == unsafe.Pointer(&m[0][0]), unsafe.Pointer(&m[1]) == unsafe.Pointer(&m[1][0]))
		mv := (*[2][2]MyInt32)(unsafe.Pointer(&m))
		mv[1][1] = 40
		mv[0] = [2]MyInt32{10, 20}
		fmt.Println("d2:", m, &mv[1] == (*[2]MyInt32)(unsafe.Pointer(&m[1])), unsafe.Pointer(&mv[1][1]) == unsafe.Pointer(&m[1][1]))
		fmt.Println("d3:", (*[2][2]int32)(unsafe.Pointer(mv)) == &m, (*[2][2]int32)(unsafe.Pointer(&mv[0][0])) == &m)
	})
	try("struct typed", func() {
		t := &T{A: [2]int32{1, 2}, B: 1.5}
		t2 := (*T2)(unsafe.Pointer(t))
		t2.A[1] = 7
		t2.B = 2.5
		fmt.Println("s1:", t.A, t.B, fmt.Sprintf("%T %T", t2.A[0], t2.B), unsafe.Pointer(t2) == unsafe.Pointer(&t.A[0]), unsafe.Pointer(&t2.A) == unsafe.Pointer(t))
		fmt.Println("s2:", (*T)(unsafe.Pointer(&t2.A[0])) == t, &t2.A[1] == (*MyInt32)(unsafe.Pointer(&t.A[1])), (*[2]MyInt32)(unsafe.Pointer(&t.A)) == &t2.A)
	})
	try("bytes", func() {
		b := []byte("hello")
		s := unsafe.String(&b[0], len(b))
		fmt.Println("b1:", s, unsafe.StringData(s) == &b[0])
		mb := (*[5]MyByte)(unsafe.Pointer(&b[0]))
		_ = mb
		bs := *(*[]MyByte)(unsafe.Pointer(&b))
		bs[0] = 'J'
		fmt.Println("b2:", string(b), fmt.Sprintf("%T", bs[0]), &bs[1] == (*MyByte)(unsafe.Pointer(&b[1])))
	})
}
