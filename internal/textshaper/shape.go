package textshaper

import (
	"ai-text-shaper/internal/openai"
	"fmt"
	"strings"
)

func ShapeText(apiKey, prompt, input string) (string, error) {
	mergedPrmpt := fmt.Sprintf("%s\n\n%s", prompt, input)
	result, err := openai.CallOpenAI(apiKey, "gpt-4o", mergedPrmpt)
	if err != nil {
		return "", fmt.Errorf("failed to call OpenAI API: %w", err)
	}

	result = strings.TrimSuffix(result, "\n")
	result = strings.TrimSpace(result)
	lines := strings.Split(result, "\n")
	if len(lines) > 0 && strings.HasPrefix(result, "```") {
		lines = lines[1:]
	}
	if len(lines) > 0 && strings.HasSuffix(result, "```") {
		lines = lines[:len(lines)-1]
	}

	return strings.Join(lines, "\n"), nil
}
