// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"github.com/mattermost/mattermost/server/v8/channels/app/openai"
)

// FormattingRequest represents a request to format a message
type FormattingRequest struct {
	Message         string                `json:"message"`
	Profile         openai.FormattingProfile `json:"profile"`
	CustomInstructions string              `json:"custom_instructions,omitempty"`
}

// FormattingResponse represents the response from formatting
type FormattingResponse struct {
	FormattedText string                `json:"formatted_text"`
	Profile       openai.FormattingProfile `json:"profile"`
	Diff          *TextDiff              `json:"diff,omitempty"`
	ProcessingMs  int64                  `json:"processing_ms"`
}

// TextDiff represents the differences between original and formatted text
type TextDiff struct {
	Original   string   `json:"original"`
	Formatted  string   `json:"formatted"`
	Changes    []Change `json:"changes,omitempty"`
}

// Change represents a single change in the text
type Change struct {
	Type      string `json:"type"`      // "insert", "delete", "replace"
	Start     int    `json:"start"`     // Start position in original
	End       int    `json:"end"`       // End position in original
	NewText   string `json:"new_text"`  // New text (for insert/replace)
	OldText   string `json:"old_text"` // Old text (for delete/replace)
}

// FormattingProfileInfo represents metadata about a formatting profile
type FormattingProfileInfo struct {
	Id          string `json:"id"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

