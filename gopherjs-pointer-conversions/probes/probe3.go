package main

import (
	"fmt"
	"unsafe"

	"github.com/gopherjs/gopherjs/js"
)

type Inner struct{ E uint16 }
type St struct {
	A int32
	B int64
	C [4]byte
	D Inner
	F float64
}

func check(name string, ok bool) {
	if ok {
		fmt.Println("ok  ", name)
	} else {
		fmt.Println("FAIL", name)
	}
}

func fresh() *St { return &St{A: 1, B: 1 << 40, C: [4]byte{1, 2, 3, 4}, D: Inner{5}, F: 2.5} }

func intact(s *St) bool {
	return s.A == 1 && s.B == 1<<40 && s.C == [4]byte{1, 2, 3, 4} && s.D.E == 5 && s.F == 2.5
}

func kernelWrite(p unsafe.Pointer) { js.InternalObject(p).SetIndex(0, 9) }
func pass(p unsafe.Pointer) unsafe.Pointer { return p }

func main() {
	s := fresh()
	up := unsafe.Pointer(s)
	q := (*St)(up)
	check("two-step round trip identity", q == s)
	check("two-step fields intact", intact(s))
	println(s.A, s.B, len(s.C), s.D.E, s.F)
	println(js.InternalObject(s).Get("B").Get("$high").Int(), js.InternalObject(s).Get("B").Get("$low").Int(), js.InternalObject(s).Get("C").Get("constructor").Get("name").String(), js.InternalObject(s).Get("C").Get("length").Int())

	s2 := fresh()
	q2 := (*St)(pass(unsafe.Pointer(s2)))
	check("one-expression round trip identity", q2 == s2)
	check("one-expression fields intact", intact(s2))

	s3 := fresh()
	kernelWrite(unsafe.Pointer(s3))
	check("syscall output loaded back", s3.A == 9 && s3.B == 1<<40 && s3.D.E == 5 && s3.F == 2.5)

	var np *St
	check("nil to unsafe", unsafe.Pointer(np) == nil)
	var nu unsafe.Pointer
	check("unsafe nil to struct", (*St)(nu) == nil)

	var buf [32]byte
	buf[0] = 7
	buf[16] = 0xff
	b := (*St)(unsafe.Pointer(&buf))
	check("bytes loaded into new struct", b.A == 7 && b.C[0] == 0xff)
}

func init() {
	// The struct's own array keeps its identity and gets the bytes written by
	// a syscall.
	s := fresh()
	c := &s.C
	writeC(unsafe.Pointer(s))
	check("array field written in place", s.C == [4]byte{8, 2, 3, 4} && c == &s.C && len(s.C) == 4)
}

func writeC(p unsafe.Pointer) { js.InternalObject(p).SetIndex(16, 8) }
