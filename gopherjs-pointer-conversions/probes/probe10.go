package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type MyInt64 int64
type MyInt int

type Raw struct {
	n int64
	p *[2]int64
	a [2]*[2]int64
	k int
}
type Named struct {
	n MyInt64
	p *[2]MyInt64
	a [2]*[2]MyInt64
	k MyInt
}

func check(name string, ok bool) {
	if ok {
		fmt.Println("ok", name)
	} else {
		fmt.Println("FAIL", name)
	}
}

func main() {
	// Field addresses through a struct view, view first.
	s := &Raw{n: 1, k: 2}
	view := (*Named)(unsafe.Pointer(s))
	check("F1 &view.n == &s.n (unsafe)", unsafe.Pointer(&view.n) == unsafe.Pointer(&s.n))
	check("F2 &view.k == &s.k typed", (*int)(unsafe.Pointer(&view.k)) == &s.k)
	check("F3 &view.n repeated", &view.n == &view.n)
	*(&view.n) = 7
	check("F4 write via &view.n", s.n == 7)
	check("F5 %T", fmt.Sprintf("%T", &view.n) == "*main.MyInt64")
	// Original first.
	s2 := &Raw{n: 3}
	pn := &s2.n
	v2 := (*Named)(unsafe.Pointer(s2))
	check("F6 orig first", unsafe.Pointer(&v2.n) == unsafe.Pointer(pn))
	check("F7 typed back", (*int64)(unsafe.Pointer(&v2.n)) == pn)
	check("F8 reflect Pointer", reflect.ValueOf(&v2.n).Pointer() == reflect.ValueOf(pn).Pointer())
	check("F9 ptr field addr", unsafe.Pointer(&v2.p) == unsafe.Pointer(&s2.p))

	// Nil pointers to arrays in arrays.
	var a [2]*[2]int64
	a[1] = &[2]int64{5, 6}
	q := (*[2]*[2]MyInt64)(unsafe.Pointer(&a))
	check("N1 q[0] == nil", q[0] == nil)
	check("N2 q[1][1]", q[1][1] == 6 && fmt.Sprintf("%T", q[1]) == "*[2]main.MyInt64")
	var a1 [1]*[2]int64
	q1 := (*[1]*[2]MyInt64)(unsafe.Pointer(&a1))
	check("N3 all-nil q1[0] == nil", q1[0] == nil)
	check("N4 %T", fmt.Sprintf("%T", q1[0]) == "*[2]main.MyInt64")
	q[1] = nil
	check("N5 write nil keeps source nil", a[1] == nil)
	q[0] = &[2]MyInt64{8, 9}
	check("N6 write non-nil", a[0] != nil && a[0][1] == 9)
	q[0] = nil
	check("N7 write nil over non-nil", a[0] == nil)
	// struct fields
	s3 := &Raw{}
	v3 := (*Named)(unsafe.Pointer(s3))
	check("N8 field nil", v3.p == nil && v3.a[0] == nil)
	s3.p = &[2]int64{1, 2}
	check("N9 field non-nil", v3.p != nil && v3.p[1] == 2)
	v3.p = nil
	check("N10 field write nil", s3.p == nil)
	v3.a[1] = &[2]MyInt64{3, 4}
	check("N11 nested field array write", s3.a[1] != nil && s3.a[1][0] == 3)
	v3.a[1] = nil
	check("N12 nested field array write nil", s3.a[1] == nil)
	// element pointers in views
	check("N13 &q[0] identity", unsafe.Pointer(&q[0]) == unsafe.Pointer(&a[0]))
	check("N14 &q[1] identity", unsafe.Pointer(&q[1]) == unsafe.Pointer(&a[1]))
	a[1] = &[2]int64{1, 1}
	check("N15 &q[1] identity non-nil", unsafe.Pointer(&q[1]) == unsafe.Pointer(&a[1]))

	// Slices read through a converted pointer.
	sv := []int64{1, 2, 3}
	sq := (*[]MyInt64)(unsafe.Pointer(&sv))
	check("S1 elem type", fmt.Sprintf("%T", any((*sq)[1])) == "main.MyInt64" && (*sq)[1] == 2)
	(*sq)[2] = 9
	check("S2 write elem", sv[2] == 9 && fmt.Sprintf("%T", any(sv[2])) == "int64")
	check("S3 elem addr", unsafe.Pointer(&(*sq)[0]) == unsafe.Pointer(&sv[0]))
	t := append((*sq)[:1], 5)
	check("S4 append shares", sv[1] == 5 && fmt.Sprintf("%T", any(t[1])) == "main.MyInt64")
	t2 := append(*sq, 6)
	check("S5 append grow", len(t2) == 4 && t2[3] == 6 && fmt.Sprintf("%T", any(t2[0])) == "main.MyInt64")
	*sq = []MyInt64{7}
	check("S6 whole write", len(sv) == 1 && sv[0] == 7 && fmt.Sprintf("%T", any(sv[0])) == "int64")
	var ns []int64
	nq := (*[]MyInt64)(unsafe.Pointer(&ns))
	check("S7 nil slice", *nq == nil)
	cp := make([]MyInt64, 1)
	sv = []int64{4}
	copy(cp, *sq)
	check("S8 copy", cp[0] == 4)
	type RS struct{ s []int64 }
	type NS struct{ s []MyInt64 }
	rs := &RS{s: []int64{3}}
	ns2 := (*NS)(unsafe.Pointer(rs))
	check("S9 field slice", fmt.Sprintf("%T", any(ns2.s[0])) == "main.MyInt64" && ns2.s[0] == 3)
	ns2.s[0] = 8
	check("S10 field slice write", rs.s[0] == 8)
}
