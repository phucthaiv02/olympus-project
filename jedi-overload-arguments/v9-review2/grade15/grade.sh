T=test/test_inference/test_gradual/test_overload_arguments_9c3e7a.py
rm -rf /tmp/b && cp -a /app /tmp/b && cd /tmp/b && git -c safe.directory='*' apply /w/test.patch
echo "base: $(python -m pytest -q -p no:cacheprovider $T 2>&1 | tail -1)"
rm -rf /tmp/r && cp -a /app /tmp/r && cd /tmp/r && git -c safe.directory='*' apply /w/solution.patch /w/test.patch
for s in 0 1 2 3 7 42; do echo "ref seed $s: $(PYTHONHASHSEED=$s python -m pytest -q -p no:cacheprovider $T 2>&1 | tail -1)"; done
echo "== reviewer probe (8 issues) on ref"; cd /w && PYTHONPATH=/tmp/r python probe.py 2>&1 | tail -20
for f in /w/rd*.patch /w/nova-*.patch; do
  rm -rf /tmp/a && cp -a /app /tmp/a && cd /tmp/a
  git -c safe.directory='*' apply --include='jedi/*' "$f" 2>/dev/null || { echo "$(basename $f .patch) APPLY-FAIL"; continue; }
  git -c safe.directory='*' apply /w/test.patch
  out=$(python -m pytest -q -p no:cacheprovider $T 2>&1)
  echo "$(basename $f .patch) | $(echo "$out" | tail -1 | sed 's/ in .*//') | $(echo "$out" | grep ^FAILED | sed 's/.*::test_//;s/ - .*//;s/\[.*//' | sort -u | tr '\n' ' ')"
done
