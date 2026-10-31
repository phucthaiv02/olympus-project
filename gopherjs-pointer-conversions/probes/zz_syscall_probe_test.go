package compiler

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/gopherjs/gopherjs/compiler/sources"
	"go/types"
)

func TestZZSyscallProbe(t *testing.T) {
	src := `package syscall
import "unsafe"
type Stat struct{ A, B int32 }
func call(p unsafe.Pointer) {}
func ToUnsafe(s *Stat) { call(unsafe.Pointer(s)) }
func FromUnsafe(p unsafe.Pointer) *Stat { return (*Stat)(p) }
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "x.go", src, 0)
	if err != nil {
		t.Fatal(err)
	}
	srcs := &sources.Sources{ImportPath: "syscall", Files: []*ast.File{f}, FileSet: fset}
	importer := func(path, dir string) (*sources.Sources, error) { t.Fatal(path); return nil, nil }
	tc := types.NewContext()
	PrepareAllSources([]*sources.Sources{srcs}, importer, tc)
	a, err := Compile(srcs, tc, false)
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range a.Declarations {
		t.Log(string(d.FuncDeclCode))
	}
}
