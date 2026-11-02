from lib import *
c = []
def case(n, prefix, anns, call, exp): c.append((n, prefix + ov(*anns) + call, exp))
def L(n): return '[' + ', '.join(str(i) for i in range(n)) + ']'
case('3lvl small', '', ['x: dict[int, type[C]]', 'x: object'], 'f({a*b*c: C for a in [1, 2] for b in [1] for c in [1]})', ['RA'])
for n in (5, 10, 15, 20):
    case('3lvl %d' % n, '', ['x: dict[int, type[C]]', 'x: object'], 'f({a: C for a in %s for b in %s for c in %s})' % (L(n), L(n), L(n)), ['RA'])
case('2lvl 60', '', ['x: dict[int, type[C]]', 'x: object'], 'f({a: C for a in %s for b in %s})' % (L(60), L(60)), ['RA'])
case('list 3lvl 20', '', ['x: list[type[C]]', 'x: object'], 'f([C for a in %s for b in %s for c in %s])' % (L(20), L(20), L(20)), ['RA'])
case('pm', 'class K:\n    @overload\n    def m(self, x: type[A]) -> RA: ...\n    @overload\n    def m(self, x: type[C]) -> RB: ...\n    def m(self, x): ...\n    pm = functools.partialmethod(m, C)\n', [], 'K().pm()', ['RB'])
run(c)
