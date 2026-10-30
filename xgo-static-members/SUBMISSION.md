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
https://github.com/goplus/xgo/issues/2513
```

Draft: `kh78nr222wf7f785y260bfyke98g9fqp` (reused on 2026-10-21 from gopherjs evaluation order, which the Scope Gate rejected). No platform checks were run.

# Challenge Details

**Title**

```
Implement static constants and variables, and make static methods resolvable
```

**Language:** Go

**Category:** Feature Request

**Description**

```
Implement static members, that is static constants, static variables and static methods, which are declarations scoped to a type. In regular XGo, `const T.name = ...` declares a static constant of `T` and `var T.name V = ...` declares a static variable of type `V` that belongs to `T`. In a classfile there is also the shorthand `.name`, meaning `T` is the type of the class. In the `var` block of a classfile, the names `T.name` and `.name` declare static variables, while names without a prefix are still fields of the class. `T` must be a type defined (not aliased) in the same package, otherwise it is a compile error. A static member cannot have the same name as a method of `T`, otherwise it is a compile error. Different types can have static members with the same name, and static members do not conflict with package-level names.

Static constants behave like regular constants. They can be used in constant expressions, can use `iota` in a `const` group, and can be declared mixed with regular names in the same group or the same spec, for example `const T.A, B = 1, 2`. Assigning a value to a static constant is a compile error. Static variables are package-level variables: they can be assigned, have their address taken, and are initialized like package variables, and their initializers can refer to other static members regardless of the order of declaration. Static methods declared with `func T.Method(...)` can be called with `T.Method(...)` from anywhere in the package, including from other static methods and from the initializers of static variables.

Static members are accessed using `T.name`. Inside a classfile, static members of the class can be used without a type prefix, from instance methods, static methods, and initializers. They then take precedence over their namesakes at the package level, and local variables and parameters still shadow them as usual.
```

# Research notes

- Proposal #2513 "Static Members" (xushiwei, Proposal-Accepted, FullSpec). Parser/AST work is merged (#2783, #2785, #2793). Open PR #2792 is parser-only, and its parsing already exists at this commit; it says "the cl/gogen semantic implementation is split out into a follow-up PR", and that PR doesn't exist. In gogen, #638/#639 (merged, and part of v1.26.3, the version xgo uses) add `NewStaticMember`, and InitXGoPackage registers `XGos_T_name` objects as static members. cl uses none of this for values.
- Behavior at base:
  - `const T.name` / `var T.name` produce invalid Go output (`const T.name = "x"`), and any use is "T.name undefined".
  - Static methods only resolve when T is loaded before InitXGoPackage. `func T.Two() int` fails while `func T.New() *T` works, because the signature that mentions T is what loads it.
  - In a classfile, `.name` values are unsupported and unqualified static calls (`New(...)`) don't resolve.
- The platform has 3/50 other submissions on goplus/xgo. One of them is virtual fields (derivative of our draft); the others are unknown, so there is some risk that static members were also taken.
- Solution (cl only; 3 files: new `cl/static.go` plus compile.go and expr.go):
  - Lower `T.name` / `.name` to `XGos_T_name` at preload and load T first (lbinames), which also fixes the static method registration.
  - Check at load time that T is a defined, non-alias type of the package and that there is no method with the same name (otherwise report an error and declare the value as `_`).
  - Split classfile var blocks into class fields and package-level static vars.
  - Resolve `T.name` (selector, LHS too) and unqualified names in classfiles directly to the lowered object, so it works in static method bodies and var initializers compiled before gogen registers the members.
  - In classfiles, static members take precedence over receiver members (once registered they also appear as methods of the receiver) and over globals.
- Tests: `cl/staticmember_8b23ff_test.go`, 7 tests. Six compile the program, run it and compare output; one checks compile errors on the marked line. They cover constants (typed, iota, mixed specs, array length), variables (++/+=, address, maps/slices, initializers in any order), static methods whose signatures don't mention T, recursion and calls from initializers, same names on different types, classfile `.name`/`T.name` with unqualified use from instance and static methods, precedence/shadowing, and errors (undefined type, non-type, method conflict, assignment to a static constant qualified and unqualified).
- Dockerfile/test.sh are the same as in xgo-go-expression.
- Verified in a clean clone with `--network none`: base 1165 passes without the solution and 1166 with it; new 7/7 fail -> 7/7 pass.
