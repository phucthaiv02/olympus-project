package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type MyInt int
type MyInt64 int64
type P *int
type PM *MyInt
type S struct {
	A int64
	B int
}
type T struct {
	A MyInt64
	B MyInt
}
type PS *S

func takeP(p P) P { return p }
func retP(p *int) P  { return p }

func main() {
	var x int
	px := &x
	fmt.Println("1 explicit eq", P(px) == P(&x), P((*int)(nil)) == nil)
	var ip P = px
	fmt.Println("2 implicit eq", ip == P(px), takeP(px) == P(&x), retP(px) == ip)
	ch := make(chan P, 1)
	ch <- px
	fmt.Println("3 chan", <-ch == ip)
	m := map[int]P{1: px}
	fmt.Println("4 map", m[1] == ip)
	sl := []P{px}
	sl = append(sl, px)
	fmt.Println("5 slice", sl[0] == ip, sl[1] == ip)
	type H struct{ p P }
	h := H{px}
	fmt.Println("6 complit", h.p == ip)
	var ifc interface{} = ip
	fmt.Println("7 iface", ifc == interface{}(P(&x)))

	u1 := (*MyInt)(unsafe.Pointer(&x))
	u2 := (*MyInt)(unsafe.Pointer(px))
	fmt.Println("8 unsafe", u1 == u2, unsafe.Pointer(u1) == unsafe.Pointer(&x))
	*u1 = 7
	fmt.Println("9 write", x, reflect.TypeOf(*u1))
	var pm PM = u1
	fmt.Println("10 PM", pm == PM(u2), (*int)(unsafe.Pointer(pm)) == &x)

	var a [2]int64 = [2]int64{1, 2}
	pa := (*[2]MyInt64)(unsafe.Pointer(&a))
	fmt.Println("11 arr", reflect.TypeOf(pa[0]), pa == (*[2]MyInt64)(unsafe.Pointer(&a)), unsafe.Pointer(pa) == unsafe.Pointer(&a))
	fmt.Println("12 elem", (*MyInt64)(unsafe.Pointer(&a[1])) == &pa[1], unsafe.Pointer(&pa[0]) == unsafe.Pointer(&a))
	pa[1] = 9
	fmt.Println("13 write", a, reflect.TypeOf(a[1]))

	s := S{1, 2}
	pt := (*T)(unsafe.Pointer(&s))
	fmt.Println("14 struct", reflect.TypeOf(pt.A), &pt.A == (*MyInt64)(unsafe.Pointer(&s.A)), unsafe.Pointer(pt) == unsafe.Pointer(&s))
	fmt.Println("15 struct eq", pt == (*T)(unsafe.Pointer(&s)), PS(&s) == PS(&s), (*S)(unsafe.Pointer(pt)) == &s)
	pt.A = 42
	fmt.Println("16 swrite", s.A, reflect.TypeOf(s.A))

	// pointer to pointer
	pp := &px
	qq := (**MyInt)(unsafe.Pointer(pp))
	fmt.Println("17 pp", *qq == u1, reflect.TypeOf(*qq), unsafe.Pointer(*qq) == unsafe.Pointer(&x))
	var pp2 *P = (*P)(unsafe.Pointer(pp))
	fmt.Println("18 pp2", *pp2 == P(px), (**int)(unsafe.Pointer(pp2)) == pp)

	// reflect
	rv := reflect.ValueOf(&x).Convert(reflect.TypeOf(P(nil))).Interface().(P)
	fmt.Println("19 rconv", rv == ip)
	rn := reflect.NewAt(reflect.TypeOf(MyInt(0)), unsafe.Pointer(&x)).Interface().(*MyInt)
	fmt.Println("20 newat", rn == u1)
	var sp P
	reflect.ValueOf(&sp).Elem().Set(reflect.ValueOf(&x))
	fmt.Println("21 set", sp == ip)
	fmt.Println("22 ptr", reflect.ValueOf(&x).Pointer() == uintptr(unsafe.Pointer(&x)), reflect.ValueOf(u1).Pointer() == uintptr(unsafe.Pointer(&x)))
	fmt.Println("23 uaddr", reflect.ValueOf(&x).Elem().UnsafeAddr() == uintptr(unsafe.Pointer(&x)))
	fmt.Println("24 addr", reflect.ValueOf(&s).Elem().Field(0).Addr().Interface().(*int64) == &s.A)
	fmt.Println("25 addr2", reflect.ValueOf(&a).Elem().Index(1).Addr().Interface().(*int64) == &a[1])
	nilc := reflect.ValueOf((*int)(nil)).Convert(reflect.TypeOf(P(nil))).Interface().(P)
	fmt.Println("26 nilconv", nilc == nil)
	nn := reflect.NewAt(reflect.TypeOf(0), nil)
	fmt.Println("27 newat nil", nn.IsNil(), nn.Interface().(*int) == nil)
	var up unsafe.Pointer
	fmt.Println("28 zero unsafe", (*int)(up) == nil, (*S)(up) == nil, (*[2]int)(up) == nil, P(up) == nil)
	rconvS := reflect.ValueOf(&s).Convert(reflect.TypeOf(PS(nil))).Interface().(PS)
	fmt.Println("29 rconvS", rconvS == PS(&s))
	rconvA := reflect.ValueOf(&a).Convert(reflect.TypeOf((*[2]int64)(nil))).Interface().(*[2]int64)
	fmt.Println("30 rconvA", rconvA == &a)
}
