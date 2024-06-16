package process

import (
	"ai-text-shaper/internal/openai"
	"fmt"
	"regexp"
	"strings"
)

func findMarkdownFirstCodeBlock(text string) (string, error) {
	re, err := regexp.Compile("(?s)```[a-zA-Z0-9]*?\n(.*?)\n```")
	if err != nil {
		return "", fmt.Errorf("error compiling regex: %w", err)
	}
	match := re.FindStringSubmatch(text)
	if match != nil {
		return match[1], nil
	}
	return "", nil
}

func ShapeText(apiKey, prompt, input string, useFirstCodeBlock bool) (string, error) {
	mergedPrmpt := fmt.Sprintf("%s\n\n%s", prompt, input)
	resp, err := openai.SendChatMessage(apiKey, "gpt-4o", mergedPrmpt)
	if err != nil {
		return "", fmt.Errorf("failed to call OpenAI API: %w", err)
	}

	var result string
	for _, choice := range resp.Choices {
		result += choice.Message.Content
	}
	result = strings.TrimSuffix(result, "\n")
	result = strings.TrimSpace(result)
	/*
		lines := strings.Split(result, "\n")
		if len(lines) > 0 && strings.HasPrefix(result, "```") {
			lines = lines[1:]
		}
		if len(lines) > 0 && strings.HasSuffix(result, "```") {
			lines = lines[:len(lines)-1]
		}

		return strings.Join(lines, "\n") + "\n", nil
	*/
	if useFirstCodeBlock {
		codeBlock, err := findMarkdownFirstCodeBlock(result)
		if err != nil {
			return "", fmt.Errorf("error finding first code block: %w", err)
		}
		if codeBlock != "" {
			result = codeBlock
		}
	}

	return strings.TrimSuffix(result, "\n"), nil
}
