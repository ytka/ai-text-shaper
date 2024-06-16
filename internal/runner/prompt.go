package runner

import (
	"fmt"
	"os"
)

func getPromptText(prompt, promptPath string) (string, error) {
	if prompt == "" && promptPath == "" {
		return "", fmt.Errorf("prompt is required")
	}
	if prompt == "" && promptPath != "" {
		text, err := os.ReadFile(promptPath)
		if err != nil {
			return "", fmt.Errorf("error reading prompt file: %w", err)
		}
		return string(text), nil
	}
	return prompt, nil
}
