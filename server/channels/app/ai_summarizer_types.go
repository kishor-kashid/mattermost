// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"github.com/mattermost/mattermost/server/public/model"
)

// SummarizationRequest represents a request to summarize messages
type SummarizationRequest struct {
	ChannelId      string
	PostId         string // For thread summarization
	StartTime      int64  // For channel summarization
	EndTime        int64  // For channel summarization
	SummaryLevel   string // brief, standard, detailed
	MaxMessages    int
	UserId         string // User requesting the summary
	UseCache       bool   // Whether to use cached summaries
}

// SummarizationResponse represents the result of a summarization
type SummarizationResponse struct {
	Summary      *model.AISummary
	FromCache    bool
	TokensUsed   int
	ProcessingMs int64
}

// MessageContext represents a formatted message for LLM prompts
type MessageContext struct {
	Author    string
	Username  string
	Timestamp int64
	Content   string
}

// SummaryMetadata contains metadata about the summarization
type SummaryMetadata struct {
	MessageCount   int
	Participants   []string
	StartTime      int64
	EndTime        int64
	ChannelName    string
	IsThread       bool
	ThreadRootText string
}

// SummarizationOptions contains options for the summarization process
type SummarizationOptions struct {
	Level              string
	MaxMessages        int
	IncludeParticipants bool
	IncludeTimestamps   bool
	ResolveUsernames    bool
}

