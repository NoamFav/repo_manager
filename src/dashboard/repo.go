package dashboard

import (
	"encoding/json"
	"github.com/muesli/reflow/wordwrap"
	"os/exec"
	"strings"

	"github.com/NoamFav/Zvezda/src/repo_manager"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RepoModel struct {
	repos        []repo_manager.Repo
	index        int
	windowHeight int
	windowWidth  int
}

type AnimFrameMsg struct{}

func NewRepoModel() RepoModel {
	repos, err := repo_manager.FetchGithubRepos()
	if err != nil || len(repos) == 0 {
		return RepoModel{
			repos: []repo_manager.Repo{
				{
					Name:            "iris",
					Desc:            "AI Assistant",
					PrimaryLanguage: repo_manager.Language{Name: "go/python"},
				},
				{
					Name:            "zvezda",
					Desc:            "Repo Manager",
					PrimaryLanguage: repo_manager.Language{Name: "go"},
				},
				{
					Name:            "enron_classifier",
					Desc:            "NLP Classifier",
					PrimaryLanguage: repo_manager.Language{Name: "python/js"},
				},
				{
					Name:            "shadowedHunter",
					Desc:            "Stealth Game",
					PrimaryLanguage: repo_manager.Language{Name: "C#"},
				},
				{
					Name:            "apple_music",
					Desc:            "neovim Plugin",
					PrimaryLanguage: repo_manager.Language{Name: "lua"},
				},
				{
					Name:            "bitvoyager",
					Desc:            "learning app",
					PrimaryLanguage: repo_manager.Language{Name: "js"},
				},
			},
			index: 2,
		}
	}

	return RepoModel{
		repos: repos,
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

		lang := "N/A"
		if repo.PrimaryLanguage.Name != "" {
			lang = repo.PrimaryLanguage.Name
		}

		card := style.Render(TitleStyle.Render(repo.Name) + "\n" + truncateDesc(repo.Desc, 36, 1) + "\n" + lang)
		centered := lipgloss.PlaceHorizontal(width, lipgloss.Center, card)
		b.WriteString(centered + "\n")
	}

	cardsRendered := b.String()
	centered := lipgloss.PlaceVertical(m.windowHeight, lipgloss.Center, cardsRendered)
	return PanelStyle.Render(centered)
}

func FetchGithubRepos() ([]repo_manager.Repo, error) {
	cmd := exec.Command("gh", "repo", "list", "--limit", "100", "--json", "name,description,language")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var repos []repo_manager.Repo
	err = json.Unmarshal(out, &repos)
	if err != nil {
		return nil, err
	}
	return repos, nil
}

func truncateDesc(desc string, maxWidth int, maxLines int) string {
	wrapped := wordwrap.String(desc, maxWidth)
	lines := strings.Split(wrapped, "\n")
	if len(lines) > maxLines {
		lines = lines[:maxLines]
		lines[maxLines-1] += "..."

	}
	return strings.Join(lines, "\n")
}

func (m RepoModel) WithHeight(h int) RepoModel {
	m.windowHeight = h
	return m
}

func (m RepoModel) WithWidth(w int) RepoModel {
	m.windowWidth = w
	return m
}
