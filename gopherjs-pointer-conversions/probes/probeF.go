package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type MyInt64 int64
type MyInt int
type A [2]int64
type AI [2]int
type S struct{ X, Y int64 }
type W struct{ P *A }

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
	check("Q1 2D elem *[2]int round trip", func() bool {
		var a [3][2]int
		a[1][1] = 5
		p := (*[2]int)(unsafe.Pointer(&a[1]))
		return p == &a[1] && p[1] == 5
	})
	check("Q2 2D elem *[2]int -> *[2]MyInt", func() bool {
		var a [3][2]int
		a[1][1] = 5
		p := (*[2]MyInt)(unsafe.Pointer(&a[1]))
		return p[1] == 5
	})
	check("Q3 2D elem *[2]int64 -> *[2]MyInt64", func() bool {
		var a [3][2]int64
		a[1][1] = 5
		p := (*[2]MyInt64)(unsafe.Pointer(&a[1]))
		return p[1] == 5 && fmt.Sprintf("%T", p[1]) == "main.MyInt64"
	})
	check("Q4 slice of arrays elem round trip", func() bool {
		s := make([][2]int, 3)
		s[2][0] = 4
		p := (*[2]int)(unsafe.Pointer(&s[2]))
		return p == &s[2] && p[0] == 4
	})
	check("Q5 2D elem NewAt", func() bool {
		var a [3][2]int
		a[1][1] = 5
		v := reflect.NewAt(reflect.TypeOf([2]int{}), unsafe.Pointer(&a[1]))
		return v.Interface().(*[2]int) == &a[1]
	})
	check("Q6 2D elem named AI", func() bool {
		var a [3][2]int
		a[1][1] = 5
		p := (*AI)(unsafe.Pointer(&a[1]))
		return p[1] == 5
	})
	check("Q7 struct elem with array field via unsafe", func() bool {
		var a [3]S
		a[1].Y = 3
		p := (*S)(unsafe.Pointer(&a[1]))
		return p == &a[1] && p.Y == 3
	})
	check("Q8 array of arrays explicit P(p)", func() bool {
		type PA *[2]int
		var a [3][2]int
		a[1][1] = 5
		var p PA = &a[1]
		q := (*[2]int)(p)
		return q == &a[1]
	})
	check("Q9 %p records named type, round trip", func() bool {
		var a [2]int64
		_ = fmt.Sprintf("%p", (*A)(&a))
		return (*[2]int64)(unsafe.Pointer(&a)) == &a
	})
	check("Q10 print nested *A then round trip", func() bool {
		var a [2]int64
		_ = fmt.Sprintf("%v", W{(*A)(&a)})
		return (*[2]int64)(unsafe.Pointer(&a)) == &a
	})
	check("Q11 reinterpretSlice records then *A round trip", func() bool {
		var a A
		s := a[:]
		_ = (*(*[]MyInt64)(unsafe.Pointer(&s)))[0]
		return (*A)(unsafe.Pointer(&a)) == &a
	})
	check("Q12 Elem after record", func() bool {
		var a [2]int64
		_ = unsafe.Pointer((*A)(&a))
		return reflect.ValueOf(&a).Elem().Addr().Interface().(*[2]int64) == &a
	})
	check("Q13 DeepEqual then round trip", func() bool {
		var a [2]int64
		pa := (*A)(&a)
		_ = reflect.DeepEqual(pa, pa)
		return (*[2]int64)(unsafe.Pointer(&a)) == &a
	})
	check("Q14a subarray &p[1] == &a[2]", func() bool {
		a := [4]int{1, 2, 3, 4}
		p := (*[2]int)(unsafe.Pointer(&a[1]))
		return &p[1] == &a[2]
	})
	check("Q14b subarray wider conv after &p[1]", func() bool {
		a := [4]int{1, 2, 3, 4}
		p := (*[2]int)(unsafe.Pointer(&a[1]))
		_ = &p[1]
		r := (*[2]int)(unsafe.Pointer(&a[2]))
		return r[1] == 4
	})
	check("Q14c &a[2] first, then p", func() bool {
		a := [4]int{1, 2, 3, 4}
		q := &a[2]
		p := (*[2]int)(unsafe.Pointer(&a[1]))
		return &p[1] == q
	})
	check("Q15 unsafe.Slice after subarray elem ptr (string arr)", func() bool {
		a := [4]string{"a", "b", "c", "d"}
		p := (*[2]string)(unsafe.Pointer(&a[1]))
		_ = &p[1]
		s := unsafe.Slice(&a[2], 2)
		return s[1] == "d"
	})
	check("Q16 unsafe.Slice after subarray elem ptr (int)", func() bool {
		a := [4]int{1, 2, 3, 4}
		p := (*[2]int)(unsafe.Pointer(&a[1]))
		_ = &p[1]
		s := unsafe.Slice(&a[2], 2)
		return s[1] == 4
	})
	check("Q17 slice of subarray-derived then append", func() bool {
		a := [4]int{1, 2, 3, 4}
		p := (*[2]int)(unsafe.Pointer(&a[1]))
		s := unsafe.Slice(&p[0], 2)
		return &s[1] == &a[2] && s[1] == 3
	})
	check("Q18 reflect Index(i).Addr on typed after subarray", func() bool {
		a := [4]int{1, 2, 3, 4}
		p := (*[2]int)(unsafe.Pointer(&a[1]))
		_ = &p[1]
		return reflect.ValueOf(&a).Elem().Index(2).Addr().Interface().(*int) == &a[2]
	})
	check("Q19 *[2]int64 view named S field arr", func() bool {
		type T struct{ A [2]int64 }
		type T2 struct{ A [2]MyInt64 }
		var t T
		v := (*T2)(unsafe.Pointer(&t))
		v.A[1] = 3
		return t.A[1] == 3 && fmt.Sprintf("%T", t.A[1]) == "int64" && unsafe.Pointer(&v.A) == unsafe.Pointer(&t.A)
	})
}
