# Probing passing runs against the reference

## Reconstruct a run
- Take the run's workspace diff and apply it to a clean checkout at the pinned commit. Drop hunks that touch the test files, the test runner, or harness directories shipped inside the diff; apply the current test patch afterwards. Grade in the task's image, not on the host, because library and submodule versions differ.
- Use the same junit path and the same dedupe rule for duplicate names as the test runner does.

## Read the diff for non-solutions
- Edits to existing tests or expectations (a weakened existing test hides a regression).
- Special-casing test inputs, hard-coded outputs, detecting the test harness, skipping or xfailing tests, changing runner or config.
- Partial implementations that pass only because the tests exercise a subset (one syntactic form, one platform branch).
- Reuse of a repo helper that silently drops a construct the prompt mentions; each drop is a distinguishing test to add.

## Probe matrix
- One row per prompt sentence, several variants each: every syntactic form, every kind of input source (literal, variable holding a literal, annotated fixed-size, annotated unbounded, unresolvable), nested and combined cases, empty/missing/nil values, both directions of two-way rules, and every parallel code path.
- Run every probe on (a) the reference, (b) each passing run, (c) the real language/runtime as ground truth. Isolate each probe (own recover/try-catch or own process) so one crash doesn't hide the rest. Vary seeds and orders that matter.
- Print the full output; don't grep for expected lines or you'll hide the wrong rows.
- Check that each expectation is achievable (an impossible expectation is the probe's bug, not the run's).
- Put results in JSON and run `scripts/compare_probes.py results.json`.

## Who is wrong
- Reference wrong, run right: fix the reference.
- Run wrong, reference right, rule in the prompt: FP, add the distinguishing test.
- Both wrong: pick the risk and record it.
- Rule not in the prompt: unfair probe, don't test it; consider a general clause only if the budget allows.
- A probe that depends on something unknowable statically (an unpredictable branch condition) while the task is about static behavior is unfair too, as is one that the reference itself fails.

## Make new tests useful
- Add variants where a name or cache could produce a coincidental match.
- A test that passes only inside the suite (order, cache) is a bad test; run it alone and with several seeds.
- Keep assertions inside the prompt: no exact error types, formatting idempotency or inputs the prompt never mentions.
