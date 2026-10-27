package main

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"unsafe"
)

type MyInt int
type IP *int
type MIP *MyInt
type S struct{ X, Y int }
type SP *S
type AP *[2]int
type MyErr struct{ s string }

func (e *MyErr) Error() string { return e.s }

type Pair struct{ K, V int }
type Pairs []Pair

func try(name string, f func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(name, "PANIC:", r)
		}
	}()
	f()
}

func main() {
	arr := [2]int{1, 2}
	up := unsafe.Pointer(&arr)
	try("refl", func() {
		v3 := reflect.ValueOf(&arr).Convert(reflect.TypeOf(AP(nil)))
		fmt.Println("conv:", v3.Pointer() == uintptr(up), unsafe.Pointer(v3.Interface().(AP)) == up)
		w := reflect.ValueOf(&arr).Elem().Index(0).Addr()
		fmt.Println("idx addr:", w.Interface().(*int) == &arr[0], w.Pointer() == uintptr(up))
		fmt.Println("unsafeptr:", reflect.ValueOf(&arr).UnsafePointer() == up, reflect.ValueOf(arr[:]).UnsafePointer() == up)
		sl := arr[:]
		fmt.Println("slice ptr:", reflect.ValueOf(sl).Pointer() == uintptr(up), reflect.ValueOf(sl[1:]).Pointer() == uintptr(unsafe.Pointer(&arr[1])))
		fmt.Println("empty slice ptr:", reflect.ValueOf(sl[2:]).Pointer() == uintptr(unsafe.Pointer(unsafe.SliceData(sl[2:]))))
	})
	try("named ptr elem", func() {
		x := 3
		var ip IP = &x
		v := reflect.ValueOf(&ip).Elem() // IP
		e := v.Elem()
		fmt.Println("elem addr:", e.Addr().Interface().(*int) == &x, e.UnsafeAddr() == uintptr(unsafe.Pointer(&x)))
		var mp MIP = (*MyInt)(unsafe.Pointer(&x))
		fmt.Println("mip:", reflect.ValueOf(mp).Elem().Addr().Interface().(*MyInt) == (*MyInt)(unsafe.Pointer(&x)))
		s := &S{}
		var sp SP = s
		fmt.Println("sp:", reflect.ValueOf(sp).Elem().Addr().Interface().(*S) == s, reflect.ValueOf(sp).Elem().Field(1).Addr().Interface().(*int) == &s.Y)
		var ap AP = &arr
		fmt.Println("ap:", reflect.ValueOf(ap).Elem().Addr().Interface().(*[2]int) == &arr, reflect.ValueOf(ap).Pointer() == uintptr(up))
	})
	try("set", func() {
		x := 3
		var ip IP
		reflect.ValueOf(&ip).Elem().Set(reflect.ValueOf(&x))
		fmt.Println("set:", (*int)(ip) == &x)
		var p *int
		reflect.ValueOf(&p).Elem().Set(reflect.ValueOf(IP(&x)))
		fmt.Println("set2:", p == &x)
		m := map[string]IP{}
		reflect.ValueOf(m).SetMapIndex(reflect.ValueOf("a"), reflect.ValueOf(&x))
		fmt.Println("setmap:", (*int)(m["a"]) == &x)
		ch := make(chan IP, 1)
		reflect.ValueOf(ch).Send(reflect.ValueOf(&x))
		fmt.Println("send:", (*int)(<-ch) == &x)
		f := reflect.MakeFunc(reflect.TypeOf(func() IP { return nil }), func([]reflect.Value) []reflect.Value { return []reflect.Value{reflect.ValueOf(&x)} }).Interface().(func() IP)
		fmt.Println("makefunc:", (*int)(f()) == &x)
		var np *int
		reflect.ValueOf(&ip).Elem().Set(reflect.ValueOf(np))
		fmt.Println("setnil:", ip == nil)
		res := reflect.ValueOf(func(p IP) *int { return p }).Call([]reflect.Value{reflect.ValueOf(&x)})
		fmt.Println("call:", res[0].Interface().(*int) == &x)
		sl := []IP{nil}
		reflect.Copy(reflect.ValueOf(sl), reflect.ValueOf([]IP{&x}))
		fmt.Println("copy:", (*int)(sl[0]) == &x)
		reflect.ValueOf(sl).Index(0).Set(reflect.ValueOf(&x))
		fmt.Println("idxset:", (*int)(sl[0]) == &x)
		app := reflect.Append(reflect.ValueOf(sl), reflect.ValueOf(&x)).Interface().([]IP)
		fmt.Println("append:", (*int)(app[1]) == &x)
	})
	try("errors.As", func() {
		var e error = &MyErr{"x"}
		var target *MyErr
		fmt.Println("as:", errors.As(e, &target), target == e.(*MyErr))
	})
	try("sort", func() {
		ps := []Pair{{3, 0}, {1, 1}, {2, 2}}
		b := (*[3]struct{ K, V int })(unsafe.Pointer(&ps[0]))
		sort.Slice(b[:], func(i, j int) bool { return b[i].K < b[j].K })
		fmt.Println("sort:", ps)
		ps2 := Pairs{{3, 0}, {1, 1}, {2, 2}}
		sort.Slice(ps2, func(i, j int) bool { return ps2[i].K < ps2[j].K })
		fmt.Println("sort2:", ps2)
	})
}
