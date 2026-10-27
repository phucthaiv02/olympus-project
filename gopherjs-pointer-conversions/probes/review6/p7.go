package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"math"
	"math/big"
	"os"
	"reflect"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"
)

type MyInt int

func main() {
	var sb strings.Builder
	sb.WriteString("hello ")
	sb.WriteString("world")
	fmt.Println(sb.String())
	fmt.Printf("%x\n", sha256.Sum256([]byte("abc")))
	fmt.Println(crc32.ChecksumIEEE([]byte("abc")))
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, 0x0102030405060708)
	fmt.Println(b, binary.BigEndian.Uint32(b))
	j, _ := json.Marshal(map[string]any{"a": 1, "b": []int{1, 2}})
	fmt.Println(string(j))
	fmt.Println(math.Float64bits(1.5), new(big.Int).Exp(big.NewInt(3), big.NewInt(50), nil))
	os.WriteFile("/tmp/x/f.txt", []byte("data"), 0644)
	d, err := os.ReadFile("/tmp/x/f.txt")
	fmt.Println(string(d), err)
	ents, _ := os.ReadDir("/tmp/x")
	fmt.Println(len(ents) > 0)
	s := []MyInt{3, 1, 2}
	sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })
	fmt.Println(s)
	var m sync.Map
	m.Store(1, 2)
	v, _ := m.Load(1)
	fmt.Println(v)
	var ap atomic.Pointer[int]
	x := 5
	ap.Store(&x)
	fmt.Println(ap.Load() == &x, *ap.Load())
	var up unsafe.Pointer
	atomic.StorePointer(&up, unsafe.Pointer(&x))
	fmt.Println((*int)(atomic.LoadPointer(&up)) == &x)
	fmt.Println(bytes.Equal([]byte("a"), []byte("a")), unsafe.String(unsafe.SliceData([]byte("hi")), 2))
	bs := []byte("xyz")
	fmt.Println(*(*byte)(unsafe.Pointer(&bs[1])), unsafe.Pointer(&bs[0]) == unsafe.Pointer(unsafe.SliceData(bs)))
	fmt.Println(reflect.DeepEqual([]*int{&x}, []*int{&x}), fmt.Sprintf("%p", &x) == fmt.Sprintf("%p", (*MyInt)(unsafe.Pointer(&x))))
	type st struct{ a, b int }
	vv := st{1, 2}
	fmt.Println(fmt.Sprintf("%v %+v", &vv, &vv))
	sw := reflect.Swapper([]st{{1, 2}, {3, 4}})
	_ = sw
	arr := []st{{1, 2}, {3, 4}}
	reflect.Swapper(arr)(0, 1)
	fmt.Println(arr)
}
