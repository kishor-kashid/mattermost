// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {Client4} from 'mattermost-redux/client';
import {ActionFunc} from 'mattermost-redux/types/actions';

import {AIActionItemsTypes} from 'utils/constants/ai';

import type {AIActionItem, ActionItemCreateRequest, ActionItemUpdateRequest, ActionItemFilters, ActionItemStats} from 'types/ai';

export function getActionItems(filters?: ActionItemFilters): ActionFunc {
    return async (dispatch) => {
        dispatch({type: AIActionItemsTypes.GET_ACTION_ITEMS_REQUEST});

        try {
            const items = await Client4.getActionItems(filters);

            dispatch({
                type: AIActionItemsTypes.GET_ACTION_ITEMS_SUCCESS,
                data: items,
            });

            return {data: items};
        } catch (error) {
            dispatch({
                type: AIActionItemsTypes.GET_ACTION_ITEMS_FAILURE,
                error,
            });

            return {error};
        }
    };
}

export function getActionItem(actionItemId: string): ActionFunc {
    return async (dispatch) => {
        dispatch({type: AIActionItemsTypes.GET_ACTION_ITEM_REQUEST});

        try {
            const item = await Client4.getActionItem(actionItemId);

            dispatch({
                type: AIActionItemsTypes.GET_ACTION_ITEM_SUCCESS,
                data: item,
            });

            return {data: item};
        } catch (error) {
            dispatch({
                type: AIActionItemsTypes.GET_ACTION_ITEM_FAILURE,
                error,
            });

            return {error};
        }
    };
}

export function createActionItem(request: ActionItemCreateRequest): ActionFunc {
    return async (dispatch) => {
        dispatch({type: AIActionItemsTypes.CREATE_ACTION_ITEM_REQUEST});

        try {
            const item = await Client4.createActionItem(request);

            dispatch({
                type: AIActionItemsTypes.CREATE_ACTION_ITEM_SUCCESS,
                data: item,
            });

            return {data: item};
        } catch (error) {
            dispatch({
                type: AIActionItemsTypes.CREATE_ACTION_ITEM_FAILURE,
                error,
            });

            return {error};
        }
    };
}

export function updateActionItem(actionItemId: string, request: ActionItemUpdateRequest): ActionFunc {
    return async (dispatch) => {
        dispatch({type: AIActionItemsTypes.UPDATE_ACTION_ITEM_REQUEST});

        try {
            const item = await Client4.updateActionItem(actionItemId, request);

            dispatch({
                type: AIActionItemsTypes.UPDATE_ACTION_ITEM_SUCCESS,
                data: item,
            });

            return {data: item};
        } catch (error) {
            dispatch({
                type: AIActionItemsTypes.UPDATE_ACTION_ITEM_FAILURE,
                error,
            });

            return {error};
        }
    };
}

export function completeActionItem(actionItemId: string): ActionFunc {
    return async (dispatch) => {
        // Optimistic update
        dispatch({
            type: AIActionItemsTypes.COMPLETE_ACTION_ITEM_REQUEST,
            data: {id: actionItemId},
        });

        try {
            const item = await Client4.completeActionItem(actionItemId);

            dispatch({
                type: AIActionItemsTypes.COMPLETE_ACTION_ITEM_SUCCESS,
                data: item,
            });

            return {data: item};
        } catch (error) {
            dispatch({
                type: AIActionItemsTypes.COMPLETE_ACTION_ITEM_FAILURE,
                error,
                data: {id: actionItemId},
            });

            return {error};
        }
    };
}

export function deleteActionItem(actionItemId: string): ActionFunc {
    return async (dispatch) => {
        dispatch({
            type: AIActionItemsTypes.DELETE_ACTION_ITEM_REQUEST,
            data: {id: actionItemId},
        });

        try {
            await Client4.deleteActionItem(actionItemId);

            dispatch({
                type: AIActionItemsTypes.DELETE_ACTION_ITEM_SUCCESS,
                data: {id: actionItemId},
            });

            return {data: true};
        } catch (error) {
            dispatch({
                type: AIActionItemsTypes.DELETE_ACTION_ITEM_FAILURE,
                error,
                data: {id: actionItemId},
            });

            return {error};
        }
    };
}

export function getActionItemStats(userId?: string): ActionFunc {
    return async (dispatch) => {
        dispatch({type: AIActionItemsTypes.GET_ACTION_ITEM_STATS_REQUEST});

        try {
            const stats = await Client4.getActionItemStats(userId);

            dispatch({
                type: AIActionItemsTypes.GET_ACTION_ITEM_STATS_SUCCESS,
                data: stats,
            });

            return {data: stats};
        } catch (error) {
            dispatch({
                type: AIActionItemsTypes.GET_ACTION_ITEM_STATS_FAILURE,
                error,
            });

            return {error};
        }
    };
}

