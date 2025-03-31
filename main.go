package main

import (
	"fmt"
	"github.com/NoamFav/auto_commit/src"
	"strings"
)

func main() {
	fmt.Println("Getting git info...")
	prompt := src.Summary()
	if prompt == "" {
		fmt.Println("✅ Nothing to commit.")
		return
	}

	prompt = "write a commit message with the following format in one sentece to be commited: <type>(<scope>): <subject>" + "\n" + prompt
	fmt.Println("Asking Ollama...")
	resp := src.AskOllama(prompt)
	resp = strings.TrimSpace(resp)

	fmt.Println("Committing...")
	src.AddCommitPush(resp)
}
