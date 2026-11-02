import sys
sys.path.insert(0, '/w')
from probe import ov, infer
R = 'c = random.random()\n'
D = ['x: dict[str, type[A]]', 'x: dict[str, type[C]]', 'x: object']
L = ['x: list[type[A]]', 'x: list[type[C]]', 'x: object']
KW = ['**kw: type[A]', '**kw: type[C]', '**kw: object']
cases = [
    ('dict nested for', '', D, 'f({k: C for ks in [["a", "b"]] for k in ks})'),
    ('dict dep for char', '', D, 'f({k: C for k in ["a"] for _ in k})'),
    ('list nested for', '', L, 'f([x for xs in [[C]] for x in xs])'),
    ('list nested for mixed', '', L, 'f([x for xs in [[C], [A]] for x in xs])'),
    ('dict tuple target', '', D, 'f({k: v for k, v in [("a", C)]})'),
    ('dict tuple target mixed', '', D, 'f({k: v for k, v in [("a", C), ("b", A)]})'),
    ('dict items', '', D, 'f({k: v for k, v in {"a": C}.items()})'),
    ('dict over var', 'ks = ["a", "b"]\n', D, 'f({k: C for k in ks})'),
    ('dict value loop var', '', D, 'f({"k": v for v in [C]})'),
    ('dict value loop var mixed', '', D, 'f({"k": v for v in [A, C]})'),
    ('list loop var', '', L, 'f([x for x in [A, C]])'),
    ('list loop var ok', '', L, 'f([x for x in [C, C]])'),
    ('list loop star', '', L, 'f([x for x in [*[C], C]])'),
    ('dict comp in dstar', '', D, 'f({"p": A, **{k: C for k in ["p"]}})'),
    ('dict comp in dstar other', '', D, 'f({"p": A, **{k: C for k in ["q"]}})'),
    ('kwargs from comp', '', KW, 'f(**{k: C for k in ["p", "q"]})'),
    ('kwargs from comp nested', '', KW, 'f(**{k: C for ks in [["p"]] for k in ks})'),
    ('star from comp', '', ['*args: type[A]', '*args: type[C]', '*args: object'], 'f(*[x for xs in [[C]] for x in xs])'),
    ('comp cond value', R, D, 'f({"k": (A if c else C) for _ in [1]})'),
    ('comp over cond list', R + 'xs = [C] if c else [A]\n', L, 'f([x for x in xs])'),
    ('comp over maybe empty', R + 'xs = [C] if c else []\n', L, 'f([x for x in xs])'),
    ('dict comp over maybe empty', R + 'xs = ["k"] if c else []\n', D, 'f({"k": A, **{k: C for k in xs}})'),
    ('set comp nested', '', ['x: set[type[A]]', 'x: set[type[C]]', 'x: object'], 'f({x for xs in [[C]] for x in xs})'),
    ('gen to list call', '', L, 'f(list(x for x in [C]))'),
    ('comp tuple elements', '', ['x: tuple[type[C], ...]', 'x: object'], 'f(tuple(x for x in [C]))'),
]
for name, pre, anns, call in cases:
    print('%s\t%s' % (name, infer(pre + ov(*anns) + call)))
