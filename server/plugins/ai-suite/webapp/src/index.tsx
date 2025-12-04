import manifest from './manifest';
import SummaryPanel from './components/summarizer/SummaryPanel';
import {SummaryBridge} from './summarizer/bridge';
import type {ClientPlugin, PluginRegistry} from './types';

class AIProductivitySuitePlugin implements ClientPlugin {
    private bridge = new SummaryBridge();
    private showRHS?: () => void;
    private hideRHS?: () => void;

    public initialize(registry: PluginRegistry) {
        if (registry.registerRightHandSidebarComponent) {
            const rhs = registry.registerRightHandSidebarComponent(
                () => (
                    <SummaryPanel
                        bridge={this.bridge}
                        onClose={() => this.hideRHS?.()}
                    />
                ),
                'AI Summaries',
            );
            this.showRHS = rhs.showRHSPlugin;
            this.hideRHS = rhs.hideRHSPlugin;
        }

        registry.registerPostDropdownMenuAction?.(
            'Summarize Thread',
            (postId: string) => {
                this.bridge.setTarget({
                    type: 'thread',
                    postId,
                });
                this.showRHS?.();
            },
            () => true,
        );

        registry.registerChannelHeaderMenuAction?.(
            'Summarize Channel',
            (channelId: string) => {
                this.bridge.setTarget({
                    type: 'channel',
                    channelId,
                    timeRange: '24h',
                });
                this.showRHS?.();
            },
        );
    }

    public uninitialize() {
        this.bridge.setTarget(null);
        this.hideRHS?.();
    }
}

window.registerPlugin(manifest.id, new AIProductivitySuitePlugin());

