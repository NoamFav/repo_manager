package dashboard

import tea "github.com/charmbracelet/bubbletea"

type BranchesModel struct{}

func (m BranchesModel) Init() tea.Cmd {
	return nil
}

func (m BranchesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m BranchesModel) View() string {
	return "BranchesModel"
}
