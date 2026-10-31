T=test/test_inference/test_gradual/test_overload_arguments_9c3e7a.py
for f in /w/solution.patch /w/rd*.patch /w/nova-*.patch; do
  rm -rf /tmp/a && cp -a /app /tmp/a && cd /tmp/a
  git -c safe.directory='*' apply --include='jedi/*' "$f" 2>/dev/null || { echo "$(basename $f .patch) APPLY-FAIL"; continue; }
  git -c safe.directory='*' apply /w/test.patch
  echo "$(basename $f .patch) $(python -m pytest -q -p no:cacheprovider $T  2>&1 | tail -1 | sed 's/ in .*//')"
done
