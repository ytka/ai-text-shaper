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

// CreateChatCompletion represents the structure of a request to the OpenAI API Chat endpoint.
type CreateChatCompletion struct {
	Messages       []ChatMessage   `json:"messages"`
	Model          string          `json:"model"`
	ResponseFormat *ResponseFormat `json:"response_format,omitempty"`
	MaxTokens      *int            `json:"max_tokens,omitempty"`
	N              *int            `json:"n,omitempty"`
	Temperature    *float64        `json:"temperature,omitempty"`
	Seed           *int            `json:"seed,omitempty"`

	TopP             *float64           `json:"top_p,omitempty"`
	Stop             []string           `json:"stop,omitempty"`
	FrequencyPenalty *float64           `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]float64 `json:"logit_bias,omitempty"`
	User             *string            `json:"user,omitempty"`
	PresencePenalty  *float64           `json:"presence_penalty,omitempty"`
}

// ChatCompletion represents the JSON structure for the completion response
type ChatCompletion struct {
	ID                string         `json:"id"`
	Object            string         `json:"object"`
	Created           int            `json:"created"`
	Model             string         `json:"model"`
	ResponseFormat    ResponseFormat `json:"response_format"`
	SystemFingerprint string         `json:"system_fingerprint"`
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

type ResponseFormat struct {
	Type string `json:"type"`
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
	apikey   APIKey
	model    string
	logLevel string
}

func New(apikey APIKey, model string, logLevel string) *ChatClient {
	return &ChatClient{
		apikey:   apikey,
		model:    model,
		logLevel: logLevel,
	}
}

func (c *ChatClient) SendChatMessage(prompt string) (string, error) {
	resp, err := sendRawChatMessage(c.apikey, c.model, prompt, c.logLevel)
	if err != nil {
		return "", fmt.Errorf("failed to send chat message: %w", err)
	}

	var result string
	for _, choice := range resp.Choices {
		result += choice.Message.Content
	}
	return result, nil
}

func makeCreateChatCompletion(model, prompt string, responseFormatJSON bool) *CreateChatCompletion {
	n := 1
	seed := 0
	c := &CreateChatCompletion{
		Model: model,
		N:     &n,
		Seed:  &seed,
	}

	if responseFormatJSON {
		c.ResponseFormat = &ResponseFormat{Type: "json_object"}
		c.Messages = []ChatMessage{
			{Role: "system", Content: "You are a helpful assistant designed to output JSON."},
			{Role: "user", Content: prompt},
		}
	} else {
		c.Messages = []ChatMessage{
			{Role: "user", Content: prompt},
		}
	}
	return c
}

func sendRawChatMessage(apiKey APIKey, model, prompt string, logLevel string) (*ChatCompletion, error) {
	crq := makeCreateChatCompletion(model, prompt, false)
	requestBody, err := json.Marshal(crq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}
	switch logLevel {
	case "info":
		fmt.Printf("model: %s, N: %d, Seed: %d, ResponseFormat: %s\n", crq.Model, crq.N, crq.Seed, crq.ResponseFormat)
	case "debug":
		fmt.Printf("createChatCompletion: %s\n", requestBody)
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

	var comp ChatCompletion
	if err := json.Unmarshal(body, &comp); err != nil {
		return nil, err
	}
	switch logLevel {
	case "info":
		fmt.Printf("ID: %s, Object: %s, Created: %d, Model: %s, SystemFingerprint: %s, ChoicesCount:%d\n",
			comp.ID, comp.Object, comp.Created, comp.Model, comp.SystemFingerprint, len(comp.Choices))
		if len(comp.Choices) > 0 {
			fmt.Printf("[0]FinishReason: %s, Index: %d", comp.Choices[0].FinishReason, comp.Choices[0].Index)
		}
	case "debug":
		fmt.Printf("responseBody: %s\n", body)
	}
	return &comp, nil
}
