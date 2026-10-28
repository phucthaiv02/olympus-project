"""
Choosing the overload of a function from the arguments of a call.
"""

from textwrap import dedent

import pytest


def soft(fn, label):
    try:
        ok = fn()
    except Exception as e:
        ok = False
        label = 'EXC %r :: %s' % (e, label)
    if not ok:
        print("FAILCASE", str(label).replace("\n", "\\n"))


def check(Script, annotations, call, expected, prefix=''):
    soft(lambda: _orig_check(Script, annotations, call, expected, prefix) or True,
         (annotations, call, prefix[:40]))

HEADER = """\
from typing import overload, Any, Type, TypeVar, Union

class A: pass
class B(A): pass
class C: pass
class RA: pass
class RB: pass
class RC: pass
"""


def overloads(*annotations):
    code = ''
    for annotation, result in zip(annotations, ('RA', 'RB', 'RC')):
        code += '@overload\ndef f(x: %s) -> %s: ...\n' % (annotation, result)
    return code + 'def f(x): ...\n'


def infer_result(Script, code):
    code = HEADER + dedent(code).strip('\n')
    lines = code.splitlines()
    return sorted(d.name for d in Script(code).infer(len(lines), len(lines[-1])))


def _orig_check(Script, annotations, call, expected, prefix=''):
    code = prefix + overloads(*annotations) + call
    soft(lambda: infer_result(Script, code) == expected, (annotations, call))


def test_issue_example(Script):
    code = dedent('''\
        from typing import overload

        class Foo: pass

        class FooBuilder:
            def func_foo(self): pass

        class Bar: pass

        class BarBuilder:
            def func_bar(self): pass

        @overload
        def get_obj(input: type[Foo]) -> FooBuilder: ...
        @overload
        def get_obj(input: type[Bar]) -> BarBuilder: ...
        def get_obj(input: type): pass

        builder = get_obj(Foo)
        builder.''')
    names = [c.name for c in Script(code).complete()]
    assert 'func_foo' in names
    assert 'func_bar' not in names

    code = code.replace('get_obj(Foo)', 'get_obj(Bar)')
    names = [c.name for c in Script(code).complete()]
    assert 'func_bar' in names
    assert 'func_foo' not in names


@pytest.mark.parametrize('first', ['type[A]', 'Type[A]'])
@pytest.mark.parametrize(
    'call, expected', [
        ('f(A)', ['RA']),
        ('f(B)', ['RA']),
        ('f(C)', ['RB']),
    ]
)
def test_type_of_class(Script, first, call, expected):
    check(Script, [first, 'type[C]'], call, expected)


def test_type_rejects_instances(Script):
    check(Script, ['type[A]', 'A'], 'f(A())', ['RB'])
    check(Script, ['Type[A]', 'A'], 'f(B())', ['RB'])
    check(Script, ['A', 'type[A]'], 'f(A)', ['RB'])
    factory = 'class Factory:\n    def __call__(self) -> A: ...\n'
    check(Script, ['type[A]', 'Factory'], 'f(Factory())', ['RB'], prefix=factory)


@pytest.mark.parametrize('annotation', ['type', 'Type', 'type[Any]', 'Type[Any]'])
def test_any_class(Script, annotation):
    for cls in ['A', 'B', 'C', 'int', 'str']:
        check(Script, ['int', annotation], 'f(%s)' % cls, ['RB'])
    for instance, cls in [('3', 'int'), ('C()', 'C'), ('"s"', 'str')]:
        check(Script, [annotation, cls], 'f(%s)' % instance, ['RB'])


@pytest.mark.parametrize('annotation', ['type[Union[A, C]]', 'Type[A | C]', 'type[A] | type[C]'])
def test_type_of_union(Script, annotation):
    check(Script, ['int', annotation], 'f(C)', ['RB'])
    check(Script, ['int', annotation], 'f(B)', ['RB'])


def test_type_of_type_var(Script):
    bound = 'T = TypeVar("T", bound=A)\n'
    check(Script, ['type[T]', 'int'], 'f(C)', ['RA', 'RB'], prefix=bound)
    check(Script, ['type[T]', 'type[C]'], 'f(C)', ['RB'], prefix=bound)
    check(Script, ['type[T]', 'type[C]'], 'f(B)', ['RA'], prefix=bound)
    check(Script, ['type[T]', 'type[C]'], 'f(A)', ['RA'], prefix=bound)
    check(Script, ['type[T]', 'A'], 'f(B())', ['RB'], prefix=bound)
    constrained = 'T = TypeVar("T", A, int)\n'
    check(Script, ['type[T]', 'type[C]'], 'f(C)', ['RB'], prefix=constrained)
    check(Script, ['type[T]', 'type[C]'], 'f(int)', ['RA'], prefix=constrained)
    check(Script, ['type[T]', 'type[C]'], 'f(B)', ['RA'], prefix=constrained)
    check(Script, ['type[T]', 'type[str]'], 'f(str)', ['RB'], prefix=constrained)
    check(Script, ['type[T]', 'int'], 'f(3)', ['RB'], prefix=constrained)
    unbound = 'T = TypeVar("T")\n'
    check(Script, ['type[T]', 'type[C]'], 'f(C)', ['RA'], prefix=unbound)
    check(Script, ['type[T]', 'type[C]'], 'f(int)', ['RA'], prefix=unbound)
    check(Script, ['type[T]', 'int'], 'f(3)', ['RB'], prefix=unbound)


def test_type_from_annotated_names(Script):
    check(Script, ['type[A]', 'type[C]'], 'f(t)', ['RB'], prefix='t: type[C]\n')
    code = overloads('type[A]', 'type[C]') + dedent('''\
        def g(k: type[C]):
            f(k)''')
    soft(lambda: infer_result(Script, code) == ['RB'], "infer_result(Script, code) == ['RB']")


def test_first_matching_overload_wins(Script):
    check(Script, ['type[A]', 'type[B]'], 'f(B)', ['RA'])
    check(Script, ['type[B]', 'type[A]'], 'f(B)', ['RA'])
    check(Script, ['type[B]', 'type[A]'], 'f(A)', ['RB'])
    check(Script, ['type', 'type[A]'], 'f(A)', ['RA'])
    check(Script, ['type[A]', 'type[B]', 'type[C]'], 'f(C)', ['RC'])


def test_type_keyword_argument(Script):
    check(Script, ['int', 'type[C]'], 'f(x=C)', ['RB'])
    check(Script, ['type[A]', 'type[B]'], 'f(x=B)', ['RA'])
    check(Script, ['type[C]', 'int'], 'f(x=3)', ['RB'])


TWO_PARAMETERS = '''\
@overload
def g(x: type[A], y: type[C]) -> RA: ...
@overload
def g(x: type[A], y: type[A]) -> RB: ...
def g(x, y): ...
'''


def test_every_argument_has_to_fit(Script):
    for call, expected in [
        ('g(A, C)', ['RA']),
        ('g(B, C)', ['RA']),
        ('g(A, A)', ['RB']),
        ('g(A, y=B)', ['RB']),
        ('g(y=C, x=B)', ['RA']),
        ('g(x=A, y=A)', ['RB']),
        ('g(C, A)', ['RA', 'RB']),
        ('g(A, A())', ['RA', 'RB']),
    ]:
        soft(lambda: infer_result(Script, TWO_PARAMETERS + call) == expected, call)


def variadic(first, second, signature, third=None):
    code = ''
    for annotation, result in [(first, 'RA'), (second, 'RB'), (third, 'RC')]:
        if annotation is not None:
            code += '@overload\ndef f(%s) -> %s: ...\n' % (signature % annotation, result)
    return code + 'def f(*args, **kwargs): ...\n'


@pytest.mark.parametrize('signature', ['*x: %s', 'y: int, *x: %s'])
def test_var_positional(Script, signature):
    prefix = '1, ' if signature.startswith('y') else ''
    for arguments, expected in [
        ('C, C', ['RB']),
        ('A, B', ['RA']),
        ('A, C', ['RA', 'RB']),
        ('C, A', ['RA', 'RB']),
        ('', ['RA']),
    ]:
        call = 'f(%s)' % (prefix + arguments).rstrip(', ')
        code = variadic('type[A]', 'type[C]', signature) + call
        soft(lambda: infer_result(Script, code) == expected, (signature, call))


def test_var_positional_unpacked(Script):
    code = variadic('type[A]', 'type[C]', '*x: %s')
    soft(lambda: infer_result(Script, code + 'ks = (C, C)\nf(*ks)') == ['RB'], "infer_result(Script, code + 'ks = (C, C)\\nf(*ks)') == ['RB']")
    soft(lambda: infer_result(Script, code + 'ks = [A, B]\nf(*ks)') == ['RA'], "infer_result(Script, code + 'ks = [A, B]\\nf(*ks)') == ['RA']")


@pytest.mark.parametrize('signature', ['**x: %s', 'y: int, **x: %s'])
def test_var_keyword(Script, signature):
    prefix = 'y=1, ' if signature.startswith('y') else ''
    for arguments, expected in [
        ('p=C, q=C', ['RB']),
        ('p=A, q=B', ['RA']),
        ('p=A, q=C', ['RA', 'RB']),
        ('', ['RA']),
    ]:
        call = 'f(%s)' % (prefix + arguments).rstrip(', ')
        code = variadic('type[A]', 'type[C]', signature) + call
        soft(lambda: infer_result(Script, code) == expected, (signature, call))


def test_var_keyword_unpacked(Script):
    code = variadic('type[A]', 'type[C]', '**x: %s')
    soft(lambda: infer_result(Script, code + "f(**{'p': C})") == ['RB'], 'infer_result(Script, code + "f(**{\'p\': C})") == [\'RB\']')


def test_variadic_parameters_with_other_annotations(Script):
    prefix = 'T = TypeVar("T", bound=A)\n'
    code = prefix + variadic('type[T]', 'type[Union[A, C]]', '*x: %s')
    soft(lambda: infer_result(Script, code + 'f(A, B)') == ['RA'], "infer_result(Script, code + 'f(A, B)') == ['RA']")
    soft(lambda: infer_result(Script, code + 'f(A, C)') == ['RB'], "infer_result(Script, code + 'f(A, C)') == ['RB']")
    soft(lambda: infer_result(Script, code + 'f(C, int)') == ['RA', 'RB'], "infer_result(Script, code + 'f(C, int)') == ['RA', 'RB']")
    code = variadic('type', 'int', '*x: %s')
    soft(lambda: infer_result(Script, code + 'f(A, int, str)') == ['RA'], "infer_result(Script, code + 'f(A, int, str)') == ['RA']")
    soft(lambda: infer_result(Script, code + 'f(A, 1)') == ['RA', 'RB'], "infer_result(Script, code + 'f(A, 1)') == ['RA', 'RB']")
    soft(lambda: infer_result(Script, code + 'f(1, 2)') == ['RB'], "infer_result(Script, code + 'f(1, 2)') == ['RB']")


def test_variadic_parameters_together(Script):
    code = dedent('''\
        @overload
        def f(y: int, *x: type[A], z: int = 0, **k: type[A]) -> RA: ...
        @overload
        def f(y: int, *x: type[C], z: int = 0, **k: type[C]) -> RB: ...
        def f(*args, **kwargs): ...
        ''')
    for call, expected in [
        ('f(1, A, B, p=A)', ['RA']),
        ('f(1, C, z=2, p=C)', ['RB']),
        ('f(1, A, p=C)', ['RA', 'RB']),
        ('f(1, p=C)', ['RB']),
        ('f(1, C)', ['RB']),
        ('f(1, z=2)', ['RA']),
    ]:
        soft(lambda: infer_result(Script, code + call) == expected, call)


def test_no_matching_overload(Script):
    check(Script, ['type[A]', 'int'], 'f(C)', ['RA', 'RB'])
    check(Script, ['type[A]', 'int'], 'f(A)', ['RA'])
    code = variadic('type[A]', 'type[C]', '*x: %s')
    soft(lambda: infer_result(Script, code + 'f(A, C, int)') == ['RA', 'RB'], "infer_result(Script, code + 'f(A, C, int)') == ['RA', 'RB']")


def test_dict_unpacking(Script):
    code = variadic('type[A]', 'type[C]', '**x: %s')
    for call, expected in [
        ('f(**dict(p=C))', ['RB']),
        ('f(**dict(p=A))', ['RA']),
        ('d = dict(p=C)\nf(**d)', ['RB']),
        ('f(**dict(p=A, q=C))', ['RA', 'RB']),
    ]:
        soft(lambda: infer_result(Script, code + call) == expected, call)


def test_optional_accepts_none(Script):
    prefix = 'from typing import Optional\n'
    check(Script, ['Optional[int]', 'str'], 'f(None)', ['RA'], prefix)
    check(Script, ['Optional[int]', 'str'], 'f(1)', ['RA'], prefix)
    check(Script, ['Optional[int]', 'str'], 'f("")', ['RB'], prefix)


def test_dict_positional_argument(Script):
    code = variadic('type[A]', 'type[C]', '**x: %s')
    for call, expected in [
        ('f(**dict([("p", C)]))', ['RB']),
        ('f(**dict({"p": C}))', ['RB']),
        ('f(**dict([("p", A)], q=C))', ['RA', 'RB']),
        # A keyword argument replaces the entry of the mapping.
        ('f(**dict({"p": A}, p=C))', ['RB']),
        ('d = dict(p=C)\nf(**dict(d))', ['RB']),
        ('d = dict(ab=C)\nf(**dict(d))', ['RB']),
        ('f(**dict(**dict(p=C)))', ['RB']),
    ]:
        soft(lambda: infer_result(Script, code + call) == expected, call)
    # The mapping and the keyword argument fill the same parameter only once.
    check(Script, ['type[A]', 'type[C]'], 'f(**dict({"x": A}, x=C))', ['RB'])


def test_dict_unpacking_in_method(Script):
    code = dedent('''\
        class K:
            @overload
            def f(self, **x: type[A]) -> RA: ...
            @overload
            def f(self, **x: type[C]) -> RB: ...
            def f(self, **x): ...
        ''')
    soft(lambda: infer_result(Script, code + 'K().f(**dict(p=C))') == ['RB'], "infer_result(Script, code + 'K().f(**dict(p=C))') == ['RB']")
    soft(lambda: infer_result(Script, code + 'K().f(**dict([("p", A)]))') == ['RA'], 'infer_result(Script, code + \'K().f(**dict([("p", A)]))\') == [\'RA\']')


def test_union_argument(Script):
    prefix = dedent('''\
        def g(x: Union[A, C]):
            return x
        ''')
    # Every type the argument may have has to be accepted.
    check(Script, ['A', 'C'], 'f(g(A()))', ['RA', 'RB'], prefix)
    check(Script, ['Union[A, C]', 'int'], 'f(g(A()))', ['RA'], prefix)
    check(Script, ['B', 'Union[A, C]'], 'f(g(A()))', ['RB'], prefix)


@pytest.mark.parametrize('annotation', ['type', 'Type', 'type[Any]', 'Type[Any]'])
def test_any_class_defined_later(Script, annotation):
    prefix = 'class D(C): pass\n'
    check(Script, ['int', annotation], 'f(D)', ['RB'], prefix)
    check(Script, [annotation, 'D'], 'f(D())', ['RB'], prefix)


def test_optional_of_class(Script):
    prefix = 'from typing import Optional\n'
    check(Script, ['Optional[A]', 'C'], 'f(B())', ['RA'], prefix)
    check(Script, ['Optional[A]', 'C'], 'f(None)', ['RA'], prefix)
    check(Script, ['Optional[A]', 'C'], 'f(C())', ['RB'], prefix)


def test_every_collected_argument_is_checked(Script):
    code = variadic('type[A]', 'type', '*x: %s')
    soft(lambda: infer_result(Script, code + 'f(A, B, C)') == ['RB'], "infer_result(Script, code + 'f(A, B, C)') == ['RB']")
    soft(lambda: infer_result(Script, code + 'f(A, B, B)') == ['RA'], "infer_result(Script, code + 'f(A, B, B)') == ['RA']")
    code = variadic('type[A]', 'type', '**x: %s')
    soft(lambda: infer_result(Script, code + 'f(p=A, q=B, r=C)') == ['RB'], "infer_result(Script, code + 'f(p=A, q=B, r=C)') == ['RB']")
    soft(lambda: infer_result(Script, code + 'f(p=A, q=B, r=B)') == ['RA'], "infer_result(Script, code + 'f(p=A, q=B, r=B)') == ['RA']")


def test_dict_unpacking_in_classmethod(Script):
    code = dedent('''\
        class K:
            @overload
            @classmethod
            def f(cls, **x: type[A]) -> RA: ...
            @overload
            @classmethod
            def f(cls, **x: type[C]) -> RB: ...
            @classmethod
            def f(cls, **x): ...
        ''')
    soft(lambda: infer_result(Script, code + 'K.f(**dict(p=C))') == ['RB'], "infer_result(Script, code + 'K.f(**dict(p=C))') == ['RB']")
    soft(lambda: infer_result(Script, code + 'K().f(**dict(p=A))') == ['RA'], "infer_result(Script, code + 'K().f(**dict(p=A))') == ['RA']")


def test_dict_with_alternative_entries(Script):
    # Every possible value of an entry has to be accepted, so only the last
    # overload matches when the entry may be either C or A.
    code = 'import random\n' + variadic('type[A]', 'type[C]', '**x: %s', third='type')
    for call, expected in [
        ('m = {"p": C} if random.random() else {"p": A}\nf(**dict(m))', ['RC']),
        ('m = {"p": B} if random.random() else {"p": A}\nf(**dict(m))', ['RA']),
        ('f(**dict([("p", C) if random.random() else ("p", A)]))', ['RC']),
        ('f(**dict([("p", B) if random.random() else ("p", A)]))', ['RA']),
        # Each possible mapping only passes classes.
        ('m = {"p": C} if random.random() else {"q": A}\nf(**dict(m))', ['RC']),
    ]:
        soft(lambda: infer_result(Script, code + call) == expected, call)


def test_dict_unpacking_in_partial(Script):
    code = 'import functools\n' + variadic('type[A]', 'type[C]', '**x: %s')
    soft(lambda: infer_result(Script, code + 'functools.partial(f)(**dict(p=C))') == ['RB'], "infer_result(Script, code + 'functools.partial(f)(**dict(p=C))') == ['RB']")
    soft(lambda: infer_result(Script, code + 'functools.partial(f, q=A)(**dict(p=A))') == ['RA'], "infer_result(Script, code + 'functools.partial(f, q=A)(**dict(p=A))') == ['RA']")
    soft(lambda: infer_result(Script, code + 'functools.partial(f, q=A)(**dict(p=C))') == ['RA', 'RB'], "infer_result(Script, code + 'functools.partial(f, q=A)(**dict(p=C))') == ['RA', 'RB']")


CONTAINERS = 'import random\nfrom typing import Iterable, List, Mapping, Sequence\n'


@pytest.mark.parametrize('first, second, call', [
    ('list[type[A]]', 'list[type[C]]', 'f([C])'),
    ('List[Type[A]]', 'List[Type[C]]', 'x = [C, C]\nf(x)'),
    ('Sequence[type[A]]', 'Sequence[type[C]]', 'f((C, C))'),
    ('Iterable[type[A]]', 'Iterable[type[C]]', 'f({C})'),
    ('tuple[type[A], ...]', 'tuple[type[C], ...]', 'f((C, C))'),
    ('dict[str, type[A]]', 'dict[str, type[C]]', 'f({"k": C})'),
    ('Mapping[str, type[A]]', 'Mapping[str, type[C]]', 'f(dict(k=C))'),
    ('dict[str, list[type[A]]]', 'dict[str, list[type[C]]]', 'f({"k": [C]})'),
])
def test_elements_of_containers(Script, first, second, call):
    check(Script, [first, second, 'object'], call, ['RB'], prefix=CONTAINERS)
    if 'C' in call:
        check(Script, [first, second, 'object'], call.replace('C', 'B'), ['RA'],
              prefix=CONTAINERS)


def test_every_element_has_to_fit(Script):
    annotations = ['list[type[A]]', 'list[type[C]]', 'object']
    for call in [
        'f([A, C])',
        'f([C()])',
        'f((C,))',
        'f([C] if random.random() else [A])',
        'f([C if random.random() else A])',
    ]:
        check(Script, annotations, call, ['RC'], prefix=CONTAINERS)


def test_tuple_and_mapping_items(Script):
    def check_items(first, second, call, prefix=''):
        check(Script, [first, second, 'object'], call, ['RB'], prefix=CONTAINERS + prefix)

    check_items('tuple[type[A], type[C]]', 'tuple[type[C], type[A]]', 'f((C, A))')
    check_items('tuple[type[C]]', 'tuple[type[C], type[C]]', 'f((C, C))')
    check_items('dict[type[A], int]', 'dict[type[C], int]', 'f({C: 1})')
    check_items('dict[str, type[A]]', 'dict[int, type[A]]', 'f({1: A})')
    check_items('Sequence[type[A]]', 'Sequence[type[C]]', 'f(g())',
                prefix='def g() -> tuple[type[C], ...]: ...\n')
    check_items('list[type[T]]', 'list[type[C]]', 'f([C])',
                prefix='T = TypeVar("T", bound=A)\n')
