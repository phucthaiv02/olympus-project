#!/bin/bash
ulimit -s unlimited
mkdir -p /tmp/g /tmp/r; cp /p/y/main.go /tmp/g/; cp /p/y/main.go /tmp/r/
if [ "$1" = native ]; then cd /tmp/g; printf 'module p\ngo 1.21\n' > go.mod; go run main.go > /out/y.go.txt 2>&1; exit; fi
cp -r /s/ref-work/compiler/. /app/compiler/; cd /app; go install .; cd /tmp/r; gopherjs run main.go > /out/y.ref.txt 2>&1
