package llm

import (
	"fmt"
	"os/exec"
	"strings"
)

const systemPrompt = `당신은 GitHub PR 변경사항을 분석하여 팀원이 따라잡아야 할 핵심 내용을 요약하는 역할입니다.

분석 관점:
1. 주요 기능 추가/변경사항
2. 중요한 기술적 결정 및 아키텍처 변경
3. 버그 수정 및 개선사항
4. 팀 리뷰에서 나온 피드백 및 학습 포인트
5. 컨트리뷰터가 알아야 할 코드 패턴/컨벤션

출력 형식:
- 한글로 작성
- 간결하고 실용적으로
- bullet points 사용
- 핵심만 추출 (장황하게 X)`

// Summarize sends PR data to Claude CLI and returns the summary.
func Summarize(prData, repo string, prCount int) (string, error) {
	userPrompt := fmt.Sprintf(`다음은 %s 레포지토리의 최근 머지된 PR %d개입니다.
컨트리뷰터로서 따라잡아야 할 핵심 내용을 요약해주세요.

---
%s
---

위 PR들을 분석하여 다음 섹션으로 요약해주세요:

## 📦 주요 변경사항
(새 기능, 개선, 리팩토링 등)

## 🐛 버그 수정
(있는 경우만)

## 💡 학습 포인트
(리뷰 코멘트에서 얻은 인사이트, 코드 패턴 등)

## ⚠️ 주의사항
(breaking changes, 마이그레이션 필요 등 - 있는 경우만)`, repo, prCount, prData)

	cmd := exec.Command("claude", "-p", "--system-prompt", systemPrompt)
	cmd.Stdin = strings.NewReader(userPrompt)

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("claude summarize: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}
