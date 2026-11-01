#!/usr/bin/env bash
# Usage: grade2.sh <diff path> <label>; prints v1 test result and barrier results.
W=/home/tes/Desktop/ADK/olympus-project/jedi-overload-arguments/work-v2
diff=$(realpath $1); label=$2; id=$(basename $diff .diff); wt=$W/wt-g-$id
D=test/test_inference/test_gradual
git -C $W/jedi worktree remove --force $wt >/dev/null 2>&1; rm -rf $wt
git -C $W/jedi worktree add -q --detach $wt 2af1a0f3
for s in typeshed django-stubs; do rm -rf $wt/jedi/third_party/$s; ln -s $W/jedi/jedi/third_party/$s $wt/jedi/third_party/$s; done
if ! git -C $wt apply --exclude=test.sh --exclude='.<harness>-eval/*' --exclude="$D/test_overload_arguments_9c3e7a.py" $diff 2>/dev/null; then echo "$label|APPLY_FAIL"; git -C $W/jedi worktree remove --force $wt; exit; fi
cp $W/tools/test_v1_9c3e7a.py $wt/$D/test_overload_arguments_9c3e7a.py
cp $W/tools/test_barriers_4b1c2d.py $wt/$D/
cd $wt
v1=$(PYTHONHASHSEED=0 timeout 900 $W/venv/bin/python -m pytest -p no:cacheprovider -q $D/test_overload_arguments_9c3e7a.py 2>&1 | tail -1)
bar=$(PYTHONHASHSEED=0 timeout 900 $W/venv/bin/python -m pytest -p no:cacheprovider -q -rA $D/test_barriers_4b1c2d.py 2>&1 | grep -E "^(PASSED|FAILED)" | sed 's/ - .*//; s/.*::test_\(b[0-9]*\)_.*/\1/; s/^PASSED.*:://' )
res=$(PYTHONHASHSEED=0 timeout 900 $W/venv/bin/python -m pytest -p no:cacheprovider -q -rA $D/test_barriers_4b1c2d.py 2>&1 | grep -E "^(PASSED|FAILED)" | awk '{split($2,a,"::"); split(a[2],b,"_"); print $1=="PASSED" ? "+"b[2] : "-"b[2]}' | tr '\n' ' ')
echo "$label|$v1|$res"
cd $W; git -C $W/jedi worktree remove --force $wt >/dev/null 2>&1
