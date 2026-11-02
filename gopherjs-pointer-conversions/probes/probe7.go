package main

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
	"unsafe"
)

type MyInt int
type MyInt32 int32
type MyInt64 int64
type P *int
type S []int
type St struct{ X int }
type St2 St
type PSt *St

func (m *MyInt) Inc()                   { *m++ }
func (m MyInt) MarshalJSON() ([]byte, error) { return []byte(`"my"`), nil }
func (s *St2) Set(v int)                { s.X = v }

func check(name string, f func() bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("FAIL", name, "(panic)")
		}
	}()
	if f() {
		fmt.Println("ok  ", name)
	} else {
		fmt.Println("FAIL", name)
	}
}

var gi int
var two = func() (*int, *int) { return &gi, &gi }

func retTwo() (P, P)            { return two() }
func takeTwo(a, b P) string     { return fmt.Sprintf("%T%T", a, b) }
func id[T any](x T) T           { return x }
func uconv[E any](p unsafe.Pointer) *E { return (*E)(p) }
func toU[E any](p *E) unsafe.Pointer   { return unsafe.Pointer(p) }
func gp[T ~*int, U ~*int](x T) U       { return U(x) }

func main() {
	// M: multi-value and statement forms
	check("M1 return tuple", func() bool { a, b := retTwo(); return fmt.Sprintf("%T", a) == "main.P" && fmt.Sprintf("%T", b) == "main.P" })
	check("M2 assign tuple", func() bool { var a, b P; a, b = two(); return fmt.Sprintf("%T%T", a, b) == "main.Pmain.P" })
	check("M3 var tuple", func() bool { var a, b P = two(); return fmt.Sprintf("%T%T", a, b) == "main.Pmain.P" })
	check("M4 call tuple", func() bool { return takeTwo(two()) == "main.Pmain.P" })
	check("M5 range assign", func() bool { ps := []*int{&gi}; var q P; for _, q = range ps { }; return fmt.Sprintf("%T", q) == "main.P" })
	check("M6 select recv", func() bool { ch := make(chan *int, 1); ch <- &gi; var q P; select { case q = <-ch: }; return fmt.Sprintf("%T", q) == "main.P" })
	check("M7 comma-ok map", func() bool { m := map[string]*int{"a": &gi}; var q P; var ok bool; q, ok = m["a"]; return ok && fmt.Sprintf("%T", q) == "main.P" })
	check("M8 comma-ok assert", func() bool { var x any = &gi; var q P; var ok bool; q, ok = x.(*int); return ok && fmt.Sprintf("%T", q) == "main.P" })
	check("M9 comma-ok recv", func() bool { ch := make(chan *int, 1); ch <- &gi; var q P; var ok bool; q, ok = <-ch; return ok && fmt.Sprintf("%T", q) == "main.P" })
	check("M10 defer arg", func() bool {
		var got string
		func() { f := func(q P) { got = fmt.Sprintf("%T", q) }; defer f(&gi) }()
		return got == "main.P"
	})
	check("M11 go arg", func() bool {
		ch := make(chan string, 1); f := func(q P) { ch <- fmt.Sprintf("%T", q) }; go f(&gi); return <-ch == "main.P"
	})
	check("M12 generic id[P]", func() bool { return fmt.Sprintf("%T", id[P](&gi)) == "main.P" })
	check("M13 closure return", func() bool { f := func() P { return &gi }; return fmt.Sprintf("%T", f()) == "main.P" && f() == f() })
	check("M14 op-assign struct field ptr", func() bool { type H struct{ F P }; hs := []H{{}}; hs[0].F = &gi; return fmt.Sprintf("%T", hs[0].F) == "main.P" })
	check("M15 switch case", func() bool { q := P(&gi); switch q { case &gi: return true }; return false })
	// N: generics with unsafe
	check("N1 uconv type", func() bool { i := 1; return fmt.Sprintf("%T", uconv[MyInt](unsafe.Pointer(&i))) == "*main.MyInt" })
	check("N2 uconv identity", func() bool { i := 1; return uconv[MyInt](unsafe.Pointer(&i)) == (*MyInt)(&i) })
	check("N3 uconv same", func() bool { i := 1; return uconv[int](unsafe.Pointer(&i)) == &i })
	check("N4 uconv struct", func() bool { s := &St{1}; q := uconv[St2](unsafe.Pointer(s)); q.X = 5; return s.X == 5 && q == (*St2)(s) })
	check("N5 uconv nil", func() bool { return uconv[MyInt](nil) == nil })
	check("N6 toU", func() bool { i := 1; return toU(&i) == unsafe.Pointer(&i) && (*int)(toU(&i)) == &i })
	check("N7 toU nil", func() bool { var p *int; return toU(p) == nil })
	check("N8 gp", func() bool { i := 1; q := gp[*int, P](&i); return fmt.Sprintf("%T", q) == "main.P" && q == P(&i) })
	check("N9 uconv int64", func() bool { x := int64(5); q := uconv[MyInt64](unsafe.Pointer(&x)); return fmt.Sprintf("%T", any(*q)) == "main.MyInt64" })
	check("N10 uconv array", func() bool { a := [2]int{}; q := uconv[[2]int](unsafe.Pointer(&a)); q[0] = 1; return a[0] == 1 && q == &a })
	// W: write-back types and methods
	check("W1 write keeps src type int64", func() bool { x := int64(1); q := (*MyInt64)(&x); *q = 7; return fmt.Sprintf("%T", any(x)) == "int64" })
	check("W2 unsafe write keeps src type", func() bool { x := int64(1); *(*MyInt64)(unsafe.Pointer(&x)) = 7; return fmt.Sprintf("%T", any(x)) == "int64" && x == 7 })
	check("W3 slice write keeps src type", func() bool { s := []int{1}; *(*S)(&s) = S{2}; return fmt.Sprintf("%T", any(s)) == "[]int" })
	check("W4 int64 slice elem", func() bool { s := []int64{1}; *(*MyInt64)(&s[0]) = 3; return fmt.Sprintf("%T", any(s[0])) == "int64" && s[0] == 3 })
	check("W5 method on converted", func() bool { i := 1; (*MyInt)(&i).Inc(); return i == 2 })
	check("W6 method via unsafe", func() bool { i := 1; (*MyInt)(unsafe.Pointer(&i)).Inc(); return i == 2 })
	check("W7 iface method", func() bool { i := 1; var x interface{ Inc() } = (*MyInt)(&i); x.Inc(); return i == 2 })
	check("W8 struct method", func() bool { s := &St{}; (*St2)(s).Set(4); return s.X == 4 })
	check("W9 json", func() bool { i := 1; b, _ := json.Marshal((*MyInt)(&i)); return string(b) == `"my"` })
	check("W10 atomic", func() bool { var v MyInt32; atomic.AddInt32((*int32)(&v), 2); return v == 2 })
	check("W11 struct unsafe write type", func() bool { s := &St{}; *(*St2)(unsafe.Pointer(s)) = St2{9}; return s.X == 9 && fmt.Sprintf("%T", *s) == "main.St" })
	check("W12 **int write", func() bool { i, j := 1, 2; p := &i; qq := (**MyInt)(unsafe.Pointer(&p)); *qq = (*MyInt)(&j); return p == &j && fmt.Sprintf("%T", p) == "*int" })
	check("W13 *P write type", func() bool { i, j := 1, 2; p := &i; q := (*P)(&p); *q = P(&j); return fmt.Sprintf("%T", p) == "*int" && p == &j })
	check("W14 Inc via pointer to slice elem", func() bool { s := []int{1}; (*MyInt)(&s[0]).Inc(); return s[0] == 2 })
	check("W15 uintptr cast", func() bool { i := 1; return uintptr(unsafe.Pointer(&i)) == uintptr(unsafe.Pointer((*MyInt)(&i))) })
}
