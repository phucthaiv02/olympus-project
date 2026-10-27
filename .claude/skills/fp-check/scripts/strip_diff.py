#!/usr/bin/env python3
"""Drop files from a unified diff. usage: strip_diff.py RUN.diff OUT.diff [--drop-from PATCH ...] [--drop-regex RE ...]
Always drops test.sh and, for every PATCH given, every file that PATCH touches (the tests and runner are re-applied afterwards)."""
import re, sys


def files_of(text):
    return set(re.findall(r"^diff --git a/(\S+) b/", text, re.M))


def main(argv):
    src, dst, rest = argv[1], argv[2], argv[3:]
    drop_files, drop_re, i = {"test.sh"}, [], 0
    while i < len(rest):
        if rest[i] == "--drop-from":
            drop_files |= files_of(open(rest[i + 1], errors="replace").read()); i += 2
        elif rest[i] == "--drop-regex":
            drop_re.append(re.compile(rest[i + 1])); i += 2
        else:
            sys.exit(f"unknown arg {rest[i]}")
    text = open(src, errors="replace").read()
    parts = re.split(r"(?m)^(?=diff --git )", text)
    kept, dropped = [], []
    for p in parts:
        m = re.match(r"diff --git a/(\S+) b/(\S+)", p)
        if not m:
            kept.append(p); continue
        f = m.group(2)
        if f in drop_files or m.group(1) in drop_files or any(r.search(f) for r in drop_re):
            dropped.append(f)
        else:
            kept.append(p)
    open(dst, "w").write("".join(kept))
    print(f"kept {len(kept) - (0 if kept and not kept[0].startswith('diff') else 0)} parts; dropped: {', '.join(dropped) or '-'}")


main(sys.argv)
