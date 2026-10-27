#!/usr/bin/env bash
# Grade agent runs on the current tests. usage: grade_run.sh <submission-dir> <run-dir>... [--modes "base new"] [--drop-regex RE]
# --modes new skips the (slow) existing-suite run; base shows regressions in existing tests.
# A run dir holds the run's workspace diff (first *.patch or *.diff found). Needs docker, a cloned repo dir in the
# submission (see its Makefile 'clone'), Dockerfile, test.patch (adds test.sh <base|new>).
set -uo pipefail
D="${1:?usage: grade_run.sh <submission-dir> <run-dir>... [--drop-regex RE]}"; shift
RUNS=(); EXTRA=(); MODES="base new"
while [ $# -gt 0 ]; do case "$1" in --modes) MODES="$2"; shift 2;; --drop-regex) EXTRA+=(--drop-regex "$2"); shift 2;; *) RUNS+=("$(cd "$1" && pwd)"); shift;; esac; done
D="$(cd "$D" && pwd)"; HERE="$(cd "$(dirname "$0")" && pwd)"
repo=$(cd "$D" && ls -d */.git 2>/dev/null | head -1 | xargs -r dirname)
[ -n "$repo" ] || { echo "no cloned repo in $D"; exit 2; }
OUT="$D/fp-grade-$(date +%F-%H%M)"; mkdir -p "$OUT"
TAG="fp-grade-$(basename "$D" | tr 'A-Z' 'a-z' | tr -c 'a-z0-9\n-' '-')"
W=$(mktemp -d); git clone -q -c advice.detachedHead=false --no-hardlinks "$D/$repo" "$W/src" || exit 2   # standalone clone: /app needs a real .git
trap 'rm -rf "$W"' EXIT
git -C "$W/src" submodule update --init --recursive -q 2>/dev/null; rm -f "$W/src/.dockerignore"
docker build -q --platform linux/amd64 -f "$D/Dockerfile" -t "$TAG" "$W/src" >"$OUT/build.log" 2>&1 \
  || { echo "FAIL docker build, see $OUT/build.log"; tail -15 "$OUT/build.log"; exit 1; }

summ() { python3 -I - "$1" <<'PY'
import sys, xml.etree.ElementTree as ET
try: root = ET.parse(sys.argv[1]).getroot()
except Exception as e: print("    no junit:", e); sys.exit(0)
tot = bad = 0; names = []
for tc in root.iter("testcase"):
    tot += 1
    if tc.find("failure") is not None or tc.find("error") is not None or tc.find("skipped") is not None:
        bad += 1; names.append(tc.get("classname", "") + "::" + tc.get("name", ""))
print(f"    {tot - bad}/{tot} pass")
for n in names[:15]: print("      not passing:", n)
if len(names) > 15: print(f"      ... {len(names) - 15} more")
PY
}
for R in "${RUNS[@]}"; do
  name=$(basename "$R"); P=$(ls "$R"/*.patch "$R"/*.diff 2>/dev/null | head -1)
  [ -n "$P" ] || { echo "== $name: no patch found"; continue; }
  mkdir -p "$OUT/$name"
  python3 -I "$HERE/strip_diff.py" "$P" "$OUT/$name/run.patch" --drop-from "$D/test.patch" "${EXTRA[@]}"
  echo "== $name"
  for mode in $MODES; do
    cat >"$OUT/$name/inner-$mode.sh" <<INNER
cd /app && git config --global --add safe.directory /app
if git apply /out/run.patch >/out/apply.log 2>&1; then echo 'applied cleanly' >>/out/apply.log; else echo 'clean apply failed, retrying with --reject (rejected hunks are NOT applied)' >>/out/apply.log; git apply --reject /out/run.patch >>/out/apply.log 2>&1; fi
git apply /in/test.patch || exit 97
ulimit -s unlimited; chmod +x test.sh; ./test.sh --output_path /out/$mode.xml $mode
INNER
    docker run --rm --network none --user root -v "$D/test.patch:/in/test.patch:ro" -v "$OUT/$name:/out" "$TAG" bash -c "
      id model >/dev/null 2>&1 || useradd -m model; chown -R model /app /out 2>/dev/null
      su model -s /bin/bash -c 'bash /out/inner-$mode.sh'" >"$OUT/$name/$mode.log" 2>&1
    echo "  $mode (exit $?)"; summ "$OUT/$name/$mode.xml"
  done
done
echo "results in $OUT"
