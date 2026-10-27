#!/usr/bin/env python3
"""Flag suspicious changes in agents' diffs. usage: diff_audit.py RUN_DIR_OR_DIFF..."""
import glob, os, re, sys

SKIP = re.compile(r"^\+.*(pytest\.mark\.(skip|xfail)|@unittest\.skip|\bt\.Skip\w*\(|\bxit\(|\bxdescribe\(|\bit\.skip\(|\btest\.skip\(|#\[ignore\]|@Disabled|@Ignore|\bDISABLED_)")
RUNNER = re.compile(r"(^|/)(test\.sh|conftest\.py|pytest\.ini|tox\.ini|setup\.cfg|jest\.config\.\w+|vitest\.config\.\w+|\.github/)")
TESTPATH = re.compile(r"(^|/)(tests?|__tests__|spec)(/|$)|(_test\.\w+|\.test\.\w+|test_\w+\.py|\.spec\.\w+)$")

for arg in sys.argv[1:]:
    files = [arg] if os.path.isfile(arg) else sorted(glob.glob(os.path.join(arg, "*.patch")) + glob.glob(os.path.join(arg, "*.diff")))[:1]
    for f in files:
        print(f"== {arg}")
        text = open(f, errors="replace").read()
        flags = 0
        for part in re.split(r"(?m)^(?=diff --git )", text):
            m = re.match(r"diff --git a/(\S+) b/(\S+)", part)
            if not m: continue
            path = m.group(2); new = "new file mode" in part or "--- /dev/null" in part
            removed = [l for l in part.splitlines() if l.startswith("-") and not l.startswith("---")]
            if RUNNER.search(path):
                print(f"  RUNNER/CONFIG changed: {path}"); flags += 1
            if TESTPATH.search(path) and not new:
                if removed:
                    print(f"  EXISTING TEST edited, {len(removed)} line(s) removed/changed: {path}"); flags += 1
                else:
                    print(f"  info: lines appended to existing test file {path} (nothing removed)")
            elif TESTPATH.search(path) and removed:
                print(f"  test file rewritten: {path}"); flags += 1
            for l in part.splitlines():
                if SKIP.match(l):
                    print(f"  SKIP marker added in {path}: {l[1:].strip()[:90]}"); flags += 1
        print("  nothing suspicious" if not flags else f"  {flags} flag(s): read these hunks")
