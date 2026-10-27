#!/bin/bash
# usage: run.sh probe.go [extra gopherjs flags]
f=$1; shift
docker run --rm --network none -v /home/tes/.cache/olympus-gjs:/s:ro -v /tmp/claude-1000/review6:/r:ro gjs-ptrconv:v2 bash -c "cd /app && git apply /s/p8/solution.patch && go install . >/dev/null 2>&1 && mkdir -p /tmp/x && cp /r/$f /tmp/x/main.go && cd /tmp/x && go mod init x >/dev/null 2>&1; echo '== native'; go run main.go 2>&1; echo '== gopherjs'; ulimit -s unlimited; gopherjs run $* main.go 2>&1 | head -50"
