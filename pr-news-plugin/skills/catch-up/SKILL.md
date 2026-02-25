---
description: "최근 머지된 PR을 분석하여 팀원이 따라잡아야 할 내용을 요약합니다. 트리거: '최근 변경사항', 'PR 뉴스', 'what changed', 'catch up on PRs', 'what did I miss', '최근 PR 뭐가 바뀌었어', '머지된 PR 요약해줘', 'PR catch up'"
allowed-tools:
  - Bash(gh pr list *)
  - Bash(gh pr view *)
  - Bash(gh pr diff *)
  - Bash(gh api repos/*)
  - Bash(git remote get-url *)
  - Bash(date *)
  - Bash(jq *)
---

# PR News - 최근 PR 분석 및 요약

GitHub 레포지토리의 최근 머지된 PR들을 분석하여 컨트리뷰터가 따라잡아야 할 내용을 요약합니다.

## 워크플로우

다음 4단계를 순서대로 수행하세요.

### Step 1: 레포지토리 및 파라미터 결정

**먼저 사용자의 메시지를 파싱하여 이미 제공된 정보를 추출합니다:**
- 레포지토리 (`owner/repo` 형식)
- 조회 기간 (일수)
- 대상 브랜치

**레포지토리 결정:**
1. 사용자가 메시지에서 레포를 명시했으면 → 그대로 사용
2. 명시하지 않았으면 → 현재 디렉토리에서 감지:
   ```bash
   git remote get-url origin 2>/dev/null | sed -E 's/.*github\.com[:/]([^/]+\/[^/.]+)(\.git)?$/\1/'
   ```
3. 감지 결과가 있으면 → 그대로 사용 (확인 질문 없이)
4. 감지 실패 → 사용자에게 `owner/repo` 형식으로 입력 요청

**파라미터 기본값:**

| 파라미터 | 기본값 | 설명 |
|----------|--------|------|
| 조회 기간 | 7일 | 며칠 전까지의 PR을 볼지 |
| 대상 브랜치 | 전체 | 특정 브랜치만 볼지 (예: main, develop) |

**핵심 원칙**: 사용자가 이미 메시지에서 레포, 기간, 브랜치를 언급했다면 해당 항목에 대해 추가 질문 없이 바로 진행합니다. 모든 정보가 제공되었다면 확인 질문 없이 곧바로 Step 2로 넘어갑니다.

### Step 2: PR 목록 조회

날짜를 계산하고 머지된 PR 목록을 가져옵니다.

```bash
# 날짜 계산 (macOS/Linux 호환)
date -v-{DAYS}d +%Y-%m-%d 2>/dev/null || date -d "{DAYS} days ago" +%Y-%m-%d
```

```bash
gh pr list \
  --repo {REPO} \
  --state merged \
  --search "merged:>={SINCE_DATE}" \
  --limit 50 \
  --json number,title,body,additions,deletions,changedFiles,mergedAt,author,url
```

브랜치 필터가 있으면 `--search` 에 `base:{BRANCH}` 를 추가합니다.

PR이 없으면 사용자에게 알리고 기간 확대를 제안합니다.

### Step 3: PR별 상세 데이터 수집

각 PR에 대해 크기를 판단하고 적절한 수준의 정보를 수집합니다.

#### PR 크기 판단 기준

- **대형 PR**: `changedFiles > 10` 또는 `(additions + deletions) > 500`
- **소형 PR**: 그 외

#### 수집 전략

**모든 PR 공통:**
- 제목, 작성자, 머지 일시, 변경 통계 (additions/deletions/files)
- PR 본문 (description)

**소형 PR 추가:**
- diff 발췌 (처음 500줄)
  ```bash
  gh pr diff {PR_NUMBER} --repo {REPO} | head -500
  ```

**리뷰 코멘트 수집 (모든 PR):**

토론 코멘트:
```bash
gh pr view {PR_NUMBER} --repo {REPO} --json comments \
  --jq '.comments[] | select(.authorAssociation != "NONE") | {author: .author.login, body: .body}'
```

코드 리뷰 코멘트:
```bash
gh api "repos/{REPO}/pulls/{PR_NUMBER}/comments" | \
  jq '[.[] | select(.user.login as $a | ["coderabbitai","snyk-io-us","dependabot","github-actions","codecov"] | index($a) | not)] | .[:10] | .[] | {author: .user.login, path: .path, body: .body}'
```

> **중요**: 봇 필터링 - `authorAssociation == "NONE"` 인 코멘트 제외, 그리고 `coderabbitai`, `snyk-io-us`, `dependabot`, `github-actions`, `codecov` 사용자의 코멘트 제외

### Step 4: 분석 및 요약 출력

수집한 모든 PR 데이터를 종합 분석하여 아래 형식으로 한글 요약을 출력합니다.

---

**분석 관점:**
1. 주요 기능 추가/변경사항
2. 중요한 기술적 결정 및 아키텍처 변경
3. 버그 수정 및 개선사항
4. 팀 리뷰에서 나온 피드백 및 학습 포인트
5. 컨트리뷰터가 알아야 할 코드 패턴/컨벤션

**출력 형식:**

```
## 📦 주요 변경사항
(새 기능, 개선, 리팩토링 등을 bullet point로)

## 🐛 버그 수정
(있는 경우만 - 없으면 섹션 생략)

## 💡 학습 포인트
(리뷰 코멘트에서 얻은 인사이트, 코드 패턴, 기술적 결정 등)

## ⚠️ 주의사항
(breaking changes, 마이그레이션 필요 등 - 없으면 섹션 생략)
```

**작성 원칙:**
- 간결하고 실용적으로 작성
- 핵심만 추출 (장황하게 X)
- 각 항목에 해당 PR 번호 표기 (예: `#123`)
- 관련 PR의 URL을 링크로 포함
