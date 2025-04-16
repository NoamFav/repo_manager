# Contributing to Repository Manager

First off, thank you for considering contributing to Repository Manager! It's people like you that make this tool better for everyone.

## Table of Contents
- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Documentation](#documentation)
- [Feature Requests and Bug Reports](#feature-requests-and-bug-reports)

## Code of Conduct

This project and everyone participating in it is governed by the [Repository Manager Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## Getting Started

1. **Fork the Repository**: Start by forking the repository on GitHub.

2. **Clone Your Fork**: 
   ```bash
   git clone https://github.com/YOUR-USERNAME/repo_manager.git
   cd repo_manager
   ```

3. **Set Up Development Environment**:
   - For Go components (AI Commit):
     ```bash
     # Ensure you have Go installed
     go mod download
     ```
   - For Python components:
     ```bash
     # Ensure you have Python 3.8+ installed
     pip install -r requirements.txt
     ```

4. **Create a New Branch**: 
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Development Workflow

### Working on Go Components

The Go components (AI Commit) are located in the `src/ai_commit` directory. To test your changes:

1. Make your changes to the code
2. Build and test locally:
   ```bash
   cd src/ai_commit
   go build -o ai_commit
   ./ai_commit
   ```

### Working on Python Components

The Python components (Auto Commit, Pull Repos) are located in the `src/repo_manager` directory. To test your changes:

1. Make your changes to the code
2. Run the script directly:
   ```bash
   python src/repo_manager/auto_commit.py
   # or
   python src/repo_manager/pull_repos.py
   ```

## Pull Request Process

1. **Update Documentation**: Ensure you've updated any relevant documentation.
2. **Add Tests**: If applicable, add tests for your changes.
3. **Ensure Code Quality**: Run any linters or formatters to maintain code quality.
4. **Create a Pull Request**: Submit a pull request from your fork to the main repository.
5. **Describe Your Changes**: In the pull request, describe what your changes do and why they should be included.
6. **Reference Issues**: If your pull request addresses an issue, reference it in the description.

## Coding Standards

### Go

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Format your code with `gofmt`
- Use meaningful variable and function names
- Write comments for exported functions

### Python

- Follow [PEP 8](https://www.python.org/dev/peps/pep-0008/) style guide
- Use meaningful variable and function names
- Document functions using docstrings
- Keep line length under 100 characters

## Testing

- For Go components, add tests in the same directory with `_test.go` suffix
- For Python components, add tests in a `tests` directory
- Run tests before submitting your pull request

## Documentation

- Update README.md if you change functionality
- Add docstrings/comments to your code
- If you add new commands or options, update the usage documentation

## Feature Requests and Bug Reports

- Use the GitHub Issues tracker to submit feature requests and bug reports
- Clearly describe the issue or feature
- For bugs, include steps to reproduce, expected behavior, and actual behavior
- For features, explain why the feature would be useful to others

---

Thank you for taking the time to contribute to Repository Manager! Your contributions help make this project better for everyone.
