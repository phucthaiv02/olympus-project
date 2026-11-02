#!/bin/bash
# usage: run.sh <patch|none|native> <name>   (run in gjs-ptrconv:v2; /p = probes/review8, /out writable)
ulimit -s unlimited
mkdir -p /tmp/r && cp /p/main.go /tmp/r/main.go
if [ "$1" = native ]; then cd /tmp/r && cat > go.mod <<M
module p
go 1.21
M
  go run main.go > /out/$2.txt 2>&1; exit; fi
cd /app; [ "$1" != none ] && { git init -q 2>/dev/null; git apply --unsafe-paths --directory=. $1 2>/out/$2.err || patch -p1 < $1 >/out/$2.err 2>&1; }
go install . >> /out/$2.err 2>&1
cd /tmp/r && gopherjs run main.go > /out/$2.txt 2>&1
