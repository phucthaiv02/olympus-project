"""Candidate assertions, each derived from a sentence of the v10 description.
Prints one line per case: name<TAB>result. Expected values are judged by hand
against the reference column."""
import sys
sys.path.insert(0, '/w')
from probe import HEADER, ov, infer

R = 'c = random.random()\n'
D = ['x: dict[str, type[A]]', 'x: dict[str, type[C]]', 'x: object']
L = ['x: list[type[A]]', 'x: list[type[C]]', 'x: object']
cases = [
    # type[X], subclasses, instances
    ('subclass of subclass', 'class D(B): pass\n', ['x: type[A]', 'x: object'], 'f(D)'),
    ('instance to type', '', ['x: type[A]', 'x: object'], 'f(A())'),
    ('class to instance ann', '', ['x: A', 'x: object'], 'f(A)'),
    ('type of instance', '', ['x: type[A]', 'x: object'], 'f(type(A()))'),
    ('Type alias', 'from typing import Type\n', ['x: Type[A]', 'x: object'], 'f(C)'),
    ('bare type instance', '', ['x: type', 'x: object'], 'f(A())'),
    ('type[Any] instance', '', ['x: type[Any]', 'x: object'], 'f(1)'),
    ('Union of type', '', ['x: Union[type[A], type[C]]', 'x: object'], 'f(C)'),
    ('type of Union', '', ['x: type[Union[A, C]]', 'x: object'], 'f(C)'),
    ('type pipe', '', ['x: type[A] | type[C]', 'x: object'], 'f(C)'),
    ('pipe None', '', ['x: type[A] | None', 'x: object'], 'f(None)'),
    # TypeVar
    ('tv bound sub', "T = TypeVar('T', bound=A)\n", ['x: type[T]', 'x: object'], 'f(B)'),
    ('tv bound no', "T = TypeVar('T', bound=A)\n", ['x: type[T]', 'x: object'], 'f(C)'),
    ('tv constraints', "T = TypeVar('T', A, C)\n", ['x: type[T]', 'x: object'], 'f(B)'),
    ('tv constraints no', "T = TypeVar('T', A, B)\n", ['x: type[T]', 'x: object'], 'f(C)'),
    ('tv plain instance', "T = TypeVar('T')\n", ['x: type[T]', 'x: object'], 'f(A())'),
    ('tv bound union', "T = TypeVar('T', bound=Union[A, C])\n", ['x: type[T]', 'x: object'], 'f(C)'),
    ('tv bound instance ann', "T = TypeVar('T', bound=A)\n", ['x: T', 'x: object'], 'f(C())'),
    ('tv in list', "T = TypeVar('T', bound=A)\n", ['x: list[type[T]]', 'x: object'], 'f([B, C])'),
    # Optional
    ('Optional None', 'from typing import Optional\n', ['x: Optional[type[A]]', 'x: object'], 'f(None)'),
    ('Optional default', 'from typing import Optional\n', ['x: Optional[type[A]] = None', 'x: object'], 'f()'),
    ('Optional in list', 'from typing import Optional\n', ['x: list[Optional[type[A]]]', 'x: object'], 'f([None, B])'),
    ('Optional cond', 'from typing import Optional\n' + R, ['x: Optional[type[A]]', 'x: object'], 'f(A if c else None)'),
    # containers
    ('nested list', '', ['x: list[list[type[A]]]', 'x: object'], 'f([[A], [C]])'),
    ('dict of list', '', ['x: dict[str, list[type[A]]]', 'x: object'], 'f({"k": [A, C]})'),
    ('tuple pos', '', ['x: tuple[type[A], type[C]]', 'x: object'], 'f((C, A))'),
    ('tuple pos ok', '', ['x: tuple[type[A], type[C]]', 'x: object'], 'f((B, C))'),
    ('tuple len', '', ['x: tuple[type[A], type[C]]', 'x: object'], 'f((A,))'),
    ('tuple ellipsis star', R + 'xs = [A] if c else []\n', ['x: tuple[type[A], ...]', 'x: object'], 'f((A, *xs))'),
    ('set literal', '', ['x: set[type[A]]', 'x: object'], 'f({A, C})'),
    ('set star', '', ['x: set[type[A]]', 'x: object'], 'f({A, *[C]})'),
    ('frozen list in tuple', '', ['x: tuple[list[type[A]], ...]', 'x: object'], 'f(([A], [C]))'),
    ('list var', 'xs = [A, C]\n', L, 'f(xs)'),
    ('list comp', '', L, 'f([C for _ in [1, 2]])'),
    ('list comp cond val', R, L, 'f([A if c else C for _ in [1]])'),
    ('dict keys int', '', ['x: dict[int, type[A]]', 'x: object'], 'f({"a": A})'),
    ('empty list', '', L, 'f([])'),
    ('empty dict', '', D, 'f({})'),
    ('list of instance', '', L, 'f([A()])'),
    # unions / conditionals
    ('cond arg', R, ['x: type[A]', 'x: object'], 'f(A if c else C)'),
    ('cond arg ok', R, ['x: type[A]', 'x: object'], 'f(A if c else B)'),
    ('cond var', R + 'v = A if c else C\n', ['x: type[A]', 'x: type[C]', 'x: object'], 'f(v)'),
    ('cond list', R + 'v = [A] if c else [C]\n', L, 'f(v)'),
    ('cond dict', R + 'v = {"k": A} if c else {"k": C}\n', D, 'f(v)'),
    ('union ann var', 'def g() -> Union[type[A], type[C]]: ...\n', ['x: type[A]', 'x: object'], 'f(g())'),
    ('if branches var', R + 'if c:\n    v = A\nelse:\n    v = C\n', ['x: type[A]', 'x: object'], 'f(v)'),
    # *args / **kwargs
    ('args mixed', '', ['*args: type[A]', '*args: object'], 'f(A, C)'),
    ('args after pos', '', ['x: type[C], *args: type[A]', 'x: object, *args: object'], 'f(C, A, C)'),
    ('kwargs mixed', '', ['**kw: type[A]', '**kw: object'], 'f(p=A, q=C)'),
    ('kwargs with named', '', ['x: type[C], **kw: type[A]', 'x: object, **kw: object'], 'f(x=C, y=A)'),
    ('kwargs star dict', '', ['**kw: type[A]', '**kw: object'], 'f(**{"p": A, "q": C})'),
    ('kwargs star dict ok', '', ['**kw: type[A]', '**kw: object'], 'f(**{"p": A, "q": B})'),
    ('kw from star into named', '', ['x: type[C]', 'x: object'], 'f(**{"x": C})'),
    ('kw from star into named bad', '', ['x: type[C]', 'x: object'], 'f(**{"x": A})'),
    ('kw unknown key', '', ['x: type[C]', 'x: object = None, **kw: object'], 'f(**{"y": C})'),
    ('star tuple var', 'xs = (A, C)\n', ['x: type[A], y: type[C]', 'x: object, y: object'], 'f(*xs)'),
    ('star tuple var swap', 'xs = (C, A)\n', ['x: type[A], y: type[C]', 'x: object, y: object'], 'f(*xs)'),
    ('star typed tuple', 'def g() -> tuple[type[A], type[C]]: ...\n', ['x: type[A], y: type[C]', 'x: object, y: object'], 'f(*g())'),
    ('star too many', '', ['x: type[A]', '*args: object'], 'f(*[A, A])'),
    ('star to kwonly', '', ['x: type[A], *, k: type[C]', 'x: object, *, k: object'], 'f(*[A], k=C)'),
    ('star after kw', '', ['x: type[A], *, k: type[C]', 'x: object, *, k: object'], 'f(k=C, *[A])'),
    ('nested star 3', '', ['*args: type[A]', '*args: object'], 'f(*[*[*[A, C]]])'),
    ('star in tuple literal arg', '', ['x: tuple[type[A], ...]', 'x: object'], 'f((*[A], *(C,)))'),
    ('dstar nested', '', ['**kw: type[A]', '**kw: object'], 'f(**{**{"p": C}, "p": A})'),
    ('dstar nested 2', '', ['**kw: type[A]', '**kw: object'], 'f(**{"p": A, **{"p": C}})'),
    # keys and replacement
    ('key replace', '', D, 'f({"k": C, "k": A})'),
    ('key replace dstar', '', D, 'f({"k": C, **{"k": A}})'),
    ('key maybe dstar', R + 'd = {"k": A} if c else {}\n', D, 'f({"k": C, **d})'),
    ('key maybe then sure', R + 'd = {"k": A} if c else {}\n', D, 'f({**d, "k": C})'),
    ('key cond both', R + 'k = "a" if c else "b"\n', D, 'f({"a": C, k: A})'),
    ('key int float', '', ['x: dict[float, type[A]]', 'x: dict[float, type[C]]', 'x: object'], 'f({1: A, 1.0: C})'),
    ('key false zero', '', ['x: dict[int, type[A]]', 'x: dict[int, type[C]]', 'x: object'], 'f({0: A, False: C})'),
    ('key None', '', ['x: dict[None, type[A]]', 'x: dict[None, type[C]]', 'x: object'], 'f({None: A, None: C})'),
    ('key class', '', ['x: dict[type, type[A]]', 'x: dict[type, type[C]]', 'x: object'], 'f({C: A, C: C})'),
    ('key tuple', '', ['x: dict[tuple, type[A]]', 'x: dict[tuple, type[C]]', 'x: object'], 'f({(1, 2): A, (1, 2): C})'),
    ('key str concat', '', D, 'f({"ab": A, "a" "b": C})'),
    ('key var', 'K = "k"\n', D, 'f({K: A, "k": C})'),
    ('key bytes', '', ['x: dict[bytes, type[A]]', 'x: dict[bytes, type[C]]', 'x: object'], 'f({b"k": A, b"k": C})'),
    ('key str vs bytes', '', ['x: dict[object, type[A]]', 'x: dict[object, type[C]]', 'x: object'], 'f({b"k": C, "k": A})'),
    ('dict var then dstar', 'd = {"k": C}\n', D, 'f({**d, "k": A})'),
    ('nested maybe 2', R + 'd1 = {"k": A} if c else {}\nd2 = {**d1}\n', D, 'f({"k": C, **d2})'),
    ('maybe two alts', R + 'd = {"k": A} if c else {"k": C}\n', D, 'f({"k": C, **d})'),
    ('maybe empty both', R + 'd = {} if c else {}\n', D, 'f({"k": C, **d})'),
    # comprehensions
    ('empty comp list', '', L, 'f([C for _ in []])'),
    ('empty comp set', '', ['x: set[type[A]]', 'x: object'], 'f({C for _ in ()})'),
    ('empty comp var', 'e = []\n', D, 'f({"k": C for _ in e})'),
    ('comp if false', '', D, 'f({"k": C for _ in [1] if 0})'),
    ('empty comp dstar kw', '', ['**kw: type[A]', '**kw: object'], 'f(**{"k": C for _ in []})'),
    ('empty comp star', '', ['*args: type[A]', '*args: object'], 'f(*[C for _ in []])'),
    ('empty comp star pos', '', ['x: type[A] = A', 'x: object = None'], 'f(*[C for _ in []])'),
    # sequence lengths
    ('len shift kw', R + 'xs = [A] if c else []\n', ['x: type[A] = A, y: type[C] = C', 'x: object = None, y: object = None'], 'f(*xs, y=C)'),
    ('len two vars', R + 'xs = [A] if c else [A, A]\n', ['x: type[A], y: type[A] = A, z: type[C] = C', 'x: object, y: object = None, z: object = None'], 'f(*xs, C)'),
    ('len tuple var', R + 'xs = (A,) if c else ()\n', ['x: type[A] = A, y: type[C] = C', 'x: object = None, y: object = None'], 'f(*xs, C)'),
    ('len missing required', R + 'xs = [A] if c else []\n', ['x: type[A]', 'x: object = None'], 'f(*xs)'),
    ('len star args', R + 'xs = [A] if c else []\n', ['*args: type[A]', '*args: object'], 'f(*xs, C)'),
    ('len nested', R + 'xs = [A] if c else []\nys = [*xs]\n', ['x: type[A] = A, y: type[C] = C', 'x: object = None, y: object = None'], 'f(*ys, C)'),
    ('len list arg', R + 'xs = [A] if c else []\n', ['x: tuple[type[A], type[C]]', 'x: object'], 'f((*xs, C))'),
    ('len dict maybe kwarg', R + 'd = {"x": A} if c else {}\n', ['x: type[A] = A, y: type[C] = C', 'x: object = None, y: object = None'], 'f(**d, y=C)'),
    ('len dict maybe required', R + 'd = {"x": A} if c else {}\n', ['x: type[A]', 'x: object = None'], 'f(**d)'),
    # methods
    ('method overload', 'class K:\n    @overload\n    def m(self, x: type[A]) -> RA: ...\n    @overload\n    def m(self, x: object) -> RB: ...\n    def m(self, x): ...\n', [], 'K().m(C)'),
    ('classmethod overload', 'class K:\n    @overload\n    @classmethod\n    def m(cls, x: type[A]) -> RA: ...\n    @overload\n    @classmethod\n    def m(cls, x: object) -> RB: ...\n    @classmethod\n    def m(cls, x): ...\n', [], 'K.m(C)'),
    ('staticmethod overload', 'class K:\n    @overload\n    @staticmethod\n    def m(x: type[A]) -> RA: ...\n    @overload\n    @staticmethod\n    def m(x: object) -> RB: ...\n    @staticmethod\n    def m(x): ...\n', [], 'K.m(C)'),
    ('init overload', 'class K:\n    @overload\n    def __init__(self, x: type[A]) -> None: ...\n    @overload\n    def __init__(self, x: object) -> None: ...\n    def __init__(self, x): ...\n', [], 'K(C)'),
]

for name, pre, anns, call in cases:
    code = pre + (ov(*anns) if anns else '') + call
    print('%s\t%s' % (name, infer(code)))
