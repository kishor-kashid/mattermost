// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import type {AIActionItem, AIAnalytics, AIPreferences, AISummary} from 'types/ai';

export type AIState = {
    summaries: AISummariesState;
    actionItems: AIActionItemsState;
    analytics: AIAnalyticsState;
    preferences: AIPreferencesState;
    system: AISystemState;
};

export type AISummariesState = {
    byId: Record<string, AISummary>;
    byChannel: Record<string, string[]>;
    loading: boolean;
    error: string | null;
};

export type AIActionItemsState = {
    byId: Record<string, AIActionItem>;
    byUser: Record<string, string[]>;
    byChannel: Record<string, string[]>;
    loading: boolean;
    error: string | null;
};

export type AIAnalyticsState = {
    byChannel: Record<string, AIAnalytics[]>;
    loading: boolean;
    error: string | null;
};

export type AIPreferencesState = {
    byUser: Record<string, AIPreferences>;
    loading: boolean;
    error: string | null;
};

export type AISystemState = {
    enabled: boolean;
    features: {
        summarization: boolean;
        analytics: boolean;
        actionItems: boolean;
        formatting: boolean;
    };
    health: {
        available: boolean;
        openaiConfigured: boolean;
        lastCheck: number | null;
    };
};

