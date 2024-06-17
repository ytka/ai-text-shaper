package process

import (
	"ai-text-shaper/internal/openai"
	"fmt"
	"regexp"
	"strings"
)

type GenerativeAIClient interface {
	SendChatMessage(prompt string) (*openai.ChatCompletion, error)
}

type Shaper struct {
	gai                      GenerativeAIClient
	maxCompletionRepeatCount int
	useFirstCodeBlock        bool
}

func NewShaper(gai GenerativeAIClient, maxCompletionRepeatCount int, useFirstCodeBlock bool) *Shaper {
	return &Shaper{gai: gai, maxCompletionRepeatCount: maxCompletionRepeatCount, useFirstCodeBlock: useFirstCodeBlock}
}

func (s *Shaper) ShapeText(promptOrg, inputOrg string) (string, string, string, error) {
	var prompt, rawResult, result string
	var err error

	if inputOrg == "" {
		prompt, rawResult, result, err = sendRawPromptAndResponse(s.gai, promptOrg)
	} else {
		prompt, rawResult, result, err = sendOptimizedPromptAndResponse(s.gai, promptOrg, inputOrg, s.useFirstCodeBlock)
	}
	if err != nil {
		return "", "", "", err
	}

	if !strings.HasSuffix(result, "\n") {
		result += "\n"
	}
	return prompt, rawResult, result, nil
}

func optimizePrompt(prompt, input string) string {
	supplements := []string{
		"The subject of the Instruction is the area enclosed by the ai-text-shaper-input tag.",
		"The result should be returned in the language of the Instruction, but if the Instruction has a language specification, that language should be given priority.",
		"Provide additional explanations or details only if explicitly requested in the Instruction.",
		// "Only the result shall be returned.",
	}
	supplementation := strings.Join(supplements, " ")
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

func sendChatMessage(gai GenerativeAIClient, prompt string) (string, error) {
	rawResult, err := gai.SendChatMessage(prompt)
	if err != nil {
		return "", fmt.Errorf("failed to send chat message: %w", err)
	}

	var result string
	for _, choice := range rawResult.Choices {
		result += choice.Message.Content
	}
	return result, nil
}

func sendOptimizedPromptAndResponse(gai GenerativeAIClient, prompt, input string, useFirstCodeBlock bool) (string, string, string, error) {
	optimized := optimizePrompt(prompt, input)
	rawResult, err := sendChatMessage(gai, optimized)
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
	rawResult, err := sendChatMessage(gai, prompt)
	if err != nil {
		return "", "", "", err
	}
	result := rawResult
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
