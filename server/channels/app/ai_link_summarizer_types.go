// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import "github.com/mattermost/mattermost/server/v8/channels/app/openai"

// LinkSummarizationRequest represents a request to summarize a link.
type LinkSummarizationRequest struct {
	UserId       string `json:"user_id"`
	URL          string `json:"url"`
	UseCache     bool   `json:"use_cache"`
	ForceRefresh bool   `json:"force_refresh"`
}

// LinkSummarizationResponse is returned after link summarization.
type LinkSummarizationResponse struct {
	Summary      *LinkSummary `json:"summary"`
	FromCache    bool         `json:"from_cache"`
	ProcessingMs int64        `json:"processing_ms"`
}

// LinkSummary is the app-layer representation before persistence.
type LinkSummary struct {
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Summary     string   `json:"summary"`
	KeyPoints   []string `json:"key_points,omitempty"`
	ContentType string   `json:"content_type,omitempty"`
	ReadingTime int      `json:"reading_time,omitempty"`
	Domain      string   `json:"domain,omitempty"`
	FaviconURL  string   `json:"favicon_url,omitempty"`
}

// LinkContent holds fetched and extracted page data.
type LinkContent struct {
	Title       string
	Description string
	Text        string
	ContentType string
	ReadingTime int
	Domain      string
	FaviconURL  string
}

// LinkSummaryLength maps to prompt verbosity.
type LinkSummaryLength string

const (
	LinkSummaryLengthShort    LinkSummaryLength = "short"
	LinkSummaryLengthStandard LinkSummaryLength = "standard"
	LinkSummaryLengthDetailed LinkSummaryLength = "detailed"
)

// ToPromptLength converts user preference to prompt enum.
func ToPromptLength(pref string) openai.LinkSummaryLength {
	switch LinkSummaryLength(pref) {
	case LinkSummaryLengthShort:
		return openai.LinkSummaryShort
	case LinkSummaryLengthDetailed:
		return openai.LinkSummaryDetailed
	default:
		return openai.LinkSummaryStandard
	}
}
