package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type MyInt int
type MyInt64 int64
type Src2 [1]*[2]int
type Dst2 [1]*[2]MyInt
type Src [1]*[1]int64
type Dst [1]*[1]MyInt64

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
	check("K1 non-nil ptr-to-int-array elems, nil write", func() bool {
		a := Src2{&[2]int{1, 2}}
		q := (*Dst2)(unsafe.Pointer(&a))
		q[0] = nil
		return a[0] == nil
	})
	check("K2 non-nil ptr-to-int-array elems, nil read after src write", func() bool {
		a := Src2{&[2]int{1, 2}}
		q := (*Dst2)(unsafe.Pointer(&a))
		a[0] = nil
		return q[0] == nil
	})
	check("K3 slice of non-nil ptr-to-int-array, nil write", func() bool {
		s := []*[2]int{{1, 2}}
		t := *(*[]*[2]MyInt)(unsafe.Pointer(&s))
		t[0] = nil
		return s[0] == nil
	})
	check("K4 reflect NewAt same", func() bool {
		a := Src2{&[2]int{1, 2}}
		q := reflect.NewAt(reflect.TypeOf(Dst2{}), unsafe.Pointer(&a)).Interface().(*Dst2)
		a[0] = nil
		return q[0] == nil
	})
	check("K5 nested arrays non-nil", func() bool {
		var a [1][1]*[2]int
		a[0][0] = &[2]int{1, 2}
		q := (*[1][1]*[2]MyInt)(unsafe.Pointer(&a))
		q[0][0] = nil
		return a[0][0] == nil
	})
	check("L1 written ptr-to-array elem type", func() bool {
		a := Src{&[1]int64{1}}
		q := (*Dst)(unsafe.Pointer(&a))
		q[0] = &[1]MyInt64{2}
		return fmt.Sprintf("%T", any(a[0][0])) == "int64"
	})
	check("L2 written ptr-to-array elem iface eq", func() bool {
		a := Src{&[1]int64{1}}
		q := (*Dst)(unsafe.Pointer(&a))
		q[0] = &[1]MyInt64{2}
		return any(a[0][0]) == any(int64(2))
	})
	check("L3 struct field written ptr-to-array elem type", func() bool {
		type RS struct{ p *[1]int64 }
		type NS struct{ p *[1]MyInt64 }
		s := &RS{}
		n := (*NS)(unsafe.Pointer(s))
		n.p = &[1]MyInt64{3}
		return fmt.Sprintf("%T", any(s.p[0])) == "int64"
	})
	check("L4 via **ptr (setter path) elem type", func() bool {
		var p *[1]int64
		qq := (**[1]MyInt64)(unsafe.Pointer(&p))
		*qq = &[1]MyInt64{3}
		return fmt.Sprintf("%T", any(p[0])) == "int64"
	})
	check("M1 *byte -> *[4]byte (regression?)", func() bool {
		b := []byte{1, 2, 3, 4}
		q := (*[4]byte)(unsafe.Pointer(&b[0]))
		return q[1] == 2
	})
	check("M2 *int -> *[2]int", func() bool {
		s := []int{1, 2}
		q := (*[2]int)(unsafe.Pointer(&s[0]))
		return q[1] == 2
	})
	check("M3 *T -> *[1]T", func() bool {
		x := 5
		q := (*[1]int)(unsafe.Pointer(&x))
		return q[0] == 5
	})
	check("M4 new(T) struct via unsafe", func() bool {
		type T struct{ A, B int32 }
		u := unsafe.Pointer(new(T))
		(*T)(u).A = 3
		return (*T)(u).A == 3
	})
}
