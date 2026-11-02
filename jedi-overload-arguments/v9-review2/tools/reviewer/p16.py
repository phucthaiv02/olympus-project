from lib import *
import jedi
from jedi.inference.gradual.base import comprehension_items
c = []
def case(n, prefix, anns, call, exp): c.append((n, prefix + ov(*anns) + call, exp))
def L(n): return '[' + ', '.join(str(i) for i in range(n)) + ']'
for n in (40, 60, 100):
    case('listcomp %d' % n, '', ['x: list[int]', 'x: object'], 'f([a for a in %s])' % L(n), ['RA'])
    case('dictcomp %d' % n, '', ['x: dict[int, type[C]]', 'x: object'], 'f({a: C for a in %s})' % L(n), ['RA'])
    case('dictcomp val %d' % n, '', ['x: dict[str, int]', 'x: object'], 'f({"k": a for a in %s})' % L(n), ['RA'])
run(c)
for n in (40, 60, 100):
    s = jedi.Script('class C: pass\nd = {a: C for a in %s}\nd' % L(n))
    val = list(s.infer(3, 1)[0]._name.infer())[0]
    items = comprehension_items(val)
    print(n, None if items is None else (len(items), sum(1 for k, v in items if not k)))
