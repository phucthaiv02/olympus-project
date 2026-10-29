# Repository

**GitHub repo**

```
https://github.com/goplus/xgo
```

**Commit SHA**

```
c54be7d95fc232e100f5f5f8b72221891c5309ba
```

**GitHub issue**

```
https://github.com/goplus/xgo/issues/2860
```

Draft: `kh78a7rc9f9gpet5y0wecydq7h8ga4yp` (reused from the otto strict mode task, which the Scope Gate dropped as a duplicate on 2026-10-20)

# Challenge Details

**Title**

```
Add virtual fields backed by XGof_ methods
```

**Language:** Go

**Category:** Feature Request

**Description**

```
Add virtual fields to XGo. A virtual field is written like a regular struct field at the point of use (`x.f`, `x.f = v`), but behind the scenes are conventionally named methods instead of actual memory. The type declares the virtual field `f` in one of two forms. The reference form is a method `XGof_ref_f() *T`, and all read and write operations go through the pointer it returns. The getter/setter form consists of `XGof_get_f() T` for reading and `XGof_set_f(v T)` for writing. Without a setter, the field can only be read, without a getter, the field can only be written. The field name is the entire part after `XGof_ref_`, `XGof_get_` or `XGof_set_`, preserving case, and the first segment always determines the form. `XGof_` methods are still normal methods and can still be called directly.

A virtual field can be used anywhere a field of type `T` can be used: read, assign, multiple assign, combined assign (`x.f op= v`), `++` and `--`. With the getter/setter form, the combined assignment and `++`/`--` read with the getter and then write the result with the setter. In these operations the `x` operand is evaluated only once, just like with a real field. In a multiple-value assignment with virtual fields, all values on the right are calculated before assignment, and then the targets are assigned from left to right. `&x.f` is only valid for the reference form and yields the pointer that the method returns.

Virtual fields are in the same namespace as fields and methods, and are promoted through embedded structs according to Go's depth rules. The member with the smallest depth wins. A virtual field and a field or method with the same name at the same depth make the selector ambiguous, as does a virtual field whose methods are ambiguous at that depth. A virtual field consists of its `XGof_` methods found at the smallest depth. Go's receiver rule still applies: if the method has a pointer receiver then `x` must be addressable or already a pointer.

Ambiguous selectors and the following cases are compilation errors, reported on the line where the field is used: the getter returns a type different from the setter's parameter type; the same field has both a reference and a getter or setter; assignment, combined assignment or `++`/`--` to a read-only field; reading a write-only field, including via combined assignment or `++`/`--`; `&x.f` for the getter/setter form; using a virtual field on a non-addressable value when its method has a pointer receiver; using a virtual field as a key in struct literals.
```

# Research notes

- Platform cap: goplus/xgo 3/50 used, not reserved. Apache-2.0, 9.5k stars, last commit 2026-10-18.
- Issue #2860 "Proposal: Virtual Fields" (xushiwei, 2026-09-21), labels Proposal-Accepted + FullSpec, no PR in any state. The only cross-reference is goplus/llcppg PR #937 (closed), which is the binding generator side, not the compiler.
- Duplicate search (dup_search.sh): "virtual field", "XGof", "bit-field", "bitfield", "union member", "getter setter", "property". Only hit: PR #2548 (merged), where "virtual fields" means gogen's named tuple fields (`tupleFields`), an unrelated mechanism. No `XGof_` anywhere in xgo or gogen v1.26.3.
- Avoided the issues that sank the old drafts: no reference implementation exists to port (unlike golua's ldebug.c), and the idea is recent with a fresh repo (unlike otto strict mode).
- Description: Vietnamese source in `description_vi.md`, translated with gtx and then lightly fixed (paragraph -> segment, format -> form, zero-width spaces removed).
- Solution (3 files, all in `cl/`; gogen is untouched):
  - vfield.go: lookup with go/types `LookupFieldOrMethod` (depth, ambiguity), form and type checks, read/ref/address, compound and ++/-- through a single-evaluation temp (`_xgo_vf_recv`), tuple assignment through typed temps (values first, then targets from left to right).
  - expr.go: hooks in compileSelectorExpr (read), compileSelectorExprLHS (reference form as lvalue), compileUnaryExpr (`&x.f`).
  - stmt.go: hooks in compileAssignStmt and compileIncDecStmt for the getter/setter form.
- Tests: `cl/vfield_a46d06_test.go`, 27 tests. Run tests compile the XGo source, `go run` the generated Go and compare stdout, so any lowering is accepted. Error tests only check that compilation fails with an error on the marked line. Every error test also has a positive check, so all 27 fail on base.
- test.sh: base = `./...` skipping `_a46d06$` and `TestTestxgo/domaintext-md` (it imports github.com/xushiwei/markdown, which go.mod doesn't require, so it needs the network); new = `-run '_a46d06$' ./cl/`.
- Dockerfile: `go build ./...` fails on demo dirs that only have .xgo files, so it builds cl/cmd/parser/printer and warms the cl test build.
- Verified in a clean clone with `--network none`: base passes with and without the solution; new 0/27 -> 27/27. The full `go test ./...` passes locally with the solution.
