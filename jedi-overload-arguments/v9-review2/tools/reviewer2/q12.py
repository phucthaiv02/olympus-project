from lib import *
import time
c = []
def case(n, prefix, anns, call, exp): c.append((n, prefix + ov(*anns) + call, exp))
classes = ''.join('class K%d: pass\n' % i for i in range(10))
def perf(nt, nc):
    tvs = ''.join('T%d = TypeVar("T%d", %s)\n' % (j, j, ', '.join('K%d' % i for i in range(nc))) for j in range(nt))
    params = ', '.join('a%d: T%d, b%d: T%d' % (j, j, j, j) for j in range(nt))
    case('perf %dx%d mismatch' % (nt, nc), classes + tvs, [params, '*a: object'], 'f(%s)' % ', '.join(['K0(), K1()'] * nt), ['RB'])
perf(4, 10); perf(5, 6); perf(6, 5); perf(5, 10)
for n, code, exp in c:
    t = time.time(); got = infer(code); dt = time.time() - t
    print('%-4s %-28s %.1fs got=%s expected=%s' % ('OK' if got == exp else 'BAD', n, dt, got, exp), flush=True)
