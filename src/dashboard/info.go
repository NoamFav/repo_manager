package dashboard

import (
	"fmt"
	"strings"

	"github.com/NoamFav/Zvezda/src/repo_manager"
	tea "github.com/charmbracelet/bubbletea"
)

type InfoModel struct {
	Info         repo_manager.Info
	windowHeight int
	windowWidth  int
}

func NewInfoModel(repo RepoModel) InfoModel {
	current_repo_name := repo.Repos[repo.Index].Name
	info, err := repo_manager.FetchInfoRepo(current_repo_name)
	if err != nil {
		return InfoModel{}
	}
	return InfoModel{
		Info: info,
	}
}

func (m InfoModel) Init() tea.Cmd {

	return nil
}

func (m InfoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m InfoModel) View() string {
	info := m.Info
	var b strings.Builder

	// Static logo
	b.WriteString(" Zvezda Repository Overview\n\n") // you can replace logo with anything else

	// Project name
	b.WriteString(fmt.Sprintf(" Name: %s\n", info.Name))

	// Head commit
	b.WriteString(fmt.Sprintf(" Head: %s \n", info.Head))

	// Pending changes
	p := info.Pending
	b.WriteString(fmt.Sprintf("󱓞 Pending: %s \n", p))

	// Authors
	var authors []string
	for _, a := range info.Authors {
		authors = append(authors, a)
	}
	b.WriteString(fmt.Sprintf(" Authors: %s\n", strings.Join(authors, ", ")))

	// URL
	b.WriteString(fmt.Sprintf("󰌷 URL: %s\n", info.URL))

	// Commits
	b.WriteString(fmt.Sprintf(" Commits: %d\n", info.Commits))

	// Lines of code
	b.WriteString(fmt.Sprintf("󰯱 LOC: %d\n", info.Lines))

	// Size
	b.WriteString(fmt.Sprintf(" Size: %d\n", info.Size))

	// License
	b.WriteString(fmt.Sprintf("󰿃 License: %s\n", info.License))

	// Last change
	b.WriteString(fmt.Sprintf("󰄉 Last Change: %s\n", info.LastChange))

	return PanelStyle.Render(b.String())
}
