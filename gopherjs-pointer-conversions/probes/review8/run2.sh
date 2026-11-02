#!/bin/bash
ulimit -s unlimited
mkdir -p /tmp/r && cp /p/main.go /tmp/r/main.go
cd /app; git init -q 2>/dev/null; git apply $1 2>/out/$2.err || patch -p1 --forward < $1 >>/out/$2.err 2>&1
go install . >> /out/$2.err 2>&1
cd /tmp/r && timeout 200 gopherjs run main.go > /out/$2.txt 2>&1
mkdir -p /app/zzjsn && cp /s/probes/probeJSN/main.go /app/zzjsn/ && cd /app/zzjsn && timeout 200 gopherjs run main.go > /out/$2.jsn.txt 2>&1
