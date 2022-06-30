package styles

import "github.com/charmbracelet/lipgloss"

var (
	Primary     = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	Error       = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	Header      = lipgloss.NewStyle().Margin(1, 1, 0, 2)
	Row         = lipgloss.NewStyle().MarginLeft(2)
	SelectedRow = lipgloss.NewStyle().MarginLeft(2).Foreground(lipgloss.Color("205"))
)
