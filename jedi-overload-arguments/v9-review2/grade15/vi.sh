T=test/test_inference/test_gradual/test_overload_arguments_9c3e7a.py
for f in /w/solution.patch /w/rd749ksxw2r7g5yscck7sh4k7n8gymkw.patch /w/rd771brznqxkkkwghjyyejrhan8gzm2x.patch; do
  rm -rf /tmp/a && cp -a /app /tmp/a && cd /tmp/a
  git -c safe.directory='*' apply --include='jedi/*' "$f" && git -c safe.directory='*' apply /w/vi_proposed.patch || { echo APPLY-FAIL; continue; }
  out=$(python -m pytest -q -p no:cacheprovider $T 2>&1)
  echo "$(basename $f .patch) | $(echo "$out" | tail -1 | sed 's/ in .*//') | $(echo "$out" | grep ^FAILED | sed 's/.*::test_//;s/ - .*//' | sort -u | tr '\n' ' ')"
done
