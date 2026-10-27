#!/usr/bin/env bash
# Usage: repo_check.sh owner/repo [owner/repo ...]
# Prints the hard eligibility facts for each repo: stars, language,
# license, archived flag, last commit on the default branch, and its SHA.
set -uo pipefail

ALLOWED_LICENSES="MIT Apache-2.0 BSD-2-Clause BSD-3-Clause ISC 0BSD Unlicense Zlib"
cutoff=$(date -u -d '12 months ago' +%Y-%m-%d 2>/dev/null || date -u -v-12m +%Y-%m-%d)

for repo in "$@"; do
  info=$(gh api "repos/$repo" --jq '[.stargazers_count, (.language // "?"), (.license.spdx_id // "NONE"), .archived, .default_branch] | @tsv' 2>/dev/null) || {
    echo "$repo  ERROR: cannot fetch repo"; continue; }
  IFS=$'\t' read -r stars lang license archived branch <<<"$info"
  commit=$(gh api "repos/$repo/commits/$branch" --jq '[.sha, .commit.committer.date[0:10]] | @tsv' 2>/dev/null)
  IFS=$'\t' read -r sha date <<<"$commit"

  problems=()
  [ "$stars" -lt 500 ] && problems+=("stars<500")
  [ "$archived" = "true" ] && problems+=("archived")
  [[ "$date" < "$cutoff" ]] && problems+=("no commit since $cutoff")
  case " TypeScript JavaScript Python Go Rust C++ Java " in
    *" $lang "*) ;; *) problems+=("language=$lang") ;; esac
  case " $ALLOWED_LICENSES " in
    *" $license "*) ;; *) problems+=("license=$license (confirm this license is acceptable)") ;; esac

  verdict="OK"; [ ${#problems[@]} -gt 0 ] && verdict="CHECK: ${problems[*]}"
  printf '%s\n  stars=%s lang=%s license=%s branch=%s\n  last_commit=%s sha=%s\n  %s\n' \
    "$repo" "$stars" "$lang" "$license" "$branch" "$date" "$sha" "$verdict"
done
