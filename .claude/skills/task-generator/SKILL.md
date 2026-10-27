---
name: task-generator
description: Generate a coding task from a real GitHub repo and open issue - repo, pinned commit, issue URL, title, category and a task description that is written in Vietnamese first and then translated to English through a translation API. Finds a new repo and issue only when the previous task was rejected (duplicate, already solved, declined, ineligible repo, too small); otherwise regenerates the task and its description from the findings of earlier checks and reviews. Use whenever the user wants to create a new task, asks to "tạo task", "generate task", "tìm repo", "tìm issue", "tìm kiếm repo phù hợp", "mô tả task", "viết task prompt", "viết description", "viết lại description", "task bị reject", "sửa task theo lỗi", wants task ideas hard enough for strong coding agents, or wants to check whether a repo or idea is eligible (stars, license, activity, existing PRs, duplicates) - even if they don't say "generate".
---

# Task generator

Produce a task: one pinned repo commit, one real open issue, a title, a category, and a task description. You stop before tests, solution and environment; the `quality-check` and `fp-check` skills review those afterwards, and their findings are the input for regenerating.

## Where this sits in the pipeline

This is the first stage of a pipeline: generate, then several automated and external checks, runs, reviews and finally submission. A failure at **any** later stage sends the task back here. So you are entered in two ways: with nothing (a new task), or with a task whose latest stage failed.

When entered with an existing task:
1. Read the task's `CHECK-LOG.md` (in the task folder) completely. The newest round's entries are the findings to resolve; older rounds show what was already tried, so do not reintroduce a rejected idea or an earlier problem. Findings can come from any stage (automated checks, agent runs, false-positive checks, reviews by a tool or a person), each with a severity and evidence.
2. Resolve every HIGH finding in the new version. Resolve MEDIUM findings unless they conflict with a stricter finding (say so). Touch LOW findings only where you are editing anyway.
3. Choose Mode A or B (next section) from the findings, then write the result.
4. Append your own entry to the log: round (previous + 1), stage "task-generator", the mode, each finding with the change made for it (or why it was not addressed), the new version id, and which tests or solution parts must change to match. The entry's Next line is "quality check, from the start": every earlier result is out of date once the description changes.

When entered with nothing, create the log with a first entry (round 1) after you finish.

Talk to the user in the language they use (often Vietnamese). The description is always written in Vietnamese first and then translated to English; all other text may be in the user's language.

## Why the bar is high

A good task is hard for strong coding agents and still fair. Hard means the change crosses several layers of a large codebase, so the solver has to explore a lot, and it has subtle semantics that a careless implementation gets wrong. Fair means a developer who reads only the description and the repo can implement it without guessing, and everything a test could check is stated or follows obvious repo conventions. Most bad tasks fail in one of two ways: someone already did the idea (a PR, a merged change, a declined request), or the change is too small or too broad. Spend effort on those two.

## Decide the mode first

Ask (or infer from what the user gave you) whether a task already exists.

| Situation | Action |
|---|---|
| No task yet | Mode A: find a repo and issue, then write the task |
| A task exists and was **rejected** | Mode A again, with a different repo or idea |
| A task exists and has **findings** (review notes, quality-check or fp-check output, a check log) but is not rejected | Mode B: regenerate the description (and, if needed, the scope) from the findings. Do not search for a new repo |

Rejected means the idea itself cannot be saved:
- an existing PR (open, merged or closed) implements it, or the capability already ships;
- maintainers declined it or called it out of scope;
- the new task is a close variant of earlier work in the same repo, or its core can be ported line by line from a public reference implementation;
- the repo is ineligible (inactive, archived, too small, non-permissive license) or the idea does not fit the project;
- it is irreparably too small (lives in one or two files) with no deeper version of the same issue available.

Everything else is fixable and goes to Mode B: a description that is too long, too prescriptive or has too many requirements; tests that assert something the description does not say (fix the description or the scope, not by adding sentences per edge case); a solution gap that exposes an unstated rule; a false positive that needs one more stated rule; a task that is too easy but can be made deeper within the same issue; too hard because the scope is broad.

When unsure whether a finding is a rejection, say which reading you take and why, and go with the cheaper one (Mode B).

## Mode A: find a repo, an issue and an idea

### A1. Constraints
If the user hasn't said, ask briefly for the language (TypeScript, JavaScript, Python, Go, Rust, C++ or Java), a domain they know, and repos to avoid. A language is enough to start; don't over-interview.

### A2. Candidate repos
Use `gh search repos` (for example `gh search repos --language=go --stars=">1500" --sort=updated --limit 40 "<topic>"`) together with your own knowledge of well-architected projects. Aim for **depth**: several layers that a single feature must flow through. Good shapes: language and DSL implementations (lexer, parser, analyzer, evaluator, formatter, linter); query engines (parser, planner, executor); compilers, bundlers, transpilers; schema and validation libraries; serializers with code generation; protocol implementations; CLI tools with config, core and output layers. Avoid thin repos (wrappers, UI kits, small utilities); they cap difficulty no matter how good the idea is.

Run the eligibility script on every candidate:

```bash
.claude/skills/task-generator/scripts/repo_check.sh owner/repo owner/repo2 ...
```

It reports stars, language, license, archived flag, last default-branch commit date and SHA, and flags anything below the floor: 500 or more stars, a commit in the last 12 months, a supported language, a permissive license (MIT, Apache-2.0, BSD-2/3, ISC, 0BSD, Unlicense, Zlib are built in). If a repo shows another license, tell the user to confirm it rather than silently dropping or accepting it.

### A3. Find issues and turn them into ideas
Every task is anchored to a real GitHub issue; its URL is part of the task. An open issue written by a user shows the idea fits the project and isn't done yet. Don't present an idea without one.

For each repo that passes, read the README and skim the architecture (`gh api repos/O/R/git/trees/HEAD?recursive=1`, key files through `gh api .../contents/...` or a shallow clone into scratch space). Then list open issues, most-reacted first, and keep those asking for a substantial feature or a deep bug fix:

```bash
.claude/skills/task-generator/scripts/issue_check.sh find owner/repo
.claude/skills/task-generator/scripts/issue_check.sh find owner/repo "label:enhancement"
.claude/skills/task-generator/scripts/issue_check.sh show owner/repo 123
```

`show` prints state, labels, reactions, **every PR that references the issue in any state**, and maintainer comments. Drop the issue if any PR already implements it (open, merged, or closed unmerged) or a maintainer turned it down. An `enhancement`/`help wanted` label, or a maintainer saying "PRs welcome", is a strong plus.

The issue gives the idea and the user-visible behavior. The task can be a precise, fully specified version of it; narrow or sharpen the scope, as long as it still clearly solves what the issue asks. Don't stretch an issue into something unrelated just to have a link.

A good task is one coherent change that:
- **fits the project's philosophy**: something a maintainer would plausibly accept;
- **cuts across several layers** (parser, analysis, runtime, formatter...); a change in one or two files is solved by any strong agent;
- **has subtle, testable semantics**: interactions with existing features, null and empty handling, ordering, error versus result cases. A sentence in the description that makes the subtle rule the contract is what makes a failure fair;
- **is one thing**, not a bundle. A family of closely related behaviors sharing one mechanism is fine; unrelated asks are not, and reviewers count independent feature groups;
- is **fully verifiable offline** (no network, timing or randomness).

### A4. Duplicate and decline check (every idea you present)
`issue_check.sh show` only finds PRs linking the issue. Search by keywords too:

```bash
.claude/skills/task-generator/scripts/dup_search.sh owner/repo "keyword" "synonym" "other phrasing"
```

Search several phrasings (common name, synonyms, the syntax or operator, the API name). It covers PRs in all states, issues and Discussions, where design rulings often live. Search also by the mechanism and the functions the change would touch, not only by the feature's name, because a PR can solve most of an idea without ever mentioning it. Check the changelog and docs to be sure the capability isn't shipped or deliberately removed. Record what you searched and found.

### A5. Present 3–5 candidates and wait
Present a compact comparison, then **stop and wait for the user's pick**. Per candidate: repo (stars, language, license, last commit); issue (URL, title, reactions, labels, what maintainers said); the idea in one or two sentences; layers touched; why it is hard; duplicate check (keywords and result); risks; a difficulty guess. Flag weak candidates honestly; an idea that is probably too easy is worth flagging, not hiding.

## Mode B: regenerate from findings

1. **Collect the findings**: what the user pasted, the review notes, the task's check log (`CHECK-LOG.md` in the task folder) and any quality-check or fp-check output. List each finding with its severity (HIGH, MEDIUM, LOW) if it has one.
2. **Map each finding to a description change** using the table below; resolve conflicts between findings by keeping the stricter requirement.
3. **Re-read the issue and the repo** only as far as needed to apply the changes. Don't change the repo, commit or issue unless a finding says so (that is a rejection, so Mode A).
4. Go to the description pipeline and save the new version without deleting the previous one.

Findings may also come from runs and reviews rather than checks. Read those as follows: a review saying the idea is a duplicate, derivative, declined upstream or out of scope is a rejection (Mode A); a review about wording, precision, unfairness or scope size is Mode B; agent runs that all pass mean the task is too easy, none ever pass means too hard (see the last two rows).

| Finding | Change |
|---|---|
| too many requirements / feature groups | drop a whole feature group from the description (and tell the user the tests and solution must drop it too); or merge related behaviors into one general rule |
| prescriptive, names internals, describes the algorithm | rewrite as behavior only |
| restates the current bug, repeats a rule, says what stays unchanged | delete the sentence |
| a test asserts something the description never states | decide: state it with one general rule if it is a natural consequence of the feature; otherwise tell the user to drop that assertion. Never add a sentence per edge case |
| description promises more than the tests or solution deliver ("any form that X supports") | narrow to what is actually tested |
| a passing run (false positive) differs on a rule | state that rule in the description, usually by widening an existing general sentence |
| a solution gap shows an ambiguous rule | make the rule explicit in one sentence, with the same wording style as its neighbours |
| too easy, every strong agent solves it | go deeper on the same issue: more demanding semantics of the same rule, a second layer the feature must reach; do not add unrelated features |
| too hard, nobody can solve it | remove the rarest rule or group, or state the hardest rule more precisely |

## The description pipeline (both modes)

The description is always generated in Vietnamese, then translated to English with the API script. Writing English directly tends to read as machine-written and drifts into long, over-specified, prescriptive text; writing in Vietnamese first keeps it natural, short and focused on behavior.

1. **Write `description_vi.md`** in natural Vietnamese, following the rules below.
2. **Translate and gate**: `python3 .claude/skills/task-generator/scripts/make_description.py description_vi.md --out description.md`. It translates paragraph by paragraph through the public translation endpoint, protects inline code spans, writes the English file and prints a gate: word count, rough clause count and flagged phrases. Needs network access.
3. **Read the English** for meaning errors (technical terms, negations, operator names, "any/every" scope). Fix the meaning by editing the Vietnamese and translating again; don't hand-edit the English, because the two versions must stay in sync. If a term cannot be translated well, put it in backticks in the Vietnamese.
4. **Count requirement clauses by hand** and list them. Stay at about 15 clauses or fewer and about 330 words or fewer. A general rule that covers many cases counts once; a list of cases counts each item.
5. **Self-review** as a skeptical reviewer and as a solver seeing the repo for the first time: could a competent developer implement this without guessing? Where would two reasonable implementations differ, and does the description decide it? Is anything prescriptive? Is it one coherent change? Does it still clearly solve the linked issue? Treat every sentence as a promise: a reviewer will probe the reference solution with a variant of each one, so expect to hand the description to `quality-check`.

### Rules for the description
- Open with the ask itself ("Add X to Y", "Fix Z when ..."). The first line must stand without the title. No motivation preamble and no paragraph about current broken behavior.
- Natural prose in full paragraphs. No bullet lists, headings or code blocks doing the describing. Inline code for syntax or names users write is fine. No URLs; the issue link has its own field. Don't copy the issue's wording.
- Describe **behavior, not implementation**. Don't name internal classes, helpers, files or fields a developer would discover, and don't hint at the approach.
- State every load-bearing detail of the contract: edge cases that tests will check (null, empty, errors versus results, precedence, interaction with existing features, accepted types, both directions of a two-way rule). A sentence such as "in all other respects it behaves exactly like X" pins many edge cases fairly.
- Don't state more than the tests and the reference do: no "in any form that X supports" unless every form is tested. Don't leave out anything a test asserts.
- Don't restate what must keep working, don't say "other cases behave as before" or "X is not part of this change", and don't repeat a rule as an example.
- Answer reviewer feedback by merging edge cases into one general rule that every case follows ("the contents are what the language builds at run time; every possibility must be accepted"), never by adding one sentence per edge case.
- Avoid terms whose technical meaning is broader than the task ("same memory layout", "any call form"). Don't promise exact error messages or formatting unless tests check them.
- No mention of tests, hidden tests, benchmarks or agents.

## Output

Create a folder named after the task (kebab-case slug) and write:

- `description_vi.md`: the Vietnamese source
- `description.md`: the English translation
- `TASK.md`: repo URL, full 40-char commit SHA (the latest default-branch commit from `repo_check.sh` unless the user specifies one), issue URL, title (short, issue-style), language, category (Feature Request, Bug Fix...), the English description, and research notes: duplicate-check keywords and results, layers expected to change, behaviors the tests should cover (one line each, as a sanity check, not a test plan), risks, and the clause list
- when regenerating: keep the previous files as `description_vi.vN.md` and `description.vN.md`, and add a "Regeneration" section to `TASK.md` that lists each finding and the change made for it
- the log entry in `CHECK-LOG.md` described at the top (see `references/severity-and-log.md` of the `quality-check` skill for the entry format)

Finish by telling the user where the files are and the next steps: tests with a runner script, the reference solution and the environment, then `quality-check` (and `fp-check` once agent runs exist). After a regeneration, remind them which tests or solution parts must change to match.
