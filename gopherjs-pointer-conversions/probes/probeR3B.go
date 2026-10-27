package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type MyInt int
type MyInt64 int64
type S struct{ X int }
type Full struct{ A, B int }
type Hdr struct{ X int }
type Pair struct{ A, B int32 }
type Rec struct {
	H [2]int32
	C int32
}
type Inner struct{ A int }
type O struct{ In Inner }

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
	check("B1 &p.X == &p.X", func() bool {
		x := 9
		p := (*S)(unsafe.Pointer(&x))
		return &p.X == &p.X
	})
	check("B2 *&p.X", func() bool {
		x := 9
		p := (*S)(unsafe.Pointer(&x))
		return *&p.X == 9
	})
	check("B3 unsafe(&p.X) == unsafe(&x)", func() bool {
		x := 9
		p := (*S)(unsafe.Pointer(&x))
		return unsafe.Pointer(&p.X) == unsafe.Pointer(&x)
	})
	check("B4 &p.X == &x", func() bool {
		x := 9
		p := (*S)(unsafe.Pointer(&x))
		return &p.X == &x
	})
	check("B5 prefix struct diff names", func() bool {
		f := Full{1, 2}
		h := (*Hdr)(unsafe.Pointer(&f))
		h.X = 5
		return h.X == 1 || (f.A == 5 && h.X == 5)
	})
	check("B5b prefix struct read", func() bool {
		f := Full{1, 2}
		h := (*Hdr)(unsafe.Pointer(&f))
		return h.X == 1
	})
	check("B6 flat array -> struct with array field", func() bool {
		a := [3]int32{1, 2, 3}
		r := (*Rec)(unsafe.Pointer(&a))
		return r.H[1] == 2 && r.C == 3
	})
	check("B7 *O -> *int", func() bool {
		o := O{Inner{4}}
		p := (*int)(unsafe.Pointer(&o))
		return *p == 4
	})
	check("B8 *[2]Pair -> *[4]int32", func() bool {
		a := [2]Pair{{1, 2}, {3, 4}}
		b := (*[4]int32)(unsafe.Pointer(&a))
		return b[3] == 4
	})
	check("B9 NewAt S from &arr", func() bool {
		a := [2]int{5, 6}
		v := reflect.NewAt(reflect.TypeOf(S{}), unsafe.Pointer(&a))
		return v.Elem().Field(0).Int() == 5 && v.Pointer() == uintptr(unsafe.Pointer(&a))
	})
	check("B10 foreign Value.Convert", func() bool {
		x := 9
		p := (*S)(unsafe.Pointer(&x))
		v := reflect.ValueOf(unsafe.Pointer(p)).Convert(reflect.TypeOf(unsafe.Pointer(nil)))
		return (*int)(v.Interface().(unsafe.Pointer)) == &x
	})
	check("B11 foreign Elem().Addr()", func() bool {
		x := 9
		p := (*S)(unsafe.Pointer(&x))
		return reflect.ValueOf(p).Elem().Addr().Interface().(*S) == p
	})
	check("B12 foreign Field.Addr == &p.X", func() bool {
		x := 9
		p := (*S)(unsafe.Pointer(&x))
		return reflect.ValueOf(p).Elem().Field(0).Addr().Interface().(*int) == &p.X
	})
	check("B13 foreign %T and Println", func() bool {
		x := 9
		p := (*S)(unsafe.Pointer(&x))
		return fmt.Sprintf("%T %v", p, p) == "*main.S &{9}"
	})
	check("B14 foreign ptr to array %T", func() bool {
		x := Pair{1, 2}
		p := (*[2]int32)(unsafe.Pointer(&x))
		return fmt.Sprintf("%T %v", p, p) == "*[2]int32 &[1 2]"
	})
	check("B15 foreign array slice p[:]", func() bool {
		x := Pair{1, 2}
		p := (*[2]int32)(unsafe.Pointer(&x))
		s := p[:]
		s[1] = 9
		return x.B == 9 && len(s) == 2
	})
	check("B16 foreign array range", func() bool {
		x := Pair{1, 2}
		p := (*[2]int32)(unsafe.Pointer(&x))
		sum := int32(0)
		for _, v := range p {
			sum += v
		}
		return sum == 3
	})
	check("B17 foreign array &p[1] identity", func() bool {
		x := Pair{1, 2}
		p := (*[2]int32)(unsafe.Pointer(&x))
		return &p[1] == &p[1] && unsafe.Pointer(&p[1]) == unsafe.Pointer(&x.B)
	})
	check("B18 foreign array copy into", func() bool {
		x := Pair{1, 2}
		p := (*[2]int32)(unsafe.Pointer(&x))
		*p = [2]int32{7, 8}
		return x.A == 7 && x.B == 8
	})
	check("B19 foreign array deref value", func() bool {
		x := Pair{1, 2}
		p := (*[2]int32)(unsafe.Pointer(&x))
		v := *p
		x.A = 5
		return v[0] == 1 && v == [2]int32{1, 2}
	})
	check("B20 foreign array reflect Index", func() bool {
		x := Pair{1, 2}
		p := (*[2]int32)(unsafe.Pointer(&x))
		v := reflect.ValueOf(p).Elem()
		v.Index(1).SetInt(6)
		return x.B == 6 && v.Len() == 2
	})
	check("B21 foreign array to slice conv", func() bool {
		x := Pair{1, 2}
		p := (*[2]int32)(unsafe.Pointer(&x))
		s := unsafe.Slice(&p[0], 2)
		return s[1] == 2
	})
	check("B22 foreign struct DeepEqual", func() bool {
		x := 9
		p := (*S)(unsafe.Pointer(&x))
		return reflect.DeepEqual(p, &S{9})
	})
	check("B23 *int -> *MyInt from foreign int", func() bool {
		s := S{3}
		p := (*int)(unsafe.Pointer(&s))
		q := (*MyInt)(p)
		*q = 4
		return s.X == 4 && (*int)(q) == p
	})
	check("B24 foreign array len/cap", func() bool {
		x := Pair{1, 2}
		p := (*[2]int32)(unsafe.Pointer(&x))
		return len(p) == 2 && cap(p[:1]) == 2
	})
}
