# Per-check failure patterns

## 1. Tests fail without the solution
- Every new test must fail in `new` mode without the solution; a passing new test does not capture the missing behavior. Check test by test, not the suite verdict.
- Base and new must run disjoint tests; every wildcard in base mode (`./...`, `pytest .`, globs) must exclude the new test directories.
- A test that crashes the whole test binary at base should run last or be wrapped in recover/try-catch that reports a normal failure, so other tests still report.
- Do not edit existing tests.
- Run as an isolated environment would: unprivileged user, unlimited stack (some runtimes exit with no output otherwise), no network (tests needing network must skip at base), clean clone, submodules present.

## 2. Solution passes everything
- Each `package + test name` appears exactly once in the base junit. A duplicated name (same package run by two runners) yields conflicting results and one is dropped.
- If test IDs embed fixture line numbers, base mode should restore the test directory from the baseline commit so edits don't produce false "missing" tests.
- The test runner must not swallow errors (`|| true`, ignored list/collect failures, relative paths after `cd`) and must use only valid flags.
- Source in folders named `build`, `dist`, `target`, `node_modules`, `__pycache__`, `.venv` is easily dropped by tooling; avoid putting real source there or make sure patches include it.

## 3. No flakiness
- Repeat each mode several times and compare per-test results. Vary hash seeds, test order and timezone where the runner allows (see `languages.md`).
- Suspects: map/set iteration order, test order and caches (a test that passes in the suite but fails alone: run each test in isolation from a script, not `xargs`, which splits ids containing spaces), time, ports, shared temp files, network.

## 4. Tests are fair
Every assertion must be justified by the prompt or discoverable in the repo.
- An assertion is justified if any reasonable reading of the prompt implies it, including through a general rule ("normal typing rules", "every possibility", "as the language does"). Flag it only when no sentence, general rule or repo convention implies it. For each flagged assertion quote the closest sentence and say why it does not cover it; a flag without that is noise.
- Typical unfair assertions: exact error type or message not stated, formatter idempotency, behavior of inputs the prompt never mentions, results through a source the prompt doesn't list (e.g. a helper's return value when it names only literals).
- Two-way rules ("or the other way around") need both directions tested; so do rules about nil/empty/missing values.
- Rules about unpacking, spreading or recursion need mixed-kind inputs, otherwise an implementation that skips the construct still passes.
- Cover plausible mutants: double evaluation, errors before short-circuit, slice steps, large indices, every call form the prompt lists.
- Don't add tests only to prove a bug in the solution is fixed; fix the solution. Over-testing makes every implementation fail.
- Comment non-obvious test-runner choices so a reviewer doesn't guess.

## 5. Prompt and tests agree
- A prompt changed after the tests were written leaves the tests asserting old behavior. After changing either, re-check both directions.
- Everything a test asserts must appear in the prompt.

## 6. Solution complete and clean
Reviews tend to find a new gap each time, so audit before declaring done instead of fixing only what was found.

Do it systematically, not by sampling probes you happen to think of:
1. List the constructs the prompt talks about (each sentence gives one or more) and the source kinds each can take (literal, variable holding a literal, annotated fixed-size, annotated unbounded/unknown length, conditional/alternative values, unresolvable). Write the full construct × source-kind table to a file, then probe every cell on the solution and on the language's real behavior. A construct with only a literal-source probe is not covered; this is where the commonest misses hide (for example unpacking a value of unknown length against overloads that take a fixed number of arguments).
2. For every helper the solution adds or extends that returns more than a value (a flag such as `exact`/`complete`/`maybe_missing`, a list of alternatives), grep all of its callers and check that each one uses the extra information. A caller that ignores it is a defect even if no test shows it.
3. For every limit, fallback or `except` in the new code, state what happens past the limit and check it fails closed (no match), never approximates.
4. For every construct that has several syntactic variants (comprehension with several `for` clauses, nested unpack), probe each variant, not only the first.
5. Only then read the rest of the patch for the points below.

- The promised behavior holds in every equivalent form, not just the tested one. Enumerate forms: syntax variants, input source kinds (literal, variable holding literal, annotated fixed-size, annotated unbounded, unresolvable), and every parallel code path (package- or OS-specific branches, reflection, generics, compiler vs runtime helpers). Parallel paths must call shared helpers with the same arguments.
- Merging alternatives keeps presence ("may be missing") as well as values; "cannot be determined" means "no match", never "no arguments".
- Size and depth limits fail closed. Don't use a fixed recursion cap that silently turns inferable into unknown; use a visited set.
- A helper that returns extra information (an `exact` flag, a "may be partial" marker) must be honored by every caller; grep all call sites. Check that reused library helpers don't silently drop syntax the prompt mentions.
- Identity or type promises hold at every place objects are created or cached; grep every constructor and cache.
- Secondary features added to keep old tests green get reviewed as hard as the main one: do them completely or keep the old behavior.
- No global change for a local feature; gate it on the new path. A new AST node or field means auditing every tree walker.
- Probe types that lack runtime wrappers as well as ones that have them.
- Run the whole probe set on the old and new version of the solution and diff results, not only the cases just fixed. Review remaining differences yourself; a reviewer's expectation can be wrong.
- Lint clean, no comments that only narrate code, patch built from every touched directory, patch file itself tested on a clean worktree in both modes plus the full suite.
- Keep probes, patches and scripts in the project directory, not in temporary scratch.

## 7. Prompt is clean
- Count feature groups first: the independent areas a solver must change (for example one for `type[...]` matching, one for `None`/`Optional`, one for `Callable`). Name them. More than two groups is a red flag even for a short prompt, because each group brings its own cluster of clauses and the "too many requirements" reviewers count behaviors, not words. If the prompt names three or more, recommend dropping a whole group from prompt, tests and solution.
- Then count requirement clauses: independently testable behaviors the prompt states. A general rule that covers many cases counts once ("every element must be accepted"); a list of cases counts each item. About 330 words and ≈15 clauses is the budget; treat 16–18 as a warning to tighten, and only fail clearly above that or when the groups check fails. Do not fail a prompt that is under 330 words and expresses edge cases through a few general rules.
- Describe behavior, not method. Draft in the author's own language and translate; a raw English draft reads as machine-written.
- The description check includes an AI detector (category `ai_generated`, reported as "too prescriptive — focus on behavior instead of implementation details, write it like how you would prompt an LLM"). It cannot be run locally and its verdict can change between runs on the same text, so judge the risk by reading, passage by passage. Flag as HIGH a passage that reads like a spec sheet rather than a maintainer's issue: a chain of rules joined by colons, semicolons or commas in one sentence; a list of constructs or types ("including X, Y, Z and W"); passive rule phrasing ("is checked according to", "is considered according to"); several sentences in the same "X should Y; if Z then W" rhythm. Flag as MEDIUM any English edited by hand after translation. Prefer: open with the ask, then the issue's example, then one general principle per paragraph followed by short concrete examples with their expected result. Evidence: jedi v10 (Vietnamese → gtx) passed on 2026-10-26 and failed on 2026-10-27 with identical text; the two flagged passages were a rule chain and a construct list. Gopherjs passed only after its hand-written English was redone as Vietnamese → gtx. A later jedi rewrite (one principle per paragraph plus examples) still got the opening sentence and the principle sentences flagged while the example-only paragraph passed: keep principle sentences short and phrased as direct requests to the solver, and let code (a fenced issue example) and concrete cases carry most of the text. Rewriting the same content as short imperative requests also failed, and a paragraph that passed once was flagged on the next run, so no passage is safe by itself: write the whole prompt as a user's bug report (symptom with the issue example, what was expected, a few cases the user tried).
- Remove: examples that only restate a rule, "other cases keep the old behavior", description of the current bug, "X is not part of this change", re-statements of what must keep working (regression tests cover that).
- Don't state more than tests and solution do ("any form X supports"), and don't leave out anything a test asserts.
- Answer a reviewer by merging edge cases into one general rule, not one sentence per edge case. If the clause budget is exceeded, drop a whole feature group from both prompt and tests.
- Avoid terms with broader technical meaning than the task intends.
- Every sentence is a promise the solution will be probed on; after a rewrite, probe each new sentence.

## Environment (Dockerfile)
- Exact pins or a lockfile; for Python install the project with `pip install --no-deps -e .`. No tests run inside the image. A second toolchain gets a comment explaining why. Assert tools exist (`RUN test -x ...`). Include submodules and data the tests need.
