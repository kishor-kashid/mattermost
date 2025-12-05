// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package ai_action_item_reminders

import (
	"fmt"
	"time"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/shared/mlog"
	"github.com/mattermost/mattermost/server/public/shared/request"
	"github.com/mattermost/mattermost/server/v8/channels/jobs"
)

func MakeWorker(jobServer *jobs.JobServer, sendReminders func(c request.CTX) error) *jobs.SimpleWorker {
	const workerName = "AIActionItemReminders"

	isEnabled := func(cfg *model.Config) bool {
		if cfg.AISettings.Enable == nil {
			return false
		}
		if cfg.AISettings.EnableActionItems == nil {
			return false
		}
		return *cfg.AISettings.Enable && *cfg.AISettings.EnableActionItems
	}

	execute := func(logger mlog.LoggerIFace, job *model.Job) error {
		defer jobServer.HandleJobPanic(logger, job)

		c := request.EmptyContext(logger)
		
		logger.Info("Starting AI Action Item Reminders job")
		
		if err := sendReminders(c); err != nil {
			logger.Error("Failed to send action item reminders", mlog.Err(err))
			return err
		}
		
		logger.Info("Completed AI Action Item Reminders job")
		return nil
	}

	return jobs.NewSimpleWorker(workerName, jobServer, execute, isEnabled)
}

// SendActionItemReminders sends reminders for overdue and due-soon action items
func SendActionItemReminders(c request.CTX, app interface {
	GetOverdueActionItems(c request.CTX) ([]*model.AIActionItem, error)
	GetDueSoonActionItems(c request.CTX, hours int) ([]*model.AIActionItem, error)
	GetUser(userID string) (*model.User, *model.AppError)
	GetOrCreateDirectChannel(c request.CTX, userId1, userId2 string, opts ...model.ChannelOption) (*model.Channel, *model.AppError)
	CreatePost(c request.CTX, post *model.Post, channel *model.Channel, flags model.CreatePostFlags) (*model.Post, *model.AppError)
	GetChannel(c request.CTX, channelID string) (*model.Channel, *model.AppError)
}) error {
	c.Logger().Debug("Checking for action items requiring reminders")

	now := time.Now()

	// Get overdue items
	overdueItems, err := app.GetOverdueActionItems(c)
	if err != nil {
		return fmt.Errorf("failed to get overdue action items: %w", err)
	}

	// Get items due in next 24 hours
	dueSoonItems, err := app.GetDueSoonActionItems(c, 24)
	if err != nil {
		return fmt.Errorf("failed to get due soon action items: %w", err)
	}

	c.Logger().Debug("Found action items for reminders",
		mlog.Int("overdue", len(overdueItems)),
		mlog.Int("due_soon", len(dueSoonItems)),
	)

	// Send overdue reminders
	for _, item := range overdueItems {
		if err := sendReminderNotification(c, app, item, "overdue", now); err != nil {
			c.Logger().Error("Failed to send overdue reminder",
				mlog.String("action_item_id", item.Id),
				mlog.Err(err),
			)
		}
	}

	// Send due soon reminders
	for _, item := range dueSoonItems {
		if err := sendReminderNotification(c, app, item, "due_soon", now); err != nil {
			c.Logger().Error("Failed to send due soon reminder",
				mlog.String("action_item_id", item.Id),
				mlog.Err(err),
			)
		}
	}

	return nil
}

func sendReminderNotification(
	c request.CTX,
	app interface {
		GetUser(userID string) (*model.User, *model.AppError)
		GetOrCreateDirectChannel(c request.CTX, userId1, userId2 string, opts ...model.ChannelOption) (*model.Channel, *model.AppError)
		CreatePost(c request.CTX, post *model.Post, channel *model.Channel, flags model.CreatePostFlags) (*model.Post, *model.AppError)
		GetChannel(c request.CTX, channelID string) (*model.Channel, *model.AppError)
	},
	item *model.AIActionItem,
	reminderType string,
	now time.Time,
) error {
	// Don't send reminders for completed or dismissed items
	if item.Status == "completed" || item.Status == "dismissed" {
		return nil
	}

	// Get assignee
	assignee, err := app.GetUser(item.AssigneeId)
	if err != nil {
		return fmt.Errorf("failed to get assignee: %w", err)
	}

	// Get source channel
	channel, err := app.GetChannel(c, item.ChannelId)
	if err != nil {
		return fmt.Errorf("failed to get channel: %w", err)
	}

	// Build reminder message
	var message string
	var emoji string

	if reminderType == "overdue" {
		dueTime := time.UnixMilli(item.DueDate)
		overdueDuration := now.Sub(dueTime)
		
		emoji = "🔴"
		message = fmt.Sprintf(`#### %s Overdue Action Item

> %s

**Due:** %s (%s ago)  
**Priority:** %s  
**Channel:** #%s

[View in context](/_redirect/pl/%s)`,
			emoji,
			item.Description,
			dueTime.Format("Mon Jan 2, 2006 at 3:04 PM"),
			formatDuration(overdueDuration),
			capitalizeFirst(item.Priority),
			channel.DisplayName,
			item.PostId,
		)
	} else {
		dueTime := time.UnixMilli(item.DueDate)
		timeUntilDue := dueTime.Sub(now)
		
		emoji = "⏰"
		message = fmt.Sprintf(`#### %s Action Item Due Soon

> %s

**Due:** %s (in %s)  
**Priority:** %s  
**Channel:** #%s

[View in context](/_redirect/pl/%s)`,
			emoji,
			item.Description,
			dueTime.Format("Mon Jan 2, 2006 at 3:04 PM"),
			formatDuration(timeUntilDue),
			capitalizeFirst(item.Priority),
			channel.DisplayName,
			item.PostId,
		)
	}

	// Get or create DM channel with assignee
	botUser, err := app.GetUser("system") // TODO: Use a proper bot user
	if err != nil {
		// Fallback: send from the creator
		botUser, err = app.GetUser(item.CreatedBy)
		if err != nil {
			return fmt.Errorf("failed to get bot user: %w", err)
		}
	}

	dmChannel, err := app.GetOrCreateDirectChannel(c, botUser.Id, assignee.Id)
	if err != nil {
		return fmt.Errorf("failed to get DM channel: %w", err)
	}

	// Create reminder post
	reminderPost := &model.Post{
		UserId:    botUser.Id,
		ChannelId: dmChannel.Id,
		Message:   message,
	}

	_, appErr := app.CreatePost(c, reminderPost, dmChannel, model.CreatePostFlags{})
	if appErr != nil {
		return fmt.Errorf("failed to create reminder post: %w", appErr)
	}

	c.Logger().Info("Sent action item reminder",
		mlog.String("action_item_id", item.Id),
		mlog.String("assignee_id", assignee.Id),
		mlog.String("reminder_type", reminderType),
	)

	return nil
}

func formatDuration(d time.Duration) string {
	if d < 0 {
		d = -d
	}

	hours := int(d.Hours())
	if hours < 24 {
		return fmt.Sprintf("%d hours", hours)
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
	return string(s[0]-32) + s[1:]
}

