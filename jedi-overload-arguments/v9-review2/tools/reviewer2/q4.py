from lib import *
import time
c = []
def case(n, prefix, anns, call, exp): c.append((n, prefix + ov(*anns) + call, exp))
TC = 'T = TypeVar("T", A, C)\n'
case('quoted list', TC, ['x: list["T"], y: T', 'x: object, y: object'], 'f([A()], C())', ['RB'])
case('quoted inner', TC, ['x: "list[T]", y: T', 'x: object, y: object'], 'f([A()], C())', ['RB'])
case('quoted Type', TC, ['x: type["T"], y: T', 'x: object, y: object'], 'f(A, C())', ['RB'])
case('TypeVar quoted args', TC, ['*a: "T"', '*a: object'], 'f(A(), C())', ['RB'])
# perf: many constrained typevars
classes = ''.join('class K%d: pass\n' % i for i in range(10))
tvs = ''.join('T%d = TypeVar("T%d", %s)\n' % (j, j, ', '.join('K%d' % i for i in range(10))) for j in range(6))
params = ', '.join('a%d: T%d' % (j, j) for j in range(6))
case('perf 6x10', classes + tvs, [params, '*a: object'], 'f(%s)' % ', '.join(['K9()'] * 6), ['RA'])
params2 = ', '.join('a%d: T%d, b%d: T%d' % (j, j, j, j) for j in range(6))
case('perf 6x10 mismatch', classes + tvs, [params2, '*a: object'], 'f(%s)' % ', '.join(['K9(), K8()'] * 6), ['RB'])
# promotions
case('int->float', '', ['x: float', 'x: object'], 'f(1)', ['RA'])
case('bool->float', '', ['x: float', 'x: object'], 'f(True)', ['RA'])
case('float->int no', '', ['x: int', 'x: object'], 'f(1.5)', ['RB'])
case('complex->float no', '', ['x: float', 'x: object'], 'f(1j)', ['RB'])
case('int->complex', '', ['x: complex', 'x: object'], 'f(1)', ['RA'])
case('list[float] ints', '', ['x: list[float]', 'x: object'], 'f([1, 2.5])', ['RA'])
case('type[float] int', '', ['x: type[float]', 'x: object'], 'f(int)', ['RA'])
case('str->float no', '', ['x: float', 'x: object'], 'f("1")', ['RB'])
case('user float', 'class float: pass\n', ['x: float', 'x: object'], 'f(1)', ['RB'])
case('int ann var', 'def g() -> int: ...\n', ['x: float', 'x: object'], 'f(g())', ['RA'])
case('float bound', 'T = TypeVar("T", bound=float)\n', ['x: T', 'x: object'], 'f(1)', ['RA'])
case('dict[float,..] int keys', '', ['x: dict[float, int]', 'x: object'], 'f({1: 1})', ['RA'])
case('float int subclass', 'class MyInt(int): pass\n', ['x: float', 'x: object'], 'f(MyInt())', ['RA'])
for n, code, exp in c:
    t = time.time(); got = infer(code); dt = time.time() - t
    print('%-4s %-28s %.1fs got=%s expected=%s' % ('OK' if got == exp else 'BAD', n, dt, got, exp), flush=True)
