#!/bin/bash
# usage: runsys.sh <patch|none> <outprefix>: treat package main like syscall
cd /app; if [ "$1" != none ]; then git apply $1 || exit 9; fi
sed -i 's/Pkg.Path() == "syscall"/Pkg.Path() == "main"/' compiler/expressions.go
go install . 2>&1|tail -2
mkdir -p /tmp/sys && cp /s/probes/probeSYS/main.go /tmp/sys/ && cd /tmp/sys && gopherjs run main.go > /out/$2.sys.txt 2>&1
cd /tmp/sys && go run main.go > /out/$2.gosys.txt 2>&1
gopherjs build -o /out/$2.sys.js main.go 2>/dev/null
