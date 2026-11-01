import sys
sys.path.insert(0, sys.argv[1])
import jedi
code = """from typing import overload
class RA: pass
class RB: pass
@overload
def f(x: None) -> RA: ...
@overload
def f(x: object) -> RB: ...
def f(x): ...
f(None)"""
jedi.set_debug_function(lambda color, s: print(s[:200]))
print(jedi.Script(code).infer(9, 7))
