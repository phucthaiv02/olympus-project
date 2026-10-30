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
https://github.com/goplus/xgo/issues/1470
```

Draft: `kh78a7rc9f9gpet5y0wecydq7h8ga4yp` (reused on 2026-10-21 from xgo virtual fields, which the plagiarism review marked "derivative" of an older submission). No platform checks were run.

# Challenge Details

**Title**

```
Allow spreading slices into slice literals with ...
```

**Language:** Go

**Category:** Feature Request

**Description**

```
Allow spreading a slice into an XGo slice literal by writing `...` after an element, for example `[a..., b, c...]`. Each element with `...` contributes all its elements in order, and each other element contributes itself. The element being spread must be of type slice or array, including named types. Spreading a value of another type, such as number, string or map, is a compile error.

When the context expects a slice type, such as a typed variable declaration, an assignment, a function argument, a return value, or the fields and elements of a composite literal, the result is of that type. Then every element, including each element of the spread slices, must be assignable to its element type, otherwise there is a compile error. When the expected type is not present, the result has the type that a regular slice literal would have if each spread element were replaced by a value of its element type. So for `a := [1, 2, 3]` and `b := [4, 5]`, `[a..., b...]` is `[]int` with value `[1 2 3 4 5]`, and `[a..., b]` is `[]any` with value `[1 2 3 [4 5]]`.

Each element is evaluated exactly once, in order from left to right. The result is always a new, non-nil slice that does not share memory with any slice being spread, including `[a...]` or when all slices being spread are empty. Slice literals with spread elements can be used wherever a slice literal is accepted, and `xgo fmt` keeps them as written.
```

# Research notes

- Issue #1470 "slice compose" (2023, labels proposal and Proposal-Accepted) has no PR. Duplicate searches covered "slice compose", "spread", "SliceLit ellipsis", "a...", the open PR list and gogen. Nothing implements it. The only related code is `ast.ElemEllipsis`, which is used for MatrixLit rows; MatrixLit compiling is commented out in cl.
- At base, `[a..., b]` is a syntax error ("expected ']'"), and `[x, a...]` parses but cl fails with "unknown - *ast.ElemEllipsis".
- The platform has 3/50 other submissions on goplus/xgo, plus 3 of ours (go expression, virtual fields, now this one).
- Solution (parser + cl):
  - The parser accepts `...` after the first element.
  - `compileSliceLitWithSpread` compiles each element once and works out the result type: the expected slice type when there is one, otherwise gogen's own SliceLit inference with each spread replaced by a value of its element type.
  - It then emits a closure that takes the operands as arguments, which gives left-to-right evaluation exactly once. The closure body is `make` with the exact capacity, then `append(..., x...)` when element types are identical and a range loop otherwise (conversion into `[]any`, or arrays).
  - The printer already prints ElemEllipsis.
- Tests: `cl/slicespread_da8d69_test.go`, 9 tests. Seven compile the program, run it and compare output; one checks compile errors on the marked line; one checks the format round trip. They cover same type, inference (int8/float32/named types/mixed → `[]any`), expected types (typed var, named slice, func arg, return, struct field, map value, nested literal, assignment), arrays, single left-to-right evaluation, new non-nil non-aliasing results, use inside comprehension / for-in / variadic / nested spread, and errors (int, string and map operands; incompatible element types in typed var, func arg and `[]int`).
- Dockerfile/test.sh are the same as in xgo-go-expression: no `go test` during the build, and base skips TestTestxgo/domaintext-md, which needs the network.
- Verified in a clean clone with `--network none`: base 1165 passes without the solution and 1166 with it; new 9/9 fail -> 9/9 pass.
- Risk: the solution is small (about 140 lines in 2 files), so Nova may solve it. The hidden depth is in the inference rules, the typed-context conversion and the parser fix for the first element.
