// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

interface Props {
    variant?: 'default' | 'beta' | 'new';
    text?: string;
}

export default function AIFeatureBadge({variant = 'default', text}: Props) {
    const badgeText = text || {
        default: 'AI',
        beta: 'AI Beta',
        new: 'AI New',
    }[variant];

    const variantClass = `ai-feature-badge--${variant}`;

    return (
        <span className={`ai-feature-badge ${variantClass}`}>
            <i className='icon icon-robot'/>
            {badgeText}
        </span>
    );
}

