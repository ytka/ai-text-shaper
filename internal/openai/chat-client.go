package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var ErrUnexpectedStatusCode = errors.New("unexpected status code")

type ChatClient struct {
	apikey    APIKey
	model     string
	logLevel  string
	maxTokens *int
}

func New(apikey APIKey, model, logLevel string, maxTokens *int) *ChatClient {
	return &ChatClient{
		apikey:    apikey,
		model:     model,
		logLevel:  logLevel,
		maxTokens: maxTokens,
	}
}

func (c *ChatClient) MakeCreateChatCompletion(prompt string) *CreateChatCompletion {
	return newCreateChatCompletion(c.model, prompt, c.maxTokens, false)
}

func (c *ChatClient) sendChatCompletionsRequest(ccc *CreateChatCompletion) (*http.Response, error) {
	requestBody, err := json.Marshal(ccc)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}
	switch c.logLevel {
	case "info":
		fmt.Printf("model: %s, N: %d, Seed: %d, ResponseFormat: %s\n", ccc.Model, ccc.N, ccc.Seed, ccc.ResponseFormat)
	case "debug":
		fmt.Printf("createChatCompletion: %s\n", requestBody)
	}

	cnt := context.Background()
	req, err := http.NewRequestWithContext(cnt, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apikey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

func (c *ChatClient) makeCatCompletions(body []byte) (*ChatCompletion, error) {
	var comp ChatCompletion
	if err := json.Unmarshal(body, &comp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	switch c.logLevel {
	case "info":
		fmt.Printf("ID: %s, Object: %s, Created: %d, Model: %s, SystemFingerprint: %s, ChoicesCount:%d\n",
			comp.ID, comp.Object, comp.Created, comp.Model, comp.SystemFingerprint, len(comp.Choices))
		if len(comp.Choices) > 0 {
			fmt.Printf("[0]FinishReason: %s, Index: %d\n", comp.Choices[0].FinishReason, comp.Choices[0].Index)
		}
	case "debug":
		fmt.Printf("responseBody: %s\n", body)
	}

	return &comp, nil
}

func (c *ChatClient) RequestCreateChatCompletion(ccc *CreateChatCompletion) (*ChatCompletion, error) {
	resp, err := c.sendChatCompletionsRequest(ccc)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			fmt.Printf("failed to close response body: %s\n", cerr)
		}
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode > 299 {
		var errorResponse ErrorResponse
		if err := json.Unmarshal(respBody, &errorResponse); err != nil {
			return nil, fmt.Errorf("failed to unmarshal error response: %w", err)
		}
		return nil, fmt.Errorf("%w: %d '%s'", ErrUnexpectedStatusCode, resp.StatusCode, errorResponse.Error.Message)
	}

	return c.makeCatCompletions(respBody)
}
