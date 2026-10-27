#!/bin/bash
# usage: runprobes.sh <patchfile-in-container> <outprefix> <probe numbers...>
cd /app; git apply $1 || exit 9; go install . 2>&1|tail -2
out=$2; shift 2
for n in "$@"; do mkdir -p /tmp/pr$n && cp /s/probes/probe$n/main.go /tmp/pr$n/ && (cd /tmp/pr$n && gopherjs run main.go > /out/$out.p$n.txt 2>&1); done
