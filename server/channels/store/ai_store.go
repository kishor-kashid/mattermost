// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package store

import (
	"github.com/mattermost/mattermost/server/public/model"
)

// AIActionItemStore represents a store for managing AI-detected action items
type AIActionItemStore interface {
	// Save creates or updates an action item
	Save(actionItem *model.AIActionItem) (*model.AIActionItem, error)
	
	// Get retrieves an action item by ID
	Get(id string) (*model.AIActionItem, error)
	
	// GetByChannel retrieves all action items for a channel
	GetByChannel(channelId string, includeCompleted bool, offset, limit int) ([]*model.AIActionItem, error)
	
	// GetByUser retrieves all action items created by or assigned to a user
	GetByUser(userId string, includeCompleted bool, offset, limit int) ([]*model.AIActionItem, error)
	
	// GetByAssignee retrieves all action items for a specific assignee
	GetByAssignee(assigneeId string, offset, limit int) ([]*model.AIActionItem, error)
	
	// GetOverdue retrieves action items that are past their due date
	GetOverdue(currentTime int64) ([]*model.AIActionItem, error)
	
	// GetDueSoon retrieves action items due within a time range
	GetDueSoon(startTime, endTime int64) ([]*model.AIActionItem, error)
	
	// Update updates an existing action item
	Update(actionItem *model.AIActionItem) (*model.AIActionItem, error)
	
	// Delete soft-deletes an action item
	Delete(id string, deleteAt int64) error
	
	// PermanentDelete permanently removes an action item
	PermanentDelete(id string) error
}

// AISummaryStore represents a store for managing AI-generated summaries
type AISummaryStore interface {
	// Save creates a new summary
	Save(summary *model.AISummary) (*model.AISummary, error)
	
	// Get retrieves a summary by ID
	Get(id string) (*model.AISummary, error)
	
	// GetByChannel retrieves summaries for a channel
	GetByChannel(channelId string, offset, limit int) ([]*model.AISummary, error)
	
	// GetCachedSummary retrieves a cached summary if not expired
	GetCachedSummary(channelId, summaryType string, startTime, endTime int64) (*model.AISummary, error)
	
	// GetByCacheKey retrieves a summary by its cache key if not expired
	GetByCacheKey(cacheKey string) (*model.AISummary, error)
	
	// DeleteExpired removes expired summaries
	DeleteExpired(currentTime int64) (int64, error)
	
	// Delete removes a specific summary
	Delete(id string) error
}

// AIAnalyticsStore represents a store for managing AI analytics data
type AIAnalyticsStore interface {
	// Save creates or updates analytics data
	Save(analytics *model.AIAnalytics) (*model.AIAnalytics, error)
	
	// Get retrieves analytics by ID
	Get(id string) (*model.AIAnalytics, error)
	
	// GetByChannel retrieves analytics for a channel
	GetByChannel(channelId string, startDate, endDate string) ([]*model.AIAnalytics, error)
	
	// GetByChannelAndDate retrieves analytics for a specific channel and date
	GetByChannelAndDate(channelId, date string) (*model.AIAnalytics, error)
	
	// Update updates existing analytics
	Update(analytics *model.AIAnalytics) (*model.AIAnalytics, error)
	
	// DeleteOlderThan removes analytics older than specified date
	DeleteOlderThan(date string) (int64, error)
}

// AIPreferencesStore represents a store for managing user AI preferences
type AIPreferencesStore interface {
	// Save creates or updates user preferences
	Save(preferences *model.AIPreferences) (*model.AIPreferences, error)
	
	// Get retrieves preferences by ID
	Get(id string) (*model.AIPreferences, error)
	
	// GetByUser retrieves preferences for a specific user
	GetByUser(userId string) (*model.AIPreferences, error)
	
	// Update updates existing preferences
	Update(preferences *model.AIPreferences) (*model.AIPreferences, error)
	
	// Delete removes user preferences
	Delete(userId string) error
	
	// GetFormatterPreferences retrieves formatting preferences for a user
	GetFormatterPreferences(userId string) (defaultProfile string, autoSuggest bool, err error)
	
	// SetFormatterPreferences updates formatting preferences for a user
	SetFormatterPreferences(userId string, defaultProfile string, autoSuggest bool) error
}

