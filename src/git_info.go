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
	if strings.TrimSpace(string(out)) == "" {
		fmt.Println("✅ No changes.")
	}
	return string(out)
}

func GitStatus() string {
	cmd := exec.Command("git", "status", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	status := string(out)
	if strings.TrimSpace(status) == "" {
		fmt.Println("✅ Nothing to commit.")
		status = "Nothing to commit."
	}
	return status
}

func Summary() string {
	return fmt.Sprintf("Git Diff: %s\nGit Status: %s", GitDiff(), GitStatus())
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
