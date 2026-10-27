package main

import (
	"fmt"
	"reflect"
	"runtime"
)

type P *int
type S struct{ X int }
type PS *S

var gi = 1

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
	fmt.Println(runtime.Version())
	check("S1 SetIterValue *int -> P", func() bool {
		m := map[string]*int{"a": &gi}
		var q P
		it := reflect.ValueOf(m).MapRange()
		it.Next()
		reflect.ValueOf(&q).Elem().SetIterValue(it)
		return q == P(&gi) && fmt.Sprintf("%T", q) == "main.P"
	})
	check("S2 SetIterKey *int -> P", func() bool {
		m := map[*int]bool{&gi: true}
		var q P
		it := reflect.ValueOf(m).MapRange()
		it.Next()
		reflect.ValueOf(&q).Elem().SetIterKey(it)
		return q == P(&gi)
	})
	check("S3 SetIterValue *S -> PS", func() bool {
		s := &S{3}
		m := map[string]*S{"a": s}
		var q PS
		it := reflect.ValueOf(m).MapRange()
		it.Next()
		reflect.ValueOf(&q).Elem().SetIterValue(it)
		return q == PS(s) && q.X == 3
	})
	check("S4 SetIterValue nil *int -> P", func() bool {
		m := map[string]*int{"a": nil}
		q := P(&gi)
		it := reflect.ValueOf(m).MapRange()
		it.Next()
		reflect.ValueOf(&q).Elem().SetIterValue(it)
		return q == nil
	})
	check("S5 Send *int on chan P", func() bool {
		c := make(chan P, 1)
		reflect.ValueOf(c).Send(reflect.ValueOf(&gi))
		q := <-c
		return q == P(&gi) && fmt.Sprintf("%T", q) == "main.P"
	})
	check("S6 Select send *int on chan P", func() bool {
		c := make(chan P, 1)
		reflect.Select([]reflect.SelectCase{{Dir: reflect.SelectSend, Chan: reflect.ValueOf(c), Send: reflect.ValueOf(&gi)}})
		q := <-c
		return q == P(&gi)
	})
	check("S7 Set P into same-typed var via Value from MapIndex", func() bool {
		m := map[string]*int{"a": &gi}
		var q P
		reflect.ValueOf(&q).Elem().Set(reflect.ValueOf(m).MapIndex(reflect.ValueOf("a")))
		return q == P(&gi)
	})
}
