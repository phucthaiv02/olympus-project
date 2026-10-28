#!/usr/bin/env bash
# Usage: seeds.sh <name> <solution.patch|""> <soft_test.py> <seeds...>
# Applies a solution to a clean jedi worktree and prints the failing hidden
# cases for each PYTHONHASHSEED.
S=/tmp/claude-1000/-home-tes-Desktop-ADK-olympus-project/234e53cc-0fe8-411f-8112-e2c6f4c4a3db/scratchpad
name=$1; patch=$2; test=$3; shift 3
wt=$S/wt-$name
git -C $S/jedi-base worktree remove --force $wt >/dev/null 2>&1; rm -rf $wt
git -C $S/jedi-base worktree add -q --detach $wt HEAD
for sub in typeshed django-stubs; do rm -rf $wt/jedi/third_party/$sub; ln -s $S/jedi-base/jedi/third_party/$sub $wt/jedi/third_party/$sub; done
if [ -n "$patch" ]; then
  git -C $wt apply --exclude=test.sh "$patch" 2>/dev/null || { echo "[$name] APPLY FAILED"; exit 1; }
fi
cp $test $wt/test/test_inference/test_gradual/test_overload_arguments_9c3e7a.py
cd $wt
for seed in "$@"; do
  res=$(PYTHONHASHSEED=$seed python3 -m pytest -p no:cacheprovider -q -s test/test_inference/test_gradual/test_overload_arguments_9c3e7a.py 2>&1)
  fails=$(echo "$res" | grep FAILCASE | sed 's/^.*FAILCASE //' | tr '\n' '|')
  summary=$(echo "$res" | grep -E '[0-9]+ (passed|failed)' | tail -1)
  echo "[$name seed=$seed] ${fails:-ALL PASS}  ($summary)"
done
