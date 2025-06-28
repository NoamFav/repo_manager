package dashboard

import tea "github.com/charmbracelet/bubbletea"

type CommitsModel struct{}

func (m CommitsModel) Init() tea.Cmd {
	return nil
}

func (m CommitsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m CommitsModel) View() string {
	return "CommitsModel"
}
