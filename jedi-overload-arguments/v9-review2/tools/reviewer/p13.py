import jedi
from jedi.inference.gradual.base import comprehension_items
from jedi.inference.arguments import mapping_items
def L(n): return '[' + ', '.join(str(i) for i in range(n)) + ']'
for n in (2, 10, 20, 25, 30, 35):
    code = 'class C: pass\nd = {a: C for a in %s for b in %s}\nd' % (L(n), L(n))
    s = jedi.Script(code)
    v = s._get_module_context().infer_node(s._module_node.children[-1].children[0]) if False else None
    defs = s.infer(3, 1)
    val = defs[0]._name.infer()
    val = list(val)[0]
    items = comprehension_items(val)
    keys, values = mapping_items(val)
    print(n, None if items is None else len(items), len(keys), sum(1 for k in keys if not k.infer()), sum(1 for v in values if not v.infer()))
