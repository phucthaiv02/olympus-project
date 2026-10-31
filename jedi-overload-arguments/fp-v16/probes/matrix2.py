# Probe matrix from description v13 plus the proposed TypeVar sentence: one group per sentence.
import json, sys
from textwrap import dedent
import jedi

HEADER = """\
from typing import (overload, Any, Callable, Dict, Iterable, List, Mapping, Optional, Protocol,
                    Sequence, Tuple, Type, TypeVar, Union)
import random
class A: pass
class B(A): pass
class C: pass
class RA: pass
class RB: pass
class RC: pass
c = random.random()
"""
def ov(*sigs):
    code = ''
    for s, r in zip(sigs, ('RA', 'RB', 'RC')):
        code += '@overload\ndef f(%s) -> %s: ...\n' % (s, r)
    return code + 'def f(*args, **kwargs): ...\n'
def infer(code):
    code = HEADER + dedent(code).strip('\n')
    lines = code.splitlines()
    try:
        return ','.join(sorted(d.name for d in jedi.Script(code).infer(len(lines), len(lines[-1]))))
    except Exception as e:
        return 'CRASH ' + type(e).__name__

P = {}
def p(id, sigs, call, pre=''):
    P[id] = pre + ov(*sigs) + call

from_ty = 'from typing import Literal, NewType, Final, Annotated, Generic, NamedTuple, TypedDict, ClassVar\nimport collections.abc as abc\n'
p('cb-lambda', ['x: Callable[[], int]', 'x: object'], 'f(lambda: 1)')
p('cb-int', ['x: Callable[[], int]', 'x: object'], 'f(1)')
p('cb-func', ['x: Callable[[], int]', 'x: object'], 'def g(): return 1\nf(g)')
p('cb-func-first-int', ['x: int', 'x: Callable[[], int]'], 'def g(): return 1\nf(g)')
p('cb-class', ['x: Callable[[], A]', 'x: object'], 'f(A)')
p('cb-bare', ['x: Callable', 'x: object'], 'f(len)')
p('cb-bare-int', ['x: Callable', 'x: object'], 'f(1)')
p('cb-order', ['x: Callable[..., int]', 'x: str'], 'f("s")')
p('cb-abc', ['x: abc.Callable[[], int]', 'x: object'], 'f(1)', pre=from_ty)
p('proto-ok', ['x: P', 'x: object'], 'f(Impl())', pre='class P(Protocol):\n    def m(self) -> int: ...\nclass Impl:\n    def m(self) -> int: return 1\n')
p('proto-bad', ['x: P', 'x: object'], 'f(C())', pre='class P(Protocol):\n    def m(self) -> int: ...\n')
p('proto-attr', ['x: P', 'x: object'], 'f(Impl())', pre='class P(Protocol):\n    a: int\nclass Impl:\n    a: int = 1\n')
p('proto-type', ['x: type[P]', 'x: object'], 'f(Impl)', pre='class P(Protocol):\n    def m(self) -> int: ...\nclass Impl:\n    def m(self) -> int: return 1\n')
p('lit-ok', ['x: Literal[1]', 'x: object'], 'f(1)', pre=from_ty)
p('lit-bad', ['x: Literal[1]', 'x: object'], 'f(2)', pre=from_ty)
p('newtype', ['x: N', 'x: object'], 'f(N(1))', pre=from_ty + 'N = NewType("N", int)\n')
p('newtype-raw', ['x: N', 'x: object'], 'f(1)', pre=from_ty + 'N = NewType("N", int)\n')
p('annotated', ['x: Annotated[int, 1]', 'x: object'], 'f(1)', pre=from_ty)
p('annotated-bad', ['x: Annotated[int, 1]', 'x: object'], 'f("s")', pre=from_ty)
p('typeddict', ['x: TD', 'x: object'], 'f({"a": 1})', pre=from_ty + 'class TD(TypedDict):\n    a: int\n')
p('namedtuple', ['x: tuple[int, str]', 'x: object'], 'f(NT(1, "s"))', pre=from_ty + 'class NT(NamedTuple):\n    a: int\n    b: str\n')
p('generic-user', ['x: Box[int]', 'x: object'], 'f(Box[int]())', pre=from_ty + 'T = TypeVar("T")\nclass Box(Generic[T]): pass\n')
p('generic-user-bad', ['x: Box[int]', 'x: object'], 'f(Box[str]())', pre=from_ty + 'T = TypeVar("T")\nclass Box(Generic[T]): pass\n')
p('sub-list', ['x: list[int]', 'x: object'], 'f(L())', pre='class L(list): pass\n')
p('list-ann-var', ['x: list[A]', 'x: list[C]'], 'v: list[C]\nf(v)')
p('list-invariant', ['x: list[A]', 'x: object'], 'v: list[B]\nf(v)')
p('seq-covariant', ['x: Sequence[A]', 'x: object'], 'v: Sequence[B]\nf(v)')
p('abc-seq', ['x: abc.Sequence[type[A]]', 'x: abc.Sequence[type[C]]'], 'f([C])', pre=from_ty)
p('abc-map', ['x: abc.Mapping[str, type[A]]', 'x: abc.Mapping[str, type[C]]'], 'f({"k": C})', pre=from_ty)
p('union-none-pipe', ['x: int | None', 'x: object'], 'f(None)')
p('ret-call', ['x: int', 'x: str'], 'def g() -> str: ...\nf(g())')
p('attr-arg', ['x: int', 'x: str'], 'class K:\n    a: str = "s"\nf(K.a)')
p('self-method', ['self, x: type[A]', 'self, x: type[C]'], 'f', pre='')
p('bytes-float', ['x: bytes', 'x: str'], 'f(b"s")')
p('none-opt-order', ['x: int', 'x: Optional[int]'], 'f(None)')
p('object-first', ['x: object', 'x: int'], 'f(1)')
p('str-seq-char', ['x: Sequence[int]', 'x: object'], 'f("ab")')
p('dict-iter', ['x: Iterable[str]', 'x: object'], 'f({"a": 1})')
p('dict-iter-bad', ['x: Iterable[int]', 'x: object'], 'f({"a": 1})')
p('range-iter', ['x: Iterable[int]', 'x: object'], 'f(range(3))')
p('gen-iter', ['x: Iterable[int]', 'x: object'], 'f(i for i in [1])')
res = {k: infer(v) for k, v in P.items()}
print('JSON' + json.dumps(res))
