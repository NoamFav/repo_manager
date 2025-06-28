package dashboard

import (
	"math"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RepoModel struct {
	repos      []Repo
	index      int
	animOffset float64
	animDir    int
	animating  bool
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

	case AnimFrameMsg:
		if m.animating {
			m.animOffset += 0.2
			if m.animOffset >= 1.0 {
				nextIndex := m.index + m.animDir

				// apply only if in bounds
				if nextIndex >= 0 && nextIndex < len(m.repos) {
					m.index = nextIndex
				}

				// stop animating no matter what
				m.animOffset = 0
				m.animDir = 0
				m.animating = false
				return m, nil
			}
			return m, animateTick()
		}

	case tea.KeyMsg:
		if m.animating {
			return m, nil
		}

		switch msg.String() {
		case "up", "k":
			if m.index > 0 {
				m.animDir = -1
				m.animOffset = 0.0
				m.animating = true
				return m, animateTick()
			}
		case "down", "j":
			if m.index < len(m.repos)-1 {
				m.animDir = +1
				m.animOffset = 0.0
				m.animating = true
				return m, animateTick()
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
		relative := float64(i - m.index)
		relative -= float64(m.animDir) * m.animOffset
		scale := scaleForDistance(relative)

		if scale == 0 {
			continue
		}

		style := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Align(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Padding(0, 2).
			Margin(1, 0).
			Width(int(30 * scale)).
			Height(int(5 * scale))

		card := style.Render(TitleStyle.Render(repo.name) + "\n" + repo.desc + "\n" + repo.language)
		centered := lipgloss.PlaceHorizontal(40, lipgloss.Center, card)
		b.WriteString(centered + "\n")
	}

	return PanelStyle.Render(b.String())
}

func animateTick() tea.Cmd {
	return tea.Tick(16*time.Millisecond, func(t time.Time) tea.Msg {
		return AnimFrameMsg{}
	})
}

func scaleForDistance(d float64) float64 {
	abs := math.Abs(d)
	switch {
	case abs <= 1:
		return 0.75 + 0.25*(1-abs) // 1.0 to 0.75
	case abs <= 2:
		return 0.5 + 0.25*(2-abs) // 0.75 to 0.5
	default:
		return 0.0
	}
}
