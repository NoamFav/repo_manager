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

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Background(lipgloss.Color("57")).
			Padding(0, 1).
			Bold(true)

	headerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("86")).
			Padding(1, 2).
			Bold(true).
			Align(lipgloss.Center)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("46")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Bold(true)

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("226")).
			Bold(true)

	repoStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("57")).
			Padding(0, 1).
			Margin(1, 0)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")).
			Italic(true)

	branchStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("211")).
			Bold(true)

	configTableStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("86")).
				Padding(1, 2).
				Margin(1, 0)
)

// Icons
const (
	IconGit        = "🔗"
	IconFolder     = "📁"
	IconSuccess    = "✅"
	IconError      = "❌"
	IconInfo       = "ℹ️"
	IconWarning    = "⚠️"
	IconCommit     = "📝"
	IconPush       = "☁️"
	IconPull       = "⬇️"
	IconBranch     = "🌿"
	IconMainBranch = "🌳"
	IconAdd        = "➕"
	IconRemove     = "➖"
	IconClock      = "⏰"
	IconSparkles   = "✨"
	IconRocket     = "🚀"
	IconConfig     = "⚙️"
	IconCheck      = "✓"
	IconDot        = "•"
)

// File type icons
func getFileIcon(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".py":
		return "🐍"
	case ".js", ".jsx", ".ts", ".tsx":
		return "📜"
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
	case ".png", ".jpg", ".jpeg", ".gif":
		return "🖼️"
	case ".mp3", ".wav":
		return "🎵"
	case ".mp4", ".mov":
		return "🎬"
	case ".zip", ".tar", ".gz":
		return "📦"
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
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	p := progress.New(progress.WithDefaultGradient())

	return Model{
		config:    config,
		state:     "scanning",
		spinner:   s,
		progress:  p,
		startTime: time.Now(),
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
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case scanCompleteMsg:
		m.repositories = msg.repos
		m.state = "processing"
		if len(m.repositories) == 0 {
			m.state = "done"
			return m, nil
		}
		// Start processing the first repository
		return m, processNextRepository(m.repositories[0], m.config)

	case repoProcessedMsg:
		// Add result to our list
		if msg.success {
			m.results = append(m.results, fmt.Sprintf("%s %s: %s", 
				IconSuccess, msg.repo.Name, msg.message))
		} else {
			m.results = append(m.results, fmt.Sprintf("%s %s: %s", 
				IconError, msg.repo.Name, msg.message))
		}
		
		m.currentRepo++
		
		// Check if we have more repositories to process
		if m.currentRepo >= len(m.repositories) {
			m.state = "done"
			return m, tea.Sequence(tea.Tick(time.Second, func(t time.Time) tea.Msg {
				return allDoneMsg{}
			}))
		}
		
		// Process next repository
		nextRepo := m.repositories[m.currentRepo]
		return m, processNextRepository(nextRepo, m.config)

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

	// Header
	header := headerStyle.Render(fmt.Sprintf("%s Git Repository Manager %s", 
		IconRocket, IconSparkles))
	b.WriteString(header + "\n\n")

	// Configuration table
	if m.state == "scanning" || m.state == "processing" {
		configTable := m.renderConfigTable()
		b.WriteString(configTable + "\n\n")
	}

	switch m.state {
	case "scanning":
		b.WriteString(fmt.Sprintf("%s Scanning for Git repositories...\n", 
			m.spinner.View()))

	case "processing":
		if len(m.repositories) > 0 {
			progressPercent := float64(m.currentRepo) / float64(len(m.repositories))
			b.WriteString(fmt.Sprintf("Processing repositories: %d/%d\n", 
				m.currentRepo, len(m.repositories)))
			b.WriteString(m.progress.ViewAs(progressPercent) + "\n\n")

			// Show current repository being processed
			if m.currentRepo < len(m.repositories) {
				currentRepo := m.repositories[m.currentRepo]
				repoInfo := repoStyle.Render(fmt.Sprintf("%s %s\n%s Branch: %s", 
					IconFolder, currentRepo.Name, 
					IconBranch, branchStyle.Render(currentRepo.Branch)))
				b.WriteString(repoInfo + "\n")
			}
			
			// Show completed results so far
			if len(m.results) > 0 {
				b.WriteString("\nCompleted:\n")
				for _, result := range m.results {
					b.WriteString(result + "\n")
				}
			}
		}

	case "done":
		// Summary
		elapsed := time.Since(m.startTime)
		successCount := 0
		for _, result := range m.results {
			if strings.Contains(result, IconSuccess) {
				successCount++
			}
		}

		summary := fmt.Sprintf("%s Processing Complete!\n", IconSparkles)
		summary += fmt.Sprintf("Successfully processed: %d/%d repositories\n", 
			successCount, len(m.repositories))
		summary += fmt.Sprintf("%s Total time: %.2f seconds\n", 
			IconClock, elapsed.Seconds())

		if successCount == len(m.repositories) {
			b.WriteString(successStyle.Render(summary) + "\n\n")
		} else {
			b.WriteString(warningStyle.Render(summary) + "\n\n")
		}

		// Show results
		for _, result := range m.results {
			b.WriteString(result + "\n")
		}
	}

	b.WriteString("\n" + statusStyle.Render("Press 'q' or Ctrl+C to quit"))

	return b.String()
}

func (m Model) renderConfigTable() string {
	var table strings.Builder

	table.WriteString(titleStyle.Render("Configuration") + "\n")
	table.WriteString(fmt.Sprintf("Base Directory: %s\n", m.config.BaseDir))
	table.WriteString(fmt.Sprintf("Pull Changes: %s\n", boolToYesNo(m.config.Pull)))
	table.WriteString(fmt.Sprintf("Handle .gitignore: %s\n", boolToYesNo(m.config.HandleGitignore)))
	table.WriteString(fmt.Sprintf("Remove .DS_Store: %s\n", boolToYesNo(m.config.RemoveDSStore)))
	table.WriteString(fmt.Sprintf("Using AI Commit: %s\n", boolToYesNo(m.config.UseAICommit)))
	
	commitMsg := m.config.CommitMessage
	if commitMsg == "auto-commit" {
		commitMsg = "AI Generated"
	}
	table.WriteString(fmt.Sprintf("Commit Message: %s\n", commitMsg))

	if len(m.config.ExcludeList) > 0 {
		table.WriteString(fmt.Sprintf("Excluded: %s\n", strings.Join(m.config.ExcludeList, ", ")))
	}
	if len(m.config.OnlyList) > 0 {
		table.WriteString(fmt.Sprintf("Including Only: %s\n", strings.Join(m.config.OnlyList, ", ")))
	}

	return configTableStyle.Render(table.String())
}

func boolToYesNo(b bool) string {
	if b {
		return successStyle.Render("Yes")
	}
	return errorStyle.Render("No")
}

// Commands
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

// New function to process repositories one at a time
func processNextRepository(repo Repository, config Config) tea.Cmd {
	return func() tea.Msg {
		success, message := processRepository(repo, config)
		return repoProcessedMsg{
			repo:    repo,
			success: success,
			message: message,
		}
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

// Helper functions
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

func generateCommitMessage() string {
	prefixes := []string{
		"Update", "Enhance", "Fix", "Refactor", "Improve", "Optimize",
		"Add", "Remove", "Modify", "Restructure", "Clean up",
	}
	
	areas := []string{
		"codebase", "functionality", "structure", "design", "performance",
		"documentation", "configuration", "dependencies", "features", "UI",
	}
	
	details := []string{
		"for better maintainability", "to improve user experience",
		"for compatibility with latest standards", "to address technical debt",
		"for enhanced security", "to optimize resource usage",
		"based on feedback", "following best practices",
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