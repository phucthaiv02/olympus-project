package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type MyInt int
type P *int
type PA *[2]int
type E struct{ F int }
type O struct{ E }
type O2 struct{ *E }
type S1 struct{ p *int }
type S2 struct{ p P }

func conv[T any, PT ~*T](p *T) PT { return PT(p) }
func back[T any, PT ~*T](p PT) *T { return (*T)(p) }
func viaUnsafe[T, U any](p *T) *U  { return (*U)(unsafe.Pointer(p)) }

func f(p P) P { return p }

func main() {
	var x int
	px := &x
	var ip P = px
	fmt.Println("1 &*p", &*ip == px, (*int)(ip) == px)
	var a [2]int
	fmt.Println("2 PA", PA(&a) == PA(&a), PA((*[2]int)(nil)) == nil, (*[2]int)(PA(nil)) == nil)
	var npa PA
	fmt.Println("3 PA nil", npa == nil, (*[2]int)(npa) == nil)
	var o O
	fmt.Println("4 promoted", &o.F == &o.E.F)
	o2 := O2{&E{}}
	fmt.Println("5 promoted ptr", &o2.F == &o2.E.F)
	fmt.Println("6 generic", conv[int, P](px) == ip, back[int, P](ip) == px, conv[int, P](nil) == nil)
	fmt.Println("7 generic unsafe", viaUnsafe[int, MyInt](px) == (*MyInt)(unsafe.Pointer(px)), viaUnsafe[int, MyInt](nil) == nil)
	s1 := S1{px}
	s2 := (*S2)(unsafe.Pointer(&s1))
	fmt.Println("8 struct ptr field", s2.p == ip, reflect.TypeOf(s2.p))
	ap := [2]*int{px, nil}
	ap2 := (*[2]P)(unsafe.Pointer(&ap))
	fmt.Println("9 arr ptr", ap2[0] == ip, ap2[1] == nil, reflect.TypeOf(ap2[0]))
	fmt.Println("10 deref", *(*P)(unsafe.Pointer(&px)) == ip)
	// reflect
	ra := reflect.Append(reflect.ValueOf([]P{}), reflect.ValueOf(px)).Interface().([]P)
	fmt.Println("11 append", ra[0] == ip)
	res := reflect.ValueOf(f).Call([]reflect.Value{reflect.ValueOf(px)})[0].Interface().(P)
	fmt.Println("12 call", res == ip)
	mf := reflect.MakeFunc(reflect.TypeOf(f), func(args []reflect.Value) []reflect.Value {
		return []reflect.Value{reflect.ValueOf(px)}
	}).Interface().(func(P) P)
	fmt.Println("13 makefunc", mf(nil) == ip)
	m := map[P]int{}
	reflect.ValueOf(m).SetMapIndex(reflect.ValueOf(px), reflect.ValueOf(1))
	fmt.Println("14 mapkey", m[ip], m[P(px)])
	ch := make(chan P, 1)
	reflect.ValueOf(ch).Send(reflect.ValueOf(px))
	fmt.Println("15 send", <-ch == ip)
	sl := make([]P, 1)
	reflect.ValueOf(sl).Index(0).Set(reflect.ValueOf(px))
	fmt.Println("16 index set", sl[0] == ip)
	fmt.Println("17 zero", reflect.Zero(reflect.TypeOf(P(nil))).Interface().(P) == nil)
	var st S2
	reflect.ValueOf(&st).Elem().Field(0)
	sv := reflect.ValueOf(&st).Elem()
	_ = sv
	nilPA := reflect.ValueOf((*[2]int)(nil)).Convert(reflect.TypeOf(PA(nil))).Interface().(PA)
	fmt.Println("18 nil PA reflect", nilPA == nil)
	var iface interface{} = ip
	fmt.Println("19 typeassert", iface.(P) == P(px))
	type PP *P
	var pp PP = &ip
	fmt.Println("20 PP", (*P)(pp) == &ip, **pp == 0)
	mp := map[P]bool{P(px): true}
	fmt.Println("21 mapkey2", mp[ip], mp[P(&x)])
	// unsafe.Pointer to named pointer and back
	up := unsafe.Pointer(px)
	fmt.Println("22", P(up) == ip, PA(unsafe.Pointer(&a)) == PA(&a), (*MyInt)(up) == (*MyInt)(unsafe.Pointer(ip)))
	// pointer to slice element
	s := []int{1, 2, 3}
	ms := unsafe.Slice((*MyInt)(unsafe.Pointer(&s[0])), 3)
	fmt.Println("23 slice", &ms[1] == (*MyInt)(unsafe.Pointer(&s[1])), reflect.TypeOf(ms[1]))
	fmt.Println("24 slice2", unsafe.Pointer(&ms[2]) == unsafe.Pointer(&s[2]), (*int)(unsafe.Pointer(&ms[2])) == &s[2])
}
