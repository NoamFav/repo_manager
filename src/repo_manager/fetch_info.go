package repo_manager

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type Title struct {
	GitUsername string `json:"gitUsername"`
	GitVersion  string `json:"gitVersion"`
}

type ProjectInfo struct {
	RepoName         string `json:"repoName"`
	NumberOfBranches int    `json:"numberOfBranches"`
	NumberOfTags     int    `json:"numberOfTags"`
}

type HeadRefs struct {
	ShortCommitId string   `json:"shortCommitId"`
	Refs          []string `json:"refs"`
}

type HeadInfo struct {
	HeadRefs HeadRefs `json:"headRefs"`
}

type PendingInfo struct {
	Added    int `json:"added"`
	Deleted  int `json:"deleted"`
	Modified int `json:"modified"`
}

type Author struct {
	Name         string  `json:"name"`
	Email        *string `json:"email"` // null-safe
	NbrOfCommits int     `json:"nbrOfCommits"`
	Contribution int     `json:"contribution"`
}

type AuthorsInfo struct {
	Authors []Author `json:"authors"`
}

type UrlInfo struct {
	RepoUrl string `json:"repoUrl"`
}

type CommitsInfo struct {
	NumberOfCommits int  `json:"numberOfCommits"`
	IsShallow       bool `json:"isShallow"`
}

type LocInfo struct {
	LinesOfCode int `json:"linesOfCode"`
}

type SizeInfo struct {
	RepoSize  string `json:"repoSize"`
	FileCount int    `json:"fileCount"`
}

type LicenseInfoRepo struct {
	License string `json:"license"`
}

type LastChangeInfo struct {
	LastChange string `json:"lastChange"`
}

// Define a tagged union with each field optional
type InfoFields struct {
	ProjectInfo     *ProjectInfo `json:"ProjectInfo,omitempty"`
	DescriptionInfo *struct {
		Description *string `json:"description"`
	} `json:"DescriptionInfo,omitempty"`
	HeadInfo    *HeadInfo    `json:"HeadInfo,omitempty"`
	PendingInfo *PendingInfo `json:"PendingInfo,omitempty"`
	VersionInfo *struct {
		Version string `json:"version"`
	} `json:"VersionInfo,omitempty"`
	CreatedInfo *struct {
		CreationDate string `json:"creationDate"`
	} `json:"CreatedInfo,omitempty"`
	AuthorsInfo     *AuthorsInfo     `json:"AuthorsInfo,omitempty"`
	UrlInfo         *UrlInfo         `json:"UrlInfo,omitempty"`
	CommitsInfo     *CommitsInfo     `json:"CommitsInfo,omitempty"`
	LocInfo         *LocInfo         `json:"LocInfo,omitempty"`
	SizeInfo        *SizeInfo        `json:"SizeInfo,omitempty"`
	LicenseInfoRepo *LicenseInfoRepo `json:"LicenseInfo,omitempty"`
	LastChangeInfo  *LastChangeInfo  `json:"LastChangeInfo,omitempty"`
}

type Root struct {
	Title      Title        `json:"title"`
	InfoFields []InfoFields `json:"infoFields"`
}

type Info struct {
	Logo       string
	Name       string
	Head       string
	Pending    string
	Authors    []string
	URL        string
	Commits    int
	Lines      int
	Size       int
	License    string
	LastChange string
}

func FetchInfoRepo(repo string) (Info, error) {
	cmd := exec.Command("onefetch", repo, "-o", "json")
	out, err := cmd.Output()
	if err != nil {
		return Info{}, err
	}

	var root Root
	err = json.Unmarshal(out, &root)
	if err != nil {
		return Info{}, err
	}

	var result Info
	result.Logo = root.Title.GitUsername // thats entirely fucked for sure

	for _, field := range root.InfoFields {
		if field.ProjectInfo != nil {
			result.Name = field.ProjectInfo.RepoName
		}
		if field.HeadInfo != nil {
			result.Head = field.HeadInfo.HeadRefs.ShortCommitId
		}
		if field.PendingInfo != nil {
			result.Pending = fmt.Sprintf("added %d, modified %d, deleted %d",
				field.PendingInfo.Added, field.PendingInfo.Modified, field.PendingInfo.Deleted)
		}
		if field.AuthorsInfo != nil {
			for _, a := range field.AuthorsInfo.Authors {
				result.Authors = append(result.Authors, a.Name)
			}
		}
		if field.UrlInfo != nil {
			result.URL = field.UrlInfo.RepoUrl
		}
		if field.CommitsInfo != nil {
			result.Commits = field.CommitsInfo.NumberOfCommits
		}
		if field.LocInfo != nil {
			result.Lines = field.LocInfo.LinesOfCode
		}
		if field.SizeInfo != nil {
			sizeParts := strings.Fields(field.SizeInfo.RepoSize)
			if len(sizeParts) > 0 {
				if sizeVal, err := strconv.ParseFloat(sizeParts[0], 64); err == nil {
					result.Size = int(sizeVal)
				}
			}
		}
		if field.LicenseInfoRepo != nil {
			result.License = field.LicenseInfoRepo.License
		}
		if field.LastChangeInfo != nil {
			result.LastChange = field.LastChangeInfo.LastChange
		}
	}

	return result, nil
}
