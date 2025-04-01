package main

import (
	"fmt"
	"github.com/NoamFav/auto_commit/src"
	"strings"
)

func main() {
	fmt.Println("Getting git info...")
	prompt := fmt.Sprintf(`You are an AI Git assistant. Your task is to write a single-line, conventional commit message in the format:
<type>(<scope>): <subject>

Be concise, avoid bullet points or summaries. ONLY return the commit message, nothing else.

Git diff and status:
%s`, src.Summary())

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
