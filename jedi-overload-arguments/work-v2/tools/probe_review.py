# Usage: python -I probe_review.py <jedi worktree> [filter]
# Probes for every issue of the v1 review (reviewer + Auto Review).
import sys
sys.path.insert(0, sys.argv[1])
import jedi
FILTER = sys.argv[2] if len(sys.argv) > 2 else ''
HEADER = """\
import random, functools
from typing import overload, Any, Type, TypeVar, Union, Iterable, List, Mapping, Sequence, Protocol
class A: pass
class B(A): pass
class C: pass
class RA: pass
class RB: pass
class RC: pass
"""
def ov(*anns, sig='x: %s'):
    code = ''
    for a, r in zip(anns, ('RA', 'RB', 'RC')):
        code += '@overload\ndef f(%s) -> %s: ...\n' % (sig % a, r)
    return code + 'def f(*args, **kwargs): ...\n'
def infer(code):
    code = HEADER + code
    lines = code.splitlines()
    return sorted(d.name for d in jedi.Script(code).infer(len(lines), len(lines[-1])))
def nested(n, inner='p=C'):
    s = 'dict(%s)' % inner
    for _ in range(n - 1):
        s = 'dict(**%s)' % s
    return s
LIST = ov('list[type[A]]', 'list[type[C]]', 'object')
SEQ = ov('Sequence[type[A]]', 'Sequence[type[C]]', 'object')
DICT = ov('dict[str, type[A]]', 'dict[str, type[C]]', 'object')
VAR = ov('type[A]', 'type[C]', sig='*x: %s')
KW = ov('type[A]', 'type[C]', sig='**x: %s')
ONE = ov('type[A]', 'type[C]')
MAPPING = '''class M:
    def keys(self) -> list[str]: ...
    def __getitem__(self, key: str) -> type[C]: ...
'''
MAPPING2 = '''class M2:
    def keys(self):
        return ["x"]
    def __getitem__(self, key):
        return C
'''
MAPPING3 = '''class M3:
    def keys(self):
        return ["x", "y"]
    def __getitem__(self, key):
        return C if key == "x" else A
'''
PROTO = '''class P(Protocol):
    def meth(self, a: int) -> int: ...
'''
probes = [
    # Auto Review: starred elements inside literals
    ('bot1 list [C, *[A]]', LIST + 'f([C, *[A]])', ['RC']),
    ('bot1 list [*[A], C]', LIST + 'f([*[A], C])', ['RC']),
    ('bot1 list [C, *[C]] pos', LIST + 'f([C, *[C]])', ['RB']),
    ('bot1 seq (C, *[A])', SEQ + 'f((C, *[A]))', ['RC']),
    ('bot2 dict {k:C, **{m:A}}', DICT + 'f({"k": C, **{"m": A}})', ['RC']),
    ('bot2 dict {**{m:A}, k:C}', DICT + 'f({**{"m": A}, "k": C})', ['RC']),
    ('bot2 dict {k:C, **dict(m=A)}', DICT + 'f({"k": C, **dict(m=A)})', ['RC']),
    # Reviewer 2 / Auto Review 3: nested starred contents at the call
    ('rev2 var *[*[A], C]', VAR + 'f(*[*[A], C])', ['RA', 'RB']),
    ('rev2 var *[A, *[C]]', VAR + 'f(*[A, *[C]])', ['RA', 'RB']),
    ('rev2 var *[*[C]]', VAR + 'f(*[*[C]])', ['RB']),
    ('rev2 var *(*[C], C) pos', VAR + 'f(*(*[C], C))', ['RB']),
    ('rev2 one *[*[C]]', ONE + 'f(*[*[C]])', ['RB']),
    ('rev2 one *[*[A]]', ONE + 'f(*[*[A]])', ['RA']),
    # Auto Review 4: depth
    ('bot4 kw nested 5', KW + 'f(**%s)' % nested(5), ['RB']),
    ('bot4 kw nested 12', KW + 'f(**%s)' % nested(12), ['RB']),
    ('bot4 dict nested 12', DICT + 'f(%s)' % nested(12, 'k=C'), ['RB']),
    # Reviewer 1: forward-referenced TypeVar bounds/constraints
    ('rev1 bound "A" with B', 'T = TypeVar("T", bound="A")\n' + ov('type[T]', 'object') + 'f(B)', ['RA']),
    ('rev1 bound "A" with C', 'T = TypeVar("T", bound="A")\n' + ov('type[T]', 'object') + 'f(C)', ['RB']),
    ('rev1 constr "A","int" B', 'T = TypeVar("T", "A", "int")\n' + ov('type[T]', 'object') + 'f(B)', ['RA']),
    ('rev1 constr "A","int" C', 'T = TypeVar("T", "A", "int")\n' + ov('type[T]', 'object') + 'f(C)', ['RB']),
    ('rev1 bound "D" later', 'T = TypeVar("T", bound="D")\n' + ov('type[T]', 'object') + 'class D: pass\nclass E(D): pass\nf(E)', ['RA']),
    # Reviewer 3: custom mapping in dict(...), empty comprehension
    ('rev3 dict(M()) kw', MAPPING + KW + 'f(**dict(M()))', ['RB']),
    ('rev3 dict(M()) one unknown key', MAPPING + ONE + 'f(**dict(M()))', ['RA', 'RB']),
    ('rev3 dict(M2()) one', MAPPING2 + ONE + 'f(**dict(M2()))', ['RB']),
    ('rev3 dict(M2()) kw', MAPPING2 + KW + 'f(**dict(M2()))', ['RB']),
    ('rev3 dict(M3()) one', MAPPING3 + ONE + 'f(**dict(M3()))', ['RA', 'RB']),
    ('rev3 dict(M2()) dict ann', MAPPING2 + DICT + 'f(dict(M2()))', ['RB']),
    ('rev3 dict(M3()) dict ann', MAPPING3 + DICT + 'f(dict(M3()))', ['RC']),
    ('rev3 **M2()', MAPPING2 + ONE + 'f(**M2())', ['RB']),
    ('rev3 ann dict kw', 'def g() -> dict[str, type[C]]: ...\n' + KW + 'f(**g())', ['RB']),
    ('rev3 ann dict kw A', 'def g() -> dict[str, type[A]]: ...\n' + KW + 'f(**g())', ['RA']),
    ('rev3 empty comp kw', KW + 'f(**{k: C for k in []})', ['RA']),
    ('rev3 empty comp list', VAR + 'f(*[C for k in []])', ['RA']),
    ('rev3 nonempty comp kw', KW + 'f(**{k: C for k in "ab"})', ['RB']),
    # Reviewer 4: conditional dict overwritten by a later keyword
    ('rev4 dict(cond, p=C)', DICT + 'f(dict({"p": A} if random.random() else {}, p=C))', ['RB']),
    ('rev4 kw dict(cond, p=C)', KW + 'f(**dict({"p": A} if random.random() else {}, p=C))', ['RB']),
    ('rev4 one dict(cond, x=C)', ONE + 'f(**dict({"x": A} if random.random() else {}, x=C))', ['RB']),
    ('rev4 {**cond, p: C}', DICT + 'f({**({"p": A} if random.random() else {}), "p": C})', ['RB']),
    # Reviewer 5: fixed vs variable length tuples
    ('rev5 tuple[int,...] arg', ov('tuple[int]', 'tuple[int, ...]') + 'def g() -> tuple[int, ...]: ...\nf(g())', ['RB']),
    ('rev5 tuple[int] arg', ov('tuple[int]', 'tuple[int, ...]') + 'def g() -> tuple[int]: ...\nf(g())', ['RA']),
    ('rev5 literal (1,)', ov('tuple[int]', 'tuple[int, ...]') + 'f((1,))', ['RA']),
    ('rev5 tuple(list)', ov('tuple[int]', 'tuple[int, ...]') + 'f(tuple([1]))', ['RB']),
    ('rev5 list arg', ov('tuple[int]', 'object') + 'def g() -> list[int]: ...\nf(g())', ['RB']),
    # Reviewer 6: partial keyword uncertainty
    ('rev6 partial unknown keys kw', 'def k1() -> str: ...\ndef k2() -> str: ...\n' + KW + 'functools.partial(f, **{k1(): A})(**{k2(): C})', ['RA', 'RB']),
    ('rev6 two unknown keys kw', 'def k1() -> str: ...\ndef k2() -> str: ...\n' + KW + 'f(**{k1(): A}, **{k2(): C})', ['RA', 'RB']),
    ('rev6 partial unknown keys', 'def k1() -> str: ...\ndef k2() -> str: ...\n' + ONE + 'functools.partial(f, **{k1(): A})(**{k2(): C})', ['RA', 'RB']),
    ('rev6 partial both maybe', 'm = {"x": A} if random.random() else {}\nn = {"x": A} if random.random() else {}\n' + ov('type[A]', 'object', sig='x: %s') + 'functools.partial(f, **m)(**n)', ['RA', 'RB']),
    ('rev6 partial both maybe dflt', 'm = {"x": A} if random.random() else {}\nn = {"x": A} if random.random() else {}\n' + ov('type[A] = A', 'object', sig='x: %s') + 'functools.partial(f, **m)(**n)', ['RA']),
    ('rev6 partial kw replaced', ONE + 'functools.partial(f, x=A)(x=C)', ['RB']),
    # Reviewer 7: protocol members
    ('rev7 proto ok', PROTO + 'class X:\n    def meth(self, a: int) -> int: ...\n' + ov('P', 'object') + 'f(X())', ['RA']),
    ('rev7 proto int attr', PROTO + 'class X:\n    meth = 1\n' + ov('P', 'object') + 'f(X())', ['RB']),
    ('rev7 proto extra param', PROTO + 'class X:\n    def meth(self, a: int, b: int) -> int: ...\n' + ov('P', 'object') + 'f(X())', ['RB']),
    ('rev7 proto extra default', PROTO + 'class X:\n    def meth(self, a: int, b: int = 1) -> int: ...\n' + ov('P', 'object') + 'f(X())', ['RA']),
    ('rev7 proto fewer params', PROTO + 'class X:\n    def meth(self) -> int: ...\n' + ov('P', 'object') + 'f(X())', ['RB']),
    ('rev7 proto missing', PROTO + 'class X: pass\n' + ov('P', 'object') + 'f(X())', ['RB']),
    ('rev7 proto *args', PROTO + 'class X:\n    def meth(self, *args) -> int: ...\n' + ov('P', 'object') + 'f(X())', ['RA']),
    ('rev7 proto kwonly req', PROTO + 'class X:\n    def meth(self, a: int, *, b: int) -> int: ...\n' + ov('P', 'object') + 'f(X())', ['RB']),
    ('sq4 proto kwonly narrower', 'class P2(Protocol):\n    def meth(self, *, x: object) -> int: ...\nclass X:\n    def meth(self, *, x: int) -> int: ...\n' + ov('P2', 'object') + 'f(X())', ['RB']),
    ('sq4 proto kwonly same', 'class P2(Protocol):\n    def meth(self, *, x: int) -> int: ...\nclass X:\n    def meth(self, *, x: object) -> int: ...\n' + ov('P2', 'object') + 'f(X())', ['RA']),
    ('sq4 proto kwonly as positional', 'class P2(Protocol):\n    def meth(self, *, x: object) -> int: ...\nclass X:\n    def meth(self, x: int) -> int: ...\n' + ov('P2', 'object') + 'f(X())', ['RB']),
    ('rev7 proto attr ok', 'class Q(Protocol):\n    name: str\nclass X:\n    name: str = ""\n' + ov('Q', 'object') + 'f(X())', ['RA']),
    ('rev7 proto attr bad', 'class Q(Protocol):\n    name: str\nclass X:\n    name: int = 1\n' + ov('Q', 'object') + 'f(X())', ['RB']),
    # solutionQuality v2: alternatives that may lack a key or a position
    ('sq1 {k: C, **maybe A}', DICT + 'd = {"k": A} if random.random() else {}\nf({"k": C, **d})', ['RC']),
    ('sq1 {k: C, **maybe C}', DICT + 'd = {"k": C} if random.random() else {}\nf({"k": C, **d})', ['RB']),
    ('sq1 {k: A, **maybe C}', DICT + 'd = {"k": C} if random.random() else {}\nf({"k": A, **d})', ['RC']),
    ('sq1 {**maybe A, k: C}', DICT + 'd = {"k": A} if random.random() else {}\nf({**d, "k": C})', ['RB']),
    ('sq1 kw {k: C, **maybe A}', KW + 'd = {"k": A} if random.random() else {}\nf(**{"k": C, **d})', ['RA', 'RB']),
    ('sq1 nested {k: C, **{**d}}', DICT + 'd = {"k": A} if random.random() else {}\nf({"k": C, **{**d}})', ['RC']),
    ('sq1 nested {k: C, **dict(d)}', DICT + 'd = {"k": A} if random.random() else {}\nf({"k": C, **dict(d)})', ['RC']),
    ('sq1 dict({k: C}, **d)', DICT + 'd = {"k": A} if random.random() else {}\nf(dict({"k": C}, **d))', ['RC']),
    ('sq1 dict(d, k=C)', DICT + 'd = {"k": A} if random.random() else {}\nf(dict(d, k=C))', ['RB']),
    ('sq1 kw f(**{**d}) missing x', ov('type[A]', 'object') + 'd = {"x": A} if random.random() else {}\nf(**{**d})', ['RA', 'RB']),
    ('sq2 nested *[*xs]', ov('type[A], y: type[A]', 'type[A]') + 'xs = (A,) if random.random() else (A, A)\nf(*[*xs])', ['RA', 'RB']),
    ('sq3 *undefined', VAR + 'f(*undefined_name)', ['RA', 'RB']),
    ('sq3 **undefined', KW + 'f(**undefined_name)', ['RA', 'RB']),
    ('sq2 two params maybe', ov('type[A], y: type[A]', 'type[A]') + 'xs = (A,) if random.random() else (A, A)\nf(*xs)', ['RA', 'RB']),
    ('sq2 one param maybe', ov('type[A]', 'type[A], y: type[A]') + 'xs = (A,) if random.random() else (A, A)\nf(*xs)', ['RA', 'RB']),
    ('sq2 default ok', ov('type[A], y: type[A] = A', 'object') + 'xs = (A,) if random.random() else (A, A)\nf(*xs)', ['RA']),
    ('sq2 var maybe', VAR + 'xs = (A,) if random.random() else (A, C)\nf(*xs)', ['RA', 'RB']),
    ('sq2 same length', ov('type[A], y: type[A]', 'object') + 'xs = (A, B) if random.random() else (A, A)\nf(*xs)', ['RA']),
    # FP panel: X | None
    ('fp pipe None', ov('type[A] | None', 'object') + 'f(None)', ['RA']),
    ('fp pipe C', ov('type[A] | None', 'object') + 'f(C)', ['RB']),
    ('fp pipe B', ov('type[A] | None', 'object') + 'f(B)', ['RA']),
]
fails = 0
for name, code, exp in probes:
    if FILTER not in name:
        continue
    try:
        got = infer(code)
    except Exception as e:
        got = 'ERR %s: %s' % (type(e).__name__, e)
    ok = got == exp
    fails += not ok
    print('%s %-30s got=%s exp=%s' % ('PASS' if ok else 'FAIL', name, got, exp))
print('FAILS', fails)
