for who in base ref; do
  rm -rf /tmp/a && cp -a /app /tmp/a && cd /tmp/a
  [ $who = ref ] && git -c safe.directory='*' apply /w/solution.patch
  echo "===== $who"; cd /w && PYTHONPATH=/tmp/a timeout 900 python probe11.py 2>&1 | tail -n +16
done
