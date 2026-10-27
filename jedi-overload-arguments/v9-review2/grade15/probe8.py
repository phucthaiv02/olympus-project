"""The defects of the independent review (2026-10-26) with variants."""
import sys
sys.path.insert(0, '/w')
from probe import ov, infer

R = 'c = random.random()\n'
L55 = '[' + ', '.join(str(i) for i in range(55)) + ']'
L10 = '[' + ', '.join(str(i) for i in range(10)) + ']'
cases = [
    # 1 the earlier key object stays, only the value is replaced
    ('1 int key stays', '', ['x: dict[bool, type[C]]', 'x: dict[int, type[C]]', 'x: object'], 'f({1: A, True: C})', ['RB']),
    ('1 bool key stays', '', ['x: dict[bool, type[C]]', 'x: dict[int, type[C]]', 'x: object'], 'f({True: A, 1: C})', ['RA']),
    ('1 key stays dstar', '', ['x: dict[bool, type[C]]', 'x: dict[int, type[C]]', 'x: object'], 'f({1: A, **{True: C}})', ['RB']),
    ('1 key stays var', 'd = {1: A, True: C}\n', ['x: dict[bool, type[C]]', 'x: dict[int, type[C]]', 'x: object'], 'f(d)', ['RB']),
    ('1 key maybe', R + 'e = {True: C} if c else {}\n', ['x: dict[bool, object]', 'x: dict[int, object]', 'x: object'], 'f({1: A, **e})', ['RB']),
    ('1 float key stays', '', ['x: dict[float, type[C]]', 'x: object'], 'f({1.0: A, 1: C})', ['RA']),
    # 2 * of a dict literal gives its keys
    ('2 star dict', '', ['x: str, y: str', 'x: object'], 'f(*{"a": A, "b": C})', ['RA']),
    ('2 star dict var', 'd = {"a": A, "b": C}\n', ['x: str, y: str', 'x: object'], 'f(*d)', ['RA']),
    ('2 star dict nested', '', ['x: str, y: str', 'x: object'], 'f(*[*{"a": A, "b": C}])', ['RA']),
    ('2 star dict dup', '', ['x: str', 'x: object, y: object'], 'f(*{"a": A, "a": C})', ['RA']),
    ('2 star dict comp dup', '', ['x: str', 'x: object, y: object'], 'f(*{k: C for k in ["a", "a"]})', ['RA']),
    ('2 star dict maybe', R + 'e = {"b": C} if c else {}\n', ['x: str', 'x: str, y: str = ""', 'x: object'], 'f(*{"a": A, **e})', ['RB']),
    ('2 dict in tuple ann', '', ['x: tuple[str, str]', 'x: object'], 'f((*{"a": A, "b": C},))', ['RA']),
    # 3 a str key that is not known may be any parameter name
    ('3 unknown key', 'k = str(1)\n', ['x: type[A] = A, **kw: object', 'x: object = None, **kw: object'], 'f(**{k: C})', ['RB']),
    ('3 ann dict', 'def g() -> dict[str, type[C]]: ...\n', ['x: type[A] = A, **kw: object', 'x: object = None, **kw: object'], 'f(**g())', ['RB']),
    ('3 ann dict ok', 'def g() -> dict[str, type[C]]: ...\n', ['x: type[C] = C, **kw: type[C]', 'x: object = None, **kw: object'], 'f(**g())', ['RA']),
    ('3 ann dict kwonly', 'def g() -> dict[str, type[C]]: ...\n', ['*, k: type[A] = A, **kw: object', '**kw: object'], 'f(**g())', ['RB']),
    ('3 ann dict posonly', 'def g() -> dict[str, type[C]]: ...\n', ['x: type[A] = A, /, **kw: object', '**kw: object'], 'f(**g())', ['RA']),
    ('3 ann dict given', 'def g() -> dict[str, type[C]]: ...\n', ['x: type[A], **kw: type[C]', '**kw: object'], 'f(A, **g())', ['RA']),
    # 4 a constrained TypeVar is one of its constraints for the whole call
    ('4 two params', "T = TypeVar('T', A, C)\n", ['x: T, y: T', 'x: object, y: object'], 'f(A(), C())', ['RB']),
    ('4 two params same', "T = TypeVar('T', A, C)\n", ['x: T, y: T', 'x: object, y: object'], 'f(C(), C())', ['RA']),
    ('4 two params sub', "T = TypeVar('T', A, C)\n", ['x: T, y: T', 'x: object, y: object'], 'f(A(), B())', ['RA']),
    ('4 list', "T = TypeVar('T', A, C)\n", ['x: list[T]', 'x: object'], 'f([A(), C()])', ['RB']),
    ('4 list type', "T = TypeVar('T', A, C)\n", ['x: list[type[T]]', 'x: object'], 'f([A, C])', ['RB']),
    ('4 type params', "T = TypeVar('T', A, C)\n", ['x: type[T], y: type[T]', 'x: object, y: object'], 'f(A, C)', ['RB']),
    ('4 args', "T = TypeVar('T', A, C)\n", ['*args: T', '*args: object'], 'f(A(), C())', ['RB']),
    ('4 bound two', "T = TypeVar('T', bound=object)\n", ['x: T, y: T', 'x: object, y: object'], 'f(A(), C())', ['RA']),
    ('4 single ok', "T = TypeVar('T', A, C)\n", ['x: T', 'x: object'], 'f(C())', ['RA']),
    ('4 union arg', "T = TypeVar('T', A, C)\n" + R, ['x: T', 'x: object'], 'f(A() if c else C())', ['RB']),
    # 5 int is accepted as float, int and float as complex
    ('5 int float', '', ['x: float', 'x: object'], 'f(1)', ['RA']),
    ('5 float complex', '', ['x: complex', 'x: object'], 'f(1.5)', ['RA']),
    ('5 int complex', '', ['x: complex', 'x: object'], 'f(1)', ['RA']),
    ('5 bool float', '', ['x: float', 'x: object'], 'f(True)', ['RA']),
    ('5 key float', '', ['x: dict[float, type[C]]', 'x: object'], 'f({1: C})', ['RA']),
    ('5 str float', '', ['x: float', 'x: object'], 'f("a")', ['RB']),
    ('5 float int', '', ['x: int', 'x: object'], 'f(1.5)', ['RB']),
    # 6 positional-only and keyword-only parameters
    ('6 posonly by name', '', ['x: type[A], /', 'x: object'], 'f(x=A)', ['RB']),
    ('6 kwonly by position', '', ['*, x: type[A]', 'x: object'], 'f(A)', ['RB']),
    ('6 posonly ok', '', ['x: type[A], /', 'x: object'], 'f(A)', ['RA']),
    ('6 kwonly ok', '', ['*, x: type[A]', 'x: object'], 'f(x=A)', ['RA']),
    ('6 posonly kw into kwargs', '', ['x: type[A] = A, /, **kw: type[C]', 'x: object = None, /, **kw: object'], 'f(x=C)', ['RA']),
    ('6 kwonly after args', '', ['*args: type[A], k: type[C]', '*args: object, k: object'], 'f(A, k=C)', ['RA']),
    # 7 a walrus argument is positional
    ('7 walrus', '', ['x: type[A]', 'x: object'], 'f(y := C)', ['RB']),
    ('7 walrus second', '', ['x: type[C], y: type[A] = A', 'x: object, y: object = None'], 'f(y := C)', ['RA']),
    # 8 large comprehensions
    ('8 dict comp 55', '', ['x: dict[int, int]', 'x: object'], 'f({a: a for a in %s})' % L55, ['RA']),
    ('8 dict comp 10x10', '', ['x: dict[int, type[C]]', 'x: object'], 'f({a: C for a in %s for b in %s})' % (L10, L10), ['RA']),
    ('8 list comp 55', '', ['x: list[int]', 'x: object'], 'f([a for a in %s])' % L55, ['RA']),
    ('8 dict comp not literal', '', ['x: dict[str, type[A]]', 'x: dict[str, type[C]]', 'x: object'], 'f({"k": v for v in [A] * 1 + [%s]})' % ', '.join(['C'] * 55), ['RC']),
]
bad = 0
for name, pre, anns, call, exp in cases:
    got = infer(pre + ov(*anns) + call)
    bad += got != exp
    print('%-4s %-26s got=%s exp=%s' % ('OK' if got == exp else 'BAD', name, got, exp))
print('bad', bad)
