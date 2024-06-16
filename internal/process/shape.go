package process

import (
	"fmt"
	"regexp"
	"strings"
)

type GenerativeAIClient interface {
	SendChatMessage(prompt string) (string, error)
}

// Return only the results
func ShapeText(gai GenerativeAIClient, prompt, input string, useFirstCodeBlock bool) (string, string, error) {
	mergedPrmpt := fmt.Sprintf(`<Instruction>%s. (Return only the results and remove ai-text-shaper-input tag)</Instruction>
<ai-text-shaper-input>%s</ai-text-shaper-input>`, prompt, input)

	result, err := gai.SendChatMessage(mergedPrmpt)
	if err != nil {
		return "", "", err
	}
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
			return "", "", fmt.Errorf("error finding first code block: %w", err)
		}
		if codeBlock != "" {
			result = codeBlock
		}
	}

	return mergedPrmpt, strings.TrimSuffix(result, "\n"), nil
}

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
