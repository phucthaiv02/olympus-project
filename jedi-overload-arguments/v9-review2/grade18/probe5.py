import sys
sys.path.insert(0, '/w')
from probe import ov, infer
D = ['x: dict[str, type[A]]', 'x: dict[str, type[C]]', 'x: object']
L = ['x: list[type[A]]', 'x: list[type[C]]', 'x: object']
cases = [
    ('dict over str', '', D, 'f({c: C for c in "ab"})'),
    ('list over str', '', L, 'f([C for _ in "ab"])'),
    ('dict over empty str', '', D, 'f({"k": C for _ in ""})'),
    ('list over empty str', '', L, 'f([C for _ in ""])'),
    ('dstar empty str comp', '', D, 'f({"k": A, **{"k": C for _ in ""}})'),
    ('dict over empty bytes', '', D, 'f({"k": C for _ in b""})'),
    ('dict over str var', 's = "ab"\n', D, 'f({c: C for c in s})'),
    ('dict over empty str var', 's = ""\n', D, 'f({"k": C for _ in s})'),
    ('dict over empty tuple', '', D, 'f({"k": C for _ in ()})'),
    ('dict over empty set call', '', D, 'f({"k": C for _ in set()})'),
    ('dict over empty dict', '', D, 'f({"k": C for _ in {}})'),
    ('list over empty dict', '', L, 'f([C for _ in {}])'),
    ('dict dep for char', '', D, 'f({k: C for k in ["a"] for _ in k})'),
    ('dict comp in dstar', '', D, 'f({"p": A, **{k: C for k in ["p"]}})'),
    ('dict value loop var mixed', '', D, 'f({"k": v for v in [A, C]})'),
]
for name, pre, anns, call in cases:
    print('%s\t%s' % (name, infer(pre + ov(*anns) + call)))
