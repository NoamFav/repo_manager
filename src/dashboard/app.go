package dashboard

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type MainModel struct {
	repoSlide RepoModel
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	updated, cmd := m.repoSlide.Update(msg)
	m.repoSlide = updated.(RepoModel)
	return m, cmd
}

func (m MainModel) View() string {
	left := m.repoSlide.View()

	return left
}

func Start() {
	p := tea.NewProgram(MainModel{repoSlide: NewRepoModel()})
	if _, err := p.Run(); err != nil {
		fmt.Println("Zvezda crashed:", err)
		os.Exit(1)
	}
}
