package main

import "fmt"

type P *string
type A *[2]int
type S struct{ f P }

var gs = "g"
var g P = &gs
var garr A = (*[2]int)(nil)

func check(name string, ok bool) {
	if !ok {
		fmt.Println("FAIL", name)
	} else {
		fmt.Println("ok  ", name)
	}
}

func main() {
	check("pkg var type", fmt.Sprintf("%T", g) == "main.P")
	check("pkg var nil array", garr == nil)
	s := "x"
	p := &s
	var st S
	st.f = p
	check("field assign", fmt.Sprintf("%T", st.f) == "main.P")
	arr := [1]P{}
	arr[0] = p
	check("array index assign", fmt.Sprintf("%T", arr[0]) == "main.P")
	f := func() P { return p }
	check("closure return", fmt.Sprintf("%T", f()) == "main.P")
	var got string
	func() {
		defer func(q P) { got = fmt.Sprintf("%T", q) }(p)
	}()
	check("defer arg", got == "main.P")
	var ip interface{} = P(p)
	check("iface eq", ip == interface{}(P(p)))
	var x, y P = p, p
	check("multi var decl", x == y && fmt.Sprintf("%T", y) == "main.P")
	var na *[2]int
	var a2 A
	a2 = na
	check("assign nil arr", a2 == nil)
	m := map[P]int{}
	m[p] = 1
	m[P(p)]++
	check("map key", len(m) == 1 && m[p] == 2)
	pp := &p
	type PP **string
	q := PP(pp)
	**q = "w"
	check("ptrptr write", s == "w" && (**string)(q) == pp)
}
