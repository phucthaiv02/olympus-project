#!/usr/bin/env bash
# Checks 1-3: tests fail without solution, solution passes all, no flakiness. Runs in an isolated container.
# usage: verify.sh <submission-dir> [--repeat N] [--new-pkg PATTERN]...
# Language-agnostic: relies only on test.sh <base|new> emitting junit. Needs: docker, a cloned repo dir (see Makefile 'clone'), test.patch containing test.sh, solution.patch.
set -uo pipefail
D="${1:?usage: verify.sh <dir> [--repeat N] [--new-pkg PATTERN]...}"; shift
REPEAT=3; PKGS=()
while [ $# -gt 0 ]; do case "$1" in
  --repeat) REPEAT="$2"; shift 2;;
  --new-pkg) PKGS+=(--new-pkg "$2"); shift 2;;
  *) echo "unknown arg $1"; exit 2;; esac; done
HERE="$(cd "$(dirname "$0")" && pwd)"
cd "$D" || exit 2
D=$PWD
repo=$(ls -d */.git 2>/dev/null | head -1 | xargs -r dirname)
[ -n "$repo" ] || { echo "no cloned repo here; run: make clone"; exit 2; }
OUT="$D/checks/verify-$(date +%F-%H%M)"; mkdir -p "$OUT"
TAG="qc-verify-$(basename "$D" | tr 'A-Z' 'a-z' | tr -c 'a-z0-9\n-' '-')"

# standalone clone at the pinned commit (no patch leftovers; /app needs a real .git)
W=$(mktemp -d); git clone -q -c advice.detachedHead=false --no-hardlinks "$D/$repo" "$W/src" || exit 2
trap 'rm -rf "$W"' EXIT
git -C "$W/src" submodule update --init --recursive -q 2>/dev/null
rm -f "$W/src/.dockerignore"
docker build -q --platform linux/amd64 -f "$D/Dockerfile" -t "$TAG" "$W/src" >"$OUT/build.log" 2>&1 \
  || { echo "FAIL  docker build (see $OUT/build.log)"; tail -20 "$OUT/build.log"; exit 1; }
echo "PASS  docker build"

run_mode() { # label withSolution mode run#
  local label=$1 sol=$2 mode=$3 n=$4 xml="$OUT/$1-$3-$4.xml"
  cat >"$OUT/inner-$label-$mode-$n.sh" <<INNER
set -u; cd /app && git config --global --add safe.directory /app
git apply /patches/test.patch || exit 97
[ "$sol" = 1 ] && { git apply /patches/solution.patch || exit 98; }
ulimit -s unlimited
chmod +x test.sh; ./test.sh --output_path /out/$(basename "$xml") $mode
INNER
  docker run --rm --network none --user root -v "$D:/patches:ro" -v "$OUT:/out" "$TAG" bash -c "
    id model >/dev/null 2>&1 || useradd -m model; chown -R model /app /out 2>/dev/null
    su model -s /bin/bash -c 'bash /out/inner-$label-$mode-$n.sh'" >"$OUT/$label-$mode-$n.log" 2>&1
  echo "      $label/$mode/run$n exit=$?"
}
for sol in 0 1; do
  lab=$([ $sol = 1 ] && echo solution || echo tests)
  for mode in base new; do run_mode $lab $sol $mode 1; done
done
echo "--- verify contract"
python3 -I "$HERE/junit_check.py" --base-tests-only "$OUT/tests-base-1.xml" --new-tests-only "$OUT/tests-new-1.xml" \
  --base-solution "$OUT/solution-base-1.xml" --new-solution "$OUT/solution-new-1.xml" "${PKGS[@]}"
rc=$?

echo "--- flakiness ($REPEAT repeats)"
flaky=0
for sol in 0 1; do
  lab=$([ $sol = 1 ] && echo solution || echo tests)
  for mode in base new; do
    for i in $(seq 2 "$REPEAT"); do run_mode $lab $sol $mode $i >/dev/null; done
    python3 -I - "$OUT" "$lab" "$mode" "$REPEAT" <<'PY' || flaky=1
import sys, xml.etree.ElementTree as ET
out, lab, mode, n = sys.argv[1], sys.argv[2], sys.argv[3], int(sys.argv[4])
res = {}
for i in range(1, n + 1):
    try:
        root = ET.parse(f"{out}/{lab}-{mode}-{i}.xml").getroot()
    except Exception as e:
        print(f"FAIL  {lab}/{mode}: run {i} produced no junit ({e})"); sys.exit(1)
    for tc in root.iter("testcase"):
        k = tc.get("classname", "") + "::" + tc.get("name", "")
        st = "fail" if tc.find("failure") is not None or tc.find("error") is not None else ("skip" if tc.find("skipped") is not None else "pass")
        res.setdefault(k, []).append(st)
flaky = {k: v for k, v in res.items() if len(set(v)) > 1}
print(("FAIL  " if flaky else "PASS  ") + f"{lab}/{mode}: {len(flaky)} flaky of {len(res)} tests")
for k, v in list(flaky.items())[:10]: print("      ", k, v)
sys.exit(1 if flaky else 0)
PY
  done
done
echo "Artifacts: $OUT"
[ $rc -eq 0 ] && [ $flaky -eq 0 ]
