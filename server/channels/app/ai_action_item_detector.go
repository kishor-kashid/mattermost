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
	likelyContains := a.likelyContainsActionItem(post.Message)
	c.Logger().Info("Heuristic check result", 
		mlog.Bool("likely_contains", likelyContains),
		mlog.String("message", post.Message))
	
	if !likelyContains {
		c.Logger().Info("Post unlikely to contain action items - skipping AI call", 
			mlog.String("post_id", post.Id))
		return &ActionItemDetectionResult{Detected: false}, nil
	}
	
	c.Logger().Info("Post passed heuristic check - calling OpenAI", 
		mlog.String("post_id", post.Id))

	// Get channel and user context
	channel, err := a.GetChannel(c, post.ChannelId)
	if err != nil {
		return nil, err
	}

	user, err := a.GetUser(post.UserId)
	if err != nil {
		return nil, err
	}

	// Build context - include parent post if this is a reply
	messageContext := post.Message
	if post.RootId != "" {
		parentPost, err := a.GetSinglePost(c, post.RootId, false)
		if err == nil {
			parentUser, err := a.GetUser(parentPost.UserId)
			if err == nil {
				messageContext = fmt.Sprintf("Previous message from %s: \"%s\"\n\nCurrent message: %s", 
					parentUser.Username, parentPost.Message, post.Message)
			}
		}
	}

	// Call OpenAI for extraction
	items, appErr := a.extractActionItemsWithAI(c, messageContext, user.Username, channel.DisplayName)
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

	// Parse JSON - match the structure from the prompt template
	var aiResponse struct {
		HasActionItems bool `json:"has_action_items"`
		ActionItems    []struct {
			Description string `json:"description"`
			Assignee    string `json:"assignee"`
			Deadline    string `json:"deadline"`  // OpenAI uses "deadline", not "due_date"
			Priority    string `json:"priority"`
		} `json:"action_items"`
	}

	if err := json.Unmarshal([]byte(response), &aiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Check if AI detected any action items
	if !aiResponse.HasActionItems || len(aiResponse.ActionItems) == 0 {
		return []*model.AIActionItem{}, nil
	}

	// Convert to action items
	items := make([]*model.AIActionItem, 0, len(aiResponse.ActionItems))
	
	for _, raw := range aiResponse.ActionItems {
		if raw.Description == "" {
			continue
		}

		item := &model.AIActionItem{
			Description: raw.Description,
			Priority:    a.normalizePriority(raw.Priority),
			Status:      "open",
		}

		// Parse deadline - try to convert natural language to timestamp
		if raw.Deadline != "" && raw.Deadline != "unspecified" {
			dueTime := a.parseDeadline(raw.Deadline)
			if !dueTime.IsZero() {
				item.DueDate = dueTime.UnixMilli()
				// Log the parsed deadline for debugging
				fmt.Printf("DEBUG: Parsed deadline '%s' to %v (timestamp: %d)\n", raw.Deadline, dueTime, item.DueDate)
			} else {
				fmt.Printf("DEBUG: Failed to parse deadline '%s'\n", raw.Deadline)
			}
		}

		// Assignee will be resolved later based on username/mention
		// Store the raw assignee for now
		if raw.Assignee != "" && raw.Assignee != "unspecified" {
			// This would need to be resolved to a user ID
			// For now, we'll handle this in the detection flow
		}

		items = append(items, item)
	}

	return items, nil
}

// parseDeadline tries to parse natural language deadlines into timestamps
func (a *App) parseDeadline(deadline string) time.Time {
	deadline = strings.ToLower(strings.TrimSpace(deadline))
	now := time.Now()

	// Try ISO format first
	if t, err := time.Parse(time.RFC3339, deadline); err == nil {
		return t
	}
	if t, err := time.Parse("2006-01-02", deadline); err == nil {
		return t
	}

	// Natural language parsing
	switch {
	case strings.Contains(deadline, "today"):
		return now
	case strings.Contains(deadline, "tomorrow"):
		return now.AddDate(0, 0, 1)
	case strings.Contains(deadline, "end of this week"), strings.Contains(deadline, "end of week"):
		// Friday of current week
		daysUntilFriday := (5 - int(now.Weekday()) + 7) % 7
		if daysUntilFriday == 0 {
			daysUntilFriday = 7
		}
		return now.AddDate(0, 0, daysUntilFriday)
	case strings.Contains(deadline, "next week"):
		return now.AddDate(0, 0, 7)
	case strings.Contains(deadline, "end of month"):
		return time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 0, now.Location())
	}

	// If we can't parse it, return zero time
	return time.Time{}
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
	c.Logger().Info("AutoDetectAndCreateActionItems triggered", 
		mlog.String("post_id", post.Id),
		mlog.Int("message_length", len(post.Message)))
	
	if !a.IsAIFeatureEnabled("action_items") {
		c.Logger().Info("Action items feature check failed")
		return nil
	}
	
	c.Logger().Info("Action items feature enabled, proceeding with detection")

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

		dueStr := "no due date"
		if created.DueDate > 0 {
			dueStr = time.UnixMilli(created.DueDate).Format("2006-01-02 15:04")
		}
		c.Logger().Info("Auto-created action item",
			mlog.String("action_item_id", created.Id),
			mlog.String("post_id", post.Id),
			mlog.String("assignee_id", created.AssigneeId),
			mlog.String("description", created.Description),
			mlog.String("due_date", dueStr),
			mlog.String("priority", created.Priority),
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

