// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import type {DispatchFunc, GetStateFunc} from '@mattermost/types/actions';
import type {FormatMessageRequest, FormatMessageResponse, FormattingProfileInfo} from 'types/ai';

import {Client4} from '@mattermost/client';

import {aiClient} from 'client/ai';

import AIActionTypes from 'utils/constants/ai';
import {logError} from 'mattermost-redux/actions/errors';

export function formatPreview(message: string, profile?: string, customInstructions?: string) {
    return async (dispatch: DispatchFunc, getState: GetStateFunc) => {
        dispatch({type: AIActionTypes.AI_FORMAT_PREVIEW_REQUEST});

        try {
            const request: FormatMessageRequest = {
                message,
                profile,
                custom_instructions: customInstructions,
            };

            const response = await aiClient.formatPreview(request);
            dispatch({
                type: AIActionTypes.AI_FORMAT_PREVIEW_SUCCESS,
                data: response,
            });

            return {data: response};
        } catch (error) {
            dispatch({
                type: AIActionTypes.AI_FORMAT_PREVIEW_FAILURE,
                error: error instanceof Error ? error.message : 'Failed to preview formatting',
            });
            dispatch(logError(error as Error));
            return {error};
        }
    };
}

export function formatApply(message: string, profile?: string, customInstructions?: string) {
    return async (dispatch: DispatchFunc, getState: GetStateFunc) => {
        dispatch({type: AIActionTypes.AI_FORMAT_APPLY_REQUEST});

        try {
            const request: FormatMessageRequest = {
                message,
                profile,
                custom_instructions: customInstructions,
            };

            const response = await aiClient.formatApply(request);
            dispatch({
                type: AIActionTypes.AI_FORMAT_APPLY_SUCCESS,
                data: response,
            });

            return {data: response};
        } catch (error) {
            dispatch({
                type: AIActionTypes.AI_FORMAT_APPLY_FAILURE,
                error: error instanceof Error ? error.message : 'Failed to apply formatting',
            });
            dispatch(logError(error as Error));
            return {error};
        }
    };
}

export function getFormattingProfiles() {
    return async (dispatch: DispatchFunc, getState: GetStateFunc) => {
        dispatch({type: AIActionTypes.AI_GET_FORMATTING_PROFILES_REQUEST});

        try {
            const profiles = await aiClient.getFormattingProfiles();
            dispatch({
                type: AIActionTypes.AI_GET_FORMATTING_PROFILES_SUCCESS,
                data: profiles,
            });

            return {data: profiles};
        } catch (error) {
            dispatch({
                type: AIActionTypes.AI_GET_FORMATTING_PROFILES_FAILURE,
                error: error instanceof Error ? error.message : 'Failed to load formatting profiles',
            });
            dispatch(logError(error as Error));
            return {error};
        }
    };
}

export function clearFormatPreview() {
    return {
        type: AIActionTypes.AI_CLEAR_FORMAT_PREVIEW,
    };
}

