# pr-news - Claude Code Plugin

GitHub 레포지토리의 최근 머지된 PR들을 분석하여 팀원이 따라잡아야 할 내용을 요약하는 Claude Code 플러그인입니다.

## 설치

```bash
# Claude Code에서 플러그인 설치
claude plugin add /path/to/pr-news-plugin
```

또는 `.claude/settings.json`에 직접 추가:

```json
{
  "plugins": [
    "/path/to/pr-news-plugin"
  ]
}
```

## 요구사항

- **gh** (GitHub CLI) - 인증 완료 상태 (`gh auth login`)
- **jq** - JSON 파싱

> `claude` CLI는 불필요합니다. Claude Code 자체가 분석을 수행합니다.

## 사용법

Claude Code 세션에서 자연어로 트리거합니다:

```
"최근 PR 뭐가 바뀌었어?"
"catch up on PRs"
"what did I miss"
"최근 변경사항 알려줘"
"PR 뉴스"
"지난 2주간 머지된 PR 요약해줘"
```

## 기본 설정

| 항목 | 기본값 | 설명 |
|------|--------|------|
| 조회 기간 | 7일 | 분석할 PR 범위 |
| 대상 브랜치 | 전체 | 특정 브랜치 필터 |
| 대형 PR 기준 | 파일 10개 또는 500줄 이상 | diff 포함 여부 결정 |
| 봇 필터 | coderabbitai, dependabot 등 | 리뷰 코멘트에서 제외 |

## 출력 예시

```
## 📦 주요 변경사항
- 인증 시스템을 JWT에서 OAuth2로 마이그레이션 (#142)
- 대시보드에 실시간 알림 기능 추가 (#138)

## 🐛 버그 수정
- 로그인 시 세션 만료 처리 오류 수정 (#140)

## 💡 학습 포인트
- React Query v5의 새 캐싱 전략 도입 (#138)
- E2E 테스트에서 MSW 대신 Playwright의 route 사용 권장 (#139)

## ⚠️ 주의사항
- API 응답 형식 변경: `data.items` → `data.results` (#142)
```

## 구조

```
pr-news-plugin/
├── .claude-plugin/
│   └── plugin.json          # 플러그인 메타데이터
├── skills/
│   └── catch-up/
│       ├── SKILL.md         # 핵심 워크플로우
│       └── reference/
│           └── gh-commands.md   # gh CLI 레퍼런스
└── README.md
```
