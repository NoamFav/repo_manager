package main

import (
	"fmt"
	"github.com/NoamFav/repo_manager/src/ai_commit"
	"strings"
)

func main() {
	fmt.Println("Getting git info...")
	prompt := src.GenerateCommitPrompt()

	if prompt == "" {
		fmt.Println(" Nothing to commit.")
		return // This should stop execution here
	}

	fmt.Println("Asking Ollama...")
	resp := src.AskOllama(prompt)
	resp = strings.TrimSpace(resp)

	fmt.Println("Committing...")
	src.AddCommitPush(resp)
}
