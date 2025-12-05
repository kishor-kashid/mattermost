// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package model

import (
	"encoding/json"
	"net/http"
)

const (
	AIActionItemStatusOpen       = "open"
	AIActionItemStatusInProgress = "in_progress"
	AIActionItemStatusCompleted  = "completed"
	AIActionItemStatusDismissed  = "dismissed"

	AISummaryTypeThread  = "thread"
	AISummaryTypeChannel = "channel"
)

// AIActionItem represents an AI-detected action item from a message
type AIActionItem struct {
	Id          string `json:"id" db:"id"`
	ChannelId   string `json:"channel_id" db:"channelid"`
	PostId      string `json:"post_id,omitempty" db:"postid"`
	CreatedBy   string `json:"created_by" db:"createdby"`
	AssigneeId  string `json:"assignee_id,omitempty" db:"assigneeid"`
	Description string `json:"description" db:"description"`
	DueDate     int64  `json:"due_date,omitempty" db:"duedate"`
	Priority    string `json:"priority" db:"priority"`
	Status      string `json:"status" db:"status"`
	CompletedAt int64  `json:"completed_at,omitempty" db:"completedat"`
	CreateAt    int64  `json:"create_at" db:"createdat"`
	UpdateAt    int64  `json:"update_at" db:"updatedat"`
	DeleteAt    int64  `json:"delete_at" db:"deletedat"`
}

func (a *AIActionItem) IsValid() *AppError {
	if !IsValidId(a.Id) {
		return NewAppError("AIActionItem.IsValid", "model.ai_action_item.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}

	if !IsValidId(a.ChannelId) {
		return NewAppError("AIActionItem.IsValid", "model.ai_action_item.is_valid.channel_id.app_error", nil, "", http.StatusBadRequest)
	}

	if !IsValidId(a.CreatedBy) {
		return NewAppError("AIActionItem.IsValid", "model.ai_action_item.is_valid.created_by.app_error", nil, "", http.StatusBadRequest)
	}

	if a.AssigneeId != "" && !IsValidId(a.AssigneeId) {
		return NewAppError("AIActionItem.IsValid", "model.ai_action_item.is_valid.assignee_id.app_error", nil, "", http.StatusBadRequest)
	}

	if a.Description == "" {
		return NewAppError("AIActionItem.IsValid", "model.ai_action_item.is_valid.description.app_error", nil, "", http.StatusBadRequest)
	}

	if a.Status != "open" && a.Status != "in_progress" && a.Status != "completed" && a.Status != "dismissed" {
		return NewAppError("AIActionItem.IsValid", "model.ai_action_item.is_valid.status.app_error", nil, "", http.StatusBadRequest)
	}
	
	if a.Priority != "" && a.Priority != "low" && a.Priority != "medium" && a.Priority != "high" && a.Priority != "urgent" {
		return NewAppError("AIActionItem.IsValid", "model.ai_action_item.is_valid.priority.app_error", nil, "", http.StatusBadRequest)
	}

	if a.CreateAt == 0 {
		return NewAppError("AIActionItem.IsValid", "model.ai_action_item.is_valid.create_at.app_error", nil, "", http.StatusBadRequest)
	}

	if a.UpdateAt == 0 {
		return NewAppError("AIActionItem.IsValid", "model.ai_action_item.is_valid.update_at.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func (a *AIActionItem) PreSave() {
	if a.Id == "" {
		a.Id = NewId()
	}

	if a.Status == "" {
		a.Status = AIActionItemStatusOpen
	}

	a.CreateAt = GetMillis()
	a.UpdateAt = a.CreateAt
	a.DeleteAt = 0
}

func (a *AIActionItem) PreUpdate() {
	a.UpdateAt = GetMillis()
}

// AISummary represents an AI-generated summary of messages
type AISummary struct {
	Id           string `json:"id"`
	ChannelId    string `json:"channel_id"`
	PostId       string `json:"post_id,omitempty"`
	SummaryType  string `json:"summary_type"`
	Summary      string `json:"summary"`
	MessageCount int    `json:"message_count"`
	StartTime    int64  `json:"start_time"`
	EndTime      int64  `json:"end_time"`
	UserId       string `json:"user_id"`
	Participants string `json:"participants,omitempty"`
	CacheKey     string `json:"cache_key,omitempty"`
	ChannelName  string `json:"channel_name,omitempty"`
	CreateAt     int64  `json:"create_at"`
	ExpiresAt    int64  `json:"expires_at"`
}

func (s *AISummary) IsValid() *AppError {
	if !IsValidId(s.Id) {
		return NewAppError("AISummary.IsValid", "model.ai_summary.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}

	if !IsValidId(s.ChannelId) {
		return NewAppError("AISummary.IsValid", "model.ai_summary.is_valid.channel_id.app_error", nil, "", http.StatusBadRequest)
	}

	if !IsValidId(s.UserId) {
		return NewAppError("AISummary.IsValid", "model.ai_summary.is_valid.user_id.app_error", nil, "", http.StatusBadRequest)
	}

	if s.SummaryType != AISummaryTypeThread && s.SummaryType != AISummaryTypeChannel {
		return NewAppError("AISummary.IsValid", "model.ai_summary.is_valid.summary_type.app_error", nil, "", http.StatusBadRequest)
	}

	if s.Summary == "" {
		return NewAppError("AISummary.IsValid", "model.ai_summary.is_valid.summary.app_error", nil, "", http.StatusBadRequest)
	}

	if s.MessageCount <= 0 {
		return NewAppError("AISummary.IsValid", "model.ai_summary.is_valid.message_count.app_error", nil, "", http.StatusBadRequest)
	}

	if s.CreateAt == 0 {
		return NewAppError("AISummary.IsValid", "model.ai_summary.is_valid.create_at.app_error", nil, "", http.StatusBadRequest)
	}

	if s.ExpiresAt == 0 {
		return NewAppError("AISummary.IsValid", "model.ai_summary.is_valid.expires_at.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func (s *AISummary) PreSave() {
	if s.Id == "" {
		s.Id = NewId()
	}

	s.CreateAt = GetMillis()

	// Default expiry is 24 hours from creation
	if s.ExpiresAt == 0 {
		s.ExpiresAt = s.CreateAt + (24 * 60 * 60 * 1000)
	}
}

// AIAnalytics represents aggregated analytics data for a channel
type AIAnalytics struct {
	Id                  string                 `json:"id"`
	ChannelId           string                 `json:"channel_id"`
	Date                string                 `json:"date"` // Format: YYYY-MM-DD
	MessageCount        int                    `json:"message_count"`
	UserCount           int                    `json:"user_count"`
	AvgResponseTime     int64                  `json:"avg_response_time"`
	TopContributors     map[string]interface{} `json:"top_contributors"`
	HourlyDistribution  map[string]interface{} `json:"hourly_distribution"`
	CreateAt            int64                  `json:"create_at"`
	UpdateAt            int64                  `json:"update_at"`
}

func (a *AIAnalytics) IsValid() *AppError {
	if !IsValidId(a.Id) {
		return NewAppError("AIAnalytics.IsValid", "model.ai_analytics.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}

	if !IsValidId(a.ChannelId) {
		return NewAppError("AIAnalytics.IsValid", "model.ai_analytics.is_valid.channel_id.app_error", nil, "", http.StatusBadRequest)
	}

	if a.Date == "" {
		return NewAppError("AIAnalytics.IsValid", "model.ai_analytics.is_valid.date.app_error", nil, "", http.StatusBadRequest)
	}

	if a.CreateAt == 0 {
		return NewAppError("AIAnalytics.IsValid", "model.ai_analytics.is_valid.create_at.app_error", nil, "", http.StatusBadRequest)
	}

	if a.UpdateAt == 0 {
		return NewAppError("AIAnalytics.IsValid", "model.ai_analytics.is_valid.update_at.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func (a *AIAnalytics) PreSave() {
	if a.Id == "" {
		a.Id = NewId()
	}

	a.CreateAt = GetMillis()
	a.UpdateAt = a.CreateAt

	if a.TopContributors == nil {
		a.TopContributors = make(map[string]interface{})
	}

	if a.HourlyDistribution == nil {
		a.HourlyDistribution = make(map[string]interface{})
	}
}

func (a *AIAnalytics) PreUpdate() {
	a.UpdateAt = GetMillis()
}

// ToJSON converts AIAnalytics to JSON string
func (a *AIAnalytics) ToJSON() string {
	b, _ := json.Marshal(a)
	return string(b)
}

// AIPreferences represents user preferences for AI features
type AIPreferences struct {
	Id                  string `json:"id"`
	UserId              string `json:"user_id"`
	EnableSummarization bool   `json:"enable_summarization"`
	EnableAnalytics     bool   `json:"enable_analytics"`
	EnableActionItems   bool   `json:"enable_action_items"`
	EnableFormatting    bool   `json:"enable_formatting"`
	DefaultModel        string `json:"default_model"`
	FormattingProfile   string `json:"formatting_profile"`
	CreateAt            int64  `json:"create_at"`
	UpdateAt            int64  `json:"update_at"`
}

func (p *AIPreferences) IsValid() *AppError {
	if !IsValidId(p.Id) {
		return NewAppError("AIPreferences.IsValid", "model.ai_preferences.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}

	if !IsValidId(p.UserId) {
		return NewAppError("AIPreferences.IsValid", "model.ai_preferences.is_valid.user_id.app_error", nil, "", http.StatusBadRequest)
	}

	if p.CreateAt == 0 {
		return NewAppError("AIPreferences.IsValid", "model.ai_preferences.is_valid.create_at.app_error", nil, "", http.StatusBadRequest)
	}

	if p.UpdateAt == 0 {
		return NewAppError("AIPreferences.IsValid", "model.ai_preferences.is_valid.update_at.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func (p *AIPreferences) PreSave() {
	if p.Id == "" {
		p.Id = NewId()
	}

	p.CreateAt = GetMillis()
	p.UpdateAt = p.CreateAt

	// Set defaults
	if p.DefaultModel == "" {
		p.DefaultModel = "gpt-3.5-turbo"
	}

	if p.FormattingProfile == "" {
		p.FormattingProfile = "professional"
	}
}

func (p *AIPreferences) PreUpdate() {
	p.UpdateAt = GetMillis()
}

