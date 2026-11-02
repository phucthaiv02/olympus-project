set -e
cd /app
git -c safe.directory=/app apply /rv/solution.patch
echo "== base? $1"
python /rv/probe.py
