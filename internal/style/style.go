package style

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	Cyan    = lipgloss.Color("#00BFFF")
	Green   = lipgloss.Color("#00FF7F")
	Yellow  = lipgloss.Color("#FFD700")
	Magenta = lipgloss.Color("#FF69B4")
	Red     = lipgloss.Color("#FF4444")
	Dim     = lipgloss.Color("#666666")
	White   = lipgloss.Color("#FFFFFF")

	// Panel borders
	InputPanel = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Cyan).
			Padding(1, 2)

	OutputPanel = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Magenta).
			Padding(1, 2)

	// Title styles
	PanelTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Cyan)

	// Status
	StatusText = lipgloss.NewStyle().Foreground(Dim)
	ErrorText  = lipgloss.NewStyle().Foreground(Red).Bold(true)
	SuccessText = lipgloss.NewStyle().Foreground(Green)

	// List items
	SelectedItem   = lipgloss.NewStyle().Foreground(Cyan).Bold(true)
	UnselectedItem = lipgloss.NewStyle().Foreground(White)
	CursorStyle    = lipgloss.NewStyle().Foreground(Cyan)

	// Labels
	Label     = lipgloss.NewStyle().Foreground(Yellow).Bold(true)
	HelpStyle = lipgloss.NewStyle().Foreground(Dim)
)
