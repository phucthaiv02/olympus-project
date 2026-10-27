"""Second independent review round (2026-10-26) with variants."""
import sys
sys.path.insert(0, '/w')
from probe import ov, infer

R = 'c = random.random()\n'
P = 'import functools\n'
cases = [
    # 1 key order when a possibly missing earlier key is replaced
    ('1 order maybe', R, ['x: int, y: str', 'x: object, y: object'], 'f(*{**({1: A} if c else {}), "b": B, 1: C})', ['RB']),
    ('1 order maybe tuple', R, ['x: tuple[int, str]', 'x: object'], 'f((*{**({1: A} if c else {}), "b": B, 1: C},))', ['RB']),
    ('1 order sure', '', ['x: int, y: str', 'x: object, y: object'], 'f(*{1: A, "b": B, 1: C})', ['RA']),
    ('1 order sure dstar', '', ['x: int, y: str', 'x: object, y: object'], 'f(*{**{1: A}, "b": B, 1: C})', ['RA']),
    ('1 values ok', R, ['x: dict[object, type[C]]', 'x: object'], 'f({**({1: A} if c else {}), "b": C, 1: C})', ['RA']),
    # 2 keys of a comprehension whose values are not known
    ('2 comp keys', '', ['x: str, y: str', 'x: object'], 'f(*{k: undefined_v for k in ["a", "b"]})', ['RA']),
    ('2 comp keys nested', '', ['x: str, y: str', 'x: object'], 'f(*[*{k: undefined_v for k in ["a", "b"]}])', ['RA']),
    ('2 comp keys dup', '', ['x: str', 'x: object, y: object'], 'f(*{k: undefined_v for k in ["a", "a"]})', ['RA']),
    ('2 comp values unknown', '', ['x: dict[str, type[C]]', 'x: object'], 'f({k: undefined_v for k in ["a"]})', ['RB']),
    # 3 a required parameter without an argument
    ('3 posonly missing', '', ['x, /, y: type[A] = A', 'y: object = None'], 'f(y=A)', ['RB']),
    ('3 unannotated missing', '', ['x, y: type[A] = A', 'y: object = None'], 'f(y=A)', ['RB']),
    ('3 partial missing', P, ['x: type[A], y, z: object = None', '*a: object, **k: object'], 'p = functools.partial(f, A)\np(z=1)', ['RB']),
    ('3 kwonly missing', '', ['*, x, y: type[A] = A', 'y: object = None'], 'f(y=A)', ['RB']),
    ('3 ok control', '', ['x, y: type[A] = A', 'y: object = None'], 'f(1, y=A)', ['RA']),
    # 4 functools.partial passes f(*partial_args, *args, **{**partial_kw, **kw})
    ('4 partial kw then pos', P, ['y: type[C] = C, x: type[A] = A', '*a: object, **k: object'], 'p = functools.partial(f, y=C)\np(A)', ['RB']),
    ('4 partial pos then kw', P, ['y: object = None, x: type[A] = A', '*a: object, **k: object'], 'p = functools.partial(f, A)\np(y=1)', ['RB']),
    ('4 partial ok', P, ['y: type[C] = C, x: type[A] = A', '*a: object, **k: object'], 'p = functools.partial(f, C)\np(x=A)', ['RA']),
    ('4 partial kw replaced', P, ['x: type[C]', 'x: object'], 'p = functools.partial(f, x=A)\np(x=C)', ['RA']),
    ('4 partial pos order', P, ['x: type[A], y: type[C]', 'x: object, y: object'], 'p = functools.partial(f, A)\np(C)', ['RA']),
    # 5 quoted TypeVar annotations
    ('5 quoted T', "T = TypeVar('T', A, C)\n", ['x: "T", y: "T"', 'x: object, y: object'], 'f(A(), C())', ['RB']),
    ('5 quoted T args', "T = TypeVar('T', A, C)\n", ['*a: "T"', '*a: object'], 'f(A(), C())', ['RB']),
    ('5 quoted T ok', "T = TypeVar('T', A, C)\n", ['x: "T", y: "T"', 'x: object, y: object'], 'f(C(), C())', ['RA']),
    ('5 quoted list T', "T = TypeVar('T', A, C)\n", ['x: "list[T]"', 'x: object'], 'f([A(), C()])', ['RB']),
    # many constraints stay fast and exact
    ('many constraints', ''.join("T%d = TypeVar('T%d', %s)\n" % (i, i, ', '.join(['A', 'C', 'int', 'str', 'bytes', 'float', 'list', 'dict', 'set', 'tuple'])) for i in range(6)),
     ['a: T0, b: T1, c: T2, d: T3, e: T4, g: T5', '*a: object'], 'f(A(), C(), 1, "s", b"", 1.5)', ['RA']),
    ('many constraints no', ''.join("T%d = TypeVar('T%d', %s)\n" % (i, i, ', '.join(['A', 'C', 'int', 'str', 'bytes', 'float', 'list', 'dict', 'set', 'tuple'])) for i in range(6)),
     ['a: T0, b: T0, c: T2, d: T3, e: T4, g: T5', '*a: object'], 'f(A(), C(), 1, "s", b"", 1.5)', ['RB']),
]
bad = 0
for name, pre, anns, call, exp in cases:
    got = infer(pre + ov(*anns) + call)
    bad += got != exp
    print('%-4s %-24s got=%s exp=%s' % ('OK' if got == exp else 'BAD', name, got, exp))
print('bad', bad)
