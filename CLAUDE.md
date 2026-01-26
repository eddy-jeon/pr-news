# PR News - Project Documentation

## Overview

GitHub 레포지토리의 최근 머지된 PR들을 분석하여 컨트리뷰터가 따라잡아야 할 내용을 LLM 기반으로 요약해주는 TUI CLI 도구

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    pr-news (main)                       │
├─────────────────────────────────────────────────────────┤
│  [Phase 0] Repository Selection (interactive)          │
│  [Phase 1] PR Data Collection (gh CLI)                  │
│  [Phase 2] LLM Summarization (claude CLI)               │
│  [Phase 3] TUI Output                                   │
└─────────────────────────────────────────────────────────┘
         │
         ├── lib/config.sh   # Configuration management
         ├── lib/ui.sh       # TUI rendering (gum/fzf/bash)
         ├── lib/github.sh   # GitHub API wrapper
         └── lib/llm.sh      # Claude CLI wrapper
```

## Module Responsibilities

### `lib/config.sh`
- 설정 파일 로드 (`~/.pr-news.conf`)
- 기본값 관리
- 의존성 체크 (gh, claude, jq)

### `lib/ui.sh`
- TUI 렌더링 (Graceful Fallback: gum > fzf > bash)
- 색상 코드 및 스타일 관리
- 대화형 선택, 스피너, 박스 등

### `lib/github.sh`
- GitHub CLI 래퍼
- PR 목록/상세/diff 조회
- 리뷰 코멘트 조회 (봇 필터링)
- PR 크기 판단 로직

### `lib/llm.sh`
- Claude CLI 래퍼
- 요약 프롬프트 관리

## Key Design Decisions

### 1. Shell Script 선택
- 의존성 최소화 (gh, claude, jq만 필요)
- 빠른 개발 및 수정 가능
- 파이프라인 조합에 최적화

### 2. Graceful TUI Fallback
```
gum (설치됨) → fzf (설치됨) → bash select (항상)
```
- 어떤 환경에서도 동작 보장
- 선택적으로 더 나은 UX 제공

### 3. PR 크기 기반 분석 전략
- **큰 PR** (파일 >10 또는 변경 >500줄): 제목 + 본문만
- **작은 PR**: 제목 + 본문 + diff 발췌
- 토큰 효율성과 분석 품질 균형

### 4. 봇 필터링
- `authorAssociation != "NONE"` 또는 명시적 봇 목록으로 필터
- 의미 있는 팀원 피드백만 포함

## Configuration

| 변수 | 기본값 | 설명 |
|------|--------|------|
| `DAYS` | 7 | 조회 기간 |
| `THRESHOLD_FILES` | 10 | 큰 PR 기준 (파일) |
| `THRESHOLD_CHANGES` | 500 | 큰 PR 기준 (라인) |
| `INCLUDE_REVIEW_COMMENTS` | true | 리뷰 코멘트 포함 |
| `BOT_FILTER` | (목록) | 제외할 봇 |

## Dependencies

| 도구 | 필수 | 용도 |
|------|------|------|
| `gh` | ✅ | GitHub API 접근 |
| `claude` | ✅ | LLM 요약 생성 |
| `jq` | ✅ | JSON 파싱 |
| `gum` | ❌ | 예쁜 TUI (선택) |
| `fzf` | ❌ | 퍼지 선택 (선택) |

## Development

### 문법 체크
```bash
bash -n pr-news
bash -n lib/*.sh
```

### 모듈 테스트
```bash
source lib/config.sh && config::load && echo "Config OK"
source lib/ui.sh && ui::header
source lib/github.sh && gh::list_repos 5
```

### 디버깅
```bash
bash -x ./pr-news  # 실행 추적
```

## File Conventions

- 함수명: `module::function_name` (예: `ui::header`, `gh::list_repos`)
- 상수: `UPPER_SNAKE_CASE`
- 지역 변수: `lower_snake_case`
- `set -euo pipefail` 사용

## Future Improvements

- [ ] 캐싱 (중복 API 호출 방지)
- [ ] Markdown 파일 출력 옵션
- [ ] 특정 라벨/브랜치 필터링
- [ ] 다중 레포 동시 분석
