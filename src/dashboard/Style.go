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

	RepoCard100 = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FFD700")).
			Width(50).
			Height(5).
			Padding(0, 2).
			Margin(1, 0).
			Align(lipgloss.Center).
			AlignVertical(lipgloss.Center)

	RepoCard75 = RepoCard100.
			BorderForeground(lipgloss.Color("#AAAAAA")).
			Width(37).
			Height(4)

	RepoCard50 = RepoCard100.
			BorderForeground(lipgloss.Color("#555555")).
			Width(25).
			Height(3)

	RepoHidden = lipgloss.NewStyle().Width(0).Height(0)
)
