---
name: fp-check
description: False-positive check for agent runs that passed a task's hidden tests. Given the task prompt, the tests, the reference solution and the agents' patches (runs), find the runs that pass without genuinely solving the task, and strengthen the tests so they can't. Use whenever the user supplies agent runs/patches/diffs together with a prompt, tests and solution and asks for "FP check", "false positive", "run pass nhưng sai", "run này có thật sự giải đúng không", "test có quá lỏng không", "siết test", or wants to know whether passing runs are legitimate — even if they don't say "false positive".
---

# FP check

A false positive (FP) is an agent run that passes the tests without really solving the task. It means the tests are too permissive: the agent found a gap, and a fake pass looks identical to a real one. The fix is always in the tests (and, if needed, one general clause in the prompt), never in the agent's patch. A prompt change goes back through the task generator; new tests are written by the task's author, and every earlier result is out of date afterwards.

Input: the task prompt, the tests (`test.patch` with its runner), the reference solution, the environment (Dockerfile + repo at the pinned commit), and the agents' runs: one patch/diff per run, normally only the passing ones. Talk to the user in the language they use (often Vietnamese). Work autonomously. If the tests, prompt or solution are themselves in doubt, run `quality-check` first.

## Workflow

1. **Collect the runs.** One folder per run holding its patch (any `*.patch`/`*.diff` of the workspace). Only passing runs matter, but a failing run is useful in step 5.

2. **Grade every run on the current tests**: `scripts/grade_run.sh <submission-dir> <run-dir>...`. It rebuilds each run on a clean checkout (the run's patch minus anything touching the test files or test runner), applies the current test patch, and runs base and new in an isolated container. It prints, per run, which tests fail. A run that passed earlier but fails now was graded on weaker tests, so it's not a legitimate pass. A run that still passes goes to step 3.

3. **Audit the diff**: `scripts/diff_audit.py <run-dir>...` flags edits to existing tests, skip/xfail markers, runner or config changes, and hard-coded expectations. Then read the diff itself for partial implementations (one syntax form, one branch), special-casing of test inputs, and reuse of a helper that silently drops a construct the prompt mentions.

4. **Probe against the reference.** Follow `references/probing.md`: write a probe matrix from the prompt (every sentence, every equivalent form, every input source, both directions of two-way rules, empty/missing values, every parallel code path), run it on the reference and on each passing run, and compare with the language's real behavior. `scripts/compare_probes.py` lists the probes where a run differs from the reference. Decide who is wrong: reference wrong → fix the reference; run wrong on a rule the prompt states → FP; both wrong → note the risk; rule not in the prompt → not a fair probe.

5. **Strengthen the tests.** For every FP: add a test that fails for the run's behavior and passes for the reference.
   - The test may only assert what the prompt states or a developer would find in the repo. If the distinguishing behavior isn't stated, add one general clause to the prompt (staying within the clause budget of `quality-check`), not a sentence per edge case.
   - Look at the failing runs too: their typical mistakes usually exist silently in the passing ones.
   - Prefer a few tests that each make one or two runs fail over one test that kills them all. A test that may crash the whole binary should recover and report a normal failure. Run new tests alone and inside the suite, with different seeds, so order and caches can't hide or fake a result.

6. **Re-grade everything** (step 2) with the new tests: every FP must now fail, the reference must still pass, and no fair run should fail for an unfair reason. If strengthening drops every run to zero, the tests may have become unfair or too hard; re-check with `quality-check` before keeping them.

## Report format

Give every finding a severity HIGH / MEDIUM / LOW and an evidence line (definitions and finding format in `references/severity-and-log.md`). Append the entry described there to the `CHECK-LOG.md` in the task's folder. Per run: verdict (legitimate / stale pass / differs from reference on rule X / tampering / shortcut), evidence (failing test, probe id, diff hunk), and the test or prompt change that closes it. Then the list of new tests, the re-grade table, and what was not verified (for example, no container available).
