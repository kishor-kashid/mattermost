// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/mattermost/mattermost/server/public/shared/mlog"
)

const (
	defaultBaseURL     = "https://api.openai.com/v1"
	defaultTimeout     = 30 * time.Second
	defaultMaxRetries  = 3
	defaultRetryDelay  = 1 * time.Second
)

// Client represents an OpenAI API client
type Client struct {
	apiKey      string
	baseURL     string
	httpClient  *http.Client
	logger      mlog.LoggerIFace
	maxRetries  int
	retryDelay  time.Duration
}

// ClientConfig holds configuration for the OpenAI client
type ClientConfig struct {
	APIKey     string
	BaseURL    string
	Timeout    time.Duration
	Logger     mlog.LoggerIFace
	MaxRetries int
	RetryDelay time.Duration
}

// NewClient creates a new OpenAI client
func NewClient(config ClientConfig) *Client {
	if config.BaseURL == "" {
		config.BaseURL = defaultBaseURL
	}
	if config.Timeout == 0 {
		config.Timeout = defaultTimeout
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = defaultMaxRetries
	}
	if config.RetryDelay == 0 {
		config.RetryDelay = defaultRetryDelay
	}

	return &Client{
		apiKey:  config.APIKey,
		baseURL: config.BaseURL,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		logger:     config.Logger,
		maxRetries: config.MaxRetries,
		retryDelay: config.RetryDelay,
	}
}

// CreateChatCompletion creates a chat completion request
func (c *Client) CreateChatCompletion(ctx context.Context, request ChatCompletionRequest) (*ChatCompletionResponse, error) {
	url := fmt.Sprintf("%s/chat/completions", c.baseURL)

	var lastErr error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			delay := c.retryDelay * time.Duration(1<<uint(attempt-1))
			if c.logger != nil {
				c.logger.Debug("Retrying OpenAI request", mlog.Int("attempt", attempt), mlog.Duration("delay", delay))
			}
			time.Sleep(delay)
		}

		response, err := c.doRequest(ctx, url, request)
		if err == nil {
			return response, nil
		}

		lastErr = err

		// Don't retry on authentication errors or client errors (4xx except 429)
		if clientErr, ok := err.(*ClientError); ok {
			if clientErr.StatusCode == 401 || (clientErr.StatusCode >= 400 && clientErr.StatusCode < 500 && clientErr.StatusCode != 429) {
				return nil, err
			}
		}
	}

	return nil, fmt.Errorf("failed after %d retries: %w", c.maxRetries, lastErr)
}

func (c *Client) doRequest(ctx context.Context, url string, request ChatCompletionRequest) (*ChatCompletionResponse, error) {
	// Marshal request body
	bodyBytes, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	// Log request (without sensitive data)
	if c.logger != nil {
		c.logger.Debug("OpenAI API request",
			mlog.String("model", request.Model),
			mlog.Int("message_count", len(request.Messages)),
		)
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &ClientError{
			Message: fmt.Sprintf("request failed: %v", err),
		}
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &ClientError{
			Message: fmt.Sprintf("failed to read response: %v", err),
		}
	}

	// Handle non-200 responses
	if resp.StatusCode != http.StatusOK {
		var apiErr OpenAIError
		if err := json.Unmarshal(respBody, &apiErr); err == nil && apiErr.Error.Message != "" {
			clientErr := &ClientError{
				Message:    apiErr.Error.Message,
				StatusCode: resp.StatusCode,
			}

			// Parse Retry-After header for rate limiting
			if resp.StatusCode == 429 {
				if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
					if seconds, err := strconv.Atoi(retryAfter); err == nil {
						clientErr.RetryAfter = seconds
					}
				}
			}

			return nil, clientErr
		}

		return nil, &ClientError{
			Message:    string(respBody),
			StatusCode: resp.StatusCode,
		}
	}

	// Parse success response
	var completion ChatCompletionResponse
	if err := json.Unmarshal(respBody, &completion); err != nil {
		return nil, &ClientError{
			Message: fmt.Sprintf("failed to parse response: %v", err),
		}
	}

	// Log response
	if c.logger != nil {
		c.logger.Debug("OpenAI API response",
			mlog.String("id", completion.ID),
			mlog.String("model", completion.Model),
			mlog.Int("prompt_tokens", completion.Usage.PromptTokens),
			mlog.Int("completion_tokens", completion.Usage.CompletionTokens),
			mlog.Int("total_tokens", completion.Usage.TotalTokens),
		)
	}

	return &completion, nil
}

// Simple convenience method to create a chat completion with a simple prompt
func (c *Client) SimpleCompletion(ctx context.Context, model, systemPrompt, userPrompt string) (string, error) {
	messages := []ChatCompletionMessage{
		{
			Role:    "user",
			Content: userPrompt,
		},
	}

	if systemPrompt != "" {
		messages = []ChatCompletionMessage{
			{
				Role:    "system",
				Content: systemPrompt,
			},
			{
				Role:    "user",
				Content: userPrompt,
			},
		}
	}

	request := ChatCompletionRequest{
		Model:    model,
		Messages: messages,
	}

	response, err := c.CreateChatCompletion(ctx, request)
	if err != nil {
		return "", err
	}

	if len(response.Choices) == 0 {
		return "", &ClientError{
			Message: "no choices returned from API",
		}
	}

	return response.Choices[0].Message.Content, nil
}

