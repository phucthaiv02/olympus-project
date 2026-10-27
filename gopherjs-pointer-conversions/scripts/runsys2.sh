#!/bin/bash
cd /app; git apply $1 || exit 9
sed -i 's/Pkg.Path() == "syscall"/Pkg.Path() == "main"/' compiler/expressions.go
go install . 2>&1|tail -2
for n in SYS RS1 RP1; do mkdir -p /tmp/$n && cp /s/probes/probe$n/main.go /tmp/$n/ && cd /tmp/$n && gopherjs run main.go > /out/$n.gjs.txt 2>&1; go run main.go > /out/$n.go.txt 2>&1; done
