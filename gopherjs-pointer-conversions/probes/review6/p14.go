package main

import "fmt"

type P *int
type H struct{ p P }
type I interface{ M(P) bool }
type R struct{}

var gx = new(int)

func (R) M(p P) bool { return p == P(gx) }

func check(name string, p P) { fmt.Println(name, p == P(gx)) }

func main() {
	defer check("1 defer", gx)
	var i I = R{}
	fmt.Println("2 iface call", i.M(gx))
	done := make(chan bool)
	go func(p P) { fmt.Println("3 go", p == P(gx)); done <- true }(gx)
	<-done
	var h H
	h.p = gx
	ph := &h
	ph.p = gx
	fmt.Println("4 field", h.p == P(gx))
	var q P
	pq := &q
	*pq = gx
	fmt.Println("5 deref assign", q == P(gx))
	fn := func() P { return gx }
	fmt.Println("6 funclit ret", fn() == P(gx))
	chP := make(chan P, 1)
	select {
	case chP <- gx:
	}
	fmt.Println("7 select send", <-chP == P(gx))
	arr := []P{nil}
	arr[0] = gx
	fmt.Println("8 index assign", arr[0] == P(gx))
	m := map[string]P{}
	m["a"] = gx
	fmt.Println("9 map assign", m["a"] == P(gx))
	var q2 P
	q2 = gx
	fmt.Println("10 cmp", q2 == gx, gx == q2)
	switch q2 {
	case gx:
		fmt.Println("11 switch true")
	default:
		fmt.Println("11 switch false")
	}
	var iface interface{} = q2
	_, isP := iface.(P)
	fmt.Println("12", isP)
	mi := map[interface{}]int{P(gx): 1}
	fmt.Println("13", mi[q2])
}
