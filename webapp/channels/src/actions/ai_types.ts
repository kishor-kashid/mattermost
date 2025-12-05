// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import type {AIActionItem, AIAnalytics, AIPreferences, AISummary} from 'types/ai';

export const AIActionTypes = {
    // Action Items
    RECEIVED_AI_ACTION_ITEM: 'RECEIVED_AI_ACTION_ITEM',
    RECEIVED_AI_ACTION_ITEMS: 'RECEIVED_AI_ACTION_ITEMS',
    UPDATED_AI_ACTION_ITEM: 'UPDATED_AI_ACTION_ITEM',
    DELETED_AI_ACTION_ITEM: 'DELETED_AI_ACTION_ITEM',

    // Summaries
    RECEIVED_AI_SUMMARY: 'RECEIVED_AI_SUMMARY',
    RECEIVED_AI_SUMMARIES: 'RECEIVED_AI_SUMMARIES',
    DELETED_AI_SUMMARY: 'DELETED_AI_SUMMARY',

    // Analytics
    RECEIVED_AI_ANALYTICS: 'RECEIVED_AI_ANALYTICS',
    RECEIVED_AI_ANALYTICS_LIST: 'RECEIVED_AI_ANALYTICS_LIST',

    // Preferences
    RECEIVED_AI_PREFERENCES: 'RECEIVED_AI_PREFERENCES',
    UPDATED_AI_PREFERENCES: 'UPDATED_AI_PREFERENCES',

    // Loading states
    AI_REQUEST_STARTED: 'AI_REQUEST_STARTED',
    AI_REQUEST_SUCCESS: 'AI_REQUEST_SUCCESS',
    AI_REQUEST_FAILURE: 'AI_REQUEST_FAILURE',
} as const;

export type AIAction =
    | ReceivedAIActionItemAction
    | ReceivedAIActionItemsAction
    | UpdatedAIActionItemAction
    | DeletedAIActionItemAction
    | ReceivedAISummaryAction
    | ReceivedAISummariesAction
    | DeletedAISummaryAction
    | ReceivedAIAnalyticsAction
    | ReceivedAIAnalyticsListAction
    | ReceivedAIPreferencesAction
    | UpdatedAIPreferencesAction
    | AIRequestStartedAction
    | AIRequestSuccessAction
    | AIRequestFailureAction;

export interface ReceivedAIActionItemAction {
    type: typeof AIActionTypes.RECEIVED_AI_ACTION_ITEM;
    data: AIActionItem;
}

export interface ReceivedAIActionItemsAction {
    type: typeof AIActionTypes.RECEIVED_AI_ACTION_ITEMS;
    data: AIActionItem[];
}

export interface UpdatedAIActionItemAction {
    type: typeof AIActionTypes.UPDATED_AI_ACTION_ITEM;
    data: AIActionItem;
}

export interface DeletedAIActionItemAction {
    type: typeof AIActionTypes.DELETED_AI_ACTION_ITEM;
    data: {id: string};
}

export interface ReceivedAISummaryAction {
    type: typeof AIActionTypes.RECEIVED_AI_SUMMARY;
    data: AISummary;
}

export interface ReceivedAISummariesAction {
    type: typeof AIActionTypes.RECEIVED_AI_SUMMARIES;
    data: AISummary[];
}

export interface DeletedAISummaryAction {
    type: typeof AIActionTypes.DELETED_AI_SUMMARY;
    data: {id: string};
}

export interface ReceivedAIAnalyticsAction {
    type: typeof AIActionTypes.RECEIVED_AI_ANALYTICS;
    data: AIAnalytics;
}

export interface ReceivedAIAnalyticsListAction {
    type: typeof AIActionTypes.RECEIVED_AI_ANALYTICS_LIST;
    data: AIAnalytics[];
}

export interface ReceivedAIPreferencesAction {
    type: typeof AIActionTypes.RECEIVED_AI_PREFERENCES;
    data: AIPreferences;
}

export interface UpdatedAIPreferencesAction {
    type: typeof AIActionTypes.UPDATED_AI_PREFERENCES;
    data: AIPreferences;
}

export interface AIRequestStartedAction {
    type: typeof AIActionTypes.AI_REQUEST_STARTED;
    data: {key: string};
}

export interface AIRequestSuccessAction {
    type: typeof AIActionTypes.AI_REQUEST_SUCCESS;
    data: {key: string};
}

export interface AIRequestFailureAction {
    type: typeof AIActionTypes.AI_REQUEST_FAILURE;
    data: {
        key: string;
        error: string;
    };
}

