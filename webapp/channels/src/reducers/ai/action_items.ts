// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import type {AnyAction} from 'redux';

import type {AIActionItemsState} from 'types/store/ai';
import type {AIActionItem} from 'types/ai';

import {AIActionItemsTypes} from 'utils/constants/ai';

const initialState: AIActionItemsState = {
    items: {},
    loading: false,
    error: null,
    stats: null,
};

export default function actionItemsReducer(state = initialState, action: AnyAction): AIActionItemsState {
    switch (action.type) {
    case AIActionItemsTypes.GET_ACTION_ITEMS_REQUEST:
    case AIActionItemsTypes.GET_ACTION_ITEM_REQUEST:
    case AIActionItemsTypes.CREATE_ACTION_ITEM_REQUEST:
    case AIActionItemsTypes.UPDATE_ACTION_ITEM_REQUEST:
    case AIActionItemsTypes.GET_ACTION_ITEM_STATS_REQUEST:
        return {
            ...state,
            loading: true,
            error: null,
        };

    case AIActionItemsTypes.GET_ACTION_ITEMS_SUCCESS: {
        const items: AIActionItem[] = action.data;
        const itemsById: Record<string, AIActionItem> = {};
        
        items.forEach((item) => {
            itemsById[item.id] = item;
        });

        return {
            ...state,
            items: {
                ...state.items,
                ...itemsById,
            },
            loading: false,
            error: null,
        };
    }

    case AIActionItemsTypes.GET_ACTION_ITEM_SUCCESS:
    case AIActionItemsTypes.CREATE_ACTION_ITEM_SUCCESS:
    case AIActionItemsTypes.UPDATE_ACTION_ITEM_SUCCESS:
    case AIActionItemsTypes.COMPLETE_ACTION_ITEM_SUCCESS: {
        const item: AIActionItem = action.data;
        return {
            ...state,
            items: {
                ...state.items,
                [item.id]: item,
            },
            loading: false,
            error: null,
        };
    }

    case AIActionItemsTypes.COMPLETE_ACTION_ITEM_REQUEST: {
        // Optimistic update
        const {id} = action.data;
        const item = state.items[id];
        if (item) {
            return {
                ...state,
                items: {
                    ...state.items,
                    [id]: {
                        ...item,
                        status: 'completed',
                        completed_at: Date.now(),
                    },
                },
            };
        }
        return state;
    }

    case AIActionItemsTypes.DELETE_ACTION_ITEM_REQUEST: {
        // Optimistic delete
        const {id} = action.data;
        const {[id]: deleted, ...remainingItems} = state.items;

        return {
            ...state,
            items: remainingItems,
        };
    }

    case AIActionItemsTypes.DELETE_ACTION_ITEM_SUCCESS:
        return {
            ...state,
            loading: false,
            error: null,
        };

    case AIActionItemsTypes.COMPLETE_ACTION_ITEM_FAILURE:
    case AIActionItemsTypes.DELETE_ACTION_ITEM_FAILURE: {
        // Revert optimistic update
        // For now, just refetch - in a real app, you'd revert the specific change
        return {
            ...state,
            loading: false,
            error: action.error,
        };
    }

    case AIActionItemsTypes.GET_ACTION_ITEM_STATS_SUCCESS:
        return {
            ...state,
            stats: action.data,
            loading: false,
            error: null,
        };

    case AIActionItemsTypes.GET_ACTION_ITEMS_FAILURE:
    case AIActionItemsTypes.GET_ACTION_ITEM_FAILURE:
    case AIActionItemsTypes.CREATE_ACTION_ITEM_FAILURE:
    case AIActionItemsTypes.UPDATE_ACTION_ITEM_FAILURE:
    case AIActionItemsTypes.GET_ACTION_ITEM_STATS_FAILURE:
        return {
            ...state,
            loading: false,
            error: action.error,
        };

    default:
        return state;
    }
}
