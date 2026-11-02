import sys, signal
from textwrap import dedent
import jedi

HEADER = """\
from typing import overload, Any, Type, TypeVar, Union, Protocol, Optional, Dict, List, Tuple, Mapping, Sequence, Iterable
import random
import functools
class A: pass
class B(A): pass
class C: pass
class RA: pass
class RB: pass
class RC: pass
"""

def ov(*anns, name='f', indent='', deco=''):
    code = ''
    for a, r in zip(anns, ('RA', 'RB', 'RC')):
        if deco:
            code += indent + deco + '\n'
        code += indent + '@overload\n' + indent + 'def %s(%s) -> %s: ...\n' % (name, a, r)
    if deco:
        code += indent + deco + '\n'
    return code + indent + 'def %s(*args, **kwargs): ...\n' % name

class TO(Exception): pass
def _h(*a): raise TO()
signal.signal(signal.SIGALRM, _h)

def infer(code):
    code = HEADER + code.strip('\n')
    lines = code.splitlines()
    signal.alarm(30)
    try:
        return sorted(d.name for d in jedi.Script(code).infer(len(lines), len(lines[-1])))
    except TO:
        return 'TIMEOUT'
    except Exception as e:
        import traceback; tb = traceback.format_exc().splitlines()[-4:]
        return 'CRASH %s: %s | %s' % (type(e).__name__, e, ' / '.join(tb))
    finally:
        signal.alarm(0)

def run(cases):
    for name, code, exp in cases:
        got = infer(code)
        print('%-4s %-40s got=%s expected=%s' % ('OK' if got == exp else 'BAD', name, got, exp), flush=True)
