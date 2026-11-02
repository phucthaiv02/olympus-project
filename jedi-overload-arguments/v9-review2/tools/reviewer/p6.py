from lib import *
G = '''from typing import Generic
T = TypeVar("T")
class K(Generic[T]):
    @overload
    def __init__(self, x: type[A], y: T) -> None: ...
    @overload
    def __init__(self, x: type[C], y: list[T]) -> None: ...
    def __init__(self, x, y): ...
    def get(self) -> T: ...
'''
c = [
 ('init generic C', G + 'K(C, [RA()]).get()', ['RA']),
 ('init generic A', G + 'K(A, [RA()]).get()', ['list']),
 ('init generic kw', G + 'K(x=C, y=[RA()]).get()', ['RA']),
 ('init generic dstar', G + 'K(**{"x": C, "y": [RA()]}).get()', ['RA']),
]
G2 = '''from typing import Generic
T = TypeVar("T")
class K(Generic[T]):
    @overload
    def __init__(self, x: type[A], y: T) -> None: ...
    @overload
    def __init__(self, x: type[C], y: T) -> None: ...
    def __init__(self, x, y): ...
    def get(self) -> T: ...
'''
run(c)
