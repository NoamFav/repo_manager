package dashboard

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RepoModel struct {
	repos []Repo
	index int
}

type Repo struct {
	name     string
	desc     string
	language string
}

type AnimFrameMsg struct{}

func NewRepoModel() RepoModel {
	return RepoModel{
		repos: []Repo{
			{"iris", "AI Assistant", "go/python"},
			{"zvezda", "Repo Manager", "go"},
			{"enron_classifier", "NLP Classifier", "python/js"},
			{"shadowedHunter", "Stealth Game", "C#"},
			{"apple_music", "neovim Plugin", "lua"},
			{"bitvoyage", "learning app", "js"},
		},
		index: 2,
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
	title := lipgloss.NewStyle().
		Padding(0, 1).
		Render(TitleStyle.Render("Repositories"))

	b.WriteString(title + "\n\n")

	for i := m.index - 2; i <= m.index+2; i++ {
		if i < 0 || i >= len(m.repos) {
			continue
		}

		repo := m.repos[i]
		var style lipgloss.Style

		switch {
		case i == m.index:
			style = RepoCard100
		case i == m.index-1 || i == m.index+1:
			style = RepoCard75
		case i == m.index-2 || i == m.index+2:
			style = RepoCard50
		default:
			style = RepoHidden
		}

		width := 40

		card := style.Render(TitleStyle.Render(repo.name) + "\n" + repo.desc + "\n" + repo.language)
		centered := lipgloss.PlaceHorizontal(width, lipgloss.Center, card)
		b.WriteString(centered + "\n")
	}

	return PanelStyle.Render(b.String())
}
