import type {Dispatch} from 'redux';

import {aiClient} from 'client/ai';
import AIActionTypes from 'utils/constants/ai';

export function getAIPreferences(userId: string) {
    return async (dispatch: Dispatch) => {
        dispatch({type: AIActionTypes.AI_PREFERENCES_GET_REQUEST});
        try {
            const prefs = await aiClient.getPreferences(userId);
            dispatch({type: AIActionTypes.AI_PREFERENCES_GET_SUCCESS, data: prefs});
            return prefs;
        } catch (error: any) {
            dispatch({type: AIActionTypes.AI_PREFERENCES_GET_FAILURE, error: error?.message || 'Unable to load AI preferences'});
            throw error;
        }
    };
}

export function updateAIPreferences(userId: string, preferences: Record<string, unknown>) {
    return async (dispatch: Dispatch) => {
        dispatch({type: AIActionTypes.AI_PREFERENCES_UPDATE_REQUEST});
        try {
            const prefs = await aiClient.updatePreferences(userId, preferences);
            dispatch({type: AIActionTypes.AI_PREFERENCES_UPDATE_SUCCESS, data: prefs});
            return prefs;
        } catch (error: any) {
            dispatch({type: AIActionTypes.AI_PREFERENCES_UPDATE_FAILURE, error: error?.message || 'Unable to update AI preferences'});
            throw error;
        }
    };
}

