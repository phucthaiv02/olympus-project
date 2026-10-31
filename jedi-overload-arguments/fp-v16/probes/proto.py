import sys
sys.path.insert(0, '/rv')
from probe import infer, ov
# (protocol method, impl method, expected) ; RA = protocol matches, RB = falls to object
C = [
 ("def m(self, x: int) -> None: ...", "def m(self, x: int) -> None: ...", 'RA'),
 ("def m(self, x: int = 1) -> None: ...", "def m(self, x: int) -> None: ...", 'RB'),
 ("def m(self, x: int = 1) -> None: ...", "def m(self, x: int = 2) -> None: ...", 'RA'),
 ("def m(self, x: int) -> None: ...", "def m(self, y: int) -> None: ...", 'RB'),
 ("def m(self, x: int, /) -> None: ...", "def m(self, y: int, /) -> None: ...", 'RA'),
 ("def m(self, x: int, /) -> None: ...", "def m(self, y: int) -> None: ...", 'RA'),
 ("def m(self, x: int) -> None: ...", "def m(self, x: int, /) -> None: ...", 'RB'),
 ("def m(self, x: int) -> None: ...", "def m(self, x: int, y: int) -> None: ...", 'RB'),
 ("def m(self, x: int) -> None: ...", "def m(self, x: int, y: int = 0) -> None: ...", 'RA'),
 ("def m(self, x: int) -> None: ...", "def m(self, *args: int) -> None: ...", 'RB'),
 ("def m(self, x: int) -> None: ...", "def m(self, *args: int, **kwargs: int) -> None: ...", 'RA'),
 ("def m(self, *args: int) -> None: ...", "def m(self) -> None: ...", 'RB'),
 ("def m(self, *args: int) -> None: ...", "def m(self, *args: str) -> None: ...", 'RB'),
 ("def m(self, *args: int) -> None: ...", "def m(self, *args: int) -> None: ...", 'RA'),
 ("def m(self, **kw: int) -> None: ...", "def m(self) -> None: ...", 'RB'),
 ("def m(self, **kw: int) -> None: ...", "def m(self, **kw: str) -> None: ...", 'RB'),
 ("def m(self, *, k: int) -> None: ...", "def m(self, *, k: int) -> None: ...", 'RA'),
 ("def m(self, *, k: int) -> None: ...", "def m(self, k: int) -> None: ...", 'RA'),
 ("def m(self, *, k: int) -> None: ...", "def m(self, *, j: int) -> None: ...", 'RB'),
 ("def m(self, *, k: int = 1) -> None: ...", "def m(self, *, k: int) -> None: ...", 'RB'),
 ("def m(self, x: 'int') -> None: ...", "def m(self, x: 'str') -> None: ...", 'RB'),
 ("def m(self, x: 'int') -> None: ...", "def m(self, x: 'int') -> None: ...", 'RA'),
 ("def m(self, x: int) -> None: ...", "def m(self, x: str) -> None: ...", 'RB'),
 ("def m(self) -> 'int': ...", "def m(self) -> 'str': ...", 'RB'),
]
bad = 0
for p, i, e in C:
    pre = "class P(Protocol):\n    %s\nclass Impl:\n    %s\n" % (p, i)
    got = infer(pre + ov('x: P', 'x: object') + 'f(Impl())')
    ok = got == [e]
    bad += not ok
    print('OK ' if ok else 'BAD', p, '|', i, got)
print('bad', bad)
