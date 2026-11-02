from lib import *
c = []
def case(n, prefix, anns, call, exp): c.append((n, prefix + ov(*anns) + call, exp))
for n in (150, 200, 300):
    case('big star %d' % n, '', ['*a: type[C]', '*a: object'], 'f(*[%s])' % ', '.join(['C'] * n), ['RA'])
    case('big pos %d' % n, '', ['*a: type[C]', '*a: object'], 'f(%s)' % ', '.join(['C'] * n), ['RA'])
case('unk stars 7', 'def g() -> tuple[type[A], ...]: ...\n', ['a: type[C]', '*a: object'], 'f(%s)' % ', '.join(['*g()'] * 7), ['RB'])
case('unk stars 8 single ov', 'def g() -> tuple[type[A], ...]: ...\n', ['*a: object', 'a: type[C]'], 'f(%s)' % ', '.join(['*g()'] * 8), ['RA'])
case('unk stars 8 var', 'def g() -> tuple[type[A], ...]: ...\nxs = g()\n', ['a: type[C]', '*a: object'], 'f(%s)' % ', '.join(['*xs'] * 8), ['RB'])
case('unk stars 8 ann param', 'def h(xs: tuple[type[A], ...]):\n    return f(%s)\n' % ', '.join(['*xs'] * 8), ['a: type[C]', '*a: object'], 'h(())', ['RB'])
case('lists 8', 'def g() -> list[type[A]]: ...\n', ['a: type[C]', '*a: object'], 'f(%s)' % ', '.join(['*g()'] * 8), ['RB'])
run(c)
