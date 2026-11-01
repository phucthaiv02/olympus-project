#!/usr/bin/env bash
# Usage: grade.sh <run-id> <label> ; applies the agent diff to a fresh worktree and runs the new hidden tests.
W=/home/tes/Desktop/ADK/olympus-project/jedi-overload-arguments/work-v2
id=$1; label=$2; wt=$W/wt-$id
T=test/test_inference/test_gradual/test_overload_arguments_9c3e7a.py
git -C $W/jedi worktree remove --force $wt >/dev/null 2>&1; rm -rf $wt
git -C $W/jedi worktree add -q --detach $wt 2af1a0f3
for s in typeshed django-stubs; do rm -rf $wt/jedi/third_party/$s; ln -s $W/jedi/jedi/third_party/$s $wt/jedi/third_party/$s; done
if ! git -C $wt apply --exclude=test.sh --exclude='.<harness>-eval/*' --exclude=$T $W/runs/$id.diff 2>/dev/null; then echo "$label APPLY_FAIL"; exit; fi
cp $W/jedi/$T $wt/$T
out=$(cd $wt && PYTHONHASHSEED=0 timeout 900 $W/venv/bin/python -m pytest -p no:cacheprovider -q $T 2>&1)
summary=$(echo "$out" | tail -1)
fails=$(echo "$out" | grep ^FAILED | sed 's/.*:://; s/ - .*//' | tr '\n' ' ')
echo "$label | $summary | $fails"
git -C $W/jedi worktree remove --force $wt >/dev/null 2>&1
