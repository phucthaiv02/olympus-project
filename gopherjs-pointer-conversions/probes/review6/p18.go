package main

import "fmt"

type MyInt int
type PM *MyInt

func (m *MyInt) Self() *MyInt { return m }

func mk() (func() *MyInt, func() *MyInt, *MyInt, PM) {
	var x MyInt
	g := x.Self
	return func() *MyInt { return x.Self() }, g, &x, func() PM { return x.Self() }()
}

func main() {
	f, g, p, pm := mk()
	fmt.Println(f() == p, g() == p, pm == PM(p))
}
