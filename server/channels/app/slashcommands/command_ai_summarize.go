// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package slashcommands

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/shared/i18n"
	"github.com/mattermost/mattermost/server/public/shared/mlog"
	"github.com/mattermost/mattermost/server/public/shared/request"
	"github.com/mattermost/mattermost/server/v8/channels/app"
)

type SummarizeProvider struct {
}

const (
	CommandTriggerSummarize = "summarize"
)

func init() {
	app.RegisterCommandProvider(&SummarizeProvider{})
}

func (sp *SummarizeProvider) GetTrigger() string {
	return CommandTriggerSummarize
}

func (sp *SummarizeProvider) GetCommand(a *app.App, T i18n.TranslateFunc) *model.Command {
	return &model.Command{
		Trigger:          CommandTriggerSummarize,
		AutoComplete:     true,
		AutoCompleteDesc: T("api.command_summarize.desc"),
		AutoCompleteHint: T("api.command_summarize.hint"),
		DisplayName:      T("api.command_summarize.name"),
		Description:      T("api.command_summarize.description"),
	}
}

func (sp *SummarizeProvider) DoCommand(a *app.App, rctx request.CTX, args *model.CommandArgs, message string) *model.CommandResponse {
	// Check if AI summarization is enabled
	mlog.Debug("Slash command /summarize invoked",
		mlog.String("user_id", args.UserId),
		mlog.String("channel_id", args.ChannelId),
		mlog.String("message", message))
	
	isEnabled := a.IsAIFeatureEnabled("summarization")
	mlog.Debug("AI summarization feature check result",
		mlog.Bool("is_enabled", isEnabled))
	
	if !isEnabled {
		mlog.Warn("AI summarization command rejected - feature not enabled",
			mlog.String("user_id", args.UserId))
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "AI summarization is not enabled on this server.",
		}
	}

	// Parse command arguments
	parts := strings.Fields(message)
	
	var summaryLevel string
	var postId string
	useThread := false
	
	// Default to standard level
	summaryLevel = "standard"

	// Parse arguments
	for i, part := range parts {
		switch part {
		case "thread":
			useThread = true
			// Next argument might be the post ID or we use current thread
			if i+1 < len(parts) && !strings.HasPrefix(parts[i+1], "-") {
				postId = parts[i+1]
			}
		case "brief", "standard", "detailed":
			summaryLevel = part
		}
	}

	// If thread mode but no post ID specified, check if we're in a thread
	if useThread && postId == "" {
		// Get the parent post if we're in a thread
		if args.ParentId != "" {
			postId = args.ParentId
		} else if args.RootId != "" {
			postId = args.RootId
		} else {
			return &model.CommandResponse{
				ResponseType: model.CommandResponseTypeEphemeral,
				Text:         "Please specify a thread to summarize or use this command in a thread.",
			}
		}
	}

	// If still no post ID, we're doing a channel summary
	if postId == "" && !useThread {
		return sp.summarizeChannel(a, rctx, args, summaryLevel)
	}

	// Otherwise, summarize the thread
	return sp.summarizeThread(a, rctx, args, postId, summaryLevel)
}

func (sp *SummarizeProvider) summarizeThread(a *app.App, rctx request.CTX, args *model.CommandArgs, postId, level string) *model.CommandResponse {
	// Build summarization request
	req := &app.SummarizationRequest{
		ChannelId:    args.ChannelId,
		PostId:       postId,
		SummaryLevel: level,
		MaxMessages:  a.GetAIMaxMessageLimit(),
		UserId:       args.UserId,
		UseCache:     true,
	}

	// Execute summarization
	result, err := a.SummarizeThread(rctx, req)
	if err != nil {
		mlog.Error("Failed to summarize thread via slash command", mlog.Err(err))
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         fmt.Sprintf("Failed to generate summary: %v", err),
		}
	}

	// Format the response
	responseText := formatSummaryResponse(result.Summary, result.FromCache)

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         responseText,
	}
}

func (sp *SummarizeProvider) summarizeChannel(a *app.App, rctx request.CTX, args *model.CommandArgs, level string) *model.CommandResponse {
	// Default to last 24 hours
	endTime := model.GetMillis()
	startTime := endTime - (24 * 60 * 60 * 1000)

	// Build summarization request
	req := &app.SummarizationRequest{
		ChannelId:    args.ChannelId,
		StartTime:    startTime,
		EndTime:      endTime,
		SummaryLevel: level,
		MaxMessages:  a.GetAIMaxMessageLimit(),
		UserId:       args.UserId,
		UseCache:     true,
	}

	// Execute summarization
	result, err := a.SummarizeChannel(rctx, req)
	if err != nil {
		mlog.Error("Failed to summarize channel via slash command", mlog.Err(err))
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         fmt.Sprintf("Failed to generate summary: %v", err),
		}
	}

	// Format the response
	responseText := formatSummaryResponse(result.Summary, result.FromCache)

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         responseText,
	}
}

func formatSummaryResponse(summary *model.AISummary, fromCache bool) string {
	var builder strings.Builder

	builder.WriteString("### AI Summary\n\n")
	
	if fromCache {
		builder.WriteString("*From cache*\n\n")
	}

	builder.WriteString(summary.Summary)
	builder.WriteString("\n\n---\n\n")
	
	builder.WriteString(fmt.Sprintf("**Messages**: %d | **Participants**: %s\n",
		summary.MessageCount, summary.Participants))

	return builder.String()
}

