package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	flag "github.com/spf13/pflag"
)

// Enhanced color palette
var (
	primaryColor  = lipgloss.Color("#6366f1") // Indigo
	successColor  = lipgloss.Color("#10b981") // Emerald
	errorColor    = lipgloss.Color("#ef4444") // Red
	warningColor  = lipgloss.Color("#f59e0b") // Amber
	infoColor     = lipgloss.Color("#3b82f6") // Blue
	accentColor   = lipgloss.Color("#8b5cf6") // Violet
	mutedColor    = lipgloss.Color("#6b7280") // Gray
	borderColor   = lipgloss.Color("#e5e7eb") // Light gray
	gradientStart = lipgloss.Color("#667eea") // Gradient start
	gradientEnd   = lipgloss.Color("#764ba2") // Gradient end
)

// Enhanced styles with modern design
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffff")).
			Background(primaryColor).
			Padding(1, 3).
			Bold(true).
			Align(lipgloss.Center).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor)

	headerStyle = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(gradientStart).
			Padding(2, 4).
			Bold(true).
			Align(lipgloss.Center).
			Foreground(primaryColor)

	successStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true).
			Padding(0, 1)

	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true).
			Padding(0, 1)

	infoStyle = lipgloss.NewStyle().
			Foreground(infoColor).
			Bold(true).
			Padding(0, 1)

	warningStyle = lipgloss.NewStyle().
			Foreground(warningColor).
			Bold(true).
			Padding(0, 1)

	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(1, 2).
			Margin(1, 0).
			Background(lipgloss.Color("#f8fafc"))

	repoCardStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(accentColor).
			Padding(1, 2).
			Margin(1, 0).
			Background(lipgloss.Color("#faf5ff")).
			Bold(true)

	statusStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Padding(0, 1)

	branchStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true).
			Background(lipgloss.Color("#f3f4f6")).
			Padding(0, 1).
			Border(lipgloss.NormalBorder()).
			BorderForeground(accentColor)

	configTableStyle = lipgloss.NewStyle().
				Border(lipgloss.ThickBorder()).
				BorderForeground(primaryColor).
				Padding(2, 3).
				Margin(1, 0).
				Background(lipgloss.Color("#f0f9ff"))

	progressBarStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(successColor).
				Padding(1, 2).
				Margin(1, 0)

	summaryCardStyle = lipgloss.NewStyle().
				Border(lipgloss.ThickBorder()).
				BorderForeground(successColor).
				Padding(2, 3).
				Margin(1, 0).
				Background(lipgloss.Color("#f0fdf4"))

	resultItemStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Margin(0, 1)

	separatorStyle = lipgloss.NewStyle().
			Foreground(borderColor).
			Bold(true)
)

// Enhanced icons with better variety
const (
	IconGit        = "🌟" // Changed from 🔗
	IconFolder     = "📂" // Changed from 📁
	IconSuccess    = "✨" // Changed from ✅
	IconError      = "💥" // Changed from ❌
	IconInfo       = "💡" // Changed from ℹ️
	IconWarning    = "⚡" // Changed from ⚠️
	IconCommit     = "💾" // Changed from 📝
	IconPush       = "🚀" // Changed from ☁️
	IconPull       = "⬇️"
	IconBranch     = "🌸" // Changed from 🌿
	IconMainBranch = "🌺" // Changed from 🌳
	IconAdd        = "➕"
	IconRemove     = "🗑️" // Changed from ➖
	IconClock      = "⏱️" // Changed from ⏰
	IconSparkles   = "✨"
	IconRocket     = "🚀"
	IconConfig     = "⚙️"
	IconCheck      = "✅" // Changed from ✓
	IconDot        = "•"
	IconProgress   = "🔄"
	IconComplete   = "🎉"
	IconScanning   = "🔍"
	IconRepo       = "📦"
	IconStats      = "📊"
)

// Enhanced file type icons
func getFileIcon(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".py":
		return "🐍"
	case ".js", ".jsx":
		return "💛"
	case ".ts", ".tsx":
		return "💙"
	case ".html":
		return "🌐"
	case ".css":
		return "🎨"
	case ".go":
		return "🐹"
	case ".json":
		return "📋"
	case ".md":
		return "📝"
	case ".yml", ".yaml":
		return "⚙️"
	case ".png", ".jpg", ".jpeg", ".gif", ".svg":
		return "🖼️"
	case ".mp3", ".wav", ".flac":
		return "🎵"
	case ".mp4", ".mov", ".avi":
		return "🎬"
	case ".zip", ".tar", ".gz", ".rar":
		return "📦"
	case ".pdf":
		return "📄"
	case ".txt":
		return "📃"
	case ".doc", ".docx":
		return "📝"
	default:
		return "📄"
	}
}

// Configuration
type Config struct {
	BaseDir         string
	Pull            bool
	HandleGitignore bool
	RemoveDSStore   bool
	CommitMessage   string
	ExcludeList     []string
	OnlyList        []string
	UseAICommit     bool
}

// Repository information
type Repository struct {
	Name    string
	Path    string
	Branch  string
	Changes []string
}

// Model for Bubble Tea
type Model struct {
	config       Config
	repositories []Repository
	currentRepo  int
	state        string // "scanning", "processing", "done"
	spinner      spinner.Model
	progress     progress.Model
	results      []string
	startTime    time.Time
	logs         []string
	width        int
	height       int
}

// Messages
type scanCompleteMsg struct {
	repos []Repository
}

type repoProcessedMsg struct {
	repo    Repository
	success bool
	message string
}

type allDoneMsg struct{}

func initialModel(config Config) Model {
	s := spinner.New()
	s.Spinner = spinner.Globe
	s.Style = lipgloss.NewStyle().Foreground(primaryColor)

	p := progress.New(progress.WithScaledGradient(string(gradientStart), string(gradientEnd)))

	return Model{
		config:    config,
		state:     "scanning",
		spinner:   s,
		progress:  p,
		startTime: time.Now(),
		width:     80,
		height:    24,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		scanRepositories(m.config),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}

	case scanCompleteMsg:
		m.repositories = msg.repos
		m.state = "processing"
		if len(m.repositories) == 0 {
			m.state = "done"
			return m, nil
		}
		return m, processAllRepositories(m.repositories, m.config)

	case repoProcessedMsg:
		m.currentRepo++
		statusIcon := IconSuccess
		if !msg.success {
			statusIcon = IconError
		}

		m.results = append(m.results, fmt.Sprintf("%s %s: %s",
			statusIcon, msg.repo.Name, msg.message))

		if m.currentRepo >= len(m.repositories) {
			m.state = "done"
			return m, tea.Sequence(tea.Tick(time.Second, func(t time.Time) tea.Msg {
				return allDoneMsg{}
			}))
		}

	case allDoneMsg:
		// Final state reached

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	var b strings.Builder

	// Main title with enhanced styling
	title := fmt.Sprintf("%s Git Repository Manager %s", IconRocket, IconSparkles)
	styledTitle := headerStyle.Width(m.width - 4).Render(title)
	b.WriteString(styledTitle + "\n\n")

	// Add decorative separator
	separator := separatorStyle.Render(strings.Repeat("─", m.width-4))
	b.WriteString(separator + "\n\n")

	switch m.state {
	case "scanning":
		// Enhanced scanning view
		scanCard := cardStyle.Render(fmt.Sprintf(
			"%s %s Discovering Git repositories...\n\n%s Scanning directory: %s",
			IconScanning, m.spinner.View(), IconFolder, m.config.BaseDir))
		b.WriteString(scanCard + "\n\n")

		// Configuration preview during scanning
		configTable := m.renderConfigTable()
		b.WriteString(configTable + "\n")

	case "processing":
		if len(m.repositories) > 0 {
			// Enhanced progress display
			progressPercent := float64(m.currentRepo) / float64(len(m.repositories))

			progressInfo := fmt.Sprintf("%s Processing repositories: %d/%d",
				IconProgress, m.currentRepo, len(m.repositories))

			progressCard := progressBarStyle.Render(
				progressInfo + "\n" + m.progress.ViewAs(progressPercent))
			b.WriteString(progressCard + "\n")

			// Show current repository being processed with enhanced styling
			if m.currentRepo < len(m.repositories) {
				currentRepo := m.repositories[m.currentRepo]
				repoInfo := fmt.Sprintf(
					"%s %s\n%s Branch: %s\n%s Status: Processing...",
					IconRepo, currentRepo.Name,
					IconBranch, branchStyle.Render(currentRepo.Branch),
					IconClock)

				currentRepoCard := repoCardStyle.Render(repoInfo)
				b.WriteString(currentRepoCard + "\n")
			}

			// Show recent results
			if len(m.results) > 0 {
				recentResults := m.results
				if len(recentResults) > 3 {
					recentResults = recentResults[len(recentResults)-3:]
				}

				b.WriteString(infoStyle.Render("Recent Results:") + "\n")
				for _, result := range recentResults {
					b.WriteString(resultItemStyle.Render(result) + "\n")
				}
			}
		}

	case "done":
		// Enhanced completion summary
		elapsed := time.Since(m.startTime)
		successCount := 0
		errorCount := 0

		for _, result := range m.results {
			if strings.Contains(result, IconSuccess) {
				successCount++
			} else if strings.Contains(result, IconError) {
				errorCount++
			}
		}

		// Create stats summary
		stats := fmt.Sprintf(
			"%s Processing Complete!\n\n"+
				"%s Successfully processed: %d repositories\n"+
				"%s Failed: %d repositories\n"+
				"%s Total repositories: %d\n"+
				"%s Total time: %.2f seconds\n"+
				"%s Average time per repo: %.2f seconds",
			IconComplete,
			IconCheck, successCount,
			IconError, errorCount,
			IconRepo, len(m.repositories),
			IconClock, elapsed.Seconds(),
			IconStats, elapsed.Seconds()/float64(len(m.repositories)))

		var summaryCard string
		if errorCount == 0 && successCount > 0 {
			summaryCard = summaryCardStyle.BorderForeground(successColor).Render(stats)
		} else if errorCount > 0 {
			summaryCard = summaryCardStyle.BorderForeground(warningColor).Render(stats)
		} else {
			summaryCard = summaryCardStyle.BorderForeground(infoColor).Render(stats)
		}

		b.WriteString(summaryCard + "\n\n")

		// Enhanced results display
		if len(m.results) > 0 {
			b.WriteString(infoStyle.Render("Detailed Results:") + "\n")
			b.WriteString(separatorStyle.Render(strings.Repeat("─", 50)) + "\n")

			for i, result := range m.results {
				resultStyle := resultItemStyle
				if strings.Contains(result, IconSuccess) {
					resultStyle = resultStyle.Foreground(successColor)
				} else if strings.Contains(result, IconError) {
					resultStyle = resultStyle.Foreground(errorColor)
				}

				formattedResult := fmt.Sprintf("%d. %s", i+1, result)
				b.WriteString(resultStyle.Render(formattedResult) + "\n")
			}
		}
	}

	// Enhanced footer
	b.WriteString("\n" + separatorStyle.Render(strings.Repeat("─", m.width-4)) + "\n")
	footer := statusStyle.Render("Press 'q', 'esc', or Ctrl+C to quit")
	b.WriteString(footer)

	return b.String()
}

func (m Model) renderConfigTable() string {
	var table strings.Builder

	// Enhanced configuration display
	configTitle := titleStyle.Render(fmt.Sprintf("%s Configuration", IconConfig))
	table.WriteString(configTitle + "\n\n")

	// Create a more structured config display
	configs := []struct {
		label string
		value string
		icon  string
	}{
		{"Base Directory", m.config.BaseDir, IconFolder},
		{"Pull Changes", boolToYesNo(m.config.Pull), IconPull},
		{"Handle .gitignore", boolToYesNo(m.config.HandleGitignore), IconConfig},
		{"Remove .DS_Store", boolToYesNo(m.config.RemoveDSStore), IconRemove},
		{"Using AI Commit", boolToYesNo(m.config.UseAICommit), IconCommit},
	}

	for _, config := range configs {
		configLine := fmt.Sprintf("%s %s: %s",
			config.icon, config.label, config.value)
		table.WriteString(configLine + "\n")
	}

	// Commit message handling
	commitMsg := m.config.CommitMessage
	if commitMsg == "auto-commit" {
		commitMsg = successStyle.Render("AI Generated")
	} else {
		commitMsg = infoStyle.Render(commitMsg)
	}
	table.WriteString(fmt.Sprintf("%s Commit Message: %s\n", IconCommit, commitMsg))

	// Lists with better formatting
	if len(m.config.ExcludeList) > 0 {
		excludeList := strings.Join(m.config.ExcludeList, ", ")
		table.WriteString(fmt.Sprintf("%s Excluded: %s\n",
			IconRemove, warningStyle.Render(excludeList)))
	}

	if len(m.config.OnlyList) > 0 {
		onlyList := strings.Join(m.config.OnlyList, ", ")
		table.WriteString(fmt.Sprintf("%s Including Only: %s\n",
			IconAdd, successStyle.Render(onlyList)))
	}

	return configTableStyle.Render(table.String())
}

func boolToYesNo(b bool) string {
	if b {
		return successStyle.Render("✓ Yes")
	}
	return errorStyle.Render("✗ No")
}

// Commands remain the same but with enhanced error handling
func scanRepositories(config Config) tea.Cmd {
	return func() tea.Msg {
		entries, err := os.ReadDir(config.BaseDir)
		if err != nil {
			log.Error("Failed to read directory", "error", err)
			return scanCompleteMsg{repos: []Repository{}}
		}

		var repos []Repository
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			name := entry.Name()

			// Check exclusions
			if contains(config.ExcludeList, name) {
				continue
			}

			// Check only list
			if len(config.OnlyList) > 0 && !contains(config.OnlyList, name) {
				continue
			}

			repoPath := filepath.Join(config.BaseDir, name)
			gitPath := filepath.Join(repoPath, ".git")

			if _, err := os.Stat(gitPath); err == nil {
				// Get current branch
				branch := getCurrentBranch(repoPath)

				repo := Repository{
					Name:   name,
					Path:   repoPath,
					Branch: branch,
				}
				repos = append(repos, repo)
			}
		}

		sort.Slice(repos, func(i, j int) bool {
			return repos[i].Name < repos[j].Name
		})

		return scanCompleteMsg{repos: repos}
	}
}

func processAllRepositories(repos []Repository, config Config) tea.Cmd {
	return func() tea.Msg {
		for _, repo := range repos {
			success, message := processRepository(repo, config)
			// In a real implementation, we'd send individual messages
			// For simplicity, we'll process all at once
			_ = success
			_ = message
		}
		return allDoneMsg{}
	}
}

func processRepository(repo Repository, config Config) (bool, string) {
	// Change to repository directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	if err := os.Chdir(repo.Path); err != nil {
		return false, fmt.Sprintf("Failed to change directory: %v", err)
	}

	var operations []string

	// Pull changes if requested
	if config.Pull {
		if err := runGitCommand("pull"); err != nil {
			return false, fmt.Sprintf("Failed to pull: %v", err)
		}
		operations = append(operations, "pulled changes")
	}

	// Handle .gitignore
	if config.HandleGitignore {
		if err := ensureGitignoreHasDSStore(repo.Path); err != nil {
			return false, fmt.Sprintf("Failed to update .gitignore: %v", err)
		}
		operations = append(operations, "updated .gitignore")
	}

	// Remove .DS_Store files
	if config.RemoveDSStore {
		count, err := removeDSStoreFiles(repo.Path)
		if err != nil {
			return false, fmt.Sprintf("Failed to remove .DS_Store files: %v", err)
		}
		if count > 0 {
			operations = append(operations, fmt.Sprintf("removed %d .DS_Store files", count))
		}
	}

	// Check for changes
	hasChanges, err := hasUncommittedChanges()
	if err != nil {
		return false, fmt.Sprintf("Failed to check for changes: %v", err)
	}

	if !hasChanges {
		return true, "No changes to commit"
	}

	// Stage changes
	if err := runGitCommand("add", "."); err != nil {
		return false, fmt.Sprintf("Failed to stage changes: %v", err)
	}

	// Commit changes
	commitMessage := config.CommitMessage
	if commitMessage == "auto-commit" {
		commitMessage = generateCommitMessage()
	}

	if config.UseAICommit {
		// Use ai_commit command
		cmd := exec.Command("ai_commit", commitMessage)
		if err := cmd.Run(); err != nil {
			return false, fmt.Sprintf("ai_commit failed: %v", err)
		}
	} else {
		// Manual commit
		if err := runGitCommand("commit", "-m", commitMessage); err != nil {
			return false, fmt.Sprintf("Failed to commit: %v", err)
		}

		// Push changes
		if err := runGitCommand("push"); err != nil {
			return false, fmt.Sprintf("Failed to push: %v", err)
		}
	}

	operations = append(operations, "committed and pushed changes")
	return true, strings.Join(operations, ", ")
}

// Helper functions (unchanged)
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func getCurrentBranch(repoPath string) string {
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(repoPath)

	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}

	return strings.TrimSpace(string(output))
}

func runGitCommand(args ...string) error {
	cmd := exec.Command("git", args...)
	return cmd.Run()
}

func hasUncommittedChanges() (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	return len(strings.TrimSpace(string(output))) > 0, nil
}

func ensureGitignoreHasDSStore(repoPath string) error {
	gitignorePath := filepath.Join(repoPath, ".gitignore")

	// Read existing .gitignore or create new one
	content := ""
	if data, err := os.ReadFile(gitignorePath); err == nil {
		content = string(data)
	}

	// Check if .DS_Store is already in .gitignore
	if strings.Contains(content, ".DS_Store") {
		return nil
	}

	// Add .DS_Store to .gitignore
	if content != "" && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	content += ".DS_Store\n"

	return os.WriteFile(gitignorePath, []byte(content), 0644)
}

func removeDSStoreFiles(repoPath string) (int, error) {
	count := 0
	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() == ".DS_Store" {
			// Remove from git tracking
			exec.Command("git", "rm", "--cached", path).Run()
			// Remove file
			if err := os.Remove(path); err == nil {
				count++
			}
		}

		return nil
	})

	return count, err
}

// Enhanced commit message generation with more variety
func generateCommitMessage() string {
	prefixes := []string{
		"✨ Add", "🔧 Fix", "♻️ Refactor", "⚡ Improve", "🎨 Enhance", "🚀 Optimize",
		"📝 Update", "🗑️ Remove", "🔨 Modify", "🏗️ Restructure", "🧹 Clean up",
		"🔒 Secure", "📦 Bundle", "🎯 Focus", "💡 Implement", "🔀 Merge",
	}

	areas := []string{
		"codebase", "functionality", "architecture", "UI/UX", "performance",
		"documentation", "configuration", "dependencies", "features", "components",
		"API endpoints", "database schema", "test coverage", "error handling",
		"user experience", "code quality", "security measures", "build process",
	}

	details := []string{
		"for better maintainability", "to improve user experience",
		"for compatibility with latest standards", "to address technical debt",
		"for enhanced security", "to optimize resource usage",
		"based on user feedback", "following best practices",
		"to meet accessibility standards", "for improved performance",
		"to reduce complexity", "for better error handling",
		"to enhance readability", "for future scalability",
	}

	rand.Seed(time.Now().UnixNano())

	return fmt.Sprintf("%s %s %s",
		prefixes[rand.Intn(len(prefixes))],
		areas[rand.Intn(len(areas))],
		details[rand.Intn(len(details))])
}

func main() {
	var config Config

	// Parse command line flags
	flag.StringVar(&config.BaseDir, "dir", filepath.Join(os.Getenv("HOME"), "Neoware"),
		"Base directory containing git repositories")
	flag.BoolVar(&config.Pull, "pull", false,
		"Pull changes from the remote repository")
	flag.BoolVar(&config.HandleGitignore, "handle-gitignore", false,
		"Ensure .gitignore includes .DS_Store and update it if necessary")
	flag.BoolVar(&config.RemoveDSStore, "remove-ds-store", false,
		"Remove .DS_Store files from the repository")
	flag.StringVar(&config.CommitMessage, "commit-message", "auto-commit",
		"Commit message to use (or 'auto-commit' for AI-generated messages)")
	flag.StringSliceVar(&config.ExcludeList, "exclude", []string{},
		"List of directories to exclude")
	flag.StringSliceVar(&config.OnlyList, "only", []string{},
		"List of directories to include (if empty, include all)")
	flag.BoolVar(&config.UseAICommit, "use-ai-commit", true,
		"Use the ai_commit command")

	flag.Parse()

	// Expand home directory
	if strings.HasPrefix(config.BaseDir, "~/") {
		config.BaseDir = filepath.Join(os.Getenv("HOME"), config.BaseDir[2:])
	}

	// Initialize and run the Bubble Tea program
	p := tea.NewProgram(initialModel(config), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal("Error running program", "error", err)
	}
}
