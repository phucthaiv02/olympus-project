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
https://github.com/gopherjs/gopherjs/issues/358
```

Draft: `kh78nr222wf7f785y260bfyke98g9fqp`

# Challenge Details

**Title**

```
Assignments, calls and composite literals don't follow Go's order of evaluation
```

**Language:** Go

**Category:** Bug Fix

**Description**

```
Fix the order in which GopherJS evaluates assignments, calls and composite literals so that programs behave as the Go specification's order of evaluation requires. Programs that depend on this order currently produce different results, or call functions in a different order, than they do with the gc compiler.

Assignments must happen in two phases. First, the operands of index expressions and pointer indirections on the left-hand side (including the implicit indirection in a selector such as `p.f`) and the expressions on the right-hand side are evaluated, with all function calls, method calls and channel receives happening in lexical left-to-right order. Only then are the assignments carried out, left to right. As a result `i, xs[i] = 1, 9` stores into the element chosen by the old value of `i`, `p, p.i = new(T), 4` sets the field of the struct `p` pointed to before the assignment, and `k, m[k] = k2, v` uses the old key. This applies to assignments with several values, to assignments from a multi-value call such as `i, xs[i] = f()`, and to range loops that assign to existing variables, as in `for i, xs[i] = range s`. In a map element assignment the map operand is evaluated before the key. A run-time panic caused by the assignment itself, such as an index out of range or an assignment to an entry in a nil map, only happens after the right-hand side has been evaluated, and assignments that come before the failing one in the same statement still take effect.

In a function or method call the function value, including any calls needed to compute it or its receiver, is evaluated before the arguments. This must also hold when the arguments are the results of a multi-value call, as in `f()(g())` or `obj().m(g())`, and when an argument itself contains such a call. The elements of composite literals, including struct literals with field names and array or slice literals with explicit indices, are evaluated in the order they are written rather than in field or index order.

All of this must also hold when some of the calls involved block, for example on a channel receive. Panic messages and output formatting are not part of this change.
```

# Research notes

- Issues, all open with no PR: #358 (multi-assignment LHS order; neelance: "the solution is not straightforward"), #837 (`x[f()] = g()`, NeedsFix/NeedsHelp), #1271 (function expression before tuple args). #1063 was closed as a duplicate of #358. #608 (return order) is unspecified in the spec and left out.
- Duplicate search: "order of evaluation", "evaluation order", "evaluated before", "issue23188". No PRs found.
- Probed 25 cases against native Go at HEAD; 8 differ. Multi-assign LHS operands are evaluated after the RHS. Map assignment evaluates the key before the map. Tuple args are evaluated before the function value. `i, xs[i] = ...`, range assignment to `xs[i]`, keyed struct literal order, and panics before RHS evaluation (OOB, nil map) also differ.
- Solution: two-phase assignment via `evaluateAssignOperands` in statements.go (statement, tuple, range and single assignments with side-effecting RHS). Map assign evaluates the map first. `translateCall` hoists the callee when args emit statements (`isHoisted` in utils.go, also used for preserveOrder). Composite literals pre-evaluate keyed elements in source order. Removed fixedbugs/issue23188.go and issue50672.go from gorepo knownFails, since both now pass.
- Tests: `tests/evalorder/evalorder_test.go`, 9 tests. All pass on native Go. They cover blocking (channel) variants too.
- Verified clean clone --network none (after precheck fixes, 2026-10-19): base 2508 pass with/without solution; new 9/9 fail -> 9/9 pass.
