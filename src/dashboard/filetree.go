package dashboard

import tea "github.com/charmbracelet/bubbletea"

type FiletreeModel struct{}

func (m FiletreeModel) Init() tea.Cmd {
	return nil
}

func (m FiletreeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m FiletreeModel) View() string {
	return "FiletreeModel"
}
