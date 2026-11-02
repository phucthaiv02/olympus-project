#!/bin/bash
# usage: run.sh <patch> <name>
ulimit -s unlimited
mkdir -p /tmp/g /tmp/r; cp /p/m/main.go /tmp/g/; cp /p/m/main.go /tmp/r/
if [ "$1" = native ]; then cd /tmp/g; printf 'module p\ngo 1.21\n' > go.mod; go run main.go > /out/xm.$2.txt 2>&1; exit; fi
cd /app; git init -q; git apply "$1"; go install .; cd /tmp/r; gopherjs run main.go > /out/xm.$2.txt 2>&1
