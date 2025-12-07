import type {AnyAction} from 'redux';

import AIActionTypes from 'utils/constants/ai';
import type {AILinkSummariesState} from 'types/store/ai';

const initialState: AILinkSummariesState = {
    byUrl: {},
    loading: {},
    error: {},
};

export default function linkSummariesReducer(state = initialState, action: AnyAction): AILinkSummariesState {
    switch (action.type) {
    case AIActionTypes.AI_LINK_SUMMARY_REQUEST: {
        const {url} = action.data;
        console.log('[Reducer] AI_LINK_SUMMARY_REQUEST for:', url);
        return {
            ...state,
            loading: {...state.loading, [url]: true},
            error: {...state.error, [url]: null},
        };
    }
    case AIActionTypes.AI_LINK_SUMMARY_SUCCESS: {
        const {url, summary, fromCache} = action.data;
        console.log('[Reducer] AI_LINK_SUMMARY_SUCCESS for:', url, summary);
        return {
            ...state,
            byUrl: {
                ...state.byUrl,
                [url]: {summary, fromCache},
            },
            loading: {...state.loading, [url]: false},
            error: {...state.error, [url]: null},
        };
    }
    case AIActionTypes.AI_LINK_SUMMARY_FAILURE: {
        const {url, error} = action.data;
        console.log('[Reducer] AI_LINK_SUMMARY_FAILURE for:', url, error);
        return {
            ...state,
            loading: {...state.loading, [url]: false},
            error: {...state.error, [url]: error},
        };
    }
    default:
        return state;
    }
}

