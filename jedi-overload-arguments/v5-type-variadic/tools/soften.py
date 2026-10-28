"""Turn the hidden test file into one that reports every failing case.

Usage: soften.py <test.patch> <out.py>
"""
import re
import sys

patch = open(sys.argv[1]).read()
part = patch[patch.index('+++ b/test/test_inference/test_gradual/test_overload_arguments_9c3e7a.py'):]
src = '\n'.join(l[1:] for l in part.split('\n')[2:] if l.startswith('+'))
src = src.replace("def check(", "def _orig_check(", 1)
helper = r'''

def soft(fn, label):
    try:
        ok = fn()
    except Exception as e:
        ok = False
        label = 'EXC %r :: %s' % (e, label)
    if not ok:
        print("FAILCASE", str(label).replace("\n", "\\n"))


def check(Script, annotations, call, expected, prefix=''):
    soft(lambda: _orig_check(Script, annotations, call, expected, prefix) or True,
         (annotations, call, prefix[:40]))
'''
src = src.replace("\nimport pytest\n", "\nimport pytest\n" + helper, 1)
out = []
for line in src.split("\n"):
    m = re.match(r"^(\s+)assert (infer_result\(.*\) == .*?)(?:, (call|\(.*\)))?$", line)
    if m:
        ind, expr, label = m.groups()
        out.append("%ssoft(lambda: %s, %s)" % (ind, expr, label or repr(expr)))
    else:
        out.append(line)
open(sys.argv[2], "w").write("\n".join(out) + "\n")
