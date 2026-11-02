from lib import *
c = []
def case(n, prefix, anns, call, exp): c.append((n, prefix + ov(*anns) + call, exp))
D = "def g() -> dict[str, type[C]]: ...\nd = g()\n"
case('ann dict passed', D, ['x: dict[str, type[A]]', 'x: dict[str, type[C]]', 'x: object'], 'f(d)', ['RB'])
case('ann dict param passed', "def h(d: dict[str, type[C]]):\n    return f(d)\n", ['x: dict[str, type[A]]', 'x: dict[str, type[C]]', 'x: object'], 'h({})', ['RB'])
case('ann Dict passed', "def g() -> Dict[str, type[C]]: ...\nd = g()\n", ['x: dict[str, type[A]]', 'x: dict[str, type[C]]', 'x: object'], 'f(d)', ['RB'])
case('ann dict in literal', D, ['x: dict[str, type[A]]', 'x: dict[str, type[C]]', 'x: object'], 'f({**d})', ['RB'])
case('ann dict in literal2', D, ['x: dict[str, type[A]]', 'x: object'], 'f({"k": A, **d})', ['RB'])
case('ann Mapping passed', "def g() -> Mapping[str, type[C]]: ...\nd = g()\n", ['x: Mapping[str, type[A]]', 'x: Mapping[str, type[C]]', 'x: object'], 'f(d)', ['RB'])
case('ann list passed', "def g() -> list[type[C]]: ...\nd = g()\n", ['x: list[type[A]]', 'x: list[type[C]]', 'x: object'], 'f(d)', ['RB'])
case('int key for float', '', ['x: dict[float, type[C]]', 'x: object'], 'f({1: C})', ['RA'])
case('int for float', '', ['x: float', 'x: object'], 'f(1)', ['RA'])
case('first key kept int', '', ['x: dict[int, type[C]]', 'x: object'], 'f({1: A, 1.0: C})', ['RA'])
case('first key kept float', '', ['x: dict[str, type[C]]', 'x: dict[float, type[C]]', 'x: object'], 'f({1.0: A, 1: C})', ['RB'])
case('star dict keys', '', ['x: str', 'x: str, y: str'], 'f(*{"a": A, "a": C})', ['RA'])
case('star dict keys2', '', ['x: int', 'x: object, y: object'], 'f(*{1: A, True: C})', ['RA'])
case('star dict var', 'd = {"a": A, "b": C}\n', ['x: str, y: str', 'x: object'], 'f(*d)', ['RA'])
run(c)
