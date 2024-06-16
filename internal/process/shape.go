package process

import (
	"fmt"
	"regexp"
	"strings"
)

type GenerativeAIClient interface {
	SendChatMessage(prompt string) (string, error)
}

func ShapeText(gai GenerativeAIClient, prompt, input string, useFirstCodeBlock bool) (string, string, string, error) {
	mergedPrmpt := fmt.Sprintf(`<Instruction>%s. (Return only the results and remove ai-text-shaper-input tag)</Instruction>
<ai-text-shaper-input>%s</ai-text-shaper-input>`, prompt, input)

	rawResult, err := gai.SendChatMessage(mergedPrmpt)
	if err != nil {
		return "", "", "", err
	}
	result := rawResult
	if strings.HasPrefix(result, "```") && strings.HasSuffix(result, "```") {
		// remove first and last line
		lines := strings.Split(result, "\n")
		if len(lines) > 2 {
			result = strings.Join(lines[1:len(lines)-1], "\n")
		}
	}

	if useFirstCodeBlock {
		codeBlock, err := findMarkdownFirstCodeBlock(result)
		if err != nil {
			return "", "", "", fmt.Errorf("error finding first code block: %w", err)
		}
		if codeBlock != "" {
			result = codeBlock
		}
	}

	if !strings.HasSuffix(result, "\n") {
		result += "\n"
	}
	return mergedPrmpt, rawResult, result, nil
}

func findMarkdownFirstCodeBlock(text string) (string, error) {
	re, err := regexp.Compile("(?s)```[a-zA-Z0-9]*?\n(.*?\n)```")
	if err != nil {
		return "", fmt.Errorf("error compiling regex: %w", err)
	}
	match := re.FindStringSubmatch(text)
	if match != nil {
		return match[1], nil
	}
	return "", nil
}
