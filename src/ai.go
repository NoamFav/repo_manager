package src

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func AskOllama(prompt string) string {
	requestBody, _ := json.Marshal(map[string]interface{}{
		"model":  "mistral",
		"prompt": prompt,
		"stream": true,
	})

	resp, err := http.Post("http://127.0.0.1:11434/api/generate", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatal("Failed to connect to Ollama:", err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	var builder strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		var chunk struct {
			Response string `json:"response"`
		}
		if err := json.Unmarshal([]byte(line), &chunk); err == nil {
			fmt.Print(chunk.Response)           // real-time print
			builder.WriteString(chunk.Response) // also save it
		}
	}

	fmt.Println() // line break after the stream ends
	return builder.String()
}
