// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import type {AnyAction} from 'redux';

import type {AISystemState} from 'types/store/ai';

import AIActionTypes from 'utils/constants/ai';

const initialState: AISystemState = {
    enabled: false,
    features: {
        summarization: false,
        analytics: false,
        actionItems: false,
        formatting: false,
    },
    health: {
        available: false,
        openaiConfigured: false,
        lastCheck: null,
    },
};

export default function systemReducer(state = initialState, action: AnyAction): AISystemState {
    switch (action.type) {
    case AIActionTypes.AI_HEALTH_CHECK_SUCCESS: {
        const {enabled, service_available, openai_configured, features} = action.data;
        return {
            ...state,
            enabled: enabled || false,
            features: {
                summarization: features?.summarization || false,
                analytics: features?.analytics || false,
                actionItems: features?.action_items || false,
                formatting: features?.formatting || false,
            },
            health: {
                available: service_available || false,
                openaiConfigured: openai_configured || false,
                lastCheck: Date.now(),
            },
        };
    }

    case AIActionTypes.AI_HEALTH_CHECK_FAILURE:
        return {
            ...state,
            health: {
                available: false,
                openaiConfigured: false,
                lastCheck: Date.now(),
            },
        };

    default:
        return state;
    }
}

