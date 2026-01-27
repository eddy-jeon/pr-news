# 🗞️ PR News

GitHub 레포지토리의 최근 머지된 PR들을 분석하여 컨트리뷰터가 따라잡아야 할 내용을 LLM 기반으로 요약해주는 TUI CLI 도구

> "휴가 다녀왔는데 뭐가 바뀌었지?" 를 해결합니다.

## Quick Start

```bash
# 1. 클론 및 실행 권한
git clone https://github.com/your-username/pr-news.git && cd pr-news
chmod +x pr-news

# 2. 실행
./pr-news
```

## Features

- 📦 대화형 레포지토리 선택 (개인 + 조직)
- 📊 PR 크기에 따른 스마트 분석 (큰 PR은 요약만, 작은 PR은 diff 포함)
- 💬 팀원 리뷰 코멘트 포함 (봇 자동 필터링)
- 🤖 Claude LLM으로 종합 요약 생성
- 🎨 Graceful TUI (gum > fzf > bash fallback)

## Requirements

### 필수
- [GitHub CLI (`gh`)](https://cli.github.com/) - 인증 필요
- [Claude CLI (`claude`)](https://claude.ai/code) - Anthropic 계정 필요
- `jq` - JSON 처리

### 선택적 (더 예쁜 TUI)
- [`gum`](https://github.com/charmbracelet/gum) - Charm TUI toolkit
- [`fzf`](https://github.com/junegunn/fzf) - Fuzzy finder

## Installation

```bash
# Clone repository
git clone https://github.com/your-username/pr-news.git
cd pr-news

# Make executable
chmod +x pr-news

# Optional: Add to PATH
ln -s $(pwd)/pr-news /usr/local/bin/pr-news

# Optional: Install gum for better TUI
brew install gum
```

## Usage

```bash
# Run interactively
./pr-news

# Or if added to PATH
pr-news
```

### Flow

1. **Repository Selection** - 접근 가능한 레포 중 선택
2. **Options** - 조회 기간(일) 입력 및 대상 브랜치 선택
3. **PR Fetching** - 머지된 PR 조회
4. **Data Collection** - PR 상세 정보 수집
5. **LLM Analysis** - Claude로 종합 요약 생성
6. **Summary Output** - 터미널에 결과 출력

## Configuration

설정 파일을 `~/.pr-news.conf`에 생성:

```bash
cp .pr-news.conf.example ~/.pr-news.conf
```

### Options

| Option | Default | Description |
|--------|---------|-------------|
| `DAYS` | 7 | 조회할 기간 (일) |
| `THRESHOLD_FILES` | 10 | 큰 PR 기준 (파일 수) |
| `THRESHOLD_CHANGES` | 500 | 큰 PR 기준 (변경 라인) |
| `INCLUDE_REVIEW_COMMENTS` | true | 리뷰 코멘트 포함 여부 |
| `BOT_FILTER` | (see file) | 제외할 봇 목록 (쉼표 구분) |

### Environment Variables

설정 파일 대신 환경 변수로도 지정 가능:

```bash
DAYS=14 ./pr-news
```

## Output Example

```
╭──────────────────────────────────────────────────────────╮
│ 🗞️  PR News - GitHub PR Learning Tool                   │
╰──────────────────────────────────────────────────────────╯

█ Repository Selection
──────────────────────────────────────────────────
✓ Selected: chequer-io/querypie-mono

█ Fetching PRs (last 7 days)
──────────────────────────────────────────────────
✓ Found 5 merged PRs

█ Analyzing PRs
──────────────────────────────────────────────────
[1/5] PR #14821: feat(apps/api): SAC Admin MCP Tool...
...
✓ Collected data from 5 PRs

█ PR News Summary for chequer-io/querypie-mono
──────────────────────────────────────────────────

## 📦 주요 변경사항
- SAC Admin MCP Tool에 Role CRUD 기능 추가
- Redis TLS 지원 및 클러스터 자동 감지
...

## 💡 학습 포인트
- 문자열 상수는 분리하는 것이 좋음
...
```

## TUI Modes

| 환경 | 선택 UI | 스피너 | 스타일 |
|------|---------|--------|--------|
| gum 설치됨 | gum filter (검색) | gum spin | gum style |
| fzf만 있음 | fzf (vim keys) | 텍스트 | ANSI 색상 |
| 둘 다 없음 | bash select | 텍스트 | 기본 |

## Keyboard Shortcuts (Vim-style)

레포 선택 시 사용 가능한 키:

| Key | Action |
|-----|--------|
| `j` / `↓` | 아래로 이동 |
| `k` / `↑` | 위로 이동 |
| `Ctrl+d` | 반 페이지 아래 |
| `Ctrl+u` | 반 페이지 위 |
| `Ctrl+f` | 한 페이지 아래 |
| `Ctrl+b` | 한 페이지 위 |
| `/` | 검색 토글 |
| `Enter` | 선택 |
| `Esc` | 취소 |

**gum 사용 시**: 바로 타이핑하면 검색됩니다 (fuzzy filter)

## Troubleshooting

### "Missing required dependencies" 에러

```bash
# gh 설치
brew install gh
gh auth login

# jq 설치
brew install jq

# claude 설치 (https://claude.ai/code 참고)
```

### "GitHub CLI not authenticated" 에러

```bash
gh auth login
# 브라우저에서 인증 진행
```

### PR이 없다고 나오는 경우

- `DAYS` 값을 늘려보세요: `DAYS=30 ./pr-news`
- 해당 레포에 머지된 PR이 있는지 확인: `gh pr list --repo OWNER/REPO --state merged`

### gum이 설치되어 있는데 기본 UI가 나오는 경우

```bash
# gum 경로 확인
which gum

# PATH에 gum이 있는지 확인
echo $PATH
```

## Project Structure

```
pr-news/
├── pr-news                    # 메인 실행 스크립트
├── lib/
│   ├── config.sh              # 설정 로드 및 의존성 체크
│   ├── ui.sh                  # TUI 렌더링 (gum/fzf/bash fallback)
│   ├── github.sh              # GitHub API 래퍼 (gh CLI)
│   └── llm.sh                 # Claude CLI 래퍼
├── .pr-news.conf.example      # 설정 파일 템플릿
├── CLAUDE.md                  # 프로젝트 기술 문서
└── README.md                  # 사용자 가이드
```

## How It Works

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Select    │───▶│   Fetch     │───▶│   Analyze   │───▶│  Summarize  │
│    Repo     │    │    PRs      │    │    Data     │    │   (LLM)     │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
     gum/fzf           gh CLI         크기별 전략       claude CLI
```

1. **레포 선택**: 접근 가능한 개인/조직 레포 목록에서 대화형 선택
2. **PR 조회**: 최근 N일간 머지된 PR 목록 가져오기
3. **데이터 수집**:
   - 작은 PR: 제목 + 본문 + diff + 리뷰 코멘트
   - 큰 PR: 제목 + 본문만 (토큰 효율성)
4. **LLM 요약**: Claude가 전체 내용을 분석하여 학습 포인트 도출

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT
