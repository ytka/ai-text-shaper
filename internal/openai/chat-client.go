package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type APIKey string

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents the structure of a request to the OpenAI API Chat endpoint
type ChatRequest struct {
	Messages       []ChatMessage `json:"messages"`
	Model          string        `json:"model"`
	MaxTokens      int           `json:"max_tokens,omitempty"`
	N              int           `json:"n,omitempty"`
	ResponseFormat string        `json:"response_format,omitempty"`
	Temperature    float64       `json:"temperature,omitempty"`
	Seed           int           `json:"seed,omitempty"`

	TopP             float64            `json:"top_p,omitempty"`
	Stop             []string           `json:"stop,omitempty"`
	FrequencyPenalty float64            `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]float64 `json:"logit_bias,omitempty"`
	User             string             `json:"user,omitempty"`
	PresencePenalty  float64            `json:"presence_penalty,omitempty"`
}

type ChatCompletion struct {
	ID                string `json:"id"`
	Object            string `json:"object"`
	Created           int    `json:"created"`
	Model             string `json:"model"`
	SystemFingerprint string `json:"system_fingerprint"`
	Choices           []struct {
		FinishReason string      `json:"finish_reason"`
		Index        int         `json:"index"`
		Message      ChatMessage `json:"message"`
		// Logprobs     interface{} `json:"logprobs"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// ErrorResponse represents the JSON structure for the error response
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail represents the details of the error
type ErrorDetail struct {
	Message string      `json:"message"`
	Type    string      `json:"type"`
	Param   interface{} `json:"param"` // Param can be of any type, hence using interface{}
	Code    string      `json:"code"`
}

type ChatClient struct {
	apikey          APIKey
	model           string
	printAPIMessage bool
}

func New(apikey APIKey, model string, printAPIMessage bool) *ChatClient {
	return &ChatClient{
		apikey:          apikey,
		model:           model,
		printAPIMessage: printAPIMessage,
	}
}

func (c *ChatClient) SendChatMessage(prompt string) (string, error) {
	resp, err := sendRawChatMessage(c.apikey, c.model, prompt, c.printAPIMessage)
	if err != nil {
		return "", fmt.Errorf("failed to send chat message: %w", err)
	}

	var result string
	for _, choice := range resp.Choices {
		result += choice.Message.Content
	}
	return result, nil
}

func sendRawChatMessage(apiKey APIKey, model, prompt string, printAPIMessage bool) (*ChatCompletion, error) {
	requestBody, err := json.Marshal(ChatRequest{
		Model:    model,
		Messages: []ChatMessage{{Role: "user", Content: prompt}},
		N:        1,
		Seed:     0,
		// ResponseFormat
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}
	if printAPIMessage {
		fmt.Printf("requestBody: %s\n", requestBody)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if printAPIMessage {
		fmt.Printf("responseBody: %s\n", body)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	if resp.StatusCode > 299 {
		var errorResponse ErrorResponse
		if err := json.Unmarshal(body, &errorResponse); err != nil {
			return nil, fmt.Errorf("unexpected status code: %d %s", resp.StatusCode, body)
		}
		return nil, fmt.Errorf("unexpected status code: %d '%s'", resp.StatusCode, errorResponse.Error.Message)
	}

	var jsonResponse ChatCompletion
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		return nil, err
	}
	return &jsonResponse, nil
}
