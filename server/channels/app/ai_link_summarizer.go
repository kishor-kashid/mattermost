// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/shared/mlog"
	"github.com/mattermost/mattermost/server/public/shared/request"
	"github.com/mattermost/mattermost/server/v8/channels/app/openai"
)

// SummarizeLink generates or retrieves a summary for a given URL.
func (a *App) SummarizeLink(c request.CTX, req *LinkSummarizationRequest) (*LinkSummarizationResponse, *model.AppError) {
	start := time.Now()

	if !a.IsAIFeatureEnabled("link_summarization") {
		return nil, model.NewAppError("SummarizeLink", "app.ai.link_summarization_disabled", nil, "", 403)
	}

	if req == nil || req.URL == "" {
		return nil, model.NewAppError("SummarizeLink", "app.ai.invalid_url", nil, "", 400)
	}

	parsed, err := url.Parse(req.URL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return nil, model.NewAppError("SummarizeLink", "app.ai.invalid_url", nil, "", 400)
	}

	urlHash := hashURL(req.URL)

	// Cache check
	if req.UseCache && !req.ForceRefresh {
		if cached, err := a.Srv().Store().AILinkSummary().GetByURLHash(urlHash); err == nil && cached != nil {
			return &LinkSummarizationResponse{
				Summary:      modelToLinkSummary(cached),
				FromCache:    true,
				ProcessingMs: time.Since(start).Milliseconds(),
			}, nil
		}
	}

	// Fetch content
	ctx, cancel := context.WithTimeout(c.Context(), defaultFetchTimeout)
	defer cancel()

	body, contentType, fetchErr := a.fetchURL(ctx, req.URL)
	if fetchErr != nil {
		return nil, model.NewAppError("SummarizeLink", "app.ai.fetch_failed", nil, fetchErr.Error(), 500)
	}

	content := a.extractContent(req.URL, body, contentType)

	// Build prompts
	prefLength := LinkSummaryLengthStandard
	if prefs, _ := a.GetOrCreateAIPreferences(req.UserId); prefs != nil && prefs.LinkSummaryLength != "" {
		prefLength = LinkSummaryLength(prefs.LinkSummaryLength)
	}

	systemPrompt, userPrompt := openai.GetLinkSummaryPrompt(ToPromptLength(string(prefLength))).Substitute(map[string]string{
		"title":        safeFallback(content.Title, "Untitled"),
		"description":  safeFallback(content.Description, "No description provided."),
		"text":         truncate(content.Text, 4000),
		"domain":       content.Domain,
		"content_type": content.ContentType,
		"reading_time": formatMinutes(content.ReadingTime),
	})

	aiService := a.GetAIService()
	if aiService == nil {
		return nil, model.NewAppError("SummarizeLink", "app.ai.service_not_available", nil, "", 500)
	}

	summaryText, openaiErr := aiService.client.SimpleCompletion(c.Context(), a.GetAIModel(), systemPrompt, userPrompt)
	if openaiErr != nil {
		a.Log().Error("Failed to generate link summary", mlog.Err(openaiErr))
		return nil, model.NewAppError("SummarizeLink", "app.ai.openai_error", nil, openaiErr.Error(), 500)
	}

	// Derive key points (simple split by lines beginning with "-")
	keyPoints := extractKeyPoints(summaryText)

	linkSummary := &model.AILinkSummary{
		Url:         req.URL,
		UrlHash:     urlHash,
		Title:       content.Title,
		Description: content.Description,
		Summary:     strings.TrimSpace(summaryText),
		KeyPoints:   keyPoints,
		ContentType: content.ContentType,
		ReadingTime: content.ReadingTime,
		Domain:      content.Domain,
		FaviconURL:  content.FaviconURL,
	}
	linkSummary.PreSave()

	// Best-effort cleanup of expired cache entries
	_, _ = a.Srv().Store().AILinkSummary().DeleteExpired(model.GetMillis())

	saved, saveErr := a.Srv().Store().AILinkSummary().Save(linkSummary)
	if saveErr != nil {
		a.Log().Warn("Failed to cache link summary", mlog.Err(saveErr))
	}

	return &LinkSummarizationResponse{
		Summary:      modelToLinkSummary(saved),
		FromCache:    false,
		ProcessingMs: time.Since(start).Milliseconds(),
	}, nil
}

func hashURL(u string) string {
	h := sha256.Sum256([]byte(strings.ToLower(strings.TrimSpace(u))))
	return hex.EncodeToString(h[:])
}

func safeFallback(s, fallback string) string {
	if strings.TrimSpace(s) == "" {
		return fallback
	}
	return s
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max]
}

func formatMinutes(min int) string {
	if min <= 0 {
		return "Unknown"
	}
	return strings.TrimSpace(strings.Join([]string{strconv.Itoa(min), "minutes"}, " "))
}

func extractKeyPoints(summary string) []string {
	var points []string
	lines := strings.Split(summary, "\n")
	for _, line := range lines {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "-") {
			points = append(points, strings.TrimSpace(strings.TrimPrefix(trim, "-")))
		}
	}
	return points
}

func modelToLinkSummary(m *model.AILinkSummary) *LinkSummary {
	if m == nil {
		return nil
	}
	return &LinkSummary{
		Title:       m.Title,
		Description: m.Description,
		Summary:     m.Summary,
		KeyPoints:   m.KeyPoints,
		ContentType: m.ContentType,
		ReadingTime: m.ReadingTime,
		Domain:      m.Domain,
		FaviconURL:  m.FaviconURL,
	}
}
