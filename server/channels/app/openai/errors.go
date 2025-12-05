// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package openai

import (
	"fmt"
)

// ClientError represents an error from the OpenAI client
type ClientError struct {
	Message    string
	StatusCode int
	RetryAfter int
}

func (e *ClientError) Error() string {
	if e.StatusCode > 0 {
		return fmt.Sprintf("OpenAI API error (status %d): %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("OpenAI client error: %s", e.Message)
}

// IsRateLimitError checks if the error is a rate limit error
func IsRateLimitError(err error) bool {
	if clientErr, ok := err.(*ClientError); ok {
		return clientErr.StatusCode == 429
	}
	return false
}

// IsAuthenticationError checks if the error is an authentication error
func IsAuthenticationError(err error) bool {
	if clientErr, ok := err.(*ClientError); ok {
		return clientErr.StatusCode == 401
	}
	return false
}

// IsServerError checks if the error is a server error (5xx)
func IsServerError(err error) bool {
	if clientErr, ok := err.(*ClientError); ok {
		return clientErr.StatusCode >= 500 && clientErr.StatusCode < 600
	}
	return false
}

