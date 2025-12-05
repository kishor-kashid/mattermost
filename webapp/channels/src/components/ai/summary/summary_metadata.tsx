// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import {FormattedMessage, FormattedDate, FormattedTime} from 'react-intl';

import type {AISummary} from '@mattermost/types/ai';

interface Props {
    summary: AISummary;
}

const SummaryMetadata: React.FC<Props> = ({summary}) => {
    const startDate = new Date(summary.start_time);
    const endDate = new Date(summary.end_time);

    return (
        <div className='ai-summary-metadata'>
            <div className='metadata-row'>
                <div className='metadata-item'>
                    <span className='metadata-label'>
                        <FormattedMessage
                            id='ai.summary.metadata.messages'
                            defaultMessage='Messages'
                        />
                        {': '}
                    </span>
                    <span className='metadata-value'>{summary.message_count}</span>
                </div>
                <div className='metadata-item'>
                    <span className='metadata-label'>
                        <FormattedMessage
                            id='ai.summary.metadata.participants'
                            defaultMessage='Participants'
                        />
                        {': '}
                    </span>
                    <span className='metadata-value'>{summary.participants}</span>
                </div>
            </div>
            <div className='metadata-row'>
                <div className='metadata-item'>
                    <span className='metadata-label'>
                        <FormattedMessage
                            id='ai.summary.metadata.timeRange'
                            defaultMessage='Time Range'
                        />
                        {': '}
                    </span>
                    <span className='metadata-value'>
                        <FormattedDate value={startDate}/>{' '}
                        <FormattedTime value={startDate}/>{' - '}
                        <FormattedDate value={endDate}/>{' '}
                        <FormattedTime value={endDate}/>
                    </span>
                </div>
            </div>
            {summary.channel_name && (
                <div className='metadata-row'>
                    <div className='metadata-item'>
                        <span className='metadata-label'>
                            <FormattedMessage
                                id='ai.summary.metadata.channel'
                                defaultMessage='Channel'
                            />
                            {': '}
                        </span>
                        <span className='metadata-value'>{summary.channel_name}</span>
                    </div>
                </div>
            )}
        </div>
    );
};

export default SummaryMetadata;

