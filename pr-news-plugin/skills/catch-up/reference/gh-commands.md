# gh CLI Commands Reference

## PR 목록 조회 (머지된 PR)

```bash
# 특정 날짜 이후 머지된 PR 목록
gh pr list \
  --repo OWNER/REPO \
  --state merged \
  --search "merged:>=YYYY-MM-DD" \
  --limit 50 \
  --json number,title,body,additions,deletions,changedFiles,mergedAt,author,url

# 특정 브랜치에 머지된 PR만
gh pr list \
  --repo OWNER/REPO \
  --state merged \
  --search "merged:>=YYYY-MM-DD base:main" \
  --limit 50 \
  --json number,title,body,additions,deletions,changedFiles,mergedAt,author,url
```

## PR 상세 조회

```bash
gh pr view PR_NUMBER \
  --repo OWNER/REPO \
  --json number,title,body,additions,deletions,changedFiles,files,mergedAt,author,url
```

## PR diff 조회

```bash
# 전체 diff (작은 PR에만 사용)
gh pr diff PR_NUMBER --repo OWNER/REPO
```

## 리뷰 코멘트 조회

```bash
# PR 토론 코멘트 (봇 제외)
gh pr view PR_NUMBER --repo OWNER/REPO --json comments \
  --jq '.comments[] | select(.authorAssociation != "NONE") | {author: .author.login, body: .body}'

# 코드 리뷰 코멘트 (라인 레벨)
gh api "repos/OWNER/REPO/pulls/PR_NUMBER/comments" | \
  jq '.[] | {author: .user.login, path: .path, body: .body}'
```

## 날짜 계산

```bash
# macOS
date -v-7d +%Y-%m-%d

# Linux
date -d "7 days ago" +%Y-%m-%d
```

## 레포 감지

```bash
# 현재 디렉토리의 remote URL에서 owner/repo 추출
git remote get-url origin | sed -E 's/.*github\.com[:/]([^/]+\/[^/.]+)(\.git)?$/\1/'
```

## 봇 필터링

알려진 봇 목록 (authorAssociation 체크와 병행):
- `coderabbitai`
- `snyk-io-us`
- `dependabot`
- `github-actions`
- `codecov`
