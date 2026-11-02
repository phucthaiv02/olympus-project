from lib import *
G0 = '''from typing import Generic
T = TypeVar("T")
class K(Generic[T]):
    def __init__(self, x: type[A], y: T) -> None: ...
    def get(self) -> T: ...
'''
G = '''from typing import Generic
T = TypeVar("T")
class K(Generic[T]):
    @overload
    def __init__(self, x: type[A], y: T) -> None: ...
    @overload
    def __init__(self, x: type[C], y: list[T]) -> None: ...
    def get(self) -> T: ...
'''
c = [
 ('plain', G0 + 'K(A, RA()).get()', ['RA']),
 ('ov stub-like C', G + 'K(C, [RA()]).get()', ['RA']),
 ('ov stub-like A', G + 'K(A, [RA()]).get()', ['list']),
]
run(c)
