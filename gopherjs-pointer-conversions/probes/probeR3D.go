package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type MyInt64 int64
type Inner struct{ A int }
type Emb struct{ Inner }
type Pair struct{ A, B int32 }
type BB struct{ A, B byte }
type T2 struct {
	A int64
	B MyInt64
}

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

func main() {
	check("D3 *Emb -> *Inner", func() bool {
		e := Emb{Inner{5}}
		in := (*Inner)(unsafe.Pointer(&e))
		in.A = 6
		return e.A == 6 && in == &e.Inner
	})
	check("D3b *Emb -> *Inner read", func() bool {
		e := Emb{Inner{5}}
		in := (*Inner)(unsafe.Pointer(&e))
		return in.A == 5
	})
	check("D4 foreign [2]byte copy", func() bool {
		var s BB
		p := (*[2]byte)(unsafe.Pointer(&s))
		copy(p[:], "hi")
		return s.A == 'h' && s.B == 'i'
	})
	check("D4b foreign [2]byte string", func() bool {
		s := BB{'h', 'i'}
		p := (*[2]byte)(unsafe.Pointer(&s))
		return string(p[:]) == "hi"
	})
	check("D5 foreign array eq twice", func() bool {
		x := Pair{1, 2}
		return (*[2]int32)(unsafe.Pointer(&x)) == (*[2]int32)(unsafe.Pointer(&x))
	})
	check("D6 foreign array NewAt Elem Index", func() bool {
		x := Pair{1, 2}
		v := reflect.NewAt(reflect.TypeOf([2]int32{}), unsafe.Pointer(&x))
		return v.Elem().Index(1).Int() == 2
	})
	check("D7 foreign array Convert unsafe then Elem", func() bool {
		x := Pair{1, 2}
		p := (*[2]int32)(unsafe.Pointer(&x))
		return reflect.ValueOf(p).Elem().Index(0).Int() == 1
	})
	check("D8 *[2]int64 -> *T2 field types", func() bool {
		a := [2]int64{1, 2}
		t := (*T2)(unsafe.Pointer(&a))
		var i any = t.B
		_, ok := i.(MyInt64)
		return ok && t.B == 2
	})
	check("D9 foreign struct write keeps source type", func() bool {
		a := [2]int64{1, 2}
		t := (*T2)(unsafe.Pointer(&a))
		t.B = 7
		var i any = a[1]
		_, ok := i.(int64)
		return ok && a[1] == 7
	})
	check("D10 foreign array json-ish fmt %v", func() bool {
		x := Pair{1, 2}
		p := (*[2]int32)(unsafe.Pointer(&x))
		return fmt.Sprint(*p) == "[1 2]" && fmt.Sprint(p) == "&[1 2]"
	})
	check("D11 foreign array p[:] write", func() bool {
		x := Pair{1, 2}
		p := (*[2]int32)(unsafe.Pointer(&x))
		s := p[:]
		s[0] = 9
		return x.A == 9
	})
}
