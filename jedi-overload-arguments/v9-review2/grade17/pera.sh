T=test/test_inference/test_gradual/test_overload_arguments_9c3e7a.py
for f in /w/solution.patch /w/rd*.patch; do
  rm -rf /tmp/a && cp -a /app /tmp/a && cd /tmp/a
  git -c safe.directory='*' apply --include='jedi/*' "$f" 2>/dev/null || { echo "$(basename $f .patch) APPLY-FAIL"; continue; }
  git -c safe.directory='*' apply /w/test.patch
  out=$(cd /tmp && PYTHONPATH=/tmp/a python /w/pera.py /tmp/a/$T 2>&1 | grep '^JSON\|Error' | tail -1)
  echo "$(basename $f .patch) $out"
done
