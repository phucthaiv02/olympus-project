# Repository

**GitHub repo**

```
https://github.com/gopherjs/gopherjs
```

**Commit SHA**

```
490705b1d6fc7d5bd9202ac41888e146183328eb
```

**GitHub issue**

```
https://github.com/gopherjs/gopherjs/issues/1240
```

Draft: `kh7cnm6syf87kby102esrpc5858gbdzb` (uploaded 2026-10-19, no checks run)

# Challenge Details

**Title**

```
Pointer conversions lose the destination type, nil and pointer identity
```

**Language:** Go

**Category:** Bug Fix

**Description**

```
Fix conversions between pointer types so that the result behaves like a pointer of the destination type, the way it does in Go. This covers explicit conversions such as `P(p)` or `(*MyInt)(p)`, and the implicit conversions that happen wherever a `*T` value is used where an assignable named pointer type such as `type P *T` is expected, or the other way around.

A converted pointer has the destination type, as seen by `%T`, reflection, type assertions and type switches. It is nil exactly when the source is nil, and it refers to the same variable, so writes through either pointer are visible through the other. Converting a pointer to another type and back gives a pointer equal to the original, and converting the same pointer to the same type twice gives equal pointers, so converted pointers also work as map keys. This applies to pointers to basic types, including element types that differ only by name such as `*int` and `*MyInt`, to pointers to pointers and to pointers to arrays, including pointers to slice elements and struct fields.

`unsafe.Pointer` follows the same rule for nil. Converting the zero `unsafe.Pointer` to any pointer type gives that type's nil pointer, and converting a nil pointer of any type to `unsafe.Pointer` gives a value equal to nil. A non-nil pointer that goes through `unsafe.Pointer` and back to its own type keeps its identity. When it is converted to a different pointer type whose element type has the same representation, such as from `*int64` to `*MyInt64` or from `**int` to `**MyInt`, the result has the destination type and refers to the same variable, and values read through it have the destination element type. Conversions done with `reflect.Value.Convert` follow the same rules as conversions in code.
```

# Research notes

- Platform cap: gopherjs/gopherjs 4/50 used (46 left), yours 0/12. License BSD-2-Clause, 13k stars, last commit 2026-07-27.
- Related issues, all open with no PR: #1240 (named pointer in struct literal), #1241 (pointer-to-array conversion loses nil), #1060 (unsafe.Pointer nil). #843 was fixed by #1435 (nil struct-pointer casts), so struct pointers are already correct. They serve as the regression control.
- Duplicate search: "pointer conversion", "named pointer", "pointer type nil", "translateImplicitConversion", "ptrType", "conversion nil". #1199, #1208 and #1435 were merged and fixed other cases. #669 is an open WIP PR from 2017 (structs/arrays/interfaces copy semantics), a different area.
- Probed at HEAD with go1.21.13 and node 24: about 20 behaviors differ from native Go. Implicit conversions in arguments, returns, literals, map keys, chan sends, append and variadic keep the source type. Nil pointer-to-array stays non-nil. Explicit conversions break identity (`(*string)(P(p)) != p`) and map keys. `unsafe.Pointer` nil is not nil, and `(*Struct)(unsafe nil)` throws a TypeError.
- Solution: in `compiler/expressions.go`, implicit pointer conversions go through `translateConversion`; array, basic and unsafe cases use the new prelude helper `$pointerConversion`, which maps nil and caches proxies for identity. `$pointerToUnsafe` maps nil to 0.
- Tests: `tests/ptrconv/ptrconv_test.go`, 11 tests. They also pass on native Go. The struct-from-unsafe case is last because at base it aborts the test binary.
- test.sh base: Go unit tests (skipping Test_JSSourceMap_*, which is path-dependent and fails at base too) plus `gopherjs test` of ./tests and 13 stdlib packages. Result: 2510 tests, all pass with and without the solution.
- Verified in a clean clone with `--network none`: new is 11/11 fail before the solution and 11/11 pass after.
- Precheck fixes 2026-10-19: tests moved to tests/ptrconv_658733/pointer_conversion_658733_test.go; test.sh runs the compiler via 'go run' wrapper (no go install); description trimmed. Re-verified: base 2510 pass, new 11/11 fail -> pass.
- Quality-check fix 2026-10-21: verifyTests/verifySolution failed because every `gopherjs test` run exited 255 with no output. The platform sandbox runs as user `model` with `ulimit -s unlimited`, so gopherjs passes node a huge `--stack_size`. Fix: the Dockerfile puts a `node` wrapper in /usr/local/bin that uses an 8 MiB stack when the limit is unlimited (this also covers rollout agents). Tests extended per the verifyFairness gaps: named-to-unnamed implicit conversions (TestPtrConvImplicitFromNamed), non-nil `**T` (TestPtrConvPointerToPointer), struct-to-struct pointer regression check, and unsafe.Pointer with `**int`/`*[3]int`. 13 tests now; the old patch is in test.patch.orig. Re-verified as `model`, `ulimit -s unlimited`, `--network none`: base 2512 pass (with and without the solution), new 13/13 fail -> 13/13 pass, and 13/13 pass on native Go. Uploaded test patch + Dockerfile; the solution is unchanged.
- Precheck fixes 2026-10-22: the Dockerfile check failed claiming the Go slim base has no Node; it does (Debian nodejs, /usr/bin/node v24), so the Dockerfile now says so and asserts `/usr/bin/node --version` at build. test.sh comments that `./tests` excludes the tests/ptrconv_658733 subpackage (base/new warning). Description: dropped the "struct pointers must keep working" sentence (HIGH conciseness).
- Assessment 2026-10-22: new 13/13 FAIL x3 at base and 13/13 PASS x5 with the solution (stable). Probes outside the tests (package var, field/array-index assign, closure return, defer arg, multi var decl, **T write) pass with the reference and on native Go; 8 of 11 fail at base. Scope Gate verdict before the edit: pass (high), 0/41 lines aligned. No new upstream PRs. Reference size is 55+/8- in 2 files.
- verifySolution fix 2026-10-22: FAIL p2p_f2p_overlap because base ran `go test $(go list ./...)`, which included tests/ptrconv_658733 natively. Base now excludes that package: 2499 tests pass with and without the solution, 0 ptrconv entries in the base junit; new 13/13 fail -> pass. Uploaded the test patch only. The solutionQuality medium note (non-nil unsafe.Pointer cast to a different pointer type keeps the source type) is left as is: verdict PASS, and the description only promises nil and own-type round trips for unsafe.Pointer.
- verifySolution fix #2 2026-10-22: FAIL before_p2p_not_passing. The tests package ran both natively (go test) and under gopherjs, which gave 205 duplicate junit names; the 11 os/exec integration tests (pass natively, skip under JS) were reported as skipped. Base now runs those 11 natively only (-run) and the rest of ./tests under gopherjs only (--run list built from func Test* in tests/*_test.go). Verified: 0 duplicates, 0 base/new overlap, base 2254 pass + 38 skip-only with and without the solution, new 13 fail -> 13 pass.
- solutionQuality fix 2026-10-22: FAIL with HIGH 'syscall unsafe-to-struct conversion mishandles nil' and MEDIUM 'non-nil unsafe.Pointer to a different pointer type keeps the source type'. Both syscall struct branches now map nil (nil struct pointer -> 0, zero unsafe.Pointer -> ptrType.nil; checked in the generated JS of a package compiled as syscall). The unsafe -> pointer path uses the new prelude helper $unsafeToPointer: a cached destination-typed proxy that also rewraps 64-bit/complex pointees and nested pointers. 18 probes match native Go. No hidden test added. base 2254 pass + 38 skip, new 13/13. The old patch is in solution.patch.v1.
- solutionQuality fix #2 2026-10-22: FAIL with HIGH 'syscall struct-pointer unsafe round trips lose identity'. I reviewed every unsafe.Pointer and pointer conversion path. Fixes: (1) syscall struct ptr -> unsafe: storeStruct pre-fills the byte array and records the source pointer, the reverse returns it ($unsafeStructPointer), and the delayed reload copies array fields in place; (2) the &x[i], &_zero and syscall-only new(T) byte-array forms record their source pointer ($withSourcePointer); new(T) is now syscall-only, so elsewhere unsafe.Pointer(new(T)) is the pointer itself; (3) reflect cvtDirect pointer conversions use the same helpers ($toPointerType, formerly $unsafeToPointer, and $pointerOfStructConversion) for identity and aliasing. 56 probes in probes/ match native Go; probe3 runs the syscall branches by treating main as syscall. base 2254 pass + 38 skip, 0 duplicate names; new 13/13 x2; 16 extra stdlib packages identical with and without the solution (1172 pass). The solution is now 195 insertions / 25 deletions in 3 files. Older patches: solution.patch.v1, .v2.
- Difficulty fix 2026-10-24: Vega passed 3/3 and verifierIncompleteness reported INCOMPLETE (5 gaps). I applied each Vega run's compiler diff and ran 27 isolated probes (probes/probe5.go). Gaps: Vega #1 and #2 fail non-nil unsafe.Pointer to a different pointer type and reflect struct-pointer aliasing, Vega #1 and #3 fail unsafe.Pointer(new(T)) round trips, Vega #1 fails unsafe.Pointer(&s[i]) stored in a variable. Added TestPtrConvUnsafeOtherType, TestPtrConvUnsafeElements (includes evaluate-once), TestPtrConvUnsafeNew and TestPtrConvReflectConvert, each with a panic-to-error guard. The description gained a sentence for different-type unsafe conversions and reflect.Value.Convert. Re-graded locally: base 17/17 fail, reference 17/17 pass (x3), native Go 17/17 pass, Vega #1/#2/#3 all fail now (0/3). Previous versions: test.patch.v3, description.md.v3.

## 2026-10-25: tighten after Nova #5 / Orion #4 / Vega #2 passed

- Probes `probes/probe6.go`..`probe9.go` (185 checks) vs native, ref and the 5 run diffs. Passing runs handled all ordinary edge cases; the differentiators were:
  - unsafe.Pointer equality across pointer types (Vega fails),
  - `reflect.NewAt` (Orion, Vega fail),
  - deep reinterpretation via unsafe: `*[2]int64`->`*[2]ptcInt64`, struct fields with int64/pointer/nested struct/array fields (only Nova #5 handles; Nova #5 fails nested array inside a struct field).
- Reference upgraded: `$pointerProxies` with `$original`, canonical `$pointerToUnsafe`, `$needsReinterpret`/`$reinterpretPointee`/`$reinterpretInto`, `$arrayPointerConversion` (JS Proxy views), positional field mapping in `$pointerOfStructConversion`, `$indexPtr` on views, `&x[i]`/`_zero` byte-array forms only in syscall, `*(*T)(unsafe.Pointer(&x))` shortcut only for identical types; reflect: `cvtDirect` Ptr -> `$toPointerType`, `assignTo` pointer conversion, `Pointer()` canonical, `NewAt` override.
- New tests: UnsafeSameAddress, ReflectNewAt, ReflectAssign, UnsafeArrayElements, UnsafeStructFields (22 total).
- Local: base 22/22 fail, ref 22/22 pass (x3), native 22/22 pass; base mode 2254 pass / 38 skip, 0 dups; 26 extra stdlib packages ok.
- All 8 historical run diffs now fail (Nova #5 fails only UnsafeStructFields; Orion #4 fails 3).
- Backups: test.patch.v4, description.md.v4, solution.patch.v3.

## 2026-10-25 (2): fix descriptionQuality "too many requirements" + solutionQuality HIGH
- Description rewritten as fewer general rules (405 -> 273 words); backup description.md.v5.
- solutionQuality HIGH: explicit pointer-to-array conversion path now also uses `$toPointerType` (array view when needed); `$toPointerType` handles array targets before `$unsafeSource` so syscall byte arrays stay arrays. Backup solution.patch.v4.
- Re-verified: new 22/22, probes unchanged, base 2254 pass/38 skip, 0 dups, 0 overlap.

## 2026-10-25 (vòng 3): solutionQuality + descriptionQuality lại FAIL

- descriptionQuality: bản 273 từ vẫn bị chặn "too many requirements". Đã viết lại còn 194 từ, 3 quy tắc (phạm vi; kiểu đích, nil, cùng biến và bằng nhau; đổi kiểu cùng layout qua unsafe). Bản cũ lưu ở `description.md.v6`.
- solutionQuality (HIGH): `&view.n` qua struct view không bằng `&s.n`. Struct proxy giờ có accessor `$ptr_<field>`, trả về pointer field của struct gốc đã đổi kiểu, và cache trên struct gốc.
- solutionQuality (MEDIUM): nil `*[2]int64` trong mảng đọc qua view vẫn là nil kiểu nguồn. `$needsReinterpret` giờ xử lý nil. Kiểu mảng gốc được ghi lại (`$recordArrayType`, `$pointerToUnsafe(obj, type)`), nên khi ghi nil vào thì vẫn là nil kiểu của biến gốc. `$indexPtr` cũng dùng kiểu gốc này.
- Sửa thêm Y18: phần tử slice đọc qua `*[]MyInt64` giờ có kiểu đích nhờ `$reinterpretSlice`, dùng view của backing array.
- Kết quả: 22/22 test mới pass. Base: 2254 pass, 38 skip, 0 fail. Probe 5–10 pass 219/219 (`probes/probe10.go` mới). Bản solution cũ lưu ở `solution.patch.v5`.

## 2026-10-25 (vòng 4)
- descriptionQuality PASS với bản 194 từ.
- solutionQuality FAIL (HIGH): "same memory layout" bị đọc thành phải đổi bit (`*uint32` → `*float32`). Đã đổi câu thành "differs only in the names of the types it is made of" (204 từ) và sửa 3 comment trong prelude.js. Code không đổi, 22/22 test mới vẫn pass. Bản cũ lưu ở `description.md.v7`.

## 2026-10-26 (vòng 5)
- descriptionQuality PASS.
- solutionQuality FAIL (HIGH): `reflect.Value.Convert` đổi pointer proxy sang `unsafe.Pointer` không chuẩn hoá. Case này không xảy ra được, vì cả Go lẫn GopherJS đều panic khi Convert pointer ↔ `unsafe.Pointer` (probe11 C1–C9). Dù vậy nhánh `UnsafePointer` của `cvtDirect` vẫn đi qua `$pointerToUnsafe`, có comment giải thích.
- Tự rà thêm: `Field(i).Addr()`, `Index(i).Addr()` và `UnsafeAddr` của reflect tạo pointer mới mỗi lần, nên không trùng `&s.f` / `&a[i]`, kể cả khi không đổi kiểu. Giờ chúng dùng `$fieldPtr` (helper mới trong prelude) và `$indexPtr`, và `UnsafeAddr` đi qua `$pointerToUnsafe`.
- Kết quả: 22/22 test mới pass. Base 2254 pass, 38 skip, 0 fail. Probe 5–11 pass hết. Bản solution cũ lưu ở `solution.patch.v6`.

## 2026-10-27 (vòng 6)
- solutionQuality FAIL (HIGH): `$indexPtr` cache pointer phần tử theo vị trí, không theo kiểu. Vì vậy `&(*[2]MyInt32)(unsafe.Pointer(&a))[0]` trả lại `*int32` đã cache. Đã sửa: pointer lấy từ cache mà khác kiểu thì đổi qua `$toPointerType`.
- probe12 (alias của int32/int/byte/float64/string/bool, slice, map key, cả hai thứ tự lấy địa chỉ): bản cũ fail 8/13, bản mới pass 13/13. Probe 5–11 pass. 22/22 test mới pass. Base 2254 pass, 38 skip, 0 fail. Bản cũ lưu ở `solution.patch.v7`.
