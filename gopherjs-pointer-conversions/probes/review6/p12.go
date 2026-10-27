package main

import (
	"fmt"
	"unsafe"
)

type MyInt int

func cv[T, U any](p *T) *U            { return (*U)(unsafe.Pointer(p)) }
func up[T any](p *T) unsafe.Pointer   { return unsafe.Pointer(p) }

func main() {
	var d [2]*[2]int
	d[0] = &[2]int{1, 2}
	d[1] = &[2]int{3, 4}
	pd := cv[[2]*[2]int, [2]*[2]MyInt](&d)
	d[0] = nil
	fmt.Println(1, pd[0] == nil, pd[1] != nil)

	var e [2]*[2]int
	e[0] = &[2]int{1, 2}
	e[1] = &[2]int{3, 4}
	pe := (*[2]*[2]MyInt)(up(&e))
	e[1] = nil
	fmt.Println(2, pe[1] == nil)

	var f [2]*[2]int
	f[0] = &[2]int{1, 2}
	f[1] = &[2]int{3, 4}
	pf := (*[2]*[2]MyInt)(unsafe.Pointer(&f))
	f[1] = nil
	fmt.Println(3, pf[1] == nil)
}
