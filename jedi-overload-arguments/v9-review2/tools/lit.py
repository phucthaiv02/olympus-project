import jedi
from jedi.inference.param import ExecutedParamName
orig = ExecutedParamName._accepts
def spy(self, annotations, values):
    r = orig(self, annotations, values)
    print('ACCEPTS', [(type(a).__name__, str(a)[:70]) for a in annotations], [str(v)[:50] for v in values], r)
    return r
ExecutedParamName._accepts = spy
H = 'from typing import overload, Literal, Callable\nclass RA: pass\nclass RB: pass\n'
for sig, call in [('Literal[1]', 'f(2)'), ('Literal[1]', 'f("s")'), ('Literal["a", "b"]', 'f("a")'), ('Callable[[], int]', 'f(1)'), ('Callable', 'f(len)')]:
    code = H + '@overload\ndef f(x: %s) -> RA: ...\n@overload\ndef f(x: object) -> RB: ...\ndef f(x): ...\n%s' % (sig, call)
    lines = code.splitlines()
    print('==', sig, call, [d.name for d in jedi.Script(code).infer(len(lines), len(lines[-1]))])
