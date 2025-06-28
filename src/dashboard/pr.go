package dashboard

import tea "github.com/charmbracelet/bubbletea"

type PrModel struct{}

func (m PrModel) Init() tea.Cmd {
	return nil
}

func (m PrModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m PrModel) View() string {
	return "PrModel"
}
