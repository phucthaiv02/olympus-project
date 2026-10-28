#!/usr/bin/env python3
"""Translate text with the public Google Translate gtx endpoint.

Usage:
  translate_gtx.py <file> [--sl auto] [--tl en]

Reads <file>, splits it into paragraphs (blank-line separated), translates each
paragraph with GET https://translate.googleapis.com/translate_a/single?client=gtx&...
and prints the translated paragraphs joined by blank lines to stdout.
Paragraphs are sent separately to stay well under URL length limits.
"""
import argparse
import json
import sys
import urllib.parse
import urllib.request

ENDPOINT = "https://translate.googleapis.com/translate_a/single"


def translate(text, sl, tl):
    params = urllib.parse.urlencode(
        {"client": "gtx", "sl": sl, "tl": tl, "dt": "t", "q": text}
    )
    req = urllib.request.Request(
        f"{ENDPOINT}?{params}", headers={"User-Agent": "Mozilla/5.0"}
    )
    with urllib.request.urlopen(req, timeout=30) as resp:
        data = json.loads(resp.read().decode("utf-8"))
    # Response: [[[translated_chunk, original_chunk, ...], ...], ...]
    return "".join(chunk[0] for chunk in data[0] if chunk and chunk[0])


def main():
    ap = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("file")
    ap.add_argument("--sl", default="auto")
    ap.add_argument("--tl", default="en")
    args = ap.parse_args()

    with open(args.file, encoding="utf-8") as f:
        paragraphs = [p.strip() for p in f.read().split("\n\n") if p.strip()]

    out = [translate(p, args.sl, args.tl) for p in paragraphs]
    sys.stdout.write("\n\n".join(out) + "\n")


if __name__ == "__main__":
    main()
