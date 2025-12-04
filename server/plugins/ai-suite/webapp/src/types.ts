export type PluginManifest = {
    id: string;
    name: string;
    description: string;
    version: string;
    min_server_version: string;
};

export type PluginRegistry = {
    registerChannelHeaderButtonAction?: (...args: any[]) => void;
    registerPostDropdownMenuAction?: (...args: any[]) => void;
    registerRightHandSidebarComponent?: (...args: any[]) => void;
};

export interface ClientPlugin {
    initialize(registry: PluginRegistry): void;
    uninitialize?(): void;
}

declare global {
    interface Window {
        registerPlugin(id: string, plugin: ClientPlugin): void;
    }
}

