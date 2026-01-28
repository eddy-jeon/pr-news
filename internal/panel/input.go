package panel

import (
	"strings"

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

type InputPanel struct {
	Repos    []string
	filtered []string
	cursor   int

	Filter textinput.Model
	Days   textinput.Model
	Branch textinput.Model
	focus  FocusField

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

	return InputPanel{
		Filter: filter,
		Days:   days,
		Branch: branch,
		focus:  FocusFilter,
	}
}

func (p *InputPanel) SetRepos(repos []string) {
	p.Repos = repos
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

func (p *InputPanel) FocusNext() {
	p.focus = (p.focus + 1) % FocusFieldCount
	p.updateFocus()
}

func (p *InputPanel) FocusPrev() {
	p.focus = (p.focus - 1 + FocusFieldCount) % FocusFieldCount
	p.updateFocus()
}

func (p *InputPanel) updateFocus() {
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

func (p InputPanel) Update(msg tea.Msg) (InputPanel, tea.Cmd) {
	var cmds []tea.Cmd

	if km, ok := msg.(tea.KeyMsg); ok {
		switch km.String() {
		case "tab":
			p.FocusNext()
			return p, nil
		case "shift+tab":
			p.FocusPrev()
			return p, nil
		case "up", "k":
			if p.focus == FocusFilter && p.cursor > 0 {
				p.cursor--
			}
			return p, nil
		case "down", "j":
			if p.focus == FocusFilter && p.cursor < len(p.filtered)-1 {
				p.cursor++
			}
			return p, nil
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

	title := style.PanelTitle.Render("Search")
	b.WriteString(title + "\n\n")

	// Repository list
	b.WriteString(style.Label.Render("Repository") + "\n")
	b.WriteString(p.Filter.View() + "\n")

	maxVisible := p.Height - 14
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

	b.WriteString("\n")
	b.WriteString(style.Label.Render("Days:   ") + p.Days.View() + "\n")
	b.WriteString(style.Label.Render("Branch: ") + p.Branch.View() + "\n\n")
	b.WriteString(style.HelpStyle.Render("[Enter] Start  [Tab] Next field  [q] Quit"))

	return b.String()
}
