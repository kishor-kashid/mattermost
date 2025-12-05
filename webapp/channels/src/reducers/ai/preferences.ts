// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import type {AnyAction} from 'redux';

import type {AIPreferencesState} from 'types/store/ai';

import AIActionTypes from 'utils/constants/ai';

const initialState: AIPreferencesState = {
    byUser: {},
    loading: false,
    error: null,
};

export default function preferencesReducer(state = initialState, action: AnyAction): AIPreferencesState {
    switch (action.type) {
    case AIActionTypes.AI_PREFERENCES_GET_REQUEST:
    case AIActionTypes.AI_PREFERENCES_UPDATE_REQUEST:
        return {
            ...state,
            loading: true,
            error: null,
        };

    case AIActionTypes.AI_PREFERENCES_GET_SUCCESS:
    case AIActionTypes.AI_PREFERENCES_UPDATE_SUCCESS: {
        const preferences = action.data;
        return {
            ...state,
            byUser: {
                ...state.byUser,
                [preferences.user_id]: preferences,
            },
            loading: false,
            error: null,
        };
    }

    case AIActionTypes.AI_PREFERENCES_GET_FAILURE:
    case AIActionTypes.AI_PREFERENCES_UPDATE_FAILURE:
        return {
            ...state,
            loading: false,
            error: action.error,
        };

    default:
        return state;
    }
}

