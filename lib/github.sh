#!/bin/bash
#
# PR News - GitHub API Module
# ===========================
#
# GitHub CLI (gh)를 래핑하여 PR 데이터를 조회합니다.
#
# Functions:
#   gh::list_repos          - 접근 가능한 레포 목록 (개인 + 조직)
#   gh::list_merged_prs     - 머지된 PR 목록 (날짜 필터)
#   gh::get_pr_detail       - PR 상세 정보
#   gh::get_pr_diff         - PR diff (작은 PR용)
#   gh::get_review_comments - 리뷰 코멘트 (봇 필터링)
#   gh::is_large_pr         - PR 크기 판단
#   gh::collect_pr_data     - PR 데이터 수집 (크기별 전략)
#
# Dependencies:
#   - gh (GitHub CLI, authenticated)
#   - jq (JSON processing)
#

# List accessible repositories (personal + organizations)
gh::list_repos() {
  local limit="${1:-30}"

  {
    # Personal repos
    gh repo list --limit "$limit" --json nameWithOwner -q '.[].nameWithOwner' 2>/dev/null

    # Organization repos
    local orgs=$(gh api user/orgs --jq '.[].login' 2>/dev/null)
    for org in $orgs; do
      gh repo list "$org" --limit "$limit" --json nameWithOwner -q '.[].nameWithOwner' 2>/dev/null
    done
  } | sort -u
}

# List merged PRs in date range (with optional base branch filter)
gh::list_merged_prs() {
  local repo="$1"
  local days="${2:-$DAYS}"
  local base_branch="${3:-}"

  local since_date=$(date -v-${days}d +%Y-%m-%d 2>/dev/null || date -d "$days days ago" +%Y-%m-%d)

  # Build search query
  local search_query="merged:>=$since_date"
  if [[ -n "$base_branch" ]]; then
    search_query+=" base:$base_branch"
  fi

  gh pr list \
    --repo "$repo" \
    --state merged \
    --search "$search_query" \
    --limit 50 \
    --json number,title,body,additions,deletions,changedFiles,mergedAt,author,url
}

# Get PR detail
gh::get_pr_detail() {
  local repo="$1"
  local pr_number="$2"

  gh pr view "$pr_number" \
    --repo "$repo" \
    --json number,title,body,additions,deletions,changedFiles,files,mergedAt,author,url
}

# Get PR diff (for small PRs)
gh::get_pr_diff() {
  local repo="$1"
  local pr_number="$2"
  local max_lines="${3:-500}"

  gh pr diff "$pr_number" --repo "$repo" 2>/dev/null | head -n "$max_lines"
}

# Get review comments (excluding bots)
gh::get_review_comments() {
  local repo="$1"
  local pr_number="$2"
  local bot_filter="$BOT_FILTER"

  # Build jq filter for bots
  local jq_filter='.[] | select(.user.login as $author | ["'"${bot_filter//,/\",\"}"'"] | index($author) | not)'

  # PR discussion comments
  local discussion_comments=$(gh pr view "$pr_number" --repo "$repo" --json comments \
    --jq ".comments[] | select(.authorAssociation != \"NONE\") | {author: .author.login, body: .body}" 2>/dev/null)

  # Code review comments (line-level)
  local review_comments=$(gh api "repos/$repo/pulls/$pr_number/comments" 2>/dev/null | \
    jq -r "$jq_filter | {author: .user.login, path: .path, body: .body}" 2>/dev/null)

  # Combine and output
  {
    echo "$discussion_comments"
    echo "$review_comments"
  } | jq -s '.' 2>/dev/null
}

# Check if PR is large
gh::is_large_pr() {
  local files="$1"
  local changes="$2"

  [[ $files -gt ${THRESHOLD_FILES:-10} ]] || [[ $changes -gt ${THRESHOLD_CHANGES:-500} ]]
}

# Collect PR data with appropriate detail level
gh::collect_pr_data() {
  local repo="$1"
  local pr_json="$2"

  local number=$(echo "$pr_json" | jq -r '.number')
  local title=$(echo "$pr_json" | jq -r '.title')
  local body=$(echo "$pr_json" | jq -r '.body // ""')
  local additions=$(echo "$pr_json" | jq -r '.additions')
  local deletions=$(echo "$pr_json" | jq -r '.deletions')
  local files=$(echo "$pr_json" | jq -r '.changedFiles')
  local author=$(echo "$pr_json" | jq -r '.author.login')
  local merged_at=$(echo "$pr_json" | jq -r '.mergedAt')
  local url=$(echo "$pr_json" | jq -r '.url')

  local changes=$((additions + deletions))

  # Start building output
  local output=""
  output+="## PR #$number: $title\n"
  output+="- Author: $author\n"
  output+="- Merged: $merged_at\n"
  output+="- Stats: +$additions -$deletions ($files files)\n"
  output+="- URL: $url\n"
  output+="\n### Description\n$body\n"

  # Include diff for small PRs
  if ! gh::is_large_pr "$files" "$changes"; then
    local diff=$(gh::get_pr_diff "$repo" "$number")
    if [[ -n "$diff" ]]; then
      output+="\n### Code Changes (excerpt)\n\`\`\`diff\n$diff\n\`\`\`\n"
    fi
  else
    output+="\n> Large PR - showing summary only\n"
  fi

  # Include review comments if enabled
  if [[ "${INCLUDE_REVIEW_COMMENTS:-true}" == "true" ]]; then
    local comments=$(gh::get_review_comments "$repo" "$number")
    local comment_count=$(echo "$comments" | jq 'length' 2>/dev/null || echo 0)

    if [[ "$comment_count" -gt 0 && "$comment_count" != "null" ]]; then
      output+="\n### Review Comments ($comment_count)\n"
      output+=$(echo "$comments" | jq -r '.[] | "- **\(.author)**: \(.body | split("\n")[0])"' 2>/dev/null | head -10)
      output+="\n"
    fi
  fi

  echo -e "$output"
}
