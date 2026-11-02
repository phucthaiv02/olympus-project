from lib import *
c = []
def case(n, prefix, anns, call, exp): c.append((n, prefix + ov(*anns) + call, exp))
for n in (10, 20, 30, 40, 60, 100):
    case('big star %d' % n, '', ['*a: type[C]', '*a: object'], 'f(*[%s])' % ', '.join(['C'] * n), ['RA'])
    case('big pos %d' % n, '', ['*a: type[C]', '*a: object'], 'f(%s)' % ', '.join(['C'] * n), ['RA'])
    case('big list %d' % n, '', ['x: list[type[C]]', 'x: object'], 'f([%s])' % ', '.join(['C'] * n), ['RA'])
for k in (1, 2, 3, 4, 6, 8):
    case('unk stars %d' % k, 'def g() -> tuple[type[A], ...]: ...\n', ['a: type[C]', '*a: object'], 'f(%s)' % ', '.join(['*g()'] * k), ['RB'])
case('comp mult', '', ['x: dict[int, type[C]]', 'x: object'], 'f({a*b: C for a in [1, 2] for b in [3]})', ['RA'])
case('comp plain', '', ['x: dict[int, type[C]]', 'x: object'], 'f({a: C for a in [1, 2]})', ['RA'])
L = '[' + ', '.join(str(i) for i in range(25)) + ']'
case('comp 25x25', '', ['x: dict[int, type[C]]', 'x: object'], 'f({a: C for a in %s for b in %s})' % (L, L), ['RA'])
case('comp 0 key', '', ['x: dict[int, type[C]]', 'x: object'], 'f({a*b: C for a in [0, 1] for b in [0]})', ['RA'])
run(c)
