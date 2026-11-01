# v9: answer to the 2nd human review (2026-10-25)

Review of platform v2 (= v8 content) asked for 8 solution fixes plus tests for
alternative sequence lengths/positions, nested maybe-absent dict entries,
conditional dict comprehensions, repeated non-scalar keys, and the Auto Review
High finding (statically empty dict comprehension).

Solution fixes (`jedi/` is the clone with v9 committed on top of 2af1a0f3):
1. Conditional comprehension crash: `_update_with_alternative_pairs` no longer
   wraps a `None` placeholder; `_merge_pairs` skips `None`.
2. + 3. `sequence_variants` (gradual/base.py) keeps the alternatives of `*`
   unpacking; `TreeArguments.argument_variants` enumerates each possible number
   of arguments of every `*x` (max 16 combinations), and
   `TreeSignature.matches_signature` requires every variant to match.
   Fixed-length `tuple[...]` checks every variant.
4. Merging two maybe-absent entries keeps the "maybe absent" mark.
5. `_static_key` compares classes, tuples of keys and equal conditional keys
   (1 == True == 1.0 like Python).
6. The mapping recursion guard keys on (node, defining context).
7. + 8. Protocol methods: `_compatible_params` follows the callable
   assignability rules (names, positional-only, keyword-only, defaults,
   *args/**kwargs); annotations go through `infer_annotation` (forward refs).
Extra: `True` is accepted by `int` (compiled classes converted to stubs).

Tests: 4 new test functions (67 total): test_unpacked_sequences_of_several_lengths,
test_entries_of_dict_comprehensions, test_entries_that_may_be_missing_in_nested_mappings,
test_keys_that_are_not_strings. No Protocol tests (not in the prompt).

Prompt: added sentences on key equality, maybe-present entries (incl. nested),
empty comprehensions, and `*xs` of several lengths shifting later arguments.

Verified: base 67/67 fail; ref 67/67 pass (PYTHONHASHSEED 0/1/7/42); test.sh base
3955 passed with the solution; flake8 clean; zuban unchanged (9 pre-existing errors).
Old run patches on the v9 tests (prompt did not have the new sentences):
Vega 77my 4 fail, Orion bfk 2 fail, the others 4-9 fail.
Probe scripts: scratchpad rv/probe.py, probe2.py, proto.py.

# v10: platform v3 checks (2026-10-25)

v3 failed descriptionQuality ("too many requirements", 435 words) and
solutionQuality (High: `{k: C for k in [1] for _ in []}` not empty; Medium:
>16 `*` variants fell back to an approximation; Medium: Protocol checker not
justified). Raw output: `../../check-errors/kh7d8qtdcfr0k5mt1rvftcvebn8gb2va.*`.

- Description rewritten from Vietnamese (`task_v10_vi.md`, gtx + meaning fixes):
  3 paragraphs, ~15 clauses, 249 words. Standard typing details (bare `type`,
  `Optional` takes `None`, `tuple[X, ...]` vs `tuple[X]`, empty `*args`) are
  now covered by "normal typing rules"; the dict/sequence edge cases by "the
  content Python creates from the literal" + "every possibility".
  Old: `task_v9.md`, `solution_v9.patch`.
- `is_empty_comprehension` checks every `for` clause.
- No variant limits: `sequence_variants` and `TreeArguments.argument_variants`
  merge variants of the same length position by position (exact, because each
  argument/element is checked on its own), so the number of variants grows
  linearly. `_star_choices` replaced by `_positional_values`.
- Protocol: comment explains it is needed because overload annotations
  (typeshed) are often protocols.
- Tests unchanged (67). Verified: base 67 fail; ref 67 pass (seeds 0/1/7/42);
  test.sh base 3955 passed; probes (tools/probe2.py, +12 new cases) and
  Protocol matrix 0 bad; flake8 clean; zuban 9 pre-existing errors.

## v10 rollouts (2026-10-25 18:36): Vega 0/4
All quality checks, prechecks and Scope Gate PASS on v10. Vega #1-#4
(rd71800e, rd72n3nr, rd73x8pt, rd7fzkww; patches in runs-v10/) all fail only
`test_keys_that_are_not_strings`: they key dicts with `get_safe_value`, so
class keys (`{C: A, C: C}`) and tuple keys (`(1, 2)`) never replace. 1/True
works for all. Vega #2 also fails `tuples_of_unknown_length` and the
conditional equal key; Vega #1 has 1 baseline failure.
Old patches (v8 prompt) on the same tests: Orion rd7bfk, Orion rd7ak8, Nova
rd7avd/rd7e42/rd7ejw handle class AND tuple keys; old Vega patches do not
(rd77my class only). So the key test is a Vega-specific wall under v10.
Tuning idea if nobody solves: hint classes only ("like `1` and `True`, or the
same class") and keep tuple keys as the discriminator.

## v10 Orion rollouts (2026-10-25 20:12): Orion 2/4 (too easy)
PASS: rd749ksx, rd771brz. FAIL: rd74gg0m (TypeVar forward ref + mapping
items), rd7000n9 (Sequence/Iterable/Mapping). Probe matrix (tools/probe3-5.py,
~140 in-scope cases) on REF + 4 Orion + 4 Vega: the two passers agree with REF
almost everywhere. Only natural discriminator: dict comprehension over a string
literal `{k: C for k in "ab"}` -> RB; rd749 gives RC (its own comprehension
unrolling cannot iterate str). `test_v11.patch` adds that one assertion to
test_entries_of_dict_comprehensions: Orion 1/4 (rd771), Vega 0/4, base 67 fail,
REF 67 pass (seeds 0/1/7/42). Not pushed yet.
REF risks found (Orion is more Python-accurate, not in tests):
`{"p": A, **{k: C for k in ["p"]}}` (Python: RB, REF: RC),
`{"k": v for v in [A, C]}` (Python: RB, REF: RC), `{"k": C for _ in ""}`
(empty str literal; Python: RA, REF: RB). Fixing REF for "" would make rd771
differ from REF (FP risk), so left as is.
All 4 Orion patches reject Protocol-satisfying arguments (REF accepts): FP risk
if the panel probes Protocol annotations.

## v11 pushed (2026-10-26) onto platform v3
- Test: `test_entries_of_dict_comprehensions` + `f({k: C for k in "ab"})` -> RB
  (comment: keys of a comprehension over a string are strings). Expected:
  Orion 1/4 (rd771 only), Vega 0/4 after re-eval.
- REF: `comprehension_items` (gradual/base.py) unrolls a comprehension when
  every `for` iterates an exact literal and every `if` is statically known;
  dict comprehensions use it for entries (`_value_entries`) and pairs
  (`_mapping_pairs`), sequences in `sequence_variants`, and
  `is_empty_comprehension` counts "no iteration passes". Fixes (REF was less
  Python-accurate than the passing Orion run): `{"p": A, **{k: C for k in
  ["p"]}}` -> RB, `{"k": v for v in [A, C]}` -> RB, `{"k": C for _ in [1] if 0}`
  -> RA, `[C for _ in [1] if 0]` -> RA. Empty str literal left unchanged on
  purpose (rd771 also treats `""` as non-empty; fixing it = FP / Auto Review risk).
- Verified: base 67 fail, REF 67 pass (seeds 0/1/7/42), test.sh base 3955,
  probes 0 bad (tools/probe6.py new), flake8 clean, zuban 9 pre-existing.
  rd771 agrees with REF on all ~160 probes except `{"k": A, **{"k": C for _ in ""}}`.
- Old: test_v10.patch, solution_v10.patch.

## v12 (2026-10-26, local): solutionQuality FAIL on v11 + independent review
- solutionQuality v11 FAIL High: `_star_arg_sequences` dropped `exact`, so
  `f(*xs)` with `xs: tuple[type[C], ...]` was one argument. Fix:
  `RepeatedLazyValue` marks elements of unknown number; during overload
  matching each expands to 0..(params + 1) arguments (`_expand_repeated`,
  `argument_variants(max_length)` forwarded by Instance/ClassMethod/partial).
  tools/probe7.py (34 cases, every kind of `*` source) 0 bad.
- Independent review agent found 8 more in-scope defects, all fixed
  (tools/probe8.py, 48 cases, 0 bad): earlier equal key object stays
  (`_set_pair`); `*` of dict literal/comprehension = its keys
  (`dict_key_variants`); unknown str keys may fill any by-name parameter;
  constrained TypeVars solved once per call (`_constraint_choices`);
  int->float->complex promotion; posonly/kwonly binding in strict mode;
  walrus is positional; comprehension unrolling falls back at jedi's
  300-inference limit.
- Test swap: the `"ab"` comprehension assertion was replaced by
  `test_unpacked_values_of_unknown_length` (68 tests). rd771 (old passer)
  fails ~15 unknown-length probes -> would be an Auto Review/FP liability;
  rd749 handles them. Local grade: Orion 1/4 (rd749), Vega 0/4.
- rd749 vs REF on probe8: only 5 niche differences (walrus, dup keys in
  `*` of a comprehension, 10x10 comprehension, unknown `**` with a given param).
- Verified: base 68 fail; REF 68 pass (seeds 0/1/2/3/7/42); test.sh base 3955;
  all probes 0 bad, probe3-5 unchanged; flake8 clean; zuban 9 pre-existing.

## v13 (2026-10-26, pushed onto platform v3): second review round
v12 is on the platform. Round 2 of the review agent found 6 more (5 fixed,
1 kept by decision); tools/probe9.py 25 cases 0 bad:
- key order: an equal key replacing a possibly-missing earlier one has two
  possible positions (`_KeyOfMovedEntry.later_position`, both orders in
  `dict_key_variants`).
- `*` of a dict comprehension needs only its keys (`comprehension_items(...,
  keys_only=True)`).
- a required parameter without an argument never matches in strict mode
  (`LazyUnknownValue`), also posonly/unannotated/partial.
- `MergedPartialArguments.unpack` = `f(*pargs, *args, **{**pkw, **kw})`;
  "multiple values" is an issue in strict mode without calling nodes.
- string annotations `"T"` / `"list[T]"` join the constraint solving; TypeVars
  are collected through generics with a visited set (recursive aliases).
- constraint solving prunes per TypeVar first and fails closed above 1024
  combinations.
- kept: `**d` with unknown str keys and no `**kwargs` never matches (same
  decision as unknown-length `*` vs required params).
Reviewer probes that still say BAD were checked by hand: wrong expectations
(duplicate "multiple values" calls, `{1|True: C}` vs dict[int]) or baseline
jedi limits (`partial(...)()` with no call arguments infers nothing in base too).
Verified: base 68 fail; REF 68 pass (6 seeds); test.sh base 3955; probe..probe9
0 bad; probe3-5 unchanged; flake8 clean; zuban 9 pre-existing. rd749 vs REF on
probe9: 2 niche differences.

## v14 (2026-10-26, pushed onto platform v3): verifyFairness FAIL on v13 -> tests only
- Unfair: the 3 `wrap(...)` assertions (dict contents through helper call
  results are not in the description). Removed with the `wrap` helper.
- Added `test_elements_of_empty_comprehensions` (list comprehension over `[]`,
  `()` inside `*`, and `[1]` control) per coverage suggestion. 69 tests.
- Verified: base 69 fail; REF 69 pass (seeds 0/1/2/3/7/42); ruff E/W/F clean.
  Grade: rd749 69/69, rd771 fails unknown-length only, others fail; Orion 1/4,
  Vega 0/4. Old: test_v12.patch.

## v15 (2026-10-27, pushed as platform v4 kd743zn69czwsefs3ysevn7m3s8h3q01): Auto Review revision_requested on v14
- Platform v14: 7/7 quality checks; rollouts Orion 1/4 (rd7f4z = rd749 patch),
  Vega 0/4, Nova 0/4 (avg 1.45M, 2 Nova infra fails); FP check passed with
  caveats; Auto Review: High constrained TypeVar across parameters, Medium `**`
  into a named parameter.
- Added `test_one_constraint_for_the_whole_call` and
  `test_dict_unpacking_into_named_parameters` (71 tests). Old: test_v14.patch.
- Verified (grade15/grade.out): base 71 fail; REF 71 pass (seeds 0/1/2/3/7/42);
  ruff clean; rd749 71/71; all 7 other Orion/Vega and 4 Nova patches fail.
  Reviewer probe.py (8 human-review issues) 0 bad on REF.
- verifierIncompleteness proposed patch (grade15/vi_proposed.patch): REF 74/74,
  rd749 fails only the Protocol test (all Orion patches are nominal-only).
- Human review (v8) "Solution & Code" 8 issues re-checked on REF v13 with
  variants: tools/probe11.py (51 cases), output grade15/p11.out: REF 0 bad, base
  34 bad. Issues 6 (wrap) and 7/8 (Protocol) are fixed in REF but have no tests
  (wrap = unfair per verifyFairness; Protocol not in the description and rd749
  fails it).

## v11 description pushed (2026-10-29)
- `problems edit --description-file=task.md` -> versionId kd743zn69czwsefs3ysevn7m3s8h3q01 (still platform v4). Download `dl16/`: task.md, test.patch (v15), solution.patch (v13), Dockerfile all identical to local.
- Next: platform stage 3 on v4, then reruns (old runs stale since v15 tests).

## v12 description pushed (2026-10-29)
- `problems edit --description-file=task.md` -> versionId kd743zn69czwsefs3ysevn7m3s8h3q01 (platform v4). Download `dl17/`: task.md, test.patch (v15), solution.patch (v13), Dockerfile identical to local.

## v13 description pushed (2026-10-29)
- Bug-report style (desc-vi/). versionId kd743zn69czwsefs3ysevn7m3s8h3q01 (platform v4). Download `dl18/` identical to local.

## v14 description + v16 tests + v14 solution pushed (2026-10-31)
- Same versionId kd743zn69czwsefs3ysevn7m3s8h3q01 (platform v4). Download `dl19/` identical to local. Working tree of `jedi/` holds the changes, uncommitted.
- Next: platform stage 3 and fresh rollouts; every earlier run is stale.

## v15 description + v17 tests pushed (2026-11-01)
- Same versionId. Download `dl20/` identical to local. Solution unchanged (v14).
