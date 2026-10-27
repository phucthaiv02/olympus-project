#!/usr/bin/env python3
"""usage: compare_probes.py RESULTS.json | RESULTS_DIR
Compare probe results of passing runs against the reference (and optional ground truth).

Input JSON:
{
  "probes": ["id1", "id2", ...],              # optional ordering
  "reference": {"id1": "A", "id2": "B"},
  "truth":     {"id1": "A", "id2": "B"},      # optional: real language behavior
  "runs": {"run-1": {"id1": "A", "id2": "C"}, "run-2": {...}}
}
Values are any comparable strings (outputs, verdicts). A missing probe in a run is reported as MISSING.
Exit 1 if any run differs from the reference.
"""
import json, sys

import os


def load_dir(path):
    """<dir>/reference.jsonl, optional truth.jsonl, and <label>.jsonl per run; lines {"probe": id, "result": str}."""
    def rd(f):
        out = {}
        for line in open(f):
            line = line.strip()
            if line:
                j = json.loads(line); out[j["probe"]] = j["result"]
        return out
    ref = rd(os.path.join(path, "reference.jsonl"))
    truth = rd(os.path.join(path, "truth.jsonl")) if os.path.exists(os.path.join(path, "truth.jsonl")) else {}
    runs = {f[:-6]: rd(os.path.join(path, f)) for f in sorted(os.listdir(path)) if f.endswith(".jsonl") and f not in ("reference.jsonl", "truth.jsonl")}
    return {"reference": ref, "truth": truth, "runs": runs}


d = load_dir(sys.argv[1]) if os.path.isdir(sys.argv[1]) else json.load(open(sys.argv[1]))
ref, truth, runs = d["reference"], d.get("truth", {}), d["runs"]
ids = d.get("probes") or sorted(ref)
bad = 0
print(f"{'probe':30} {'reference':>14} {'truth':>14}  differing runs")
for p in ids:
    diff = {r: v.get(p, "MISSING") for r, v in runs.items() if v.get(p, "MISSING") != ref.get(p)}
    t = truth.get(p, "-")
    flag = ""
    if t != "-" and ref.get(p) != t:
        flag = "  <-- REFERENCE != TRUTH"
        bad += 1
    if diff or flag:
        print(f"{p:30} {str(ref.get(p)):>14} {str(t):>14}  {diff}{flag}")
        bad += bool(diff)
n = len(ids)
print(f"\n{n} probes, {len(runs)} runs, {bad} probe(s) with a difference")
for r, v in runs.items():
    k = sum(1 for p in ids if v.get(p, "MISSING") != ref.get(p))
    print(f"  {r}: differs from reference on {k}/{n}")
sys.exit(1 if bad else 0)
