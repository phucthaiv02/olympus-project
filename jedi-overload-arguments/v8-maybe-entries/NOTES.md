# v8: tests for may-be-missing entries from ** (2026-10-23, NOT pushed yet)

Same task.md / solution.patch / Dockerfile as v7. Only test.patch changes (3 assertions, still 63 tests):
- test_items_of_mappings_and_other_containers:
  `f({"k": C, **({"k": A} if random.random() else {})})` -> RC
  `f({**{"k": A}, **({"k": C} if random.random() else {})})` -> RC
- test_dict_unpacking: `f(**{"p": C, **({"p": A} if random.random() else {})})` -> ['RA', 'RB']
Basis in the prompt: the later-key-replaces sentence + "every element ... must be accepted" + "If an argument can
have multiple types ..., the annotation must accept all of those types".

Why: v7 runs gave Vega 2/4 (too easy). Probe of 26 candidate assertions on the solution patches of the
8 v7 runs (pp/probe.py in the session scratchpad):
- vegaP1 (rd7ehqadn52s) passes all candidates; vegaP2 (rd7638agy1ma) treats a maybe-present later ** entry
  as a definite overwrite (fails D1/D4/D14/D15/K5/K6).

Results with the v8 test.patch (new tests):
| patch | v7 | v8 |
| ref | 63/63 | 63/63 |
| Vega rd7ehqadn52s | pass | 63/63 |
| Vega rd7638agy1ma | pass | 61/63 |
| Vega rd7ay1kcr99h, rd7fmta5maf4 | 62/63 | 61/63 |
| Nova rd7bg8qwjy7c | 62/63 | 62/63 (forward-ref TypeVar) |
| Nova rd750h2ws0te | 62/63 | 62/63 |
| Orion rd79s1wssrdw | 58/63 | 58/63 |
| Orion rd7a6hxz88d9 | 61/63 | 60/63 |
Stable under PYTHONHASHSEED 0,1,2,3,4,5,42,1234 for ref, vegaP1 and vegaP2.
Observed Vega pass rate would be 1/4 = 25%.
