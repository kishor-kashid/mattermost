// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

interface Props {
    error: string | Error | null;
    onRetry?: () => void;
    onDismiss?: () => void;
}

export default function AIErrorDisplay({error, onRetry, onDismiss}: Props) {
    if (!error) {
        return null;
    }

    const errorMessage = typeof error === 'string' ? error : error.message;

    return (
        <div className='ai-error-display'>
            <div className='ai-error-icon'>
                <i className='icon icon-alert-outline'/>
            </div>
            <div className='ai-error-content'>
                <h4 className='ai-error-title'>{'Something went wrong'}</h4>
                <p className='ai-error-message'>{errorMessage}</p>
                <div className='ai-error-actions'>
                    {onRetry && (
                        <button
                            className='btn btn-primary'
                            onClick={onRetry}
                        >
                            {'Try Again'}
                        </button>
                    )}
                    {onDismiss && (
                        <button
                            className='btn btn-tertiary'
                            onClick={onDismiss}
                        >
                            {'Dismiss'}
                        </button>
                    )}
                </div>
            </div>
        </div>
    );
}

