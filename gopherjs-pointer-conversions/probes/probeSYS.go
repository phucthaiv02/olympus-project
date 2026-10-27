package main

import (
	"fmt"
	"unsafe"
)

type S struct {
	A int32
	B uint8
}

type MyInt32 int32

type S2 struct {
	A MyInt32
	B uint8
}

type Raw struct {
	X [5]uint8
}

var _zero uintptr

func main() {
	s := []byte{1, 2, 3, 4}
	p1 := unsafe.Pointer(&s[1])
	p2 := unsafe.Pointer(&s[1])
	fmt.Println("S1 elem eq", p1 == p2)
	t := s[:2]
	fmt.Println("S2 shorter slice eq", unsafe.Pointer(&t[1]) == p1)
	fmt.Println("S3 back", (*byte)(p1) == &s[1])
	i := 1
	fmt.Println("S4 var index", unsafe.Pointer(&s[i]) == p1)
	var i64 int64 = 1
	fmt.Println("S5 int64 index", unsafe.Pointer(&s[i64]) == p1)
	var st S
	q1 := unsafe.Pointer(&st)
	q2 := unsafe.Pointer(&st)
	fmt.Println("S6 struct eq", q1 == q2, (*S)(q1) == &st)
	var ps *S
	fmt.Println("S7 nil", unsafe.Pointer(ps) == nil)
	a := [3]byte{1, 2, 3}
	fmt.Println("S8 array elem", unsafe.Pointer(&a[0]) == unsafe.Pointer(&a[0]))
	fmt.Println("S9 new differ", unsafe.Pointer(new(S)) != unsafe.Pointer(new(S)))
	z1 := unsafe.Pointer(&_zero)
	fmt.Println("S10 zero", z1 == unsafe.Pointer(&_zero))
	st.A = 5
	b := (*S)(unsafe.Pointer(&st))
	b.B = 7
	fmt.Println("S11 roundtrip", st.A, st.B, b == &st)
	v2 := (*S2)(unsafe.Pointer(&st))
	v2.A = 42
	fmt.Println("S15 names-only struct same var", st.A == 42, v2 == (*S2)(unsafe.Pointer(&st)))
	fmt.Println("S16 names-only unsafe eq", unsafe.Pointer(v2) == unsafe.Pointer(&st))
	fmt.Println("S17 back to *S", (*S)(unsafe.Pointer(v2)) == &st)
	fmt.Printf("S18 field type %T\n", v2.A)
	bs := []byte{1, 2, 3, 4, 5}
	r := (*Raw)(unsafe.Pointer(&bs[0]))
	fmt.Println("S19 bytes to struct", r.X)
	w := []int32{10, 20, 30}
	fmt.Println("S12 int32 elem", unsafe.Pointer(&w[2]) == unsafe.Pointer(&w[2]))
	k := 0
	f := func() int { k++; return 1 }
	_ = unsafe.Pointer(&s[f()])
	fmt.Println("S13 index evaluated once", k)
	defer func() { fmt.Println("S14 out of range panics", recover() != nil) }()
	j := 9
	_ = unsafe.Pointer(&s[j])
}
