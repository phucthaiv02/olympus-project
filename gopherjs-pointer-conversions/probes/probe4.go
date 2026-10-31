package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type MyInt int
type P *int
type Q *int
type PS *[]int
type PM *map[string]int
type SA struct{ X int }
type SB SA
type MyInt64 int64

func check(name string, ok bool) {
	if ok {
		fmt.Println("ok  ", name)
	} else {
		fmt.Println("FAIL", name)
	}
}

func conv[T ~*int](p *int) T   { return T(p) }
func back[T ~*int](p T) *int   { return (*int)(p) }
func convAny[T any, U ~*int](p *int) any { return U(p) }

func main() {
	sl := []int{1, 2, 3}
	e := &sl[1]
	check("slice elem round trip", (*int)(unsafe.Pointer(&sl[1])) == e)
	up := unsafe.Pointer(&sl[1])
	check("slice elem via var", (*int)(up) == e && *(*int)(up) == 2)
	np := new(int)
	*np = 4
	check("new round trip", *(*int)(unsafe.Pointer(new(int))) == 0)

	i := 5
	p := &i
	gp := conv[P](p)
	check("generic type", fmt.Sprintf("%T", gp) == "main.P")
	check("generic identity", conv[P](p) == gp && back(gp) == p)
	var nilp *int
	check("generic nil", conv[P](nilp) == nil)
	check("generic any", fmt.Sprintf("%T", convAny[int, Q](p)) == "main.Q")

	rv := reflect.ValueOf(p).Convert(reflect.TypeOf(P(nil)))
	check("reflect convert type", fmt.Sprintf("%T", rv.Interface()) == "main.P")
	check("reflect convert identity", rv.Interface().(P) == P(p))
	check("reflect convert nil", reflect.ValueOf(nilp).Convert(reflect.TypeOf(P(nil))).IsNil())
	mi := reflect.ValueOf(p).Convert(reflect.TypeOf((*MyInt)(nil))).Interface().(*MyInt)
	check("reflect convert elem", mi == (*MyInt)(p))

	a := &SA{1}
	rb := reflect.ValueOf(a).Convert(reflect.TypeOf((*SB)(nil))).Interface().(*SB)
	rb.X = 7
	check("reflect struct ptr aliasing", a.X == 7 && rb == (*SB)(a))
	x := int64(3)
	r64 := reflect.ValueOf(&x).Convert(reflect.TypeOf((*MyInt64)(nil))).Interface().(*MyInt64)
	check("reflect int64 ptr", fmt.Sprintf("%T %v", any(*r64), *r64) == "main.MyInt64 3" && r64 == (*MyInt64)(&x))
	check("reflect slice elem ptr", reflect.ValueOf(&sl[0]).Convert(reflect.TypeOf((*MyInt)(nil))).Interface().(*MyInt) == (*MyInt)(&sl[0]))
	s := []int{1}
	ps := PS(&s)
	(*ps)[0] = 9
	check("ptr to slice", s[0] == 9 && fmt.Sprintf("%T", ps) == "main.PS" && (*[]int)(ps) == &s)
	m := map[string]int{}
	pm := PM(&m)
	(*pm)["a"] = 1
	check("ptr to map", m["a"] == 1 && (*map[string]int)(pm) == &m)
	check("named to named", Q(P(p)) == Q(p) && fmt.Sprintf("%T", Q(P(p))) == "main.Q")
}

var _zero int

func init() {
	check("_zero round trip", (*int)(unsafe.Pointer(&_zero)) == &_zero)
}
