# v7: dict overwrite rule in the prompt + stable baseline ids (2026-10-22)

Pushed to kh7d8qtdcfr0k5mt1rvftcvebn8gb2va (v2 draft) with `olympus problems edit --description-file task.md --test-patch test.patch`.

Why:
- v2 runs: 0/10 solved (Vega 0/5, Orion 0/3, Nova 0/2). 8/9 failed only the assertion in
  test_items_of_mappings_and_other_containers: `f({**({"k": A} if random.random() else {}), "k": C})` -> RB.
- The grading wrapper matches baseline by node id. Integration test ids contain line numbers
  (`functions:113`), and agents edit test/completion/functions.py and pep0484_overload.py, so ids shift
  and are reported as missing (17, 30, 42 false baseline failures across runs).

Changes:
- task.md: after "With dict, both key and value must be accepted." added
  "The entries of a dict literal are the ones it holds once built: a key written later replaces an
  earlier entry with the same key, whether that entry was written directly or came from `**`."
- test.sh base: `git -c safe.directory=/app checkout HEAD -- test` before pytest (only when HEAD exists;
  warns and continues on failure). test.patch touches no tracked files; solution.patch touches nothing in test/.

Local verification (image built from clean 2af1a0f3 with .git, --network none):
- plain / fixtures edited / solution / solution+edited / non-root uid 1000 (ulimit -s unlimited):
  base = 3981 ids, 0 failures, identical id set in all scenarios; new = 0/63 without, 63/63 with solution.
- Control with the old test.sh + edited fixtures: 116 ids missing.
- Reference solution on overwrite cases: {"k": A, "k": C} -> RB, {"k": C, "k": A} -> RA, {**{"k": A}, "k": C} -> RB,
  {"k": A, **{"k": C}} -> RB, {**{"k": A}, **{"k": C}} -> RB, conditional ** then "k": C -> RB,
  "k": C then conditional ** -> RC. All match the new sentence.
