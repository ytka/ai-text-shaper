package steps

import (
	"errors"
	"fmt"
	"os"
)

// ErrPromptRequired is exported and uses CamelCase.
var ErrPromptRequired = errors.New("prompt is required")

// GetPromptText retrieves the prompt text from the specified source.
func GetPromptText(prompt, promptPath string) (string, error) {
	if prompt == "" && promptPath == "" {
		return "", ErrPromptRequired
	}

	if promptPath == "-" {
		return getInputFromStdin()
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
