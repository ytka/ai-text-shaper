package steps

import (
	"fmt"
	"regexp"
	"strings"

	"github/ytka/ai-text-shaper/internal/openai"
)

// GenerativeAIClient represents an interface for generating AI client operations.
type GenerativeAIClient interface {
	RequestCreateChatCompletion(*openai.CreateChatCompletion) (*openai.ChatCompletion, error)
	MakeCreateChatCompletion(prompt string) *openai.CreateChatCompletion
}

// ShapeResult represents the result of a text shaping operation.
type ShapeResult struct {
	Prompt    string
	RawResult string
	Result    string
}

// NewShapeResult creates a new ShapeResult.
func NewShapeResult(prompt, rawResult, result string) *ShapeResult {
	if !strings.HasSuffix(result, "\n") {
		result += "\n"
	}
	return &ShapeResult{
		Prompt:    prompt,
		RawResult: rawResult,
		Result:    result,
	}
}

// Shaper is responsible for shaping the text by interacting with GenerativeAIClient.
type Shaper struct {
	gai                      GenerativeAIClient
	maxCompletionRepeatCount int
	useFirstCodeBlock        bool
	promptOptimize           bool
}

// NewShaper creates a new Shaper.
func NewShaper(gai GenerativeAIClient, maxCompletionRepeatCount int, useFirstCodeBlock, promptOptimize bool) *Shaper {
	return &Shaper{
		gai:                      gai,
		maxCompletionRepeatCount: maxCompletionRepeatCount,
		useFirstCodeBlock:        useFirstCodeBlock,
		promptOptimize:           promptOptimize,
	}
}

// ShapeText shapes the text based on the given prompts.
func (s *Shaper) ShapeText(inputFilePath, promptOrg, inputOrg string) (*ShapeResult, error) {
	if inputOrg == "" && !s.promptOptimize {
		rawResult, err := s.requestCreateChatCompletion(promptOrg)
		if err != nil {
			return nil, err
		}
		return NewShapeResult(promptOrg, rawResult, rawResult), nil
	}

	optimized := optimizePrompt(inputFilePath, promptOrg, inputOrg)
	rawResult, err := s.requestCreateChatCompletion(optimized)
	if err != nil {
		return nil, err
	}
	result, err := optimizeResponseResult(rawResult, s.useFirstCodeBlock)
	if err != nil {
		return nil, err
	}
	return NewShapeResult(optimized, rawResult, result), nil
}

func (s *Shaper) requestCreateChatCompletion(prompt string) (string, error) {
	var result string
	cr := s.gai.MakeCreateChatCompletion(prompt)
	maxCount := 1
	for i := 0; i < maxCount; i++ {
		comp, err := s.gai.RequestCreateChatCompletion(cr)
		if err != nil {
			return "", fmt.Errorf("failed to send chat message: %w", err)
		}

		if comp.Choices == nil || len(comp.Choices) == 0 {
			break
		}

		choice := comp.Choices[0]
		result += choice.Message.Content
		if choice.FinishReason != "length" {
			break
		}
	}

	return result, nil
}

func optimizePrompt(inputFilePath, prompt, input string) string {
	supplements := []string{
		"The subject of the Instruction is the area enclosed by the ai-text-shaper-input tag.",
		"The result should be returned in the language of the Instruction, but if the Instruction has a language specification, that language should be given priority.",
		"Only results should be returned and no explanation or supplementary information is required, but additional explanation or details should be provided if explicitly requested in the instructions.",
	}
	supplementation := strings.Join(supplements, " ")
	header := ""
	if inputFilePath != "" && inputFilePath != "-" {
		header = fmt.Sprintf("----%s----\n", inputFilePath)
	}
	return fmt.Sprintf("<Instruction>%s. (%s)</Instruction>\n%s<ai-text-shaper-input>\n%s\n</ai-text-shaper-input>", prompt, supplementation, header, input)
}

func optimizeResponseResult(rawResult string, useFirstCodeBlock bool) (string, error) {
	result := rawResult
	if strings.HasPrefix(result, "```") && strings.HasSuffix(result, "```") {
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
