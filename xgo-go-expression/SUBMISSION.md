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
https://github.com/goplus/xgo/issues/2558
```

Draft: `kh76r4rqt81rye1rma8q5rza198g697d` (reused from the go-lua debug library task, which the Scope Gate dropped as publicly-solved on 2026-10-20)

# Challenge Details

**Title**

```
Allow go to be used as an expression that yields a result channel
```

**Language:** Go

**Category:** Feature Request

**Description**

```
Allow `go` to be used as an expression. In expression position, `go f(args)` launches the call in a new goroutine like the `go` statement, and its value is a channel `chan R` of capacity 1. Exactly one value is sent on this channel: the result of the call when it returns. The function value and arguments, including the method receiver, are evaluated in the current goroutine at the time of evaluating the expression, as with the `go` statement. If the function has one result, then R is the type of that result. If there are multiple results, then R is a tuple `(T1, ..., Tn)` containing the results in order. If there is no result then R is `struct{}` and the value is sent after the function returns. The statement `go f(args)` retains its original meaning.

The operand must be a function or method call, including generic functions, variadic calls with `...`, and function literals. Calling a built-in function or a type conversion is a compilation error, and a missing call (`go f`, `go 1`) is a syntax error. `go` attaches to the call immediately following it as a unary operator, so `<-go f()` waits and receives the result, and the `go` expression can be used anywhere an expression is accepted, including in the `case` of `select`. For a function with no result, receiving directly from the `go` expression can only stand as a statement (`<-go f()` or `case <-go f():`), and using the received value is a compilation error.

When a receive `<-ch` from a channel whose element type is a tuple with n fields (n of 2 or more) must provide exactly n values, the fields of the received tuple are unpacked in order, for example `a, b, err := <-go split(s)`. This applies to assignments, short variable declarations, `var` declarations and `select` cases, and takes precedence over the comma-ok form when n is 2. A receive case in `select` can therefore have more than two variables on the left-hand side. In other cases, the received value is still a tuple.

`xgo fmt` preserves the `go` expression as written.
```

# Research notes

- Platform cap: goplus/xgo 3/50 (the virtual fields task is on the same repo, so it will become 4-5/50). Apache-2.0.
- Issue #2558 "Proposal: `go` Expression Syntax" (cpunion, 2026-01-24), labels Proposal-Accepted + FullSpec, no PR in any state. It replaced the closed #2534 (await proposal).
- Duplicate search: "go expression", "await", "future", "GoExpr", "async", "goroutine result" over PRs/issues/discussions, plus `gh search code GoExpr --owner goplus`. Nothing implements it.
- How the task narrows and sharpens the issue:
  - The issue's lowering evaluates the arguments inside the new goroutine. The task requires go-statement semantics instead: the function value, receiver and arguments are evaluated by the caller.
  - Multiple results use XGo's existing tuple type `(T1, ..., Tn)` instead of a hidden `_gop_ret` struct, so `chan (int, string, error)` can be written by users.
  - Destructuring is defined by type: receiving from a tuple channel into n targets (n >= 2) unpacks it, also in `var` declarations and `select` cases, and it wins over comma-ok when n == 2.
  - The issue's optimizations and its "fut received only once" restriction are left out.
- Layers: ast (GoExpr node, Walk), parser (unary `go`, select receive cases with more than 2 targets), printer, x/format, cl (lowering, tuple receive, void check, select), cl/recorder.
- Solution (`cl/goexpr.go` and hooks):
  - When the callee is a plain function value (variables, methods, func literals, explicitly typed values), lower to `func(f F, a0 A0, ...) chan R {...}(f, args...)`.
  - For declared functions (overloaded, generic with inference, pkg.Func), compile the call inside the goroutine. Non-trivial arguments are hoisted into closure parameters; constants and lambdas stay inline.
- Tests: `cl/goexpr_39ed2d_test.go`, 15 tests. Run tests compile, `go run` and compare stdout. They cover:
  - concurrency (a synchronous implementation deadlocks) and the buffer;
  - argument, receiver and func-value evaluation time;
  - call forms, use inside expressions, and select;
  - tuple receive forms, including closed channels;
  - void / builtin / conversion errors, syntax errors, the unchanged go statement, and formatter round-trips.
  Error tests check only the line, and each one has a positive part, so all 15 fail on base. Stable over -count=5 and with GOMAXPROCS=1.
- test.sh and Dockerfile are the same as for xgo-virtual-fields (hash 39ed2d). TestTestxgo/domaintext-md is skipped in base because it needs the network.
- Verified in a clean clone with `--network none`: base passes with and without the solution; new 0/15 -> 15/15. The full `go test ./...` passes locally with the solution.
