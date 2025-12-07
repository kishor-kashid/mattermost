// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"encoding/json"
	"net/http"

	"github.com/mattermost/mattermost/server/public/shared/mlog"
	"github.com/mattermost/mattermost/server/v8/channels/app"
)

type LinkSummarizeRequest struct {
	URL          string `json:"url"`
	UseCache     bool   `json:"use_cache"`
	ForceRefresh bool   `json:"force_refresh"`
}

func (api *API) initLinkSummarizerRoutes() {
	api.BaseRoutes.AI.Handle("/links/summarize", api.APISessionRequired(summarizeLink)).Methods(http.MethodPost)
	api.BaseRoutes.AI.Handle("/links/summary", api.APISessionRequired(getLinkSummary)).Methods(http.MethodGet)
}

// summarizeLink handles POST /api/v4/ai/links/summarize
func summarizeLink(c *Context, w http.ResponseWriter, r *http.Request) {
	if !requireAIEnabled(c) {
		return
	}

	var req LinkSummarizeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.SetInvalidParamWithErr("body", err)
		return
	}

	if req.URL == "" {
		c.SetInvalidParam("url")
		return
	}

	linkReq := &app.LinkSummarizationRequest{
		UserId:       c.AppContext.Session().UserId,
		URL:          req.URL,
		UseCache:     req.UseCache,
		ForceRefresh: req.ForceRefresh,
	}

	resp, err := c.App.SummarizeLink(c.AppContext, linkReq)
	if err != nil {
		c.Err = err
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		c.Logger.Error("Failed to encode link summary response", mlog.Err(err))
	}
}

// getLinkSummary handles GET /api/v4/ai/links/summary?url=
func getLinkSummary(c *Context, w http.ResponseWriter, r *http.Request) {
	if !requireAIEnabled(c) {
		return
	}

	urlParam := r.URL.Query().Get("url")
	if urlParam == "" {
		c.SetInvalidParam("url")
		return
	}

	useCache := r.URL.Query().Get("use_cache") != "false"

	linkReq := &app.LinkSummarizationRequest{
		UserId:   c.AppContext.Session().UserId,
		URL:      urlParam,
		UseCache: useCache,
	}

	resp, err := c.App.SummarizeLink(c.AppContext, linkReq)
	if err != nil {
		c.Err = err
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		c.Logger.Error("Failed to encode link summary response", mlog.Err(err))
	}
}
