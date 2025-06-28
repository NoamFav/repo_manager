package dashboard

import tea "github.com/charmbracelet/bubbletea"

type InfoModel struct{}

func (m InfoModel) Init() tea.Cmd {
	return nil
}

func (m InfoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m InfoModel) View() string {
	return "InfoModel"
}
