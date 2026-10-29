import sys
sys.path.insert(0, '.')
from cases import PRELUDE, CASES
out = [PRELUDE]
for i, (g, n, code) in enumerate(CASES):
    assert ']==]' not in code
    out.append(f'print("@@CASE {i}")\nprint(run(assert(load([==[{code}]==], "=case"))))')
open('all.lua', 'w').write("\n".join(out) + "\n")
