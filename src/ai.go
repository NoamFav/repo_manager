package src

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func AskOllama(prompt string) {
	requestBody, _ := json.Marshal(map[string]interface{}{
		"model":  "mistral", // your working model
		"prompt": prompt,
		"stream": true,
	})

	resp, err := http.Post("http://127.0.0.1:11434/api/generate", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatal("Failed to connect to Ollama:", err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()

		// Each line is a tiny JSON chunk like {"response":"token"}
		var chunk struct {
			Response string `json:"response"`
		}
		if err := json.Unmarshal([]byte(line), &chunk); err == nil {
			fmt.Print(chunk.Response) // print as it arrives, no newline
		}
	}

	fmt.Println() // add newline after full output
}
