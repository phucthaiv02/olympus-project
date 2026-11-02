from lib import *
c = []
def case(n, prefix, anns, call, exp): c.append((n, prefix + ov(*anns) + call, exp))
D = "def g() -> dict[str, type[C]]: ...\nd = g()\n"
case('ukey to named+kw', D, ['x: type[A] = A, **kw: object', 'x: object = None, **kw: object'], 'f(**d)', ['RB'])
case('ukey to named+kw lit', 'k = str(1)\n', ['x: type[A] = A, **kw: object', 'x: object = None, **kw: object'], 'f(**{k: C})', ['RB'])
case('ukey comp', '', ['x: type[A] = A, **kw: object', 'x: object = None, **kw: object'], 'f(**{k: C for k in unknown_names})', ['RB'])
P = 'p = functools.partial(f, '
case('partial pos', '', ['x: type[A], y: type[C]', 'x: object, y: object'], P + 'A)\np(C)', ['RA'])
case('partial pos bad', '', ['x: type[A], y: type[A]', 'x: object, y: object'], P + 'A)\np(C)', ['RB'])
case('partial kw', '', ['x: type[A], y: type[C]', 'x: object, y: object'], P + 'y=C)\np(A)', ['RA'])
case('partial kw override', '', ['x: type[A], y: type[C]', 'x: object, y: object'], P + 'y=A)\np(A, y=C)', ['RA'])
case('partial kw override bad', '', ['x: type[A], y: type[A]', 'x: object, y: object'], P + 'y=A)\np(A, y=C)', ['RB'])
case('partial star', 'xs = [A] if random.random() else []\n', ['x: type[A] = A, y: type[C] = C', 'x: object = None, y: object = None'], P + '*xs)\np(C)', ['RB'])
case('partial star call', 'xs = [A] if random.random() else []\n', ['x: type[A] = A, y: type[C] = C', 'x: object = None, y: object = None'], 'p = functools.partial(f)\np(*xs, C)', ['RB'])
case('partial dstar', '', ['x: type[A], y: type[C]', 'x: object, y: object'], P + '**{"y": C})\np(A)', ['RA'])
case('partial dstar override', '', ['x: type[A], y: type[C]', 'x: object, y: object'], P + '**{"y": A})\np(A, **{"y": C})', ['RA'])
case('partial dstar maybe', '', ['x: type[A], y: type[C]', 'x: object, y: object'], P + 'y=C)\np(A, **({"y": A} if random.random() else {}))', ['RB'])
case('partial kw var', 'd = {"y": A}\n', ['x: type[A], y: type[C]', 'x: object, y: object'], P + 'y=C)\np(A, **d)', ['RB'])
case('partial both pos', '', ['x: type[A], y: type[C]', 'x: object, y: object'], P + 'A)\np(*[C])', ['RA'])
run(c)
