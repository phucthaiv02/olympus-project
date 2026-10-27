import sys
sys.argv = ['x']
from textwrap import dedent
import jedi
from probe import HEADER, ov

def infer(code):
    code = HEADER + dedent(code).strip('\n')
    lines = code.splitlines()
    try:
        return sorted(d.name for d in jedi.Script(code).infer(len(lines), len(lines[-1])))
    except Exception as e:
        import traceback; traceback.print_exc()
        return 'CRASH %s: %s' % (type(e).__name__, e)

D = ['x: dict[str, type[A]]', 'x: dict[str, type[C]]', 'x: object']
K = lambda kt: ['x: dict[%s, type[A]]' % kt, 'x: dict[%s, type[C]]' % kt, 'x: object']
R = 'import random\nc = random.random()\n'
cases = [
 # sequences / positions
 ('seq empty-or-A then C', R+"xs = [A] if c else []\n", ['x: type[A] = A, y: type[C] = C', 'x: object = None, y: object = None'], 'f(*xs, C)', ['RB']),
 ('seq C-or-empty then C', R+"xs = [C] if c else []\n", ['x: type[C], y: type[C] = C', 'x: object = None, y: object = None'], 'f(*xs, C)', ['RA']),
 ('seq len 1/2 both type C', R+"xs = [C] if c else [C, C]\n", ['x: type[C], y: type[C] = C', 'x: object'], 'f(*xs)', ['RA']),
 ('seq len 1/2 second A', R+"xs = [C] if c else [C, A]\n", ['x: type[C], y: type[C] = C', 'x: object, y: object = None'], 'f(*xs)', ['RB']),
 ('seq len 1/2 overload needs 2', R+"xs = [C] if c else [C, C]\n", ['x: type[C], y: type[C]', 'x: object, y: object = None'], 'f(*xs)', ['RB']),
 ('two stars', R+"xs = [A] if c else []\nys = [C] if c else []\n", ['*args: type[A]', '*args: object'], 'f(*xs, *ys)', ['RB']),
 ('two stars ok', R+"xs = [A] if c else []\nys = [B] if c else []\n", ['*args: type[A]', '*args: object'], 'f(*xs, *ys)', ['RA']),
 ('nested literal positions', R+"xs = [A] if c else []\n", ['x: type[A] = A, y: type[C] = C', 'x: object = None, y: object = None'], 'f(*[*xs, C])', ['RB']),
 ('nested alt same', R+"xs = [C] if c else [C]\n", ['x: type[C]', 'x: object, y: object'], 'f(*[*xs])', ['RA']),
 ('nested alt tuple', R+"xs = (C,) if c else (C,)\n", ['x: type[C]', 'x: object, y: object'], 'f(*(*xs,))', ['RA']),
 ('tuple annotation alt lengths', R+"xs = [C] if c else [C, C]\n", ['x: tuple[type[C]]', 'x: object'], 'f((*xs,))', ['RB']),
 ('tuple annotation alt same len', R+"xs = [C] if c else [B]\n", ['x: tuple[type[C]]', 'x: tuple[Union[type[C], type[A]]]', 'x: object'], 'f((*xs,))', ['RB']),
 ('list alt any length', R+"xs = [C] if c else []\n", ['x: list[type[C]]', 'x: object'], 'f([*xs])', ['RA']),
 ('star var list typed', 'def g() -> list[type[C]]: ...\n', ['*args: type[C]', '*args: object'], 'f(*g())', ['RA']),
 ('star unknown', '', ['x: type[C]', 'x: object'], 'f(*unknown_name)', ['RA', 'RB']),
 ('kw after shift', R+"xs = [A] if c else []\n", ['x: type[A] = A, *, k: type[C]', 'x: object = None, *, k: object'], 'f(*xs, k=C)', ['RA']),
 # dict comprehensions
 ('empty comp', '', D, 'f({"k": C for _ in []})', ['RA']),
 ('empty comp in **', '', D, 'f({"k": A, **{"k": C for _ in ()}})', ['RA']),
 ('cond comp in **', R, D, "f({**({k: C for k in ['p']} if c else {})})", ['RB']),
 ('cond comp mixed', R, D, "f({**({k: C for k in ['p']} if c else {'p': A})})", ['RC']),
 ('comp values', '', D, "f({k: C for k in ['a', 'b']})", ['RB']),
 ('comp keys int', '', D, "f({k: C for k in [1]})", ['RC']),
 ('kwargs cond comp', R, ['**kw: type[C]', '**kw: object'], "f(**({k: C for k in ['p']} if c else {}))", ['RA']),
 # comprehension clauses (solutionQuality v3 High)
 ('empty 2nd for', '', D, 'f({k: C for k in ["a"] for _ in []})', ['RA']),
 ('empty 2nd for after if', '', D, 'f({k: C for k in ["a"] if k for _ in ()})', ['RA']),
 ('empty 3rd for', '', D, 'f({k: C for k in ["a"] for j in [1] for _ in {}})', ['RA']),
 ('nonempty 2 fors', '', D, 'f({k: C for k in ["a"] for _ in [1]})', ['RB']),
 ('empty 2nd for in **', '', D, 'f({"k": A, **{k: C for k in ["k"] for _ in []}})', ['RA']),
 ('dep 2nd for', '', D, 'f({k: C for k in ["a"] for _ in k})', ['RB']),
 ('async-free list comp 2nd', '', ['x: list[type[A]]', 'x: object'], 'f([C for _ in [1] for _ in []])', ['RA']),
 # variant limits (solutionQuality v3 Medium)
 ('5 stars then C', R+"xs = [A] if c else []\n", ['a: type[A] = A, b: type[A] = A, c: type[A] = A, d: type[A] = A, e: type[A] = A, g: type[C] = C', 'a: object = None, b: object = None, c: object = None, d: object = None, e: object = None, g: object = None'], 'f(*xs, *xs, *xs, *xs, *xs, C)', ['RB']),
 ('5 nested stars then C', R+"xs = [A] if c else []\n", ['a: type[A] = A, b: type[A] = A, c: type[A] = A, d: type[A] = A, e: type[A] = A, g: type[C] = C', 'a: object = None, b: object = None, c: object = None, d: object = None, e: object = None, g: object = None'], 'f(*[*xs, *xs, *xs, *xs, *xs, C])', ['RB']),
 ('4 stars then C', R+"xs = [A] if c else []\n", ['a: type[A] = A, b: type[A] = A, c: type[A] = A, d: type[A] = A, g: type[C] = C', 'a: object = None, b: object = None, c: object = None, d: object = None, g: object = None'], 'f(*xs, *xs, *xs, *xs, C)', ['RB']),
 ('5 stars all fine', R+"xs = [A] if c else []\n", ['*args: type[A]', '*args: object'], 'f(*xs, *xs, *xs, *xs, *xs, A)', ['RA']),
 ('5 tuple stars', R+"xs = [C] if c else []\n", ['x: tuple[type[C], ...]', 'x: object'], 'f((*xs, *xs, *xs, *xs, *xs))', ['RA']),
 # nested absence
 ('nested absence', R+"d1 = {'x': A} if c else {}\nd2 = {'x': A} if c else {}\nd3 = {**d1, **d2}\n", D, "f({'x': C, **d3})", ['RC']),
 ('nested absence one sure', R+"d1 = {'x': A} if c else {}\nd2 = {'x': A}\nd3 = {**d1, **d2}\n", D, "f({'x': C, **d3})", ['RA']),
 ('nested absence 3 levels', R+"d1 = {'x': A} if c else {}\nd2 = {**d1}\nd3 = {**d2}\n", D, "f({'x': C, **d3})", ['RC']),
 ('kwargs nested absence', R+"d1 = {'x': A} if c else {}\nd2 = {'x': A} if c else {}\n", ['x: type[C]', 'x: object = None'], "f(**{'x': C, **{**d1, **d2}})", ['RB']),
 # keys
 ('class key', '', K('type[C]'), 'f({C: A, C: C})', ['RB']),
 ('tuple key', '', K('tuple[int, int]'), 'f({(1, 2): A, (1, 2): C})', ['RB']),
 ('tuple key differs', '', K('tuple[int, int]'), 'f({(1, 2): A, (1, 3): C})', ['RC']),
 ('nested tuple key', '', K('tuple'), 'f({(1, (2, "a")): A, (1, (2, "a")): C})', ['RB']),
 ('1 True key', '', K('int'), 'f({1: C, True: A})', ['RA']),
 ('1.0 key', '', K('float'), 'f({1: A, 1.0: C})', ['RB']),
 ('cond equal keys', R+"k = 'a' if c else 'a'\n", D, "f({'a': A, k: C})", ['RB']),
 ('cond different keys', R+"k = 'a' if c else 'b'\n", D, "f({'a': A, k: C})", ['RC']),
 ('None key', '', K('None'), 'f({None: A, None: C})', ['RB']),
 ('class key in **', '', K('type[C]'), 'f({C: A, **{C: C}})', ['RB']),
 ('tuple key in dict()', '', K('tuple[int, int]'), 'f(dict([((1, 2), A), ((1, 2), C)]))', ['RB']),
 ('class key differs', '', K('type'), 'f({C: A, B: C})', ['RC']),
 # recursion
 ('wrap wrap', "def wrap(d):\n    return {**d}\n", D, "f(wrap(wrap({'x': C})))", ['RB']),
 ('wrap wrap A', "def wrap(d):\n    return {**d}\n", D, "f(wrap(wrap({'x': A})))", ['RA']),
 ('wrap overwrite', "def wrap(d):\n    return {'x': A, **d}\n", D, "f(wrap(wrap({'x': C})))", ['RB']),
 ('dict wrap', "def wrap(d):\n    return dict(d)\n", D, "f(wrap(wrap({'x': C})))", ['RB']),
 ('recursive fn', "def g(d):\n    return g({**d})\n", D, "f(g({'x': C}))", None),
 ('loop self', "d = {'x': C}\nfor _ in range(3):\n    d = {**d}\n", D, "f(d)", None),
]
bad = 0
for name, pre, anns, call, exp in cases:
    got = infer(pre + ov(*anns) + call)
    ok = exp is None and not isinstance(got, str) or got == exp
    bad += not ok
    print('%-4s %-30s got=%s exp=%s' % ('OK' if ok else 'BAD', name, got, exp))
print('bad', bad)
