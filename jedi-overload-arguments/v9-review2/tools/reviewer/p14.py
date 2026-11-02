import jedi
from jedi.inference.gradual.base import comprehension_items
from jedi.inference.arguments import mapping_items
def L(n): return '[' + ', '.join(str(i) for i in range(n)) + ']'
for desc, comp in [('1lvl10', '{a: C for a in %s}' % L(10)), ('1lvl30', '{a: C for a in %s}' % L(30)),
                   ('2lvl 3x3', '{a: C for a in %s for b in %s}' % (L(3), L(3))),
                   ('2lvl 7x2', '{a: C for a in %s for b in %s}' % (L(7), L(2))),
                   ('2lvl 2x7', '{a: C for a in %s for b in %s}' % (L(2), L(7))),
                   ('list 2lvl 10x10', '[a for a in %s for b in %s]' % (L(10), L(10))),
                   ('2lvl10 abc', '{a: C for a in "abcdefghij" for b in "abcdefghij"}'),
                   ]:
    code = 'class C: pass\nd = %s\nd' % comp
    s = jedi.Script(code)
    val = list(s.infer(3, 1)[0]._name.infer())[0]
    items = comprehension_items(val)
    print(desc, None if items is None else len(items))
