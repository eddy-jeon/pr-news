package panel

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/eddy/pr-news/internal/style"
)

type OutputState int

const (
	OutputLoading OutputState = iota
	OutputIdle
	OutputFetching
	OutputSummarizing
	OutputDone
	OutputError
)

type OutputPanel struct {
	State    OutputState
	Status   string
	Progress string
	Content  string
	Error    string

	spinner  spinner.Model
	viewport viewport.Model
	Width    int
	Height   int
	ready    bool
}

func NewOutputPanel() OutputPanel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = style.CursorStyle

	return OutputPanel{
		State:   OutputIdle,
		spinner: s,
	}
}

func (p *OutputPanel) SetSize(w, h int) {
	p.Width = w
	p.Height = h
	if !p.ready {
		p.viewport = viewport.New(w, h-3)
		p.ready = true
	} else {
		p.viewport.Width = w
		p.viewport.Height = h - 3
	}
}

func (p *OutputPanel) SetContent(md string) {
	rendered, err := glamour.Render(md, "dark")
	if err != nil {
		rendered = md
	}
	p.Content = rendered
	if p.ready {
		p.viewport.SetContent(rendered)
	}
}

func (p OutputPanel) Update(msg tea.Msg) (OutputPanel, tea.Cmd) {
	var cmds []tea.Cmd

	if p.State == OutputLoading || p.State == OutputFetching || p.State == OutputSummarizing {
		var cmd tea.Cmd
		p.spinner, cmd = p.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	if p.State == OutputDone && p.ready {
		var cmd tea.Cmd
		p.viewport, cmd = p.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return p, tea.Batch(cmds...)
}

func (p OutputPanel) Init() tea.Cmd {
	return p.spinner.Tick
}

func (p OutputPanel) View() string {
	var b strings.Builder

	b.WriteString(style.PanelTitle.Render("Output") + "\n")

	switch p.State {
	case OutputLoading:
		b.WriteString(p.spinner.View() + " " + style.StatusText.Render("Loading repositories..."))

	case OutputIdle:
		b.WriteString(style.StatusText.Render("Select a repository and press Enter to start."))

	case OutputFetching, OutputSummarizing:
		b.WriteString(p.spinner.View() + " " + p.Status + "\n")
		if p.Progress != "" {
			b.WriteString(style.StatusText.Render(p.Progress))
		}

	case OutputDone:
		if p.ready {
			b.WriteString(p.viewport.View() + "\n")
			b.WriteString(style.HelpStyle.Render(
				fmt.Sprintf("j/k scroll  r restart  %d%%", int(p.viewport.ScrollPercent()*100))))
		}

	case OutputError:
		b.WriteString(style.ErrorText.Render("Error: "+p.Error) + "\n\n")
		b.WriteString(style.HelpStyle.Render("r retry  q quit"))
	}

	return b.String()
}
