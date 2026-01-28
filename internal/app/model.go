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
	prData  string
	prCount int
	repo    string

	width  int
	height int
}

func NewModel() Model {
	return Model{
		State:  StateLoading,
		Input:  panel.NewInputPanel(),
		Output: panel.NewOutputPanel(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.Output.Init(),
		loadReposCmd(),
	)
}
