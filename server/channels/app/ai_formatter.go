// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"strings"
	"time"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/shared/mlog"
	"github.com/mattermost/mattermost/server/public/shared/request"
	"github.com/mattermost/mattermost/server/v8/channels/app/openai"
)

// FormatMessage formats a message using the specified profile
func (a *App) FormatMessage(c request.CTX, req *FormattingRequest) (*FormattingResponse, *model.AppError) {
	startTime := time.Now()

	// Check if feature is enabled
	if !a.IsAIFeatureEnabled("formatting") {
		return nil, model.NewAppError("FormatMessage", "app.ai.formatting_disabled", nil, "", 403)
	}

	// Validate request
	if req.Message == "" {
		return nil, model.NewAppError("FormatMessage", "app.ai.invalid_message", nil, "", 400)
	}

	if req.Profile == "" {
		req.Profile = openai.FormattingProfessional
	}

	// Get AI service
	aiService := a.GetAIService()
	if aiService == nil {
		return nil, model.NewAppError("FormatMessage", "app.ai.service_not_available", nil, "", 500)
	}

	// Get prompt template
	promptTemplate := openai.GetMessageFormattingPrompt(req.Profile)
	if promptTemplate == nil {
		return nil, model.NewAppError("FormatMessage", "app.ai.invalid_profile", nil, "", 400)
	}

	// Build user prompt
	userPrompt := openai.BuildMessageFormattingUserPrompt(req.Message, req.Profile)
	
	// Add custom instructions if provided
	if req.CustomInstructions != "" {
		userPrompt += "\n\nAdditional instructions: " + req.CustomInstructions
	}

	// Call OpenAI
	systemPrompt, _ := promptTemplate.Substitute(nil)
	
	formattedText, err := aiService.client.SimpleCompletion(c.Context(), a.GetAIModel(), systemPrompt, userPrompt)
	if err != nil {
		c.Logger().Error("Failed to format message", mlog.Err(err))
		return nil, model.NewAppError("FormatMessage", "app.ai.formatting_failed", nil, err.Error(), 500)
	}

	formattedText = strings.TrimSpace(formattedText)

	// Generate diff for preview
	diff := a.generateTextDiff(req.Message, formattedText)

	processingMs := time.Since(startTime).Milliseconds()

	c.Logger().Debug("Message formatted successfully",
		mlog.String("profile", string(req.Profile)),
		mlog.Int("processing_ms", int(processingMs)),
		mlog.Int("original_length", len(req.Message)),
		mlog.Int("formatted_length", len(formattedText)))

	return &FormattingResponse{
		FormattedText: formattedText,
		Profile:       req.Profile,
		Diff:          diff,
		ProcessingMs:  processingMs,
	}, nil
}

// PreviewFormatting returns a preview of formatted text without applying it
func (a *App) PreviewFormatting(c request.CTX, req *FormattingRequest) (*FormattingResponse, *model.AppError) {
	// Preview is the same as format, but we always include diff
	response, err := a.FormatMessage(c, req)
	if err != nil {
		return nil, err
	}

	// Ensure diff is included
	if response.Diff == nil {
		response.Diff = a.generateTextDiff(req.Message, response.FormattedText)
	}

	return response, nil
}

// generateTextDiff generates a simple diff between original and formatted text
// This is a basic implementation - could be enhanced with a proper diff library
func (a *App) generateTextDiff(original, formatted string) *TextDiff {
	// Simple character-by-character comparison for now
	// In production, you might want to use a proper diff algorithm
	changes := []Change{}

	// If texts are identical, return empty diff
	if original == formatted {
		return &TextDiff{
			Original:  original,
			Formatted: formatted,
			Changes:   changes,
		}
	}

	// Simple approach: if texts differ, mark as replacement
	// A more sophisticated implementation would use a diff algorithm
	if original != formatted {
		changes = append(changes, Change{
			Type:     "replace",
			Start:    0,
			End:      len(original),
			OldText:  original,
			NewText:  formatted,
		})
	}

	return &TextDiff{
		Original:  original,
		Formatted: formatted,
		Changes:   changes,
	}
}

// GetFormattingProfiles returns available formatting profiles
func (a *App) GetFormattingProfiles(c request.CTX) ([]FormattingProfileInfo, *model.AppError) {
	metadata := GetFormattingProfileMetadata()
	profiles := make([]FormattingProfileInfo, 0, len(metadata))
	
	for _, info := range metadata {
		profiles = append(profiles, info)
	}

	return profiles, nil
}

