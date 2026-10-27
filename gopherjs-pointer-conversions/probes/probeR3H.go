package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type S struct{ X int }
type S2 S
type PS *S
type Pair struct{ A, B int32 }
type Arr [2]int32
type PArr *[2]int32

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
	check("H1 foreign *S -> (*S2)(p)", func() bool {
		x := 5
		p := (*S)(unsafe.Pointer(&x))
		q := (*S2)(p)
		return q.X == 5
	})
	check("H2 foreign *S -> (*S2)(p) write", func() bool {
		x := 5
		p := (*S)(unsafe.Pointer(&x))
		q := (*S2)(p)
		q.X = 6
		return x == 6
	})
	check("H3 foreign *S -> PS implicit", func() bool {
		x := 5
		p := (*S)(unsafe.Pointer(&x))
		var q PS = p
		return q.X == 5
	})
	check("H4 foreign *S -> PS(p)", func() bool {
		x := 5
		p := (*S)(unsafe.Pointer(&x))
		q := PS(p)
		return q.X == 5 && (*S)(q) == p
	})
	check("H5 foreign *S reflect Convert *S2", func() bool {
		x := 5
		p := (*S)(unsafe.Pointer(&x))
		q := reflect.ValueOf(p).Convert(reflect.TypeOf((*S2)(nil))).Interface().(*S2)
		return q.X == 5
	})
	check("H6 foreign from array *S -> *S2", func() bool {
		a := [1]int{5}
		p := (*S)(unsafe.Pointer(&a))
		q := (*S2)(p)
		return q.X == 5
	})
	check("H7 foreign *[2]int32 -> (*Arr)(p)", func() bool {
		x := Pair{1, 2}
		p := (*[2]int32)(unsafe.Pointer(&x))
		q := (*Arr)(p)
		return q[1] == 2
	})
	check("H8 foreign *[2]int32 -> PArr implicit", func() bool {
		x := Pair{1, 2}
		p := (*[2]int32)(unsafe.Pointer(&x))
		var q PArr = p
		return q[1] == 2
	})
	check("H9 foreign *S generic", func() bool {
		x := 5
		p := (*S)(unsafe.Pointer(&x))
		q := conv[*S2](p)
		return q.X == 5
	})
	check("H10 foreign *S Set into PS var via reflect", func() bool {
		x := 5
		p := (*S)(unsafe.Pointer(&x))
		var q PS
		reflect.ValueOf(&q).Elem().Set(reflect.ValueOf(p))
		return q.X == 5
	})
	check("H11 *int -> *MyInt foreign roundtrip", func() bool {
		x := 5
		p := (*S)(unsafe.Pointer(&x))
		q := (*S2)(p)
		return unsafe.Pointer(q) == unsafe.Pointer(&x)
	})
}

func conv[T ~*S2](p *S) T { return T((*S2)(p)) }
