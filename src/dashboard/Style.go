package dashboard

import "github.com/charmbracelet/lipgloss"

var (
	borderColor = lipgloss.Color("#7D56F4")
	titleColor  = lipgloss.Color("#FF6AC1")

	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(0, 1).
			Margin(0, 1)

	FocusedPanelStyle = PanelStyle.
				BorderForeground(titleColor)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(titleColor).
			MarginBottom(1)
)
