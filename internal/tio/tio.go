package tio

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type logger interface {
	Msg(msg string, args ...interface{})
}

type TIO struct {
	logger logger
}

func NewTIO(logger logger) *TIO {
	return &TIO{logger: logger}
}

func GetAPIKey() (string, error) {
	apiKeyFilePath := os.Getenv("HOME") + "/.openai-apikey"
	bytes, err := os.ReadFile(apiKeyFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read API key: %w", err)
	}
	return strings.TrimSuffix(string(bytes), "\n"), nil
}

func GetPromptText(prompt, promptPath string) (string, error) {
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

func GetInputText(inputFilePath string) (string, error) {
	if inputFilePath == "-" {
		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", fmt.Errorf("error reading input from stdin: %w", err)
		}
		return string(input), nil
	}

	input, err := os.ReadFile(inputFilePath)
	if err != nil {
		return "", fmt.Errorf("error reading input file: %w", err)
	}
	return string(input), nil
}

func WriteToFile(path, content string) error {
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	return nil
}
