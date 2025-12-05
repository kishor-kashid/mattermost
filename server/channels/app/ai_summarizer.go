// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/shared/mlog"
	"github.com/mattermost/mattermost/server/public/shared/request"
	"github.com/mattermost/mattermost/server/v8/channels/app/openai"
)

// SummarizeThread generates an AI summary of a thread
func (a *App) SummarizeThread(c request.CTX, req *SummarizationRequest) (*SummarizationResponse, *model.AppError) {
	startTime := time.Now()

	// Check if feature is enabled
	if !a.IsAIFeatureEnabled("summarization") {
		return nil, model.NewAppError("SummarizeThread", "app.ai.summarization_disabled", nil, "", 403)
	}

	// Validate request
	if req.ChannelId == "" {
		return nil, model.NewAppError("SummarizeThread", "app.ai.invalid_channel_id", nil, "", 400)
	}
	if req.PostId == "" {
		return nil, model.NewAppError("SummarizeThread", "app.ai.invalid_post_id", nil, "", 400)
	}

	// Check channel membership permission
	if !a.HasPermissionToChannel(c, req.UserId, req.ChannelId, model.PermissionReadChannel) {
		return nil, model.NewAppError("SummarizeThread", "app.ai.no_channel_permission", nil, "", 403)
	}

	// Check cache first
	if req.UseCache {
		cacheKey := a.generateSummaryCacheKey(req.PostId, model.AISummaryTypeThread, req.SummaryLevel)
		cachedSummary, err := a.Srv().Store().AISummary().GetByCacheKey(cacheKey)
		if err == nil && cachedSummary != nil {
			// Check if cache is still valid (not expired)
			if cachedSummary.ExpiresAt > model.GetMillis() {
				a.Log().Debug("Returning cached thread summary", mlog.String("post_id", req.PostId))
				return &SummarizationResponse{
					Summary:      cachedSummary,
					FromCache:    true,
					ProcessingMs: time.Since(startTime).Milliseconds(),
				}, nil
			}
		}
	}

	// Fetch thread posts
	posts, err := a.getThreadPosts(c, req.PostId, req.MaxMessages)
	if err != nil {
		return nil, err
	}

	if len(posts) == 0 {
		return nil, model.NewAppError("SummarizeThread", "app.ai.no_messages_found", nil, "", 404)
	}

	// Get root post to understand thread context
	rootPost := posts[0]
	if rootPost.RootId != "" {
		// This isn't actually the root, fetch the real root
		actualRoot, err := a.GetSinglePost(c, rootPost.RootId, false)
		if err == nil {
			rootPost = actualRoot
			// Add root to beginning if not already there
			if posts[0].Id != rootPost.Id {
				posts = append([]*model.Post{rootPost}, posts...)
			}
		}
	}

	// Format messages for LLM
	messageContexts, participants, err := a.formatMessagesForLLM(c, posts)
	if err != nil {
		return nil, err
	}

	// Build the prompt
	level := req.SummaryLevel
	if level == "" {
		level = string(openai.SummarizationStandard)
	}
	
	messages := a.buildMessageText(messageContexts)
	participantList := strings.Join(participants, ", ")
	
	userPrompt := openai.BuildSummarizationUserPrompt(messages, participantList, len(posts), true)
	systemPrompt, _ := openai.GetSummarizationPrompt(openai.SummarizationLevel(level)).Substitute(nil)

	// Get AI service
	aiService := a.GetAIService()
	if aiService == nil {
		return nil, model.NewAppError("SummarizeThread", "app.ai.service_not_available", nil, "", 500)
	}

	// Generate summary via OpenAI
	summaryText, openaiErr := aiService.client.SimpleCompletion(c.Context(), a.GetAIModel(), systemPrompt, userPrompt)
	if openaiErr != nil {
		a.Log().Error("Failed to generate thread summary", mlog.Err(openaiErr))
		return nil, model.NewAppError("SummarizeThread", "app.ai.openai_error", nil, openaiErr.Error(), 500)
	}

	// Get channel info
	channel, err := a.GetChannel(c, req.ChannelId)
	if err != nil {
		channel = &model.Channel{DisplayName: "Unknown Channel"}
	}

	// Create summary record
	summary := &model.AISummary{
		ChannelId:    req.ChannelId,
		PostId:       req.PostId,
		SummaryType:  model.AISummaryTypeThread,
		Summary:      summaryText,
		MessageCount: len(posts),
		StartTime:    posts[len(posts)-1].CreateAt,
		EndTime:      posts[0].CreateAt,
		UserId:       req.UserId,
		Participants: participantList,
		CacheKey:     a.generateSummaryCacheKey(req.PostId, model.AISummaryTypeThread, level),
		ExpiresAt:    model.GetMillis() + (24 * 60 * 60 * 1000), // 24 hours
		ChannelName:  channel.DisplayName,
	}

	summary.PreSave()

	// Save to cache
	savedSummary, saveErr := a.Srv().Store().AISummary().Save(summary)
	if saveErr != nil {
		a.Log().Warn("Failed to cache thread summary", mlog.Err(saveErr))
		// Don't fail the request, just log the warning
		savedSummary = summary
	}

	return &SummarizationResponse{
		Summary:      savedSummary,
		FromCache:    false,
		ProcessingMs: time.Since(startTime).Milliseconds(),
	}, nil
}

// SummarizeChannel generates an AI summary of channel messages
func (a *App) SummarizeChannel(c request.CTX, req *SummarizationRequest) (*SummarizationResponse, *model.AppError) {
	startTime := time.Now()

	// Check if feature is enabled
	if !a.IsAIFeatureEnabled("summarization") {
		return nil, model.NewAppError("SummarizeChannel", "app.ai.summarization_disabled", nil, "", 403)
	}

	// Validate request
	if req.ChannelId == "" {
		return nil, model.NewAppError("SummarizeChannel", "app.ai.invalid_channel_id", nil, "", 400)
	}

	// Check channel membership permission
	if !a.HasPermissionToChannel(c, req.UserId, req.ChannelId, model.PermissionReadChannel) {
		return nil, model.NewAppError("SummarizeChannel", "app.ai.no_channel_permission", nil, "", 403)
	}

	// Set defaults for time range
	if req.EndTime == 0 {
		req.EndTime = model.GetMillis()
	}
	if req.StartTime == 0 {
		// Default to last 24 hours
		req.StartTime = req.EndTime - (24 * 60 * 60 * 1000)
	}

	// Check cache first
	if req.UseCache {
		cacheKey := a.generateChannelSummaryCacheKey(req.ChannelId, req.StartTime, req.EndTime, req.SummaryLevel)
		cachedSummary, err := a.Srv().Store().AISummary().GetByCacheKey(cacheKey)
		if err == nil && cachedSummary != nil {
			// Check if cache is still valid (not expired)
			if cachedSummary.ExpiresAt > model.GetMillis() {
				a.Log().Debug("Returning cached channel summary", mlog.String("channel_id", req.ChannelId))
				return &SummarizationResponse{
					Summary:      cachedSummary,
					FromCache:    true,
					ProcessingMs: time.Since(startTime).Milliseconds(),
				}, nil
			}
		}
	}

	// Fetch channel posts with pagination
	maxMessages := req.MaxMessages
	if maxMessages == 0 {
		maxMessages = a.GetAIMaxMessageLimit()
	}

	posts, err := a.getChannelPostsInRange(c, req.ChannelId, req.StartTime, req.EndTime, maxMessages)
	if err != nil {
		return nil, err
	}

	if len(posts) == 0 {
		return nil, model.NewAppError("SummarizeChannel", "app.ai.no_messages_found", nil, "", 404)
	}

	// Format messages for LLM
	messageContexts, participants, err := a.formatMessagesForLLM(c, posts)
	if err != nil {
		return nil, err
	}

	// Build the prompt
	level := req.SummaryLevel
	if level == "" {
		level = string(openai.SummarizationStandard)
	}
	
	messages := a.buildMessageText(messageContexts)
	participantList := strings.Join(participants, ", ")
	
	userPrompt := openai.BuildSummarizationUserPrompt(messages, participantList, len(posts), false)
	systemPrompt, _ := openai.GetSummarizationPrompt(openai.SummarizationLevel(level)).Substitute(nil)

	// Get AI service
	aiService := a.GetAIService()
	if aiService == nil {
		return nil, model.NewAppError("SummarizeChannel", "app.ai.service_not_available", nil, "", 500)
	}

	// Generate summary via OpenAI
	summaryText, openaiErr := aiService.client.SimpleCompletion(c.Context(), a.GetAIModel(), systemPrompt, userPrompt)
	if openaiErr != nil {
		a.Log().Error("Failed to generate channel summary", mlog.Err(openaiErr))
		return nil, model.NewAppError("SummarizeChannel", "app.ai.openai_error", nil, openaiErr.Error(), 500)
	}

	// Get channel info
	channel, err := a.GetChannel(c, req.ChannelId)
	if err != nil {
		channel = &model.Channel{DisplayName: "Unknown Channel"}
	}

	// Create summary record
	cacheKey := a.generateChannelSummaryCacheKey(req.ChannelId, req.StartTime, req.EndTime, level)
	summary := &model.AISummary{
		ChannelId:    req.ChannelId,
		SummaryType:  model.AISummaryTypeChannel,
		Summary:      summaryText,
		MessageCount: len(posts),
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		UserId:       req.UserId,
		Participants: participantList,
		CacheKey:     cacheKey,
		ExpiresAt:    model.GetMillis() + (24 * 60 * 60 * 1000), // 24 hours
		ChannelName:  channel.DisplayName,
	}

	summary.PreSave()

	// Save to cache
	savedSummary, saveErr := a.Srv().Store().AISummary().Save(summary)
	if saveErr != nil {
		a.Log().Warn("Failed to cache channel summary", mlog.Err(saveErr))
		// Don't fail the request, just log the warning
		savedSummary = summary
	}

	return &SummarizationResponse{
		Summary:      savedSummary,
		FromCache:    false,
		ProcessingMs: time.Since(startTime).Milliseconds(),
	}, nil
}

// getThreadPosts fetches all posts in a thread
func (a *App) getThreadPosts(c request.CTX, postId string, maxMessages int) ([]*model.Post, *model.AppError) {
	// Get the post to determine if it's a root or reply
	post, err := a.GetSinglePost(c, postId, false)
	if err != nil {
		return nil, err
	}

	// Determine the root ID
	rootId := post.Id
	if post.RootId != "" {
		rootId = post.RootId
	}

	// Fetch the thread
	opts := model.GetPostsOptions{
		SkipFetchThreads: false,
	}
	postList, err := a.GetPostThread(c, rootId, opts, "")
	if err != nil {
		return nil, err
	}

	// Convert to slice and sort by creation time (oldest first for context)
	posts := make([]*model.Post, 0, len(postList.Posts))
	for _, p := range postList.Posts {
		posts = append(posts, p)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreateAt < posts[j].CreateAt
	})

	// Limit to max messages
	if maxMessages > 0 && len(posts) > maxMessages {
		posts = posts[:maxMessages]
	}

	return posts, nil
}

// getChannelPostsInRange fetches posts from a channel within a time range
func (a *App) getChannelPostsInRange(c request.CTX, channelId string, startTime, endTime int64, maxMessages int) ([]*model.Post, *model.AppError) {
	// Fetch posts in reverse chronological order
	posts := make([]*model.Post, 0)
	page := 0
	perPage := 100
	
	for len(posts) < maxMessages {
		postList, err := a.GetPosts(c, channelId, page, perPage)
		if err != nil {
			return nil, err
		}

		if len(postList.Posts) == 0 {
			break
		}

		for _, post := range postList.Posts {
			// Filter by time range
			if post.CreateAt >= startTime && post.CreateAt <= endTime {
				posts = append(posts, post)
				
				if len(posts) >= maxMessages {
					break
				}
			}
		}

		// If we've gone past the start time, we're done
		if len(postList.Order) > 0 {
			oldestPost := postList.Posts[postList.Order[len(postList.Order)-1]]
			if oldestPost.CreateAt < startTime {
				break
			}
		}

		page++
	}

	// Sort by creation time (oldest first for context)
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreateAt < posts[j].CreateAt
	})

	return posts, nil
}

// formatMessagesForLLM formats posts into message contexts for the LLM
func (a *App) formatMessagesForLLM(c request.CTX, posts []*model.Post) ([]*MessageContext, []string, *model.AppError) {
	contexts := make([]*MessageContext, 0, len(posts))
	participantIds := make(map[string]bool)

	for _, post := range posts {
		user, err := a.GetUser(post.UserId)
		if err != nil {
			// If we can't get the user, use a placeholder
			user = &model.User{
				Username:  "unknown",
				FirstName: "Unknown",
				LastName:  "User",
			}
		}

		participantIds[post.UserId] = true

		// Resolve mentions in message
		message := a.resolveMentions(post.Message, post.ChannelId)

		contexts = append(contexts, &MessageContext{
			Author:    user.GetDisplayName(model.ShowFullName),
			Username:  "@" + user.Username,
			Timestamp: post.CreateAt,
			Content:   message,
		})
	}

	// Build participant list
	participants := make([]string, 0, len(participantIds))
	for userId := range participantIds {
		user, err := a.GetUser(userId)
		if err == nil {
			participants = append(participants, user.GetDisplayName(model.ShowFullName))
		}
	}
	sort.Strings(participants)

	return contexts, participants, nil
}

// resolveMentions resolves @username mentions to display names
func (a *App) resolveMentions(message, channelId string) string {
	// Simple implementation - could be enhanced to actually resolve mentions
	// For now, just return the message as-is
	return message
}

// buildMessageText builds the formatted message text for the LLM prompt
func (a *App) buildMessageText(contexts []*MessageContext) string {
	var builder strings.Builder

	for i, ctx := range contexts {
		timestamp := time.UnixMilli(ctx.Timestamp).Format("Jan 02, 15:04")
		builder.WriteString(fmt.Sprintf("[%s] %s %s:\n%s\n\n",
			timestamp, ctx.Author, ctx.Username, ctx.Content))

		// Limit total output to prevent excessive token usage
		if i > 0 && i%100 == 0 && builder.Len() > 50000 {
			builder.WriteString("... (additional messages truncated for length) ...\n")
			break
		}
	}

	return builder.String()
}

// generateSummaryCacheKey generates a cache key for thread summaries
func (a *App) generateSummaryCacheKey(postId, summaryType, level string) string {
	data := fmt.Sprintf("%s:%s:%s", postId, summaryType, level)
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("summary:%s:%s", summaryType, hex.EncodeToString(hash[:])[:16])
}

// generateChannelSummaryCacheKey generates a cache key for channel summaries
func (a *App) generateChannelSummaryCacheKey(channelId string, startTime, endTime int64, level string) string {
	data := fmt.Sprintf("%s:%d:%d:%s", channelId, startTime, endTime, level)
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("summary:channel:%s", hex.EncodeToString(hash[:])[:16])
}

