run() { # name patch
  rm -rf /tmp/a && cp -a /app /tmp/a && cd /tmp/a
  if [ -n "$2" ] && ! git -c safe.directory='*' apply "$2" 2>/dev/null; then echo "### $1 APPLY-FAIL"; return; fi
  git -c safe.directory='*' apply /rv/test9.patch
  out=$(python -m pytest -q -p no:cacheprovider test/test_inference/test_gradual/test_overload_arguments_9c3e7a.py 2>&1)
  echo "### $1 [$(echo "$out" | tail -1)] $(echo "$out" | grep ^FAILED | sed 's/.*:://' | tr '\n' ' ')"
}
run BASE ""
run REF /rv/sol9.patch
for s in 0 1 7 42; do PYTHONHASHSEED=$s run REF-seed$s /rv/sol9.patch; done
for f in /rv/p/*.patch; do run $(basename $f .patch) $f; done
