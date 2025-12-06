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
	"github.com/mattermost/mattermost/server/v8/channels/app/openai"
)

type FormatProvider struct{}

const (
	CommandTriggerFormat = "format"
)

func init() {
	app.RegisterCommandProvider(&FormatProvider{})
}

func (fp *FormatProvider) GetTrigger() string {
	return CommandTriggerFormat
}

func (fp *FormatProvider) GetCommand(a *app.App, T i18n.TranslateFunc) *model.Command {
	return &model.Command{
		Trigger:          CommandTriggerFormat,
		AutoComplete:     true,
		AutoCompleteDesc: T("api.command_format.desc"),
		AutoCompleteHint: T("api.command_format.hint"),
		DisplayName:      T("api.command_format.name"),
		Description:      T("api.command_format.description"),
	}
}

func (fp *FormatProvider) DoCommand(a *app.App, rctx request.CTX, args *model.CommandArgs, message string) *model.CommandResponse {
	// Check if AI formatting is enabled
	if !a.IsAIFeatureEnabled("formatting") {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "AI formatting is not enabled on this server.",
		}
	}

	// Parse command arguments
	parts := strings.Fields(message)
	
	var profile openai.FormattingProfile = openai.FormattingProfessional
	var textToFormat string
	
	// Parse profile if specified
	for i, part := range parts {
		switch strings.ToLower(part) {
		case "professional", "casual", "technical", "concise":
			profile = openai.FormattingProfile(strings.ToLower(part))
			// Everything after the profile is the text to format
			if i+1 < len(parts) {
				textToFormat = strings.Join(parts[i+1:], " ")
			}
			break
		}
	}
	
	// If no profile found, use the entire message as text
	if textToFormat == "" {
		textToFormat = message
	}
	
	if textToFormat == "" {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "Please provide text to format. Usage: `/format [professional|casual|technical|concise] <your message>`",
		}
	}

	// Format the message
	formattingReq := &app.FormattingRequest{
		Message: textToFormat,
		Profile:  profile,
	}

	response, err := a.FormatMessage(rctx, formattingReq)
	if err != nil {
		rctx.Logger().Error("Failed to format message", mlog.Err(err))
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         fmt.Sprintf("Failed to format message: %s", err.Message),
		}
	}

	// Return formatted text with instructions
	profileName := strings.Title(string(profile))
	responseText := fmt.Sprintf("**Formatted (%s):**\n\n%s\n\n*Processing time: %dms*", 
		profileName, response.FormattedText, response.ProcessingMs)

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         responseText,
	}
}

