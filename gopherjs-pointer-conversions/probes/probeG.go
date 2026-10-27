package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"unsafe"
)

type MyInt int
type AI [2]int
type PA *[2]int
type St struct{ Arr [2]int }

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

func takePA(p PA) int { return p[1] }

func main() {
	check("R1 implicit var PA = &a[1]", func() bool {
		var a [3][2]int
		a[1][1] = 5
		var p PA = &a[1]
		return p[1] == 5
	})
	check("R2 arg PA from &a[1]", func() bool {
		var a [3][2]int
		a[1][1] = 5
		return takePA(&a[1]) == 5
	})
	check("R3 (*AI)(&a[1])", func() bool {
		var a [3][2]int
		a[1][1] = 5
		p := (*AI)(&a[1])
		p[0] = 7
		return p[1] == 5 && a[1][0] == 7
	})
	check("R4 (*AI)(&s[i]) slice", func() bool {
		s := make([][2]int, 3)
		for i := range s {
			r := (*AI)(&s[i])
			r[0] = i
		}
		return s[2][0] == 2
	})
	check("R5 fmt.Sprint(&a[1])", func() bool {
		var a [3][2]int
		a[1][1] = 5
		return fmt.Sprint(&a[1]) == "&[0 5]"
	})
	check("R6 json.Marshal(&s[1])", func() bool {
		s := make([][2]int, 3)
		s[1][1] = 5
		b, err := json.Marshal(&s[1])
		return err == nil && string(b) == "[0,5]"
	})
	check("R7 reflect Elem of &a[1]", func() bool {
		var a [3][2]int
		a[1][1] = 5
		return reflect.ValueOf(&a[1]).Elem().Index(1).Int() == 5
	})
	check("R8 json.Unmarshal into &s[1]", func() bool {
		s := make([][2]int, 3)
		err := json.Unmarshal([]byte("[3,4]"), &s[1])
		return err == nil && s[1][1] == 4
	})
	check("R9 (*AI)(&st.Arr)", func() bool {
		var st St
		p := (*AI)(&st.Arr)
		p[1] = 2
		return st.Arr[1] == 2
	})
	check("R10 (*AI)(&sts[1].Arr)", func() bool {
		sts := make([]St, 2)
		p := (*AI)(&sts[1].Arr)
		p[1] = 2
		return sts[1].Arr[1] == 2
	})
	check("R11 3D (*[2][2]int)(&a[1])", func() bool {
		var a [2][2][2]int
		a[1][1][1] = 9
		type A2 [2][2]int
		p := (*A2)(&a[1])
		return p[1][1] == 9
	})
	check("R12 reflect Convert &a[1] to *AI", func() bool {
		var a [3][2]int
		a[1][1] = 5
		v := reflect.ValueOf(&a[1]).Convert(reflect.TypeOf((*AI)(nil)))
		return v.Interface().(*AI)[1] == 5
	})
	check("R13 Swapper on [][2]int then &s[i]", func() bool {
		s := [][2]int{{1, 2}, {3, 4}}
		_ = fmt.Sprint(&s[0])
		reflect.Swapper(s)(0, 1)
		return s[0][0] == 3 && s[1][0] == 1
	})
	check("R14 (*[2]MyInt)(unsafe) of 2D typed elem twice equal", func() bool {
		var a [3][2]int
		p1 := (*[2]MyInt)(unsafe.Pointer(&a[1]))
		p2 := (*[2]MyInt)(unsafe.Pointer(&a[1]))
		return p1 == p2
	})
}
