package main

import "unsafe"

type Pair struct{ A, B int32 }

func main() {
	x := Pair{1, 2}
	p := (*[2]int32)(unsafe.Pointer(&x))
	s := p[:]
	s[0] = 9
	println(x.A, s[0], p[0])
}
