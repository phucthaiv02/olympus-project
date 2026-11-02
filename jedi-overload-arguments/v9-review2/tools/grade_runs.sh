for f in /w/rd*.patch; do
  rm -rf /tmp/a && cp -a /app /tmp/a && cd /tmp/a
  git -c safe.directory='*' apply --include='jedi/*' "$f" 2>/dev/null || { echo "$(basename $f .patch) APPLY-FAIL"; continue; }
  git -c safe.directory='*' apply /w/test9.patch
  out=$(python -m pytest -q -p no:cacheprovider test/test_inference/test_gradual/test_overload_arguments_9c3e7a.py 2>&1)
  keys=$(cd /w && PYTHONPATH=/tmp/a python probe2.py 2>/dev/null | grep -E "class key |tuple key |1 True|cond equal|class key in" | awk '{printf "%s:%s ", $1, $3$4}')
  echo "$(basename $f .patch) | $(echo "$out" | tail -1 | sed 's/ in .*//') | $(echo "$out" | grep ^FAILED | sed 's/.*::test_//' | tr '\n' ' ') | $keys"
done
