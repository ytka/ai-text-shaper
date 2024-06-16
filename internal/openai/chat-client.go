package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type APIKey string

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Response struct {
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type ChatClient struct {
	apikey APIKey
	model  string
}

func New(apikey APIKey, model string) *ChatClient {
	return &ChatClient{
		apikey: apikey,
		model:  model,
	}
}

func (c *ChatClient) SendChatMessage(prompt string) (string, error) {
	resp, err := sendRawChatMessage(c.apikey, c.model, prompt)
	if err != nil {
		return "", fmt.Errorf("failed to send chat message: %w", err)
	}

	var result string
	for _, choice := range resp.Choices {
		result += choice.Message.Content
	}
	result = strings.TrimSuffix(result, "\n")
	result = strings.TrimSpace(result)
	return result, nil
}

func sendRawChatMessage(apiKey APIKey, model, prompt string) (*Response, error) {
	requestBody, err := json.Marshal(Request{
		Model:    model,
		Messages: []Message{{Role: "user", Content: prompt}},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var openAIResponse Response
	err = json.Unmarshal(body, &openAIResponse)
	if err != nil {
		return nil, err
	}
	return &openAIResponse, nil
}
