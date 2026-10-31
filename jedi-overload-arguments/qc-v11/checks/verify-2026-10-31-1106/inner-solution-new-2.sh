set -u; cd /app && git config --global --add safe.directory /app
git apply /patches/test.patch || exit 97
[ "1" = 1 ] && { git apply /patches/solution.patch || exit 98; }
ulimit -s unlimited
chmod +x test.sh; ./test.sh --output_path /out/solution-new-2.xml new
