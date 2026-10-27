#!/usr/bin/env bash
# Static gate for a task submission folder (prompt, test.patch, solution.patch, Dockerfile). Prints PASS/WARN/FAIL lines, exits 1 on any FAIL.
set -u
D="${1:?usage: static_gate.sh <submission-dir>}"
cd "$D" || exit 2
fail=0
f() { echo "FAIL  $*"; fail=1; }
w() { echo "WARN  $*"; }
p() { echo "PASS  $*"; }

desc=""
for c in description.md task.md; do [ -f "$c" ] && desc="$c" && break; done
if [ -n "$desc" ]; then
  words=$(wc -w < "$desc")
  # rough clause estimate: sentences plus semicolon/colon-separated parts
  clauses=$(tr '\n' ' ' < "$desc" | grep -oE '[.;:][ ]' | wc -l)
  if [ "$words" -gt 330 ]; then f "$desc has $words words (limit 330; >330)"; else p "$desc has $words words"; fi
  echo "INFO  ~$clauses sentence/clause breaks in $desc; count requirement clauses by hand, keep <= 15"
  if grep -qiE 'other cases (keep|remain)|unchanged|currently (fails|returns|throws)|not part of this' "$desc"; then
    w "$desc contains phrasing commonly flagged by prompt reviewers (unchanged behavior / current bug / out-of-scope sentence)"
  fi
else
  f "no description.md or task.md"
fi

for x in Dockerfile test.patch solution.patch; do [ -f "$x" ] || f "missing $x"; done

if [ -f test.patch ]; then
  # test.sh lives inside test.patch for these tasks: inspect its added lines
  added=$(grep -E '^\+' test.patch | grep -vE '^\+\+\+')
  if echo "$added" | grep -qE 'go install|go build|pip install|npm (install|i |ci)|yarn (add|install)|pnpm (add|install)|cargo (install|add)|apt(-get)? install'; then f "install/build step added in test.patch (test.sh must not install)"; fi
  echo "$added" | grep -qE 'cargo build|mvn (install|package)|gradle(w)? build|cmake --build|make( |$)' && w "build step in test.patch: only OK if the test runner needs it and the image can't prebuild it"
  if echo "$added" | grep -qE '\|\| true'; then w "'|| true' in test.patch (fine after a junit converter if the exit code is propagated elsewhere, otherwise it swallows failures): $(echo "$added" | grep -m1 '|| true' | cut -c1-100)"; fi
  if echo "$added" | grep -qE 'go test.* -skip'; then w "'go test -skip' is often rejected as an invalid flag"; fi
  newfiles=$(grep -E '^\+\+\+ b/' test.patch | sed 's#^+++ b/##' | grep -vE '(^|/)test\.sh$')
  for n in $newfiles; do
    if ! echo "$n" | grep -qE '[0-9a-f]{6}'; then w "new test path without a random hash suffix: $n"; fi
  done
  if grep -E '^--- a/' test.patch | grep -vE '/dev/null' | grep -qE 'test'; then
    w "test.patch modifies an existing file; do not edit existing repo tests: $(grep -E '^--- a/' test.patch | head -3 | tr '\n' ' ')"
  fi
fi

if [ -f Dockerfile ]; then
  if grep -E '^[[:space:]]*(RUN|&&|;)|&&|;' Dockerfile | grep -qE '(^|[;&|[:space:]])((go|npm|yarn|pnpm|cargo|mvn|gradle|make) test|(python[0-9.]* -m )?(pytest|jest|vitest|ctest)([[:space:]]|$))'; then f "Dockerfile runs tests (hard fail)"; fi
  grep -qE '(npm install|yarn install|pnpm install)' Dockerfile && w "use lockfile installs (npm ci / yarn --frozen-lockfile / pnpm --frozen-lockfile)"
  grep -qE 'cargo (build|fetch)' Dockerfile && ! grep -qE -- '--locked' Dockerfile && w "cargo without --locked"
  grep -qE 'pip install' Dockerfile && grep -qE '(>=|~=)' Dockerfile && w "unpinned pip versions"
  if grep -qE 'pip install' Dockerfile; then
    grep -qE 'pip install.*(--no-deps).*-e|pip install.*-e.*--no-deps' Dockerfile || f "pip install without '--no-deps -e .' editable project install"
    grep -qE '(>=|~=|\*)' Dockerfile && w "unpinned versions in Dockerfile (>=, ~=, *)"
  fi
  if [ "$(grep -cE 'go[0-9.]+\.linux-amd64|openjdk|rustup|nvm install' Dockerfile)" -gt 0 ]; then
    w "extra toolchain installed: explain it in a Dockerfile comment (WARN otherwise)"
  fi
fi

# patches apply on a clean checkout of the pinned commit
repo=$(ls -d */.git 2>/dev/null | head -1 | xargs -r dirname)
if [ -n "$repo" ]; then
  tmp=$(mktemp -d); git -C "$repo" worktree add -q --detach "$tmp/w" HEAD 2>/dev/null
  if [ -d "$tmp/w" ]; then
    git -C "$tmp/w" apply --check "$PWD/test.patch" 2>/dev/null && p "test.patch applies to clean checkout" || f "test.patch does not apply"
    git -C "$tmp/w" apply "$PWD/test.patch" 2>/dev/null
    git -C "$tmp/w" apply --check "$PWD/solution.patch" 2>/dev/null && p "solution.patch applies on top of test.patch" || f "solution.patch does not apply on test.patch"
    git -C "$repo" worktree remove --force "$tmp/w" 2>/dev/null
  fi
else
  echo "INFO  no cloned repo dir here (run 'make clone' in the folder) – skipped patch apply check"
fi
exit $fail
