from lib import *
import time
c = []
def case(n, prefix, anns, call, exp):
    pre = [l for l in prefix.splitlines(True) if 'partial' not in l]
    post = [l for l in prefix.splitlines(True) if 'partial' in l]
    c.append((n, ''.join(pre) + ov(*anns) + ''.join(post) + call, exp))
case('partial posonly', 'p = functools.partial(f, x=A)\n', ['x: type[A], /, y: object = None, **kw: object', 'x: object = None, y: object = None, **kw: object'], 'p(y=1)', ['RB'])
case('partial posonly ok', 'p = functools.partial(f, A)\n', ['x: type[A], /, y: object = None, **kw: object', 'x: object = None, y: object = None, **kw: object'], 'p(y=1)', ['RA'])
case('partial kwonly', 'p = functools.partial(f, A)\n', ['y: object = None, *, x: type[A] = A', '*a: object'], 'p(y=1)', ['RB'])
case('partial kwonly ok', 'p = functools.partial(f, x=A)\n', ['y: object = None, *, x: type[A] = A', '*a: object, **k: object'], 'p(1)', ['RA'])
case('partial ukey', 'def g() -> dict[str, type[C]]: ...\np = functools.partial(f, **g())\n', ['y: object, x: type[A] = A, **kw: object', 'y: object, x: object = None, **kw: object'], 'p(1)', ['RB'])
case('partial ukey call', 'def g() -> dict[str, type[C]]: ...\np = functools.partial(f, 1)\n', ['y: object, x: type[A] = A, **kw: object', 'y: object, x: object = None, **kw: object'], 'p(**g())', ['RB'])
case('partial ukey then kw', 'def g() -> dict[str, type[C]]: ...\np = functools.partial(f, 1, **g())\n', ['y: object, x: type[A] = A, **kw: object', 'y: object, x: object = None, **kw: object'], 'p(x=A)', ['RA'])
case('partial posonly kw', 'p = functools.partial(f, x=A)\n', ['x: type[A], /, **kw: type[C]', 'x: object, /, **kw: object'], 'p(A)', ['RB'])
case('partial c constraint', 'T = TypeVar("T", A, C)\np = functools.partial(f, A())\n', ['x: T, y: T', 'x: object, y: object'], 'p(C())', ['RB'])
run(c)
