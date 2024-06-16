package runner

import (
	"fmt"
	"os"
	"strings"
)

type APIKey string

func getAPIKey() (APIKey, error) {
	apiKeyFilePath := os.Getenv("HOME") + "/.ai-text-shaper-apikey"
	bytes, err := os.ReadFile(apiKeyFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read API key: %w", err)
	}
	return APIKey(strings.TrimSuffix(string(bytes), "\n")), nil
}
