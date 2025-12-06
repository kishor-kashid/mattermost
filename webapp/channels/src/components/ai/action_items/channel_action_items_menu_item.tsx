// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import {useDispatch} from 'react-redux';
import {FormattedMessage} from 'react-intl';

import {CheckCircleOutlineIcon} from '@mattermost/compass-icons/components';

import {openRHSForActionItems} from 'actions/views/rhs';

import * as Menu from 'components/menu';

interface Props extends Menu.FirstMenuItemProps {
    channelId: string;
    onClose?: () => void;
}

/**
 * Menu item component for viewing channel action items
 * This can be added to the channel header menu
 */
const ChannelActionItemsMenuItem: React.FC<Props> = ({channelId, onClose, ...rest}) => {
    const dispatch = useDispatch();

    const handleClick = () => {
        dispatch(openRHSForActionItems(channelId));
        if (onClose) {
            onClose();
        }
    };

    return (
        <Menu.Item
            leadingElement={<CheckCircleOutlineIcon size={18}/>}
            onClick={handleClick}
            labels={
                <FormattedMessage
                    id='ai.action_items.viewChannelItems'
                    defaultMessage='View Action Items'
                />
            }
            {...rest}
        />
    );
};

export default ChannelActionItemsMenuItem;

