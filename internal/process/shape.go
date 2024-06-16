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
	supplementation := "Only the result shall be returned and the ai-text-shaper-input tag shall be removed. The result should be returned in the language of the Instruction, but if the Instruction has a language specification, that language should be given priority."
	mergedPrmpt := fmt.Sprintf(`<Instruction>%s. (%s)</Instruction>
<ai-text-shaper-input>%s</ai-text-shaper-input>`, prompt, supplementation, input)

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
