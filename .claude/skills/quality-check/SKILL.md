---
name: quality-check
description: Review a coding-task submission made of a task prompt, hidden tests and a reference solution, and report whether it passes seven quality checks - tests fail without the solution, solution passes everything, no flakiness, tests are fair, prompt and tests agree, solution is complete, prompt is clean. Use whenever the user provides (or points at) a prompt/description, test patch and solution patch and asks for a "quality check", "check chất lượng", "verify tests/solution", "flaky", "test có fair không", "prompt có ổn không", "solution có đủ không", or wants to know if a task is ready, even if they don't say "quality".
---

# Quality check

Input: a task prompt (description), the tests (usually a patch that adds a test runner script and test files), the reference solution (patch), and an environment to run them (Dockerfile + repo at a pinned commit). Output: a verdict per check with evidence, and concrete fixes.

Talk to the user in the language they use (often Vietnamese); the task prompt itself stays English. Work autonomously: choose the recommended option and continue.

## The seven checks

| Check | Passing means |
|---|---|
| 1. Tests fail without the solution | on the clean repo the existing tests pass (`base`) and every new test fails (`new`) |
| 2. Solution passes everything | with the solution applied, base and new both pass fully; no existing test broken |
| 3. No flakiness | repeating identical runs, with and without the solution, never changes any test's result |
| 4. Tests are fair | every assertion checks something the prompt states or a developer would find in the repo |
| 5. Prompt and tests agree | everything the tests assert is in the prompt, and every prompt requirement is tested |
| 6. Solution is complete and clean | the solution does what the prompt promises in every equivalent form, without unrelated changes |
| 7. Prompt is clean | behavior not method, no redundancy or discoverable detail, not too many requirements |

Details, typical failures and fixes for each check: `references/checks.md`. Language-specific details (how to repeat runs, wildcard pitfalls, environment pins): `references/languages.md`. Read the sections that apply.

## Workflow

1. **Collect the pieces and the language.** Folder with `Dockerfile`, prompt (`description.md`/`task.md`), `test.patch`, `solution.patch`, and the repo checked out at the pinned commit. If something is missing, say what, and still do every check that doesn't need it.

2. **Static gate** (seconds): `scripts/static_gate.sh <dir>`. Checks prompt size, forbidden commands in the test runner (installs, error-swallowing), hash suffixes on new test paths, Dockerfile pinning, and that the patches apply to a clean checkout. Fix every FAIL first.

3. **Run checks 1–3**: `scripts/verify.sh <dir> [--repeat N] [--new-pkg PATTERN]`. It builds the image from a clean worktree and runs tests-only and solution runs in base and new modes in an isolated container (unprivileged user, unlimited stack, no network), repeated N times. `scripts/junit_check.py` then checks the contract on the junit files: duplicate test names, base/new overlap, new tests passing without the solution. If docker or the repo is unavailable, say checks 1–3 were not run instead of guessing.

4. **Judge checks 4–5 by reading.** For each assertion in the tests, name the prompt sentence or repo fact that justifies it; flag any with none. For each prompt sentence, find the test that covers it, in both directions of any two-way rule. Test by test, which new tests would still pass for a plausible wrong implementation (mutants)?

5. **Judge check 6 by probing, systematically.** Treat the prompt as a contract. Follow the numbered procedure in `references/checks.md` check 6: write the full table of constructs × source kinds to a file, probe every cell on the solution and on the language's real behavior, and grep every caller of any helper that returns extra information (flags, alternatives). Sampling a few probes is how real defects are missed. For a large change give this to an independent subagent that sees only the prompt and the solution, so it isn't anchored on your reading.

6. **Judge check 7.** First name the feature groups, then count requirement clauses (≈15 at most, ≈330 words at most), and apply the rules in `references/checks.md`. Flag a prompt only for what the rules actually support; over-flagging a clean prompt wastes the author's time as much as missing a real problem.

7. **Fix the whole class of a problem, not only the cited line.** Findings about the prompt (checks 4, 5, 7) go back to the task generator, which rewrites the description; findings about tests and the solution are fixed by their author. Either way a fix makes every earlier result out of date. Re-run from step 2 after every edit, because a change to prompt, tests or solution can break another check. If a review of the solution found one gap, expect sibling gaps and look for them.

## Report format

Give every finding a severity HIGH / MEDIUM / LOW and an evidence line (definitions and finding format in `references/severity-and-log.md`). Append the entry described there to the `CHECK-LOG.md` in the task's folder. Then a table: each of the seven checks, status (PASS / FAIL / not run), evidence (command, file, probe), next action. Then the list of fixes in priority order. State what was not run and why. Don't say "ready" unless all seven pass on the final versions of the files. Result in the log is FAIL when any HIGH finding exists; the next step is then "back to the task generator".
