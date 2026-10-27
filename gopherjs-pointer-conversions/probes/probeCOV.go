package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type MyInt int
type P *int
type PM *MyInt
type PP **int
type S struct{ X int }
type S2 S
type PS *S
type PS2 *S2
type A [3]int
type PA *[3]int
type PNA *A
type Z struct{}
type PZ *Z
type A0 [0]int
type PA0 *[0]int
type H struct{ f P }
type HP struct{ p *int }

var g int

func check(name string, f func() bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("FAIL", name, "panic:", r)
		}
	}()
	if f() {
		fmt.Println("ok", name)
	} else {
		fmt.Println("FAIL", name)
	}
}

func tn(v any) string { return fmt.Sprintf("%T", v) }

func conv[T ~*int](p *int) T     { return T(p) }
func back[T ~*int](p T) *int      { return (*int)(p) }
func takeP(p P) P                 { return p }
func takePS(p PS) PS              { return p }
func retP() P                     { return &g }
func two() (*int, P)              { return &g, &g }
func (h *H) set(p P)              { h.f = p }
func mk() *int                    { return &g }

func main() {
	// explicit
	check("E1 PS(&local) type+id", func() bool {
		s := &S{1}
		ps := PS(s)
		return tn(ps) == "main.PS" && (*S)(ps) == s && PS(s) == ps
	})
	check("E2 P(&global)", func() bool { return tn(P(&g)) == "main.P" && P(&g) == P(&g) && (*int)(P(&g)) == &g })
	check("E3 P(&s.f) field", func() bool {
		s := HP{}
		var x S
		p := P(&x.X)
		*p = 4
		_ = s
		return x.X == 4 && P(&x.X) == p && (*int)(p) == &x.X && tn(p) == "main.P"
	})
	check("E4 P(&a[1]) array elem", func() bool {
		a := [3]int{1, 2, 3}
		p := P(&a[1])
		*p = 7
		return a[1] == 7 && P(&a[1]) == p && (*int)(p) == &a[1] && tn(p) == "main.P"
	})
	check("E5 P(&sl[1]) slice elem", func() bool {
		sl := []int{1, 2, 3}
		p := P(&sl[1])
		*p = 7
		return sl[1] == 7 && P(&sl[1]) == p && (*int)(p) == &sl[1] && tn(p) == "main.P"
	})
	check("E6 P(new(int))", func() bool {
		n := new(int)
		p := P(n)
		*p = 3
		return *n == 3 && (*int)(p) == n && tn(p) == "main.P"
	})
	check("E7 closure var", func() bool {
		x := 1
		var q P
		func() { q = P(&x) }()
		*q = 5
		return x == 5 && q == P(&x) && (*int)(q) == &x
	})
	check("E8 (*MyInt)(&a[1])", func() bool {
		a := [3]int{1, 2, 3}
		q := (*MyInt)(&a[1])
		*q = 9
		return a[1] == 9 && q == (*MyInt)(&a[1]) && (*int)(q) == &a[1] && tn(q) == "*main.MyInt" && tn(*q) == "main.MyInt"
	})
	check("E9 (*MyInt)(&s.X)", func() bool {
		var s S
		q := (*MyInt)(&s.X)
		*q = 2
		return s.X == 2 && q == (*MyInt)(&s.X) && (*int)(q) == &s.X && tn(*q) == "main.MyInt"
	})
	check("E10 PP(&p) local", func() bool {
		n := 1
		p := &n
		pp := PP(&p)
		return tn(pp) == "main.PP" && (**int)(pp) == &p && PP(&p) == pp && **pp == 1
	})
	check("E11 PA(&a2[1])", func() bool {
		var a2 [2][3]int
		pa := PA(&a2[1])
		pa[2] = 5
		return a2[1][2] == 5 && PA(&a2[1]) == pa && (*[3]int)(pa) == &a2[1] && tn(pa) == "main.PA"
	})
	check("E12 (*S2)(&ss[0])", func() bool {
		ss := []S{{1}}
		q := (*S2)(&ss[0])
		q.X = 3
		return ss[0].X == 3 && q == (*S2)(&ss[0]) && (*S)(q) == &ss[0] && tn(q) == "*main.S2"
	})
	check("E13 generic T(p)", func() bool {
		n := 1
		q := conv[P](&n)
		return tn(q) == "main.P" && q == P(&n) && back(q) == &n
	})
	check("E14 *ps = S{5}", func() bool {
		s := &S{1}
		ps := PS(s)
		*ps = S{5}
		return s.X == 5
	})
	check("E15 any(P(&n)) == any(P(&n))", func() bool {
		n := 1
		return any(P(&n)) == any(P(&n)) && any(P(&n)) != any(&n)
	})
	check("E16 P(h.p)", func() bool {
		n := 1
		h := HP{&n}
		return P(h.p) == P(&n) && tn(P(h.p)) == "main.P"
	})
	check("E17 P(mk())", func() bool { return P(mk()) == P(&g) && tn(P(mk())) == "main.P" })
	check("E18 PNA(&A)", func() bool {
		var a A
		q := PNA(&a)
		q[0] = 1
		return a[0] == 1 && tn(q) == "main.PNA" && (*A)(q) == &a && PNA(&a) == q
	})
	check("E19 PS2((*S2)(s))", func() bool {
		s := &S{1}
		q := PS2((*S2)(s))
		return tn(q) == "main.PS2" && (*S)((*S2)(q)) == s
	})
	check("E20 P(&n) %v vs &n", func() bool {
		n := 1
		return fmt.Sprintf("%p", P(&n)) == fmt.Sprintf("%p", &n)
	})
	check("E21 reflect Pointer equal", func() bool {
		n := 1
		return reflect.ValueOf(P(&n)).Pointer() == reflect.ValueOf(&n).Pointer()
	})
	check("E22 (*MyInt)(&n) %p", func() bool {
		n := 1
		return fmt.Sprintf("%p", (*MyInt)(&n)) == fmt.Sprintf("%p", &n)
	})

	// implicit
	check("I1 range assign", func() bool {
		n := 1
		var q P
		for _, q = range []*int{&n} {
		}
		return tn(q) == "main.P" && q == P(&n)
	})
	check("I2 closure return", func() bool {
		n := 1
		q := func() P { return &n }()
		return tn(q) == "main.P" && q == P(&n)
	})
	check("I3 method arg", func() bool {
		n := 1
		var h H
		h.set(&n)
		return tn(h.f) == "main.P" && h.f == P(&n)
	})
	check("I4 field assign", func() bool {
		n := 1
		var h H
		h.f = &n
		return tn(h.f) == "main.P" && h.f == P(&n)
	})
	check("I5 *pp = &n", func() bool {
		n := 1
		var q P
		pq := &q
		*pq = &n
		return tn(q) == "main.P" && q == P(&n)
	})
	check("I6 arr[0] = &n", func() bool {
		n := 1
		var arr [1]P
		arr[0] = &n
		return tn(arr[0]) == "main.P" && arr[0] == P(&n)
	})
	check("I7 converse int var", func() bool {
		n := 1
		q := P(&n)
		var p *int = q
		return tn(p) == "*int" && p == &n
	})
	check("I8 converse map/chan/append/literal", func() bool {
		n := 1
		q := P(&n)
		m := map[string]*int{"a": q}
		ch := make(chan *int, 1)
		ch <- q
		sl := append([]*int{}, q)
		st := struct{ p *int }{q}
		c := <-ch
		return tn(m["a"]) == "*int" && tn(c) == "*int" && tn(sl[0]) == "*int" && tn(st.p) == "*int" && c == &n && m["a"] == &n
	})
	check("I9 converse struct ptr", func() bool {
		s := &S{1}
		ps := PS(s)
		var p *S = ps
		return tn(p) == "*main.S" && p == s
	})
	check("I10 tuple assign", func() bool {
		n := 1
		var a *int
		var b P
		a, b = P(&n), &n
		return tn(a) == "*int" && tn(b) == "main.P" && a == &n && b == P(&n)
	})
	check("I11 multi return", func() bool {
		var a P
		var b *int
		a, b = two()
		return tn(a) == "main.P" && tn(b) == "*int" && a == P(&g)
	})
	check("I12 defer arg", func() (ok bool) {
		n := 1
		func() {
			defer func(p P) { ok = tn(p) == "main.P" && p == P(&n) }(&n)
		}()
		return
	})
	check("I13 compare q == &n", func() bool {
		n := 1
		q := P(&n)
		return q == &n && &n == q
	})
	check("I14 switch case", func() bool {
		n := 1
		q := P(&n)
		switch q {
		case &n:
			return true
		}
		return false
	})
	check("I15 map key implicit", func() bool {
		n := 1
		m := map[P]int{}
		m[&n] = 1
		m[P(&n)]++
		delete(m, &n)
		return len(m) == 0
	})
	check("I16 retP", func() bool { return tn(retP()) == "main.P" && retP() == P(&g) })
	check("I17 struct ptr arg", func() bool {
		s := &S{1}
		q := takePS(s)
		return tn(q) == "main.PS" && q == PS(s)
	})
	check("I18 interface from implicit", func() bool {
		n := 1
		var q P = &n
		var i any = q
		_, ok := i.(P)
		return ok && i == any(P(&n))
	})
	check("I19 go stmt arg", func() bool {
		n := 1
		done := make(chan bool)
		go func(p P) { done <- tn(p) == "main.P" && p == P(&n) }(&n)
		return <-done
	})
	check("I20 converse struct field arr", func() bool {
		var a [3]int
		pa := PA(&a)
		x := struct{ a *[3]int }{pa}
		return tn(x.a) == "*[3]int" && x.a == &a
	})

	// nil
	check("N1 PS(nil)", func() bool {
		var s *S
		return PS(s) == nil && (*S2)(s) == nil && (*S)(PS(s)) == nil
	})
	check("N2 unsafe nil arr", func() bool {
		var a *[3]int
		return (*[3]MyInt)(unsafe.Pointer(a)) == nil && (*A)(unsafe.Pointer(a)) == nil
	})
	check("N3 unsafe nil **", func() bool {
		var pp **int
		return (**MyInt)(unsafe.Pointer(pp)) == nil && (*P)(unsafe.Pointer(pp)) == nil
	})
	check("N4 reflect Convert nil *S", func() bool {
		var s *S
		v := reflect.ValueOf(s).Convert(reflect.TypeOf((*S2)(nil)))
		q, ok := v.Interface().(*S2)
		return v.IsNil() && ok && q == nil
	})
	check("N5 reflect Set nil into PS", func() bool {
		ps := PS(&S{})
		reflect.ValueOf(&ps).Elem().Set(reflect.ValueOf((*S)(nil)))
		return ps == nil
	})
	check("N6 NewAt nil struct/array", func() bool {
		return reflect.NewAt(reflect.TypeOf(S{}), nil).IsNil() && reflect.NewAt(reflect.TypeOf([3]int{}), nil).IsNil() &&
			reflect.NewAt(reflect.TypeOf(S{}), nil).Interface().(*S) == nil
	})
	check("N7 implicit nil struct", func() bool {
		var s *S
		var ps PS = s
		var b *S = ps
		return ps == nil && b == nil
	})
	check("N8 generic nil", func() bool { return conv[P](nil) == nil && back(P(nil)) == nil })
	check("N9 typed nil in iface", func() bool {
		var n *int
		var i any = P(n)
		q, ok := i.(P)
		return i != nil && ok && q == nil
	})
	check("N10 zero-size struct non-nil", func() bool {
		var z Z
		return PZ(&z) != nil && (*Z)(PZ(&z)) == &z
	})
	check("N11 zero-size array non-nil", func() bool {
		var a [0]int
		return PA0(&a) != nil && (*A0)(&a) != nil && (*A0)(unsafe.Pointer(&a)) != nil
	})
	check("N12 unsafe.Pointer(P(nil)) == nil", func() bool {
		var n *int
		return unsafe.Pointer(P(n)) == nil && unsafe.Pointer((*MyInt)(n)) == nil
	})
	check("N13 reflect Convert nil P -> *int", func() bool {
		var q P
		v := reflect.ValueOf(q).Convert(reflect.TypeOf((*int)(nil)))
		return v.IsNil() && v.Interface().(*int) == nil
	})
	check("N14 unsafe nil struct to *S2", func() bool {
		var s *S
		return (*S2)(unsafe.Pointer(s)) == nil && PS2((*S2)(unsafe.Pointer(s))) == nil
	})
	check("N15 reflect Call nil", func() bool {
		var n *int
		out := reflect.ValueOf(takeP).Call([]reflect.Value{reflect.ValueOf(n)})
		q, ok := out[0].Interface().(P)
		return ok && q == nil
	})
	check("N16 (*MyInt)(nil field ptr)", func() bool {
		h := HP{}
		return (*MyInt)(h.p) == nil && P(h.p) == nil
	})
	check("N17 nil via unsafe from named", func() bool {
		var q PS
		return (*S2)(unsafe.Pointer(q)) == nil
	})
}
