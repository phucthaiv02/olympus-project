cd /app && git config --global --add safe.directory /app
if git apply /out/run.patch >/out/apply.log 2>&1; then echo 'applied cleanly' >>/out/apply.log; else echo 'clean apply failed, retrying with --reject (rejected hunks are NOT applied)' >>/out/apply.log; git apply --reject /out/run.patch >>/out/apply.log 2>&1; fi
git apply /in/test.patch || exit 97
ulimit -s unlimited; chmod +x test.sh; ./test.sh --output_path /out/base.xml base
