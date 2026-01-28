package app

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/eddy/pr-news/internal/github"
	"github.com/eddy/pr-news/internal/llm"
	"github.com/eddy/pr-news/internal/panel"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		inputW := msg.Width*4/10 - 4
		outputW := msg.Width*6/10 - 4
		m.Input.Width = inputW
		m.Input.Height = msg.Height - 2
		m.Output.SetSize(outputW, msg.Height-2)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			if m.State == StateDone || m.State == StateError {
				return m, tea.Quit
			}
		case "r":
			if m.State == StateDone || m.State == StateError {
				m.State = StateInput
				m.Output.State = panel.OutputIdle
				m.prData = ""
				return m, nil
			}
		}

	case panel.StartSearchMsg:
		if m.State == StateInput {
			return m, m.startFetch()
		}

	case ReposLoadedMsg:
		if msg.Err != nil {
			m.State = StateError
			m.Output.State = panel.OutputError
			m.Output.Error = msg.Err.Error()
			return m, nil
		}
		m.Input.SetRepos(msg.Repos)
		m.State = StateInput
		m.Output.State = panel.OutputIdle
		return m, nil

	case PRsFetchedMsg:
		if msg.Err != nil {
			m.State = StateError
			m.Output.State = panel.OutputError
			m.Output.Error = msg.Err.Error()
			return m, nil
		}
		m.prCount = len(msg.PRs)
		if m.prCount == 0 {
			m.State = StateError
			m.Output.State = panel.OutputError
			m.Output.Error = "No merged PRs found"
			return m, nil
		}
		m.Output.Status = fmt.Sprintf("Collecting data from %d PRs...", m.prCount)
		return m, collectPRDataCmd(m.repo, msg.PRs)

	case PRDataCollectedMsg:
		m.prData = msg.Data
		m.State = StateSummarizing
		m.Output.State = panel.OutputSummarizing
		m.Output.Status = "Claude is analyzing..."
		m.Output.Progress = fmt.Sprintf("%d PRs collected", msg.Total)
		return m, summarizeCmd(m.prData, m.repo, m.prCount)

	case SummaryDoneMsg:
		if msg.Err != nil {
			m.State = StateError
			m.Output.State = panel.OutputError
			m.Output.Error = msg.Err.Error()
			return m, nil
		}
		m.State = StateDone
		m.Output.State = panel.OutputDone
		m.Output.SetContent(msg.Summary)
		return m, nil
	}

	// Delegate to panels
	var cmds []tea.Cmd
	if m.State == StateInput || m.State == StateLoading {
		var cmd tea.Cmd
		m.Input, cmd = m.Input.Update(msg)
		cmds = append(cmds, cmd)
	}

	var cmd tea.Cmd
	m.Output, cmd = m.Output.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) startFetch() tea.Cmd {
	repo := m.Input.SelectedRepo()
	if repo == "" {
		return nil
	}
	m.repo = repo
	m.State = StateFetching
	m.Output.State = panel.OutputFetching
	m.Output.Status = fmt.Sprintf("Fetching merged PRs from %s...", repo)

	daysStr := m.Input.Days.Value()
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		days = 7
	}
	branch := strings.TrimSpace(m.Input.Branch.Value())

	return fetchPRsCmd(repo, days, branch)
}

func loadReposCmd() tea.Cmd {
	return func() tea.Msg {
		repos, err := github.ListRepos(30)
		return ReposLoadedMsg{Repos: repos, Err: err}
	}
}

func fetchPRsCmd(repo string, days int, branch string) tea.Cmd {
	return func() tea.Msg {
		prs, err := github.ListMergedPRs(repo, days, branch)
		return PRsFetchedMsg{PRs: prs, Err: err}
	}
}

func collectPRDataCmd(repo string, prs []github.PR) tea.Cmd {
	return func() tea.Msg {
		var b strings.Builder
		for i, pr := range prs {
			_ = i
			data := github.CollectPRData(repo, pr)
			b.WriteString(data)
			b.WriteString("\n---\n")
		}
		return PRDataCollectedMsg{
			Data:    b.String(),
			Current: len(prs),
			Total:   len(prs),
		}
	}
}

func summarizeCmd(prData, repo string, count int) tea.Cmd {
	return func() tea.Msg {
		summary, err := llm.Summarize(prData, repo, count)
		return SummaryDoneMsg{Summary: summary, Err: err}
	}
}
