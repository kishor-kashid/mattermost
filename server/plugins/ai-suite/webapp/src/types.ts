import type {ComponentType, ReactNode} from 'react';

export type PluginManifest = {
    id: string;
    name: string;
    description: string;
    version: string;
    min_server_version: string;
};

export type RHSRegistration = {
    id: string;
    showRHSPlugin: () => void;
    hideRHSPlugin: () => void;
    toggleRHSPlugin: () => void;
};

export type PluginRegistry = {
    registerChannelHeaderButtonAction?: (...args: any[]) => void;
    registerChannelHeaderMenuAction?: (text: ReactNode, action: (channelId: string) => void) => string;
    registerPostDropdownMenuAction?: (text: ReactNode, action: (postId: string) => void, filter: (postId: string) => boolean) => string;
    registerRightHandSidebarComponent?: (component: ComponentType, title: ReactNode) => RHSRegistration;
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

