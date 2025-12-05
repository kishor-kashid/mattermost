// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useState} from 'react';
import {FormattedMessage} from 'react-intl';
import {useDispatch} from 'react-redux';

import type {AISummary} from '@mattermost/types/ai';

import SummaryContent from './summary_content';
import SummaryMetadata from './summary_metadata';

import LoadingState from '../common/loading_state';
import ErrorDisplay from '../common/error_display';
import FeatureBadge from '../common/feature_badge';

interface Props {
    summary?: AISummary;
    loading?: boolean;
    error?: any;
    fromCache?: boolean;
    processingMs?: number;
    onClose?: () => void;
    onRegenerate?: () => void;
}

const SummaryPanel: React.FC<Props> = ({
    summary,
    loading,
    error,
    fromCache,
    processingMs,
    onClose,
    onRegenerate,
}) => {
    const dispatch = useDispatch();
    const [copied, setCopied] = useState(false);

    const handleCopy = () => {
        if (summary) {
            navigator.clipboard.writeText(summary.summary);
            setCopied(true);
            setTimeout(() => setCopied(false), 2000);
        }
    };

    const handleShare = () => {
        // TODO: Implement share functionality
        // Could post summary to channel as a message
    };

    if (loading) {
        return (
            <div className='ai-summary-panel'>
                <div className='panel-header'>
                    <h3>
                        <FormattedMessage
                            id='ai.summary.panel.title'
                            defaultMessage='AI Summary'
                        />
                    </h3>
                    <FeatureBadge feature='summarization'/>
                    {onClose && (
                        <button
                            className='close-btn'
                            onClick={onClose}
                            aria-label='Close'
                        >
                            <i className='icon icon-close'/>
                        </button>
                    )}
                </div>
                <div className='panel-content'>
                    <LoadingState
                        size='large'
                        message='Generating summary...'
                    />
                </div>
            </div>
        );
    }

    if (error) {
        return (
            <div className='ai-summary-panel'>
                <div className='panel-header'>
                    <h3>
                        <FormattedMessage
                            id='ai.summary.panel.title'
                            defaultMessage='AI Summary'
                        />
                    </h3>
                    <FeatureBadge feature='summarization'/>
                    {onClose && (
                        <button
                            className='close-btn'
                            onClick={onClose}
                            aria-label='Close'
                        >
                            <i className='icon icon-close'/>
                        </button>
                    )}
                </div>
                <div className='panel-content'>
                    <ErrorDisplay
                        error={error}
                        onRetry={onRegenerate}
                    />
                </div>
            </div>
        );
    }

    if (!summary) {
        return null;
    }

    return (
        <div className='ai-summary-panel'>
            <div className='panel-header'>
                <h3>
                    <FormattedMessage
                        id='ai.summary.panel.title'
                        defaultMessage='AI Summary'
                    />
                </h3>
                <FeatureBadge feature='summarization'/>
                {onClose && (
                    <button
                        className='close-btn'
                        onClick={onClose}
                        aria-label='Close'
                    >
                        <i className='icon icon-close'/>
                    </button>
                )}
            </div>

            {fromCache && (
                <div className='cache-indicator'>
                    <i className='icon icon-cached'/>
                    <FormattedMessage
                        id='ai.summary.fromCache'
                        defaultMessage='From cache'
                    />
                </div>
            )}

            <div className='panel-content'>
                <SummaryContent summary={summary}/>
                
                <div className='divider'/>
                
                <SummaryMetadata summary={summary}/>
                
                {processingMs !== undefined && (
                    <div className='processing-time'>
                        <FormattedMessage
                            id='ai.summary.processingTime'
                            defaultMessage='Processed in {time}ms'
                            values={{time: processingMs}}
                        />
                    </div>
                )}
            </div>

            <div className='panel-actions'>
                <button
                    className='btn btn-tertiary'
                    onClick={handleCopy}
                >
                    <i className={`icon ${copied ? 'icon-check' : 'icon-content-copy'}`}/>
                    <FormattedMessage
                        id={copied ? 'ai.summary.copied' : 'ai.summary.copy'}
                        defaultMessage={copied ? 'Copied!' : 'Copy'}
                    />
                </button>
                
                {onRegenerate && (
                    <button
                        className='btn btn-tertiary'
                        onClick={onRegenerate}
                    >
                        <i className='icon icon-refresh'/>
                        <FormattedMessage
                            id='ai.summary.regenerate'
                            defaultMessage='Regenerate'
                        />
                    </button>
                )}
            </div>
        </div>
    );
};

export default SummaryPanel;

