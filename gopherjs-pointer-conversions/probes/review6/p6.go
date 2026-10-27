package main

import "fmt"

type E struct{ F int }
type O struct{ E }
type A struct{ *E }
type B struct{ A }
type GE[T any] struct{ F T }
type GO[T any] struct{ GE[T] }

func try(name string, f func()) {
	defer func() { fmt.Println(name, "recovered:", recover() != nil) }()
	f()
}

func gen[T any]() bool {
	var o GO[T]
	p := &o
	return &o.F == &o.GE.F && &p.F == &o.F
}

func main() {
	try("nil *O", func() { var p *O; _ = &p.F })
	try("nil embedded", func() { var a A; _ = &a.F })
	try("nil embedded deep", func() { var b B; _ = &b.F })
	try("nil *B", func() { var b *B; _ = &b.F })
	fmt.Println(gen[int](), gen[string]())
	e := &E{}
	b := B{A{e}}
	pb := &b
	fmt.Println(&b.F == &e.F, &pb.F == &e.F, &pb.A.F == &b.E.F)
}
