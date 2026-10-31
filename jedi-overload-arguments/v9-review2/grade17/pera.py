# Per-assertion results for the hard tests plus candidate discriminators.
import importlib.util, json, sys, jedi
spec = importlib.util.spec_from_file_location('t', sys.argv[1]); m = importlib.util.module_from_spec(spec); spec.loader.exec_module(m)
S = jedi.Script; res = {}
def rec(name, code, expected):
    try: got = m.infer_result(S, code)
    except Exception as e: got = 'CRASH ' + type(e).__name__
    res[name] = got == expected
n = {'i': 0}; label = {'v': ''}
def fake_check(Script, annotations, call, expected, prefix=''):
    n['i'] += 1; rec(f"{label['v']}{n['i']} {call}", prefix + m.overloads(*annotations) + call, expected)
m.check = fake_check
for lab, fn in [('F', m.test_type_var_with_forward_references), ('K', m.test_keys_that_are_not_strings), ('I', m.test_items_of_mappings_and_other_containers)]:
    label['v'] = lab; n['i'] = 0; fn(S)
p = 'T = TypeVar("T", A, int)\n'
c = p + m.signatures('x: type[T], y: type[T]', 'x: object, y: object')
for i, (call, e) in enumerate([('f(A, int)', ['RB']), ('f(int, B)', ['RB']), ('f(y=int, x=A)', ['RB']), ('f(A, B)', ['RA']), ('f(int, int)', ['RA'])], 1):
    rec(f'C{i} {call}', c + call, e)
c = p + m.signatures('x: T, y: T', 'x: object, y: object'); rec('C6 T,T f(A(), 1)', c + 'f(A(), 1)', ['RB'])
c = p + m.variadic('type[T]', 'object', '*x: %s'); rec('C7 *x f(A, int)', c + 'f(A, int)', ['RB']); rec('C8 *x f(A, B)', c + 'f(A, B)', ['RA'])
c = p + m.signatures('x: list[type[T]], y: type[T]', 'x: object, y: object'); rec('C9 list f([A], int)', c + 'f([A], int)', ['RB'])
H = lambda *sig: ''.join('@overload\ndef f(%s) -> %s: ...\n' % (s, r) for s, r in zip(sig, ('RA', 'RB'))) + 'def f(*a, **k): ...\n'
for name, code, e in [
    ('G1 kwonly', H('*, x: int', 'x: int') + 'f(1)', ['RB']),
    ('G2 posonly', H('x: int, /', '*, x: int') + 'f(x=1)', ['RB']),
    ('G3 int->float', H('x: float', 'x: str') + 'f(1)', ['RA']),
    ('G4 float->complex', H('x: complex', 'x: str') + 'f(1.5)', ['RA']),
    ('G5 None', H('x: None', 'x: int') + 'f(None)', ['RA']),
    ('G6 tuple order', H('x: tuple[int, str]', 'x: object') + 'f(("s", 1))', ['RB']),
    ('G7 key kept bool', H('x: dict[bool, C]', 'x: dict[int, C]') + 'f({1: C(), True: C()})', ['RB']),
    ('G8 key kept int', H('x: dict[bool, C]', 'x: dict[int, C]') + 'f({True: C(), 1: C()})', ['RA']),
]:
    rec(name, code, e)
print('JSON' + json.dumps(res))
