package main

import (
	"fmt"
	"unsafe"
)

type MyInt int
type MyInt64 int64

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
	check("M1 map value type via unsafe", func() bool {
		m := map[string]int64{"a": 1}
		q := (*map[string]MyInt64)(unsafe.Pointer(&m))
		return fmt.Sprintf("%T", any((*q)["a"])) == "main.MyInt64"
	})
	check("M2 map ptr value identity via unsafe", func() bool {
		x := 1
		m := map[string]*int{"a": &x}
		q := (*map[string]*MyInt)(unsafe.Pointer(&m))
		return (*q)["a"] == (*MyInt)(&x)
	})
	check("M3 map write keeps orig type", func() bool {
		m := map[string]int64{}
		q := (*map[string]MyInt64)(unsafe.Pointer(&m))
		(*q)["b"] = 5
		return fmt.Sprintf("%T", any(m["b"])) == "int64"
	})
	check("M4 map key type via unsafe", func() bool {
		m := map[int64]bool{3: true}
		q := (*map[MyInt64]bool)(unsafe.Pointer(&m))
		for k := range *q {
			return fmt.Sprintf("%T", any(k)) == "main.MyInt64"
		}
		return false
	})
	check("C1 chan elem via unsafe", func() bool {
		c := make(chan int64, 1)
		c <- 7
		q := (*chan MyInt64)(unsafe.Pointer(&c))
		v := <-*q
		return fmt.Sprintf("%T", any(v)) == "main.MyInt64"
	})
	check("C2 chan of ptr via unsafe", func() bool {
		x := 1
		c := make(chan *int, 1)
		c <- &x
		q := (*chan *MyInt)(unsafe.Pointer(&c))
		return <-*q == (*MyInt)(&x)
	})
	check("F1 func result via unsafe", func() bool {
		f := func() int64 { return 3 }
		q := (*func() MyInt64)(unsafe.Pointer(&f))
		return fmt.Sprintf("%T", any((*q)())) == "main.MyInt64"
	})
	check("S1 slice of maps via unsafe", func() bool {
		s := []map[string]int64{{"a": 1}}
		q := (*[]map[string]MyInt64)(unsafe.Pointer(&s))
		return fmt.Sprintf("%T", any((*q)[0]["a"])) == "main.MyInt64"
	})
	check("A1 array of maps via unsafe", func() bool {
		a := [1]map[string]int64{{"a": 1}}
		q := (*[1]map[string]MyInt64)(unsafe.Pointer(&a))
		return fmt.Sprintf("%T", any(q[0]["a"])) == "main.MyInt64"
	})
	check("ST1 struct field map via unsafe", func() bool {
		type A struct{ M map[string]int64 }
		type B struct{ M map[string]MyInt64 }
		a := &A{M: map[string]int64{"a": 1}}
		b := (*B)(unsafe.Pointer(a))
		return fmt.Sprintf("%T", any(b.M["a"])) == "main.MyInt64"
	})
}
