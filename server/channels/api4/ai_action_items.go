// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/shared/mlog"
	"github.com/mattermost/mattermost/server/v8/channels/app"
)

func (api *API) InitAIActionItemsRoutes() {
	api.BaseRoutes.AI.Handle("/actionitems", api.APISessionRequired(createActionItem)).Methods("POST")
	api.BaseRoutes.AI.Handle("/actionitems", api.APISessionRequired(getActionItems)).Methods("GET")
	api.BaseRoutes.AI.Handle("/actionitems/{action_item_id:[A-Za-z0-9]+}", api.APISessionRequired(getActionItem)).Methods("GET")
	api.BaseRoutes.AI.Handle("/actionitems/{action_item_id:[A-Za-z0-9]+}", api.APISessionRequired(updateActionItem)).Methods("PUT")
	api.BaseRoutes.AI.Handle("/actionitems/{action_item_id:[A-Za-z0-9]+}", api.APISessionRequired(deleteActionItem)).Methods("DELETE")
	api.BaseRoutes.AI.Handle("/actionitems/{action_item_id:[A-Za-z0-9]+}/complete", api.APISessionRequired(completeActionItem)).Methods("POST")
	api.BaseRoutes.AI.Handle("/actionitems/stats", api.APISessionRequired(getActionItemStats)).Methods("GET")
}

func createActionItem(c *Context, w http.ResponseWriter, r *http.Request) {
	var req app.ActionItemCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.SetInvalidParamWithErr("body", err)
		return
	}

	// Build action item from request
	actionItem := &model.AIActionItem{
		Description: req.Description,
		AssigneeId:  req.AssigneeID,
		ChannelId:   req.ChannelID,
		PostId:      req.PostID,
		Priority:    req.Priority,
		Status:      req.Status,
		CreatedBy:   c.AppContext.Session().UserId,
	}

	if req.DueDate != nil {
		actionItem.DueDate = req.DueDate.UnixMilli()
	}

	created, err := c.App.CreateActionItem(c.AppContext, actionItem)
	if err != nil {
		c.Err = model.NewAppError("createActionItem", "api.action_item.create.app_error", nil, err.Error(), http.StatusInternalServerError)
		return
	}

	c.LogAudit("action_item_id=" + created.Id)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(created); err != nil {
		c.Logger.Warn("Error writing response", mlog.Err(err))
	}
}

func getActionItem(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireActionItemId()
	if c.Err != nil {
		return
	}

	actionItem, err := c.App.GetActionItem(c.AppContext, c.Params.ActionItemId, c.AppContext.Session().UserId)
	if err != nil {
		c.Err = model.NewAppError("getActionItem", "api.action_item.get.app_error", nil, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(actionItem); err != nil {
		c.Logger.Warn("Error writing response", mlog.Err(err))
	}
}

func getActionItems(c *Context, w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	userID := r.URL.Query().Get("user_id")
	channelID := r.URL.Query().Get("channel_id")
	status := r.URL.Query().Get("status")
	priority := r.URL.Query().Get("priority")
	includeCompleted := r.URL.Query().Get("include_completed") == "true"
	
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage == 0 {
		perPage = 60
	}

	filters := &app.ActionItemFilters{
		Status:           status,
		Priority:         priority,
		IncludeCompleted: includeCompleted,
		Page:             page,
		PerPage:          perPage,
	}

	var items []*model.AIActionItem
	var err error

	// If user_id is not specified, default to current user
	if userID == "" {
		userID = c.AppContext.Session().UserId
	}

	// Determine which items to fetch
	if channelID != "" {
		// Get items for a specific channel
		items, err = c.App.GetActionItemsForChannel(c.AppContext, channelID, c.AppContext.Session().UserId, filters)
	} else {
		// Get items for a user
		items, err = c.App.GetActionItemsForUser(c.AppContext, userID, filters)
	}

	if err != nil {
		c.Err = model.NewAppError("getActionItems", "api.action_item.list.app_error", nil, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(items); err != nil {
		c.Logger.Warn("Error writing response", mlog.Err(err))
	}
}

func updateActionItem(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireActionItemId()
	if c.Err != nil {
		return
	}

	var req app.ActionItemUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.SetInvalidParamWithErr("body", err)
		return
	}

	updated, err := c.App.UpdateActionItem(c.AppContext, c.Params.ActionItemId, c.AppContext.Session().UserId, &req)
	if err != nil {
		c.Err = model.NewAppError("updateActionItem", "api.action_item.update.app_error", nil, err.Error(), http.StatusInternalServerError)
		return
	}

	c.LogAudit("action_item_id=" + updated.Id)
	if err := json.NewEncoder(w).Encode(updated); err != nil {
		c.Logger.Warn("Error writing response", mlog.Err(err))
	}
}

func completeActionItem(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireActionItemId()
	if c.Err != nil {
		return
	}

	completed, err := c.App.CompleteActionItem(c.AppContext, c.Params.ActionItemId, c.AppContext.Session().UserId)
	if err != nil {
		c.Err = model.NewAppError("completeActionItem", "api.action_item.complete.app_error", nil, err.Error(), http.StatusInternalServerError)
		return
	}

	c.LogAudit("action_item_id=" + completed.Id)
	if err := json.NewEncoder(w).Encode(completed); err != nil {
		c.Logger.Warn("Error writing response", mlog.Err(err))
	}
}

func deleteActionItem(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireActionItemId()
	if c.Err != nil {
		return
	}

	if err := c.App.DeleteActionItem(c.AppContext, c.Params.ActionItemId, c.AppContext.Session().UserId); err != nil {
		c.Err = model.NewAppError("deleteActionItem", "api.action_item.delete.app_error", nil, err.Error(), http.StatusInternalServerError)
		return
	}

	c.LogAudit("action_item_id=" + c.Params.ActionItemId)
	ReturnStatusOK(w)
}

func getActionItemStats(c *Context, w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		userID = c.AppContext.Session().UserId
	}

	stats, err := c.App.GetActionItemStats(c.AppContext, userID)
	if err != nil {
		c.Err = model.NewAppError("getActionItemStats", "api.action_item.stats.app_error", nil, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(stats); err != nil {
		c.Logger.Warn("Error writing response", mlog.Err(err))
	}
}

