// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import type {AnyAction} from 'redux';

import type {AISummariesState} from 'types/store/ai';

import AIActionTypes from 'utils/constants/ai';

const initialState: AISummariesState = {
    byId: {},
    byChannel: {},
    loading: false,
    error: null,
};

export default function summariesReducer(state = initialState, action: AnyAction): AISummariesState {
    switch (action.type) {
    case AIActionTypes.AI_SUMMARY_REQUEST:
    case AIActionTypes.AI_SUMMARIES_BY_CHANNEL_REQUEST:
    case AIActionTypes.AI_SUMMARIZE_THREAD_REQUEST:
    case AIActionTypes.AI_SUMMARIZE_CHANNEL_REQUEST:
        return {
            ...state,
            loading: true,
            error: null,
        };

    case AIActionTypes.AI_SUMMARY_SUCCESS: {
        const summary = action.data;
        return {
            ...state,
            byId: {
                ...state.byId,
                [summary.id]: summary,
            },
            byChannel: {
                ...state.byChannel,
                [summary.channel_id]: [
                    ...(state.byChannel[summary.channel_id] || []).filter((id) => id !== summary.id),
                    summary.id,
                ],
            },
            loading: false,
            error: null,
        };
    }

    case AIActionTypes.AI_SUMMARIZE_THREAD_SUCCESS:
    case AIActionTypes.AI_SUMMARIZE_CHANNEL_SUCCESS: {
        const {summary} = action.data;
        return {
            ...state,
            byId: {
                ...state.byId,
                [summary.id]: summary,
            },
            byChannel: {
                ...state.byChannel,
                [summary.channel_id]: [
                    ...(state.byChannel[summary.channel_id] || []).filter((id) => id !== summary.id),
                    summary.id,
                ],
            },
            loading: false,
            error: null,
        };
    }

    case AIActionTypes.AI_SUMMARIES_BY_CHANNEL_SUCCESS: {
        const {channelId, summaries} = action.data;
        const byId = {...state.byId};
        const summaryIds: string[] = [];

        summaries.forEach((summary: any) => {
            byId[summary.id] = summary;
            summaryIds.push(summary.id);
        });

        return {
            ...state,
            byId,
            byChannel: {
                ...state.byChannel,
                [channelId]: summaryIds,
            },
            loading: false,
            error: null,
        };
    }

    case AIActionTypes.AI_SUMMARY_FAILURE:
    case AIActionTypes.AI_SUMMARIES_BY_CHANNEL_FAILURE:
    case AIActionTypes.AI_SUMMARIZE_THREAD_FAILURE:
    case AIActionTypes.AI_SUMMARIZE_CHANNEL_FAILURE:
        return {
            ...state,
            loading: false,
            error: action.error,
        };

    case AIActionTypes.AI_CLEAR_SUMMARY:
        return initialState;

    default:
        return state;
    }
}

