// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"net/http"

	"github.com/mattermost/mattermost/server/public/model"
)

// InitAI initializes the AI API routes
func (api *API) InitAI() {
	// AI root endpoints
	api.BaseRoutes.AI.Handle("/health", api.APISessionRequired(aiHealthCheck)).Methods(http.MethodGet)
	api.BaseRoutes.AI.Handle("/config/validate", api.APISessionRequired(aiValidateConfig)).Methods(http.MethodPost)
	api.BaseRoutes.AI.Handle("/test", api.APISessionRequired(aiTestConnection)).Methods(http.MethodPost)

	// Feature-specific routes
	api.initSummarizerRoutes()
}

// requireAIEnabled checks if AI features are enabled in the configuration
func requireAIEnabled(c *Context) bool {
	if c.App.Config().AISettings.Enable == nil || !*c.App.Config().AISettings.Enable {
		c.Err = model.NewAppError("requireAIEnabled", "api.ai.disabled.app_error", nil, "", http.StatusNotImplemented)
		return false
	}
	return true
}

// checkAIPermissions validates that the user has permission to use AI features in the given channel
func checkAIPermissions(c *Context, channelID string) bool {
	if !requireAIEnabled(c) {
		return false
	}

	// Check if user has access to the channel
	if !c.App.SessionHasPermissionToChannel(c.AppContext, *c.AppContext.Session(), channelID, model.PermissionReadChannel) {
		c.SetPermissionError(model.PermissionReadChannel)
		return false
	}

	return true
}

