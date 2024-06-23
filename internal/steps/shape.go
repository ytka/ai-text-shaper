package steps

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ytka/ai-text-shaper/internal/openai"
)

// GenerativeAIClient represents an interface for generating AI client operations.
type GenerativeAIClient interface {
	RequestCreateChatCompletion(*openai.CreateChatCompletion) (*openai.ChatCompletion, error)
	MakeCreateChatCompletion(prompt string) *openai.CreateChatCompletion
}

// reCodeBlock is a regular expression to find code blocks in markdown.
var reCodeBlock = regexp.MustCompile("(?s)```[a-zA-Z0-9]*?\n(.*?\n)```")

type ShapePrompt string

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

// MakeShapePrompt generates a ShapePrompt based on input parameters.
func (s *Shaper) MakeShapePrompt(inputFilePath, promptOrg, inputOrg string) ShapePrompt {
	if inputOrg == "" && !s.promptOptimize {
		return ShapePrompt(promptOrg)
	}
	return ShapePrompt(optimizePrompt(inputFilePath, promptOrg, inputOrg))
}

// Shape shapes the text based on the given prompts.
func (s *Shaper) Shape(prompt ShapePrompt) (*ShapeResult, error) {
	rawResult, err := s.requestCreateChatCompletion(string(prompt))
	if err != nil {
		return nil, err
	}

	return NewShapeResult(string(prompt), rawResult, optimizeResponseResult(rawResult, s.useFirstCodeBlock)), nil
}

// requestCreateChatCompletion requests the AI to create chat completion based on the given prompt.
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

// optimizePrompt refines the prompt by incorporating additional information.
func optimizePrompt(inputFilePath, prompt, input string) string {
	supplements := []string{
		"The subject of the Instruction is the area enclosed by the ai-text-shaper-input tag.",
		"The result should be returned in the language of the Instruction, but if the Instruction has a language specification, that language should be given priority.",
		"Only results should be returned and no explanation or supplementary information is required, but additional explanation or details should be provided if explicitly requested in the instructions.",
	}
	supplementation := strings.Join(supplements, " ")
	header := ""
	if inputFilePath != "" && inputFilePath != "-" {
		header = fmt.Sprintf("filepath=\"%s\"\n", inputFilePath)
	}
	return fmt.Sprintf("<Instruction>%s. (%s)</Instruction>\n%s<ai-text-shaper-input>\n%s\n<ai-text-shaper-input>", prompt, supplementation, header, input)
}

// optimizeResponseResult refines the AI's response, potentially extracting code blocks.
func optimizeResponseResult(rawResult string, useFirstCodeBlock bool) string {
	result := rawResult
	if strings.HasPrefix(result, "```") && strings.HasSuffix(result, "```") {
		lines := strings.Split(result, "\n")
		if len(lines) > 2 {
			result = strings.Join(lines[1:len(lines)-1], "\n")
		}
	}

	if useFirstCodeBlock {
		if codeBlock := findMarkdownFirstCodeBlock(result); codeBlock != "" {
			result = codeBlock
		}
	}
	return result
}

// findMarkdownFirstCodeBlock extracts the first code block found in text.
func findMarkdownFirstCodeBlock(text string) string {
	match := reCodeBlock.FindStringSubmatch(text)
	if match != nil {
		return match[1]
	}
	return ""
}
