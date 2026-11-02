from lib import *
c = []
def case(n, prefix, anns, call, exp): c.append((n, prefix + ov(*anns) + call, exp))
def L(n): return '[' + ', '.join(str(i) for i in range(n)) + ']'
for n in (8, 10, 12, 16, 20):
    case('dictcomp %dx%d' % (n, n), '', ['x: dict[int, type[C]]', 'x: object'], 'f({a: C for a in %s for b in %s})' % (L(n), L(n)), ['RA'])
case('dictcomp 1x100', '', ['x: dict[int, type[C]]', 'x: object'], 'f({a: C for a in %s})' % L(100), ['RA'])
case('dictcomp str 10x10', '', ['x: dict[str, type[C]]', 'x: object'], 'f({a: C for a in "abcdefghij" for b in "abcdefghij"})', ['RA'])
case('dictcomp ** kwargs 20x20', '', ['**k: type[C]', '**k: object'], 'f(**{a: C for a in %s for b in %s})' % (str(["k%d" % i for i in range(20)]), L(20)), ['RA'])
case('dictcomp ** kwargs 20x20 bad', '', ['**k: type[A]', '**k: object'], 'f(**{a: C for a in %s for b in %s})' % (str(["k%d" % i for i in range(20)]), L(20)), ['RB'])
run(c)
