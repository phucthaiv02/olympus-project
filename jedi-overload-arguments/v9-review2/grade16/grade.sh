T=test/test_inference/test_gradual/test_overload_arguments_9c3e7a.py
for f in /w/rd*.patch; do
  rm -rf /tmp/a && cp -a /app /tmp/a && cd /tmp/a
  git -c safe.directory='*' apply --include='jedi/*' "$f" 2>/dev/null || { echo "$(basename $f .patch) APPLY-FAIL"; continue; }
  git -c safe.directory='*' apply /w/test.patch
  out=$(python -m pytest -q -p no:cacheprovider $T 2>&1)
  echo "$(basename $f .patch) | $(echo "$out" | tail -1 | sed 's/ in .*//') | $(echo "$out" | grep ^FAILED | sed 's/.*::test_//;s/ - .*//' | sort -u | tr '\n' ' ')"
done
