// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useState} from 'react';
import {FormattedMessage} from 'react-intl';
import {useDispatch} from 'react-redux';

import {summarizeChannel} from 'actions/ai_summarizer';
import DateRangeModal from './date_range_modal';

interface Props {
    channelId: string;
    onClose?: () => void;
}

/**
 * Menu item component for channel summarization
 * This can be added to the channel header menu
 */
const ChannelSummarizeMenuItem: React.FC<Props> = ({channelId, onClose}) => {
    const dispatch = useDispatch();
    const [showModal, setShowModal] = useState(false);

    const handleClick = () => {
        setShowModal(true);
        if (onClose) {
            onClose();
        }
    };

    const handleConfirm = (startTime: number, endTime: number, level: string) => {
        dispatch(summarizeChannel({
            channel_id: channelId,
            start_time: startTime,
            end_time: endTime,
            summary_level: level as 'brief' | 'standard' | 'detailed',
            use_cache: true,
        }));
    };

    return (
        <>
            <li
                className='ai-summarize-channel-menu-item'
                role='menuitem'
            >
                <button
                    className='style--none'
                    role='presentation'
                    onClick={handleClick}
                >
                    <span className='icon icon-ai-sparkles'/>
                    <FormattedMessage
                        id='ai.summary.channelMenuItem'
                        defaultMessage='Summarize Channel'
                    />
                </button>
            </li>

            <DateRangeModal
                show={showModal}
                onHide={() => setShowModal(false)}
                onConfirm={handleConfirm}
            />
        </>
    );
};

export default ChannelSummarizeMenuItem;

