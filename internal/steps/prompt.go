package steps

import (
	"errors"
	"fmt"
	"os"
)

var ErrPromptRequired = errors.New("prompt is required")

func GetPromptText(prompt, promptPath string) (string, error) {
	if prompt == "" && promptPath == "" {
		return "", ErrPromptRequired
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
