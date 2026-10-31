"""``*`` unpacking from every kind of source against fixed-arity, default
and variadic overloads. An unpacked value of unknown length may give any
number of arguments, so only overloads that fit all of them match."""
import sys
sys.path.insert(0, '/w')
from probe import ov, infer

R = 'c = random.random()\n'
SRC = '''\
def ht() -> tuple[type[C], ...]: ...
def hl() -> list[type[C]]: ...
def ft() -> tuple[type[C], type[C]]: ...
def ft1() -> tuple[type[C]]: ...
def mix() -> tuple[type[A], type[C]]: ...
'''
ONE = ['x: type[C]', '*args: type[C]', '*args: object']
cases = [
    # homogeneous tuple / list: any number of arguments
    ('ht one', SRC, ONE, 'f(*ht())', ['RB']),
    ('hl one', SRC, ONE, 'f(*hl())', ['RB']),
    ('ht var', SRC + 'xs = ht()\n', ONE, 'f(*xs)', ['RB']),
    ('ht default', SRC, ['x: type[C] = C', '*args: type[C]', '*args: object'], 'f(*ht())', ['RB']),
    ('ht default+args', SRC, ['x: type[C] = C, *args: type[C]', '*args: object'], 'f(*ht())', ['RA']),
    ('ht args wrong type', SRC, ['*args: type[A]', '*args: object'], 'f(*ht())', ['RB']),
    ('ht args ok', SRC, ['*args: type[C]', '*args: object'], 'f(*ht())', ['RA']),
    ('ht two required', SRC, ['x: type[C], y: type[C]', '*args: type[C]', '*args: object'], 'f(*ht())', ['RB']),
    ('ht then C', SRC, ['x: type[C] = C, y: type[C] = C', '*args: type[C]', '*args: object'], 'f(*ht(), C)', ['RB']),
    ('ht then A args', SRC, ['*args: type[C]', '*args: object'], 'f(*ht(), A)', ['RB']),
    ('A then ht', SRC, ['x: type[A], *args: type[C]', '*args: object'], 'f(A, *ht())', ['RA']),
    ('A then ht fixed', SRC, ['x: type[A], y: type[C]', 'x: type[A], *args: type[C]', '*args: object'], 'f(A, *ht())', ['RB']),
    ('ht kwonly', SRC, ['*args: type[C], k: type[A]', '*args: object, k: object'], 'f(*ht(), k=A)', ['RA']),
    ('ht kw fills named', SRC, ['x: type[C] = C, *, k: type[A] = A', '*args: object, **kw: object'], 'f(*ht(), k=A)', ['RB']),
    # unknown length inside a literal
    ('lit of ht', SRC, ONE, 'f(*[*ht()])', ['RB']),
    ('lit A then ht', SRC, ['x: type[A], *args: type[C]', '*args: object'], 'f(*[A, *ht()])', ['RA']),
    ('lit A then ht fixed', SRC, ['x: type[A], y: type[C]', '*args: object'], 'f(*[A, *ht()])', ['RB']),
    ('lit ht then A', SRC, ['*args: type[C]', '*args: object'], 'f(*[*ht(), A])', ['RB']),
    ('tuple lit of ht', SRC, ['x: type[C]', 'x: type[C], y: type[C] = C', '*args: object'], 'f(*(*ht(),))', ['RC']),
    ('cond ht or one', SRC + R + 'xs = ht() if c else (C,)\n', ONE, 'f(*xs)', ['RB']),
    ('cond lit alt ht', SRC + R + 'xs = [C] if c else [*ht()]\n', ONE, 'f(*xs)', ['RB']),
    # fixed-length annotated tuples keep their length
    ('ft two', SRC, ['x: type[C], y: type[C]', '*args: object'], 'f(*ft())', ['RA']),
    ('ft1 one', SRC, ONE, 'f(*ft1())', ['RA']),
    ('mix positions', SRC, ['x: type[A], y: type[C]', '*args: object'], 'f(*mix())', ['RA']),
    ('mix swapped', SRC, ['x: type[C], y: type[A]', '*args: object'], 'f(*mix())', ['RB']),
    # literals and variables keep exact lengths
    ('lit one', '', ONE, 'f(*[C])', ['RA']),
    ('var one', 'xs = (C,)\n', ONE, 'f(*xs)', ['RA']),
    ('empty lit', '', ['x: type[C] = C', '*args: object'], 'f(*[])', ['RA']),
    # unknown source stays "not known"
    ('unknown', '', ['x: type[C]', '*args: object'], 'f(*undefined_name)', ['RA', 'RB']),
    # unknown length into containers
    ('ht in list arg', SRC, ['x: list[type[C]]', 'x: object'], 'f([*ht()])', ['RA']),
    ('ht in fixed tuple arg', SRC, ['x: tuple[type[C]]', 'x: tuple[type[C], ...]', 'x: object'], 'f((*ht(),))', ['RB']),
    # methods, classmethods and partial forward the same handling
    ('method ht', SRC + 'class K:\n    @overload\n    def m(self, x: type[C]) -> RA: ...\n    @overload\n    def m(self, *args: type[C]) -> RB: ...\n    def m(self, *args): ...\n', [], 'K().m(*ht())', ['RB']),
    ('classmethod ht', SRC + 'class K:\n    @overload\n    @classmethod\n    def m(cls, x: type[C]) -> RA: ...\n    @overload\n    @classmethod\n    def m(cls, *args: type[C]) -> RB: ...\n    @classmethod\n    def m(cls, *args): ...\n', [], 'K.m(*ht())', ['RB']),
    ('partial ht', SRC + 'import functools\n', ONE, 'functools.partial(f, *ht())()', ['RB']),
]
bad = 0
for name, pre, anns, call, exp in cases:
    got = infer(pre + (ov(*anns) if anns else '') + call)
    bad += got != exp
    print('%-4s %-22s got=%s exp=%s' % ('OK' if got == exp else 'BAD', name, got, exp))
print('bad', bad)
