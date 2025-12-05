// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import type {AnyAction} from 'redux';

import type {AIActionItemsState} from 'types/store/ai';

import AIActionTypes from 'utils/constants/ai';

const initialState: AIActionItemsState = {
    byId: {},
    byUser: {},
    byChannel: {},
    loading: false,
    error: null,
};

export default function actionItemsReducer(state = initialState, action: AnyAction): AIActionItemsState {
    switch (action.type) {
    case AIActionTypes.AI_ACTION_ITEM_CREATE_REQUEST:
    case AIActionTypes.AI_ACTION_ITEM_GET_REQUEST:
    case AIActionTypes.AI_ACTION_ITEMS_BY_USER_REQUEST:
    case AIActionTypes.AI_ACTION_ITEMS_BY_CHANNEL_REQUEST:
    case AIActionTypes.AI_ACTION_ITEM_UPDATE_REQUEST:
        return {
            ...state,
            loading: true,
            error: null,
        };

    case AIActionTypes.AI_ACTION_ITEM_CREATE_SUCCESS:
    case AIActionTypes.AI_ACTION_ITEM_GET_SUCCESS:
    case AIActionTypes.AI_ACTION_ITEM_UPDATE_SUCCESS: {
        const actionItem = action.data;
        return {
            ...state,
            byId: {
                ...state.byId,
                [actionItem.id]: actionItem,
            },
            byUser: {
                ...state.byUser,
                [actionItem.assignee_id || actionItem.user_id]: [
                    ...(state.byUser[actionItem.assignee_id || actionItem.user_id] || []).filter((id) => id !== actionItem.id),
                    actionItem.id,
                ],
            },
            byChannel: {
                ...state.byChannel,
                [actionItem.channel_id]: [
                    ...(state.byChannel[actionItem.channel_id] || []).filter((id) => id !== actionItem.id),
                    actionItem.id,
                ],
            },
            loading: false,
            error: null,
        };
    }

    case AIActionTypes.AI_ACTION_ITEMS_BY_USER_SUCCESS: {
        const {userId, actionItems} = action.data;
        const byId = {...state.byId};
        const itemIds: string[] = [];

        actionItems.forEach((item: any) => {
            byId[item.id] = item;
            itemIds.push(item.id);
        });

        return {
            ...state,
            byId,
            byUser: {
                ...state.byUser,
                [userId]: itemIds,
            },
            loading: false,
            error: null,
        };
    }

    case AIActionTypes.AI_ACTION_ITEMS_BY_CHANNEL_SUCCESS: {
        const {channelId, actionItems} = action.data;
        const byId = {...state.byId};
        const itemIds: string[] = [];

        actionItems.forEach((item: any) => {
            byId[item.id] = item;
            itemIds.push(item.id);
        });

        return {
            ...state,
            byId,
            byChannel: {
                ...state.byChannel,
                [channelId]: itemIds,
            },
            loading: false,
            error: null,
        };
    }

    case AIActionTypes.AI_ACTION_ITEM_DELETE_SUCCESS: {
        const {id} = action.data;
        const {[id]: deleted, ...remainingById} = state.byId;

        return {
            ...state,
            byId: remainingById,
            loading: false,
            error: null,
        };
    }

    case AIActionTypes.AI_ACTION_ITEM_CREATE_FAILURE:
    case AIActionTypes.AI_ACTION_ITEM_GET_FAILURE:
    case AIActionTypes.AI_ACTION_ITEMS_BY_USER_FAILURE:
    case AIActionTypes.AI_ACTION_ITEMS_BY_CHANNEL_FAILURE:
    case AIActionTypes.AI_ACTION_ITEM_UPDATE_FAILURE:
    case AIActionTypes.AI_ACTION_ITEM_DELETE_FAILURE:
        return {
            ...state,
            loading: false,
            error: action.error,
        };

    default:
        return state;
    }
}

