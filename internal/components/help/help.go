package help

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	separator = lipgloss.NewStyle().Foreground(lipgloss.Color("#3C3C3C")).Render(" • ")

	// Styles
    colStyle = lipgloss.NewStyle().Width(30)
	helpStyle = lipgloss.NewStyle().Margin(2, 0, 1, 4)
	keyStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))
	descStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#4A4A4A"))
)

type HelpOption struct {
	Key  string
	Help string
}

func View(options []HelpOption) string {
	var lines []string

	var line []string
	for i, help := range options {
		line = append(line, colStyle.Render(keyStyle.Render(help.Key)+" "+descStyle.Render(help.Help)))

		// splits in rows of 3 options max
		if (i+1)%3 == 0 {
			lines = append(lines, strings.Join(line, separator))
			line = []string{}
		}
	}

	// append remainer
	lines = append(lines, strings.Join(line, separator))

	return helpStyle.Render(strings.Join(lines, "\n"))
}
