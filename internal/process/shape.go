package process

import (
	"fmt"
	"regexp"
	"strings"
)

type GenerativeAIClient interface {
	SendChatMessage(prompt string) (string, error)
}

func optimizePrompt(prompt, input string) string {
	supplementation := "Only the result shall be returned and the ai-text-shaper-input tag shall be removed. The result should be returned in the language of the Instruction, but if the Instruction has a language specification, that language should be given priority."
	mergedPrmpt := fmt.Sprintf(`<Instruction>%s. (%s)</Instruction>
<ai-text-shaper-input>%s</ai-text-shaper-input>`, prompt, supplementation, input)
	return mergedPrmpt
}

func optimizeResponseResult(rawResult string, useFirstCodeBlock bool) (string, error) {
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
			return "", fmt.Errorf("error finding first code block: %w", err)
		}
		if codeBlock != "" {
			result = codeBlock
		}
	}
	return result, nil
}

func sendOptimizedPromptAndResponse(gai GenerativeAIClient, prompt, input string, useFirstCodeBlock bool) (string, string, string, error) {
	optimized := optimizePrompt(prompt, input)
	rawResult, err := gai.SendChatMessage(optimized)
	if err != nil {
		return "", "", "", err
	}
	result, err := optimizeResponseResult(rawResult, useFirstCodeBlock)
	if err != nil {
		return "", "", "", err
	}
	return optimized, rawResult, result, nil
}

func sendRawPromptAndResponse(gai GenerativeAIClient, prompt string) (string, string, string, error) {
	rawResult, err := gai.SendChatMessage(prompt)
	if err != nil {
		return "", "", "", err
	}
	result := rawResult
	return prompt, rawResult, result, nil
}

func ShapeText(gai GenerativeAIClient, promptOrg, inputOrg string, useFirstCodeBlock bool) (string, string, string, error) {
	var prompt, rawResult, result string
	var err error

	if inputOrg == "" {
		prompt, rawResult, result, err = sendRawPromptAndResponse(gai, promptOrg)
	} else {
		prompt, rawResult, result, err = sendOptimizedPromptAndResponse(gai, promptOrg, inputOrg, useFirstCodeBlock)
	}
	if err != nil {
		return "", "", "", err
	}

	if !strings.HasSuffix(result, "\n") {
		result += "\n"
	}
	return prompt, rawResult, result, nil
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
