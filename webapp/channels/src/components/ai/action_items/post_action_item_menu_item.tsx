// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useState} from 'react';
import {useDispatch} from 'react-redux';
import {FormattedMessage} from 'react-intl';

import {CheckCircleOutlineIcon} from '@mattermost/compass-icons/components';

import * as Menu from 'components/menu';

import CreateActionItemModal from './create_modal';

interface Props {
    postId: string;
    channelId: string;
    onClose?: () => void;
}

/**
 * Menu item component for creating action items from posts
 * This can be added to the post dot menu
 */
const PostActionItemMenuItem: React.FC<Props> = ({postId, channelId, onClose}) => {
    const dispatch = useDispatch();
    const [showModal, setShowModal] = useState(false);

    const handleClick = () => {
        setShowModal(true);
        if (onClose) {
            onClose();
        }
    };

    return (
        <>
            <Menu.Item
                id={`create_action_item_${postId}`}
                data-testid={`create_action_item_${postId}`}
                leadingElement={<CheckCircleOutlineIcon size={18}/>}
                labels={
                    <FormattedMessage
                        id='ai.action_items.createFromPost'
                        defaultMessage='Create Action Item'
                    />
                }
                onClick={handleClick}
            />

            {showModal && (
                <CreateActionItemModal
                    show={showModal}
                    onHide={() => setShowModal(false)}
                    channelId={channelId}
                    postId={postId}
                />
            )}
        </>
    );
};

export default PostActionItemMenuItem;

