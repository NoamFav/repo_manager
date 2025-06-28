package dashboard

import tea "github.com/charmbracelet/bubbletea"

type IssuesModel struct{}

func (m IssuesModel) Init() tea.Cmd {
	return nil
}

func (m IssuesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m IssuesModel) View() string {
	return "IssuesModel"
}
