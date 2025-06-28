package dashboard

import tea "github.com/charmbracelet/bubbletea"

type LanguageModel struct{}

func (m LanguageModel) Init() tea.Cmd {
	return nil
}

func (m LanguageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m LanguageModel) View() string {
	return "LanguageModel"
}
