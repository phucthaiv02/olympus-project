from lib import *
import time
c = []
def case(n, prefix, anns, call, exp): c.append((n, prefix + ov(*anns) + call, exp))
classes = ''.join('class K%d: pass\n' % i for i in range(10))
def perf(nt, nc):
    tvs = ''.join('T%d = TypeVar("T%d", %s)\n' % (j, j, ', '.join('K%d' % i for i in range(nc))) for j in range(nt))
    params = ', '.join('a%d: T%d, b%d: T%d' % (j, j, j, j) for j in range(nt))
    case('perf %dx%d mismatch' % (nt, nc), classes + tvs, [params, '*a: object'], 'f(%s)' % ', '.join(['K0(), K1()'] * nt), ['RB'])
perf(3, 5); perf(4, 5); perf(3, 10); perf(5, 4)
# AnyStr-like realistic: 4 AnyStr-ish params
case('anystr x4', 'from typing import AnyStr\nS = TypeVar("S", str, bytes)\nU = TypeVar("U", int, float)\n', ['a: AnyStr, b: S, c: U, d: AnyStr', '*a: object'], 'f("a", b"b", 1, b"x")', ['RB'])
# odd params / crashes
case('lambda arg', '', ['x: type[A]', 'x: object'], 'f(lambda: C)', ['RB'])
case('lambda overload-like', 'g = lambda x, /, *, y: x\n', ['x: type[A]', 'x: object'], 'f(g(C, y=1))', ['RB'])
case('only star', '', ['*, x: type[A] = A', '*a: object'], 'f()', ['RA'])
case('posonly kwonly mix', '', ['a: type[A], /, b: type[A], *, c: type[A]', 'a: object, /, b: object, *, c: object'], 'f(A, A, c=A)', ['RA'])
case('posonly kwonly mix kw b', '', ['a: type[A], /, b: type[A], *, c: type[A]', 'a: object, /, b: object, *, c: object'], 'f(A, b=A, c=A)', ['RA'])
case('posonly kwonly mix bad', '', ['a: type[A], /, b: type[A], *, c: type[A]', 'a: object, /, b: object, *, c: object'], 'f(A, A, A)', ['RA', 'RB'])
case('posonly name collision kw', '', ['a: type[A], /, **kw: type[C]', '*a: object, **kw: object'], 'f(A, **{"a": C})', ['RA'])
case('method posonly', 'class K:\n' + ov('self, x: type[A], /', 'self, x: object', name='m', indent='    '), [], 'K().m(x=A)', ['RB'])
case('method posonly ok', 'class K:\n' + ov('self, x: type[A], /', 'self, x: object', name='m', indent='    '), [], 'K().m(A)', ['RA'])
case('classmethod kwonly', 'class K:\n' + ov('cls, *, x: type[A]', 'cls, x: object', name='m', indent='    ', deco='@classmethod'), [], 'K.m(A)', ['RB'])
case('staticmethod posonly', 'class K:\n' + ov('x: type[A], /', 'x: object', name='m', indent='    ', deco='@staticmethod'), [], 'K.m(x=A)', ['RB'])
case('partial posonly', 'p = functools.partial(f, x=A)\n', ['x: type[A], /, **kw: object', 'x: object = None, **kw: object'], 'p()', ['RB'])
case('partial kwonly', 'p = functools.partial(f, A)\n', ['*, x: type[A] = A', '*a: object'], 'p()', ['RB'])
case('self named kw', 'class K:\n' + ov('self, x: type[A], **kw: type[C]', 'self, x: object, **kw: object', name='m', indent='    '), [], 'K().m(A, self=C)', ['RA'])
case('recursive alias', 'from typing import List, Union\nJ = Union[int, List["J"]]\n', ['x: J', 'x: object'], 'f([1, [2]])', None)
case('recursive class gen', 'from typing import Generic\nS = TypeVar("S")\nclass N(Generic[S]):\n    nxt: "N[S]"\n', ['x: N[int]', 'x: object'], 'f(N())', None)
case('nested typevar generic', 'from typing import Generic\nS = TypeVar("S", A, C)\nclass N(Generic[S]): pass\ndef g() -> N[C]: ...\n', ['x: N[S], y: S', 'x: object, y: object'], 'f(g(), A())', ['RB'])
case('typevar constraint is typevar', 'S = TypeVar("S")\nT = TypeVar("T", list[S], int)\n', ['x: T', 'x: object'], 'f([1])', None)
case('constraint unknown', 'T = TypeVar("T", Undefined1, Undefined2)\n', ['x: T, y: T', 'x: object, y: object'], 'f(A(), C())', None)
case('typevar bad args', 'T = TypeVar("T", A)\n', ['x: T', 'x: object'], 'f(A())', None)
case('typevar no name', 'T = TypeVar()\n', ['x: T', 'x: object'], 'f(A())', None)
case('default star empty', '', ['*args', 'x: object'], 'f(*())', ['RA'])
case('kwonly no default unknown star', 'def g() -> tuple[type[A], ...]: ...\n', ['*a: type[A], k: type[A] = A', '*a: object'], 'f(*g(), k=A)', ['RA'])
for n, code, exp in c:
    t = time.time(); got = infer(code); dt = time.time() - t
    print('%-4s %-28s %.1fs got=%s expected=%s' % ('OK' if exp is None and not str(got).startswith(('CRASH','TIMEOUT')) or got == exp else 'BAD', n, dt, got, exp), flush=True)
