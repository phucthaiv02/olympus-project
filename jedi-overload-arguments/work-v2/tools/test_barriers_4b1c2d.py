"""
Choosing the overload of a function from the arguments of a call.
"""

from textwrap import dedent

import pytest

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


def check(Script, annotations, call, expected, prefix=''):
    code = prefix + overloads(*annotations) + call
    assert infer_result(Script, code) == expected, (annotations, call)


def variadic(first, second, signature, third=None):
    code = ''
    for annotation, result in [(first, 'RA'), (second, 'RB'), (third, 'RC')]:
        if annotation is not None:
            code += '@overload\ndef f(%s) -> %s: ...\n' % (signature % annotation, result)
    return code + 'def f(*args, **kwargs): ...\n'


CONTAINERS = 'import random\nfrom typing import Iterable, List, Mapping, Sequence\n'


def test_b1_nested_star_args(Script):
    code = variadic('type[A]', 'type[C]', '*x: %s')
    for call, expected in [
        ('f(*[*[A], C])', ['RA', 'RB']),
        ('f(*[A, *[C]])', ['RA', 'RB']),
        ('f(*[*[C]])', ['RB']),
        ('f(*(*[C], C))', ['RB']),
    ]:
        assert infer_result(Script, code + call) == expected, call
    check(Script, ['type[A]', 'type[C]'], 'f(*[*[C]])', ['RB'])
    check(Script, ['type[A]', 'type[C]'], 'f(*[*[A]])', ['RA'])


def test_b2_star_in_sequence_literal(Script):
    for call in ['f([C, *[A]])', 'f([*[A], C])', 'f([*[C], *(A,)])']:
        check(Script, ['list[type[A]]', 'list[type[C]]', 'object'], call, ['RC'], prefix=CONTAINERS)
    check(Script, ['Sequence[type[A]]', 'Sequence[type[C]]', 'object'], 'f((C, *[A]))', ['RC'],
          prefix=CONTAINERS)


def test_b3_star_star_in_dict_literal(Script):
    for call in ['f({"k": C, **{"m": A}})', 'f({**{"m": A}, "k": C})']:
        check(Script, ['dict[str, type[A]]', 'dict[str, type[C]]', 'object'], call, ['RC'],
              prefix=CONTAINERS)


def test_b4_mapping_with_keys(Script):
    mapping = dedent('''\
        class M:
            def keys(self):
                return ["x"]

            def __getitem__(self, key):
                return C
        ''')
    check(Script, ['type[A]', 'type[C]'], 'f(**dict(M()))', ['RB'], prefix=mapping)
    code = mapping + variadic('type[A]', 'type[C]', '**x: %s')
    assert infer_result(Script, code + 'f(**dict(M()))') == ['RB']
    check(Script, ['dict[str, type[A]]', 'dict[str, type[C]]', 'object'], 'f(dict(M()))',
          ['RB'], prefix=CONTAINERS + mapping)


def test_b5_comprehensions(Script):
    code = variadic('type[A]', 'type[C]', '**x: %s')
    assert infer_result(Script, code + 'f(**{k: C for k in "ab"})') == ['RB']
    assert infer_result(Script, code + 'f(**{k: C for k in []})') == ['RA']


def test_b6_forward_type_var(Script):
    bound = 'T = TypeVar("T", bound="A")\n'
    check(Script, ['type[T]', 'type[C]'], 'f(B)', ['RA'], prefix=bound)
    check(Script, ['type[T]', 'type[C]'], 'f(C)', ['RB'], prefix=bound)
    constrained = 'T = TypeVar("T", "D", int)\nclass D: pass\nclass E(D): pass\n'
    check(Script, ['type[T]', 'type[C]'], 'f(E)', ['RA'], prefix=constrained)
    check(Script, ['type[T]', 'type[C]'], 'f(int)', ['RA'], prefix=constrained)
    check(Script, ['type[T]', 'type[C]'], 'f(C)', ['RB'], prefix=constrained)


def test_b7_uncertain_partial(Script):
    prefix = 'import functools\nimport random\ndef key() -> str: ...\n'
    code = prefix + variadic('type[A]', 'type[C]', '**x: %s')
    call = 'functools.partial(f, **{key(): A})(**{key(): C})'
    assert infer_result(Script, code + call) == ['RA', 'RB']
    maybe = prefix + 'm = {"x": C} if random.random() else {}\n'
    check(Script, ['type[A]', 'object'], 'functools.partial(f, **m)(x=A)', ['RA'], maybe)


def test_b8_conditional_overwrite(Script):
    for call in ['f(dict({"k": A} if random.random() else {}, k=C))',
                 'f({**({"k": A} if random.random() else {}), "k": C})']:
        check(Script, ['dict[str, type[A]]', 'dict[str, type[C]]', 'object'], call, ['RB'],
              prefix=CONTAINERS)


def test_b9_tuple_unknown_length(Script):
    prefix = 'def g() -> tuple[int, ...]: ...\ndef h() -> tuple[int]: ...\n'
    check(Script, ['tuple[int]', 'tuple[int, ...]'], 'f(g())', ['RB'], prefix=prefix)
    check(Script, ['tuple[int]', 'tuple[int, ...]'], 'f(h())', ['RA'], prefix=prefix)
    check(Script, ['tuple[int]', 'tuple[int, ...]'], 'f(tuple([1]))', ['RB'], prefix=prefix)


def test_b10_fairness_extras(Script):
    check(Script, ['type[T]', 'type[C]'], 'f(str)', ['RA'], prefix='T = TypeVar("T")\n')
    prefix = 'import functools\n'
    check(Script, ['type[A]', 'type[C]'], 'functools.partial(f, A)()', ['RA'], prefix)
    check(Script, ['type[A]', 'type[C]'], 'functools.partial(f, C)()', ['RB'], prefix)
    code = variadic('type[A]', 'type[C]', '**x: %s')
    assert infer_result(Script, code + 'f(**dict({"p": A}, **{"p": C}))') == ['RB']
    assert infer_result(Script, code + 'f(**dict([("p", C)], **{"q": A}))') == ['RA', 'RB']
