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
			Foreground(lipgloss.Color("#89b4fa")).
			Background(lipgloss.Color("#313244")).
			Padding(0, 2).
			Bold(true).
			MarginBottom(1)

	headerStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#89b4fa")).
			Padding(1, 3).
			Bold(true).
			Align(lipgloss.Center).
			Foreground(lipgloss.Color("#cdd6f4"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#a6e3a1")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f38ba8")).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#74c7ec")).
			Bold(true)

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f9e2af")).
			Bold(true)

	repoStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#89b4fa")).
			Padding(1, 2).
			Margin(1, 0).
			Background(lipgloss.Color("#1e1e2e"))

	currentRepoStyle = lipgloss.NewStyle().
				Border(lipgloss.ThickBorder()).
				BorderForeground(lipgloss.Color("#fab387")).
				Padding(1, 2).
				Margin(1, 0).
				Background(lipgloss.Color("#313244"))

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6c7086")).
			Italic(true)

	branchStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#cba6f7")).
			Bold(true)

	configTableStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#89b4fa")).
				Padding(1, 2).
				Margin(1, 0).
				Background(lipgloss.Color("#181825"))

	logStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#585b70")).
			Padding(1, 2).
			Margin(1, 0).
			Background(lipgloss.Color("#1e1e2e")).
			Height(8)

	operationStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#94e2d5")).
			PaddingLeft(2)

	progressBarStyle = lipgloss.NewStyle().
				Margin(1, 0)
)

// Nerd Font Icons
const (
	IconGit        = ""   // Git branch
	IconFolder     = ""   // Folder
	IconSuccess    = ""   // Check circle
	IconError      = ""   // X circle
	IconInfo       = ""   // Info circle
	IconWarning    = ""   // Warning triangle
	IconCommit     = ""   // Git commit
	IconPush       = ""   // Upload
	IconPull       = ""   // Download
	IconBranch     = ""   // Git branch
	IconMainBranch = ""   // Tree
	IconAdd        = ""   // Plus
	IconRemove     = ""   // Minus
	IconClock      = ""   // Clock
	IconSparkles   = ""   // Star
	IconRocket     = ""   // Rocket
	IconConfig     = ""   // Gear
	IconCheck      = ""   // Check
	IconDot        = ""   // Dot
	IconSync       = ""   // Sync
	IconFile       = ""   // File
	IconTerminal   = ""   // Terminal
	IconProcess    = ""   // Process
)

// File type icons using Nerd Fonts
func getFileIcon(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".py":
		return ""
	case ".js", ".jsx":
		return ""
	case ".ts", ".tsx":
		return ""
	case ".html":
		return ""
	case ".css":
		return ""
	case ".go":
		return ""
	case ".json":
		return ""
	case ".md":
		return ""
	case ".yml", ".yaml":
		return ""
	case ".png", ".jpg", ".jpeg", ".gif":
		return ""
	case ".mp3", ".wav":
		return ""
	case ".mp4", ".mov":
		return ""
	case ".zip", ".tar", ".gz":
		return ""
	default:
		return IconFile
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

// Operation log entry
type LogEntry struct {
	Timestamp time.Time
	Level     string
	Repo      string
	Message   string
	Icon      string
}

// Repository information
type Repository struct {
	Name    string
	Path    string
	Branch  string
	Changes []string
	Status  string
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
	logs         []LogEntry
	currentOps   []string // Current operations being performed
}

// Messages
type scanCompleteMsg struct {
	repos []Repository
}

type repoProcessedMsg struct {
	repo       Repository
	success    bool
	message    string
	operations []string
	logs       []LogEntry
}

type operationUpdateMsg struct {
	repo      string
	operation string
	success   bool
}

type allDoneMsg struct{}

func initialModel(config Config) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#f38ba8"))

	p := progress.New(progress.WithScaledGradient("#f38ba8", "#a6e3a1"))

	return Model{
		config:    config,
		state:     "scanning",
		spinner:   s,
		progress:  p,
		startTime: time.Now(),
		logs:      []LogEntry{},
	}
}

func (m *Model) addLog(level, repo, message, icon string) {
	m.logs = append(m.logs, LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Repo:      repo,
		Message:   message,
		Icon:      icon,
	})
	
	// Keep only last 20 log entries
	if len(m.logs) > 20 {
		m.logs = m.logs[len(m.logs)-20:]
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
		m.addLog("INFO", "SYSTEM", fmt.Sprintf("Found %d repositories to process", len(msg.repos)), IconInfo)
		
		if len(m.repositories) == 0 {
			m.state = "done"
			return m, nil
		}
		// Start processing the first repository
		return m, processNextRepository(m.repositories[0], m.config)

	case repoProcessedMsg:
		// Add logs from processing
		for _, logEntry := range msg.logs {
			m.logs = append(m.logs, logEntry)
		}
		
		// Add result to our list
		if msg.success {
			m.results = append(m.results, fmt.Sprintf("%s %s: %s", 
				IconSuccess, msg.repo.Name, msg.message))
			m.addLog("SUCCESS", msg.repo.Name, msg.message, IconSuccess)
		} else {
			m.results = append(m.results, fmt.Sprintf("%s %s: %s", 
				IconError, msg.repo.Name, msg.message))
			m.addLog("ERROR", msg.repo.Name, msg.message, IconError)
		}
		
		m.currentRepo++
		
		// Check if we have more repositories to process
		if m.currentRepo >= len(m.repositories) {
			m.state = "done"
			m.addLog("INFO", "SYSTEM", "All repositories processed", IconSparkles)
			return m, tea.Sequence(tea.Tick(time.Second, func(t time.Time) tea.Msg {
				return allDoneMsg{}
			}))
		}
		
		// Process next repository
		nextRepo := m.repositories[m.currentRepo]
		return m, processNextRepository(nextRepo, m.config)

	case operationUpdateMsg:
		if msg.success {
			m.addLog("INFO", msg.repo, msg.operation, IconCheck)
		} else {
			m.addLog("ERROR", msg.repo, msg.operation, IconError)
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

	// Header with enhanced styling
	header := headerStyle.Render(fmt.Sprintf("%s Git Repository Manager %s", 
		IconRocket, IconSparkles))
	b.WriteString(header + "\n\n")

	// Configuration table
	if m.state == "scanning" || m.state == "processing" {
		configTable := m.renderConfigTable()
		b.WriteString(configTable + "\n")
	}

	switch m.state {
	case "scanning":
		scanningMsg := fmt.Sprintf("%s %s Scanning for Git repositories...", 
			m.spinner.View(), IconFolder)
		b.WriteString(infoStyle.Render(scanningMsg) + "\n\n")

	case "processing":
		if len(m.repositories) > 0 {
			progressPercent := float64(m.currentRepo) / float64(len(m.repositories))
			
			// Progress section
			progressHeader := titleStyle.Render(fmt.Sprintf("%s Processing Progress", IconProcess))
			b.WriteString(progressHeader + "\n")
			
			progressInfo := fmt.Sprintf("Repository %d of %d", 
				m.currentRepo+1, len(m.repositories))
			b.WriteString(infoStyle.Render(progressInfo) + "\n")
			
			progressBar := progressBarStyle.Render(m.progress.ViewAs(progressPercent))
			b.WriteString(progressBar + "\n")

			// Current repository being processed
			if m.currentRepo < len(m.repositories) {
				currentRepo := m.repositories[m.currentRepo]
				repoHeader := fmt.Sprintf("%s Currently Processing", IconTerminal)
				repoInfo := fmt.Sprintf("%s %s\n%s Branch: %s\n%s Status: Processing...", 
					IconFolder, currentRepo.Name, 
					IconBranch, branchStyle.Render(currentRepo.Branch),
					IconSync)
				
				currentRepoBox := currentRepoStyle.Render(repoHeader + "\n" + repoInfo)
				b.WriteString(currentRepoBox + "\n")
			}
			
			// Recent logs
			if len(m.logs) > 0 {
				logsSection := m.renderLogs()
				b.WriteString(logsSection + "\n")
			}
			
			// Show completed results
			if len(m.results) > 0 {
				completedHeader := titleStyle.Render(fmt.Sprintf("%s Completed Repositories", IconCheck))
				b.WriteString(completedHeader + "\n")
				
				// Show only last 5 completed results to save space
				start := 0
				if len(m.results) > 5 {
					start = len(m.results) - 5
				}
				
				for i := start; i < len(m.results); i++ {
					result := m.results[i]
					b.WriteString(operationStyle.Render(result) + "\n")
				}
				
				if len(m.results) > 5 {
					moreCount := len(m.results) - 5
					b.WriteString(statusStyle.Render(fmt.Sprintf("... and %d more", moreCount)) + "\n")
				}
			}
		}

	case "done":
		// Summary with enhanced styling
		elapsed := time.Since(m.startTime)
		successCount := 0
		for _, result := range m.results {
			if strings.Contains(result, IconSuccess) {
				successCount++
			}
		}

		summaryHeader := titleStyle.Render(fmt.Sprintf("%s Processing Complete!", IconSparkles))
		b.WriteString(summaryHeader + "\n")

		stats := []string{
			fmt.Sprintf("%s Successfully processed: %d/%d repositories", 
				IconSuccess, successCount, len(m.repositories)),
			fmt.Sprintf("%s Total time: %.2f seconds", 
				IconClock, elapsed.Seconds()),
			fmt.Sprintf("%s Total operations logged: %d", 
				IconTerminal, len(m.logs)),
		}

		for _, stat := range stats {
			if successCount == len(m.repositories) {
				b.WriteString(successStyle.Render(stat) + "\n")
			} else {
				b.WriteString(warningStyle.Render(stat) + "\n")
			}
		}
		b.WriteString("\n")

		// Final results
		resultsHeader := titleStyle.Render(fmt.Sprintf("%s Final Results", IconCheck))
		b.WriteString(resultsHeader + "\n")
		
		for _, result := range m.results {
			b.WriteString(operationStyle.Render(result) + "\n")
		}

		// Final logs
		if len(m.logs) > 0 {
			finalLogsSection := m.renderLogs()
			b.WriteString("\n" + finalLogsSection)
		}
	}

	b.WriteString("\n" + statusStyle.Render(fmt.Sprintf("%s Press 'q' or Ctrl+C to quit", IconInfo)))

	return b.String()
}

func (m Model) renderConfigTable() string {
	var table strings.Builder

	configHeader := titleStyle.Render(fmt.Sprintf("%s Configuration", IconConfig))
	table.WriteString(configHeader + "\n")
	
	configItems := []string{
		fmt.Sprintf("%s Base Directory: %s", IconFolder, m.config.BaseDir),
		fmt.Sprintf("%s Pull Changes: %s", IconPull, boolToYesNo(m.config.Pull)),
		fmt.Sprintf("%s Handle .gitignore: %s", IconFile, boolToYesNo(m.config.HandleGitignore)),
		fmt.Sprintf("%s Remove .DS_Store: %s", IconRemove, boolToYesNo(m.config.RemoveDSStore)),
		fmt.Sprintf("%s Using AI Commit: %s", IconSparkles, boolToYesNo(m.config.UseAICommit)),
	}
	
	commitMsg := m.config.CommitMessage
	if commitMsg == "auto-commit" {
		commitMsg = "AI Generated"
	}
	configItems = append(configItems, fmt.Sprintf("%s Commit Message: %s", IconCommit, commitMsg))

	if len(m.config.ExcludeList) > 0 {
		configItems = append(configItems, fmt.Sprintf("%s Excluded: %s", IconRemove, strings.Join(m.config.ExcludeList, ", ")))
	}
	if len(m.config.OnlyList) > 0 {
		configItems = append(configItems, fmt.Sprintf("%s Including Only: %s", IconAdd, strings.Join(m.config.OnlyList, ", ")))
	}

	for _, item := range configItems {
		table.WriteString(item + "\n")
	}

	return configTableStyle.Render(table.String())
}

func (m Model) renderLogs() string {
	if len(m.logs) == 0 {
		return ""
	}

	var logContent strings.Builder
	logHeader := titleStyle.Render(fmt.Sprintf("%s Recent Activity", IconTerminal))
	logContent.WriteString(logHeader + "\n")

	// Show last 10 logs
	start := 0
	if len(m.logs) > 10 {
		start = len(m.logs) - 10
	}

	for i := start; i < len(m.logs); i++ {
		entry := m.logs[i]
		timestamp := entry.Timestamp.Format("15:04:05")
		
		var style lipgloss.Style
		switch entry.Level {
		case "SUCCESS":
			style = successStyle
		case "ERROR":
			style = errorStyle
		case "WARNING":
			style = warningStyle
		default:
			style = infoStyle
		}
		
		logLine := fmt.Sprintf("[%s] %s %s: %s", 
			timestamp, entry.Icon, entry.Repo, entry.Message)
		logContent.WriteString(style.Render(logLine) + "\n")
	}

	return logStyle.Render(logContent.String())
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
					Status: "pending",
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
		success, message, operations, logs := processRepositoryWithLogs(repo, config)
		return repoProcessedMsg{
			repo:       repo,
			success:    success,
			message:    message,
			operations: operations,
			logs:       logs,
		}
	}
}

func processRepositoryWithLogs(repo Repository, config Config) (bool, string, []string, []LogEntry) {
	var logs []LogEntry
	var operations []string
	
	addLog := func(level, message, icon string) {
		logs = append(logs, LogEntry{
			Timestamp: time.Now(),
			Level:     level,
			Repo:      repo.Name,
			Message:   message,
			Icon:      icon,
		})
	}
	
	addLog("INFO", "Starting repository processing", IconProcess)
	
	// Change to repository directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	if err := os.Chdir(repo.Path); err != nil {
		addLog("ERROR", fmt.Sprintf("Failed to change directory: %v", err), IconError)
		return false, fmt.Sprintf("Failed to change directory: %v", err), operations, logs
	}
	
	addLog("INFO", fmt.Sprintf("Changed to directory: %s", repo.Path), IconFolder)

	// Pull changes if requested
	if config.Pull {
		addLog("INFO", "Pulling changes from remote", IconPull)
		if err := runGitCommand("pull"); err != nil {
			addLog("ERROR", fmt.Sprintf("Failed to pull: %v", err), IconError)
			return false, fmt.Sprintf("Failed to pull: %v", err), operations, logs
		}
		operations = append(operations, "pulled changes")
		addLog("SUCCESS", "Successfully pulled changes", IconSuccess)
	}

	// Handle .gitignore
	if config.HandleGitignore {
		addLog("INFO", "Updating .gitignore file", IconFile)
		if err := ensureGitignoreHasDSStore(repo.Path); err != nil {
			addLog("ERROR", fmt.Sprintf("Failed to update .gitignore: %v", err), IconError)
			return false, fmt.Sprintf("Failed to update .gitignore: %v", err), operations, logs
		}
		operations = append(operations, "updated .gitignore")
		addLog("SUCCESS", "Successfully updated .gitignore", IconSuccess)
	}

	// Remove .DS_Store files
	if config.RemoveDSStore {
		addLog("INFO", "Searching for .DS_Store files", IconRemove)
		count, err := removeDSStoreFiles(repo.Path)
		if err != nil {
			addLog("ERROR", fmt.Sprintf("Failed to remove .DS_Store files: %v", err), IconError)
			return false, fmt.Sprintf("Failed to remove .DS_Store files: %v", err), operations, logs
		}
		if count > 0 {
			operations = append(operations, fmt.Sprintf("removed %d .DS_Store files", count))
			addLog("SUCCESS", fmt.Sprintf("Removed %d .DS_Store files", count), IconSuccess)
		} else {
			addLog("INFO", "No .DS_Store files found", IconCheck)
		}
	}

	// Check for changes
	addLog("INFO", "Checking for uncommitted changes", IconSync)
	hasChanges, err := hasUncommittedChanges()
	if err != nil {
		addLog("ERROR", fmt.Sprintf("Failed to check for changes: %v", err), IconError)
		return false, fmt.Sprintf("Failed to check for changes: %v", err), operations, logs
	}

	if !hasChanges {
		addLog("INFO", "No changes to commit", IconCheck)
		return true, "No changes to commit", operations, logs
	}

	addLog("INFO", "Found uncommitted changes", IconCommit)

	// Stage changes
	addLog("INFO", "Staging changes", IconAdd)
	if err := runGitCommand("add", "."); err != nil {
		addLog("ERROR", fmt.Sprintf("Failed to stage changes: %v", err), IconError)
		return false, fmt.Sprintf("Failed to stage changes: %v", err), operations, logs
	}
	addLog("SUCCESS", "Successfully staged changes", IconSuccess)

	// Commit changes
	commitMessage := config.CommitMessage
	if commitMessage == "auto-commit" {
		commitMessage = generateCommitMessage()
		addLog("INFO", fmt.Sprintf("Generated commit message: %s", commitMessage), IconSparkles)
	}

	if config.UseAICommit {
		addLog("INFO", "Using AI commit command", IconSparkles)
		// Use ai_commit command
		cmd := exec.Command("ai_commit", commitMessage)
		if err := cmd.Run(); err != nil {
			addLog("ERROR", fmt.Sprintf("ai_commit failed: %v", err), IconError)
			return false, fmt.Sprintf("ai_commit failed: %v", err), operations, logs
		}
		addLog("SUCCESS", "AI commit completed successfully", IconSuccess)
	} else {
		// Manual commit
		addLog("INFO", "Committing changes", IconCommit)
		if err := runGitCommand("commit", "-m", commitMessage); err != nil {
			addLog("ERROR", fmt.Sprintf("Failed to commit: %v", err), IconError)
			return false, fmt.Sprintf("Failed to commit: %v", err), operations, logs
		}
		addLog("SUCCESS", "Successfully committed changes", IconSuccess)

		// Push changes
		addLog("INFO", "Pushing changes to remote", IconPush)
		if err := runGitCommand("push"); err != nil {
			addLog("ERROR", fmt.Sprintf("Failed to push: %v", err), IconError)
			return false, fmt.Sprintf("Failed to push: %v", err), operations, logs
		}
		addLog("SUCCESS", "Successfully pushed changes", IconSuccess)
	}

	operations = append(operations, "committed and pushed changes")
	addLog("SUCCESS", "Repository processing completed", IconSparkles)
	
	return true, strings.Join(operations, ", "), operations, logs
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