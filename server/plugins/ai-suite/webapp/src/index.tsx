import manifest from './manifest';
import type {ClientPlugin, PluginRegistry} from './types';

class AIProductivitySuitePlugin implements ClientPlugin {
    public initialize(registry: PluginRegistry) {
        // Placeholder registrations. Future PRs will hook real UI integrations.
        if (registry.registerChannelHeaderButtonAction) {
            registry.registerChannelHeaderButtonAction(
                () => alert('AI Productivity Suite placeholder'),
                'AI Suite',
                'AI actions',
            );
        }

        // eslint-disable-next-line no-console
        console.log('AI Productivity Suite webapp initialized', manifest.version);
    }

    public uninitialize() {
        // eslint-disable-next-line no-console
        console.log('AI Productivity Suite webapp unmounted');
    }
}

window.registerPlugin(manifest.id, new AIProductivitySuitePlugin());

