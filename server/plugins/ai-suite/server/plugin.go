package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/public/pluginapi"

	"github.com/mattermost/mattermost/server/plugins/ai-suite/server/openai"
	"github.com/mattermost/mattermost/server/plugins/ai-suite/server/store"
	"github.com/mattermost/mattermost/server/plugins/ai-suite/server/summarizer"
)

// Plugin implements the Mattermost plugin interface.
type Plugin struct {
	plugin.MattermostPlugin

	configurationLock sync.RWMutex
	configuration     *configuration

	apiClient    *pluginapi.Client
	router       *APIRouter
	openAIClient *openai.Client
	store        *store.Service
	summarizer   *summarizer.Service
}

// OnActivate is invoked when the plugin is activated on a Mattermost server.
func (p *Plugin) OnActivate() error {
	p.apiClient = pluginapi.NewClient(p.API, p.Driver)
	p.router = NewAPIRouter(p)

	if err := p.OnConfigurationChange(); err != nil {
		return err
	}

	if err := p.registerCommands(); err != nil {
		p.apiClient.Log.Warn("failed to register slash command", "error", err.Error())
	}

	p.apiClient.Log.Info("AI Productivity Suite plugin activated", "version", PluginVersion)
	return nil
}

// OnDeactivate is invoked just before the plugin is removed or disabled.
func (p *Plugin) OnDeactivate() error {
	if p.apiClient != nil {
		p.apiClient.Log.Info("AI Productivity Suite plugin deactivated")
	}

	p.openAIClient = nil
	p.store = nil
	p.router = nil
	p.summarizer = nil

	return nil
}

// OnConfigurationChange is triggered whenever the admin updates the plugin configuration.
func (p *Plugin) OnConfigurationChange() error {
	cfg := newConfiguration()
	if err := p.API.LoadPluginConfiguration(cfg); err != nil {
		return err
	}

	cfg.ApplyDefaults()
	if err := cfg.Validate(); err != nil {
		return err
	}

	p.setConfiguration(cfg)

	p.store = store.NewService(store.NewKVStore(p.API, "ai-suite"))

	if cfg.OpenAIAPIKey == "" {
		p.openAIClient = nil
		p.summarizer = nil
		p.apiClient.Log.Warn("OpenAI API key not configured; AI functionality disabled")
		return nil
	}

	client, err := openai.NewClient(openai.ClientConfig{
		APIKey:               cfg.OpenAIAPIKey,
		Model:                cfg.OpenAIModel,
		BaseURL:              cfg.OpenAIBaseURL,
		RequestTimeout:       time.Duration(cfg.RequestTimeoutSecs) * time.Second,
		MaxRequestsPerMinute: cfg.APIRateLimit,
		MaxRetries:           cfg.OpenAIRetryAttempts,
	})
	if err != nil {
		return err
	}

	p.openAIClient = client
	p.summarizer = summarizer.NewService(
		p.apiClient,
		p.store,
		p.openAIClient,
		func() summarizer.Config {
			conf := p.getConfiguration()
			return summarizer.Config{
				EnableSummarization: conf.EnableSummarization,
				MaxSummaryMessages:  conf.MaxSummaryMessages,
			}
		},
	)

	return nil
}

// ServeHTTP routes incoming plugin HTTP requests.
func (p *Plugin) ServeHTTP(_ *plugin.Context, w http.ResponseWriter, r *http.Request) {
	if p.router == nil {
		http.NotFound(w, r)
		return
	}

	p.router.ServeHTTP(w, r)
}

func (p *Plugin) getConfiguration() *configuration {
	p.configurationLock.RLock()
	defer p.configurationLock.RUnlock()

	if p.configuration == nil {
		return newConfiguration()
	}

	return p.configuration.Clone()
}

func (p *Plugin) setConfiguration(cfg *configuration) {
	p.configurationLock.Lock()
	defer p.configurationLock.Unlock()
	p.configuration = cfg.Clone()
}

func main() {
	plugin.ClientMain(&Plugin{})
}

func (p *Plugin) registerCommands() error {
	if p.API == nil {
		return fmt.Errorf("plugin api unavailable")
	}

	cmd := &model.Command{
		Trigger:          summarizer.CommandTrigger,
		DisplayName:      "AI Summarize",
		Description:      "Generate AI-powered conversation summaries",
		AutoComplete:     true,
		AutoCompleteDesc: "Summarize the current thread or channel",
		AutoCompleteHint: "[thread|channel] [time-range]",
	}

	return p.API.RegisterCommand(cmd)
}

// ExecuteCommand handles /summarize.
func (p *Plugin) ExecuteCommand(_ *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	if args == nil || strings.TrimSpace(args.Command) == "" {
		return nil, model.NewAppError("ExecuteCommand", "app.command.execute.no_command.app_error", nil, "", http.StatusBadRequest)
	}

	if p.summarizer == nil {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "Summarization is not configured yet. Please ask a system admin to add the OpenAI API key.",
		}, nil
	}

	opts, err := summarizer.ParseCommand(args.Command)
	if err != nil {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         err.Error(),
		}, nil
	}

	req := summarizer.Request{
		Type:       opts.Target,
		UserID:     args.UserId,
		ChannelID:  args.ChannelId,
		RootPostID: args.RootId,
		TimeRange:  opts.Argument,
	}

	if req.Type == summarizer.TypeThread && req.RootPostID == "" {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "Run `/summarize thread` from inside the thread you want summarized.",
		}, nil
	}

	if req.Type == summarizer.TypeChannel && req.ChannelID == "" {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "Unable to determine the target channel.",
		}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	result, err := p.summarizer.Summarize(ctx, req)
	if err != nil {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         fmt.Sprintf("Unable to generate summary: %v", err),
		}, nil
	}

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("📋 **Summary** (%d messages, %d participants)\n\n", result.MessageCount, result.ParticipantCount))
	builder.WriteString(result.Summary)

	return &model.CommandResponse{
		ResponseType:     model.CommandResponseTypeEphemeral,
		Text:             builder.String(),
		SkipSlackParsing: true,
	}, nil
}
