import json, re, sys
sys.path.insert(0, '.')
from cases import PRELUDE, CASES
ref = open('refgo.out').read().split('@@CASE ')[1:]
wants = {}
for chunk in ref:
    idx, _, body = chunk.partition('\n')
    wants[int(idx)] = body.rstrip('\n')
assert len(wants) == len(CASES)
groups = []
for i, (g, n, code) in enumerate(CASES):
    if not groups or groups[-1][0] != g:
        groups.append((g, []))
    groups[-1][1].append((n, code, wants[i]))
names = {"GetInfo": "TestDebugGetInfo_a2959f", "GetInfoNames": "TestDebugGetInfoNames_a2959f", "ShortSource": "TestDebugShortSource_a2959f",
         "Errors": "TestDebugArgumentErrors_a2959f", "GetLocal": "TestDebugGetLocal_a2959f", "SetLocal": "TestDebugSetLocal_a2959f", "Hooks": "TestDebugHooks_a2959f"}
out = ['''package lua

import "testing"

// debugPrelude_a2959f defines run(f), which calls f in protected mode and
// renders its results. Tables are rendered with sorted keys, func fields as
// their type and activelines as a sorted list of lines.
const debugPrelude_a2959f = `''' + PRELUDE.strip('\n') + '''
`

type debugCase_a2959f struct {
	name string
	code string
	want string
}

func runDebugCases_a2959f(t *testing.T, cases []debugCase_a2959f) {
	t.Helper()
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			l := NewState()
			OpenLibraries(l)
			if err := DoString(l, debugPrelude_a2959f); err != nil {
				t.Fatalf("prelude: %v", err)
			}
			l.Global("run")
			if err := LoadBuffer(l, c.code, "=case", "t"); err != nil {
				t.Fatalf("load: %v", err)
			}
			if err := l.ProtectedCall(1, 1, 0); err != nil {
				t.Fatalf("run: %v", err)
			}
			got, _ := l.ToString(-1)
			if got != c.want {
				t.Errorf("got  %s\\nwant %s\\ncode:\\n%s", got, c.want, c.code)
			}
		})
	}
}
''']
for g, cs in groups:
    out.append(f'func {names[g]}(t *testing.T) {{\n\trunDebugCases_a2959f(t, []debugCase_a2959f{{')
    for n, code, want in cs:
        assert '`' not in code
        out.append(f'\t\t{{\n\t\t\tname: {json.dumps(n)},\n\t\t\tcode: `{code}\n`,\n\t\t\twant: {json.dumps(want)},\n\t\t}},')
    out.append('\t})\n}\n')
open('debug_library_a2959f_test.go', 'w').write("\n".join(out))
