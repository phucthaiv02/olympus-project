package main

import (
	"fmt"
	"reflect"

	"github.com/gopherjs/gopherjs/js"
)

type W struct {
	O *js.Object
	N int
}

func try(name string, f func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(name, "PANIC:", r)
		}
	}()
	f()
}

func main() {
	o := js.Global.Get("Object").New()
	try("slice", func() { s := fmt.Sprint([]*js.Object{nil, o}); fmt.Println("slice ok", len(s) > 0) })
	try("p", func() { s := fmt.Sprintf("%p", o); fmt.Println("p ok", len(s) > 0) })
	try("pnil", func() { var n *js.Object; s := fmt.Sprintf("%p", n); fmt.Println("pnil:", s) })
	try("struct", func() { s := fmt.Sprintf("%v", W{nil, 1}); fmt.Println("struct:", s) })
	try("structp", func() { s := fmt.Sprintf("%v", &W{o, 1}); fmt.Println("structp ok", len(s) > 0) })
	try("ptr", func() { fmt.Println("ptr eq:", reflect.ValueOf(o).Pointer() == reflect.ValueOf(o).Pointer()) })
	try("uptr", func() { fmt.Println("uptr:", reflect.ValueOf(o).UnsafePointer() != nil) })
	try("deep", func() { fmt.Println("deepequal:", reflect.DeepEqual([]*js.Object{o}, []*js.Object{o})) })
}
