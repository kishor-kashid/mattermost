// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"encoding/json"
	"net/http"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/shared/mlog"
	"github.com/mattermost/mattermost/server/v8/channels/app"
)

// SummarizeRequest represents the API request body for summarization
type SummarizeRequest struct {
	ChannelId    string `json:"channel_id"`
	PostId       string `json:"post_id,omitempty"`       // For thread summarization
	StartTime    int64  `json:"start_time,omitempty"`    // For channel summarization
	EndTime      int64  `json:"end_time,omitempty"`      // For channel summarization
	SummaryLevel string `json:"summary_level,omitempty"` // brief, standard, detailed
	UseCache     bool   `json:"use_cache"`               // Whether to use cached summaries
}

// SummarizeResponse represents the API response for summarization
type SummarizeResponse struct {
	Summary      *model.AISummary `json:"summary"`
	FromCache    bool             `json:"from_cache"`
	ProcessingMs int64            `json:"processing_ms"`
}

func (api *API) initSummarizerRoutes() {
	api.BaseRoutes.AI.Handle("/summarize", api.APISessionRequired(summarize)).Methods("POST")
	api.BaseRoutes.AI.Handle("/summarize/thread/{post_id:[A-Za-z0-9]+}", api.APISessionRequired(summarizeThread)).Methods("GET")
	api.BaseRoutes.AI.Handle("/summarize/channel/{channel_id:[A-Za-z0-9]+}", api.APISessionRequired(summarizeChannel)).Methods("POST")
}

// summarize handles POST /api/v4/ai/summarize
func summarize(c *Context, w http.ResponseWriter, r *http.Request) {
	var req SummarizeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.SetInvalidParamWithErr("body", err)
		return
	}

	// Validate request
	if req.ChannelId == "" {
		c.SetInvalidParam("channel_id")
		return
	}

	// Determine if this is thread or channel summarization
	if req.PostId != "" {
		handleThreadSummarization(c, w, &req)
	} else {
		handleChannelSummarization(c, w, &req)
	}
}

// summarizeThread handles GET /api/v4/ai/summarize/thread/{post_id}
func summarizeThread(c *Context, w http.ResponseWriter, r *http.Request) {
	postId := c.Params.PostId

	// Get optional query parameters
	summaryLevel := r.URL.Query().Get("level")
	if summaryLevel == "" {
		summaryLevel = "standard"
	}

	useCache := r.URL.Query().Get("use_cache") != "false" // Default to true

	// Get the post to find the channel
	post, err := c.App.GetSinglePost(c.AppContext, postId, false)
	if err != nil {
		c.Err = err
		return
	}

	req := SummarizeRequest{
		ChannelId:    post.ChannelId,
		PostId:       postId,
		SummaryLevel: summaryLevel,
		UseCache:     useCache,
	}

	handleThreadSummarization(c, w, &req)
}

// summarizeChannel handles POST /api/v4/ai/summarize/channel/{channel_id}
func summarizeChannel(c *Context, w http.ResponseWriter, r *http.Request) {
	channelId := c.Params.ChannelId

	var req SummarizeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.SetInvalidParamWithErr("body", err)
		return
	}

	req.ChannelId = channelId

	// Set defaults
	if req.SummaryLevel == "" {
		req.SummaryLevel = "standard"
	}

	handleChannelSummarization(c, w, &req)
}

// handleThreadSummarization processes thread summarization requests
func handleThreadSummarization(c *Context, w http.ResponseWriter, req *SummarizeRequest) {
	// Build app request
	appRequest := &app.SummarizationRequest{
		ChannelId:    req.ChannelId,
		PostId:       req.PostId,
		SummaryLevel: req.SummaryLevel,
		MaxMessages:  c.App.GetAIMaxMessageLimit(),
		UserId:       c.AppContext.Session().UserId,
		UseCache:     req.UseCache,
	}

	// Execute summarization
	result, err := c.App.SummarizeThread(c.AppContext, appRequest)
	if err != nil {
		c.Err = err
		return
	}

	// Build response
	response := SummarizeResponse{
		Summary:      result.Summary,
		FromCache:    result.FromCache,
		ProcessingMs: result.ProcessingMs,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		c.Logger.Warn("Failed to encode summarization response", mlog.Err(err))
	}
}

// handleChannelSummarization processes channel summarization requests
func handleChannelSummarization(c *Context, w http.ResponseWriter, req *SummarizeRequest) {
	// Set defaults for time range if not provided
	if req.EndTime == 0 {
		req.EndTime = model.GetMillis()
	}
	if req.StartTime == 0 {
		// Default to last 24 hours
		req.StartTime = req.EndTime - (24 * 60 * 60 * 1000)
	}

	// Build app request
	appRequest := &app.SummarizationRequest{
		ChannelId:    req.ChannelId,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		SummaryLevel: req.SummaryLevel,
		MaxMessages:  c.App.GetAIMaxMessageLimit(),
		UserId:       c.AppContext.Session().UserId,
		UseCache:     req.UseCache,
	}

	// Execute summarization
	result, err := c.App.SummarizeChannel(c.AppContext, appRequest)
	if err != nil {
		c.Err = err
		return
	}

	// Build response
	response := SummarizeResponse{
		Summary:      result.Summary,
		FromCache:    result.FromCache,
		ProcessingMs: result.ProcessingMs,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		c.Logger.Warn("Failed to encode summarization response", mlog.Err(err))
	}
}

