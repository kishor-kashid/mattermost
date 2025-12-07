import React from 'react';

import type {AILinkSummary} from 'types/ai';

import {KeyPoints} from './key_points';

type Props = {
    summary?: AILinkSummary;
    loading?: boolean;
    onRefresh?: () => void;
};

export const LinkSummaryCard = ({summary, loading, onRefresh}: Props) => {
    console.log('[LinkSummaryCard] Received props:', {summary, loading});

    // Show nothing if no summary and not loading
    if (!summary && !loading) {
        return null;
    }

    // When loading, show skeleton
    if (loading) {
        return (
            <div className='ai-link-summary__card'>
                <div className='ai-link-summary__header'>
                    <div className='ai-link-summary__title'>
                        {summary?.title || 'Loading summary...'}
                    </div>
                </div>
                <div className='ai-link-summary__body'>
                    <div className='ai-link-summary__skeleton'/>
                </div>
            </div>
        );
    }

    // Show full summary card
    return (
        <div className='ai-link-summary__card'>
            <div className='ai-link-summary__header'>
                <div className='ai-link-summary__title'>
                    {summary?.title || summary?.domain || 'Link summary'}
                </div>
                {summary?.domain && <div className='ai-link-summary__domain'>{summary.domain}</div>}
                {summary?.reading_time ? <div className='ai-link-summary__meta'>{summary.reading_time} min read</div> : null}
                {onRefresh && (
                    <button
                        type='button'
                        className='ai-link-summary__refresh'
                        onClick={onRefresh}
                    >
                        Refresh
                    </button>
                )}
            </div>
            <div className='ai-link-summary__body'>
                {summary?.summary && (
                    <div className='ai-link-summary__summary'>
                        {summary.summary}
                    </div>
                )}
                <KeyPoints points={summary?.key_points}/>
            </div>
        </div>
    );
};

