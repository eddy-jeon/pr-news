package app

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/eddy/pr-news/internal/style"
)

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	inputW := m.width*4/10 - 4
	outputW := m.width*6/10 - 4
	panelH := m.height - 2

	left := style.InputPanel.
		Width(inputW).
		Height(panelH).
		Render(m.Input.View())

	right := style.OutputPanel.
		Width(outputW).
		Height(panelH).
		Render(m.Output.View())

	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}
