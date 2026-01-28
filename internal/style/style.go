package style

import "github.com/charmbracelet/lipgloss"

var (
	// 2-color palette: soft blue + soft green
	Primary = lipgloss.Color("#82AAFF")
	Accent  = lipgloss.Color("#C3E88D")
	Dim     = lipgloss.Color("#555555")
	Text    = lipgloss.Color("#CDD6F4")
	Red     = lipgloss.Color("#F38BA8")

	// Panel borders — minimal padding
	InputPanel = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(0, 1)

	OutputPanel = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Dim).
			Padding(0, 1)

	// Titles
	PanelTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary)

	// Status
	StatusText  = lipgloss.NewStyle().Foreground(Dim)
	ErrorText   = lipgloss.NewStyle().Foreground(Red)
	SuccessText = lipgloss.NewStyle().Foreground(Accent)

	// List items
	SelectedItem   = lipgloss.NewStyle().Foreground(Primary).Bold(true)
	UnselectedItem = lipgloss.NewStyle().Foreground(Text)
	CursorStyle    = lipgloss.NewStyle().Foreground(Accent)

	// Labels
	Label     = lipgloss.NewStyle().Foreground(Primary)
	HelpStyle = lipgloss.NewStyle().Foreground(Dim)

	// Active field indicator
	ActiveLabel = lipgloss.NewStyle().Foreground(Accent).Bold(true)
)
