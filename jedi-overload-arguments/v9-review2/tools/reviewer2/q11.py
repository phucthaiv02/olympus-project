from lib import *
import jedi.inference.param as P
orig = P.get_executed_param_names_and_issues
def wrap(fv, args, *a, **k):
    r = orig(fv, args, *a, **k)
    print('  bind', type(args).__name__, k, [(e.string_name, type(e._lazy_value).__name__) for e in r[0]], 'issues=%d' % len(r[1]))
    return r
P.get_executed_param_names_and_issues = wrap
import jedi.inference.signature as S
S_orig = None
c = []
c.append(('direct', ov('x: type[A], y, z: object = None', '*a: object, **k: object') + 'f(A, z=1)', ['RB']))
c.append(('partial', ov('x: type[A], y, z: object = None', '*a: object, **k: object') + 'p = functools.partial(f, A)\np(z=1)', ['RB']))
run(c)
