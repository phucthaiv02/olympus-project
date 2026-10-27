package main

import (
	"fmt"
	"time"
)

type MyInt int
type P *int

func a(x int) bool {
	f := func() *int { return &x }
	return f() == &x
}

func b() (r int, ok bool) {
	f := func() *int { return &r }
	ok = f() == &r
	return
}

func c() (bool, bool) {
	x := 1
	var px1 *int
	g := func() bool {
		x := 2
		inner := func() *int { return &x }
		return inner() == &x
	}
	r1 := g()
	px1 = &x
	h := func() *int { return &x }
	return r1, px1 == h()
}

func d[T any](x T) bool {
	f := func() *T { return &x }
	return f() == &x
}

func e() bool {
	x := 1
	f := func() *int { time.Sleep(time.Millisecond); return &x }
	p := f()
	time.Sleep(time.Millisecond)
	return p == &x
}

var gv = func() bool {
	x := 1
	f := func() *int { return &x }
	return f() == &x && (*MyInt)(f()) == (*MyInt)(&x)
}()

func nested() bool {
	x := 1
	f := func() func() *int {
		y := 0
		_ = &y
		return func() *int { return &x }
	}
	return f()() == &x && P(f()()) == P(&x)
}

type S struct{ v int }

func (s *S) M() *S { return s }

func mv() bool {
	var s S
	f := func() *S { return s.M() }
	return f() == &s
}

func loop() bool {
	ok := true
	var ps []*int
	for i := 0; i < 3; i++ {
		x := i
		f := func() *int { return &x }
		ok = ok && f() == &x
		ps = append(ps, &x)
	}
	return ok && ps[0] != ps[1]
}

func blockingOuter() bool {
	x := 1
	time.Sleep(time.Millisecond)
	f := func() *int { return &x }
	p := &x
	time.Sleep(time.Millisecond)
	return f() == p
}

func blockingOuterNonEsc() bool {
	// not escaping: address taken only inside closure in non-blocking? outer blocking
	x := 1
	f := func() *int { return &x }
	p1 := f()
	time.Sleep(time.Millisecond)
	p2 := f()
	*p2 = 5
	return p1 == p2 && x == 5
}

func main() {
	fmt.Println(a(1))
	fmt.Println(b())
	fmt.Println(c())
	fmt.Println(d(1), d("s"))
	fmt.Println(e())
	fmt.Println(gv)
	fmt.Println(nested())
	fmt.Println(mv())
	fmt.Println(loop())
	fmt.Println(blockingOuter())
	fmt.Println(blockingOuterNonEsc())
}
