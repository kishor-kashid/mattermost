package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

const (
	defaultBaseURL    = "https://api.openai.com/v1"
	defaultMaxRetries = 3
)

// ClientConfig controls how the OpenAI client operates.
type ClientConfig struct {
	APIKey               string
	Model                string
	BaseURL              string
	RequestTimeout       time.Duration
	MaxRequestsPerMinute int
	MaxRetries           int
	HTTPClient           *http.Client
}

// Client wraps basic chat completion functionality with logging-friendly errors,
// rate limiting, and retry behavior.
type Client struct {
	cfg        ClientConfig
	httpClient *http.Client
	limiter    *rate.Limiter
}

// NewClient constructs a new Client instance. It returns nil if no API key is configured.
func NewClient(cfg ClientConfig) (*Client, error) {
	if cfg.APIKey == "" {
		return nil, errors.New("openai: api key is required")
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = defaultBaseURL
	}

	if cfg.RequestTimeout <= 0 {
		cfg.RequestTimeout = 30 * time.Second
	}

	if cfg.MaxRetries <= 0 {
		cfg.MaxRetries = defaultMaxRetries
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: cfg.RequestTimeout,
		}
	}

	var limiter *rate.Limiter
	if cfg.MaxRequestsPerMinute > 0 {
		every := time.Minute / time.Duration(cfg.MaxRequestsPerMinute)
		limiter = rate.NewLimiter(rate.Every(every), cfg.MaxRequestsPerMinute)
	}

	return &Client{
		cfg:        cfg,
		httpClient: httpClient,
		limiter:    limiter,
	}, nil
}

// ChatCompletion sends a chat completion request to OpenAI and returns the response body.
func (c *Client) ChatCompletion(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error) {
	if c == nil {
		return nil, errors.New("openai: client is not configured")
	}

	if strings.TrimSpace(req.Model) == "" {
		req.Model = c.cfg.Model
	}

	if len(req.Messages) == 0 {
		return nil, errors.New("openai: at least one message is required")
	}

	if c.limiter != nil {
		if err := c.limiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("openai: rate limiter wait failed: %w", err)
		}
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("openai: marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/chat/completions", strings.TrimSuffix(c.cfg.BaseURL, "/"))
	var lastErr error

	for attempt := 0; attempt < c.cfg.MaxRetries; attempt++ {
		resp, err := c.doRequest(ctx, url, payload)
		if err == nil {
			return resp, nil
		}

		lastErr = err

		var apiErr *APIError
		if errors.As(err, &apiErr) {
			// Retry on 429 and 5xx responses only.
			if apiErr.StatusCode != http.StatusTooManyRequests && apiErr.StatusCode < 500 {
				break
			}
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(backoffDuration(attempt)):
		}
	}

	return nil, lastErr
}

func (c *Client) doRequest(ctx context.Context, url string, payload []byte) (*ChatCompletionResponse, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("openai: build request: %w", err)
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.cfg.APIKey))
	request.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("openai: http request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("openai: read response: %w", err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var completion ChatCompletionResponse
		if err = json.Unmarshal(body, &completion); err != nil {
			return nil, fmt.Errorf("openai: decode response: %w", err)
		}
		return &completion, nil
	}

	var apiErrResp ErrorResponse
	if err = json.Unmarshal(body, &apiErrResp); err != nil {
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Details: ErrorDetails{
				Message: string(body),
			},
		}
	}

	return nil, &APIError{
		StatusCode: resp.StatusCode,
		Details:    apiErrResp.Error,
	}
}

func backoffDuration(attempt int) time.Duration {
	base := 250 * time.Millisecond
	return time.Duration(1<<attempt) * base
}
