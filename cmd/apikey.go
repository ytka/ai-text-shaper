package cmd

import (
	"fmt"
	"github.com/ytka/ai-text-shaper/internal/openai"
	"os"
	"strings"
)

func getAPIKeyFilePath() string {
	return os.Getenv("HOME") + "/.ai-text-shaper-apikey"
}

func checkAPIKeyFileExists() bool {
	_, err := os.Stat(getAPIKeyFilePath())
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func getAPIKey() (openai.APIKey, error) {
	apiKeyFilePath := getAPIKeyFilePath()
	bytes, err := os.ReadFile(apiKeyFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read API key: %w", err)
	}
	return openai.APIKey(strings.TrimSuffix(string(bytes), "\n")), nil
}