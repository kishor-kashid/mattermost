// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package slashcommands

import (
	"fmt"
	"strings"
	"time"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/shared/i18n"
	"github.com/mattermost/mattermost/server/public/shared/request"
	"github.com/mattermost/mattermost/server/v8/channels/app"
)

type AIActionItemsProvider struct {
}

const (
	CmdAIActionItems = "actionitems"
)

func init() {
	app.RegisterCommandProvider(&AIActionItemsProvider{})
}

func (*AIActionItemsProvider) GetTrigger() string {
	return CmdAIActionItems
}

func (*AIActionItemsProvider) GetCommand(a *app.App, T i18n.TranslateFunc) *model.Command {
	return &model.Command{
		Trigger:          CmdAIActionItems,
		AutoComplete:     true,
		AutoCompleteDesc: "Manage AI-detected action items",
		AutoCompleteHint: "[list|mine|team|stats|complete <id>]",
		DisplayName:      "Action Items",
		Description:      "AI-powered action item management",
	}
}

func (*AIActionItemsProvider) DoCommand(a *app.App, c request.CTX, args *model.CommandArgs, message string) *model.CommandResponse {
	if !a.IsAIFeatureEnabled("action_items") {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "AI action items feature is not enabled on this server.",
		}
	}

	parts := strings.Fields(message)
	command := "mine" // default command
	if len(parts) > 0 {
		command = strings.ToLower(parts[0])
	}

	switch command {
	case "list", "mine":
		return handleListActionItems(a, c, args, "mine")
	case "team", "channel":
		return handleListActionItems(a, c, args, "channel")
	case "stats":
		return handleActionItemStats(a, c, args)
	case "complete":
		if len(parts) < 2 {
			return getHelp()
		}
		return handleCompleteActionItem(a, c, args, parts[1])
	case "help":
		return getHelp()
	default:
		return getHelp()
	}
}

func handleListActionItems(a *app.App, c request.CTX, args *model.CommandArgs, scope string) *model.CommandResponse {
	filters := &app.ActionItemFilters{
		IncludeCompleted: false,
		Page:             0,
		PerPage:          20,
	}

	var items []*model.AIActionItem
	var err error
	var title string

	if scope == "mine" {
		items, err = a.GetActionItemsForUser(c, args.UserId, filters)
		title = "Your Action Items"
	} else {
		items, err = a.GetActionItemsForChannel(c, args.ChannelId, args.UserId, filters)
		title = "Channel Action Items"
	}

	if err != nil {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         fmt.Sprintf("Error retrieving action items: %s", err.Error()),
		}
	}

	if len(items) == 0 {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "No action items found.",
		}
	}

	// Group items by status
	overdue := []*model.AIActionItem{}
	dueSoon := []*model.AIActionItem{}
	noDueDate := []*model.AIActionItem{}

	now := time.Now()
	weekFromNow := now.Add(7 * 24 * time.Hour)

	for _, item := range items {
		if item.DueDate > 0 {
			dueTime := time.UnixMilli(item.DueDate)
			if dueTime.Before(now) {
				overdue = append(overdue, item)
			} else if dueTime.Before(weekFromNow) {
				dueSoon = append(dueSoon, item)
			} else {
				noDueDate = append(noDueDate, item)
			}
		} else {
			noDueDate = append(noDueDate, item)
		}
	}

	// Build response message
	var message strings.Builder
	message.WriteString(fmt.Sprintf("### %s\n\n", title))

	if len(overdue) > 0 {
		message.WriteString("#### 🔴 Overdue\n")
		for _, item := range overdue {
			message.WriteString(formatActionItemLine(item, now))
		}
		message.WriteString("\n")
	}

	if len(dueSoon) > 0 {
		message.WriteString("#### ⏰ Due Soon\n")
		for _, item := range dueSoon {
			message.WriteString(formatActionItemLine(item, now))
		}
		message.WriteString("\n")
	}

	if len(noDueDate) > 0 {
		message.WriteString("#### 📋 Other Items\n")
		for _, item := range noDueDate {
			message.WriteString(formatActionItemLine(item, now))
		}
	}

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         message.String(),
	}
}

func handleActionItemStats(a *app.App, c request.CTX, args *model.CommandArgs) *model.CommandResponse {
	stats, err := a.GetActionItemStats(c, args.UserId)
	if err != nil {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         fmt.Sprintf("Error retrieving stats: %s", err.Error()),
		}
	}

	var message strings.Builder
	message.WriteString("### Your Action Item Statistics\n\n")
	message.WriteString(fmt.Sprintf("**Total:** %d\n", stats.Total))
	message.WriteString(fmt.Sprintf("**Overdue:** %d\n", stats.Overdue))
	message.WriteString(fmt.Sprintf("**Due Today:** %d\n", stats.DueToday))
	message.WriteString(fmt.Sprintf("**Due This Week:** %d\n", stats.DueSoon))
	message.WriteString(fmt.Sprintf("**Completed:** %d\n", stats.Completed))
	message.WriteString(fmt.Sprintf("**No Due Date:** %d\n\n", stats.NoDueDate))

	message.WriteString("**By Priority:**\n")
	for priority, count := range stats.ByPriority {
		if count > 0 {
			message.WriteString(fmt.Sprintf("- %s: %d\n", capitalizeFirst(priority), count))
		}
	}

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         message.String(),
	}
}

func handleCompleteActionItem(a *app.App, c request.CTX, args *model.CommandArgs, itemID string) *model.CommandResponse {
	_, err := a.CompleteActionItem(c, itemID, args.UserId)
	if err != nil {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         fmt.Sprintf("Error completing action item: %s", err.Error()),
		}
	}

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         "✅ Action item marked as complete!",
	}
}

func formatActionItemLine(item *model.AIActionItem, now time.Time) string {
	priorityEmoji := getPriorityEmoji(item.Priority)
	
	dueStr := "No due date"
	if item.DueDate > 0 {
		dueTime := time.UnixMilli(item.DueDate)
		dueStr = dueTime.Format("Jan 2")
		
		// Add time until due
		duration := dueTime.Sub(now)
		if duration < 0 {
			duration = -duration
			dueStr += fmt.Sprintf(" (overdue by %s)", formatDuration(duration))
		} else {
			dueStr += fmt.Sprintf(" (in %s)", formatDuration(duration))
		}
	}

	postLink := ""
	if item.PostId != "" {
		postLink = fmt.Sprintf(" [→](/_redirect/pl/%s)", item.PostId)
	}

	return fmt.Sprintf("- %s **%s** - %s%s `%s`\n", 
		priorityEmoji,
		item.Description,
		dueStr,
		postLink,
		item.Id[:8],
	)
}

func getPriorityEmoji(priority string) string {
	switch priority {
	case "urgent":
		return "🔥"
	case "high":
		return "🔴"
	case "medium":
		return "🟡"
	case "low":
		return "🟢"
	default:
		return "⚪"
	}
}

func formatDuration(d time.Duration) string {
	if d < 0 {
		d = -d
	}

	hours := int(d.Hours())
	if hours < 24 {
		return fmt.Sprintf("%dh", hours)
	}

	days := hours / 24
	if days == 1 {
		return "1 day"
	}
	if days < 7 {
		return fmt.Sprintf("%d days", days)
	}

	weeks := days / 7
	if weeks == 1 {
		return "1 week"
	}
	return fmt.Sprintf("%d weeks", weeks)
}

func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	first := strings.ToUpper(string(s[0]))
	return first + s[1:]
}

func getHelp() *model.CommandResponse {
	help := `### Action Items Command Help

**Usage:** /actionitems [command]

**Available Commands:**
- **mine** or **list** - Show your action items
- **team** or **channel** - Show action items for the current channel
- **stats** - Show your action item statistics
- **complete <id>** - Mark an action item as complete (use the short ID shown in listings)
- **help** - Show this help message

**Examples:**
- /actionitems mine
- /actionitems channel
- /actionitems complete abc123de
`
	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         help,
	}
}

