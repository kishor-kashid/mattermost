// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/v8/channels/app/openai"
)

// InitializeAI initializes AI services based on configuration
func (a *App) InitializeAI() error {
	if a.Config().AISettings.Enable == nil || !*a.Config().AISettings.Enable {
		a.Log().Debug("AI features are disabled in configuration")
		return nil
	}

	apiKey := a.Config().AISettings.OpenAIAPIKey
	if apiKey == nil || *apiKey == "" {
		a.Log().Warn("AI features enabled but no OpenAI API key configured")
		return nil
	}

	// Create OpenAI client
	clientConfig := openai.ClientConfig{
		APIKey: *apiKey,
		Logger: a.Log(),
	}

	client := openai.NewClient(clientConfig)

	// Store the client in app context (we'll add this field to App later)
	// For now, we're just setting up the initialization pattern
	a.Log().Info("AI services initialized successfully")

	// Test connectivity (optional - can be enabled for debugging)
	// ctx := context.Background()
	// _, err := client.SimpleCompletion(ctx, "gpt-3.5-turbo", "", "Say 'hello'")
	// if err != nil {
	// 	a.Log().Error("Failed to connect to OpenAI API", mlog.Err(err))
	// 	return err
	// }

	_ = client // Suppress unused variable warning for now

	return nil
}

// GetAIModel returns the configured AI model or the default
func (a *App) GetAIModel() string {
	if a.Config().AISettings.OpenAIModel != nil && *a.Config().AISettings.OpenAIModel != "" {
		return *a.Config().AISettings.OpenAIModel
	}
	return "gpt-3.5-turbo"
}

// IsAIFeatureEnabled checks if a specific AI feature is enabled
func (a *App) IsAIFeatureEnabled(feature string) bool {
	if a.Config().AISettings.Enable == nil || !*a.Config().AISettings.Enable {
		return false
	}

	switch feature {
	case "summarization":
		return a.Config().AISettings.EnableSummarization != nil && *a.Config().AISettings.EnableSummarization
	case "analytics":
		return a.Config().AISettings.EnableAnalytics != nil && *a.Config().AISettings.EnableAnalytics
	case "actionitems":
		return a.Config().AISettings.EnableActionItems != nil && *a.Config().AISettings.EnableActionItems
	case "formatting":
		return a.Config().AISettings.EnableFormatting != nil && *a.Config().AISettings.EnableFormatting
	default:
		return false
	}
}

// GetAIMaxMessageLimit returns the maximum number of messages to process for AI operations
func (a *App) GetAIMaxMessageLimit() int {
	if a.Config().AISettings.MaxMessageLimit != nil {
		return *a.Config().AISettings.MaxMessageLimit
	}
	return 500
}

// GetAIRateLimit returns the API rate limit per minute
func (a *App) GetAIRateLimit() int {
	if a.Config().AISettings.APIRateLimit != nil {
		return *a.Config().AISettings.APIRateLimit
	}
	return 60
}

// GetOrCreateAIPreferences gets or creates AI preferences for a user
func (a *App) GetOrCreateAIPreferences(userId string) (*model.AIPreferences, error) {
	// Try to get existing preferences
	preferences, err := a.Srv().Store().AIPreferences().GetByUser(userId)
	if err == nil {
		return preferences, nil
	}

	// Create default preferences if not found
	preferences = &model.AIPreferences{
		UserId:              userId,
		EnableSummarization: true,
		EnableAnalytics:     true,
		EnableActionItems:   true,
		EnableFormatting:    true,
		DefaultModel:        a.GetAIModel(),
		FormattingProfile:   "professional",
	}

	preferences, err = a.Srv().Store().AIPreferences().Save(preferences)
	if err != nil {
		return nil, err
	}

	return preferences, nil
}

