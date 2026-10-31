T=test/test_inference/test_gradual/test_overload_arguments_9c3e7a.py
rm -rf /tmp/r && cp -a /app /tmp/r && cd /tmp/r && git -c safe.directory='*' apply /w/solution.patch /w/test.patch
for s in 0 1 2 3 7 42; do echo "ref seed $s: $(PYTHONHASHSEED=$s python -m pytest -q -p no:cacheprovider $T 2>&1 | tail -1)"; done
cd /w
for p in probe.py probe2.py probe3.py probe4.py probe5.py probe6.py probe7.py probe8.py probe9.py probe10.py probe11.py proto.py matrix.py matrix2.py; do
  out=$(PYTHONPATH=/tmp/r timeout 400 python $p 2>&1)
  echo "$out" > /o/ref-$p.txt
  echo "== $p: lines=$(echo "$out" | wc -l) BAD=$(echo "$out" | grep -c '^BAD') CRASH=$(echo "$out" | grep -c 'CRASH\|Traceback')"
  echo "$out" | grep -E '^BAD|Traceback' | cut -c1-120 | head -8
done
for f in /w/rd*.patch /w/nova-*.patch; do
  rm -rf /tmp/a && cp -a /app /tmp/a && cd /tmp/a
  git -c safe.directory='*' apply --include='jedi/*' "$f" 2>/dev/null || { echo "$(basename $f .patch) APPLY-FAIL"; continue; }
  git -c safe.directory='*' apply /w/test.patch
  out=$(python -m pytest -q -p no:cacheprovider $T 2>&1)
  echo "$(basename $f .patch) | $(echo "$out" | tail -1 | sed 's/ in .*//') | $(echo "$out" | grep ^FAILED | sed 's/.*::test_//;s/ - .*//;s/\[.*//' | sort -u | tr '\n' ' ')"
done
