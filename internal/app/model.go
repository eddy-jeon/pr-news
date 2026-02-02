package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/eddy/pr-news/internal/panel"
)

type AppState int

const (
	StateLoading AppState = iota
	StateInput
	StateFetching
	StateSummarizing
	StateDone
	StateError
)

type Model struct {
	State  AppState
	Input  panel.InputPanel
	Output panel.OutputPanel

	// collected data
	prData    string
	prCount   int
	repo      string
	dateRange string // PR 기간 (예: "2026-01-26 ~ 2026-02-02")

	width  int
	height int
}

func NewModel() Model {
	o := panel.NewOutputPanel()
	o.State = panel.OutputLoading
	return Model{
		State:  StateLoading,
		Input:  panel.NewInputPanel(),
		Output: o,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.Input.Init(),
		m.Output.Init(),
		loadReposCmd(),
	)
}
