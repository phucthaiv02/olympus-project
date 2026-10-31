import jedi
from jedi.inference.param import ExecutedParamName
orig = ExecutedParamName._accepts
def spy(self, annotations, values):
    r = orig(self, annotations, values)
    print('ACC', self.string_name, [str(a)[:60] for a in annotations], [str(v)[:60] for v in values], r)
    return r
ExecutedParamName._accepts = spy
src = open('/r/test/completion/stdlib.py').read().splitlines()
code = '\n'.join(src[:229])
print(src[225:228])
print(jedi.Script(code).infer(229, 5))
