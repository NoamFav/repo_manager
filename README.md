# Repository Manager

A comprehensive toolkit for managing multiple Git repositories with intelligent commit messages, automatic syncing, and GitHub repository cloning.

## Table of Contents

- [Overview](#overview)
- [Components](#components)
- [Repository Structure](#repository-structure)
- [Installation](#installation)
- [Usage](#usage)
  - [AI Commit](#ai-commit)
  - [Auto Commit](#auto-commit)
  - [Pull Repos](#pull-repos)
  - [Useful Aliases](#useful-aliases)
- [Features](#features)
- [Requirements](#requirements)
- [License](#license)

## Overview

This Repository Manager toolkit provides a seamless workflow for managing multiple Git repositories. It consists of three main components:

1. **AI Commit** (Go): Generates intelligent commit messages based on your code changes.
2. **Auto Commit** (Python): Batch processes multiple repositories, automatically committing and pushing changes.
3. **Pull Repos** (Python): Clones and updates multiple GitHub repositories.

All tools feature rich, colorful terminal interfaces with intuitive visualizations to make repository management a pleasant experience.

## Components

### AI Commit (Go)

- Analyzes Git diffs and staged changes
- Generates commit messages based on code analysis
- Detects the type and scope of changes for conventional commits
- Integrates with Ollama for natural language processing

### Auto Commit (Python)

- Batch processes multiple repositories in a directory
- Handles .gitignore updates and .DS_Store removal
- Beautiful terminal UI with rich visualizations
- Optionally uses AI Commit for intelligent commit messages

### Pull Repos (Python)

- Clones repositories from GitHub using the GitHub CLI
- Provides detailed visualization of repository information
- Filters repositories by stars, forks, and more
- Shows repository size and structure information

### In-Progress Components

The following components are currently under development:

- **Dashboard**: A web-based dashboard for visualizing repository metrics and status
- **Repo Init**: Tools for initializing new repositories with templates and configurations
- **Cleaner**: Utilities for cleaning and maintaining repositories
- **Schedule**: Automated scheduling for repository maintenance and updates

## Repository Structure

```
.
├──  go.mod
├──  go.sum
├──  main.go
├──  pyproject.toml
├── 󰂺 README.md
├──  requirements.txt
└──  src
    ├──  __init__.py
    ├──  __main__.py
    ├──  __pycache__
    │   ├──  __init__.cpython-312.pyc
    │   └──  __main__.cpython-312.pyc
    ├──  ai_commit
    │   ├──  ai.go
    │   └──  git_info.go
    ├──  cleaner
    │   └──  cleaner.py
    ├──  dashboard
    │   ├──  __pycache__
    │   │   └──  app.cpython-312.pyc
    │   └──  app.py
    ├──  repo_init
    │   ├──  creator.py
    │   └──  handler.py
    ├──  repo_manager
    │   ├──  auto_commit.py
    │   └──  pull_repos.py
    ├──  run.py
    └──  schedule
        ├──  alarm.py
        └──  cron.py
```

## Installation

### Prerequisites

- Go 1.15+
- Python 3.8+
- Git
- [GitHub CLI](https://cli.github.com/) (for Pull Repos)
- [Ollama](https://ollama.ai/) (for AI Commit)
- [Rich](https://github.com/Textualize/rich) Python package

### Step 1: Clone the repository

```bash
git clone https://github.com/yourusername/repo_manager.git
cd repo_manager
```

### Step 2: Install Python dependencies

```bash
pip install rich
```

### Step 3: Build the Go module

```bash
# From the repository root
cd src/ai_commit
go build -o ai_commit
sudo mv ai_commit /usr/local/bin/
```

### Step 4: Set up Python executables

Make the Python scripts executable:

```bash
# From the repository root
chmod +x src/repo_manager/auto_commit.py
chmod +x src/repo_manager/pull_repos.py
```

Create symbolic links to make them available system-wide:

```bash
# From the repository root
sudo ln -s "$(pwd)/src/repo_manager/auto_commit.py" /usr/local/bin/auto_commit
sudo ln -s "$(pwd)/src/repo_manager/pull_repos.py" /usr/local/bin/pull_repos
```

> **Note:** The above commands create symbolic links from your current working directory. Make sure you're in the repository root when running them.

## Usage

### AI Commit

AI Commit is used to generate intelligent commit messages based on your Git changes. It analyzes diffs, determines the type and scope of changes, and generates a conventional commit message.

```bash
# Basic usage - analyzes your git changes and creates a commit
ai_commit

# With a custom prompt
ai_commit "Add more context to the commit message"
```

> **Note:** AI Commit requires Ollama to be running locally with the mistral model. It sends your git diff to Ollama to generate commit messages.

### Auto Commit

Auto Commit processes multiple repositories with beautiful visualizations in your terminal.

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

**Detailed examples:**

```bash
# Process all repositories, pull changes, and update .gitignore files
auto_commit --pull --handle-gitignore

# Only process the 'myproject' repository and use AI to generate commit messages
auto_commit --only myproject

# Process all repositories except for 'temp' and 'test', with a specific commit message
auto_commit --exclude temp test --commit-message "update documentation"
```

Auto Commit provides a rich visualization experience with:

- Repository summaries and statistics
- Commit details with file changes
- Color-coded status information
- Progress tracking for batch operations

### Pull Repos

Pull Repos clones GitHub repositories with detailed visualizations.

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

**Detailed examples:**

```bash
# Clone only non-fork repositories with at least 10 stars
pull_repos --filter-forks --only-stars 10

# Clone up to 5 repositories to a specific directory
pull_repos --limit 5 --base-dir ~/Projects/new-repos

# Clone all repositories except specific ones
pull_repos --exclude organization/repo1 organization/repo2
```

Pull Repos provides:

- Repository metadata visualization (stars, forks, privacy status)
- Size and structure information for cloned repositories
- Colorful progress tracking
- Filtering options based on GitHub metadata

### Useful Aliases

Add this alias to your `.bashrc` or `.zshrc` file for quick access to auto-commit for the current repository:

```bash
# Add this to your .bashrc or .zshrc
alias gacp='auto_commit --only "${PWD##*/}"'
```

This creates a `gacp` command (Git Add, Commit, Push) that automatically determines your current repository name and only processes that repository.

## Features

- **Rich Terminal UI**: Beautiful, colorful, and informative terminal interfaces with progress tracking and visualizations.
- **Smart Detection**: Automatically detects commit types and scopes based on changes by analyzing code patterns.
- **Batch Processing**: Process multiple repositories with a single command, handling different states and conditions.
- **GitHub Integration**: Seamlessly clone and manage GitHub repositories with metadata filtering.
- **Customizable**: Many command-line options to customize behavior for different workflows.
- **AI-Powered**: Generate intelligent commit messages with Ollama integration using the mistral model.
- **Detailed Summaries**: Get comprehensive information about repositories including files changed, commit statistics, and more.
- **In-Progress Features**: Active development on dashboard, scheduling, and more comprehensive repository management tools.

## Requirements

- **Go**: Required for AI Commit
- **Python 3.8+**: Required for Auto Commit and Pull Repos
- **Rich Package**: Required for beautiful terminal UI
- **Git**: Required for all components
- **GitHub CLI**: Required for Pull Repos
- **Ollama**: Required for AI Commit to generate intelligent messages
  - You must install Ollama from [ollama.ai](https://ollama.ai)
  - You need to pull the Mistral model: `ollama pull mistral`
  - To use a different model, modify the model name in `src/ai_commit/ai.go`

```go
// In src/ai_commit/ai.go
// Change "mistral" to your preferred model name
requestBody, _ := json.Marshal(map[string]interface{}{
    "model":  "mistral", // Change this to use a different model
    "prompt": prompt,
    "stream": true,
})
```

## License

[MIT License](LICENSE)
