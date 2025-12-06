// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {createSelector} from 'mattermost-redux/selectors/create_selector';

import type {GlobalState} from 'types/store';
import type {AISummary} from '@mattermost/types/ai';

export function getSummariesState(state: GlobalState) {
    return state.ai.summaries;
}

export function getSummaryById(state: GlobalState, summaryId: string): AISummary | undefined {
    return state.ai.summaries.byId[summaryId];
}

export const getSummariesByChannel = createSelector(
    'getSummariesByChannel',
    getSummariesState,
    (_state: GlobalState, channelId: string) => channelId,
    (summariesState, channelId): AISummary[] => {
        const summaryIds = summariesState.byChannel[channelId] || [];
        return summaryIds.map((id) => summariesState.byId[id]).filter(Boolean);
    },
);

export const getLatestChannelSummary = createSelector(
    'getLatestChannelSummary',
    getSummariesByChannel,
    (summaries): AISummary | undefined => {
        if (summaries.length === 0) {
            return undefined;
        }

        return summaries.reduce((latest, current) => {
            return current.create_at > latest.create_at ? current : latest;
        });
    },
);

export function isSummariesLoading(state: GlobalState): boolean {
    return state.ai.summaries.loading;
}

export function getSummariesError(state: GlobalState): any {
    return state.ai.summaries.error;
}

export const getThreadSummary = createSelector(
    'getThreadSummary',
    getSummariesState,
    (_state: GlobalState, postId: string) => postId,
    (summariesState, postId): AISummary | undefined => {
        // Find summary where post_id matches
        return Object.values(summariesState.byId).find(
            (summary) => summary.post_id === postId && summary.summary_type === 'thread',
        );
    },
);

