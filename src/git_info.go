package src

import (
	"fmt"
	"os/exec"
	"strings"
)

func GitDiff() string {
	cmd := exec.Command("git", "diff")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}

	return string(out)
}

func GitStatus() string {
	cmd := exec.Command("git", "status", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	return string(out)
}

func Summary() string {
	status := GitStatus()

	// Check if there's nothing to commit
	if strings.TrimSpace(status) == "" {
		return ""
	}

	// Get diff only if we have something to commit
	diff := GitDiff()
	if strings.TrimSpace(diff) == "" {
		return ""
	}

	return fmt.Sprintf("Git Diff: %s\nGit Status: %s", diff, status)
}

func GitAdd() {
	cmd := exec.Command("git", "add", ".")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
}

func GitCommit(message string) {
	cmd := exec.Command("git", "commit", "-m", message)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
}

func GitPush() {
	cmd := exec.Command("git", "push")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
}

func AddCommitPush(message string) {
	GitAdd()
	GitCommit(message)
	GitPush()
}
