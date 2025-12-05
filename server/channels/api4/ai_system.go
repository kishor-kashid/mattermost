// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"encoding/json"
	"net/http"

	"github.com/mattermost/mattermost/server/public/model"
)

// aiHealthCheck checks if AI services are healthy and operational
func aiHealthCheck(c *Context, w http.ResponseWriter, r *http.Request) {
	if !requireAIEnabled(c) {
		return
	}

	// Check if OpenAI client is initialized
	aiService := c.App.GetAIService()
	if aiService == nil {
		c.Err = model.NewAppError("aiHealthCheck", "api.ai.service_not_initialized.app_error", nil, "", http.StatusServiceUnavailable)
		return
	}

	// Build health check response
	health := map[string]interface{}{
		"enabled":            true,
		"service_available":  aiService != nil,
		"openai_configured":  c.App.Config().AISettings.OpenAIAPIKey != nil && *c.App.Config().AISettings.OpenAIAPIKey != "",
		"features": map[string]bool{
			"summarization": c.App.Config().AISettings.EnableSummarization != nil && *c.App.Config().AISettings.EnableSummarization,
			"analytics":     c.App.Config().AISettings.EnableAnalytics != nil && *c.App.Config().AISettings.EnableAnalytics,
			"action_items":  c.App.Config().AISettings.EnableActionItems != nil && *c.App.Config().AISettings.EnableActionItems,
			"formatting":    c.App.Config().AISettings.EnableFormatting != nil && *c.App.Config().AISettings.EnableFormatting,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(health); err != nil {
		c.Err = model.NewAppError("aiHealthCheck", "api.marshal_error", nil, err.Error(), http.StatusInternalServerError)
		return
	}
}

// AIConfigValidateRequest represents a request to validate AI configuration
type AIConfigValidateRequest struct {
	OpenAIAPIKey string `json:"openai_api_key"`
	Model        string `json:"model"`
}

// aiValidateConfig validates AI configuration settings
func aiValidateConfig(c *Context, w http.ResponseWriter, r *http.Request) {
	// Require system admin permission to validate config
	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionManageSystem) {
		c.SetPermissionError(model.PermissionManageSystem)
		return
	}

	var req AIConfigValidateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.SetInvalidParamWithErr("body", err)
		return
	}

	// Validate API key format
	if req.OpenAIAPIKey == "" {
		c.SetInvalidParam("openai_api_key")
		return
	}

	// Validate model name
	validModels := []string{"gpt-4", "gpt-4-turbo-preview", "gpt-3.5-turbo"}
	modelValid := false
	for _, validModel := range validModels {
		if req.Model == validModel {
			modelValid = true
			break
		}
	}

	if !modelValid {
		c.SetInvalidParam("model")
		return
	}

	response := map[string]interface{}{
		"valid":   true,
		"message": "Configuration is valid",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		c.Err = model.NewAppError("aiValidateConfig", "api.marshal_error", nil, err.Error(), http.StatusInternalServerError)
		return
	}
}

// AITestConnectionRequest represents a request to test OpenAI connectivity
type AITestConnectionRequest struct {
	TestPrompt string `json:"test_prompt,omitempty"`
}

// aiTestConnection tests the connection to OpenAI API
func aiTestConnection(c *Context, w http.ResponseWriter, r *http.Request) {
	if !requireAIEnabled(c) {
		return
	}

	// Require system admin permission to test connection
	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionManageSystem) {
		c.SetPermissionError(model.PermissionManageSystem)
		return
	}

	var req AITestConnectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Empty request is valid
		req.TestPrompt = "Hello"
	}

	if req.TestPrompt == "" {
		req.TestPrompt = "Hello"
	}

	// Get AI service
	aiService := c.App.GetAIService()
	if aiService == nil {
		c.Err = model.NewAppError("aiTestConnection", "api.ai.service_not_initialized.app_error", nil, "", http.StatusServiceUnavailable)
		return
	}

	// Test OpenAI connection with a simple prompt
	result, err := aiService.TestConnection(c.AppContext.Context(), req.TestPrompt)
	if err != nil {
		response := map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
			c.Err = model.NewAppError("aiTestConnection", "api.marshal_error", nil, encodeErr.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := map[string]interface{}{
		"success": true,
		"result":  result,
		"message": "Successfully connected to OpenAI API",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		c.Err = model.NewAppError("aiTestConnection", "api.marshal_error", nil, err.Error(), http.StatusInternalServerError)
		return
	}
}

