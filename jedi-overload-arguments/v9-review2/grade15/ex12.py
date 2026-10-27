import sys
sys.path.insert(0, '/w')
from probe import ov, infer
R = 'c = random.random()\n'
for name, pre, anns, call, exp in [
    ('cond arg', R, ['x: type[A]', 'x: object'], 'f(A if c else C)', ['RB']),
    ('list star', 'def g() -> list[type[C]]: ...\nxs = g()\n', ['x: type[C]', '*args: object'], 'f(*xs)', ['RB']),
    ('list star ann', 'xs: list[type[C]] = []\n', ['x: type[C]', '*args: object'], 'f(*xs)', ['RB']),
]:
    got = infer(pre + ov(*anns) + call)
    print('OK ' if got == exp else 'BAD', name, got, exp)
