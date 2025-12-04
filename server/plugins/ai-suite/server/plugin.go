package main

import (
	"sync"

	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/public/pluginapi"
)

// Plugin implements the Mattermost plugin interface.
type Plugin struct {
	plugin.MattermostPlugin

	configurationLock sync.RWMutex
	configuration     *configuration

	apiClient *pluginapi.Client
}

// OnActivate is invoked when the plugin is activated on a Mattermost server.
func (p *Plugin) OnActivate() error {
	p.apiClient = pluginapi.NewClient(p.API, p.Driver)
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
	return nil
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
