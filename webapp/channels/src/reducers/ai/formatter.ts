// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import type {AnyAction} from 'redux';

import type {AIFormatterState} from 'types/store/ai';

import AIActionTypes from 'utils/constants/ai';

const initialState: AIFormatterState = {
    preview: null,
    profiles: [],
    loading: false,
    formatting: false,
    error: null,
};

export default function formatterReducer(state = initialState, action: AnyAction): AIFormatterState {
    switch (action.type) {
    case AIActionTypes.AI_FORMAT_PREVIEW_REQUEST:
        return {
            ...state,
            loading: true,
            formatting: true,
            error: null,
        };

    case AIActionTypes.AI_FORMAT_PREVIEW_SUCCESS: {
        const response = action.data;
        return {
            ...state,
            preview: response,
            loading: false,
            formatting: false,
            error: null,
        };
    }

    case AIActionTypes.AI_FORMAT_PREVIEW_FAILURE:
        return {
            ...state,
            loading: false,
            formatting: false,
            error: action.error || 'Failed to preview formatting',
        };

    case AIActionTypes.AI_FORMAT_APPLY_REQUEST:
        return {
            ...state,
            loading: true,
            formatting: true,
            error: null,
        };

    case AIActionTypes.AI_FORMAT_APPLY_SUCCESS:
        return {
            ...state,
            preview: null,
            loading: false,
            formatting: false,
            error: null,
        };

    case AIActionTypes.AI_FORMAT_APPLY_FAILURE:
        return {
            ...state,
            loading: false,
            formatting: false,
            error: action.error || 'Failed to apply formatting',
        };

    case AIActionTypes.AI_GET_FORMATTING_PROFILES_REQUEST:
        return {
            ...state,
            loading: true,
            error: null,
        };

    case AIActionTypes.AI_GET_FORMATTING_PROFILES_SUCCESS: {
        const profiles = action.data || [];
        return {
            ...state,
            profiles,
            loading: false,
            error: null,
        };
    }

    case AIActionTypes.AI_GET_FORMATTING_PROFILES_FAILURE:
        return {
            ...state,
            loading: false,
            error: action.error || 'Failed to load formatting profiles',
        };

    case AIActionTypes.AI_CLEAR_FORMAT_PREVIEW:
        return {
            ...state,
            preview: null,
            error: null,
        };

    default:
        return state;
    }
}

