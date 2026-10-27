package main

import (
	"fmt"
	"reflect"

	"github.com/gopherjs/gopherjs/js"
)

type W struct {
	*js.Object
	N int `js:"n"`
}

type H struct {
	O *js.Object
	K int
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
	o := js.Global.Get("Object").New()
	o.Set("n", 4)
	check("G1 reflect Elem of *js.Object", func() bool {
		v := reflect.ValueOf(o)
		return v.Elem().IsValid()
	})
	check("G2 fmt %v of *js.Object", func() bool {
		s := fmt.Sprintf("%v", o)
		fmt.Println("  ->", s)
		return true
	})
	check("G3 W field via reflect", func() bool {
		w := &W{Object: o}
		return reflect.ValueOf(w).Elem().Field(1).Int() == 4
	})
	check("G4 H field js obj via reflect", func() bool {
		h := &H{O: o, K: 2}
		v := reflect.ValueOf(h).Elem().Field(0)
		return v.Interface().(*js.Object).Get("n").Int() == 4
	})
	check("G5 H field Addr", func() bool {
		h := &H{O: o, K: 2}
		p := reflect.ValueOf(h).Elem().Field(0).Addr().Interface().(**js.Object)
		return (*p).Get("n").Int() == 4 && p == &h.O
	})
	check("G6 slice of *js.Object Index", func() bool {
		s := []*js.Object{o}
		return reflect.ValueOf(s).Index(0).Interface().(*js.Object).Get("n").Int() == 4
	})
}
