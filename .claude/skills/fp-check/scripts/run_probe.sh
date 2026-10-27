#!/usr/bin/env bash
# Run a probe file against one version of the code and store its output.
# usage: run_probe.sh <submission-dir> <patch|-> <probe-file> <out.jsonl> [command...]
#   <patch>   a run's diff (test files stripped) or the reference solution.patch; "-" = clean checkout
#   <probe-file> is mounted at /probes/<name>; the command (default: python) is run as: <command> /probes/<name>
#   The probe must print one JSON object per line: {"probe": "<id>", "result": "<short string>"}
set -uo pipefail
D="$(cd "${1:?usage: run_probe.sh <submission-dir> <patch|-> <probe-file> <out.jsonl> [command...]}" && pwd)"; PATCH="$2"; PROBE="$(cd "$(dirname "$3")" && pwd)/$(basename "$3")"; OUTF="$(cd "$(dirname "$4")" && pwd)/$(basename "$4")"; shift 4
CMD=("${@:-python}")
HERE="$(cd "$(dirname "$0")" && pwd)"
TAG="fp-grade-$(basename "$D" | tr 'A-Z' 'a-z' | tr -c 'a-z0-9\n-' '-')"
docker image inspect "$TAG" >/dev/null 2>&1 || { echo "image $TAG missing: run grade_run.sh once first (it builds it)"; exit 2; }
T=$(mktemp -d); trap 'rm -rf "$T"' EXIT
if [ "$PATCH" != "-" ]; then python3 -I "$HERE/strip_diff.py" "$PATCH" "$T/code.patch" --drop-from "$D/test.patch" >/dev/null; else : > "$T/code.patch"; fi
cat >"$T/inner.sh" <<INNER
cd /app && git config --global --add safe.directory /app
[ -s /t/code.patch ] && { git apply --3way /t/code.patch >/dev/null 2>&1 || git apply --reject /t/code.patch >/dev/null 2>&1; }
ulimit -s unlimited
${CMD[*]} /probes/$(basename "$PROBE")
INNER
docker run --rm --network none --user root -v "$T:/t" -v "$PROBE:/probes/$(basename "$PROBE"):ro" "$TAG" bash -c "
  id model >/dev/null 2>&1 || useradd -m model; chown -R model /app 2>/dev/null
  su model -s /bin/bash -c 'bash /t/inner.sh'" 2>"$OUTF.err" | grep '^{' > "$OUTF"
echo "$(wc -l < "$OUTF") probe results -> $OUTF"
