package dashboard

import tea "github.com/charmbracelet/bubbletea"

type ActionsModel struct{}

func (m ActionsModel) Init() tea.Cmd {
	return nil
}

func (m ActionsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m ActionsModel) View() string {
	return "ActionsModel"
}
