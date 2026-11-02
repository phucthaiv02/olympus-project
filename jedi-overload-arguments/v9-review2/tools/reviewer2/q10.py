from lib import *
c = []
def P(n, anns, part, call, exp): c.append((n, ov(*anns) + 'p = functools.partial(f, %s)\np(%s)' % (part, call), exp))
P('missing unannotated', ['x: type[A], y, z: object = None', '*a: object, **k: object'], 'A', 'z=1', ['RB'])
P('too many', ['x: type[A], z: object = None', '*a: object, **k: object'], 'A, A', 'z=1', ['RB'])
P('unexpected kw', ['x: type[A]', '*a: object, **k: object'], 'A', 'w=1', ['RB'])
P('dup partial kw call pos', ['y: type[A] = A, x: type[A] = A', '*a: object, **k: object'], 'y=C', 'A', ['RB'])
P('posonly by kw via partial', ['x: type[A], /, z: object = None', '*a: object, **k: object'], 'x=A', 'z=1', ['RB'])
P('dup via call', ['y: object = None, x: type[A] = A', '*a: object, **k: object'], 'A', 'y=1', ['RB'])
# method partial
c.append(('method partial dup', 'class K:\n' + ov('self, y: type[A] = A, x: type[A] = A', 'self, *a: object, **k: object', name='m', indent='    ') + 'p = functools.partial(K().m, y=C)\np(A)', ['RB']))
run(c)
