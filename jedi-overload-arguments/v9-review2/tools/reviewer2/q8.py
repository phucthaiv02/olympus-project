from lib import *
c = []
def case(n, prefix, anns, call, exp): c.append((n, prefix + ov(*anns) + call, exp))
case('dup direct kwonly', '', ['y: object = None, *, x: type[A] = A', 'y: object = None, z: object = None'], 'f(A, y=1)', ['RA', 'RB'])
case('dup direct', '', ['y: object = None, x: type[A] = A', 'z: object = None, y: object = None'], 'f(A, y=1)', ['RB'])
case('dup direct2', '', ['y: object, x: type[A] = A', 'z: object = None, y: object = None'], 'f(A, y=1)', ['RB'])
case('dup kwonly push', '', ['y: type[C] = C, *, x: type[A] = A', 'y: object = None, z: object = None'], 'f(A, y=C)', ['RB'])
case('kwonly push then kw', '', ['*, x: type[A] = A', 'z: object = None, x: object = None'], 'f(C, x=A)', ['RB'])
case('kwonly push then kw2', '', ['y: type[C] = C, *, x: type[A] = A', 'y: object = None, z: object = None, x: object = None'], 'f(C, A, x=A)', ['RB'])
run(c)
