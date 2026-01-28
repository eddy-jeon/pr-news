package panel

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/eddy/pr-news/internal/style"
)

type FocusField int

const (
	FocusFilter FocusField = iota
	FocusDays
	FocusBranch
	FocusFieldCount
)

// StartSearchMsg is sent when the user completes all fields and presses Enter.
type StartSearchMsg struct{}

type InputPanel struct {
	Repos    []string
	filtered []string
	cursor   int
	Loading  bool

	Filter textinput.Model
	Days   textinput.Model
	Branch textinput.Model
	focus  FocusField

	spinner spinner.Model

	Width  int
	Height int
}

func NewInputPanel() InputPanel {
	filter := textinput.New()
	filter.Placeholder = "type to filter..."
	filter.Focus()

	days := textinput.New()
	days.Placeholder = "7"
	days.SetValue("7")
	days.CharLimit = 4

	branch := textinput.New()
	branch.Placeholder = "all branches"

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = style.CursorStyle

	return InputPanel{
		Filter:  filter,
		Days:    days,
		Branch:  branch,
		focus:   FocusFilter,
		Loading: true,
		spinner: s,
	}
}

func (p *InputPanel) SetRepos(repos []string) {
	p.Repos = repos
	p.Loading = false
	p.applyFilter()
}

func (p *InputPanel) applyFilter() {
	q := strings.ToLower(p.Filter.Value())
	if q == "" {
		p.filtered = p.Repos
	} else {
		p.filtered = nil
		for _, r := range p.Repos {
			if strings.Contains(strings.ToLower(r), q) {
				p.filtered = append(p.filtered, r)
			}
		}
	}
	if p.cursor >= len(p.filtered) {
		p.cursor = max(0, len(p.filtered)-1)
	}
}

func (p *InputPanel) SelectedRepo() string {
	if len(p.filtered) == 0 {
		return ""
	}
	return p.filtered[p.cursor]
}

func (p *InputPanel) focusNext() {
	p.focus = (p.focus + 1) % FocusFieldCount
	p.syncFocus()
}

func (p *InputPanel) focusPrev() {
	p.focus = (p.focus - 1 + FocusFieldCount) % FocusFieldCount
	p.syncFocus()
}

func (p *InputPanel) syncFocus() {
	p.Filter.Blur()
	p.Days.Blur()
	p.Branch.Blur()
	switch p.focus {
	case FocusFilter:
		p.Filter.Focus()
	case FocusDays:
		p.Days.Focus()
	case FocusBranch:
		p.Branch.Focus()
	}
}

func (p InputPanel) Init() tea.Cmd {
	return p.spinner.Tick
}

func (p InputPanel) Update(msg tea.Msg) (InputPanel, tea.Cmd) {
	var cmds []tea.Cmd

	// spinner tick while loading
	if p.Loading {
		var cmd tea.Cmd
		p.spinner, cmd = p.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	if km, ok := msg.(tea.KeyMsg); ok {
		switch km.String() {
		case "tab":
			p.focusNext()
			return p, tea.Batch(cmds...)
		case "shift+tab":
			p.focusPrev()
			return p, tea.Batch(cmds...)
		case "up", "k":
			if p.focus == FocusFilter && p.cursor > 0 {
				p.cursor--
			}
			return p, tea.Batch(cmds...)
		case "down", "j":
			if p.focus == FocusFilter && p.cursor < len(p.filtered)-1 {
				p.cursor++
			}
			return p, tea.Batch(cmds...)
		case "enter":
			// Enter advances to next field; on last field, trigger search
			if p.focus < FocusBranch {
				p.focusNext()
				return p, tea.Batch(cmds...)
			}
			// Last field → send StartSearchMsg
			return p, func() tea.Msg { return StartSearchMsg{} }
		}
	}

	var cmd tea.Cmd
	switch p.focus {
	case FocusFilter:
		p.Filter, cmd = p.Filter.Update(msg)
		cmds = append(cmds, cmd)
		p.applyFilter()
	case FocusDays:
		p.Days, cmd = p.Days.Update(msg)
		cmds = append(cmds, cmd)
	case FocusBranch:
		p.Branch, cmd = p.Branch.Update(msg)
		cmds = append(cmds, cmd)
	}

	return p, tea.Batch(cmds...)
}

func (p InputPanel) View() string {
	var b strings.Builder

	b.WriteString(style.PanelTitle.Render("Search") + "\n")

	// Repository list
	if p.focus == FocusFilter {
		b.WriteString(style.ActiveLabel.Render("Repository") + "\n")
	} else {
		b.WriteString(style.Label.Render("Repository") + "\n")
	}
	b.WriteString(p.Filter.View() + "\n")

	if p.Loading {
		b.WriteString(p.spinner.View() + " " + style.StatusText.Render("Loading repositories...") + "\n")
	} else {
		maxVisible := p.Height - 10
		if maxVisible < 3 {
			maxVisible = 3
		}

		start := 0
		if p.cursor >= maxVisible {
			start = p.cursor - maxVisible + 1
		}

		for i := start; i < len(p.filtered) && i < start+maxVisible; i++ {
			if i == p.cursor {
				b.WriteString(style.CursorStyle.Render("> ") + style.SelectedItem.Render(p.filtered[i]) + "\n")
			} else {
				b.WriteString("  " + style.UnselectedItem.Render(p.filtered[i]) + "\n")
			}
		}
		if len(p.filtered) == 0 && len(p.Repos) > 0 {
			b.WriteString(style.StatusText.Render("  (no match)") + "\n")
		}
		if len(p.filtered) > 0 {
			b.WriteString(style.StatusText.Render(fmt.Sprintf("  %d/%d repos", len(p.filtered), len(p.Repos))) + "\n")
		}
	}

	b.WriteString("\n")

	// Days
	if p.focus == FocusDays {
		b.WriteString(style.ActiveLabel.Render("Days    ") + p.Days.View() + "\n")
	} else {
		b.WriteString(style.Label.Render("Days    ") + p.Days.View() + "\n")
	}

	// Branch
	if p.focus == FocusBranch {
		b.WriteString(style.ActiveLabel.Render("Branch  ") + p.Branch.View() + "\n")
	} else {
		b.WriteString(style.Label.Render("Branch  ") + p.Branch.View() + "\n")
	}

	b.WriteString("\n")
	b.WriteString(style.HelpStyle.Render("Enter next  Tab skip  Ctrl+C quit"))

	return b.String()
}
