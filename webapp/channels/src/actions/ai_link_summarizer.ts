import type {Dispatch} from 'redux';

import AIActionTypes from 'utils/constants/ai';
import {aiClient} from 'client/ai';
import type {GlobalState} from 'types/store';
import type {AILinkSummary} from 'types/ai';

const normalizeUrl = (url: string) => url.trim().toLowerCase();

export function summarizeLink(url: string, forceRefresh = false) {
    return async (dispatch: Dispatch, getState: () => GlobalState) => {
        const key = normalizeUrl(url);

        console.log('[AI Link Summary] Starting summarization for:', key);
        dispatch({type: AIActionTypes.AI_LINK_SUMMARY_REQUEST, data: {url: key}});

        try {
            console.log('[AI Link Summary] Calling API...');
            const resp = await aiClient.summarizeLink(url, forceRefresh);
            console.log('[AI Link Summary] API response:', resp);
            const summary: AILinkSummary = resp.summary || resp;
            console.log('[AI Link Summary] Extracted summary:', summary);
            dispatch({
                type: AIActionTypes.AI_LINK_SUMMARY_SUCCESS,
                data: {url: key, summary, fromCache: resp.from_cache},
            });
            return summary;
        } catch (error: any) {
            console.error('[AI Link Summary] Error:', error);
            dispatch({
                type: AIActionTypes.AI_LINK_SUMMARY_FAILURE,
                data: {url: key, error: error?.message || 'Link summary failed'},
            });
            throw error;
        }
    };
}

export function getLinkSummary(url: string, useCache = true) {
    return async (dispatch: Dispatch) => {
        const key = normalizeUrl(url);
        dispatch({type: AIActionTypes.AI_LINK_SUMMARY_REQUEST, data: {url: key}});

        try {
            const resp = await aiClient.getLinkSummary(url, useCache);
            const summary: AILinkSummary = resp.summary || resp;
            dispatch({
                type: AIActionTypes.AI_LINK_SUMMARY_SUCCESS,
                data: {url: key, summary, fromCache: resp.from_cache},
            });
            return summary;
        } catch (error: any) {
            dispatch({
                type: AIActionTypes.AI_LINK_SUMMARY_FAILURE,
                data: {url: key, error: error?.message || 'Link summary failed'},
            });
            throw error;
        }
    };
}

