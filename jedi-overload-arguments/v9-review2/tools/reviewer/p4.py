from lib import *
c = []
def M(anns, deco='', call='', exp=None, name='m', prefix='', n=''):
    code = prefix + 'class K:\n' + ov(*[(('self, ' if not deco else ('cls, ' if 'class' in deco else '')) + a) for a in anns], name=name, indent='    ', deco=deco)
    c.append((n or (deco + ' ' + call), code + call, exp))
for deco, obj in (('', 'K()'), ('@classmethod', 'K'), ('@classmethod', 'K()'), ('@staticmethod', 'K'), ('@staticmethod', 'K()')):
    M(['x: type[A]', 'x: type[C]'], deco, obj + '.m(C)', ['RB'])
    M(['x: type[A]', 'x: type[C]'], deco, obj + '.m(x=C)', ['RB'])
    M(['x: type[A]', 'x: type[C]'], deco, obj + '.m(*[C])', ['RB'])
    M(['x: type[A]', 'x: type[C]'], deco, obj + '.m(**{"x": C})', ['RB'])
    M(['x: type[A] = A, y: type[C] = C', 'x: object = None, y: object = None'], deco, 'xs = [A] if random.random() else []\n' + obj + '.m(*xs, C)', ['RB'])
    M(['x: type[A]', 'x: type[C]'], deco, 'def g() -> tuple[type[C], ...]: ...\n' + obj + '.m(*g())', ['RB'])
    M(['x: type[A]', 'x: type[C]'], deco, 'def g() -> dict[str, type[C]]: ...\n' + obj + '.m(**g())', ['RB'])
    M(['*a: type[A]', '*a: type[C]'], deco, 'def g() -> tuple[type[C], ...]: ...\n' + obj + '.m(*g())', ['RB'])
    M(['**k: type[A]', '**k: type[C]'], deco, 'def g() -> dict[str, type[C]]: ...\n' + obj + '.m(**g())', ['RB'])
    M(['**k: type[A]', '**k: type[C]'], deco, 'd = {"x": A} if random.random() else {"y": C}\n' + obj + '.m(**d)', ['RB'] if False else ['RB'])
run(c)
