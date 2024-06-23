package openai

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

// ErrorResponse represents the JSON structure for the error response.
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail represents the details of the error.
type ErrorDetail struct {
	Message string      `json:"message"`
	Type    string      `json:"type"`
	Param   interface{} `json:"param"` // Param can be of any type, hence using interface{}.
	Code    string      `json:"code"`
}

func newCreateChatCompletion(model, prompt string, maxTokens *int, responseFormatJSON bool) *CreateChatCompletion {
	n := 1
	seed := 0
	cr := &CreateChatCompletion{
		Model:     model,
		N:         &n,
		Seed:      &seed,
		MaxTokens: maxTokens,
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
