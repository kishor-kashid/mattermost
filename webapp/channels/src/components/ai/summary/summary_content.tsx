// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import Markdown from 'components/markdown';

import type {AISummary} from '@mattermost/types/ai';

interface Props {
    summary: AISummary;
}

const SummaryContent: React.FC<Props> = ({summary}) => {
    return (
        <div className='ai-summary-content'>
            <Markdown
                message={summary.summary}
                options={{
                    markdown: true,
                    singleline: false,
                    mentionHighlight: false,
                }}
            />
        </div>
    );
};

export default SummaryContent;

