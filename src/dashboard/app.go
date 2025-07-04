package dashboard

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MainModel struct {
	windowHeight int
	windowWidth  int
	repoModel    RepoModel
	infoModel    InfoModel
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

	case tea.WindowSizeMsg:
		m.windowHeight = msg.Height
		m.windowWidth = msg.Width
		m.repoModel = m.repoModel.WithHeight(msg.Height).WithWidth(msg.Width)
	}

	updated, cmd := m.repoModel.Update(msg)
	m.repoModel = updated.(RepoModel)
	return m, cmd
}

func (m MainModel) View() string {
	sidebar := m.repoModel.View()
	main_info := m.infoModel.View()

	return lipgloss.JoinHorizontal(lipgloss.Top, sidebar, main_info)
}

func Start() {
	repo := NewRepoModel()
	info := NewInfoModel(repo)
	p := tea.NewProgram(MainModel{repoModel: repo, infoModel: info}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Zvezda crashed:", err)
		os.Exit(1)
	}
}
