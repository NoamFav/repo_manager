package repo_manager

import (
	"encoding/json"
	"os/exec"
)

type Language struct {
	Name string `json:"name"`
}

type Repo struct {
	Name            string   `json:"name"`
	Desc            string   `json:"description"`
	PrimaryLanguage Language `json:"primaryLanguage"`
}

func FetchGithubRepos() ([]Repo, error) {
	cmd := exec.Command("gh", "repo", "list", "--limit", "100", "--json", "name,description,primaryLanguage")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var repos []Repo
	err = json.Unmarshal(out, &repos)
	if err != nil {
		return nil, err
	}

	return repos, nil
}
