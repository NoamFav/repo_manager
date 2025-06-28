package dashboard

import tea "github.com/charmbracelet/bubbletea"

type StatusModel struct{}

func (m StatusModel) Init() tea.Cmd {
	return nil
}

func (m StatusModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m StatusModel) View() string {
	return "StatusModel"
}
