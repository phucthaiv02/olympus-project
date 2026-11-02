# Repository

**GitHub repo**

```
https://github.com/goplus/xgo
```

**Commit SHA**

```
26e013d1fc7171b01d32bdb4fdda3aab5a3bcd25
```

**GitHub issue**

```
https://github.com/goplus/xgo/issues/859
```

# Challenge Details

**Title**

```
Add index operator methods so user-defined types support m[i] and m[i, j]
```

**Language:** Go

**Category:** Feature Request

**Description**

```
Add index operator methods, which let a type define what indexing its values means. A method declared as `func (m *Matrix) [](i, j int) float64` is the index getter of `Matrix`: an index expression `m[i, j]` on a value of that type calls it with the indices as arguments, and the expression has the method's result as its value. A method declared as `func (m *Matrix) []=(i, j int, v float64)` is the index setter: an assignment `m[i, j] = v` calls it with the indices followed by the assigned value. The number of indices is whatever the method declares, so a single index `v[i]` and three indices `m[i, j, k]` work the same way, and a type may define a getter, a setter, or both. A getter must have at least one parameter and exactly one result, and a setter must have at least two parameters and no results. Any other signature is a compile error.

In the generated Go code the getter and setter are methods named `XGo_Index` and `XGo_IndexSet`, following the naming of the other operator methods. Any type that has methods with these names can be indexed the same way, whether they were declared with the operator syntax or under those names.

Index operators are found the same way as ordinary methods. They are promoted through embedded fields, a getter or setter with a pointer receiver can be used on an addressable value, and one with a value receiver can be used through a pointer. The indices and the assigned value are converted to the parameter types exactly as the arguments of a method call are, so untyped constants take the parameter type, and the wrong number of indices or a value of the wrong type is a compile error. Reading an index expression when the type has no getter, assigning to one when it has no setter, and taking its address are compile errors. Using several indices on built-in slices, arrays, maps and strings remains an error, and instantiating a generic function with several type arguments, as in `f[int, string]`, keeps working.

Index operators can be used anywhere an index expression can appear in an expression or as an assignment target, including chains such as `grid[i][j]`, where each step uses whatever indexing its operand supports. This includes the init statements of `if`, `for` and `switch` statements, the post statement of a `for` loop, and the iteration variables of a `for ... = range` loop, which are assigned in each iteration like the targets of one assignment. In a compound assignment such as `m[i, j] += v` and in an increment or decrement such as `m[i, j]++`, the operand and the indices are evaluated once. Then the getter is called, the right-hand side is evaluated, and the setter is called with the same indices and the combined value. In an assignment with several targets, the operands and indices of all targets and all right-hand side values are evaluated first, in the usual left-to-right order. Only then are the setters and other assignments carried out, from left to right, so `m[0], m[1] = m[1], m[0]` swaps the two elements. The same holds when the values come from a multi-value function call.

Formatting source that contains getter and setter declarations must produce code that parses back to the same declarations.
```

# Research notes

- Draft ID reused: `kh78nr222wf7f785y260bfyke98g9fqp`. The previous idea (xgo static members, #2513) was dropped by the Scope Gate on 2026-10-22 as publicly-solved: earlier fork commits of PR #2792 (KurodaKayn/xgo d6b9c7e, 54b5e13) held the full cl implementation and covered 5/7 tests.
- Issue #859 "Index Operator Methods" (2021, labels proposal and Proposal-Accepted, 4 reactions). No PR references it in any state. In #764 the maintainer (xushiwei) wrote the syntax himself: `func (m *Matrix[T]) [] (i, j int) T` (getter, `m[i, j]` → `m.[](i, j)`) and `func (m *Matrix[T]) []= (i, j int, x T)` (setter).
- The `XGo_Index` name follows `dql/README.md` (`func (ns NodeSet) XGo_Index(index int) NodeSet // ns[n]`, documented but not implemented anywhere). Nothing defines or uses the `XGo_IndexSet` name yet.
- Duplicate searches on goplus/xgo (PRs in all states, issues, discussions): "index operator", "operator method", "subscript", "XGo_Index", "Gop_Index", "multi index", "m[i, j]", "indexer", "matrix", "overload index", "index method", "getter setter", "IndexListExpr", "[]=". Only #859/#764 are related. #2707 (open, aofei) adds `NodeSet.XGo_Index` inside the dql/maps runtime evaluator; it does not touch cl indexing. On goplus/gogen, "index operator", "XGo_Index", "IndexSet", "operator method" and "Index" found nothing related. A GitHub code search for `XGo_IndexSet`, `XGo_Index(` and `Gop_IndexSet` found only dql/README.md.
- Platform: goplus/xgo has 3/50 global submissions, not reserved.
- Base behavior (checked at 26e013d1):
  - `func (m *Matrix) [](i, j int) float64` → parse error "expected '{', found float64".
  - `x := m[0, 1]` → parse error "expected type, found 1" (the parser only accepts types after the first comma).
  - `m[i, j]` with identifiers → cl error "Index doesn't support a[i, j...] yet" (from gogen).
  - `m[0]` with an `XGo_Index` method → "type *Matrix does not support indexing".
  - With operator methods on named types, the method wins over the built-in operator (`a*10+b` inside `foo.+` calls XGo_Add).
  - Generic type declarations (`type G[T any]`) don't parse in .xgo, so generic receivers stay out of scope.
- Expected layers: parser (operator method names `[]` and `[]=` after the receiver; expressions after the comma in an index list), printer (round trip of the declarations), cl declaration (name mapping and signature checks), cl expressions (IndexExpr/IndexListExpr on values with a getter, but not on generic functions), and cl statements (assignment, op-assign, inc/dec, multi-assign two-phase order, tuple from a multi-value call). gogen cannot be modified, so it must be lowered to method calls in xgo itself.
- Behaviors tests should cover:
  - getter and setter with 1, 2 and 3 indices;
  - untyped constant conversion of indices and value;
  - methods declared by name `XGo_Index`/`XGo_IndexSet`;
  - promotion through embedded fields;
  - pointer receiver on an addressable var, value receiver through a pointer;
  - chain `grid[i][j]` (built-in then operator, and the reverse);
  - `+=`, `++`, `--` with side-effecting operand, indices and RHS (each evaluated once; order getter → RHS → setter);
  - swap `m[0], m[1] = m[1], m[0]`, `i, m[i] = ...` with the old i, and a multi-value call RHS;
  - errors: bad getter/setter signatures, read with no getter, assign with no setter, `&m[i]`, wrong index count/type, multi-index on a built-in slice/map;
  - `f[int, string]` still works;
  - format round trip.
- Not specified, so don't test: precedence when the underlying type already supports indexing, variadic getters, classfile shorthand, range-loop targets, interfaces with index methods.
- Risks: the old popular proposal may already exist as an older platform submission (but xgo has only 3 global submissions). Multi-index parsing must not break generic instantiation. Size looks moderate to large (parser + printer + three areas of cl).
- Re-verified on 2026-10-24, because the draft still holds the static-members content that the Scope Gate rejected:
  - goplus/xgo HEAD is still 26e013d1, the pinned commit.
  - #859 has no referencing PRs, and no PR on xgo has been updated since 2026-10-15 except dependency bumps and #2005 (unrelated).
  - Duplicate searches were re-run and are still clean. On xgo: "index operator", "operator method", "XGo_Index", "IndexSet", "subscript", "multi index", "IndexListExpr", "matrix". On gogen: "XGo_Index", "IndexSet", "index operator".
  - None of PR #2707's 3 commits (2a2e0c32, e3522685, aa1c13f4) touches cl indexing. In cl it only changes nodeInterp/outline instantiate; its `NodeSet.XGo_Index` is a runtime method in dql/maps.
  - Platform: goplus/xgo is at 3/50 global, 0/12 own, not reserved.
- Alternatives considered on 2026-10-24 and rejected:
  - xgo #2834 (Auto Lambda control flow): the maintainer disabled it on purpose until the semantics are settled.
  - xgo #2441 (Untyped Literal Extension): vague, multi-phase, and not implementable as one task.
  - xgo #2828, #2559: already implemented.
  - Other Go interpreter repos are saturated: expr, risor and mvdan/sh are at 50/50, tengo is reserved, goja is at 44/50 and starlark-go at 28/50, so the risk of overlapping an older platform submission is high.

## Build (2026-10-25)

Files in this folder: `solution.patch`, `test.patch` (includes `test.sh`), `Dockerfile`, `test.sh`. The working clone is in `work/` and the probe programs are in `work/_probe/`.

- Solution, about 700 lines across parser + cl; gogen is unchanged:
  - Parser:
    - `parseIndexOpName` reads `[]`/`[]=` after the receiver. `func (params) []T {...}()` is still a function literal.
    - In an index list, the elements after the first comma are parsed with `parseRHS` instead of `parseType`, so `m[0, 1]` parses. `f[int, string]` still works.
    - The printer already prints `func (m *M) [] (...)`, which parses back.
  - `cl/compile.go`, `loadFunc`: `[]`→`XGo_Index` and `[]=`→`XGo_IndexSet`. `checkIndexOpSig` checks the signature: the getter needs ≥1 parameter and 1 result, the setter needs ≥2 parameters and 0 results.
  - `cl/index_op.go`:
    - `indexOp` uses `types.LookupFieldOrMethod` after `loadMethods`. Probing with `cb.Member` must not be used: on `any` it emitted map-field-access statements and broke `TestTestxgo/map-field-access2`.
    - The getter is lowered to `x.XGo_Index(i...)`. `&m[i]` is an error.
    - Single assignment: `x.XGo_IndexSet(i..., v)`.
    - op-assign / inc / dec is lowered to a block: the operand (`&x` when the getter or setter has a pointer receiver) and the typed indices go into temporaries, then `_xgo_v := get(...)`, `_xgo_v op= rhs`, `set(..., _xgo_v)`.
    - Multiple targets: one block evaluates the operands and indices of every target (built-in index/map/array via `&`, `p.f` via `&p`/`p`, `*p`) and then every RHS into typed temporaries, from left to right, and then assigns in order. It also handles the multi-value call case.
    - In the init/post of if/for/switch the block is wrapped in a closure `func(){...}()`, via `ctx.simpleStmt`, which `compileSimpleStmt` sets.
    - `for k, a[i] = range x` is rewritten to `for _xgo_k, _xgo_v := range x { k, a[i] = _xgo_k, _xgo_v; ... }`.
  - Multiple indices on slices, maps and strings give the error "invalid operation: more than one index", with its position.
- Tests: `cl/indexop_77753e_test.go`, 15 tests.
  - Run tests: getter/setter with 1–3 indices, generated method names, methods named `XGo_Index` and through an interface, receivers and promotion, conversions, chains, op-assign order, op-assign with a value-receiver getter, multiple assignment, multiple-assignment order (with built-in, field, `*p` targets and the old `p`), statements (for post/init, if/switch init, range), generic instantiation, format round trip.
  - Error tests: 5 bad signatures and 19 other errors.
- Verified with Docker `--network none`, `ulimit -s unlimited`, on a clean checkout at 26e013d1, with the test.patch as uploaded:
  - Without the solution: base 1165 pass, new 15/15 fail.
  - With the solution: base 1166 pass, new 15/15 pass.
  - The extra base test is `cl/internal/gop-in-go/foo TestReverseMap`, because the first run in the same container generates `xgo_autogen.go`; it is not caused by the solution.
  - The junit for base contains no `77753e` names.
  - The 2 duplicate names `TestErrAssign`/`TestErrStringLit` already exist in the repo, and both pass.
- Description: added the sentence about init/post/range, which the tests check. 552 words.
- Not specified and not tested: precedence when the underlying type also supports indexing (the method wins), pointer receivers on non-addressable values (same as ordinary methods: the generated Go fails), classfiles, and variadic getters/setters.
