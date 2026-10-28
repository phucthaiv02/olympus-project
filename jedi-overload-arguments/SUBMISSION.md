# Repository

**GitHub repo**

```
https://github.com/davidhalter/jedi
```

**Commit SHA**

```
2af1a0f3943f801f5de755f8ffd251f0385a3607
```

**GitHub issue**

```
https://github.com/davidhalter/jedi/issues/2028
```

# Challenge Details

**Title**

```
Overload resolution ignores type[...] and *args/**kwargs annotations
```

**Language:** Python

**Category:** Bug Fix

**Description**

```
Fix how Jedi picks an overload when its parameters are annotated with `type[...]` or collect arguments through `*args` and `**kwargs`. Given overloads `get_obj(input: type[Foo]) -> FooBuilder` and `get_obj(input: type[Bar]) -> BarBuilder`, completing `get_obj(Foo).` offers the methods of both `FooBuilder` and `BarBuilder`, when only `FooBuilder`'s belong there.

Jedi should choose the overload the way a type checker does: take the first overload that accepts every argument of the call, and fall back to the combined results of all overloads only when none of them fits. This holds for arguments passed by position or by keyword, and for the arguments collected by `*args` and `**kwargs`, where the annotation describes each collected argument on its own. Arguments unpacked at the call site from a tuple, list or dict count as well, and a variadic parameter that receives no arguments fits.

A `type[X]` (or `Type[X]`) parameter accepts the class `X` and its subclasses, but not instances; in the other direction, a parameter annotated `X` does not accept the class `X` itself. Bare `type` or `Type`, and `type[Any]`, accept any class. A union may appear inside or outside `type[...]`. For `type[T]` where `T` is a TypeVar, the class has to satisfy `T`'s bound or one of its constraints, and any class is accepted when `T` has neither. The argument can be a class named directly, or a variable or parameter annotated `type[...]`.
```

# Research notes

Olympus draft: `kh7d8qtdcfr0k5mt1rvftcvebn8gb2va` (version `kd73n08qh9kjrnt507bysqkpmn8gb1t7`).

- Duplicate check: no PR touches overload argument matching. #2117 only recognizes `@typing.overload`, #2111 is the interpreter Literal path, #2109 is a docstring fix. The issue has no linked PRs and no maintainer rejection.
- Scope: `type[...]`/`Type[...]`/bare `type`, `None`/`Optional`, `Any`, TypeVar bound/constraints, numeric tower (bool/int/float/complex), `Callable` including its result type. Protocols stay lenient. `Literal` is out of scope.
- Solution: 7 files (`annotation_accepts` per annotation value). 29 hidden tests. Docker `--network none`: base 3955 passed with and without the solution; new 29/29 fail without, 29/29 pass with (stable over 3 runs). test.sh deselects `test_file_path_completions` (it lists the repo root, which now holds test.sh) and `test_complete_expanduser` (home folder of root).
- Prechecks (2026-10-15): GitHub WARNING (license not auto-detected; jedi is MIT); Problem & Tests WARNING (suggests testing "uninferable callable result is accepted"; conciseness suggestions); Dockerfile PASS; Solution PASS.
- Scope Gate (2026-10-15): **PASS**, high confidence. Medium finding: borderline overlap with an older platform submission (53/108 solution lines, 49.1%); the gate judged each task retains independent core work.

## Version 5: type[...] + *args/**kwargs (2026-10-16, local only, NOT uploaded)

- Version 4 had too many requirement groups. Scope review using Nova's passing patch (run rd72padr481k, 722K tokens) as a yardstick: Nova's patch passes all type[...] tests (incl. coverage), fails None/Callable/metaclass/variadic only because they were not asked. None is trivial and Callable is in the "derivative" area; metaclass is ~5 lines. Kept type[...] (issue anchor) + variadic (novel, ~25 lines, touches argument packing), folded into the overload-selection paragraph so the description has 2 groups.
- 28 hidden tests: 28/28 fail at base, 28/28 pass with the solution (3 runs); Nova's patch fails the 8 variadic tests. Base suite 3955 passed with the solution.
- Version 4 files backed up in `v4-metaclass-variadic/`.

## Version 4: + metaclass and *args/**kwargs (2026-10-16, local only, NOT uploaded)

- Nova solved version 3 on the first try, so two groups were added: parameters annotated with a metaclass accept classes created by it (inherited metaclass, metaclass subclass, stdlib ABCMeta/EnumMeta), and `*args`/`**kwargs` annotations are checked against every collected argument (incl. unpacked tuple/list/dict, empty = fits).
- 29 hidden tests: 29/29 fail at base, 29/29 pass with the solution (3 runs); 9 of them fail with the version-3 solution. Base suite 3955 passed with the solution.
- Gotcha found while building: values like `_DynamicArrayAdditions` (list()/set() typeshed hack) are not full `Value`s, so the metaclass lookup must guard `isinstance(value, Value)`; without it ~250 base tests fail.
- Version 3 files backed up in `v3-type-coverage/`.

## Version 3: type[...] only (2026-10-15)

- Narrowed again after "too many requirements": dropped the None/Optional and Callable groups (tests, solution and description). The previous version is backed up in `v1-three-groups/`.
- 19 hidden tests. Docker `--network none`: new 19/19 fail at base, 19/19 pass with the solution; base 3955 passed with the solution.
- Solution: 6 files (compiled/value.py no longer touched). Prechecks and Scope Gate not re-run for this version.
- Uploaded test/solution/title to the draft (2026-10-15); "Problem and tests are aligned" then passed.
- Coverage pass (local only, NOT uploaded): added first-match-wins, more keyword cases, a two-parameter overload (every argument must fit, mixed positional/keyword), bare `type`/`type[Any]` with more classes and rejecting instances, more bound/constraint/unbound TypeVar cases. 20 tests: 20/20 fail at base, 20/20 pass with the solution (3 runs), base 3955 passed. Previous test.patch kept in `v2-type-only/`.

## Status when the session stopped (2026-10-15, version 2)

- Draft now holds the NARROWED version (3 requirement groups: `type[...]`, `None`, `Callable` with result type), 30 hidden tests (fail 30/30 at base, pass 30/30 with the solution in Docker; base 3955 passed).
- Prechecks 5/5 pass. `plagiarism_review` WARNING: "derivative" (0.78) of an older platform submission that implements a broad overload-evaluation engine for jedi.
- Scope Gate is stale for this version (not re-run). It passed the earlier broad version (overlap 40-49%, borderline).
- Last results on the broad version: verifyTests, verifySolution, taskQuality, solutionQuality, verifyFairness PASS; descriptionQuality FAIL ("too many requirements"); verifierIncompleteness INCOMPLETE (5 unpinned behaviors, 3 of them now covered by tests). No rollouts run.
