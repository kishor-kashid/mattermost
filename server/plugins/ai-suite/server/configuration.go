package main

import (
	"fmt"
)

// configuration captures all user configurable settings for the plugin.
type configuration struct {
	OpenAIAPIKey        string `json:"openai_api_key"`
	OpenAIModel         string `json:"openai_model"`
	MaxSummaryMessages  int    `json:"max_summary_messages"`
	APIRateLimit        int    `json:"api_rate_limit"`
	RequestTimeoutSecs  int    `json:"request_timeout_secs"`
	EnableSummarization bool   `json:"enable_summarization"`
	EnableAnalytics     bool   `json:"enable_analytics"`
	EnableActionItems   bool   `json:"enable_action_items"`
	EnableFormatting    bool   `json:"enable_formatting"`
}

func newConfiguration() *configuration {
	return &configuration{
		OpenAIModel:         "gpt-3.5-turbo",
		MaxSummaryMessages:  500,
		APIRateLimit:        60,
		RequestTimeoutSecs:  30,
		EnableSummarization: true,
		EnableAnalytics:     true,
		EnableActionItems:   true,
		EnableFormatting:    true,
	}
}

func (c *configuration) Clone() *configuration {
	if c == nil {
		return newConfiguration()
	}

	cloned := *c
	return &cloned
}

func (c *configuration) ApplyDefaults() {
	if c.OpenAIModel == "" {
		c.OpenAIModel = "gpt-3.5-turbo"
	}

	if c.MaxSummaryMessages == 0 {
		c.MaxSummaryMessages = 500
	}

	if c.APIRateLimit == 0 {
		c.APIRateLimit = 60
	}

	if c.RequestTimeoutSecs == 0 {
		c.RequestTimeoutSecs = 30
	}
}

func (c *configuration) Validate() error {
	if c.MaxSummaryMessages < 100 || c.MaxSummaryMessages > 1000 {
		return fmt.Errorf("MaxSummaryMessages must be between 100 and 1000")
	}

	if c.APIRateLimit < 1 || c.APIRateLimit > 600 {
		return fmt.Errorf("APIRateLimit must be between 1 and 600")
	}

	if c.RequestTimeoutSecs < 5 || c.RequestTimeoutSecs > 60 {
		return fmt.Errorf("RequestTimeoutSecs must be between 5 and 60")
	}

	return nil
}
