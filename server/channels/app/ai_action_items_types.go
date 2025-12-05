// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"time"

	"github.com/mattermost/mattermost/server/public/model"
)

// ActionItemDetectionResult represents the result of AI action item detection
type ActionItemDetectionResult struct {
	Items     []*model.AIActionItem
	Detected  bool
	Confidence float64
}

// ActionItemFilters represents filters for querying action items
type ActionItemFilters struct {
	UserID      string
	ChannelID   string
	TeamID      string
	Status      string
	Priority    string
	DueBefore   *time.Time
	DueAfter    *time.Time
	AssignedBy  string
	IncludeCompleted bool
	Page        int
	PerPage     int
}

// ActionItemStats represents statistics about action items
type ActionItemStats struct {
	Total       int
	Overdue     int
	DueToday    int
	DueSoon     int // Next 7 days
	NoDueDate   int
	Completed   int
	ByPriority  map[string]int
	ByStatus    map[string]int
}

// ActionItemBatch represents a batch operation on action items
type ActionItemBatch struct {
	IDs    []string
	Status string
}

// ActionItemCreateRequest represents a request to create an action item
type ActionItemCreateRequest struct {
	Description string
	AssigneeID  string
	ChannelID   string
	PostID      string
	DueDate     *time.Time
	Priority    string
	Status      string
}

// ActionItemUpdateRequest represents a request to update an action item
type ActionItemUpdateRequest struct {
	Description *string
	AssigneeID  *string
	DueDate     *time.Time
	Priority    *string
	Status      *string
	CompletedAt *time.Time
}

