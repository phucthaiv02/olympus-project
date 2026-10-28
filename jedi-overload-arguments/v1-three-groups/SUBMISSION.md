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
Overload resolution ignores type[...], None and Callable annotations
```

**Language:** Python

**Category:** Bug Fix

**Description**

```
Completion is wrong after calling an overloaded function that is overloaded on classes, `None` or callables. With overloads `get_obj(input: type[Foo]) -> FooBuilder` and `get_obj(input: type[Bar]) -> BarBuilder`, `get_obj(Foo).` offers the methods of both `FooBuilder` and `BarBuilder`, because Jedi returns the results of every overload instead of picking one.

Make Jedi choose the overload the way a type checker would: take the first overload that accepts the call's arguments, and only fall back to the combined results when no overload fits.

The cases that need to work:

- `type[X]` / `Type[X]` parameters accept the class `X` or its subclasses, not instances. This includes unions, `Any`, bare `type`, and TypeVars with a bound or constraints, and the class may come from a variable or parameter annotated `type[...]`.
- `None` is accepted only where the annotation allows it (`None`, `Optional[...]`, a union containing `None`).
- `Callable` parameters accept anything callable (functions, lambdas, builtins, classes, bound methods, instances of classes that define `__call__`) and reject plain values. If the annotation states a return type, the callable's return value has to fit it; if Jedi can't infer what the callable returns, treat it as fitting. The parameter list inside `Callable[...]` doesn't need to be checked.
```

# Research notes

Olympus draft: `kh7d8qtdcfr0k5mt1rvftcvebn8gb2va` (version `kd73n08qh9kjrnt507bysqkpmn8gb1t7`).

- Duplicate check: no PR touches overload argument matching. #2117 only recognizes `@typing.overload`, #2111 is the interpreter Literal path, #2109 is a docstring fix. The issue has no linked PRs and no maintainer rejection.
- Scope: `type[...]`/`Type[...]`/bare `type`, `None`/`Optional`, `Any`, TypeVar bound/constraints, numeric tower (bool/int/float/complex), `Callable` including its result type. Protocols stay lenient. `Literal` is out of scope.
- Solution: 7 files (`annotation_accepts` per annotation value). 29 hidden tests. Docker `--network none`: base 3955 passed with and without the solution; new 29/29 fail without, 29/29 pass with (stable over 3 runs). test.sh deselects `test_file_path_completions` (it lists the repo root, which now holds test.sh) and `test_complete_expanduser` (home folder of root).
- Prechecks (2026-10-15): GitHub WARNING (license not auto-detected; jedi is MIT); Problem & Tests WARNING (suggests testing "uninferable callable result is accepted"; conciseness suggestions); Dockerfile PASS; Solution PASS.
- Scope Gate (2026-10-15): **PASS**, high confidence. Medium finding: borderline overlap with an older platform submission (53/108 solution lines, 49.1%); the gate judged each task retains independent core work.

## Status when the session stopped (2026-10-15)

- Draft now holds the NARROWED version (3 requirement groups: `type[...]`, `None`, `Callable` with result type), 30 hidden tests (fail 30/30 at base, pass 30/30 with the solution in Docker; base 3955 passed).
- Prechecks 5/5 pass. `plagiarism_review` WARNING: "derivative" (0.78) of an older platform submission that implements a broad overload-evaluation engine for jedi.
- Scope Gate is stale for this version (not re-run). It passed the earlier broad version (overlap 40-49%, borderline).
- Last results on the broad version: verifyTests, verifySolution, taskQuality, solutionQuality, verifyFairness PASS; descriptionQuality FAIL ("too many requirements"); verifierIncompleteness INCOMPLETE (5 unpinned behaviors, 3 of them now covered by tests). No rollouts run.
