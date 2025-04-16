<div align="center">

# 🚀 Repository Manager

**A powerful toolkit for managing multiple Git repositories with AI-powered features**

[![Go Version](https://img.shields.io/badge/Go-1.15%2B-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![Python Version](https://img.shields.io/badge/Python-3.8%2B-blue?style=for-the-badge&logo=python)](https://www.python.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](https://opensource.org/licenses/MIT)

</div>

---

## 🔍 Overview

Repository Manager is a comprehensive suite of tools designed to streamline the management of multiple Git repositories. With intelligent AI-powered commit generation, automatic syncing, and GitHub repository management, it transforms your Git workflow into a seamless experience.

<div align="center">

### ✨ Key Features

| Feature | Description |
|---------|-------------|
| 🧠 **AI-Powered Commits** | Generate intelligent commit messages based on code changes |
| 🔄 **Batch Operations** | Process multiple repositories with a single command |
| 📊 **Rich Visualizations** | Beautiful terminal UI with detailed repository insights |
| 🌐 **GitHub Integration** | Clone and manage repositories with powerful filtering options |
| 🎨 **Customizable Workflows** | Tailor operations to your specific needs with extensive options |

</div>

## 📋 Table of Contents

- [Components](#-components)
- [Repository Structure](#-repository-structure)
- [Installation](#-installation)
- [Usage](#-usage)
  - [AI Commit](#ai-commit)
  - [Auto Commit](#auto-commit)
  - [Pull Repos](#pull-repos)
  - [Useful Aliases](#useful-aliases)
- [Features](#-features)
- [Requirements](#-requirements)
- [Roadmap](#-roadmap)
- [License](#-license)

## 🧩 Components

<details open>
<summary><b>AI Commit (Go)</b></summary>
<br>

> Intelligence for your Git commits

- Analyzes Git diffs and staged changes with precision
- Generates contextual commit messages based on code analysis
- Automatically detects the type and scope of changes for conventional commits
- Integrates with Ollama for natural language processing
- Supports customizable prompts for specialized commit formats

</details>

<details open>
<summary><b>Auto Commit (Python)</b></summary>
<br>

> Batch processing for multiple repositories

- Processes multiple repositories in a directory simultaneously
- Handles .gitignore updates and .DS_Store removal automatically
- Features a beautiful terminal UI with rich visualizations
- Offers pull-before-commit option to prevent merge conflicts
- Optionally uses AI Commit for intelligent commit messages

</details>

<details open>
<summary><b>Pull Repos (Python)</b></summary>
<br>

> Streamlined GitHub repository management

- Clones repositories from GitHub using the GitHub CLI
- Provides detailed visualization of repository information
- Filters repositories by stars, forks, and more
- Shows repository size and structure information
- Supports batch operations with comprehensive progress tracking

</details>

<details>
<summary><b>Coming Soon</b></summary>
<br>

The following components are currently under development:

- **📊 Dashboard**: A web-based dashboard for visualizing repository metrics and status
- **🏗️ Repo Init**: Tools for initializing new repositories with templates and configurations
- **🧹 Cleaner**: Utilities for cleaning and maintaining repositories
- **⏰ Schedule**: Automated scheduling for repository maintenance and updates

</details>

## 📁 Repository Structure

```
.
├── 📄 go.mod                   # Go module definition
├── 📄 go.sum                   # Go dependencies checksum
├── 📄 main.go                  # Main Go entry point
├── 📄 pyproject.toml           # Python project configuration
├── 📄 README.md                # Project documentation
├── 📄 requirements.txt         # Python dependencies
└── 📁 src                      # Source code directory
    ├── 📄 __init__.py          # Python package initialization
    ├── 📄 __main__.py          # Python entry point
    ├── 📁 __pycache__          # Python bytecode cache
    ├── 📁 ai_commit            # AI Commit component (Go)
    │   ├── 📄 ai.go            # AI integration code
    │   └── 📄 git_info.go      # Git information handling
    ├── 📁 cleaner              # Repository cleaning utilities
    │   └── 📄 cleaner.py       # Cleaning implementation
    ├── 📁 dashboard            # Web dashboard component
    │   ├── 📁 __pycache__      # Python bytecode cache
    │   └── 📄 app.py           # Dashboard application
    ├── 📁 repo_init            # Repository initialization
    │   ├── 📄 creator.py       # Repository creation
    │   └── 📄 handler.py       # Initialization handler
    ├── 📁 repo_manager         # Core repository management
    │   ├── 📄 auto_commit.py   # Auto commit implementation
    │   └── 📄 pull_repos.py    # Repository pulling implementation
    ├── 📄 run.py               # Main runner script
    └── 📁 schedule             # Scheduling component
        ├── 📄 alarm.py         # Schedule alerting
        └── 📄 cron.py          # Cron job management
```

## 🔧 Installation

### Prerequisites

<details open>
<summary>Required software</summary>
<br>

- **Go 1.15+**: Required for AI Commit functionality
- **Python 3.8+**: Required for Auto Commit and Pull Repos
- **Git**: Core requirement for all components
- **GitHub CLI**: Required for Pull Repos functionality
- **Ollama**: Required for AI-powered commit message generation
  - Installation: [ollama.ai](https://ollama.ai)
  - Mistral model setup: `ollama pull mistral`

</details>

### Step-by-Step Installation

<details open>
<summary><b>Step 1: Clone the repository</b></summary>
<br>

```bash
git clone https://github.com/NoamFav/repo_manager.git
cd repo_manager
```

</details>

<details open>
<summary><b>Step 2: Install Python dependencies</b></summary>
<br>

```bash
# Install all required Python packages
pip install -r requirements.txt
```

</details>

<details open>
<summary><b>Step 3: Build the Go module</b></summary>
<br>

```bash
# Navigate to the AI Commit directory
cd src/ai_commit

# Build the Go binary
go build -o ai_commit

# Move the binary to your path
sudo mv ai_commit /usr/local/bin/

# Return to the repository root
cd ../..
```

</details>

<details open>
<summary><b>Step 4: Set up Python executables</b></summary>
<br>

```bash
# Make Python scripts executable
chmod +x src/repo_manager/auto_commit.py
chmod +x src/repo_manager/pull_repos.py

# Create symbolic links for system-wide access
sudo ln -s "$(pwd)/src/repo_manager/auto_commit.py" /usr/local/bin/auto_commit
sudo ln -s "$(pwd)/src/repo_manager/pull_repos.py" /usr/local/bin/pull_repos
```

> **Note:** These commands create symbolic links from your current working directory. Make sure you're in the repository root when running them.

</details>

<details open>
<summary><b>Step 5: Enable autocomplete (optional)</b></summary>
<br>

```bash
# Install argcomplete if not already installed
pip install argcomplete

# Register autocompletion for the commands
eval "$(register-python-argcomplete auto_commit)"
eval "$(register-python-argcomplete pull_repos)"
```

Add these eval lines to your `.bashrc` or `.zshrc` to make the autocompletion persistent.

</details>

## 🚀 Usage

### AI Commit

AI Commit generates intelligent commit messages based on your Git changes.

<details open>
<summary><b>Basic Usage</b></summary>
<br>

```bash
# Generate a commit message based on your git changes
ai_commit

# With a custom prompt for context
ai_commit "Add more context to the commit message"
```

</details>

<details>
<summary><b>Example Output</b></summary>
<br>

```
✨ AI Commit
┌─────────────────────────────────────────────────┐
│ Analyzing git changes in current repository...  │
│ Found 3 modified files and 1 new file          │
│                                                 │
│ 🧠 Generating commit message...                 │
└─────────────────────────────────────────────────┘

📝 Generated Commit Message:
feat(auth): implement OAuth2 authentication flow

- Add OAuth2 provider integration in auth_service.go
- Update user model to include OAuth tokens
- Fix token refresh mechanism
- Add tests for authentication workflow
```

</details>

### Auto Commit

Auto Commit processes multiple repositories with beautiful visualizations.

<details open>
<summary><b>Command Options</b></summary>
<br>

```bash
# Basic usage (processes all Git repositories in ~/Neoware)
auto_commit

# Specify a different directory
auto_commit --dir ~/Projects

# Use a specific commit message
auto_commit --commit-message "update dependencies"

# Pull changes before committing
auto_commit --pull

# Exclude specific repositories
auto_commit --exclude repo1 repo2

# Only process specific repositories
auto_commit --only repo1 repo2

# Handle .gitignore and .DS_Store files
auto_commit --handle-gitignore --remove-ds-store

# Use manual git commands instead of ai_commit
auto_commit --no-auto-commit
```

</details>

<details>
<summary><b>Advanced Examples</b></summary>
<br>

```bash
# Process all repositories, pull changes, and update .gitignore files
auto_commit --pull --handle-gitignore

# Only process the 'myproject' repository and use AI to generate commit messages
auto_commit --only myproject

# Process all repositories except for 'temp' and 'test', with a specific commit message
auto_commit --exclude temp test --commit-message "update documentation"
```

</details>

<details>
<summary><b>Example Output</b></summary>
<br>

```
🔄 Auto Commit
┌─────────────────────────────────────────────────┐
│ Processing repositories in ~/Projects           │
│ Found 5 Git repositories                        │
└─────────────────────────────────────────────────┘

📊 Repository Status:
┌────────────────┬─────────┬──────────┬──────────────┐
│ Repository     │ Status  │ Changes  │ Last Commit  │
├────────────────┼─────────┼──────────┼──────────────┤
│ project-alpha  │ ✅ OK   │ 3 files  │ 2 hours ago  │
│ project-beta   │ ✅ OK   │ 1 file   │ 5 days ago   │
│ docs-site      │ ✅ OK   │ 0 files  │ 1 week ago   │
│ api-service    │ ⚠️ WARN │ 7 files  │ 3 days ago   │
│ mobile-app     │ ✅ OK   │ 2 files  │ 1 day ago    │
└────────────────┴─────────┴──────────┴──────────────┘

🔄 Processing repositories... [3/5]
[project-alpha] ✓ Committed and pushed 3 files
[project-beta] ✓ Committed and pushed 1 file
[api-service] ⚠️ Merge conflicts detected. Skipping.

📝 Summary:
Successfully processed 4/5 repositories
```

</details>

### Pull Repos

Pull Repos clones GitHub repositories with detailed visualizations.

<details open>
<summary><b>Command Options</b></summary>
<br>

```bash
# Basic usage (clones repositories to ~/Neoware)
pull_repos

# Specify a different target directory
pull_repos --base-dir ~/Projects

# Limit the number of repositories to fetch
pull_repos --limit 10

# Filter out forked repositories
pull_repos --filter-forks

# Only clone repositories with at least 5 stars
pull_repos --only-stars 5

# Exclude specific repositories
pull_repos --exclude user/repo1 user/repo2
```

</details>

<details>
<summary><b>Advanced Examples</b></summary>
<br>

```bash
# Clone only non-fork repositories with at least 10 stars
pull_repos --filter-forks --only-stars 10

# Clone up to 5 repositories to a specific directory
pull_repos --limit 5 --base-dir ~/Projects/new-repos

# Clone all repositories except specific ones
pull_repos --exclude organization/repo1 organization/repo2
```

</details>

<details>
<summary><b>Example Output</b></summary>
<br>

```
🌐 Pull Repos
┌─────────────────────────────────────────────────┐
│ Fetching GitHub repositories...                 │
│ Using filters: --only-stars 5 --filter-forks    │
└─────────────────────────────────────────────────┘

📊 Repository Information:
┌─────────────────────┬─────────┬───────┬──────────┬─────────┐
│ Repository          │ Stars   │ Forks │ Size     │ Status  │
├─────────────────────┼─────────┼───────┼──────────┼─────────┤
│ user/awesome-app    │ ⭐ 128  │ 23    │ 4.2 MB   │ Public  │
│ user/data-tools     │ ⭐ 67   │ 12    │ 1.8 MB   │ Public  │
│ org/web-framework   │ ⭐ 3.5k │ 342   │ 12.6 MB  │ Public  │
└─────────────────────┴─────────┴───────┴──────────┴─────────┘

🔄 Cloning repositories... [2/3]
[user/awesome-app] ✓ Cloned successfully
[user/data-tools] ✓ Cloned successfully
[org/web-framework] ⏳ Cloning...

📝 Summary:
Successfully cloned 3/3 repositories to ~/Projects
```

</details>

### Useful Aliases

<details open>
<summary><b>Quick Commit Alias</b></summary>
<br>

Add this alias to your `.bashrc` or `.zshrc` file for quick access to auto-commit for the current repository:

```bash
# Add this to your .bashrc or .zshrc
alias gacp='auto_commit --only "${PWD##*/}"'
```

This creates a `gacp` command (Git Add, Commit, Push) that automatically determines your current repository name and only processes that repository.

</details>

<details>
<summary><b>Additional Aliases</b></summary>
<br>

```bash
# Quick commit with AI-generated message
alias gaic='ai_commit && git push'

# Pull all repositories in your projects directory
alias gpa='pull_repos --base-dir ~/Projects'

# Clean and update all repositories
alias gcl='auto_commit --handle-gitignore --remove-ds-store --pull'
```

</details>

## ✨ Features

<div align="center">

| Feature | Description |
|---------|-------------|
| **Rich Terminal UI** | Beautiful, colorful interfaces with progress tracking and visualizations |
| **Smart Detection** | Automatically detects commit types and scopes based on code patterns |
| **Batch Processing** | Process multiple repositories simultaneously with intelligent handling |
| **GitHub Integration** | Seamless integration with GitHub repositories and metadata |
| **Customizable Workflows** | Extensive command-line options for tailored experiences |
| **AI-Powered Intelligence** | Generate contextual commit messages with Ollama integration |
| **Detailed Analytics** | Comprehensive information about repositories and changes |
| **Cross-Platform** | Works on macOS, Linux, and Windows (with WSL) |

</div>

## 📋 Requirements

<details open>
<summary><b>Core Requirements</b></summary>
<br>

- **Go 1.15+**: Required for AI Commit
- **Python 3.8+**: Required for Auto Commit and Pull Repos
- **Rich Package**: Required for beautiful terminal UI
- **Git**: Required for all components
- **GitHub CLI**: Required for Pull Repos
- **Ollama**: Required for AI Commit to generate intelligent messages

</details>

<details>
<summary><b>Ollama Configuration</b></summary>
<br>

The AI Commit component uses Ollama with the Mistral model by default. To use a different model, modify the model name in `src/ai_commit/ai.go`:

```go
// In src/ai_commit/ai.go
// Change "mistral" to your preferred model name
requestBody, _ := json.Marshal(map[string]interface{}{
    "model":  "mistral", // Change this to use a different model
    "prompt": prompt,
    "stream": true,
})
```

Available models include:
- mistral
- llama2
- codellama
- vicuna
- orca-mini

</details>

## 🚀 Roadmap

<details open>
<summary><b>Upcoming Features</b></summary>
<br>

- **Web Dashboard**: Interactive web interface for repository insights
- **Team Collaboration**: Multi-user support for team repository management
- **Advanced Analytics**: Repository health metrics and contribution analytics
- **Integration Ecosystem**: Support for GitLab, Bitbucket, and other Git providers
- **Custom Templates**: Templating system for repository initialization
- **Automated Testing**: Integrated testing workflows for repositories

</details>

## 📜 License

This project is licensed under the [MIT License](LICENSE).

<div align="center">

---

**Made with ❤️ by your development team**

[GitHub](https://github.com/NoamFav/repo_manager) • [Documentation](https://NoamFav.github.io/repo_manager) • [Issues](https://github.com/NoamFav/repo_manager/issues)

</div>
