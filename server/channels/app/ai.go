// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"context"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/shared/mlog"
	"github.com/mattermost/mattermost/server/v8/channels/app/openai"
)

// AIService provides AI-powered features
type AIService struct {
	app    *App
	client *openai.Client
	logger mlog.LoggerIFace
}

// InitializeAI initializes AI services based on configuration
func (a *App) InitializeAI() error {
	if a.Config().AISettings.Enable == nil || !*a.Config().AISettings.Enable {
		a.Log().Debug("AI features are disabled in configuration")
		return nil
	}

	apiKey := a.Config().AISettings.OpenAIAPIKey
	if apiKey == nil || *apiKey == "" {
		a.Log().Warn("AI features enabled but no OpenAI API key configured",
			mlog.String("hint", "Set MM_AISETTINGS_OPENAIAPIKEY environment variable or configure in config.json"))
		return nil
	}

	// Mask the API key for logging (show first 7 chars only)
	maskedKey := "sk-..."
	if len(*apiKey) > 10 {
		maskedKey = (*apiKey)[:7] + "..." + (*apiKey)[len(*apiKey)-4:]
	}

	// Create OpenAI client
	clientConfig := openai.ClientConfig{
		APIKey: *apiKey,
		Logger: a.Log(),
	}

	client := openai.NewClient(clientConfig)

	// Store the AI service (note: this requires adding aiService field to App/Server struct)
	// For now, we're setting up the pattern
	a.Log().Info("AI services initialized successfully", 
		mlog.String("api_key", maskedKey),
		mlog.String("model", a.GetAIModel()))

	_ = client // Will be used by GetAIService

	return nil
}

// GetAIService returns the AI service instance
func (a *App) GetAIService() *AIService {
	if a.Config().AISettings.Enable == nil || !*a.Config().AISettings.Enable {
		a.Log().Debug("AI service requested but AI features are disabled")
		return nil
	}

	apiKey := a.Config().AISettings.OpenAIAPIKey
	if apiKey == nil || *apiKey == "" {
		a.Log().Warn("AI service requested but no OpenAI API key is configured",
			mlog.String("hint", "Set MM_AISETTINGS_OPENAIAPIKEY environment variable"))
		return nil
	}

	// Create OpenAI client
	clientConfig := openai.ClientConfig{
		APIKey: *apiKey,
		Logger: a.Log(),
	}

	client := openai.NewClient(clientConfig)

	return &AIService{
		app:    a,
		client: client,
		logger: a.Log(),
	}
}

// TestConnection tests the OpenAI API connection
func (s *AIService) TestConnection(ctx context.Context, testPrompt string) (string, error) {
	model := s.app.GetAIModel()

	result, err := s.client.SimpleCompletion(ctx, model, "", testPrompt)
	if err != nil {
		s.logger.Error("OpenAI connection test failed", mlog.Err(err))
		return "", err
	}

	s.logger.Debug("OpenAI connection test successful")
	return result, nil
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
	// Debug logging for AI feature check
	aiEnabled := a.Config().AISettings.Enable != nil && *a.Config().AISettings.Enable
	
	a.Log().Debug("Checking AI feature status",
		mlog.String("feature", feature),
		mlog.Bool("ai_enabled", aiEnabled),
		mlog.Bool("enable_ptr_nil", a.Config().AISettings.Enable == nil))
	
	if !aiEnabled {
		a.Log().Debug("AI features are disabled globally",
			mlog.String("feature_requested", feature))
		return false
	}

	var featureEnabled bool
	switch feature {
	case "summarization":
		featureEnabled = a.Config().AISettings.EnableSummarization != nil && *a.Config().AISettings.EnableSummarization
		a.Log().Debug("Summarization feature check",
			mlog.Bool("summarization_ptr_nil", a.Config().AISettings.EnableSummarization == nil),
			mlog.Bool("enabled", featureEnabled))
	case "analytics":
		featureEnabled = a.Config().AISettings.EnableAnalytics != nil && *a.Config().AISettings.EnableAnalytics
	case "actionitems":
		featureEnabled = a.Config().AISettings.EnableActionItems != nil && *a.Config().AISettings.EnableActionItems
	case "formatting":
		featureEnabled = a.Config().AISettings.EnableFormatting != nil && *a.Config().AISettings.EnableFormatting
	default:
		a.Log().Warn("Unknown AI feature requested", mlog.String("feature", feature))
		return false
	}
	
	a.Log().Debug("AI feature enabled check result",
		mlog.String("feature", feature),
		mlog.Bool("result", featureEnabled))
	
	return featureEnabled
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

