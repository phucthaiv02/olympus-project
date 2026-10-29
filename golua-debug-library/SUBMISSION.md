# Repository

**GitHub repo**

```
https://github.com/Shopify/go-lua
```

**Commit SHA**

```
3f823f845ff9c2ce15a9877b60cf374ce570bffd
```

**GitHub issue**

```
https://github.com/Shopify/go-lua/issues/89
```

Draft: `kh76r4rqt81rye1rma8q5rza198g697d` (pushed 2026-10-20, no checks run)

# Challenge Details

**Title**

```
Implement debug.getinfo, debug.getlocal and debug.setlocal
```

**Language:** Go

**Category:** Feature Request

**Description**

```
Implement `debug.getinfo`, `debug.getlocal` and `debug.setlocal` as described in the Lua 5.2 reference manual, and make debug hooks and the information they observe behave like the reference implementation. Functions written in Go keep the existing conventions of `Info`: `what` is `"Go"`, `source` is `"=[Go]"` and `short_src` is `"[Go]"`.

`debug.getinfo([thread,] f [, what])` accepts a level or a function and returns a table with the fields selected by `what` (default `"flnStu"`). Level 0 is `getinfo` itself, and a negative level or one beyond the stack returns nil. A `what` of `""` returns an empty table. Each name reported with option `n` follows how the caller invoked the function: `"global"`, `"local"`, `"method"`, `"field"` or `"upvalue"`, the event name such as `"__index"` with namewhat `"metamethod"`, or `"for iterator"` for the iterator of a generic `for`. A value called through `__call` is named like an ordinary call. When the caller is not Lua code, after a tail call, or when `f` is a function, `name` is nil and `namewhat` is `""`. Option `f` returns the function, Go functions included. Option `L` gives a Go function nil activelines.

`debug.getlocal([thread,] level, n)` returns the name and value of the n-th local variable that is active at the current instruction of that level, or nil. Slots in use that have no variable name are reported as `"(*temporary)"`, including those of Go functions, and a negative `n` selects vararg values, named `"(*vararg)"`. When the first argument is a function, only its parameter names are returned. `debug.setlocal` assigns to the same variables and returns the name, or nil without assigning anything. An invalid option is reported as bad argument #2 with `invalid option`. A first argument that is neither a level nor a function is bad argument #1 with `function or level expected`. A level out of range for `getlocal` or `setlocal` is bad argument #1 with `level out of range`.

Line hooks fire for the line of the instruction about to run when a function is entered, when execution moves to a new line, and when a jump goes backward. In line and count hooks, level 2 is the interrupted function, observed at that instruction. `getlocal`/`setlocal` at that level work on the variables active there. Call, tail call and return hooks report the function being called or returning, Go functions included. `debug.gethook` returns the count given to `sethook`.

`short_src`, which also prefixes error positions, follows the reference implementation's 60-byte limit. A `=` or `@` source of at most 60 bytes is shown without its first character. A longer `=` source shows the next 59 characters after the `=`, and a longer `@` source shows `...` followed by its last 56 characters. Any other source is shown as `[string "..."]`. A one-line source shorter than 45 characters appears unchanged. Otherwise the source is cut at its first newline and to at most 45 characters, and `...` is added before the closing quote.
```

# Research notes

- Platform cap: Shopify/go-lua 2/50 used (48 left). MIT, 3.4k stars, active (last commit 2026-10-07).
- Issue #89 "Implement debug.getInfo()" is open and has no PR. The more visible gaps (pattern matching #142 with open PR #76, coroutines #88) were avoided, both because of the open PR and because the 2 existing submissions likely cover them.
- Oracle: the reference Lua 5.2.4 was built from lua.org sources (sha256 b9e2e4aa…f4b) and used to produce every expected value. The only substitution is C→Go for Go functions (go-lua's existing `Info` convention).
- Hidden depth: a direct port of ldblib's db_getinfo/db_getlocal/db_setlocal still fails all 7 test functions, because of these existing bugs:
  - `functionName` and `findLocal` used the next pc instead of the current one, so names were wrong and locals not yet declared were visible.
  - `objectName` returned "local" for table fields, and `constantName` checked kind "c".
  - `Info` panicked on an empty `what`, and `f` pushed nil for Go functions.
  - `chunkID` did no truncation for multi-line or long sources, and kept the wrong end of `@` names.
  - Line and count hooks ran before pc was incremented. Line events came from the previous instruction, `internalHook`'s `Info("lS")` hit index -1 on function entry (any Lua `debug.sethook(..., "l")` plus a call crashed), and `callHook` read `pci.code[savedPC-1]` while inside a traced frame.
  - `DebugHookCount` returned the remaining counter.
- Solution (5 files):
  - debug.go: getinfo/getlocal/setlocal, `Local`/`SetLocal` Go API, `currentPC`, chunkID, the Info fixes, and DebugHookCount.
  - lua.go: `XMove` and the `callStatusTraced` flag (callStatus is now uint16).
  - stack.go: callHook.
  - types.go: objectName.
  - vm.go: traceExecution marks the frame as traced and computes line events for the instruction about to run.
- The existing TestProtectedCall reads `ci.code[ci.savedPC]` inside a count hook, so the engine's savedPC convention is kept and the traced flag is used instead.
- Out of scope (they also differ from 5.2.4 but are not tested): the variable names in runtime error messages, xpcall handler levels (open PR #147), and (*temporary) slots inside table.sort (go-lua's sort is implemented differently).
- Tests: `debug_library_a2959f_test.go`, 7 table-driven tests with 59 subtests. Lua snippets run in fresh states; the expected strings come from Lua 5.2.4.
- test.sh: `go test -vet=off` (vet rejects doc_test.go's Example). base = `./...` skipping `_a2959f$`; new = `-run '_a2959f$'`.
- Dockerfile: fetches the lua-tests submodule at its pinned commit when the build context lacks it (TestLua needs it).
- Verified in a clean clone without submodules, with `--network none`: base 53 tests (4 skipped, they need luac) pass with and without the solution; new 66/66 fail -> 66/66 pass.
