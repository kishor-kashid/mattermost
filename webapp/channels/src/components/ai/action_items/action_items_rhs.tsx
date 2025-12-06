// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import {useSelector, useDispatch} from 'react-redux';
import {FormattedMessage} from 'react-intl';

import {getCurrentChannelId} from 'mattermost-redux/selectors/entities/channels';

import {closeRightHandSide} from 'actions/views/rhs';

import ActionItemsDashboard from './dashboard';

import './action_items_rhs.scss';

const ActionItemsRHS: React.FC = () => {
    const dispatch = useDispatch();
    const channelId = useSelector(getCurrentChannelId);

    const handleClose = () => {
        dispatch(closeRightHandSide());
    };

    return (
        <div className='action-items-rhs'>
            <div className='rhs-header'>
                <h3>
                    <FormattedMessage
                        id='ai.action_items.rhs.title'
                        defaultMessage='Action Items'
                    />
                </h3>
                <button
                    className='close-btn'
                    onClick={handleClose}
                    aria-label='Close'
                >
                    <i className='icon icon-close'/>
                </button>
            </div>
            <div className='rhs-content'>
                <ActionItemsDashboard/>
            </div>
        </div>
    );
};

export default ActionItemsRHS;

