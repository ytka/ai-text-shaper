package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ytka/textforge/internal/openai"
)

func getAPIKeyFilePath() string {
	return os.Getenv("HOME") + "/.textforge-apikey"
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
