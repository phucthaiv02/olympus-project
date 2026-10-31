T=/w
rm -rf /tmp/r && cp -a /app /tmp/r && cd /tmp/r && git -c safe.directory='*' apply /p/solution.patch
cd /tmp && PYTHONPATH=/tmp/r python /w/matrix.py 2>&1 | grep '^JSON\|Traceback\|Error' | tail -1 > /o/reference.txt
for id in rd749ksxw2r7g5yscck7sh4k7n8gymkw rd74wc0e3hx2ss0wjp0t6vc75s8h97y2; do
  rm -rf /tmp/a && cp -a /app /tmp/a && cd /tmp/a && git -c safe.directory='*' apply --include='jedi/*' /p/runs/$id/$id.patch
  cd /tmp && PYTHONPATH=/tmp/a python /w/matrix.py 2>&1 | grep '^JSON\|Traceback\|Error' | tail -1 > /o/$id.txt
done
rm -rf /tmp/b && cp -a /app /tmp/b; cd /tmp && PYTHONPATH=/tmp/b python /w/matrix.py 2>&1 | grep '^JSON\|Traceback\|Error' | tail -1 > /o/base.txt
