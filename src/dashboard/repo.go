package dashboard

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type RepoModel struct {
	repos []string
	index int
}

func NewRepoModel() RepoModel {
	return RepoModel{
		repos: []string{"iris", "zvezda", "enron_classifier", "shadowedHunter"},
		index: 0,
	}
}

func (m RepoModel) Init() tea.Cmd {
	return nil
}

func (m RepoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.index > 0 {
				m.index--
			}
		case "down", "j":
			if m.index < len(m.repos)-1 {
				m.index++
			}
		}
	}
	return m, nil
}

func (m RepoModel) View() string {
	var b strings.Builder

	b.WriteString(TitleStyle.Render("Repos\n\n"))

	for i, repo := range m.repos {
		if i == m.index {
			b.WriteString(FocusedPanelStyle.Render("> " + repo + "\n"))
		} else {
			b.WriteString("  " + repo + "\n") // no per-line border
		}
	}

	return PanelStyle.Render(b.String()) // <- wrap everything at once
}
