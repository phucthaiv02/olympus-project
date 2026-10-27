#!/usr/bin/env bash
# Usage: dup_search.sh owner/repo "keyword one" "keyword two" ...
# Searches PRs (open, merged, closed) and issues for each keyword phrase,
# so you can see whether an idea was already implemented, proposed or declined.
set -uo pipefail
repo="$1"; shift
for kw in "$@"; do
  echo "=== \"$kw\" ==="
  echo "--- PRs (all states)"
  gh search prs --repo "$repo" "$kw" --limit 15 \
    --json number,title,state,url --jq '.[] | "#\(.number) [\(.state)] \(.title)  \(.url)"' 2>/dev/null
  echo "--- Issues (all states)"
  gh search issues --repo "$repo" "$kw" --limit 15 \
    --json number,title,state,url --jq '.[] | "#\(.number) [\(.state)] \(.title)  \(.url)"' 2>/dev/null
  echo "--- Discussions"
  gh api graphql -f q="repo:$repo $kw" -f query='query($q:String!){search(query:$q,type:DISCUSSION,first:10){nodes{... on Discussion{number title url}}}}' \
    --jq '.data.search.nodes[] | "#\(.number) \(.title)  \(.url)"' 2>/dev/null
done
