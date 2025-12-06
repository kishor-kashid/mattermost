// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"encoding/json"
	"net/http"

	"github.com/mattermost/mattermost/server/public/shared/mlog"
	"github.com/mattermost/mattermost/server/v8/channels/app"
	"github.com/mattermost/mattermost/server/v8/channels/app/openai"
)

// FormatPreviewRequest represents the API request body for formatting preview
type FormatPreviewRequest struct {
	Message            string `json:"message"`
	Profile            string `json:"profile,omitempty"`
	CustomInstructions string `json:"custom_instructions,omitempty"`
}

// FormatApplyRequest represents the API request body for applying formatting
type FormatApplyRequest struct {
	Message            string `json:"message"`
	Profile            string `json:"profile,omitempty"`
	CustomInstructions string `json:"custom_instructions,omitempty"`
}

func (api *API) initFormatterRoutes() {
	api.BaseRoutes.AI.Handle("/format/preview", api.APISessionRequired(formatPreview)).Methods("POST")
	api.BaseRoutes.AI.Handle("/format/apply", api.APISessionRequired(formatApply)).Methods("POST")
	api.BaseRoutes.AI.Handle("/format/profiles", api.APISessionRequired(getFormattingProfiles)).Methods("GET")
}

// formatPreview handles POST /api/v4/ai/format/preview
func formatPreview(c *Context, w http.ResponseWriter, r *http.Request) {
	if !requireAIEnabled(c) {
		return
	}

	var req FormatPreviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.SetInvalidParamWithErr("body", err)
		return
	}

	// Validate request
	if req.Message == "" {
		c.SetInvalidParam("message")
		return
	}

	// Convert profile string to FormattingProfile
	profile := openai.FormattingProfessional
	if req.Profile != "" {
		if !app.IsValidFormattingProfile(req.Profile) {
			c.SetInvalidParam("profile")
			return
		}
		profile = openai.FormattingProfile(req.Profile)
	} else {
		// Get user's default profile from preferences
		preferences, err := c.App.GetOrCreateAIPreferences(c.AppContext.Session().UserId)
		if err == nil && preferences.FormattingProfile != "" {
			profile = openai.FormattingProfile(preferences.FormattingProfile)
		}
	}

	formattingReq := &app.FormattingRequest{
		Message:            req.Message,
		Profile:            profile,
		CustomInstructions: req.CustomInstructions,
	}

	response, err := c.App.PreviewFormatting(c.AppContext, formattingReq)
	if err != nil {
		c.Err = err
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		c.Logger.Error("Failed to encode format preview response", mlog.Err(err))
	}
}

// formatApply handles POST /api/v4/ai/format/apply
func formatApply(c *Context, w http.ResponseWriter, r *http.Request) {
	if !requireAIEnabled(c) {
		return
	}

	var req FormatApplyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.SetInvalidParamWithErr("body", err)
		return
	}

	// Validate request
	if req.Message == "" {
		c.SetInvalidParam("message")
		return
	}

	// Convert profile string to FormattingProfile
	profile := openai.FormattingProfessional
	if req.Profile != "" {
		if !app.IsValidFormattingProfile(req.Profile) {
			c.SetInvalidParam("profile")
			return
		}
		profile = openai.FormattingProfile(req.Profile)
	} else {
		// Get user's default profile from preferences
		preferences, err := c.App.GetOrCreateAIPreferences(c.AppContext.Session().UserId)
		if err == nil && preferences.FormattingProfile != "" {
			profile = openai.FormattingProfile(preferences.FormattingProfile)
		}
	}

	formattingReq := &app.FormattingRequest{
		Message:            req.Message,
		Profile:            profile,
		CustomInstructions: req.CustomInstructions,
	}

	response, err := c.App.FormatMessage(c.AppContext, formattingReq)
	if err != nil {
		c.Err = err
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		c.Logger.Error("Failed to encode format apply response", mlog.Err(err))
	}
}

// getFormattingProfiles handles GET /api/v4/ai/format/profiles
func getFormattingProfiles(c *Context, w http.ResponseWriter, r *http.Request) {
	if !requireAIEnabled(c) {
		return
	}

	profiles, err := c.App.GetFormattingProfiles(c.AppContext)
	if err != nil {
		c.Err = err
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(profiles); err != nil {
		c.Logger.Error("Failed to encode formatting profiles response", mlog.Err(err))
	}
}

