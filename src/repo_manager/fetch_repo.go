package repo_manager

import (
	"encoding/json"
	"os/exec"
)

type LicenseInfo struct {
	Key      string `json:"key"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
}

type Language struct {
	Name string `json:"name"`
}

type PullRequests struct {
	TotalCount int `json:"totalCount"`
}

type Watchers struct {
	TotalCount int `json:"totalCount"`
}

type Repo struct {
	Name            string       `json:"name"`
	Description     string       `json:"description"`
	PrimaryLanguage Language     `json:"primaryLanguage"`
	DiskUsage       int          `json:"diskUsage"`
	ForkCount       int          `json:"forkCount"`
	LicenseInfo     LicenseInfo  `json:"licenseInfo"`
	PullRequests    PullRequests `json:"pullRequests"`
	StargazerCount  int          `json:"stargazerCount"`
	Visibility      string       `json:"visibility"`
	Watchers        Watchers     `json:"watchers"`
}

func FetchGithubRepos() ([]Repo, error) {
	cmd := exec.Command("gh", "repo", "list", "--limit", "100", "--json", "name,description,primaryLanguage,stargazerCount,forkCount,diskUsage,pullRequests,visibility,watchers,licenseInfo")
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
