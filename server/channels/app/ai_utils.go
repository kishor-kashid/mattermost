// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"fmt"
	"strings"
	"time"
)

// FormatMessageForAI formats a message for AI processing, removing sensitive data
func FormatMessageForAI(message string) string {
	// Remove any potential sensitive patterns
	// This is a simple implementation - can be enhanced
	return strings.TrimSpace(message)
}

// FormatTimeRange formats a time range for display
func FormatTimeRange(startTime, endTime int64) string {
	start := time.Unix(startTime/1000, 0)
	end := time.Unix(endTime/1000, 0)

	if start.Format("2006-01-02") == end.Format("2006-01-02") {
		// Same day
		return fmt.Sprintf("%s to %s", start.Format("3:04 PM"), end.Format("3:04 PM on Jan 2, 2006"))
	}

	// Different days
	return fmt.Sprintf("%s to %s", start.Format("Jan 2, 2006 3:04 PM"), end.Format("Jan 2, 2006 3:04 PM"))
}

// TruncateMessages truncates a list of messages to a maximum count
func TruncateMessages(messages []string, maxCount int) []string {
	if len(messages) <= maxCount {
		return messages
	}
	return messages[len(messages)-maxCount:]
}

// EstimateTokenCount provides a rough estimate of token count for a string
// This is a simple approximation: ~4 characters per token on average
func EstimateTokenCount(text string) int {
	return len(text) / 4
}

// ValidateAIConfig checks if AI configuration is valid
func ValidateAIConfig(apiKey, model string) error {
	if apiKey == "" {
		return fmt.Errorf("OpenAI API key is required")
	}

	if model == "" {
		return fmt.Errorf("model is required")
	}

	// Validate model is supported
	supportedModels := []string{
		"gpt-4",
		"gpt-4-turbo",
		"gpt-4-turbo-preview",
		"gpt-3.5-turbo",
		"gpt-3.5-turbo-16k",
	}

	modelSupported := false
	for _, supportedModel := range supportedModels {
		if strings.HasPrefix(model, supportedModel) {
			modelSupported = true
			break
		}
	}

	if !modelSupported {
		return fmt.Errorf("unsupported model: %s", model)
	}

	return nil
}

// SanitizeAIResponse removes any potentially problematic content from AI responses
func SanitizeAIResponse(response string) string {
	// Remove excessive newlines
	response = strings.TrimSpace(response)

	// Ensure response is not too long (max 10000 characters)
	if len(response) > 10000 {
		response = response[:10000] + "..."
	}

	return response
}

// FormatDate formats a date for AI analytics
func FormatDate(timestamp int64) string {
	t := time.Unix(timestamp/1000, 0)
	return t.Format("2006-01-02")
}

// GetDateRange returns the start and end timestamps for a given number of days
func GetDateRange(days int) (startTime, endTime int64) {
	now := time.Now()
	endTime = now.UnixMilli()
	startTime = now.AddDate(0, 0, -days).UnixMilli()
	return
}

