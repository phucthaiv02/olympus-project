#!/bin/bash
# usage: run.sh probe.go [jsonly]
f=$1
docker run --rm --network none -v /home/tes/.cache/olympus-gjs:/s:ro -v /tmp/claude-1000/review7:/r gjs-ptrconv:v2 bash -c "
cd /app && go install . >/dev/null 2>&1 && cp -r /app /tmp/base && mkdir -p /app/x && cp /r/$f /app/x/main.go && cd /app/x && ulimit -s unlimited && echo '==base' && gopherjs run main.go 2>&1 | tail -150
cd /app && git apply /s/p13/solution.patch && go install . >/dev/null 2>&1 && cd /app/x && echo '==patched' && gopherjs run main.go 2>&1 | tail -150
if [ -z \"$2\" ]; then echo ==native; mkdir -p /tmp/n && cp /r/$f /tmp/n/main.go && cd /tmp/n && go mod init x >/dev/null 2>&1 && go run . 2>&1 | tail -150; fi
"
