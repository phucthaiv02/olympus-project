#!/usr/bin/env python3
"""Translate a Vietnamese task description to English and run the description gate.

usage: make_description.py description_vi.md [--out description.md] [--no-translate]

Inline code spans (`like this`) are replaced by placeholders before translation and restored afterwards,
so identifiers and operators survive; paragraphs that are fenced code blocks are copied untranslated. After translating (or with --no-translate, on an existing English
file) it prints the gate: word count, requirement-clause estimate, and phrases that reviewers flag.
The English file is always produced from the Vietnamese one; to change the English, edit the Vietnamese and rerun.
"""
import argparse, os, re, sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
from translate_gtx import translate  # noqa: E402

FLAGS = [
    (r"(?m)^\s*([-*+]|\d+\.)\s", "list item (use full paragraphs)"),
    (r"(?m)^#{1,6}\s", "heading (use full paragraphs)"),
    (r"https?://", "URL in description"),
    (r"(?i)\b(currently|at the moment|today)\b", "describes current behavior"),
    (r"(?i)(unchanged|remain(s)? the same|keeps? working|still works?)", "restates what stays the same"),
    (r"(?i)(not part of this|out of scope|is not required)", "states what is out of scope"),
    (r"(?i)\b(hidden tests?|benchmark|agents?)\b", "mentions tests/benchmark/agents"),
    (r"(?i)\b(any|all|every) (form|way|kind)s? (that|which)", "possibly wider than the tests (check each claim has a test)"),
]


def protect(text):
    spans = []
    def sub(m):
        spans.append(m.group(0)); return f"XQ{len(spans)-1}X"
    return re.sub(r"`[^`\n]+`", sub, text), spans


def restore(text, spans):
    for i, s in enumerate(spans):
        text = re.sub(rf"XQ\s?{i}\s?X", lambda m: s, text)
    return text


def gate(en):
    words = len(en.split())
    sentences = [s for s in re.split(r"(?<=[.;:])\s+", en.replace("\n", " ")) if s.strip()]
    print(f"words: {words} (limit about 330)" + ("  <-- OVER" if words > 330 else ""))
    print(f"sentences/clauses (rough): {len(sentences)}; count real requirement clauses by hand, keep about 15 or fewer")
    for rx, msg in FLAGS:
        for m in re.finditer(rx, en):
            line = en[max(0, m.start() - 30): m.end() + 40].replace("\n", " ")
            print(f"flag: {msg}: ...{line}...")


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("vi"); ap.add_argument("--out", default="description.md"); ap.add_argument("--no-translate", action="store_true")
    a = ap.parse_args()
    if a.no_translate:
        gate(open(a.vi, encoding="utf-8").read()); return
    src = open(a.vi, encoding="utf-8").read()
    paras = [p.strip() for p in src.split("\n\n") if p.strip()]
    out = []
    for p in paras:
        if p.startswith("```"):  # fenced code block: keep as is
            out.append(p); continue
        q, spans = protect(p)
        out.append(restore(translate(q, "vi", "en"), spans))
    en = "\n\n".join(out) + "\n"
    open(a.out, "w", encoding="utf-8").write(en)
    print(f"wrote {a.out}")
    gate(en)


main()
