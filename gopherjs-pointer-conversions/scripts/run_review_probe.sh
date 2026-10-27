#!/bin/bash
# usage: run.sh probeX [nopatch]
P=$1
if [ "$2" = "nopatch" ]; then APPLY=true; else APPLY="git apply /s/p3/solution.patch"; fi
docker run --rm --network none -v ~/.cache/olympus-gjs:/s:ro gjs-ptrconv:v2 bash -c "cd /app && $APPLY && go install . >/dev/null && mkdir -p /tmp/x && cp /s/review/$P/main.go /tmp/x/ && cd /tmp/x && gopherjs run main.go" 2>&1
