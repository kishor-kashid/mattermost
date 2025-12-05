// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import LoadingSpinner from 'components/widgets/loading/loading_spinner';

interface Props {
    message?: string;
    size?: 'small' | 'medium' | 'large';
}

export default function AILoadingState({message = 'AI is processing...', size = 'medium'}: Props) {
    const sizeClass = {
        small: 'ai-loading-small',
        medium: 'ai-loading-medium',
        large: 'ai-loading-large',
    }[size];

    return (
        <div className={`ai-loading-state ${sizeClass}`}>
            <LoadingSpinner/>
            {message && (
                <p className='ai-loading-message'>{message}</p>
            )}
        </div>
    );
}

