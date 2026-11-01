# jedi-overload-arguments (draft kh7d8qtdcfr0k5mt1rvftcvebn8gb2va)

- Repo: https://github.com/davidhalter/jedi
- Commit: 2af1a0f3943f801f5de755f8ffd251f0385a3607
- Issue: https://github.com/davidhalter/jedi/issues/2028
- Title: Overload resolution ignores type[...] and *args/**kwargs annotations
- Language: Python. Category: Bug Fix
- Description: `description.md` (from `description_vi.md`, gtx), copied to `task.md`

## Clause list (description v11)
1. choose the overload from the types of the call's arguments, incl. `type[...]` (issue example)
2. first overload accepting all arguments
3. positional or by name (`*args`/`**kwargs` binding is standard Python)
4. no overload accepts: combine all results
5. annotations follow normal typing rules
6. a multi-valued argument: every possibility must be accepted (examples `f(A if c else C)`, `f(*xs)` with `xs: list[type[C]]`)
7. literal tuple/list/dict or a variable holding it, even unpacked or nested: use runtime content
8. equal keys: the later value replaces the earlier one (`1`/`True`)
9. possibly absent `**d` entry: both values possible
10. comprehension over an empty literal: no elements
11. `*xs` possibly empty shifts later positional arguments

## Regeneration
### v11 (2026-10-27), from platform v4 descriptionQuality FAIL (ai_generated, major x2)
- "too prescriptive ... write it like how you would prompt an LLM" on paragraph 2:
  opened with the ask; the typing enumeration (TypeVar bounds/constraints, Optional,
  containers, every element/key/value, `*args`/`**kwargs` check each) is replaced by
  "Annotations follow normal typing rules". verifyFairness v4 already rated the
  related tests "Standard external semantics".
- Same finding on paragraph 3: the chain of rules is replaced by one principle
  (runtime content of literals) plus short examples with outcomes. Every behavior
  that verifyFairness v4 cited as prompt-stated (about 12 tests) keeps an example.
- Previous: `description_vi.v10.md`, `description.v10.md` (= `task_v10.md`).
- Tests and solution: no change needed.

### v11 wording rounds (2026-10-27), local quality-check R3
- Merged unknown-length `*` into the multi-value rule (MEDIUM, test_unpacked_values_of_unknown_length).
- New check-7 AI-detector rule: paragraph 2 rewritten as one principle + two examples with results; dropped "the same way a type checker does" (broader than the tests) and the `*args`/`**kwargs` construct list.
- Previous wordings: description_vi.v11a.md, v11b.md, v11c.md.

### v12 (2026-10-29), from platform v4 descriptionQuality FAIL (ai_generated, major x3 on paragraphs 1-2)
- Flagged: the opening sentence, the first-overload/fallback sentence, and "Annotations follow normal typing rules. In the examples below ..."
- Paragraph 1: `@overload` context + issue example as a fenced code block. Paragraph 2: short imperative requests to the solver instead of repeated "Jedi should ..." sentences. Paragraph 3 unchanged (not flagged). Clause list unchanged.

### v13 (2026-10-29), from platform v4 descriptionQuality FAIL (ai_generated, major x3 on paragraphs 2-3, v12)
- Whole prompt rewritten as a first-person bug report: symptom + issue example + expected result; "I think Jedi should ..."; "Some examples I have tried: ...". Vietnamese in `desc-vi/description.vi.md`, raw gtx in `desc-vi/description.gtx.md`. Clause list unchanged.

### v15 (2026-11-01), from platform v4 FAIL (too many requirements; test_first_of_two_equal_keys_stays not fair)
- verifyFairness counted 11 requirements on v14, 9 on v13 (PASS). One TypeVar example (`T = TypeVar("T", "A", int)`) now covers the string reference and the shared constraint. The comprehension sentence moved into the general rule ("tuple, list, dict, or comprehension"). Back to 9 requirements, 293 words.
- "actual content, both key and value" makes the equal-key object test follow from a stated rule without naming the rule that Nova misses.

