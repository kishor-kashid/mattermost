package summarizer

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/pluginapi"

	"github.com/mattermost/mattermost/server/plugins/ai-suite/server/openai"
	"github.com/mattermost/mattermost/server/plugins/ai-suite/server/store"
)

const (
	defaultChannelRange = 24 * time.Hour
	maxChannelRange     = 30 * 24 * time.Hour
)

// Config exposes summarizer settings.
type Config struct {
	EnableSummarization bool
	MaxSummaryMessages  int
}

// Service orchestrates Mattermost message retrieval and OpenAI summarization.
type Service struct {
	api     *pluginapi.Client
	store   *store.Service
	openai  *openai.Client
	cfgFn   func() Config
	cache   *Cache
	timeNow func() time.Time
}

// NewService builds a summarizer instance.
func NewService(api *pluginapi.Client, store *store.Service, client *openai.Client, cfgFn func() Config) *Service {
	if api == nil || store == nil || client == nil || cfgFn == nil {
		return nil
	}

	return &Service{
		api:     api,
		store:   store,
		openai:  client,
		cfgFn:   cfgFn,
		cache:   NewCache(store, 24*time.Hour),
		timeNow: time.Now,
	}
}

// Summarize generates (or returns cached) summaries for the requested context.
func (s *Service) Summarize(ctx context.Context, req Request) (*Response, error) {
	if s == nil || s.openai == nil {
		return nil, ErrClientUnavailable
	}

	cfg := s.cfgFn()
	if !cfg.EnableSummarization {
		return nil, ErrDisabled
	}

	conv, err := s.buildConversation(ctx, cfg, req)
	if err != nil {
		return nil, err
	}

	cacheKey := s.buildCacheKey(req, conv)
	if !req.Force {
		if cached, found, err := s.cache.Get(ctx, cacheKey); err == nil && found {
			cached.ID = cacheKey
			cached.ConversationHash = conv.hash
			return cached, nil
		}
	}

	summary, err := s.generateSummary(ctx, req, conv)
	if err != nil {
		return nil, err
	}

	summary.ID = cacheKey
	summary.ConversationHash = conv.hash
	_ = s.cache.Set(ctx, cacheKey, summary)
	return summary, nil
}

// Config returns the current summarizer config (exported for handlers/tests).
func (s *Service) Config() Config {
	if s == nil || s.cfgFn == nil {
		return Config{}
	}
	return s.cfgFn()
}

func (s *Service) maxMessages(cfg Config) int {
	if cfg.MaxSummaryMessages > 0 {
		return cfg.MaxSummaryMessages
	}
	return 500
}

func (s *Service) buildConversation(ctx context.Context, cfg Config, req Request) (*conversation, error) {
	switch req.Type {
	case TypeChannel:
		return s.buildChannelConversation(ctx, cfg, req)
	case TypeThread, "":
		return s.buildThreadConversation(ctx, cfg, req)
	default:
		return nil, fmt.Errorf("%w: unknown type %s", ErrInvalidRequest, req.Type)
	}
}

func (s *Service) generateSummary(ctx context.Context, req Request, conv *conversation) (*Response, error) {
	if conv == nil {
		return nil, ErrEmptyConversation
	}

	promptVars := map[string]string{
		"channel":      conv.channelDisplay,
		"timeframe":    conv.rangeLabel,
		"conversation": conv.formatted,
	}

	messages := openai.SummarizeConversationTemplate.Render(promptVars)
	resp, err := s.openai.ChatCompletion(ctx, openai.ChatCompletionRequest{
		Messages:    messages,
		User:        req.UserID,
		Temperature: 0.2,
	})
	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 || strings.TrimSpace(resp.Choices[0].Message.Content) == "" {
		return nil, errors.New("openai: empty response")
	}

	content := strings.TrimSpace(resp.Choices[0].Message.Content)
	now := s.timeNow().UTC()

	summary := &Response{
		Type:             req.Type,
		ChannelID:        conv.channelID,
		ChannelName:      conv.channelDisplay,
		RootPostID:       conv.rootPostID,
		Title:            conv.title,
		Summary:          content,
		MessageCount:     conv.messageCount,
		ParticipantCount: len(conv.participants),
		Participants:     conv.participants,
		GeneratedAt:      now.UnixMilli(),
		Range: SummaryRange{
			Since: conv.since.UnixMilli(),
			Until: conv.until.UnixMilli(),
			Label: conv.rangeLabel,
		},
		Context: SummaryContext{
			TypeLabel:    conv.contextLabel,
			MessageLimit: conv.messageLimit,
			Timeframe:    conv.rangeLabel,
		},
		LimitReached: conv.limitReached,
		Cached:       false,
	}

	if resp.Usage != nil {
		summary.Usage = &Usage{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		}
	}

	return summary, nil
}

func (s *Service) buildCacheKey(req Request, conv *conversation) string {
	rangePart := fmt.Sprintf("%d-%d", conv.since.UnixMilli(), conv.until.UnixMilli())
	root := conv.rootPostID
	if root == "" {
		root = req.PostID
	}

	return BuildCacheKey(
		string(req.Type),
		conv.channelID,
		root,
		rangePart,
		conv.hash,
	)
}

func (s *Service) ensureMember(channelID, userID string) error {
	if channelID == "" || userID == "" {
		return fmt.Errorf("%w: missing channel or user", ErrInvalidRequest)
	}

	if _, err := s.api.Channel.GetMember(channelID, userID); err != nil {
		return ErrUnauthorized
	}

	return nil
}

type conversation struct {
	channelID      string
	channelDisplay string
	rootPostID     string
	title          string
	contextLabel   string
	rangeLabel     string
	since          time.Time
	until          time.Time
	messageCount   int
	messageLimit   int
	limitReached   bool
	participants   []Participant
	formatted      string
	hash           string
}

func (s *Service) buildThreadConversation(ctx context.Context, cfg Config, req Request) (*conversation, error) {
	rootID, channelID, err := s.resolveRoot(req)
	if err != nil {
		return nil, err
	}

	if err := s.ensureMember(channelID, req.UserID); err != nil {
		return nil, err
	}

	thread, err := s.api.Post.GetPostThread(rootID)
	if err != nil {
		return nil, err
	}

	posts := filterPosts(thread.ToSlice())
	if len(posts) == 0 {
		return nil, ErrEmptyConversation
	}

	sortPosts(posts)
	limit := s.maxMessages(cfg)
	limitReached, trimmed := applyLimit(posts, limit, false)
	participants := s.collectParticipants(posts)

	channel, err := s.api.Channel.Get(channelID)
	if err != nil {
		return nil, err
	}

	formatted := formatConversation(trimmed, participants)
	since := time.UnixMilli(trimmed[0].CreateAt)
	until := time.UnixMilli(trimmed[len(trimmed)-1].CreateAt)
	rangeLabel := fmt.Sprintf("%s – %s", since.Format(time.RFC822), until.Format(time.RFC822))

	return &conversation{
		channelID:      channelID,
		channelDisplay: displayName(channel),
		rootPostID:     rootID,
		title:          "Thread Summary",
		contextLabel:   "Thread",
		rangeLabel:     rangeLabel,
		since:          since,
		until:          until,
		messageCount:   len(trimmed),
		messageLimit:   limit,
		limitReached:   limitReached,
		participants:   participants,
		formatted:      formatted,
		hash:           hashPosts(trimmed),
	}, nil
}

func (s *Service) buildChannelConversation(ctx context.Context, cfg Config, req Request) (*conversation, error) {
	if req.ChannelID == "" {
		return nil, fmt.Errorf("%w: channel_id required", ErrInvalidRequest)
	}

	if err := s.ensureMember(req.ChannelID, req.UserID); err != nil {
		return nil, err
	}

	channel, err := s.api.Channel.Get(req.ChannelID)
	if err != nil {
		return nil, err
	}

	since, until, label, err := s.resolveRange(req)
	if err != nil {
		return nil, err
	}

	postList, err := s.api.Post.GetPostsSince(req.ChannelID, since.UnixMilli())
	if err != nil {
		return nil, err
	}

	posts := filterPosts(postList.ToSlice())
	sortPosts(posts)

	filtered := posts[:0]
	for _, post := range posts {
		ts := time.UnixMilli(post.CreateAt)
		if ts.Before(since) || ts.After(until) {
			continue
		}
		filtered = append(filtered, post)
	}
	if len(filtered) == 0 {
		return nil, ErrEmptyConversation
	}

	limit := s.maxMessages(cfg)
	limitReached, trimmed := applyLimit(filtered, limit, true)
	participants := s.collectParticipants(trimmed)
	formatted := formatConversation(trimmed, participants)

	return &conversation{
		channelID:      req.ChannelID,
		channelDisplay: displayName(channel),
		title:          fmt.Sprintf("Channel Summary • #%s", channel.Name),
		contextLabel:   "Channel",
		rangeLabel:     label,
		since:          since,
		until:          until,
		messageCount:   len(trimmed),
		messageLimit:   limit,
		limitReached:   limitReached,
		participants:   participants,
		formatted:      formatted,
		hash:           hashPosts(trimmed),
	}, nil
}

func (s *Service) resolveRoot(req Request) (string, string, error) {
	if req.RootPostID != "" {
		post, err := s.api.Post.GetPost(req.RootPostID)
		if err != nil {
			return "", "", err
		}
		return req.RootPostID, post.ChannelId, nil
	}

	if req.PostID == "" {
		return "", "", fmt.Errorf("%w: root_post_id or post_id required", ErrInvalidRequest)
	}

	post, err := s.api.Post.GetPost(req.PostID)
	if err != nil {
		return "", "", err
	}
	root := post.Id
	if post.RootId != "" {
		root = post.RootId
	}
	return root, post.ChannelId, nil
}

func (s *Service) resolveRange(req Request) (time.Time, time.Time, string, error) {
	now := s.timeNow()
	until := now
	if req.UntilMillis > 0 {
		until = time.UnixMilli(req.UntilMillis)
	}

	var since time.Time
	switch strings.ToLower(strings.TrimSpace(req.TimeRange)) {
	case "", "24h":
		since = until.Add(-defaultChannelRange)
	case "3d":
		since = until.Add(-72 * time.Hour)
	case "7d":
		since = until.Add(-7 * 24 * time.Hour)
	case "30d":
		since = until.Add(-30 * 24 * time.Hour)
	case "today":
		y, m, d := until.Date()
		since = time.Date(y, m, d, 0, 0, 0, 0, until.Location())
	default:
		if duration, err := time.ParseDuration(req.TimeRange); err == nil {
			if duration <= 0 || duration > maxChannelRange {
				return time.Time{}, time.Time{}, "", fmt.Errorf("%w: invalid time range", ErrInvalidRequest)
			}
			since = until.Add(-duration)
		} else if parts := strings.Fields(req.TimeRange); len(parts) == 2 {
			start, errStart := time.Parse("2006-01-02", parts[0])
			end, errEnd := time.Parse("2006-01-02", parts[1])
			if errStart != nil || errEnd != nil {
				return time.Time{}, time.Time{}, "", fmt.Errorf("%w: invalid date range", ErrInvalidRequest)
			}
			since = start
			until = end.Add(24*time.Hour - time.Nanosecond)
		} else if req.SinceMillis > 0 {
			since = time.UnixMilli(req.SinceMillis)
		} else {
			return time.Time{}, time.Time{}, "", fmt.Errorf("%w: unsupported range %q", ErrInvalidRequest, req.TimeRange)
		}
	}

	if req.SinceMillis > 0 {
		since = time.UnixMilli(req.SinceMillis)
	}

	if since.After(until) {
		return time.Time{}, time.Time{}, "", fmt.Errorf("%w: since after until", ErrInvalidRequest)
	}

	if until.Sub(since) > maxChannelRange {
		return time.Time{}, time.Time{}, "", fmt.Errorf("%w: range exceeds 30 days", ErrInvalidRequest)
	}

	label := fmt.Sprintf("%s – %s", since.Format("Jan 2, 15:04"), until.Format("Jan 2, 15:04"))
	return since, until, label, nil
}

func (s *Service) collectParticipants(posts []*model.Post) []Participant {
	unique := make(map[string]struct{})
	for _, post := range posts {
		if post == nil || post.UserId == "" {
			continue
		}
		unique[post.UserId] = struct{}{}
	}

	ids := make([]string, 0, len(unique))
	for id := range unique {
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return nil
	}

	userMap := make(map[string]*model.User, len(ids))
	users, err := s.api.User.ListByUserIDs(ids)
	if err == nil {
		for _, user := range users {
			userMap[user.Id] = user
		}
	} else {
		for _, id := range ids {
			user, getErr := s.api.User.Get(id)
			if getErr == nil {
				userMap[id] = user
			}
		}
	}

	out := make([]Participant, 0, len(ids))
	for _, id := range ids {
		user := userMap[id]
		name := id
		username := ""
		if user != nil {
			if display := user.GetDisplayName(model.ShowUsername); display != "" {
				name = display
			}
			username = user.Username
		}
		out = append(out, Participant{
			ID:          id,
			Username:    username,
			DisplayName: name,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return strings.ToLower(out[i].DisplayName) < strings.ToLower(out[j].DisplayName)
	})

	return out
}

// Helper utilities

func filterPosts(posts []*model.Post) []*model.Post {
	out := make([]*model.Post, 0, len(posts))
	for _, post := range posts {
		if post == nil || post.DeleteAt != 0 || post.IsSystemMessage() {
			continue
		}
		out = append(out, post)
	}
	return out
}

func sortPosts(posts []*model.Post) {
	sort.SliceStable(posts, func(i, j int) bool {
		return posts[i].CreateAt < posts[j].CreateAt
	})
}

func applyLimit(posts []*model.Post, limit int, keepNewest bool) (bool, []*model.Post) {
	if limit <= 0 || len(posts) <= limit {
		return false, posts
	}

	if keepNewest {
		return true, posts[len(posts)-limit:]
	}
	return true, posts[:limit]
}

func formatConversation(posts []*model.Post, participants []Participant) string {
	names := make(map[string]string, len(participants))
	for _, p := range participants {
		name := p.DisplayName
		if name == "" {
			name = p.Username
		}
		if name == "" {
			name = p.ID
		}
		names[p.ID] = name
	}

	var b strings.Builder
	for _, post := range posts {
		when := time.UnixMilli(post.CreateAt).UTC().Format("2006-01-02 15:04")
		author := names[post.UserId]
		if author == "" {
			author = "Unknown"
		}
		message := sanitize(post.Message)
		if message == "" && len(post.FileIds) > 0 {
			message = fmt.Sprintf("[Attached %d file(s)]", len(post.FileIds))
		}
		if message == "" {
			continue
		}
		b.WriteString(fmt.Sprintf("[%s] %s: %s\n", when, author, message))
	}
	return strings.TrimSpace(b.String())
}

func sanitize(input string) string {
	clean := strings.ReplaceAll(input, "\r\n", "\n")
	lines := strings.Split(clean, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	return strings.TrimSpace(strings.Join(lines, " "))
}

func displayName(channel *model.Channel) string {
	if channel == nil {
		return "channel"
	}
	if channel.DisplayName != "" {
		return channel.DisplayName
	}
	return channel.Name
}

func hashPosts(posts []*model.Post) string {
	var b strings.Builder
	for _, post := range posts {
		if post == nil {
			continue
		}
		b.WriteString(post.Id)
		b.WriteRune(':')
		b.WriteString(fmt.Sprint(post.UpdateAt))
		b.WriteRune(';')
	}
	return BuildCacheKey(b.String())
}
