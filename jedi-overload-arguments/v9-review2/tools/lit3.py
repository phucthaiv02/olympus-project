import jedi
for src in ['from typing import Literal\ndef h(p: Literal[1]):\n    p.', 'from typing import Literal\nx: Literal["a"] = "a"\nx.', 'from typing import Literal\nL = Literal[1, 2]\ndef h(p: L):\n    p.', 'import typing\ndef h(p: typing.Literal[1]):\n    p.']:
    lines = src.splitlines(); print(len(jedi.Script(src).complete(len(lines), len(lines[-1]))), [d.description[:40] for d in jedi.Script(src).infer(len(lines), len(lines[-1]) - 1)])
