import sys, traceback
from textwrap import dedent
import jedi

HEADER = """\
from typing import overload, Any, Type, TypeVar, Union, Protocol
import random
class A: pass
class B(A): pass
class C: pass
class RA: pass
class RB: pass
class RC: pass
"""

def ov(*anns):
    code = ''
    for a, r in zip(anns, ('RA', 'RB', 'RC')):
        code += '@overload\ndef f(%s) -> %s: ...\n' % (a, r)
    return code + 'def f(*args, **kwargs): ...\n'

def infer(code):
    code = HEADER + dedent(code).strip('\n')
    lines = code.splitlines()
    try:
        return sorted(d.name for d in jedi.Script(code).infer(len(lines), len(lines[-1])))
    except Exception as e:
        return 'CRASH %s: %s' % (type(e).__name__, e)

CASES = []
def case(name, prefix, anns, call, expected):
    CASES.append((name, prefix + ov(*anns) + call, expected))

# 1 conditional dict comprehension crash
case('1 cond-comp crash', '', ['x: dict[str, type[A]]', 'x: dict[str, type[C]]', 'x: object'],
     "f({**({k: C for k in ['p']} if random.random() else {})})", ['RB'])
# 2 nested starred alternatives
case('2 nested star alt', "xs = [C] if random.random() else [C]\n", ['x: type[C]', 'x: object, y: object'],
     'f(*[*xs])', ['RA'])
# 3 possibly absent positional shift
case('3 absent positional shift', "xs = [A] if random.random() else []\n",
     ['x: type[A] = A, y: type[C] = C', 'x: object = None, y: object = None'], 'f(*xs, C)', ['RB'])
# 4 nested dict merges absence
case('4 nested merge absence',
     "d1 = {'x': A} if random.random() else {}\nd2 = {'x': A} if random.random() else {}\nd3 = {**d1, **d2}\n",
     ['x: dict[str, type[A]]', 'x: dict[str, type[C]]', 'x: object'], "f({'x': C, **d3})", ['RC'])
# 5 non-scalar keys
case('5a class key overwrite', '', ['x: dict[type[C], type[A]]', 'x: dict[type[C], type[C]]', 'x: object'],
     'f({C: A, C: C})', ['RB'])
case('5b tuple key overwrite', '', ['x: dict[tuple[int, int], type[A]]', 'x: dict[tuple[int, int], type[C]]', 'x: object'],
     'f({(1, 2): A, (1, 2): C})', ['RB'])
# 6 recursion guard
case('6 wrap(wrap())', "def wrap(d):\n    return {**d}\n",
     ['x: dict[str, type[A]]', 'x: dict[str, type[C]]', 'x: object'], "f(wrap(wrap({'x': C})))", ['RB'])
case('6 inline', "",
     ['x: dict[str, type[A]]', 'x: dict[str, type[C]]', 'x: object'], "f({**{**{'x': C}}})", ['RB'])
# 7 protocol signatures
P7 = "class P(Protocol):\n    def m(self, x: int = 1) -> None: ...\nclass Impl:\n    def m(self, x: int) -> None: ...\n"
case('7a protocol default', P7, ['x: P', 'x: object'], 'f(Impl())', ['RB'])
P7b = "class P(Protocol):\n    def m(self, x: int) -> None: ...\nclass Impl:\n    def m(self, y: int) -> None: ...\n"
case('7b protocol name', P7b, ['x: P', 'x: object'], 'f(Impl())', ['RB'])
P7ok = "class P(Protocol):\n    def m(self, x: int) -> None: ...\nclass Impl:\n    def m(self, x: int) -> None: ...\n"
case('7 protocol ok control', P7ok, ['x: P', 'x: object'], 'f(Impl())', ['RA'])
P7c = "class P(Protocol):\n    def m(self, x: int) -> None: ...\nclass Impl:\n    def m(self, x: str) -> None: ...\n"
case('7c protocol unquoted mismatch control', P7c, ['x: P', 'x: object'], 'f(Impl())', ['RB'])
# 8 quoted protocol annotation
P8 = "class P(Protocol):\n    def m(self, x: 'int') -> None: ...\nclass Impl:\n    def m(self, x: 'str') -> None: ...\n"
case('8 protocol quoted', P8, ['x: P', 'x: object'], 'f(Impl())', ['RB'])
# auto review: empty dict comprehension
case('AR empty comp', '', ['x: dict[str, type[A]]', 'x: dict[str, type[C]]', 'x: object'],
     'f({"k": C for _ in []})', ['RA'])
# v1 review leftovers
case('v1 fwd TypeVar bound', 'T = TypeVar("T", bound="A")\n', ['x: type[T]', 'x: object'], 'f(B)', ['RA'])

for name, code, exp in CASES:
    got = infer(code)
    print('%-5s %-38s got=%s expected=%s' % ('OK' if got == exp else 'BAD', name, got, exp))
