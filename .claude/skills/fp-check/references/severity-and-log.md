# Severity of findings and the check log

## Severity (every finding gets exactly one)
| Level | Meaning | Examples |
|---|---|---|
| HIGH | A false positive: a run passes the tests but differs from the reference on a behavior the prompt states, or it tampers with tests/runner. Must be closed with a new test. | run passes yet gives a different result on a prompt-stated rule; a run edits or weakens an existing test; a run special-cases test inputs |
| MEDIUM | Probable FP or a weakness, not proven: a run passes with an implementation that is partial on a rule the tests barely touch, or differs from the reference only on a rule the prompt implies through a general sentence; the reference itself is wrong where a run is right. | partial implementation that happens to pass; reference less accurate than a passing run on an in-scope probe; a test that passes only in the suite (order/cache) |
| LOW | Advisory: differences on behavior the prompt does not state, or noise. Do not add tests for these. | niche probe differences outside the prompt; a run that edits fixtures only by appending lines |

Rules: pick the lower level when unsure and say why; a finding without evidence (command, file:line, probe id, quoted sentence) is not reported; a probe whose expectation was not checked against the language's real behavior is at most LOW.

## Finding format
`[LEVEL] run <id> — <one-line defect> — evidence: <test/probe/diff hunk> — fix: <test or prompt change>`

## The check log
Each task keeps its own log as `CHECK-LOG.md` in the task's folder (the folder that holds its prompt, tests and solution). Every stage of the task's pipeline appends one entry per run, newest last, including this skill. Create the file from the header below if it doesn't exist. Never mix tasks in one file. Never rewrite earlier entries, except to fill the lines left empty on purpose (`Reference result`, `Agreement`) once an independent result is known.

Round: 1 for the first pass of a task; +1 each time the task is regenerated after a failure. Read the log first: the entries of the current round are the findings you must not repeat, and earlier rounds show what already failed.

```
### [R<round>] fp-check — <version id> — <YYYY-MM-DD HH:MM>
- Result: PASS | FAIL | PARTIAL (say which checks did not run)
- Inputs: <prompt file, test patch, solution patch, versions, runs>
- Run: <what ran / did not run and why>
- Findings:
  - [HIGH] run <id> — ... — evidence: ... — fix: ...
  - [MEDIUM] ...
  - [LOW] ...
- Counts: HIGH <n>, MEDIUM <n>, LOW <n>
- Next: <pass to the next stage | back to the task generator because ...>
- Reference result: <empty until known: the result an independent reference (another reviewer or a service) gave for the same checks, with its own severity if it has one>
- Agreement: <empty until known: per finding, match / skill-only / reference-only / false alarm>
```

Result is FAIL if any HIGH finding exists, PARTIAL if a check could not run, otherwise PASS.

Header of a new task log: `# Check log — <task name>`, one line with repo/commit/language and task folder, then the table `| Vòng | Stage | Ngày | Version / ca | Kết quả | HIGH | MEDIUM | LOW | Đối chiếu |` under `## Tổng kết theo thời gian` (copy another task's `CHECK-LOG.md` for the exact layout), then `## Entries`. Add one table row per entry.

Keep an entry short; put long probe output in a file next to the submission and link it.
