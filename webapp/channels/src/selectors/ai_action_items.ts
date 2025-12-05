// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {createSelector} from 'reselect';

import {GlobalState} from 'types/store';
import type {AIActionItem} from 'types/ai';

export function getActionItemsState(state: GlobalState) {
    return state.ai.actionItems;
}

export const getActionItemsArray = createSelector(
    'getActionItemsArray',
    getActionItemsState,
    (actionItemsState) => {
        return Object.values(actionItemsState.items);
    },
);

export const getActionItemById = (state: GlobalState, actionItemId: string): AIActionItem | undefined => {
    return state.ai.actionItems.items[actionItemId];
};

export const getOverdueActionItems = createSelector(
    'getOverdueActionItems',
    getActionItemsArray,
    (items) => {
        const now = Date.now();
        return items.filter((item) => {
            return item.due_date && item.due_date < now && item.status !== 'completed' && item.status !== 'dismissed';
        });
    },
);

export const getDueSoonActionItems = createSelector(
    'getDueSoonActionItems',
    getActionItemsArray,
    (items) => {
        const now = Date.now();
        const weekFromNow = now + (7 * 24 * 60 * 60 * 1000);
        return items.filter((item) => {
            return item.due_date && item.due_date >= now && item.due_date <= weekFromNow && item.status !== 'completed' && item.status !== 'dismissed';
        });
    },
);

export const getActiveActionItems = createSelector(
    'getActiveActionItems',
    getActionItemsArray,
    (items) => {
        return items.filter((item) => item.status !== 'completed' && item.status !== 'dismissed');
    },
);

export const getCompletedActionItems = createSelector(
    'getCompletedActionItems',
    getActionItemsArray,
    (items) => {
        return items.filter((item) => item.status === 'completed');
    },
);

export const getActionItemsByChannel = (state: GlobalState, channelId: string): AIActionItem[] => {
    const items = getActionItemsArray(state);
    return items.filter((item) => item.channel_id === channelId);
};

export const getActionItemsByPriority = createSelector(
    'getActionItemsByPriority',
    getActiveActionItems,
    (items) => {
        return {
            urgent: items.filter((item) => item.priority === 'urgent'),
            high: items.filter((item) => item.priority === 'high'),
            medium: items.filter((item) => item.priority === 'medium'),
            low: items.filter((item) => item.priority === 'low'),
        };
    },
);

export const getActionItemsLoading = (state: GlobalState): boolean => {
    return state.ai.actionItems.loading;
};

export const getActionItemsError = (state: GlobalState): string | null => {
    return state.ai.actionItems.error;
};

export const getActionItemStats = (state: GlobalState) => {
    return state.ai.actionItems.stats;
};

