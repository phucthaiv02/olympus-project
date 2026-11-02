mkdir -p /w/out5
for f in /w/REF.patch /w/rd*.patch /w/vega_*.pat; do
  n=$(basename $f); n=${n%.*}; n=${n:0:12}
  rm -rf /tmp/a && cp -a /app /tmp/a && cd /tmp/a
  git -c safe.directory='*' apply --include='jedi/*' "$f" 2>/dev/null || { echo "$n APPLY-FAIL"; continue; }
  cd /w && PYTHONPATH=/tmp/a timeout 900 python probe5.py > out5/$n.txt 2>/dev/null
done
