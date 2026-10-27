package main

import (
	"fmt"
	"reflect"
	"sort"
	"unsafe"
)

type MyInt int
type P *int
type St struct{ A, B int64 }

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
	check("T1 sort.Slice [][2]int", func() bool {
		s := [][2]int{{3, 0}, {1, 0}, {2, 0}}
		sort.Slice(s, func(i, j int) bool { return s[i][0] < s[j][0] })
		return s[0][0] == 1 && s[2][0] == 3
	})
	check("T2 sort.Slice []St", func() bool {
		s := []St{{3, 0}, {1, 0}, {2, 0}}
		sort.Slice(s, func(i, j int) bool { return s[i].A < s[j].A })
		return s[0].A == 1 && s[2].A == 3
	})
	check("T3 MakeFunc returning P", func() bool {
		x := 1
		var f func() P
		v := reflect.MakeFunc(reflect.TypeOf(f), func([]reflect.Value) []reflect.Value {
			return []reflect.Value{reflect.ValueOf(&x)}
		})
		reflect.ValueOf(&f).Elem().Set(v)
		p := f()
		return (*int)(p) == &x
	})
	check("T4 cap0 subslice append", func() bool {
		s := []int{1, 2, 3}
		t := s[3:]
		t = append(t, 4)
		return len(t) == 1 && t[0] == 4 && s[2] == 3
	})
	check("T5 cap0 then reslice of parent", func() bool {
		s := make([]int, 3, 5)
		t := s[1:1:1]
		u := s[1:4]
		u[2] = 7
		return len(t) == 0 && s[:4][3] == 7
	})
	check("T6 view struct elem write via 2D? &a[1] struct elem", func() bool {
		var a [3]St
		p := (*[2]St)(unsafe.Pointer(&a[1]))
		p[1].B = 4
		return a[2].B == 4 && &p[1] == &a[2]
	})
	check("T7 Set PA via reflect from &a[1]", func() bool {
		type PA *[2]int
		var a [3][2]int
		var pa PA
		reflect.ValueOf(&pa).Elem().Set(reflect.ValueOf(&a[1]))
		return pa != nil
	})
	check("T8 struct elem ptr to *[1]St", func() bool {
		var a [3]St
		a[1].A = 5
		p := (*[1]St)(unsafe.Pointer(&a[1]))
		return p[0].A == 5
	})
}
