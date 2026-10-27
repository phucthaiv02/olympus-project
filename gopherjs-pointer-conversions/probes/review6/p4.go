package main

import (
	"fmt"
	"unsafe"
)

type P *int
type MyInt int

func cv[P1, P2 ~*int](p P1) P2           { return P2(p) }
func uc[PT ~*int](p unsafe.Pointer) PT   { return PT(p) }
func tu[PT ~*int](p PT) unsafe.Pointer   { return unsafe.Pointer(p) }
func h() (*int, error)                  { return gx, nil }
func g() (P, error)                     { return h() }
func two() (*int, *int)                 { return gx, nil }
func vf(ps ...P) bool                   { return ps[0] == P(gx) && ps[1] == nil }

var gx = new(int)

func main() {
	ip := P(gx)
	fmt.Println("1 generic", cv[*int, P](gx) == ip, cv[P, *int](nil) == nil, cv[*int, P](nil) == nil)
	fmt.Println("2 generic unsafe", uc[P](unsafe.Pointer(gx)) == ip, uc[P](nil) == nil, tu[P](ip) == unsafe.Pointer(gx), tu[P](nil) == nil)
	r, _ := g()
	fmt.Println("3 tuple return", r == ip)
	var a, b P = two()
	fmt.Println("4 tuple assign", a == ip, b == nil)
	var q P
	for _, q = range []*int{gx} {
	}
	fmt.Println("5 range", q == ip)
	m := map[int]*int{1: gx}
	var ok bool
	var q2 P
	q2, ok = m[1]
	var q3 P = P(nil)
	q3, ok = m[2]
	fmt.Println("6 commaok", q2 == ip, ok, q3 == nil)
	ch := make(chan *int, 1)
	ch <- gx
	var q4 P
	select {
	case q4 = <-ch:
	}
	fmt.Println("7 select", q4 == ip)
	fmt.Println("8 variadic", vf(gx, nil))
	arr := [2]P{gx}
	fmt.Println("9 arrlit", arr[0] == ip, arr[1] == nil)
	mk := map[P]int{gx: 1}
	fmt.Println("10 mapkey lit", mk[ip])
	var nilp *int
	var q5 P = nilp
	fmt.Println("11 nil", q5 == nil)
	var q6 interface{} = P(nilp)
	fmt.Println("12 nil iface", q6.(P) == nil)
	var i2 interface{ M() }
	_ = i2
	ch2 := make(chan *int, 1)
	ch2 <- gx
	var q7 P
	var ok2 bool
	q7, ok2 = <-ch2
	fmt.Println("13 recv commaok", q7 == ip, ok2)
	var it interface{} = gx
	var q8 P
	var pp *int
	pp, _ = it.(*int)
	q8 = pp
	fmt.Println("14", q8 == ip)
	var a1 [2]int
	var a2 [0]int
	fmt.Println("15 arrays", unsafe.Pointer(&a1[0]) == unsafe.Pointer(&a1), unsafe.Pointer(&a2) != nil)
}
