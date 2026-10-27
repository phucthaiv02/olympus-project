"""The 8 Solution & Code issues of the human review (v8) with variants."""
import sys
sys.path.insert(0, '/w')
from probe import ov, infer

R = 'c = random.random()\n'
D3 = ['x: dict[str, type[A]]', 'x: dict[str, type[C]]', 'x: object']
WRAP = 'def wrap(d):\n    return {**d}\n'


def proto(pm, im, pre=''):
    return pre + 'class P(Protocol):\n    def m(self%s) -> None: ...\nclass Impl:\n    def m(self%s) -> None: ...\n' % (pm, im)


PA = ['x: P', 'x: object']
cases = [
    # 1 conditional dict comprehensions
    ('1 cond comp', R, D3, "f({**({k: C for k in ['p']} if c else {})})", ['RB']),
    ('1 cond comp two keys', R, D3, "f({**({k: C for k in ['p', 'q']} if c else {'p': C})})", ['RB']),
    ('1 cond comp mixed', R, D3, "f({**({k: C for k in ['p']} if c else {k: A for k in ['q']})})", ['RC']),
    ('1 cond comp var', R + "d = {k: C for k in ['p']} if c else {}\n", D3, 'f({**d})', ['RB']),
    ('1 cond comp kwargs', R, ['**x: type[A]', '**x: type[C]', '**x: object'], "f(**({k: C for k in ['p']} if c else {}))", ['RB']),
    ('1 cond comp arg', R, D3, "f({k: C for k in ['p']} if c else {})", ['RB']),
    # 2 nested starred alternatives keep their length
    ('2 nested star', R + 'xs = [C] if c else [C]\n', ['x: type[C]', 'x: object, y: object'], 'f(*[*xs])', ['RA']),
    ('2 nested tuple', R + 'xs = [C] if c else [C]\n', ['x: type[C]', 'x: object, y: object'], 'f(*(*xs,))', ['RA']),
    ('2 nested var', R + 'xs = [C] if c else [C]\nys = [*xs]\n', ['x: type[C]', 'x: object, y: object'], 'f(*ys)', ['RA']),
    ('2 nested twice', R + 'xs = [C] if c else [C]\n', ['x: type[C], y: type[C]', '*a: object'], 'f(*[*xs, *xs])', ['RA']),
    ('2 nested 1 or 2', R + 'xs = [C] if c else [C, C]\n', ['x: type[C]', 'x: type[C], y: type[C] = C', '*a: object'], 'f(*[*xs])', ['RB']),
    # 3 possibly absent positional arguments shift later ones
    ('3 shift', R + 'xs = [A] if c else []\n', ['x: type[A] = A, y: type[C] = C', 'x: object = None, y: object = None'], 'f(*xs, C)', ['RB']),
    ('3 shift literal', R + 'xs = [A] if c else []\n', ['x: type[A] = A, y: type[C] = C', 'x: object = None, y: object = None'], 'f(*xs, *[C])', ['RB']),
    ('3 shift by name ok', R + 'xs = [A] if c else []\n', ['x: type[A] = A, y: type[C] = C', 'x: object = None, y: object = None'], 'f(*xs, y=C)', ['RA']),
    ('3 shift twice', R + 'xs = [A] if c else []\n', ['x: type[A] = A, y: type[A] = A, z: type[C] = C', '*a: object'], 'f(*xs, *xs, C)', ['RB']),
    ('3 shift nested', R + 'xs = [A] if c else []\n', ['x: type[A] = A, y: type[C] = C', 'x: object = None, y: object = None'], 'f(*[*xs, C])', ['RB']),
    # 4 nested merges keep possible absence
    ('4 d3', R + "d1 = {'x': A} if c else {}\nd2 = {'x': A} if c else {}\nd3 = {**d1, **d2}\n", D3, "f({'x': C, **d3})", ['RC']),
    ('4 d4', R + "d1 = {'x': A} if c else {}\nd2 = {'x': A} if c else {}\nd4 = {**{**d1, **d2}}\n", D3, "f({'x': C, **d4})", ['RC']),
    ('4 inline', R + "d1 = {'x': A} if c else {}\nd2 = {'x': A} if c else {}\n", D3, "f({'x': C, **d1, **d2})", ['RC']),
    ('4 one sure', R + "d1 = {'x': A} if c else {}\nd3 = {**d1, **{'x': A}}\n", D3, "f({'x': C, **d3})", ['RA']),
    # 5 equal keys replace earlier entries
    ('5 class key', '', ['x: dict[type[C], type[A]]', 'x: dict[type[C], type[C]]', 'x: object'], 'f({C: A, C: C})', ['RB']),
    ('5 tuple key', '', ['x: dict[tuple[int, int], type[A]]', 'x: dict[tuple[int, int], type[C]]', 'x: object'], 'f({(1, 2): A, (1, 2): C})', ['RB']),
    ('5 tuple key dstar', '', ['x: dict[tuple[int, int], type[A]]', 'x: dict[tuple[int, int], type[C]]', 'x: object'], 'f({(1, 2): A, **{(1, 2): C}})', ['RB']),
    ('5 cond class key', R + 'k = C if c else C\n', ['x: dict[type[C], type[A]]', 'x: dict[type[C], type[C]]', 'x: object'], 'f({k: A, C: C})', ['RB']),
    ('5 cond str key', R + 'k = "a" if c else "a"\n', D3, 'f({k: A, "a": C})', ['RB']),
    ('5 nested tuple key', '', ['x: dict[tuple[tuple[int, str], int], type[A]]', 'x: dict[tuple[tuple[int, str], int], type[C]]', 'x: object'], 'f({((1, "a"), 2): A, ((1, "a"), 2): C})', ['RB']),
    ('5 different keys', '', ['x: dict[type, type[A]]', 'x: dict[type, type[C]]', 'x: object'], 'f({C: A, B: C})', ['RC']),
    # 6 separate executions are not a cycle
    ('6 wrap wrap', WRAP, D3, "f(wrap(wrap({'x': C})))", ['RB']),
    ('6 wrap x3', WRAP, D3, "f(wrap(wrap(wrap({'x': C}))))", ['RB']),
    ('6 two wraps', WRAP, D3, "f({**wrap({'x': A}), **wrap({'x': C})})", ['RB']),
    ('6 wrap2', WRAP + 'def wrap2(d):\n    return wrap(d)\n', D3, "f(wrap2({'x': C}))", ['RB']),
    ('6 recursive stays finite', 'def rec(d):\n    return {**rec(d)}\n', D3, "f(rec({'x': C}))", None),
    # 7 protocol method compatibility
    ('7 default missing', proto(', x: int = 1', ', x: int'), PA, 'f(Impl())', ['RB']),
    ('7 name differs', proto(', x: int', ', y: int'), PA, 'f(Impl())', ['RB']),
    ('7 kw only impl', proto(', x: int', ', *, x: int'), PA, 'f(Impl())', ['RB']),
    ('7 pos only impl ok?', proto(', x: int, /', ', x: int, /'), PA, 'f(Impl())', ['RA']),
    ('7 extra required', proto(', x: int', ', x: int, y: int'), PA, 'f(Impl())', ['RB']),
    ('7 extra default ok', proto(', x: int', ', x: int, y: int = 0'), PA, 'f(Impl())', ['RA']),
    ('7 args impl ok', proto(', x: int', ', *args: int'), PA, 'f(Impl())', ['RB']),
    ('7 args proto', proto(', *args: int', ', x: int'), PA, 'f(Impl())', ['RB']),
    ('7 args both ok', proto(', *args: int', ', *args: int'), PA, 'f(Impl())', ['RA']),
    ('7 kwargs proto', proto(', **kw: int', ', x: int = 0'), PA, 'f(Impl())', ['RB']),
    ('7 ann mismatch', proto(', x: int', ', x: str'), PA, 'f(Impl())', ['RB']),
    ('7 ann wider ok', proto(', x: int', ', x: object'), PA, 'f(Impl())', ['RA']),
    ('7 ann narrower', proto(', x: object', ', x: int'), PA, 'f(Impl())', ['RB']),
    ('7 same ok', proto(', x: int', ', x: int'), PA, 'f(Impl())', ['RA']),
    ('7 missing method', 'class P(Protocol):\n    def m(self) -> None: ...\nclass Impl: pass\n', PA, 'f(Impl())', ['RB']),
    # 8 quoted annotations
    ('8 quoted mismatch', proto(", x: 'int'", ", x: 'str'"), PA, 'f(Impl())', ['RB']),
    ('8 quoted ok', proto(", x: 'int'", ", x: 'int'"), PA, 'f(Impl())', ['RA']),
    ('8 mixed ok', proto(", x: 'int'", ', x: int'), PA, 'f(Impl())', ['RA']),
    ('8 mixed mismatch', proto(', x: int', ", x: 'str'"), PA, 'f(Impl())', ['RB']),
]
bad = 0
for name, pre, anns, call, exp in cases:
    got = infer(pre + ov(*anns) + call)
    ok = (not str(got).startswith('CRASH')) if exp is None else got == exp
    bad += not ok
    print('%-4s %-24s got=%s exp=%s' % ('OK' if ok else 'BAD', name, got, exp))
print('bad', bad)
