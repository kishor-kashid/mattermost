// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import type {AnyAction} from 'redux';

import type {AIAnalyticsState} from 'types/store/ai';

import AIActionTypes from 'utils/constants/ai';

const initialState: AIAnalyticsState = {
    byChannel: {},
    loading: false,
    error: null,
};

export default function analyticsReducer(state = initialState, action: AnyAction): AIAnalyticsState {
    switch (action.type) {
    case AIActionTypes.AI_ANALYTICS_REQUEST:
        return {
            ...state,
            loading: true,
            error: null,
        };

    case AIActionTypes.AI_ANALYTICS_SUCCESS: {
        const {channelId, analytics} = action.data;
        return {
            ...state,
            byChannel: {
                ...state.byChannel,
                [channelId]: analytics,
            },
            loading: false,
            error: null,
        };
    }

    case AIActionTypes.AI_ANALYTICS_FAILURE:
        return {
            ...state,
            loading: false,
            error: action.error,
        };

    default:
        return state;
    }
}

