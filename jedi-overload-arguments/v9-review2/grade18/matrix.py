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
TA, TC = 'x: type[A]', 'x: type[C]'
# first overload / binding
p('bind-pos', [TA, TC], 'f(C)')
p('bind-name', [TA, TC], 'f(x=C)')
p('bind-kwonly', ['*, x: int', 'x: int'], 'f(1)')
p('bind-posonly', ['x: int, /', '*, x: int'], 'f(x=1)')
p('bind-toomany', ['x: int', 'x: int, y: int'], 'f(1, 2)')
p('bind-missing', ['x: int, y: int', 'x: int'], 'f(1)')
p('bind-default', ['x: int, y: int = 1', 'x: int'], 'f(1)')
p('bind-unknown-name', ['x: int', 'y: int'], 'f(y=1)')
p('bind-dup', ['x: int', 'x: object'], 'f(1, x=2)')
p('first-wins', ['x: type[A]', 'x: type[B]'], 'f(B)')
p('none-match', [TA, TC], 'f(int)')
# typing rules
p('ty-int-float', ['x: float', 'x: str'], 'f(1)')
p('ty-float-complex', ['x: complex', 'x: str'], 'f(1.5)')
p('ty-bool-int', ['x: int', 'x: str'], 'f(True)')
p('ty-none', ['x: None', 'x: int'], 'f(None)')
p('ty-opt', ['x: Optional[A]', 'x: object'], 'f(None)')
p('ty-opt-bad', ['x: Optional[A]', 'x: object'], 'f(C())')
p('ty-union', ['x: Union[A, int]', 'x: object'], 'f(1)')
p('ty-union-type', ['x: Union[type[A], type[C]]', 'x: object'], 'f(C)')
p('ty-pipe', ['x: A | int', 'x: object'], 'f(1)')
p('ty-any', ['x: Any', 'x: object'], 'f(1)')
p('ty-bare-type', ['x: type', 'x: object'], 'f(C)')
p('ty-bare-type-inst', ['x: type', 'x: object'], 'f(C())')
p('ty-type-any', ['x: type[Any]', 'x: object'], 'f(C)')
p('ty-str-ann', ["x: 'type[A]'", "x: 'type[C]'"], 'f(C)')
p('ty-callable', ['x: Callable[[], int]', 'x: object'], 'f(lambda: 1)')
p('ty-callable-bad', ['x: Callable[[], int]', 'x: object'], 'f(1)')
p('ty-tuple-order', ['x: tuple[int, str]', 'x: object'], 'f(("s", 1))')
p('ty-tuple-len', ['x: tuple[int, str]', 'x: object'], 'f((1, "s", 2))')
p('ty-tuple-ok', ['x: tuple[int, str]', 'x: object'], 'f((1, "s"))')
p('ty-tuple-empty', ['x: tuple[()]', 'x: object'], 'f(())')
p('ty-seq-str', ['x: Sequence[str]', 'x: object'], 'f("ab")')
p('ty-mapping', ['x: Mapping[str, type[A]]', 'x: Mapping[str, type[C]]'], 'f({"k": C})')
p('ty-iter-set', ['x: Iterable[type[A]]', 'x: Iterable[type[C]]'], 'f({C})')
p('ty-set', ['x: set[type[A]]', 'x: set[type[C]]'], 'f({C})')
p('ty-frozenset-ann', ['x: frozenset[int]', 'x: object'], 'f(frozenset([1]))')
p('ty-list-class-inst', ['x: list[int]', 'x: object'], 'f(["s"])')
p('ty-subclass-inst', ['x: A', 'x: object'], 'f(B())')
p('ty-typealias', ['x: Alias', 'x: object'], 'f(C)', pre='Alias = Type[C]\n')
# TypeVar sentence
b = 'T = TypeVar("T", bound="A")\n'
k = 'T = TypeVar("T", A, int)\n'
ks = 'T = TypeVar("T", "A", int)\n'
p('tv-bound-ok', ['x: type[T]', TC], 'f(B)', pre=b)
p('tv-bound-bad', ['x: type[T]', TC], 'f(C)', pre=b)
p('tv-bound-later', ['x: type[T]', TC], 'f(B)', pre='T = TypeVar("T", bound="Later")\nclass Later: pass\nclass B2(Later): pass\n')
p('tv-bound-inst', ['x: T', 'x: object'], 'f(C())', pre=b)
p('tv-constr-str', ['x: type[T]', TC], 'f(B)', pre=ks)
p('tv-constr-bad', ['x: type[T]', TC], 'f(C)', pre=ks)
p('tv-unbound', ['x: type[T]', TC], 'f(C)', pre='T = TypeVar("T")\n')
p('tv-same', ['x: type[T], y: type[T]', 'x: object, y: object'], 'f(A, int)', pre=k)
p('tv-same-ok', ['x: type[T], y: type[T]', 'x: object, y: object'], 'f(A, B)', pre=k)
p('tv-same-name', ['x: type[T], y: type[T]', 'x: object, y: object'], 'f(y=int, x=A)', pre=k)
p('tv-same-star', ['*x: type[T]', '*x: object'], 'f(A, int)', pre=k)
p('tv-same-list', ['x: list[type[T]], y: type[T]', 'x: object, y: object'], 'f([A], int)', pre=k)
p('tv-same-tuple', ['x: tuple[type[T], type[T]]', 'x: object'], 'f((A, int))', pre=k)
p('tv-same-dict', ['x: dict[type[T], type[T]]', 'x: object'], 'f({A: int})', pre=k)
p('tv-same-star-unpack', ['x: type[T], y: type[T]', 'x: object, y: object'], 'f(*[A, int])', pre=k)
p('tv-same-dstar', ['x: type[T], y: type[T]', 'x: object, y: object'], 'f(**{"x": A, "y": int})', pre=k)
p('tv-same-alt', ['x: type[T], y: type[T]', 'x: object, y: object'], 'f(A if c else int, int)', pre=k)
p('tv-two-vars', ['x: type[T], y: type[U]', 'x: object, y: object'], 'f(A, int)', pre=k + 'U = TypeVar("U", A, int)\n')
p('tv-bound-two', ['x: type[T], y: type[T]', 'x: object, y: object'], 'f(A, B)', pre='T = TypeVar("T", bound=A)\n')
# multiple values / unknown length
p('alt-ifexp', [TA, TC], 'f(A if c else C)')
p('alt-ifexp-ok', ['x: type[A]', TC], 'f(A if c else B)')
p('alt-branch', [TA, TC], 'if c:\n    v = A\nelse:\n    v = C\nf(v)')
p('star-list', [TC, 'x: object'], 'f(*xs)', pre='xs: list[type[C]]\n')
p('star-tuple-unk', [TC, 'x: object'], 'f(*xs)', pre='xs: tuple[type[C], ...]\n')
p('star-tuple-fixed', [TC, 'x: object'], 'f(*xs)', pre='xs: tuple[type[C]]\n')
p('star-to-star', ['*x: type[C]', '*x: object'], 'f(*xs)', pre='xs: list[type[C]]\n')
p('star-unresolved', [TC, 'x: object'], 'f(*undefined)')
p('dstar-unk', [TC, 'x: object'], 'f(**d)', pre='d: dict[str, type[C]]\n')
p('dstar-to-dstar', ['**k: type[C]', '**k: object'], 'f(**d)', pre='d: dict[str, type[C]]\n')
# literal content
p('lit-tuple-var', ['x: tuple[type[C], type[C]]', 'x: object'], 't = (C, C)\nf(t)')
p('lit-star-tuple', ['x: type[A], y: type[C]', 'x: object, y: object'], 'f(*(A, C))')
p('lit-star-var', ['x: type[A], y: type[C]', 'x: object, y: object'], 't = [A, C]\nf(*t)')
p('lit-nested', ['x: type[A], y: type[C]', 'x: object, y: object'], 'b = [C]\nf(*[A, *b])')
p('lit-dstar', ['x: type[A], y: type[C]', 'x: object, y: object'], 'f(**{"y": C, "x": A})')
p('lit-dstar-nested', ['x: type[A], y: type[C]', 'x: object, y: object'], 'f(**{"x": A, **{"y": C}})')
p('key-equal', ['x: dict[int, type[A]]', 'x: dict[int, type[C]]', 'x: object'], 'f({1: A, True: C})')
p('key-equal-float', ['x: dict[int, type[A]]', 'x: dict[int, type[C]]', 'x: object'], 'f({1: A, 1.0: C})')
p('key-kept-bool', ['x: dict[bool, type[C]]', 'x: dict[int, type[C]]'], 'f({True: C, 1: C})')
p('key-kept-int', ['x: dict[bool, type[C]]', 'x: dict[int, type[C]]'], 'f({1: C, True: C})')
p('key-kept-float', ['x: dict[float, type[C]]', 'x: dict[int, type[C]]'], 'f({1.0: C, 1: C})')
p('key-kept-dstar', ['x: dict[bool, type[C]]', 'x: dict[int, type[C]]'], 'f({1: C, **{True: C}})')
p('key-maybe', ['x: dict[str, type[A]]', 'x: dict[str, type[C]]', 'x: object'], 'd = {"k": A} if c else {}\nf({"k": C, **d})')
p('key-maybe-after', ['x: dict[str, type[A]]', 'x: dict[str, type[C]]', 'x: object'], 'd = {"k": A} if c else {}\nf({**d, "k": C})')
p('key-var', ['x: dict[str, type[A]]', 'x: dict[str, type[C]]', 'x: object'], 'k = "a"\nf({"a": A, k: C})')
p('key-class', ['x: dict[type[C], type[A]]', 'x: dict[type[C], type[C]]', 'x: object'], 'f({C: A, C: C})')
p('comp-empty', ['x: list[type[C]]', 'x: object'], 'f([A for _ in []])')
p('comp-empty-tuple', ['x: list[type[C]]', 'x: object'], 'f([A for _ in ()])')
p('comp-second-for', ['x: list[type[C]]', 'x: object'], 'f([A for _ in [1] for _ in []])')
p('comp-dict', ['x: dict[str, type[C]]', 'x: object'], 'f({"k": A for _ in []})')
p('comp-set', ['x: set[type[C]]', 'x: object'], 'f({A for _ in []})')
p('comp-gen', ['x: Iterable[type[C]]', 'x: object'], 'f(A for _ in [])')
p('comp-nonempty', ['x: list[type[C]]', 'x: object'], 'f([A for _ in [1]])')
p('comp-unpack', ['x: list[type[C]]', 'x: object'], 'f([C, *[A for _ in []]])')
p('comp-var-empty', ['x: list[type[C]]', 'x: object'], 'e = []\nf([A for _ in e])')
p('maybe-star', ['x: type[C]', 'x: type[A], y: type[C]', 'x: object, y: object = 1'], 'xs = [A] if c else []\nf(*xs, C)')
p('maybe-star-fix', ['x: type[A], y: type[C]', 'x: object, y: object'], 'xs = [A] if c else [A]\nf(*xs, C)')
p('maybe-star-tuple', ['x: type[C]', 'x: object'], 'xs = (A,) if c else ()\nf(*xs, C)')
res = {k: infer(v) for k, v in P.items()}
print('JSON' + json.dumps(res))
