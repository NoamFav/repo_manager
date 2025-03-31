package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/NoamFav/auto_commit/src"
)

func main() {
	fmt.Print("Enter your prompt: ")

	reader := bufio.NewReader(os.Stdin)
	prompt, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("❌ Failed to read input:", err)
		return
	}

	prompt = strings.TrimSpace(prompt) // remove trailing newline
	fmt.Fprintln(os.Stdout, "🧠 Asking the LLM...\n")

	src.AskOllama(prompt)
}
