// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import {FormattedMessage} from 'react-intl';
import {useDispatch} from 'react-redux';

import {summarizeThread} from 'actions/ai_summarizer';

interface Props {
    postId: string;
    channelId: string;
    onClose?: () => void;
}

/**
 * Menu item component for thread summarization
 * This can be added to the post dot menu
 */
const ThreadSummarizeMenuItem: React.FC<Props> = ({postId, channelId, onClose}) => {
    const dispatch = useDispatch();

    const handleSummarize = () => {
        dispatch(summarizeThread({
            channel_id: channelId,
            post_id: postId,
            summary_level: 'standard',
            use_cache: true,
        }));

        if (onClose) {
            onClose();
        }
    };

    return (
        <li
            className='ai-summarize-thread-menu-item'
            role='menuitem'
        >
            <button
                className='style--none'
                role='presentation'
                onClick={handleSummarize}
            >
                <span className='icon icon-ai-sparkles'/>
                <FormattedMessage
                    id='ai.summary.threadMenuItem'
                    defaultMessage='Summarize Thread'
                />
            </button>
        </li>
    );
};

export default ThreadSummarizeMenuItem;

