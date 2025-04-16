package src

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// GitDiff returns the current diff output
func GitDiff() string {
	cmd := exec.Command("git", "diff")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error getting diff:", err)
	}
	return string(out)
}

// GitStagedDiff returns the diff of staged changes
func GitStagedDiff() string {
	cmd := exec.Command("git", "diff", "--staged")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error getting staged diff:", err)
	}
	return string(out)
}

// GitStatus returns the current status in porcelain format
func GitStatus() string {
	cmd := exec.Command("git", "status", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error getting status:", err)
	}
	return string(out)
}

// GitBranch returns the current branch name
func GitBranch() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error getting branch:", err)
		return ""
	}
	return strings.TrimSpace(string(out))
}

// GitLastCommit returns the last commit message
func GitLastCommit() string {
	cmd := exec.Command("git", "log", "-1", "--pretty=%B")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error getting last commit:", err)
		return ""
	}
	return strings.TrimSpace(string(out))
}

// GitChangedFiles returns a list of files that have been changed
func GitChangedFiles() []string {
	cmd := exec.Command("git", "diff", "--name-only")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error getting changed files:", err)
		return []string{}
	}

	files := strings.Split(strings.TrimSpace(string(out)), "\n")
	// Filter out empty strings
	var result []string
	for _, file := range files {
		if file != "" {
			result = append(result, file)
		}
	}
	return result
}

// GitStagedFiles returns a list of files that have been staged
func GitStagedFiles() []string {
	cmd := exec.Command("git", "diff", "--staged", "--name-only")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error getting staged files:", err)
		return []string{}
	}

	files := strings.Split(strings.TrimSpace(string(out)), "\n")
	// Filter out empty strings
	var result []string
	for _, file := range files {
		if file != "" {
			result = append(result, file)
		}
	}
	return result
}

// ExtractPackageNames attempts to find what packages were modified
func ExtractPackageNames() []string {
	files := append(GitChangedFiles(), GitStagedFiles()...)
	packages := make(map[string]bool)

	for _, file := range files {
		if strings.HasSuffix(file, ".go") {
			parts := strings.Split(file, "/")
			if len(parts) > 1 {
				// Consider the first or second directory as the package
				// This is a heuristic and might need adjustment for your project structure
				pkgIndex := 0
				if len(parts) > 2 && parts[0] == "cmd" || parts[0] == "pkg" || parts[0] == "internal" {
					pkgIndex = 1
				}
				if pkgIndex < len(parts) {
					packages[parts[pkgIndex]] = true
				}
			}
		}
	}

	// Convert map keys to slice
	var result []string
	for pkg := range packages {
		result = append(result, pkg)
	}
	return result
}

// DetectScope tries to intelligently determine the scope for conventional commits
func DetectScope() string {
	packages := ExtractPackageNames()
	if len(packages) == 1 {
		return packages[0]
	} else if len(packages) > 1 {
		return "multi"
	}

	// If we couldn't detect packages, try to determine if this is a specific type of change
	files := append(GitChangedFiles(), GitStagedFiles()...)

	// Check for common patterns
	for _, file := range files {
		if strings.Contains(file, "test") || strings.HasSuffix(file, "_test.go") {
			return "tests"
		}
		if strings.Contains(file, "docs") || strings.HasSuffix(file, ".md") {
			return "docs"
		}
		if strings.Contains(file, "config") || strings.HasSuffix(file, ".yaml") || strings.HasSuffix(file, ".yml") {
			return "config"
		}
	}

	// Default scope based on branch name
	branch := GitBranch()
	scopeRegex := regexp.MustCompile(`(feature|fix|hotfix|chore)/([a-zA-Z0-9_-]+)`)
	matches := scopeRegex.FindStringSubmatch(branch)
	if len(matches) >= 3 {
		return matches[2]
	}

	return ""
}

// DetectType tries to intelligently determine the commit type
func DetectType() string {
	// First check branch name for hints
	branch := GitBranch()
	if strings.HasPrefix(branch, "feature/") {
		return "feat"
	}
	if strings.HasPrefix(branch, "fix/") || strings.HasPrefix(branch, "hotfix/") {
		return "fix"
	}
	if strings.HasPrefix(branch, "chore/") {
		return "chore"
	}

	// Then check files
	files := append(GitChangedFiles(), GitStagedFiles()...)

	// Look for testing changes
	testCount := 0
	docCount := 0
	configCount := 0

	for _, file := range files {
		if strings.Contains(file, "test") || strings.HasSuffix(file, "_test.go") {
			testCount++
		}
		if strings.Contains(file, "docs") || strings.HasSuffix(file, ".md") {
			docCount++
		}
		if strings.Contains(file, "config") || strings.HasSuffix(file, ".yaml") || strings.HasSuffix(file, ".yml") {
			configCount++
		}
	}

	if testCount > 0 && testCount >= len(files)/2 {
		return "test"
	}
	if docCount > 0 && docCount >= len(files)/2 {
		return "docs"
	}
	if configCount > 0 && configCount >= len(files)/2 {
		return "config"
	}

	// Check diff for specific patterns
	diff := GitDiff() + GitStagedDiff()

	if strings.Contains(strings.ToLower(diff), "fix") ||
		strings.Contains(strings.ToLower(diff), "bug") ||
		strings.Contains(strings.ToLower(diff), "issue") {
		return "fix"
	}

	if strings.Contains(strings.ToLower(diff), "refactor") {
		return "refactor"
	}

	// Default to feat if we've added more lines than we've removed
	addedLines := len(regexp.MustCompile(`(?m)^\+`).FindAllString(diff, -1))
	removedLines := len(regexp.MustCompile(`(?m)^-`).FindAllString(diff, -1))

	if addedLines > removedLines {
		return "feat"
	} else if removedLines > addedLines {
		return "refactor"
	}

	return "chore"
}

// Summary provides a comprehensive summary of repository changes
func Summary() string {
	status := GitStatus()
	// Check if there's nothing to commit
	if strings.TrimSpace(status) == "" {
		return ""
	}

	stagedDiff := GitStagedDiff()
	diff := GitDiff()

	if strings.TrimSpace(diff) == "" && strings.TrimSpace(stagedDiff) == "" {
		return ""
	}

	branch := GitBranch()
	changedFiles := GitChangedFiles()
	stagedFiles := GitStagedFiles()

	// Create a more detailed summary
	summary := fmt.Sprintf("Branch: %s\n\n", branch)

	if len(stagedFiles) > 0 {
		summary += "Staged Files:\n"
		for _, file := range stagedFiles {
			summary += fmt.Sprintf("  - %s\n", file)
		}
		summary += "\n"
	}

	if len(changedFiles) > 0 {
		summary += "Unstaged Changed Files:\n"
		for _, file := range changedFiles {
			summary += fmt.Sprintf("  - %s\n", file)
		}
		summary += "\n"
	}

	if strings.TrimSpace(stagedDiff) != "" {
		// Limit the diff size to avoid overwhelming the AI
		if len(stagedDiff) > 2000 {
			summary += fmt.Sprintf("Staged Git Diff (truncated):\n%s...\n\n", stagedDiff[:2000])
		} else {
			summary += fmt.Sprintf("Staged Git Diff:\n%s\n\n", stagedDiff)
		}
	}

	if strings.TrimSpace(diff) != "" {
		// Limit the diff size to avoid overwhelming the AI
		if len(diff) > 2000 {
			summary += fmt.Sprintf("Unstaged Git Diff (truncated):\n%s...\n\n", diff[:2000])
		} else {
			summary += fmt.Sprintf("Unstaged Git Diff:\n%s\n\n", diff)
		}
	}

	// Add suggestions for the commit
	suggestedType := DetectType()
	suggestedScope := DetectScope()

	summary += "Commit Suggestions:\n"
	summary += fmt.Sprintf("  - Type: %s\n", suggestedType)
	summary += fmt.Sprintf("  - Scope: %s\n", suggestedScope)

	return summary
}

// GitAdd stages all changes for commit
func GitAdd() {
	cmd := exec.Command("git", "add", ".")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error adding files:", err)
	}
	fmt.Println(string(out))
}

// GitCommit commits staged changes with the provided message
func GitCommit(message string) {
	cmd := exec.Command("git", "commit", "-m", message)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error committing:", err)
	}
	fmt.Println(string(out))
}

// GitPush pushes commits to the remote repository
func GitPush() {
	cmd := exec.Command("git", "push")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error pushing:", err)
	}
	fmt.Println(string(out))
}

// AddCommitPush performs all three operations in sequence
func AddCommitPush(message string) {
	GitAdd()
	GitCommit(message)
	GitPush()
}

// GenerateCommitPrompt creates an improved prompt for the AI
func GenerateCommitPrompt() string {
	summary := Summary()

	// Detect type and scope to provide better context
	suggestedType := DetectType()
	suggestedScope := DetectScope()

	prompt := fmt.Sprintf(`You are an AI Git assistant. Your task is to write a conventional commit message in the format:
<type>(<scope>): <subject>

I've analyzed the changes and suggest:
- Type: %s (but choose the most appropriate from: feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert)
- Scope: %s (update if you think another scope is more appropriate)

Be concise but descriptive. The subject should:
- Use imperative, present tense (e.g., "change" not "changed" or "changes")
- Not capitalize the first letter
- No period at the end

ONLY return the commit message, nothing else.

Repository changes summary:
%s`, suggestedType, suggestedScope, summary)

	return prompt
}
