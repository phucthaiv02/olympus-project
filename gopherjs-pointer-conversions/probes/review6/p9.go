package main

import "fmt"

type MyInt int
type P *int

func mk() (func() *int, *int) {
	x := 1
	return func() *int { return &x }, &x
}

func mkLate() (func() *int, func() *int) {
	x := 1
	f := func() *int { return &x }
	g := func() *int {
		return func() *int { return &x }()
	}
	return f, g
}

func mkConv() (func() P, P) {
	x := 1
	return func() P { return &x }, P(&x)
}

func mkNamed() (r int, f func() *int, p *int) {
	f = func() *int { return &r }
	p = &r
	return
}

func mkParam(x int) (func() *int, *int) {
	return func() *int { return &x }, &x
}

func mkGen[T any](x T) (func() *T, *T) {
	return func() *T { return &x }, &x
}

func (m MyInt) mkRecv() (func() *MyInt, *MyInt) {
	return func() *MyInt { return &m }, &m
}

func mkShadow() (func() *int, *int, func() (*int, *int)) {
	x := 1
	g := func() (*int, *int) {
		x := 2
		h := func() *int { return &x }
		return h(), &x
	}
	f := func() *int { return &x }
	return f, &x, g
}

func mkDeep() (func() func() *int, *int) {
	x := 0
	return func() func() *int {
		y := 5
		_ = &y
		return func() *int { y++; return &x }
	}, &x
}

var pkgF, pkgP = func() (func() *int, *int) {
	z := 3
	return func() *int { return &z }, &z
}()

func main() {
	f, p := mk()
	fmt.Println(1, f() == p)
	f1, g1 := mkLate()
	fmt.Println(2, f1() == g1())
	fc, pc := mkConv()
	fmt.Println(3, fc() == pc)
	_, fn, pn := mkNamed()
	fmt.Println(4, fn() == pn)
	fp, pp := mkParam(3)
	fmt.Println(5, fp() == pp)
	fg, pg := mkGen("s")
	fmt.Println(6, fg() == pg)
	fr, pr := MyInt(2).mkRecv()
	fmt.Println(7, fr() == pr)
	fs, ps, gs := mkShadow()
	a, b := gs()
	fmt.Println(8, fs() == ps, a == b, a != ps)
	fd, pd := mkDeep()
	fmt.Println(9, fd()() == pd)
	fmt.Println(10, pkgF() == pkgP)
	mv := MyInt(4).mkRecv
	m1, m2 := mv()
	fmt.Println(11, m1() == m2)
}
