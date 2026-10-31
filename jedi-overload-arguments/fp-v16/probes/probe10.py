"""Auto Review v14: call-wide constrained TypeVar, ** into named parameters."""
import sys
sys.path.insert(0, '/w')
from probe import ov, infer

T = "T = TypeVar('T', A, int)\n"
cases = [
    ('1 type two diff', T, ['x: type[T], y: type[T]', 'x: object, y: object'], 'f(A, int)', ['RB']),
    ('1 type two diff rev', T, ['x: type[T], y: type[T]', 'x: object, y: object'], 'f(int, B)', ['RB']),
    ('1 type two same', T, ['x: type[T], y: type[T]', 'x: object, y: object'], 'f(A, B)', ['RA']),
    ('1 type two int', T, ['x: type[T], y: type[T]', 'x: object, y: object'], 'f(int, int)', ['RA']),
    ('1 inst two diff', T, ['x: T, y: T', 'x: object, y: object'], 'f(A(), 1)', ['RB']),
    ('1 kw two diff', T, ['x: type[T], y: type[T]', 'x: object, y: object'], 'f(y=int, x=A)', ['RB']),
    ('1 args diff', T, ['*args: type[T]', '*args: object'], 'f(A, int)', ['RB']),
    ('1 args same', T, ['*args: type[T]', '*args: object'], 'f(A, B)', ['RA']),
    ('1 list and param', T, ['x: list[type[T]], y: type[T]', 'x: object, y: object'], 'f([A], int)', ['RB']),
    ('2 dstar named C', '', ['x: type[A]', 'x: type[C]', 'x: object'], "f(**{'x': C})", ['RB']),
    ('2 dstar named A', '', ['x: type[A]', 'x: type[C]', 'x: object'], "f(**{'x': A})", ['RA']),
    ('2 dstar named var', "d = {'x': C}\n", ['x: type[A]', 'x: type[C]', 'x: object'], 'f(**d)', ['RB']),
    ('2 dstar named wrong', '', ['x: type[A]', 'x: type[C]', 'x: object'], "f(**{'x': 1})", ['RC']),
    ('2 dstar two named', '', ['x: type[A], y: type[C]', 'x: object, y: object'], "f(**{'y': C, 'x': A})", ['RA']),
    ('2 dstar two swapped', '', ['x: type[A], y: type[C]', 'x: object, y: object'], "f(**{'y': A, 'x': C})", ['RB']),
    ('2 pos and dstar', '', ['x: type[A], y: type[C]', 'x: object, y: object'], "f(A, **{'y': C})", ['RA']),
    ('2 dstar default', '', ['x: type[A] = A, y: type[C] = C', 'x: object = None, y: object = None'], "f(**{'y': A})", ['RB']),
    ('2 dstar kwonly', '', ['*, x: type[C]', '*, x: object'], "f(**{'x': C})", ['RA']),
    ('2 dstar nested', '', ['x: type[A]', 'x: type[C]', 'x: object'], "f(**{**{'x': A}, 'x': C})", ['RB']),
]
bad = 0
for name, pre, anns, call, exp in cases:
    got = infer(pre + ov(*anns) + call)
    bad += got != exp
    print('%-4s %-22s got=%s exp=%s' % ('OK' if got == exp else 'BAD', name, got, exp))
print('bad', bad)
