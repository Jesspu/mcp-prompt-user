package tui

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	PromptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F8F8F8")).
			Padding(1, 0, 0, 0)

	ContainerStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#7D56F4")).
				Padding(1, 2)

	InstructionsStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Padding(1, 0)
)
