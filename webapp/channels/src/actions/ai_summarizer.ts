// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import type {AISummary} from '@mattermost/types/ai';

import {ActionTypes} from 'utils/constants/ai';
import {Client4} from 'mattermost-redux/client';

export type SummarizeThreadRequest = {
    channel_id: string;
    post_id: string;
    summary_level?: 'brief' | 'standard' | 'detailed';
    use_cache?: boolean;
};

export type SummarizeChannelRequest = {
    channel_id: string;
    start_time?: number;
    end_time?: number;
    summary_level?: 'brief' | 'standard' | 'detailed';
    use_cache?: boolean;
};

export type SummaryResponse = {
    summary: AISummary;
    from_cache: boolean;
    processing_ms: number;
};

export function summarizeThread(request: SummarizeThreadRequest) {
    return async (dispatch: any) => {
        dispatch({
            type: ActionTypes.AI_SUMMARIZE_THREAD_REQUEST,
            data: request,
        });

        try {
            const response = await Client4.doFetch<SummaryResponse>(
                `/ai/summarize`,
                {
                    method: 'POST',
                    body: JSON.stringify(request),
                },
            );

            dispatch({
                type: ActionTypes.AI_SUMMARIZE_THREAD_SUCCESS,
                data: response,
            });

            return {data: response};
        } catch (error) {
            dispatch({
                type: ActionTypes.AI_SUMMARIZE_THREAD_FAILURE,
                error,
            });

            return {error};
        }
    };
}

export function summarizeChannel(request: SummarizeChannelRequest) {
    return async (dispatch: any) => {
        dispatch({
            type: ActionTypes.AI_SUMMARIZE_CHANNEL_REQUEST,
            data: request,
        });

        try {
            const response = await Client4.doFetch<SummaryResponse>(
                `/ai/summarize`,
                {
                    method: 'POST',
                    body: JSON.stringify(request),
                },
            );

            dispatch({
                type: ActionTypes.AI_SUMMARIZE_CHANNEL_SUCCESS,
                data: response,
            });

            return {data: response};
        } catch (error) {
            dispatch({
                type: ActionTypes.AI_SUMMARIZE_CHANNEL_FAILURE,
                error,
            });

            return {error};
        }
    };
}

export function clearSummary() {
    return {
        type: ActionTypes.AI_CLEAR_SUMMARY,
    };
}

export function setSummaryLevel(level: 'brief' | 'standard' | 'detailed') {
    return {
        type: ActionTypes.AI_SET_SUMMARY_LEVEL,
        data: level,
    };
}

