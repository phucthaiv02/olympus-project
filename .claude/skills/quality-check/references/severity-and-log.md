# Severity of findings and the check log

## Severity (every finding gets exactly one)
| Level | Meaning | Examples |
|---|---|---|
| HIGH | Blocks the task: it would fail a check, or the contract is broken. Must be fixed before the task is used. | a new test passes without the solution; the solution breaks an existing test; a test is flaky; a test asserts behavior the prompt never states and no general rule implies; the solution is wrong on a behavior the prompt promises (any equivalent form); the prompt has three or more feature groups or is clearly over the clause budget |
| MEDIUM | A reviewer would likely flag it, or it weakens the task, but the contract still holds. Fix before re-submitting anything else. | a fixed limit or fallback that approximates instead of failing closed; an extra feature in the solution that the prompt does not ask for; a plausible mutant no test catches; 16–18 clauses; a sentence that restates the bug or repeats a rule |
| LOW | Polish or advisory. Fix only if the file is being edited anyway. | wording, a missing comment, an unpinned base image tag, lint nits |

Rules: pick the lower level when unsure and say why; a finding without evidence (command, file:line, probe id, quoted sentence) is not reported; a probe whose expectation was not checked against the language's real behavior is at most LOW.

## Finding format
`[LEVEL] check <n> — <one-line defect> — evidence: <where> — fix: <what>`

## The check log
Each task keeps its own log as `CHECK-LOG.md` in the task's folder (the folder that holds its prompt, tests and solution). Every stage of the task's pipeline appends one entry per run, newest last, including this skill. Create the file from the header below if it doesn't exist. Never mix tasks in one file. Never rewrite earlier entries, except to fill the lines left empty on purpose (`Reference result`, `Agreement`) once an independent result is known.

Round: 1 for the first pass of a task; +1 each time the task is regenerated after a failure. Read the log first: the entries of the current round are the findings you must not repeat, and earlier rounds show what already failed.

```
### [R<round>] quality-check — <version id> — <YYYY-MM-DD HH:MM>
- Result: PASS | FAIL | PARTIAL (say which checks did not run)
- Inputs: <prompt file, test patch, solution patch, versions>
- Run: <what ran / did not run and why>
- Findings:
  - [HIGH] check 6 — ... — evidence: ... — fix: ...
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
