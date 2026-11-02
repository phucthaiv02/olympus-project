#!/bin/bash
# Runs inside the container: apply test patch, run base/new without and with the solution.
set -u
cd /app
git config --global --add safe.directory /app 2>/dev/null
test ! -e /app/test.sh || exit 1
git apply /x/test.patch || exit 1
echo "== base (no solution)"; ./test.sh --output_path /tmp/base0.xml base > /tmp/base0.log 2>&1; echo "exit $?"
echo "== new (no solution)"; ./test.sh --output_path /tmp/new0.xml new > /tmp/new0.log 2>&1; echo "exit $?"
git apply /x/solution.patch || exit 1
echo "== base (solution)"; ./test.sh --output_path /tmp/base1.xml base > /tmp/base1.log 2>&1; echo "exit $?"
echo "== new (solution)"; ./test.sh --output_path /tmp/new1.xml new > /tmp/new1.log 2>&1; echo "exit $?"
for f in base0 new0 base1 new1; do
  printf "%s: tests=%s failures=%s skipped=%s\n" $f \
    "$(grep -o '<testcase ' /tmp/$f.xml | wc -l)" \
    "$(grep -c '<failure' /tmp/$f.xml)" \
    "$(grep -c '<skipped' /tmp/$f.xml)"
done
echo "new test names in base junit: $(grep -c 77753e /tmp/base1.xml)"
echo "duplicate names in base junit: $(grep -o '<testcase [^>]*' /tmp/base1.xml | sed -E 's/.*classname="([^"]*)".*name="([^"]*)".*/\1 \2/;s/.*name="([^"]*)".*classname="([^"]*)".*/\2 \1/' | sort | uniq -d | wc -l)"
grep -E "^--- FAIL" /tmp/base0.log /tmp/base1.log | head
grep -E "^(--- FAIL|--- PASS)" /tmp/new0.log | sort | uniq -c | head -30
