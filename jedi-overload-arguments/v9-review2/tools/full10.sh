cd /app && git -c safe.directory=/app apply /rv/sol10.patch || exit 1
cd /rv && echo "== probe2"; PYTHONPATH=/app python probe2.py 2>&1 | grep -v "^OK"
echo "== probe"; PYTHONPATH=/app python probe.py 2>&1 | grep -v "^OK"
echo "== proto"; PYTHONPATH=/app python proto.py 2>&1 | grep -v "^OK"
cd /app && git -c safe.directory=/app apply /rv/test9.patch
for s in 0 1 7 42; do echo "seed $s: $(PYTHONHASHSEED=$s python -m pytest -q -p no:cacheprovider test/test_inference/test_gradual/test_overload_arguments_9c3e7a.py 2>&1 | tail -1)"; done
echo "== base"; ./test.sh base 2>&1 | tail -2
echo "== flake8"; python -m flake8 jedi/inference jedi/plugins 2>&1 | tail -5
