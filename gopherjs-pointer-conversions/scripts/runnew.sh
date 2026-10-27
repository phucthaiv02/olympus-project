#!/bin/bash
# usage: runnew.sh <solution patch or none|native> <name>
cd /app; git apply /p3/test.patch || exit 9
if [ "$1" = native ]; then go test -count=1 -v ./tests/ptrconv_658733 2>&1 | grep -E '^(--- |ok|FAIL|panic)' > /out/$2.txt; exit; fi
[ "$1" != none ] && { git apply $1 || exit 8; }
go install . >/dev/null 2>&1
gopherjs test -v ./tests/ptrconv_658733 2>&1 | grep -E '^(--- |ok|FAIL|panic)' > /out/$2.txt
