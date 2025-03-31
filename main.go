package main

import (
	"fmt"
	"github.com/NoamFav/auto_commit/src"
)

func main() {
	var prompt string

	fmt.Println("Getting git info...")
	prompt = src.Summary()
	fmt.Println("Prompt: ", prompt)

	prompt = "write a commit message with the following format in one sentece to be commited: <type>(<scope>): <subject>" + "\n" + prompt
	fmt.Println("Asking Ollama...")
	resp := src.AskOllama(prompt)

	fmt.Println("Committing...")
	src.AddCommitPush(resp)
}
