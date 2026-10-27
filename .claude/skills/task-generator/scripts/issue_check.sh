#!/usr/bin/env bash
# Usage:
#   issue_check.sh find owner/repo ["keywords"]   list open feature-request-like issues, most-reacted first
#   issue_check.sh show owner/repo <number>       state, labels, reactions, and every PR that references the issue
set -uo pipefail
mode="$1"; repo="$2"

case "$mode" in
  find)
    kw="${3:-}"
    gh search issues --repo "$repo" --state open --sort reactions --limit 40 $kw \
      --json number,title,labels,commentsCount,url \
      --jq '.[] | "#\(.number) [\([.labels[].name] | join(","))] (\(.commentsCount) comments) \(.title)  \(.url)"'
    ;;
  show)
    num="$3"
    gh api "repos/$repo/issues/$num" --jq '"#\(.number) [\(.state)] \(.title)\n  author=\(.user.login) (\(.author_association)) labels=\([.labels[].name] | join(",")) reactions=\(.reactions.total_count) comments=\(.comments)\n  \(.html_url)"'
    echo "  Referencing PRs (any state):"
    gh api --paginate "repos/$repo/issues/$num/timeline" \
      --jq '.[] | select(.event=="cross-referenced" and .source.issue.pull_request != null) | "    PR #\(.source.issue.number) [\(.source.issue.state)\(if .source.issue.pull_request.merged_at then ",merged" else "" end)] \(.source.issue.repository.full_name): \(.source.issue.title)"' 2>/dev/null
    echo "  Maintainer comments (OWNER/MEMBER/COLLABORATOR):"
    gh api --paginate "repos/$repo/issues/$num/comments" \
      --jq '.[] | select(.author_association=="OWNER" or .author_association=="MEMBER" or .author_association=="COLLABORATOR") | "    \(.user.login) \(.created_at[0:10]): \(.body | gsub("\\s+";" ") | .[0:200])"'
    ;;
  *) echo "unknown mode: $mode (expected find or show)" >&2; exit 2 ;;
esac
