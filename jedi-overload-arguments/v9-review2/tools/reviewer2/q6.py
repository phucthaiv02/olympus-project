from lib import *
import time
c = []
def case(n, prefix, anns, call, exp): c.append((n, prefix + ov(*anns) + call, exp))
case('partial posonly', 'p = functools.partial(f, x=A)\n', ['x: type[A], /, **kw: object', 'x: object = None, **kw: object'], 'p()', ['RB'])
case('partial posonly ok', 'p = functools.partial(f, A)\n', ['x: type[A], /, **kw: object', 'x: object = None, **kw: object'], 'p()', ['RA'])
case('partial kwonly', 'p = functools.partial(f, A)\n', ['*, x: type[A] = A', '*a: object'], 'p()', ['RB'])
case('partial kwonly ok', 'p = functools.partial(f, x=A)\n', ['*, x: type[A] = A', '*a: object'], 'p()', ['RA'])
case('partial plain', 'p = functools.partial(f, x=A)\n', ['x: type[C] = C', 'x: object = None'], 'p()', ['RB'])
case('partial plain2', 'p = functools.partial(f, C)\n', ['x: type[A]', 'x: object'], 'p()', ['RB'])
case('partial none match', 'p = functools.partial(f, C)\n', ['x: type[A]', 'x: type[B]'], 'p()', ['RA', 'RB'])
case('direct none match', '', ['x: type[A]', 'x: type[B]'], 'f(C)', ['RA', 'RB'])
case('partial too many', 'p = functools.partial(f, A, A)\n', ['x: type[A]', '*a: object'], 'p()', ['RB'])
run(c)
