package dashboard

import tea "github.com/charmbracelet/bubbletea"

type SearchModel struct{}

func (m SearchModel) Init() tea.Cmd {
	return nil
}

func (m SearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m SearchModel) View() string {
	return "SearchModel"
}
