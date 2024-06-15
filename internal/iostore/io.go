package iostore

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Logger interface {
	Msg(msg string, args ...interface{})
	VerboseMsg(msg string, args ...interface{})
}

type IOStore struct {
	verboseLog func(msg string, args ...interface{})
}

func New(logger func(msg string, args ...interface{})) *IOStore {
	return &IOStore{verboseLog: logger}
}

func (ios *IOStore) GetAPIKey() (string, error) {
	ios.verboseLog("Getting API key")

	apiKeyFilePath := os.Getenv("HOME") + "/.openai-apikey"
	bytes, err := os.ReadFile(apiKeyFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read API key: %w", err)
	}
	return strings.TrimSuffix(string(bytes), "\n"), nil
}

func (ios *IOStore) GetPromptText(prompt, promptPath string) (string, error) {
	ios.verboseLog("Getting prompt text")

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

func (ios *IOStore) GetInputText(inputFilePath string) (string, error) {
	ios.verboseLog("Getting input text")

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

func (ios *IOStore) WriteToFile(path, content string) error {
	ios.verboseLog("Writing to file: %s", path)

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	return nil
}
