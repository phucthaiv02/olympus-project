package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type MyInt int
type MyInt64 int64
type MyC complex128
type MyStr string
type P *int
type PP **int
type S []int
type St struct{ X int }
type St2 St
type PSt *St
type Arr [2]int
type PArr *[2]int
type Holder struct{ F P }

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

func gconv[T ~*int](p *int) T  { return T(p) }
func gback[T ~*int](q T) *int  { return (*int)(q) }
func gtype[T any](x T) string  { return fmt.Sprintf("%T", x) }
func retP(p *int) P            { return p }
func argP(q P) P               { return q }
func varP(qs ...P) []P         { return qs }
func gconvS[T ~*St](p *St) T   { return T(p) }
func gconvA[T ~*[2]int](p *[2]int) T { return T(p) }

func main() {
	// D: mixed comparisons
	check("D1 P(p)==p", func() bool { i := 1; p := &i; return P(p) == p })
	check("D2 p==P(p) var", func() bool { i := 1; p := &i; var q P = p; return p == q })
	check("D3 PSt mixed", func() bool { s := &St{1}; var q PSt = s; return q == s })
	check("D4 PArr mixed", func() bool { a := &[2]int{}; var q PArr = a; return q == a })
	check("D5 any(P)==any(P)", func() bool { i := 1; p := &i; return any(P(p)) == any(P(p)) })
	check("D6 any(P)!=any(p)", func() bool { i := 1; p := &i; return any(P(p)) != any(p) })
	check("D7 any(*MyInt)==any(*MyInt)", func() bool { i := 1; p := &i; return any((*MyInt)(p)) == any((*MyInt)(p)) })
	// G: generics
	check("G1 generic T(p) type", func() bool { i := 1; return fmt.Sprintf("%T", gconv[P](&i)) == "main.P" })
	check("G2 generic identity", func() bool { i := 1; p := &i; return gconv[P](p) == P(p) })
	check("G3 generic back", func() bool { i := 1; p := &i; return gback(gconv[P](p)) == p })
	check("G4 generic write", func() bool { i := 1; q := gconv[P](&i); *q = 4; return i == 4 })
	check("G5 generic nil", func() bool { return gconv[P](nil) == nil })
	check("G6 generic gtype", func() bool { i := 1; return gtype(P(&i)) == "main.P" })
	check("G7 generic struct", func() bool { s := &St{1}; q := gconvS[PSt](s); return fmt.Sprintf("%T", q) == "main.PSt" && q == PSt(s) })
	check("G8 generic array", func() bool { a := &[2]int{}; q := gconvA[PArr](a); return fmt.Sprintf("%T", q) == "main.PArr" && q == PArr(a) })
	// I: implicit contexts
	check("I1 return", func() bool { i := 1; p := &i; q := retP(p); return fmt.Sprintf("%T", q) == "main.P" && q == P(p) })
	check("I2 arg", func() bool { i := 1; p := &i; return fmt.Sprintf("%T", argP(p)) == "main.P" })
	check("I3 variadic", func() bool { i := 1; p := &i; v := varP(p, p); return fmt.Sprintf("%T", v[0]) == "main.P" && v[0] == v[1] })
	check("I4 composite field", func() bool { i := 1; p := &i; h := Holder{F: p}; return fmt.Sprintf("%T", h.F) == "main.P" })
	check("I5 slice lit", func() bool { i := 1; p := &i; v := []P{p}; return fmt.Sprintf("%T", v[0]) == "main.P" })
	check("I6 map value", func() bool { i := 1; p := &i; m := map[string]P{}; m["a"] = p; return fmt.Sprintf("%T", m["a"]) == "main.P" })
	check("I7 map key", func() bool { i := 1; p := &i; m := map[P]int{}; m[p] = 1; m[P(p)]++; return len(m) == 1 && m[p] == 2 })
	check("I8 chan send", func() bool { i := 1; p := &i; c := make(chan P, 1); c <- p; return fmt.Sprintf("%T", <-c) == "main.P" })
	check("I9 append", func() bool { i := 1; p := &i; var v []P; v = append(v, p); return fmt.Sprintf("%T", v[0]) == "main.P" })
	check("I10 multi assign", func() bool { i := 1; p := &i; var a, b P; a, b = p, p; return fmt.Sprintf("%T", a) == "main.P" && a == b })
	check("I11 struct field assign", func() bool { i := 1; p := &i; var h Holder; h.F = p; return fmt.Sprintf("%T", h.F) == "main.P" })
	check("I12 array lit", func() bool { i := 1; p := &i; v := [1]P{p}; return fmt.Sprintf("%T", v[0]) == "main.P" })
	check("I13 map lit key", func() bool { i := 1; p := &i; m := map[P]int{p: 1}; return m[P(p)] == 1 })
	check("I14 func value arg", func() bool { i := 1; p := &i; f := argP; return fmt.Sprintf("%T", f(p)) == "main.P" })
	check("I15 implicit back", func() bool { i := 1; p := &i; var q P = p; var r *int = q; return r == p })
	check("I16 PP from **int", func() bool { i := 1; p := &i; var q PP = &p; return fmt.Sprintf("%T", q) == "main.PP" && *q == p })
	check("I17 ptr of P", func() bool { i := 1; p := &i; q := (*P)(&p); return fmt.Sprintf("%T", *q) == "main.P" && *q == P(p) })
	check("I18 write via *P", func() bool { i, j := 1, 2; p := &i; q := (*P)(&p); *q = &j; return p == &j })
	// T: other element kinds
	check("T1 *[]int->*S type", func() bool { s := []int{1}; q := (*S)(&s); return fmt.Sprintf("%T", *q) == "main.S" })
	check("T2 *[]int->*S write", func() bool { s := []int{1}; q := (*S)(&s); *q = S{5, 6}; return len(s) == 2 && s[0] == 5 })
	check("T3 *[]int->*S identity", func() bool { s := []int{1}; return (*S)(&s) == (*S)(&s) && (*[]int)((*S)(&s)) == &s })
	check("T4 *string->*MyStr", func() bool { s := "a"; q := (*MyStr)(&s); *q = "b"; return s == "b" && fmt.Sprintf("%T", *q) == "main.MyStr" })
	check("T5 *int64->*MyInt64 read", func() bool { x := int64(3); q := (*MyInt64)(&x); return fmt.Sprintf("%T", any(*q)) == "main.MyInt64" && *q == 3 })
	check("T6 *int64 write", func() bool { x := int64(3); q := (*MyInt64)(&x); *q = 1 << 40; return x == 1<<40 })
	check("T7 complex", func() bool { c := 1i; q := (*MyC)(&c); *q = 2; return c == 2 && fmt.Sprintf("%T", any(*q)) == "main.MyC" })
	check("T8 field ptr", func() bool { s := St{1}; q := (*MyInt)(&s.X); *q = 9; return s.X == 9 && q == (*MyInt)(&s.X) })
	check("T9 array elem ptr", func() bool { a := [2]int{}; q := (*MyInt)(&a[1]); *q = 9; return a[1] == 9 && (*int)(q) == &a[1] })
	check("T10 *St field via conv", func() bool { s := &St{1}; q := (*St2)(s); q.X = 3; return s.X == 3 })
	check("T11 struct value type", func() bool { s := &St{1}; q := (*St2)(s); return fmt.Sprintf("%T", *q) == "main.St2" })
	check("T12 *Arr->*[2]int", func() bool { a := Arr{}; q := (*[2]int)(&a); q[0] = 3; return a[0] == 3 && (*Arr)(q) == &a })
	check("T13 *St -> *St2 -> *St", func() bool { s := &St{1}; return (*St)((*St2)(s)) == s })
	check("T14 ptr to ptr field", func() bool { var h Holder; i := 1; q := (*(*int))(&h.F); *q = &i; return h.F == P(&i) })
	// U: unsafe
	check("U1 unsafe equal across types", func() bool { i := 1; p := &i; return unsafe.Pointer(p) == unsafe.Pointer(P(p)) })
	check("U2 unsafe equal *MyInt", func() bool { i := 1; p := &i; return unsafe.Pointer(p) == unsafe.Pointer((*MyInt)(p)) })
	check("U3 unsafe map key", func() bool { i := 1; p := &i; m := map[unsafe.Pointer]int{}; m[unsafe.Pointer(p)]++; m[unsafe.Pointer((*MyInt)(p))]++; return len(m) == 1 })
	check("U4 unsafe struct equal", func() bool { s := &St{}; return unsafe.Pointer(s) == unsafe.Pointer((*St2)(s)) })
	check("U5 unsafe from P back", func() bool { i := 1; p := &i; return (*int)(unsafe.Pointer(P(p))) == p })
	check("U6 unsafe any", func() bool { i := 1; p := &i; return any(unsafe.Pointer(p)) == any(unsafe.Pointer((*MyInt)(p))) })
	check("U7 unsafe distinct", func() bool { i, j := 1, 2; return unsafe.Pointer(&i) != unsafe.Pointer(&j) })
	check("U8 unsafe nil eq", func() bool { var p P; return unsafe.Pointer(p) == nil })
	check("U9 unsafe arr equal", func() bool { a := &[2]int{}; return unsafe.Pointer(a) == unsafe.Pointer((*Arr)(a)) })
	check("U10 unsafe elem equal", func() bool { s := []int{1, 2}; return unsafe.Pointer(&s[1]) == unsafe.Pointer((*MyInt)(&s[1])) })
	// R: reflect
	check("R1 Pointer() equal", func() bool { i := 1; p := &i; return reflect.ValueOf(p).Pointer() == reflect.ValueOf(P(p)).Pointer() })
	check("R2 Set assignable", func() bool { i := 1; p := &i; var q P; reflect.ValueOf(&q).Elem().Set(reflect.ValueOf(p)); return fmt.Sprintf("%T", q) == "main.P" && q == P(p) })
	check("R3 Elem type", func() bool { i := 1; v := reflect.ValueOf((*MyInt)(&i)).Elem(); return fmt.Sprintf("%T", v.Interface()) == "main.MyInt" })
	check("R4 Elem SetInt", func() bool { i := 1; reflect.ValueOf((*MyInt)(&i)).Elem().SetInt(6); return i == 6 })
	check("R5 Call arg", func() bool { i := 1; p := &i; out := reflect.ValueOf(argP).Call([]reflect.Value{reflect.ValueOf(p)}); return fmt.Sprintf("%T", out[0].Interface()) == "main.P" })
	check("R6 MapIndex key", func() bool { i := 1; p := &i; m := map[P]int{P(p): 3}; return reflect.ValueOf(m).MapIndex(reflect.ValueOf(P(p))).Int() == 3 })
	check("R7 DeepEqual", func() bool { i := 1; p := &i; return reflect.DeepEqual(P(p), P(p)) })
	check("R8 %p equal", func() bool { i := 1; p := &i; return fmt.Sprintf("%p", p) == fmt.Sprintf("%p", P(p)) })
	check("R12 Set struct ptr", func() bool { s := &St{}; var q PSt; reflect.ValueOf(&q).Elem().Set(reflect.ValueOf(s)); return fmt.Sprintf("%T", q) == "main.PSt" })
	check("R13 Append", func() bool { i := 1; p := &i; v := reflect.Append(reflect.ValueOf([]P{}), reflect.ValueOf(P(p))); return v.Index(0).Interface().(P) == P(p) })
	check("R14 Convert *St->PSt identity", func() bool { s := &St{}; v := reflect.ValueOf(s).Convert(reflect.TypeOf(PSt(nil))).Interface().(PSt); return v == PSt(s) })
	check("R15 Convert *[2]int->PArr", func() bool { a := &[2]int{}; v := reflect.ValueOf(a).Convert(reflect.TypeOf(PArr(nil))).Interface(); return fmt.Sprintf("%T", v) == "main.PArr" && v.(PArr) == PArr(a) })
	// E: evaluation order / once
	check("E1 conv eval once", func() bool { n := 0; i := 1; f := func() *int { n++; return &i }; q := P(f()); _ = q; return n == 1 })
	check("E2 implicit eval once", func() bool { n := 0; i := 1; f := func() *int { n++; return &i }; var q P = f(); _ = q; return n == 1 })
	check("E3 unsafe eval once", func() bool { n := 0; i := 1; f := func() *int { n++; return &i }; q := (*MyInt)(unsafe.Pointer(f())); _ = q; return n == 1 })
	check("E4 reflect eval", func() bool { n := 0; i := 1; f := func() *int { n++; return &i }; _ = (*MyInt)(f()); return n == 1 })
}
