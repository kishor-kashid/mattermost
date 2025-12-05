// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/shared/mlog"
	"github.com/mattermost/mattermost/server/public/shared/request"
	"github.com/mattermost/mattermost/server/v8/channels/app/openai"
)

// DetectActionItems analyzes a post for potential action items using AI
func (a *App) DetectActionItems(c request.CTX, post *model.Post) (*ActionItemDetectionResult, error) {
	if !a.IsAIFeatureEnabled("action_items") {
		return &ActionItemDetectionResult{Detected: false}, nil
	}

	// Quick heuristic check before calling AI
	if !a.likelyContainsActionItem(post.Message) {
		c.Logger().Debug("Post unlikely to contain action items", 
			mlog.String("post_id", post.Id))
		return &ActionItemDetectionResult{Detected: false}, nil
	}

	// Get channel and user context
	channel, err := a.GetChannel(c, post.ChannelId)
	if err != nil {
		return nil, err
	}

	user, err := a.GetUser(post.UserId)
	if err != nil {
		return nil, err
	}

	// Call OpenAI for extraction
	items, appErr := a.extractActionItemsWithAI(c, post.Message, user.Username, channel.DisplayName)
	if appErr != nil {
		c.Logger().Error("Failed to extract action items with AI",
			mlog.String("post_id", post.Id),
			mlog.Err(appErr),
		)
		return &ActionItemDetectionResult{Detected: false}, appErr
	}

	if len(items) == 0 {
		return &ActionItemDetectionResult{Detected: false}, nil
	}

	// Enrich items with post context
	for _, item := range items {
		item.Id = model.NewId()
		item.PostId = post.Id
		item.ChannelId = post.ChannelId
		item.CreatedBy = post.UserId
		item.CreateAt = model.GetMillis()
		item.UpdateAt = item.CreateAt
		
		// Set default status and priority if not set
		if item.Status == "" {
			item.Status = "open"
		}
		if item.Priority == "" {
			item.Priority = "medium"
		}
	}

	return &ActionItemDetectionResult{
		Items:      items,
		Detected:   true,
		Confidence: 0.85, // Could be improved with AI confidence scores
	}, nil
}

// likelyContainsActionItem performs quick heuristic checks
func (a *App) likelyContainsActionItem(message string) bool {
	message = strings.ToLower(message)

	// Check for commitment indicators
	commitmentPhrases := []string{
		"i will", "i'll", "i can", "i'm going to",
		"will do", "will handle", "will take care",
		"let me", "i'll get", "i'll send", "i'll update",
		"by tomorrow", "by next week", "by end of",
		"todo:", "to do:", "action:", "task:",
		"need to", "have to", "should",
		"@", // Mentions often indicate assignments
	}

	for _, phrase := range commitmentPhrases {
		if strings.Contains(message, phrase) {
			return true
		}
	}

	// Check for deadline indicators
	deadlinePatterns := []*regexp.Regexp{
		regexp.MustCompile(`\b(today|tomorrow|tonight)\b`),
		regexp.MustCompile(`\b(monday|tuesday|wednesday|thursday|friday|saturday|sunday)\b`),
		regexp.MustCompile(`\b(jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)[a-z]*\s+\d{1,2}\b`),
		regexp.MustCompile(`\d{1,2}[/-]\d{1,2}[/-]\d{2,4}`),
		regexp.MustCompile(`\bby\s+(end\s+of|eod|eow|eom)\b`),
	}

	for _, pattern := range deadlinePatterns {
		if pattern.MatchString(message) {
			return true
		}
	}

	return false
}

// extractActionItemsWithAI uses OpenAI to extract structured action items
func (a *App) extractActionItemsWithAI(c request.CTX, message string, authorName string, channelName string) ([]*model.AIActionItem, *model.AppError) {
	aiService := a.GetAIService()
	if aiService == nil {
		return nil, model.NewAppError("extractActionItemsWithAI", "app.ai.service_not_available", nil, "", 500)
	}

	// Build the prompts using the prompt builder
	prompt := openai.GetActionItemExtractionPrompt()
	systemPrompt := prompt.System
	userPrompt := openai.BuildActionItemExtractionUserPrompt(message, authorName, channelName)

	// Call OpenAI
	aiModel := a.GetAIModel()
	response, err := aiService.client.SimpleCompletion(c.Context(), aiModel, systemPrompt, userPrompt)
	if err != nil {
		return nil, model.NewAppError("extractActionItemsWithAI", "app.ai.extraction_failed", nil, err.Error(), 500)
	}

	// Parse response
	items, err := a.parseActionItemsResponse(response)
	if err != nil {
		c.Logger().Warn("Failed to parse action items response",
			mlog.Err(err),
			mlog.String("response", response),
		)
		return nil, model.NewAppError("extractActionItemsWithAI", "app.ai.parse_failed", nil, err.Error(), 500)
	}

	return items, nil
}

// parseActionItemsResponse parses the AI response into action items
func (a *App) parseActionItemsResponse(response string) ([]*model.AIActionItem, error) {
	// Clean up response - remove markdown code blocks if present
	response = strings.TrimSpace(response)
	response = strings.TrimPrefix(response, "```json")
	response = strings.TrimPrefix(response, "```")
	response = strings.TrimSuffix(response, "```")
	response = strings.TrimSpace(response)

	// Parse JSON
	var rawItems []struct {
		Description string `json:"description"`
		Assignee    string `json:"assignee"`
		DueDate     string `json:"due_date"`
		Priority    string `json:"priority"`
	}

	if err := json.Unmarshal([]byte(response), &rawItems); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Convert to action items
	items := make([]*model.AIActionItem, 0, len(rawItems))
	
	for _, raw := range rawItems {
		if raw.Description == "" {
			continue
		}

		item := &model.AIActionItem{
			Description: raw.Description,
			Priority:    a.normalizePriority(raw.Priority),
			Status:      "open",
		}

		// Parse due date
		if raw.DueDate != "" {
			dueTime, err := time.Parse(time.RFC3339, raw.DueDate)
			if err != nil {
				// Try other common formats
				dueTime, err = time.Parse("2006-01-02", raw.DueDate)
			}
			if err == nil {
				item.DueDate = dueTime.UnixMilli()
			}
		}

		// Assignee will be resolved later based on username/mention
		// Store the raw assignee for now
		if raw.Assignee != "" {
			// This would need to be resolved to a user ID
			// For now, we'll handle this in the detection flow
		}

		items = append(items, item)
	}

	return items, nil
}

// normalizePriority normalizes priority strings
func (a *App) normalizePriority(priority string) string {
	priority = strings.ToLower(strings.TrimSpace(priority))
	
	validPriorities := map[string]bool{
		"low": true, "medium": true, "high": true, "urgent": true,
	}

	if validPriorities[priority] {
		return priority
	}

	// Default to medium
	return "medium"
}

// extractMentionedUsers extracts @mentioned usernames from a message
func (a *App) extractMentionedUsers(c request.CTX, message string) ([]string, error) {
	// Find all @mentions
	mentionRegex := regexp.MustCompile(`@([a-zA-Z0-9._-]+)`)
	matches := mentionRegex.FindAllStringSubmatch(message, -1)

	usernames := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) > 1 {
			usernames = append(usernames, match[1])
		}
	}

	return usernames, nil
}

// resolveAssignee resolves an assignee username to a user ID
func (a *App) resolveAssignee(c request.CTX, assignee string, channelID string, fallbackUserID string) string {
	// If assignee is empty, use the message author
	if assignee == "" {
		return fallbackUserID
	}

	// Try to find user by username
	user, err := a.GetUserByUsername(assignee)
	if err == nil && user != nil {
		return user.Id
	}

	// Fallback to message author
	return fallbackUserID
}

// AutoDetectAndCreateActionItems automatically detects and creates action items from a post
func (a *App) AutoDetectAndCreateActionItems(c request.CTX, post *model.Post) error {
	if !a.IsAIFeatureEnabled("action_items") {
		return nil
	}

	// Detect action items
	result, err := a.DetectActionItems(c, post)
	if err != nil {
		c.Logger().Error("Failed to detect action items",
			mlog.String("post_id", post.Id),
			mlog.Err(err),
		)
		return err
	}

	if !result.Detected || len(result.Items) == 0 {
		return nil
	}

	// Create each detected action item
	for _, item := range result.Items {
		// Resolve assignee - default to post author
		if item.AssigneeId == "" {
			item.AssigneeId = post.UserId
		}

		// Create the action item
		created, err := a.CreateActionItem(c, item)
		if err != nil {
			c.Logger().Error("Failed to create detected action item",
				mlog.String("post_id", post.Id),
				mlog.String("description", item.Description),
				mlog.Err(err),
			)
			continue
		}

		c.Logger().Info("Auto-created action item",
			mlog.String("action_item_id", created.Id),
			mlog.String("post_id", post.Id),
			mlog.String("assignee_id", created.AssigneeId),
			mlog.String("description", created.Description),
		)

		// Send a notification to the assignee
		a.notifyActionItemCreated(c, created, post)
	}

	return nil
}

// notifyActionItemCreated sends a DM notification about a new action item
func (a *App) notifyActionItemCreated(c request.CTX, item *model.AIActionItem, sourcePost *model.Post) {
	// Only notify if assignee is different from creator
	if item.AssigneeId == item.CreatedBy {
		return
	}

	// Get DM channel
	dmChannel, err := a.GetOrCreateDirectChannel(c, item.CreatedBy, item.AssigneeId)
	if err != nil {
		c.Logger().Error("Failed to get DM channel for action item notification",
			mlog.Err(err),
		)
		return
	}

	// Build notification message
	creator, _ := a.GetUser(item.CreatedBy)
	creatorName := "Someone"
	if creator != nil {
		creatorName = creator.Username
	}

	dueDateStr := "No due date"
	if item.DueDate > 0 {
		dueTime := time.UnixMilli(item.DueDate)
		dueDateStr = dueTime.Format("Mon Jan 2, 2006 at 3:04 PM")
	}

	message := fmt.Sprintf(`#### 📋 New Action Item Assigned

**%s** assigned you an action item:

> %s

**Priority:** %s  
**Due:** %s  
**Channel:** <#%s>

[View original message](/_redirect/pl/%s)`,
		creatorName,
		item.Description,
		strings.Title(item.Priority),
		dueDateStr,
		item.ChannelId,
		sourcePost.Id,
	)

	// Send DM
	notificationPost := &model.Post{
		UserId:    item.CreatedBy,
		ChannelId: dmChannel.Id,
		Message:   message,
	}

	_, err = a.CreatePost(c, notificationPost, dmChannel, model.CreatePostFlags{})
	if err != nil {
		c.Logger().Error("Failed to send action item notification",
			mlog.Err(err),
		)
	}
}

