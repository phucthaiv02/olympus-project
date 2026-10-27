package main

import (
	"fmt"
	"unsafe"
)

type MyInt32 int32
type MyInt int
type MyByte byte
type MyFloat float64
type MyStr string
type MyBool bool

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
	check("I1 int32 array elem type+identity (src first)", func() bool {
		var a [2]int32
		p := &a[0]
		q := &(*[2]MyInt32)(unsafe.Pointer(&a))[0]
		var i any = q
		_, ok := i.(*MyInt32)
		return ok && q == (*MyInt32)(unsafe.Pointer(p)) && unsafe.Pointer(q) == unsafe.Pointer(p) && fmt.Sprintf("%T", q) == "*main.MyInt32"
	})
	check("I2 int32 array (dst first)", func() bool {
		var a [2]int32
		q := &(*[2]MyInt32)(unsafe.Pointer(&a))[1]
		p := &a[1]
		var i any = p
		_, ok := i.(*int32)
		return ok && (*int32)(unsafe.Pointer(q)) == p && unsafe.Pointer(q) == unsafe.Pointer(p) && fmt.Sprintf("%T", p) == "*int32"
	})
	check("I3 write through", func() bool {
		var a [2]int32
		p := &a[0]
		q := &(*[2]MyInt32)(unsafe.Pointer(&a))[0]
		*q = 5
		*p += 1
		return a[0] == 6 && *q == 6
	})
	check("I4 repeated", func() bool {
		var a [2]int32
		_ = &a[0]
		v := (*[2]MyInt32)(unsafe.Pointer(&a))
		return &v[0] == &v[0] && &v[0] != &v[1]
	})
	check("I5 int array", func() bool {
		var a [3]int
		p := &a[2]
		q := &(*[3]MyInt)(unsafe.Pointer(&a))[2]
		return fmt.Sprintf("%T", q) == "*main.MyInt" && q == (*MyInt)(unsafe.Pointer(p))
	})
	check("I6 byte array", func() bool {
		var a [4]byte
		p := &a[1]
		q := &(*[4]MyByte)(unsafe.Pointer(&a))[1]
		return fmt.Sprintf("%T", q) == "*main.MyByte" && q == (*MyByte)(unsafe.Pointer(p))
	})
	check("I7 float64 array", func() bool {
		var a [2]float64
		p := &a[0]
		q := &(*[2]MyFloat)(unsafe.Pointer(&a))[0]
		return fmt.Sprintf("%T", q) == "*main.MyFloat" && q == (*MyFloat)(unsafe.Pointer(p))
	})
	check("I8 string array (plain Array)", func() bool {
		var a [2]string
		p := &a[0]
		q := &(*[2]MyStr)(unsafe.Pointer(&a))[0]
		*q = "x"
		return fmt.Sprintf("%T", q) == "*main.MyStr" && q == (*MyStr)(unsafe.Pointer(p)) && a[0] == "x"
	})
	check("I9 bool array", func() bool {
		var a [2]bool
		p := &a[1]
		q := &(*[2]MyBool)(unsafe.Pointer(&a))[1]
		return fmt.Sprintf("%T", q) == "*main.MyBool" && q == (*MyBool)(unsafe.Pointer(p))
	})
	check("I10 slice elem", func() bool {
		s := []int32{1, 2, 3}
		p := &s[1]
		ms := *(*[]MyInt32)(unsafe.Pointer(&s))
		q := &ms[1]
		return fmt.Sprintf("%T", q) == "*main.MyInt32" && q == (*MyInt32)(unsafe.Pointer(p)) && *q == 2
	})
	check("I11 subarray view of slice", func() bool {
		s := make([]int32, 4)
		p := &s[2]
		v := (*[2]MyInt32)(unsafe.Pointer(&s[1]))
		_ = v
		q := (*MyInt32)(unsafe.Pointer(p))
		return fmt.Sprintf("%T", q) == "*main.MyInt32"
	})
	check("I12 map key", func() bool {
		var a [2]int32
		m := map[*MyInt32]int{}
		m[&(*[2]MyInt32)(unsafe.Pointer(&a))[0]] = 1
		m[(*MyInt32)(unsafe.Pointer(&a[0]))]++
		return len(m) == 1 && m[(*MyInt32)(unsafe.Pointer(&a[0]))] == 2
	})
	check("I13 same type unaffected", func() bool {
		var a [2]int32
		p := &a[0]
		var i any = &a[0]
		_, ok := i.(*int32)
		return p == &a[0] && ok
	})
}
