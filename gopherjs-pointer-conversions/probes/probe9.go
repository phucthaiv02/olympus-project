package main

import (
	"fmt"
	"unsafe"
)

type MyInt int
type MyInt64 int64
type MyC complex128
type A64 [2]int64
type B64 [2]MyInt64
type SA struct {
	X int64
	Y *int
	Z [2]int64
}
type SB struct {
	X MyInt64
	Y *MyInt
	Z [2]MyInt64
}
type Inner struct{ V int64 }
type Inner2 struct{ V MyInt64 }
type OA struct{ In Inner }
type OB struct{ In Inner2 }

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
	check("Y1 arr int64 elem type", func() bool { a := [2]int64{1, 2}; q := (*[2]MyInt64)(unsafe.Pointer(&a)); return fmt.Sprintf("%T", any(q[1])) == "main.MyInt64" && q[1] == 2 })
	check("Y2 arr int64 write", func() bool { a := [2]int64{1, 2}; q := (*[2]MyInt64)(unsafe.Pointer(&a)); q[0] = 5; return a[0] == 5 && fmt.Sprintf("%T", any(a[0])) == "int64" })
	check("Y3 arr int64 identity", func() bool { a := [2]int64{}; q := (*[2]MyInt64)(unsafe.Pointer(&a)); return (*[2]int64)(unsafe.Pointer(q)) == &a && q == (*[2]MyInt64)(unsafe.Pointer(&a)) })
	check("Y4 named arr B64", func() bool { a := A64{1, 2}; q := (*B64)(unsafe.Pointer(&a)); return fmt.Sprintf("%T %T", q, any(q[0])) == "*main.B64 main.MyInt64" })
	check("Y5 arr of ptr", func() bool { i := 1; a := [1]*int{&i}; q := (*[1]*MyInt)(unsafe.Pointer(&a)); return fmt.Sprintf("%T", q[0]) == "*main.MyInt" && q[0] == (*MyInt)(&i) })
	check("Y6 arr complex", func() bool { a := [1]complex128{1i}; q := (*[1]MyC)(unsafe.Pointer(&a)); return fmt.Sprintf("%T", any(q[0])) == "main.MyC" })
	check("Y7 arr int elem (wrapped)", func() bool { a := [2]int{1, 2}; q := (*[2]MyInt)(unsafe.Pointer(&a)); q[1] = 7; return a[1] == 7 && fmt.Sprintf("%T", q[1]) == "main.MyInt" })
	check("Y8 struct field int64", func() bool { s := &SA{X: 3}; q := (*SB)(unsafe.Pointer(s)); return fmt.Sprintf("%T", any(q.X)) == "main.MyInt64" && q.X == 3 })
	check("Y9 struct field write", func() bool { s := &SA{}; q := (*SB)(unsafe.Pointer(s)); q.X = 9; return s.X == 9 && fmt.Sprintf("%T", any(s.X)) == "int64" })
	check("Y10 struct field ptr", func() bool { i := 1; s := &SA{Y: &i}; q := (*SB)(unsafe.Pointer(s)); return fmt.Sprintf("%T", q.Y) == "*main.MyInt" && q.Y == (*MyInt)(&i) })
	check("Y11 struct field arr", func() bool { s := &SA{Z: [2]int64{4, 5}}; q := (*SB)(unsafe.Pointer(s)); return fmt.Sprintf("%T", any(q.Z[1])) == "main.MyInt64" && q.Z[1] == 5 })
	check("Y12 struct identity", func() bool { s := &SA{}; q := (*SB)(unsafe.Pointer(s)); return (*SA)(unsafe.Pointer(q)) == s && q == (*SB)(unsafe.Pointer(s)) })
	check("Y13 nested struct", func() bool { o := &OA{Inner{2}}; q := (*OB)(unsafe.Pointer(o)); return fmt.Sprintf("%T", any(q.In.V)) == "main.MyInt64" && q.In.V == 2 })
	check("Y14 nested struct write", func() bool { o := &OA{}; q := (*OB)(unsafe.Pointer(o)); q.In.V = 6; return o.In.V == 6 })
	check("Y15 struct value copy type", func() bool { s := &SA{X: 1}; q := (*SB)(unsafe.Pointer(s)); v := *q; return fmt.Sprintf("%T", any(v.X)) == "main.MyInt64" })
	check("Y16 struct whole write", func() bool { s := &SA{}; q := (*SB)(unsafe.Pointer(s)); *q = SB{X: 4}; return s.X == 4 && fmt.Sprintf("%T", any(s.X)) == "int64" })
	check("Y17 slice elem int64 via unsafe", func() bool { v := []int64{1, 2}; q := (*MyInt64)(unsafe.Pointer(&v[1])); return fmt.Sprintf("%T", any(*q)) == "main.MyInt64" && q == (*MyInt64)(unsafe.Pointer(&v[1])) })
	check("Y18 *[]int64->*[]MyInt64", func() bool { v := []int64{1}; q := (*[]MyInt64)(unsafe.Pointer(&v)); return fmt.Sprintf("%T", any((*q)[0])) == "main.MyInt64" })
	check("Y19 field addr after conv", func() bool { s := &SA{}; q := (*SB)(unsafe.Pointer(s)); return unsafe.Pointer(&q.X) == unsafe.Pointer(&s.X) })
	check("Y20 elem addr after conv", func() bool { a := [2]int64{}; q := (*[2]MyInt64)(unsafe.Pointer(&a)); return unsafe.Pointer(&q[1]) == unsafe.Pointer(&a[1]) })
}
