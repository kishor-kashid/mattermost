package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/public/pluginapi"

	"github.com/mattermost/mattermost/server/plugins/ai-suite/server/openai"
	"github.com/mattermost/mattermost/server/plugins/ai-suite/server/store"
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
}

// OnActivate is invoked when the plugin is activated on a Mattermost server.
func (p *Plugin) OnActivate() error {
	p.apiClient = pluginapi.NewClient(p.API, p.Driver)
	p.router = NewAPIRouter(p)

	if err := p.OnConfigurationChange(); err != nil {
		return err
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
