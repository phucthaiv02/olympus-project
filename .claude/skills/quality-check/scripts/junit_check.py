#!/usr/bin/env python3
"""Check junit files of a submission against the verify contract.

usage: junit_check.py --base-tests-only B.xml --new-tests-only N.xml \
           [--base-solution B2.xml --new-solution N2.xml] [--new-pkg PATTERN ...]
Files are optional; each rule runs only when its inputs are present.
"""
import argparse, collections, sys
import xml.etree.ElementTree as ET


def load(path):
    out = []  # (classname, name, status)
    for tc in ET.parse(path).getroot().iter("testcase"):
        st = "pass"
        if tc.find("failure") is not None or tc.find("error") is not None:
            st = "fail"
        elif tc.find("skipped") is not None:
            st = "skip"
        out.append((tc.get("classname", ""), tc.get("name", ""), st))
    return out


def ids(rows):
    return {f"{c}::{n}" for c, n, _ in rows}


ap = argparse.ArgumentParser()
ap.add_argument("--base-tests-only"); ap.add_argument("--new-tests-only")
ap.add_argument("--base-solution"); ap.add_argument("--new-solution")
ap.add_argument("--new-pkg", action="append", default=[])
a = ap.parse_args()
bad = 0


def report(ok, msg):
    global bad
    print(("PASS  " if ok else "FAIL  ") + msg)
    bad += (not ok)


def dups(rows, label):
    cnt = collections.defaultdict(list)
    for c, n, s in rows:
        cnt[f"{c}::{n}"].append(s)
    d = {k: v for k, v in cnt.items() if len(v) > 1}
    conflict = {k: v for k, v in d.items() if len(set(v)) > 1}
    report(not conflict, f"{label}: {len(d)} duplicate names, {len(conflict)} with conflicting results (before_p2p_not_passing)")
    for k, v in list(conflict.items())[:10]:
        print("      ", k, v)


if a.base_tests_only:
    b = load(a.base_tests_only)
    nonpass = [f"{c}::{n}" for c, n, s in b if s == "fail"]
    skipped = {f"{c}::{n}" for c, n, s in b if s == "skip"}
    report(not nonpass, f"base without solution: no failures among {len(b)} tests ({len(nonpass)} failing, {len(skipped)} skipped)")
    dups(b, "base junit")
    for pat in a.new_pkg:
        hit = [x for x in ids(b) if pat in x]
        report(not hit, f"base junit has no test from new package '{pat}' ({len(hit)} found; p2p_f2p_overlap)")
if a.new_tests_only:
    n = load(a.new_tests_only)
    passing = [f"{c}::{x}" for c, x, s in n if s == "pass" and any(p in c or p in x for p in a.new_pkg)] if a.new_pkg else []
    report(not passing, f"new without solution: no new test passes ({len(passing)} pass; they do not capture the missing behavior)")
    for t in passing[:10]:
        print("      ", t)
    dups(n, "new junit (no solution)")
if a.base_solution:
    b2 = load(a.base_solution)
    bad_t = [f"{c}::{n}" for c, n, s in b2 if s == "fail"]
    sk2 = {f"{c}::{n}" for c, n, s in b2 if s == "skip"}
    report(not bad_t, f"base with solution: {len(bad_t)} failing, {len(sk2)} skipped")
    if a.base_tests_only:
        new_skips = sk2 - skipped
        report(not new_skips, f"solution does not add skipped base tests ({len(new_skips)} newly skipped)")
    dups(b2, "base junit (solution)")
if a.new_solution:
    n2 = load(a.new_solution)
    bad_t = [f"{c}::{n}" for c, n, s in n2 if s != "pass"]
    report(not bad_t, f"new with solution: {len(bad_t)} not passing (failing or skipped)")
    for t in bad_t[:10]:
        print("      ", t)
    if a.base_solution:
        both = ids(load(a.base_solution)) & ids(n2)
        report(not both, f"base and new test sets disjoint ({len(both)} shared)")
sys.exit(1 if bad else 0)
