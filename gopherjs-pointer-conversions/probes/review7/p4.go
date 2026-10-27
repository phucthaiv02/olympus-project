package main

import (
	"fmt"
	"reflect"
	"sync/atomic"
	"unsafe"
)

type E struct{ A int }
type Emb struct {
	E
	Y int
}
type ZB struct {
	_ [0]int
	X int
}
type ZS struct {
	Z struct{}
	X int
}
type PA struct {
	P *[2]int
	Y int
}
type A struct{ X int64 }
type B struct{ X int64 }
type G[T any] struct{ V T }
type MyInt int
type MyInt64 int64

func try(name string, f func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(name, "PANIC:", r)
		}
	}()
	f()
}

func main() {
	try("emb", func() {
		s := &Emb{}
		fmt.Println("emb:", unsafe.Pointer(s) == unsafe.Pointer(&s.E), unsafe.Pointer(&s.E) == unsafe.Pointer(&s.A), unsafe.Pointer(&s.A) == unsafe.Pointer(&s.E.A), &s.A == &s.E.A)
		fmt.Println("emb back:", (*Emb)(unsafe.Pointer(&s.A)) == s, (*E)(unsafe.Pointer(&s.A)) == &s.E, (*E)(unsafe.Pointer(s)) == &s.E)
	})
	try("zb", func() {
		s := &ZB{}
		fmt.Println("zb:", unsafe.Pointer(s) == unsafe.Pointer(&s.X))
		fmt.Println("zb back:", (*ZB)(unsafe.Pointer(&s.X)) == s)
	})
	try("zs", func() {
		s := &ZS{}
		fmt.Println("zs:", unsafe.Pointer(s) == unsafe.Pointer(&s.X))
	})
	try("pa", func() {
		arr := [2]int{1, 2}
		s := &PA{P: &arr}
		fmt.Println("pa:", unsafe.Pointer(s) == unsafe.Pointer(&s.P), (*PA)(unsafe.Pointer(&s.P)) == s, unsafe.Pointer(s.P) == unsafe.Pointer(&arr[0]))
		ps := [2]*[2]int{&arr, nil}
		fmt.Println("arr of arrptr:", unsafe.Pointer(&ps) == unsafe.Pointer(&ps[0]), (*[2]*[2]int)(unsafe.Pointer(&ps[0])) == &ps)
		v := reflect.ValueOf(&ps).Pointer() == uintptr(unsafe.Pointer(&ps[0]))
		fmt.Println("refl arr of arrptr:", v)
	})
	try("arr structs", func() {
		a := [2]A{{1}, {2}}
		b := (*[2]B)(unsafe.Pointer(&a))
		b[1].X = 5
		fmt.Println("as:", a[1].X, (*A)(unsafe.Pointer(&b[1])) == &a[1], unsafe.Pointer(&b[1]) == unsafe.Pointer(&a[1]), unsafe.Pointer(&b[1].X) == unsafe.Pointer(&a[1].X))
		fmt.Println("as2:", (*[2]A)(unsafe.Pointer(b)) == &a, (*[2]A)(unsafe.Pointer(&b[0])) == &a, (*B)(unsafe.Pointer(&a[0])) == &b[0])
		bs := b[:]
		bs[0] = B{9}
		fmt.Println("as3:", a[0].X, fmt.Sprintf("%T", b[0]), fmt.Sprintf("%T", bs[1]))
		var ifc interface{} = b[1]
		_, ok := ifc.(B)
		fmt.Println("as4:", ok)
	})
	try("generic", func() {
		g := &G[int64]{V: 3}
		h := (*G[MyInt64])(unsafe.Pointer(g))
		h.V = 4
		fmt.Println("gen:", g.V, (*G[int64])(unsafe.Pointer(h)) == g, unsafe.Pointer(&h.V) == unsafe.Pointer(&g.V), fmt.Sprintf("%T", h.V))
		ga := &G[[2]int64]{}
		gb := (*G[[2]MyInt64])(unsafe.Pointer(ga))
		gb.V[1] = 7
		fmt.Println("gen arr:", ga.V[1], unsafe.Pointer(&gb.V[1]) == unsafe.Pointer(&ga.V[1]), fmt.Sprintf("%T", gb.V[0]))
	})
	try("atomic", func() {
		var ap atomic.Pointer[int]
		x, y := 5, 6
		ap.Store(&x)
		fmt.Println("atomic:", ap.Load() == &x, ap.CompareAndSwap(&x, &y), ap.Load() == &y)
		var aa atomic.Pointer[[2]int]
		arr := [2]int{}
		aa.Store(&arr)
		fmt.Println("atomic arr:", aa.Load() == &arr, aa.CompareAndSwap(&arr, nil), aa.Load() == nil)
		var as atomic.Pointer[E]
		e := &E{}
		as.Store(e)
		fmt.Println("atomic struct:", as.Load() == e, as.Swap(nil) == e)
		var au unsafe.Pointer
		atomic.StorePointer(&au, unsafe.Pointer(&arr))
		fmt.Println("atomic sp:", atomic.LoadPointer(&au) == unsafe.Pointer(&arr[0]), atomic.CompareAndSwapPointer(&au, unsafe.Pointer(&arr[0]), nil))
	})
	try("setpointer", func() {
		var up unsafe.Pointer
		arr := [2]int{}
		reflect.ValueOf(&up).Elem().SetPointer(unsafe.Pointer(&arr))
		fmt.Println("setptr:", up == unsafe.Pointer(&arr[0]), (*[2]int)(up) == &arr)
		v := reflect.NewAt(reflect.TypeOf(arr), up)
		fmt.Println("newat:", v.Interface().(*[2]int) == &arr, v.Pointer() == uintptr(up), v.UnsafePointer() == up)
		v2 := reflect.NewAt(reflect.TypeOf(0), up)
		fmt.Println("newat int:", v2.Interface().(*int) == &arr[0], v2.Elem().UnsafeAddr() == uintptr(up))
		v3 := reflect.ValueOf(&arr).Convert(reflect.TypeOf((*[2]MyInt)(nil)))
		fmt.Println("conv:", v3.Pointer() == uintptr(up), unsafe.Pointer(v3.Interface().(*[2]MyInt)) == up)
		w := reflect.ValueOf(&arr).Elem().Index(0).Addr()
		fmt.Println("idx addr:", w.Interface().(*int) == &arr[0], w.Pointer() == uintptr(up))
		fmt.Println("unsafeptr:", reflect.ValueOf(&arr).UnsafePointer() == up, reflect.ValueOf(arr[:]).UnsafePointer() == up)
	})
	try("map", func() {
		s := &Emb{}
		m := map[unsafe.Pointer]int{}
		m[unsafe.Pointer(s)]++
		m[unsafe.Pointer(&s.E)]++
		m[unsafe.Pointer(&s.A)]++
		m[unsafe.Pointer(&s.Y)]++
		fmt.Println("map:", len(m))
	})
}
