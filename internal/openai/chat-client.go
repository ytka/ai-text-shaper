package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ChatClient struct {
	apikey    APIKey
	model     string
	logLevel  string
	maxTokens *int
}

func New(apikey APIKey, model string, logLevel string, maxTokens *int) *ChatClient {
	return &ChatClient{
		apikey:    apikey,
		model:     model,
		logLevel:  logLevel,
		maxTokens: maxTokens,
	}
}

func (c *ChatClient) makeCreateChatCompletion(prompt string, responseFormatJSON bool) *CreateChatCompletion {
	n := 1
	seed := 0
	cr := &CreateChatCompletion{
		Model: c.model,
		N:     &n,
		Seed:  &seed,
	}
	if c.maxTokens != nil {
		cr.MaxTokens = c.maxTokens
	}

	if responseFormatJSON {
		cr.ResponseFormat = &ResponseFormat{Type: "json_object"}
		cr.Messages = []ChatMessage{
			{Role: "system", Content: "You are a helpful assistant designed to output JSON."},
			{Role: "user", Content: prompt},
		}
	} else {
		cr.Messages = []ChatMessage{
			{Role: "user", Content: prompt},
		}
	}
	return cr
}

func (c *ChatClient) SendChatMessage(prompt string) (*ChatCompletion, error) {
	crq := c.makeCreateChatCompletion(prompt, false)
	requestBody, err := json.Marshal(crq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}
	switch c.logLevel {
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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apikey))

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
	switch c.logLevel {
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
