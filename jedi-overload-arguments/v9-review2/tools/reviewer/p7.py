from lib import *
import time
c = []
def case(n, prefix, anns, call, exp): c.append((n, prefix + ov(*anns) + call, exp))
two = ['x: type[A]', 'x: object']
case('self append', 'xs = [A]\nxs.append(xs)\n', two, 'f(*xs)', None)
case('rebind star', 'xs = []\nxs = [*xs, C]\n', two, 'f(*xs)', ['RB'])
case('rebind dstar', 'd = {}\nd = {**d, "x": C}\n', two, 'f(**d)', ['RB'])
case('self dict', 'd = {"x": C}\nd["y"] = d\n', two, 'f(**d)', None)
case('star None', '', two, 'f(*None)', None)
case('dstar None', '', two, 'f(**None)', None)
case('star int', '', two, 'f(*1)', None)
case('dstar list', '', two, 'f(**[1, 2])', None)
case('star class', '', two, 'f(*A)', None)
case('dstar int key', '', two, 'f(**{1: A})', ['RA', 'RB'])
case('lambda', '', two, 'f((lambda: 1)())', ['RB'])
case('walrus', '', two, 'f(y := C)', ['RB'])
case('kw walrus', '', two, 'f(x=(y := C))', ['RB'])
case('starred genexpr', '', two, 'f(*(C for _ in [1]))', ['RB'])
case('starred dictcomp', '', ['x: str', 'x: object'], 'f(*{k: C for k in ["a"]})', ['RA'])
case('starred setcomp', '', ['x: type[C]', 'x: object, y: object'], 'f(*{C for _ in [1, 2]})', ['RA'])
case('starred listcomp', '', ['x: type[C], y: type[C]', 'x: object'], 'f(*[C for _ in [1, 2]])', ['RA'])
L = '[' + ', '.join(str(i) for i in range(25)) + ']'
case('big nested comp', '', ['x: dict[int, type[C]]', 'x: object'], 'f({a*b*c: C for a in %s for b in %s for c in %s})' % (L, L, L), ['RA'])
case('big star', '', ['*a: type[C]', '*a: object'], 'f(*[%s])' % ', '.join(['C'] * 300), ['RA'])
case('many cond stars', ''.join('x%d = [A] if random.random() else []\n' % i for i in range(12)), ['*a: type[A]', '*a: object'], 'f(%s)' % ', '.join('*x%d' % i for i in range(12)), ['RA'])
case('many unknown stars', 'def g() -> tuple[type[A], ...]: ...\n', ['a: type[A], b: type[A]=A, c: type[A]=A, d: type[A]=A, e: type[A]=A, *r: type[A]', '*a: object'], 'f(%s)' % ', '.join(['*g()'] * 8), ['RA'])
for n, code, exp in c:
    t = time.time(); got = infer(code); dt = time.time() - t
    print('%-4s %-24s %.1fs got=%s expected=%s' % ('OK' if exp is None or got == exp else 'BAD', n, dt, got, exp), flush=True)
