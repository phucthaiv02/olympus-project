import sys
sys.path.insert(0, '/w')
from probe import ov, infer
R = 'c = random.random()\n'
D = ['x: dict[str, type[A]]', 'x: dict[str, type[C]]', 'x: object']
L = ['x: list[type[A]]', 'x: list[type[C]]', 'x: object']
KW = ['**kw: type[A]', '**kw: type[C]', '**kw: object']
cases = [
    ('tuple target replace', '', D, 'f({k: v for k, v in [("a", A), ("a", C)]})', ['RB']),
    ('tuple target two keys', '', D, 'f({k: v for k, v in [("a", A), ("b", C)]})', ['RC']),
    ('nested for replace', '', D, 'f({"k": v for xs in [[A], [C]] for v in xs})', ['RB']),
    ('if on loop var', '', D, 'f({k: C for k in ["a", ""] if k})', ['RB']),
    ('if false list', '', L, 'f([C for _ in [1] if 0])', ['RA']),
    ('if unknown', R, D, 'f({"k": C for _ in [1] if c})', ['RB']),
    ('if unknown replace', R, D, 'f({"k": A, **{"k": C for _ in [1] if c}})', ['RC']),
    ('kwargs replace in comp', '', KW, 'f(**{"p": v for v in [A, C]})', ['RB']),
    ('kwargs dstar comp key', '', KW, 'f(**{"p": A, **{k: C for k in ["p"]}})', ['RB']),
    ('kwargs comp two keys', '', KW, 'f(**{k: C for k in ["p", "q"]})', ['RB']),
    ('comp over var list', 'ks = ["p"]\n', D, 'f({"p": A, **{k: C for k in ks}})', ['RB']),
    ('comp over cond list', R + 'ks = ["p"] if c else []\n', D, 'f({"p": A, **{k: C for k in ks}})', ['RC']),
    ('comp over unknown', 'def g() -> list[str]: ...\n', D, 'f({"p": A, **{k: C for k in g()}})', ['RC']),
    ('comp over str stays', '', D, 'f({k: C for k in "ab"})', ['RB']),
    ('comp over range', '', D, 'f({str(i): C for i in range(2)})', ['RB']),
    ('comp key class', '', ['x: dict[type, type[A]]', 'x: dict[type, type[C]]', 'x: object'], 'f({k: v for k, v in [(C, A), (C, C)]})', ['RB']),
]
bad = 0
for name, pre, anns, call, exp in cases:
    got = infer(pre + ov(*anns) + call)
    bad += got != exp
    print('%-4s %-26s got=%s exp=%s' % ('OK' if got == exp else 'BAD', name, got, exp))
print('bad', bad)
