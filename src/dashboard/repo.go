package dashboard

import (
	"fmt"
	"strings"

	"github.com/muesli/reflow/wordwrap"

	"github.com/NoamFav/Zvezda/src/repo_manager"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RepoModel struct {
	Repos        []repo_manager.Repo
	Index        int
	windowHeight int
	windowWidth  int
}

type AnimFrameMsg struct{}

func NewRepoModel() RepoModel {
	repos, err := repo_manager.FetchGithubRepos()
	if err != nil || len(repos) == 0 {
		return RepoModel{
			Repos: []repo_manager.Repo{
				{
					Name:            "iris",
					Description:     "AI Assistant",
					PrimaryLanguage: repo_manager.Language{Name: "go/python"},
				},
				{
					Name:            "zvezda",
					Description:     "Repo Manager",
					PrimaryLanguage: repo_manager.Language{Name: "go"},
				},
				{
					Name:            "enron_classifier",
					Description:     "NLP Classifier",
					PrimaryLanguage: repo_manager.Language{Name: "python/js"},
				},
				{
					Name:            "shadowedHunter",
					Description:     "Stealth Game",
					PrimaryLanguage: repo_manager.Language{Name: "C#"},
				},
				{
					Name:            "apple_music",
					Description:     "neovim Plugin",
					PrimaryLanguage: repo_manager.Language{Name: "lua"},
				},
				{
					Name:            "bitvoyager",
					Description:     "learning app",
					PrimaryLanguage: repo_manager.Language{Name: "js"},
				},
			},
			Index: 2,
		}
	}

	return RepoModel{
		Repos: repos,
		Index: 2,
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
			if m.Index > 0 {
				m.Index--
			}
		case "down", "j":
			if m.Index < len(m.Repos)-1 {
				m.Index++
			}
		}
	}

	return m, nil
}

func (m RepoModel) View() string {
	var (
		renderedCards     []string
		cardHeights       []int
		repoIndices       []int
		selectedCardIndex = -1
	)

	for i := m.Index - 2; i <= m.Index+2; i++ {
		if i < 0 || i >= len(m.Repos) {
			continue
		}

		repo := m.Repos[i]
		style := RepoHidden
		switch {
		case i == m.Index:
			style = RepoCard100
			selectedCardIndex = len(renderedCards)
		case i == m.Index-1 || i == m.Index+1:
			style = RepoCard75
		case i == m.Index-2 || i == m.Index+2:
			style = RepoCard50
		}

		// Data fallbacks
		lang := "N/A"
		if repo.PrimaryLanguage.Name != "" {
			lang = repo.PrimaryLanguage.Name
		}

		license := "No license"
		if repo.LicenseInfo.Name != "" {
			license = repo.LicenseInfo.Name
		}

		// Icons and metadata
		desc := truncateDesc(repo.Description, 36, 1)
		stars := fmt.Sprintf("󰓎 %d", repo.StargazerCount)
		forks := fmt.Sprintf("󰓂 %d", repo.ForkCount)
		watchers := fmt.Sprintf("  %d", repo.Watchers.TotalCount)

		visIcon := "󰹇"
		if repo.Visibility == "PRIVATE" {
			visIcon = "󱠱"
		}
		vis := fmt.Sprintf("%s %s", visIcon, repo.Visibility)

		metaLine := fmt.Sprintf("󰗀 %s  •  󰿃 %s  •  %s", lang, license, vis)
		countsLine := fmt.Sprintf("%s  %s  %s", stars, forks, watchers)
		title := TitleStyle.Render(" " + repo.Name)

		card := style.Render(fmt.Sprintf("%s\n%s\n%s\n%s", title, desc, metaLine, countsLine))
		centeredCard := lipgloss.PlaceHorizontal(RepoCard100.GetWidth(), lipgloss.Center, card)

		renderedCards = append(renderedCards, centeredCard)
		cardHeights = append(cardHeights, lipgloss.Height(centeredCard))
		repoIndices = append(repoIndices, i)
	}

	if len(renderedCards) == 0 || selectedCardIndex == -1 {
		return PanelStyle.Render("N/A")
	}

	// Vertical centering
	selectedCardHeight := cardHeights[selectedCardIndex]
	centerLine := m.windowHeight / 2
	selectedCardStart := centerLine - (selectedCardHeight / 2) + 1

	finalLines := make([]string, m.windowHeight)
	for i, card := range renderedCards {
		var cardStart int
		if i == selectedCardIndex {
			cardStart = selectedCardStart
		} else {
			spacing := 1
			offset := i - selectedCardIndex
			if offset < 0 {
				cardStart = selectedCardStart
				for j := selectedCardIndex - 1; j >= i; j-- {
					cardStart -= (cardHeights[j] + spacing)
				}
			} else {
				cardStart = selectedCardStart + selectedCardHeight + spacing
				for j := selectedCardIndex + 1; j < i; j++ {
					cardStart += (cardHeights[j] + spacing)
				}
			}
		}

		for j, line := range strings.Split(card, "\n") {
			linePos := cardStart + j
			if linePos >= 0 && linePos < len(finalLines) {
				finalLines[linePos] = line
			}
		}
	}

	var result strings.Builder
	for _, line := range finalLines {
		result.WriteString(line + "\n")
	}
	return result.String()
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
