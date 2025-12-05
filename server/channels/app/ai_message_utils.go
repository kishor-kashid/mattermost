// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/shared/mlog"
)

// FormatMessagesForPrompt formats a list of messages for AI prompt input
func FormatMessagesForPrompt(posts []*model.Post, users map[string]*model.User) string {
	if len(posts) == 0 {
		return ""
	}

	// Sort posts by creation time
	sortedPosts := make([]*model.Post, len(posts))
	copy(sortedPosts, posts)
	sort.Slice(sortedPosts, func(i, j int) bool {
		return sortedPosts[i].CreateAt < sortedPosts[j].CreateAt
	})

	var builder strings.Builder

	for _, post := range sortedPosts {
		username := "Unknown User"
		if user, ok := users[post.UserId]; ok {
			username = user.Username
		}

		// Format: [Username] Message
		builder.WriteString(fmt.Sprintf("[%s] %s\n", username, post.Message))
	}

	return builder.String()
}

// ExtractParticipants extracts unique participant usernames from a list of posts
func ExtractParticipants(posts []*model.Post, users map[string]*model.User) []string {
	participantMap := make(map[string]bool)
	participants := []string{}

	for _, post := range posts {
		if _, exists := participantMap[post.UserId]; !exists {
			participantMap[post.UserId] = true

			username := "Unknown User"
			if user, ok := users[post.UserId]; ok {
				username = user.Username
			}

			participants = append(participants, username)
		}
	}

	sort.Strings(participants)
	return participants
}

// GetUsersForPosts fetches user information for all posts
func (a *App) GetUsersForPosts(posts []*model.Post) (map[string]*model.User, error) {
	userIds := make(map[string]bool)
	for _, post := range posts {
		userIds[post.UserId] = true
	}

	users := make(map[string]*model.User)
	for userId := range userIds {
		user, err := a.GetUser(userId)
		if err != nil {
			// Log error but continue with other users
			a.Log().Warn("Failed to fetch user for AI processing", mlog.String("user_id", userId), mlog.Err(err))
			continue
		}
		users[userId] = user
	}

	return users, nil
}

// TruncateMessageList truncates a list of messages to a maximum count
func TruncateMessageList(posts []*model.Post, maxCount int) []*model.Post {
	if len(posts) <= maxCount {
		return posts
	}

	// Return the most recent messages
	sorted := make([]*model.Post, len(posts))
	copy(sorted, posts)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].CreateAt > sorted[j].CreateAt
	})

	return sorted[:maxCount]
}

// FormatMessageCount returns a human-readable message count string
func FormatMessageCount(count int) string {
	if count == 1 {
		return "1 message"
	}
	return fmt.Sprintf("%d messages", count)
}

// FormatParticipantList returns a human-readable participant list
func FormatParticipantList(participants []string) string {
	switch len(participants) {
	case 0:
		return "No participants"
	case 1:
		return participants[0]
	case 2:
		return fmt.Sprintf("%s and %s", participants[0], participants[1])
	default:
		return fmt.Sprintf("%s, and %d others", strings.Join(participants[:2], ", "), len(participants)-2)
	}
}

// SanitizeMessageForAI removes potentially sensitive or problematic content
func SanitizeMessageForAI(message string) string {
	// For now, just trim whitespace
	// In the future, we might want to:
	// - Remove API keys or tokens
	// - Redact sensitive patterns
	// - Remove excessive whitespace
	return strings.TrimSpace(message)
}

// BuildMessageContext builds context information about a set of messages
func BuildMessageContext(posts []*model.Post, users map[string]*model.User) map[string]interface{} {
	participants := ExtractParticipants(posts, users)

	// Calculate time range
	var minTime, maxTime int64
	if len(posts) > 0 {
		minTime = posts[0].CreateAt
		maxTime = posts[0].CreateAt

		for _, post := range posts {
			if post.CreateAt < minTime {
				minTime = post.CreateAt
			}
			if post.CreateAt > maxTime {
				maxTime = post.CreateAt
			}
		}
	}

	return map[string]interface{}{
		"message_count":    len(posts),
		"participant_count": len(participants),
		"participants":     participants,
		"time_range_start": minTime,
		"time_range_end":   maxTime,
		"duration_ms":      maxTime - minTime,
	}
}

