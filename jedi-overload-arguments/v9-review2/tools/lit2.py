import jedi
H = '''from typing import overload, Literal, Callable, Optional, Union, List
import functools
class RA: pass
class RB: pass
class K:
    def m(self): pass
    def __call__(self): pass
class E: pass
def g(): return 1
'''
C = [('Literal[1]', 'f(1)', 'RA'), ('Literal[1]', 'f(2)', 'RB'), ('Literal[1]', 'f("s")', 'RB'), ('Literal[1]', 'f(True)', 'RB'),
     ('Literal["a", "b"]', 'f("a")', 'RA'), ('Literal["a", "b"]', 'f("b")', 'RA'), ('Literal["a", "b"]', 'f("c")', 'RB'),
     ('Literal[None]', 'f(None)', 'RA'), ('Literal[True]', 'f(True)', 'RA'), ('Literal[True]', 'f(False)', 'RB'), ('Literal[True]', 'f(1)', 'RB'),
     ('Optional[Literal[1]]', 'f(None)', 'RA'), ('Optional[Literal[1]]', 'f(1)', 'RA'), ('Optional[Literal[1]]', 'f(2)', 'RB'),
     ('Union[Literal[1], str]', 'f("s")', 'RA'), ('Union[Literal[1], str]', 'f(2)', 'RB'),
     ('list[Literal[1]]', 'f([1])', 'RA'), ('list[Literal[1]]', 'f([2])', 'RB'),
     ("'Literal[1]'", 'f(2)', 'RB'), ('Literal[1]', 'v = 1\nf(v)', 'RA'), ('Literal[1]', 'f(1 if g() else 2)', 'RB'),
     ('Literal[1]', 'def h(p: int):\n    f(p)', None),
     ('Callable[[], int]', 'f(g)', 'RA'), ('Callable[[], int]', 'f(lambda: 1)', 'RA'), ('Callable[[], int]', 'f(1)', 'RB'), ('Callable[[], int]', 'f("s")', 'RB'), ('Callable[[], int]', 'f(None)', 'RB'),
     ('Callable[[], int]', 'f(K)', 'RA'), ('Callable[[], int]', 'f(K())', 'RA'), ('Callable[[], int]', 'f(K().m)', 'RA'), ('Callable[[], int]', 'f(K.m)', 'RA'), ('Callable[[], int]', 'f(len)', 'RA'), ('Callable[[], int]', 'f(E())', 'RB'),
     ('Callable[[], int]', 'f(functools.partial(g))', 'RA'), ('Callable', 'f(len)', 'RA'), ('Callable', 'f(1)', 'RB'), ('Callable[..., int]', 'f(g)', 'RA'), ('Optional[Callable[[], int]]', 'f(None)', 'RA'), ('Optional[Callable[[], int]]', 'f(g)', 'RA'), ('Optional[Callable[[], int]]', 'f(1)', 'RB'),
     ('list[Callable[[], int]]', 'f([g])', 'RA'), ('list[Callable[[], int]]', 'f([1])', 'RB')]
bad = 0
for sig, call, exp in C:
    code = H + '@overload\ndef f(x: %s) -> RA: ...\n@overload\ndef f(x: object) -> RB: ...\ndef f(x): ...\n%s' % (sig, call)
    lines = code.splitlines()
    if lines[-1].startswith('    '): pass
    try: got = ','.join(sorted(d.name for d in jedi.Script(code).infer(len(lines), len(lines[-1]))))
    except Exception as e: got = 'CRASH ' + repr(e)[:90]
    ok = exp is None or got == exp
    bad += not ok
    print('OK ' if ok else 'BAD', sig, '|', call.replace('\n', ' ; '), '|', got, exp)
print('bad', bad)
print([c.name for c in jedi.Script('from typing import Literal\ndef h(p: Literal[1]):\n    p.').complete(3, 6)][:3])
H += '''from typing import Iterator, Generator, TypeVar, Any
T = TypeVar('T')
def gi() -> int: ...
def gs() -> str: ...
def gn() -> None: ...
def gf() -> Iterator[float]: ...
def gg() -> Generator[int, None, None]: ...
def ga() -> Any: ...
def gu(): return 1
def gb() -> bool: ...
'''
bad = 0
for sig, call, exp in [('Callable[[], int]', 'f(gi)', 'RA'), ('Callable[[], int]', 'f(gs)', 'RB'), ('Callable[[], int]', 'f(gb)', 'RA'), ('Callable[[], int]', 'f(gu)', 'RA'), ('Callable[[], int]', 'f(ga)', 'RA'),
    ('Callable[[], float]', 'f(gi)', 'RA'), ('Callable[[], None]', 'f(gn)', 'RA'), ('Callable[[], None]', 'f(gi)', 'RB'), ('Callable[[], object]', 'f(gs)', 'RA'), ('Callable[[], T]', 'f(gs)', 'RA'), ('Callable[[], Any]', 'f(gs)', 'RA'),
    ('Callable[[], Iterator[float]]', 'f(gf)', 'RA'), ('Callable[[], Generator[float, None, None]]', 'f(gf)', 'RB'), ('Callable[[], Iterator[int]]', 'f(gg)', 'RA'), ('Callable[..., str]', 'f(gi)', 'RB'), ('Callable[[], int]', 'f(K().m)', 'RA'), ('Callable[[], int]', 'f(lambda: "s")', 'RA')]:
    code = H + '@overload\ndef f(x: %s) -> RA: ...\n@overload\ndef f(x: object) -> RB: ...\ndef f(x): ...\n%s' % (sig, call)
    lines = code.splitlines()
    try: got = ','.join(sorted(d.name for d in jedi.Script(code).infer(len(lines), len(lines[-1]))))
    except Exception as e: got = 'CRASH ' + repr(e)[:90]
    ok = got == exp; bad += not ok
    print('OK ' if ok else 'BAD', sig, '|', call, '|', got, exp)
print('bad2', bad)
