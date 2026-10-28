# Repository

**GitHub repo**

```
https://github.com/robertkrimen/otto
```

**Commit SHA**

```
f5fb48544f35425b37d4399b4ad00e3319ea1f5f
```

**GitHub issue**

```
https://github.com/robertkrimen/otto/issues/492
```

Draft: `kh78a7rc9f9gpet5y0wecydq7h8ga4yp` (pushed 2026-10-20, no checks run)

# Challenge Details

**Title**

```
"use strict" has no effect: implement ES5 strict mode
```

**Language:** Go

**Category:** Feature Request

**Description**

```
Implement ECMAScript 5 strict mode. Code that starts with a `"use strict"` directive must follow the strict mode semantics of ES5 (Annex C) instead of running as ordinary code. Code that isn't strict keeps its current behavior.

Code is strict when its directive prologue, the sequence of plain string literal statements at the start of a script, function body or eval code, contains `"use strict"` or `'use strict'` written exactly so (no escape sequences, no parentheses). Strictness applies to every function nested inside strict code. Functions created by `Function(...)` are strict only if their own body has the directive. The code of a direct `eval` call made from strict code is strict.

In a strict function `this` is used exactly as passed. `undefined` and `null` are not replaced by the global object, and primitive values are not wrapped. This covers plain calls, `call`, `apply`, `bind`, callbacks of built-in functions, calls from Go, and methods or accessors invoked on primitive values. A strict function's arguments object is not linked to its parameters. Reading or writing its `callee` throws a TypeError, as does accessing the `caller` or `arguments` property of a strict function itself.

Strict code throws instead of failing silently. Assigning to an undeclared identifier throws a ReferenceError. An assignment that can't be carried out throws a TypeError. This includes non-writable properties, accessors without a setter, non-extensible objects, read-only globals such as `NaN`, and properties of primitive values. An inherited setter still receives such an assignment, with the primitive value as its `this`. Assigning to the name of a function expression from inside that function also throws a TypeError, and so does deleting a non-configurable property. Strict eval code gets its own variable environment, so its declarations don't leak into the calling code or the global scope, whether the eval is direct or indirect.

A SyntaxError is reported for strict code that contains a `with` statement, an octal numeric literal, or an octal escape sequence in a string (a `\0` not followed by a digit stays allowed). The same applies to duplicate parameter names, to `eval` or `arguments` used as the name of a variable, function, parameter or catch parameter or as the target of an assignment, `++` or `--`, to `delete` applied to a plain identifier, and to a strict mode reserved word (`implements`, `interface`, `let`, `package`, `private`, `protected`, `public`, `static`, `yield`) used as an identifier. These rules also apply to the name and parameters of a function whose own body contains the directive.
```

# Research notes

- Platform cap: robertkrimen/otto 6/50 used (44 left, just under the ">45 left" preference). MIT, 8.4k stars, active (last commit 2026-10-08).
- Issue #492 (enhancement, help wanted), open, with no PR. Maintainer stevenh: "Happy for someone to put in a PR that addresses this current short coming." The README caveat says "use strict will parse, but does nothing".
- Duplicate search: "strict mode", "use strict", "strict", "this undefined", "TypeError frozen". Only #492 matches. Open PR #556 (Object.assign TypeError) is unrelated.
- Risk: this is the most visible gap in otto (README caveat), and 6 other submissions exist on the repo, so the plagiarism check may find an older strict-mode task.
- Found by differential testing against node: 38/40 strict-mode probes differed at base. After the solution, the only divergences are a few ES2015+ items that otto doesn't target.
- Solution (18 files):
  - parser: directive prologue detection, inherited strictness, the `StrictMode` parse mode for eval, and the early SyntaxErrors.
  - ast: `Strict` on Program and FunctionLiteral.
  - cmpl: strict flags.
  - runtime: `scope.this` is now a Value and coercion happens only for non-strict code; implicit this is undefined for calls by name (JS functions only); primitive base references (getters, setters and methods get the primitive); ReferenceError on unresolvable assignment; TypeError through the [[Put]]/[[Delete]] throw flag; immutable function-expression name binding; unmapped arguments with a poisoned callee; poisoned caller/arguments on strict functions.
  - eval: strict eval gets its own declarative environment.
  - builtins: apply/call/bind pass thisArg through to JS functions.
  - docs: README/otto.go caveat removed.
- Native functions keep getting an object as this (wrapper or global), so the existing builtins and tests are unchanged.
- Tests: `strict_mode_d193ed_test.go`, 10 tests. They check behavior and error names only, not messages; non-strict behavior is checked alongside. arguments.caller is not tested (ES5 vs modern engines differ).
- test.sh: `go test -vet=off` (the vet "tests" analyzer rejects the existing documentation_test.go on current Go). base = packages `.`, `./parser`, `./ast` minus the new tests.
- Verified in a clean clone with `--network none`: base 588 pass with and without the solution; new 10/10 fail -> 10/10 pass.
