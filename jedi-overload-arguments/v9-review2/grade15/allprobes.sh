rm -rf /tmp/r && cp -a /app /tmp/r && cd /tmp/r && git -c safe.directory='*' apply /w/solution.patch
cd /w
for p in probe.py probe2.py probe3.py probe4.py probe5.py probe6.py probe7.py probe8.py probe9.py probe10.py probe11.py proto.py; do
  out=$(PYTHONPATH=/tmp/r timeout 900 python $p 2>&1)
  echo "== $p: lines=$(echo "$out" | wc -l) BAD=$(echo "$out" | grep -c '^BAD') CRASH=$(echo "$out" | grep -c 'CRASH\|Traceback') $(echo "$out" | grep -E '^bad [0-9]+' | tail -1)"
  echo "$out" | grep -E '^BAD|CRASH|Traceback' | head -8
done
