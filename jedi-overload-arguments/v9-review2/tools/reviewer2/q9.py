from lib import *
c = []
O1 = ['y: object = None, *, x: type[A] = A', '*a: object']
O2 = ['y: type[C] = C, x: type[A] = A', '*a: object, **k: object']
c.append(('direct', ov(*O1) + 'f(A, y=1)', ['RA', 'RB']))
c.append(('partial', ov(*O1) + 'p = functools.partial(f, A)\np(y=1)', ['RA', 'RB']))
c.append(('direct2', ov(*O2) + 'f(A, y=C)', ['RB']))
c.append(('partial2', ov(*O2) + 'p = functools.partial(f, A)\np(y=C)', ['RB']))
c.append(('partial3', ov(*O2) + 'p = functools.partial(f, y=C)\np(A)', ['RB']))
c.append(('partial4 ok', ov(*O2) + 'p = functools.partial(f, y=C)\np(x=A)', ['RA']))
run(c)
