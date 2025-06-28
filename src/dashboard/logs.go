package dashboard

import tea "github.com/charmbracelet/bubbletea"

type LogsModel struct{}

func (m LogsModel) Init() tea.Cmd {
	return nil
}

func (m LogsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m LogsModel) View() string {
	return "LogsModel"
}
