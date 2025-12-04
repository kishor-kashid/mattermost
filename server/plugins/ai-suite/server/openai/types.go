package openai

import "fmt"

// MessageRole represents the speaker role in a chat completion request.
type MessageRole string

const (
	RoleSystem    MessageRole = "system"
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
)

// Message represents a single chat completion message.
type Message struct {
	Role    MessageRole `json:"role"`
	Content string      `json:"content"`
}

// ChatCompletionRequest is the payload sent to the OpenAI chat completions endpoint.
type ChatCompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float32   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	TopP        float32   `json:"top_p,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
	User        string    `json:"user,omitempty"`
}

// ChatCompletionResponse captures the subset of OpenAI fields we rely on.
type ChatCompletionResponse struct {
	ID      string               `json:"id"`
	Object  string               `json:"object"`
	Choices []ChatChoice         `json:"choices"`
	Usage   *ChatCompletionUsage `json:"usage,omitempty"`
}

// ChatChoice represents a single LLM response option.
type ChatChoice struct {
	Index        int     `json:"index"`
	FinishReason string  `json:"finish_reason"`
	Message      Message `json:"message"`
}

// ChatCompletionUsage contains token accounting metadata.
type ChatCompletionUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ErrorResponse represents an error payload returned by OpenAI.
type ErrorResponse struct {
	Error ErrorDetails `json:"error"`
}

// ErrorDetails contains metadata about an OpenAI API failure.
type ErrorDetails struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    string `json:"code"`
}

// APIError wraps error responses from OpenAI with HTTP metadata.
type APIError struct {
	StatusCode int
	Details    ErrorDetails
}

func (e *APIError) Error() string {
	if e == nil {
		return "openai: nil error"
	}

	if e.Details.Message != "" {
		return fmt.Sprintf("openai: %s (code=%s status=%d)", e.Details.Message, e.Details.Code, e.StatusCode)
	}

	return fmt.Sprintf("openai: unexpected status %d", e.StatusCode)
}
