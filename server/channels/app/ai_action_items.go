// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"fmt"
	"time"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/shared/mlog"
	"github.com/mattermost/mattermost/server/public/shared/request"
)

// CreateActionItem creates a new action item
func (a *App) CreateActionItem(c request.CTX, item *model.AIActionItem) (*model.AIActionItem, error) {
	if !a.IsAIFeatureEnabled("action_items") {
		return nil, fmt.Errorf("action items feature is not enabled")
	}

	// Validate action item
	if err := a.ValidateActionItem(item); err != nil {
		return nil, err
	}

	// Check permissions - user must have access to the channel
	if item.ChannelId != "" {
		if !a.HasPermissionToChannel(c, item.CreatedBy, item.ChannelId, model.PermissionReadChannel) {
			return nil, model.NewAppError("CreateActionItem", "api.action_item.create.permission_denied", nil, "", 403)
		}
	}

	// Set defaults
	item.Id = model.NewId()
	item.CreateAt = model.GetMillis()
	item.UpdateAt = item.CreateAt

	if item.Status == "" {
		item.Status = "open"
	}

	if item.Priority == "" {
		item.Priority = "medium"
	}

	// Save to database
	created, err := a.Srv().Store().AIActionItem().Save(item)
	if err != nil {
		return nil, err
	}

	c.Logger().Debug("Created action item",
		mlog.String("action_item_id", created.Id),
		mlog.String("assignee_id", created.AssigneeId),
		mlog.String("channel_id", created.ChannelId),
	)

	return created, nil
}

// GetActionItem retrieves a single action item by ID
func (a *App) GetActionItem(c request.CTX, actionItemID string, userID string) (*model.AIActionItem, error) {
	if !a.IsAIFeatureEnabled("action_items") {
		return nil, fmt.Errorf("action items feature is not enabled")
	}

	item, err := a.Srv().Store().AIActionItem().Get(actionItemID)
	if err != nil {
		return nil, err
	}

	// Check permissions - user must be assignee, creator, or have access to channel
	if !a.CanUserAccessActionItem(c, userID, item) {
		return nil, model.NewAppError("GetActionItem", "api.action_item.get.permission_denied", nil, "", 403)
	}

	return item, nil
}

// GetActionItemsForUser retrieves action items for a specific user
func (a *App) GetActionItemsForUser(c request.CTX, userID string, filters *ActionItemFilters) ([]*model.AIActionItem, error) {
	if !a.IsAIFeatureEnabled("action_items") {
		return nil, fmt.Errorf("action items feature is not enabled")
	}

	// Set default pagination
	if filters.PerPage == 0 {
		filters.PerPage = 60
	}
	if filters.PerPage > 200 {
		filters.PerPage = 200
	}

	items, err := a.Srv().Store().AIActionItem().GetByUser(userID, filters.IncludeCompleted, filters.Page, filters.PerPage)
	if err != nil {
		return nil, err
	}

	// Apply additional filters
	filtered := a.filterActionItems(items, filters)

	return filtered, nil
}

// GetActionItemsForChannel retrieves action items for a specific channel
func (a *App) GetActionItemsForChannel(c request.CTX, channelID string, userID string, filters *ActionItemFilters) ([]*model.AIActionItem, error) {
	if !a.IsAIFeatureEnabled("action_items") {
		return nil, fmt.Errorf("action items feature is not enabled")
	}

	// Check channel permissions
	if !a.HasPermissionToChannel(c, userID, channelID, model.PermissionReadChannel) {
		return nil, model.NewAppError("GetActionItemsForChannel", "api.action_item.get_channel.permission_denied", nil, "", 403)
	}

	// Set default pagination
	if filters.PerPage == 0 {
		filters.PerPage = 60
	}

	items, err := a.Srv().Store().AIActionItem().GetByChannel(channelID, filters.IncludeCompleted, filters.Page, filters.PerPage)
	if err != nil {
		return nil, err
	}

	// Apply additional filters
	filtered := a.filterActionItems(items, filters)

	return filtered, nil
}

// UpdateActionItem updates an existing action item
func (a *App) UpdateActionItem(c request.CTX, actionItemID string, userID string, update *ActionItemUpdateRequest) (*model.AIActionItem, error) {
	if !a.IsAIFeatureEnabled("action_items") {
		return nil, fmt.Errorf("action items feature is not enabled")
	}

	// Get existing item
	item, err := a.Srv().Store().AIActionItem().Get(actionItemID)
	if err != nil {
		return nil, err
	}

	// Check permissions
	if !a.CanUserModifyActionItem(c, userID, item) {
		return nil, model.NewAppError("UpdateActionItem", "api.action_item.update.permission_denied", nil, "", 403)
	}

	// Apply updates
	if update.Description != nil {
		item.Description = *update.Description
	}
	if update.AssigneeID != nil {
		item.AssigneeId = *update.AssigneeID
	}
	if update.DueDate != nil {
		item.DueDate = update.DueDate.UnixMilli()
	}
	if update.Priority != nil {
		item.Priority = *update.Priority
	}
	if update.Status != nil {
		item.Status = *update.Status
		
		// If marking as completed, set CompletedAt
		if *update.Status == "completed" && item.CompletedAt == 0 {
			item.CompletedAt = model.GetMillis()
		}
	}
	if update.CompletedAt != nil {
		item.CompletedAt = update.CompletedAt.UnixMilli()
	}

	item.UpdateAt = model.GetMillis()

	// Validate
	if err := a.ValidateActionItem(item); err != nil {
		return nil, err
	}

	// Save
	updated, err := a.Srv().Store().AIActionItem().Update(item)
	if err != nil {
		return nil, err
	}

	c.Logger().Debug("Updated action item",
		mlog.String("action_item_id", updated.Id),
		mlog.String("status", updated.Status),
	)

	return updated, nil
}

// CompleteActionItem marks an action item as completed
func (a *App) CompleteActionItem(c request.CTX, actionItemID string, userID string) (*model.AIActionItem, error) {
	status := "completed"
	now := time.Now()
	
	return a.UpdateActionItem(c, actionItemID, userID, &ActionItemUpdateRequest{
		Status:      &status,
		CompletedAt: &now,
	})
}

// DeleteActionItem deletes an action item
func (a *App) DeleteActionItem(c request.CTX, actionItemID string, userID string) error {
	if !a.IsAIFeatureEnabled("action_items") {
		return fmt.Errorf("action items feature is not enabled")
	}

	// Get existing item
	item, err := a.Srv().Store().AIActionItem().Get(actionItemID)
	if err != nil {
		return err
	}

	// Check permissions
	if !a.CanUserModifyActionItem(c, userID, item) {
		return model.NewAppError("DeleteActionItem", "api.action_item.delete.permission_denied", nil, "", 403)
	}

	// Delete
	deleteAt := time.Now().Unix() * 1000 // Convert to milliseconds
	err = a.Srv().Store().AIActionItem().Delete(actionItemID, deleteAt)
	if err != nil {
		return err
	}

	c.Logger().Debug("Deleted action item", mlog.String("action_item_id", actionItemID))

	return nil
}

// GetOverdueActionItems retrieves all overdue action items
func (a *App) GetOverdueActionItems(c request.CTX) ([]*model.AIActionItem, error) {
	if !a.IsAIFeatureEnabled("action_items") {
		return nil, fmt.Errorf("action items feature is not enabled")
	}

	now := time.Now()
	items, err := a.Srv().Store().AIActionItem().GetOverdue(now.UnixMilli())
	if err != nil {
		return nil, err
	}

	return items, nil
}

// GetDueSoonActionItems retrieves action items due within the next 24 hours
func (a *App) GetDueSoonActionItems(c request.CTX, hours int) ([]*model.AIActionItem, error) {
	if !a.IsAIFeatureEnabled("action_items") {
		return nil, fmt.Errorf("action items feature is not enabled")
	}

	now := time.Now()
	future := now.Add(time.Duration(hours) * time.Hour)
	
	items, err := a.Srv().Store().AIActionItem().GetDueSoon(now.UnixMilli(), future.UnixMilli())
	if err != nil {
		return nil, err
	}

	return items, nil
}

// ValidateActionItem validates an action item
func (a *App) ValidateActionItem(item *model.AIActionItem) error {
	if item.Description == "" {
		return fmt.Errorf("description is required")
	}

	if len(item.Description) > 500 {
		return fmt.Errorf("description too long (max 500 characters)")
	}

	if item.AssigneeId == "" {
		return fmt.Errorf("assignee is required")
	}

	// Validate priority
	validPriorities := map[string]bool{"low": true, "medium": true, "high": true, "urgent": true}
	if item.Priority != "" && !validPriorities[item.Priority] {
		return fmt.Errorf("invalid priority: %s", item.Priority)
	}

	// Validate status
	validStatuses := map[string]bool{"open": true, "in_progress": true, "completed": true, "dismissed": true}
	if item.Status != "" && !validStatuses[item.Status] {
		return fmt.Errorf("invalid status: %s", item.Status)
	}

	return nil
}

// CanUserAccessActionItem checks if a user can access an action item
func (a *App) CanUserAccessActionItem(c request.CTX, userID string, item *model.AIActionItem) bool {
	// User is the assignee
	if item.AssigneeId == userID {
		return true
	}

	// User created the item
	if item.CreatedBy == userID {
		return true
	}

	// User has access to the channel
	if item.ChannelId != "" {
		return a.HasPermissionToChannel(c, userID, item.ChannelId, model.PermissionReadChannel)
	}

	return false
}

// CanUserModifyActionItem checks if a user can modify an action item
func (a *App) CanUserModifyActionItem(c request.CTX, userID string, item *model.AIActionItem) bool {
	// User is the assignee (can update their own items)
	if item.AssigneeId == userID {
		return true
	}

	// User created the item
	if item.CreatedBy == userID {
		return true
	}

	// User is a channel admin
	if item.ChannelId != "" {
		return a.HasPermissionToChannel(c, userID, item.ChannelId, model.PermissionManageChannelRoles)
	}

	return false
}

// filterActionItems applies additional filters to action items
func (a *App) filterActionItems(items []*model.AIActionItem, filters *ActionItemFilters) []*model.AIActionItem {
	if filters == nil {
		return items
	}

	filtered := make([]*model.AIActionItem, 0, len(items))

	for _, item := range items {
		// Filter by status
		if filters.Status != "" && item.Status != filters.Status {
			continue
		}

		// Filter by priority
		if filters.Priority != "" && item.Priority != filters.Priority {
			continue
		}

		// Filter by due date range
		if filters.DueBefore != nil && item.DueDate > filters.DueBefore.UnixMilli() {
			continue
		}
		if filters.DueAfter != nil && item.DueDate < filters.DueAfter.UnixMilli() {
			continue
		}

		// Filter by assigned by
		if filters.AssignedBy != "" && item.CreatedBy != filters.AssignedBy {
			continue
		}

		filtered = append(filtered, item)
	}

	return filtered
}

// GetActionItemStats retrieves statistics about action items for a user
func (a *App) GetActionItemStats(c request.CTX, userID string) (*ActionItemStats, error) {
	if !a.IsAIFeatureEnabled("action_items") {
		return nil, fmt.Errorf("action items feature is not enabled")
	}

	items, err := a.Srv().Store().AIActionItem().GetByUser(userID, true, 0, 1000)
	if err != nil {
		return nil, err
	}

	stats := &ActionItemStats{
		Total:      0,
		Overdue:    0,
		DueToday:   0,
		DueSoon:    0,
		NoDueDate:  0,
		Completed:  0,
		ByPriority: make(map[string]int),
		ByStatus:   make(map[string]int),
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	weekFromNow := now.Add(7 * 24 * time.Hour)

	for _, item := range items {
		stats.Total++
		stats.ByPriority[item.Priority]++
		stats.ByStatus[item.Status]++

		if item.Status == "completed" {
			stats.Completed++
		}

		if item.DueDate > 0 {
			dueTime := time.UnixMilli(item.DueDate)
			
			if item.Status != "completed" {
				if dueTime.Before(now) {
					stats.Overdue++
				} else if dueTime.Before(today) {
					stats.DueToday++
				} else if dueTime.Before(weekFromNow) {
					stats.DueSoon++
				}
			}
		} else {
			stats.NoDueDate++
		}
	}

	return stats, nil
}

