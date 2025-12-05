// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {combineReducers} from 'redux';

import summaries from './summaries';
import actionItems from './action_items';
import analytics from './analytics';
import preferences from './preferences';
import system from './system';

export default combineReducers({
    summaries,
    actionItems,
    analytics,
    preferences,
    system,
});

