#!/usr/bin/env bash
# Usage: ./test.sh [--output_path <junit.xml>] <base|new>
#
#   base  the existing tests.
#   new   the index operator method tests.
set -uo pipefail

cd /app

OUTPUT_PATH=""
if [ "${1:-}" = "--output_path" ]; then
  OUTPUT_PATH="$2"
  shift 2
fi

MODE="${1:-new}"

TEST_LOG="$(mktemp)"

case "$MODE" in
  base)
    # TestTestxgo/domaintext-md imports a module that go.mod doesn't require,
    # so it needs the network; it is left out.
    go test -count=1 -v -skip '_77753e$|^TestTestxgo$' ./... 2>&1 | tee "$TEST_LOG"
    STATUS=${PIPESTATUS[0]}
    go test -count=1 -v -run '^TestTestxgo$' -skip '^TestTestxgo$/^domaintext-md$' ./cl/ 2>&1 | tee -a "$TEST_LOG"
    [ "${PIPESTATUS[0]}" -eq 0 ] || STATUS=1
    ;;
  new)
    go test -count=1 -v -run '_77753e$' ./cl/ 2>&1 | tee "$TEST_LOG"
    STATUS=${PIPESTATUS[0]}
    ;;
  *)
    echo "unknown mode: $MODE (expected base or new)" >&2
    exit 2
    ;;
esac

if [ -n "$OUTPUT_PATH" ]; then
  go-junit-report -set-exit-code < "$TEST_LOG" > "$OUTPUT_PATH" || true
fi

exit "$STATUS"
