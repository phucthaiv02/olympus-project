from lib import *
c = []
def case(n, prefix, anns, call, exp): c.append((n, prefix + ov(*anns) + call, exp))
def L(n): return '[' + ', '.join(str(i) for i in range(n)) + ']'
for n in (45, 50, 55, 101, 150):
    case('listcomp %d' % n, '', ['x: list[int]', 'x: object'], 'f([a for a in %s])' % L(n), ['RA'])
    case('dictcomp %d' % n, '', ['x: dict[int, int]', 'x: object'], 'f({a: a for a in %s})' % L(n), ['RA'])
run(c)
