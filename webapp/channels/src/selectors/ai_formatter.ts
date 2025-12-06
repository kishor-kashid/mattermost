// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {createSelector} from 'mattermost-redux/selectors/create_selector';

import type {GlobalState} from '@mattermost/types/store';

const getAIState = (state: GlobalState) => state.entities.ai;

export const getFormatterState = createSelector(
    'getFormatterState',
    getAIState,
    (ai) => ai?.formatter || {
        preview: null,
        profiles: [],
        loading: false,
        formatting: false,
        error: null,
    },
);

export const getFormatPreview = createSelector(
    'getFormatPreview',
    getFormatterState,
    (formatter) => formatter.preview,
);

export const getFormattingProfiles = createSelector(
    'getFormattingProfiles',
    getFormatterState,
    (formatter) => formatter.profiles,
);

export const isFormatting = createSelector(
    'isFormatting',
    getFormatterState,
    (formatter) => formatter.formatting || formatter.loading,
);

export const getFormatterError = createSelector(
    'getFormatterError',
    getFormatterState,
    (formatter) => formatter.error,
);

